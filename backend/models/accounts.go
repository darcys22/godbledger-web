package models

import (
	"context"
	//"crypto/tls"
	//"crypto/x509"
	"flag"
	"fmt"
	//"io/ioutil"
	//"math"
	//"strconv"
	"time"

	"github.com/darcys22/godbledger/godbledger/cmd"
	"github.com/darcys22/godbledger/godbledger/ledger"
	pb "github.com/darcys22/godbledger/proto/transaction"

	"google.golang.org/grpc"
	//"google.golang.org/grpc/credentials"

	//"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

type Account struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}

type PostAccountCommand struct {
	Name string     `json:"name" binding:"required"`
	Tags []string `json:"tags"`
}

type GetAccounts struct {
	Results []Account `json:"results"`
}

func NewAccountsListing() *GetAccounts {
	accounts := []Account{}
	return &GetAccounts{accounts}
}

func (a *GetAccounts) SearchAccounts() error {
	set := flag.NewFlagSet("getAccountListing", 0)
	set.String("config", "", "doc")

	ctx := cli.NewContext(nil, set, nil)
	err, cfg := cmd.MakeConfig(ctx)
	if err != nil {
		log.Errorf("Could not make config (%v)", err)
	}

	ledger, err := ledger.New(ctx, cfg)
	if err != nil {
		log.Errorf("Could not make new ledger (%v)", err)
	}

	queryDB := `
	SELECT 
		accounts.name 
	FROM accounts 
	;`

	log.Debug("Querying Database")
	rows, err := ledger.LedgerDb.Query(queryDB)
	if err != nil {
		log.Errorf("Could not query database (%v)", err)
	}
	defer rows.Close()

	a.Results = []Account{}

	index := 0

	for rows.Next() {
		//Scan one account record
		index++
		t := Account{ID: index}
		if err := rows.Scan(&t.Text); err != nil {
			log.Errorf("Could not scan rows of query (%v)", err)
		}
		a.Results = append(a.Results, t)
	}
	if rows.Err() != nil {
		log.Errorf("rows errored with (%v)", rows.Err())
	}
	return nil
}

func (j *PostAccountCommand) Save() error {
	log.Trace("Calling Save Account function")
	set := flag.NewFlagSet("PostAccount", 0)
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

	req := &pb.AccountTagRequest{
		Account:  j.Name,
		Tag:      j.Tags,
	}
	r, err := client.AddAccount(ctxTimeout, req)
	if err != nil {
		return fmt.Errorf("Could not call Add Account Method (%v)", err)
	}
	log.Infof("Add Account Response: %s", r.GetMessage())
	return nil
}

func DeleteAccountCommand(id string) error {
	log.Trace("Calling Delete Account function")
	set := flag.NewFlagSet("DeleteAccount", 0)
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

	req := &pb.DeleteAccountTagRequest{
		Account: id,
	}
	r, err := client.DeleteAccount(ctxTimeout, req)
	if err != nil {
		return fmt.Errorf("Could not call Delete Transaction Method (%v)", err)
	}
	log.Infof("Delete Transaction Response: %s", r.GetMessage())

	return nil
}

func GetAccountCommand(id string) (PostAccountCommand, error) {
	log.Trace("Calling Get Journal Command function")
	j := PostAccountCommand{}
	j.Tags = []string{}

	//set := flag.NewFlagSet("getJournal", 0)
	//set.String("config", "", "doc")

	//ctx := cli.NewContext(nil, set, nil)
	//err, cfg := cmd.MakeConfig(ctx)
	//if err != nil {
		//return j, fmt.Errorf("Could not make config (%v)", err)
	//}

	//ledger, err := ledger.New(ctx, cfg)
	//if err != nil {
		//return j, fmt.Errorf("Could not make new ledger (%v)", err)
	//}

	//queryDB := `
		//SELECT
			//transactions.transaction_id,
			//splits.split_date,
			//splits.description,
			//splits.currency,
			//currency.decimals,
			//splits.amount,
			//split_accounts.account_id,
			//transactions.description
		//FROM
			//splits
			//JOIN split_accounts ON splits.split_id = split_accounts.split_id
			//JOIN transactions on splits.transaction_id = transactions.transaction_id
			//JOIN currencies AS currency ON splits.currency = currency.NAME
		//WHERE
			//transactions.transaction_id = ?
		//LIMIT 50
	//;`

	//log.Debug("Querying Database")
	//rows, err := ledger.LedgerDb.Query(queryDB, id)

	//if err != nil {
		//return j, fmt.Errorf("Could not query database (%v)", err)
	//}
	//defer rows.Close()

	//for rows.Next() {
		//var t LineItem
		//var decimals float64
		//var narration string
		//if err := rows.Scan(&t.ID, &t.Date, &t.Description, &t.Currency, &decimals, &t.Amount, &t.Account, &narration); err != nil {
			//return j, fmt.Errorf("Could not scan rows of query (%v)", err)
		//}
		//centsAmount, err := strconv.ParseFloat(t.Amount, 64)
		//if err != nil {
			//return j, fmt.Errorf("Could not process the amount as a float (%v)", err)
		//}
		//t.Amount = fmt.Sprintf("%.2f", centsAmount/math.Pow(10, decimals))
		//j.LineItems = append(j.LineItems, t)
		//j.Narration = narration
		//j.Date = t.Date
	//}
	//if rows.Err() != nil {
		//return j, fmt.Errorf("rows errored with (%v)", rows.Err())
	//}

	//j.LineItemCount = len(j.LineItems)
	return j, nil
}
