package models

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"strconv"
	"time"

	"github.com/darcys22/godbledger/godbledger/cmd"
	"github.com/darcys22/godbledger/godbledger/ledger"
	pb "github.com/darcys22/godbledger/proto/transaction"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var log = logrus.WithField("prefix", "Model")

type LineItem struct {
	ID          string `json:"id"`
	Date        string `json:"_date" binding:"required"`
	Description string `json:"_description"`
	Account     string `json:"_account" binding:"required"`
	Amount      string `json:"_amount" binding:"required"`
	Currency    string `json:"_currency" binding:"required"`
}

type PostJournalCommand struct {
	Date          string     `json:"_date" binding:"required"`
	Narration     string     `json:"_narration"`
	LineItemCount int        `json:"_lineItemCount" binding:"required"`
	LineItems     []LineItem `json:"_lineItems" binding:"required"`
}

type GetJournals struct {
	Journals []LineItem
}

func NewJournalsListing() *GetJournals {
	lineItems := []LineItem{}
	return &GetJournals{lineItems}
}

func (j *GetJournals) SearchJournals() error {
	log.Trace("Calling Search Journals function")
	j.Journals = []LineItem{}

	set := flag.NewFlagSet("getJournalListing", 0)
	set.String("config", "", "doc")

	ctx := cli.NewContext(nil, set, nil)
	err, cfg := cmd.MakeConfig(ctx)
	if err != nil {
		return fmt.Errorf("Could not make config (%v)", err)
	}

	ledger, err := ledger.New(ctx, cfg)
	if err != nil {
		return fmt.Errorf("Could not make new ledger (%v)", err)
	}

	queryDateStart := time.Now().Add(time.Hour * 24 * 365 * -100)
	queryDateEnd := time.Now().Add(time.Hour * 24 * 365 * 100)

	queryDB := `
		SELECT
			transactions.transaction_id,
			max(splits.split_date) as jdate,
			transactions.brief,
			sum(case when splits.amount > 0 then splits.amount else 0 end),
			currency.decimals,
			currency.NAME
		FROM
			splits
			JOIN split_accounts ON splits.split_id = split_accounts.split_id
			JOIN transactions on splits.transaction_id = transactions.transaction_id
			JOIN currencies AS currency ON splits.currency = currency.NAME
		WHERE
			splits.split_date >= ?
			AND splits.split_date <= ?
			AND "void" NOT IN(
				SELECT
					t.tag_name
				FROM
					tags AS t
					JOIN transaction_tag AS tt ON tt.tag_id = t.tag_id
				WHERE
					tt.transaction_id = splits.transaction_id
			)
		GROUP BY transactions.transaction_id, splits.currency
		ORDER BY jdate DESC
		LIMIT 50
	;`

	log.Debug("Querying Database")
	rows, err := ledger.LedgerDb.Query(queryDB, queryDateStart, queryDateEnd)

	if err != nil {
		return fmt.Errorf("Could not query database (%v)", err)
	}
	defer rows.Close()

	for rows.Next() {
		var t LineItem
		var decimals float64
		if err := rows.Scan(&t.ID, &t.Date, &t.Description, &t.Amount, &decimals, &t.Currency); err != nil {
			return fmt.Errorf("Could not scan rows of query (%v)", err)
		}
		centsAmount, err := strconv.ParseFloat(t.Amount, 64)
		if err != nil {
			return fmt.Errorf("Could not process the amount as a float (%v)", err)
		}
		t.Amount = fmt.Sprintf("%.2f", centsAmount/math.Pow(10, decimals))
		j.Journals = append(j.Journals, t)
	}
	if rows.Err() != nil {
		return fmt.Errorf("rows errored with (%v)", rows.Err())
	}

	return nil
}

var currenciesDecimals = map[string]int{
	"USD": 2,
}

