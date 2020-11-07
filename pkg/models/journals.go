package models

import (
	"flag"
	"fmt"
	//"math"
	//"strconv"
	"time"

	"github.com/darcys22/godbledger/godbledger/cmd"
	"github.com/darcys22/godbledger/godbledger/ledger"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var log = logrus.WithField("prefix", "JournalsModel")

type LineItem struct {
	ID          string `json:"id"`
	Date        string `json:"_date" binding:"required"`
	Description string `json:"_description"`
	Account     string `json:"_account" binding:"required"`
	Amount      int    `json:"_amount" binding:"required"`
	Currency    string `json:"_currency" binding:"required"`
}

type PostJournalCommand struct {
	Date          string     `json:"_date" binding:"required"`
	Narration     string     `json:"narration"`
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
			splits.split_date,
			splits.description,
			splits.currency,
			currency.decimals,
			splits.amount,
			split_accounts.account_id
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
		LIMIT 50
	;`

	log.Debug("Querying Database")
	rows, err := ledger.LedgerDb.Query(queryDB, queryDateStart, queryDateEnd)

	if err != nil {
		return fmt.Errorf("Could not query database (%v)", err)
	}
	defer rows.Close()

	for rows.Next() {
		// Scan one customer record
		var t LineItem
		var decimals float64
		if err := rows.Scan(&t.ID, &t.Date, &t.Description, &t.Currency, &decimals, &t.Amount, &t.Account); err != nil {
			return fmt.Errorf("Could not scan rows of query (%v)", err)
		}
		//centsAmount := float64(t.Amount)
		//if err != nil {
		//return fmt.Errorf("Could not process the amount as a float (%v)", err)
		//}
		//t.Amount = fmt.Sprintf("%.2f", centsAmount/math.Pow(10, decimals))
		j.Journals = append(j.Journals, t)
	}
	if rows.Err() != nil {
		return fmt.Errorf("rows errored with (%v)", rows.Err())
	}

	return nil
}

func (j *PostJournalCommand) Save(transaction) error {
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	transactionLines := make([]*pb.LineItem, transaction.lineitemcount)

	for i, accChange := range transaction.lineitems {
		amountInt64 := accChange.Balance.Num().Int64() * int64(100) / accChange.Balance.Denom().Int64()
		transactionLines[i] = &pb.LineItem{
			Accountname: accChange.Name,
			Description: accChange.Description,
			Amount:      amountInt64,
			Currency:    accChange.Currency,
		}
	}

	req := &pb.TransactionRequest{
		Date:        t.Date.Format("2006-01-02"),
		Description: t.Payee,
		Lines:       transactionLines,
	}
	r, err := client.AddTransaction(ctx, req)
	if err != nil {
		return fmt.Errorf("Could not call Add Transaction Method (%v)", err)
	}
	log.Infof("Add Transaction Response: %s", r.GetMessage())
	return nil
}