func (j *PostJournalCommand) Save() error {
	log.Trace("Calling Save Journal function")
	set := flag.NewFlagSet("PostJournal", 0)
	set.String("config", "", "doc")

	ctx := cli.NewContext(nil, set, nil)
	err, cfg := cmd.MakeConfig(ctx)
	if err != nil {
		return fmt.Errorf("Could not make config (%v)", err)
	}

	address := fmt.Sprintf("%s:%s", cfg.Host, cfg.RPCPort)
	log.WithField("address", address).Info("GRPC Dialing on port")
	opts := []grpc.DialOption{}

	if cfg.CACert != "" && cfg.Cert != "" && cfg.Key != "" {
		tlsCredentials, err := loadTLSCredentials(cfg)
		if err != nil {
			return fmt.Errorf("Could not load TLS credentials (%v)", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(tlsCredentials))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		return fmt.Errorf("Could not connect to GRPC (%v)", err)
	}
	defer conn.Close()
	client := pb.NewTransactorClient(conn)

	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	transactionLines := make([]*pb.LineItem, j.LineItemCount)

	layout := "2006-01-02T15:04:05-07:00"
	t, err := time.Parse(layout, j.Date)
	if err != nil {
		return fmt.Errorf("Could not parse date", err)
	}

	for i, accChange := range j.LineItems {

		//TODO sean get this from somewhere
		decimals := float64(currenciesDecimals["USD"])
		dollarsAmount, err := strconv.ParseFloat(accChange.Amount, 64)
		if err != nil {
			return fmt.Errorf("Could not process the amount as a float (%v)", err)
		}
		amount := int64(dollarsAmount * math.Pow(10, decimals))
		transactionLines[i] = &pb.LineItem{
			Accountname: accChange.Account,
			Description: accChange.Description,
			Amount:      amount,
			Currency:    accChange.Currency,
		}
	}

	req := &pb.TransactionRequest{
		Date:        t.Format("2006-01-02"),
		Description: j.Narration,
		Lines:       transactionLines,
	}
	r, err := client.AddTransaction(ctxTimeout, req)
	if err != nil {
		return fmt.Errorf("Could not call Add Transaction Method (%v)", err)
	}
	log.Infof("Add Transaction Response: %s", r.GetMessage())
	return nil
}

func DeleteJournalCommand(id string) error {
	log.Trace("Calling Delete Journal function")
	set := flag.NewFlagSet("DeleteJournal", 0)
	set.String("config", "", "doc")

	ctx := cli.NewContext(nil, set, nil)
	err, cfg := cmd.MakeConfig(ctx)
	if err != nil {
		return fmt.Errorf("Could not make config (%v)", err)
	}

	address := fmt.Sprintf("%s:%s", cfg.Host, cfg.RPCPort)
	log.WithField("address", address).Info("GRPC Dialing on port")
	opts := []grpc.DialOption{}

	if cfg.CACert != "" && cfg.Cert != "" && cfg.Key != "" {
		tlsCredentials, err := loadTLSCredentials(cfg)
		if err != nil {
			return fmt.Errorf("Could not load TLS credentials (%v)", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(tlsCredentials))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		return fmt.Errorf("Could not connect to GRPC (%v)", err)
	}
	defer conn.Close()
	client := pb.NewTransactorClient(conn)

	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.DeleteRequest{
		Identifier: id,
	}
	r, err := client.DeleteTransaction(ctxTimeout, req)
	if err != nil {
		return fmt.Errorf("Could not call Delete Transaction Method (%v)", err)
	}
	log.Infof("Delete Transaction Response: %s", r.GetMessage())

	return nil
}

func GetJournalCommand(id string) (PostJournalCommand, error) {
	log.Trace("Calling Get Journal Command function")
	j := PostJournalCommand{}
	j.LineItems = []LineItem{}

	set := flag.NewFlagSet("getJournal", 0)
	set.String("config", "", "doc")

	ctx := cli.NewContext(nil, set, nil)
	err, cfg := cmd.MakeConfig(ctx)
	if err != nil {
		return j, fmt.Errorf("Could not make config (%v)", err)
	}

	ledger, err := ledger.New(ctx, cfg)
	if err != nil {
		return j, fmt.Errorf("Could not make new ledger (%v)", err)
	}

	queryDB := `
		SELECT
			transactions.transaction_id,
			splits.split_date,
			splits.description,
			splits.currency,
			currency.decimals,
			splits.amount,
			split_accounts.account_id,
			transactions.brief
		FROM
			splits
			JOIN split_accounts ON splits.split_id = split_accounts.split_id
			JOIN transactions on splits.transaction_id = transactions.transaction_id
			JOIN currencies AS currency ON splits.currency = currency.NAME
		WHERE
			transactions.transaction_id = ?
		LIMIT 50
	;`

	log.Debug("Querying Database")
	rows, err := ledger.LedgerDb.Query(queryDB, id)

	if err != nil {
		return j, fmt.Errorf("Could not query database (%v)", err)
	}
	defer rows.Close()

	for rows.Next() {
		var t LineItem
		var decimals float64
		var narration string
		if err := rows.Scan(&t.ID, &t.Date, &t.Description, &t.Currency, &decimals, &t.Amount, &t.Account, &narration); err != nil {
			return j, fmt.Errorf("Could not scan rows of query (%v)", err)
		}
		centsAmount, err := strconv.ParseFloat(t.Amount, 64)
		if err != nil {
			return j, fmt.Errorf("Could not process the amount as a float (%v)", err)
		}
		t.Amount = fmt.Sprintf("%.2f", centsAmount/math.Pow(10, decimals))
		j.LineItems = append(j.LineItems, t)
		j.Narration = narration
		j.Date = t.Date
	}
	if rows.Err() != nil {
		return j, fmt.Errorf("rows errored with (%v)", rows.Err())
	}

	j.LineItemCount = len(j.LineItems)
	return j, nil
}

func loadTLSCredentials(cfg *cmd.LedgerConfig) (credentials.TransportCredentials, error) {
	// Load certificate of the CA who signed server's certificate
	pemServerCA, err := ioutil.ReadFile(cfg.CACert)
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemServerCA) {
		return nil, fmt.Errorf("failed to add server CA's certificate")
	}

	// Load client's certificate and private key
	clientCert, err := tls.LoadX509KeyPair(cfg.Cert, cfg.Key)
	if err != nil {
		return nil, err
	}

	// Create the credentials and return it
	config := &tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      certPool,
	}

	return credentials.NewTLS(config), nil
}
