package models

import (
	"fmt"
	"context"
	"time"
	"math"

	"github.com/darcys22/godbledger-web/backend/models/backend"
	"github.com/darcys22/godbledger-web/backend/setting"
	"github.com/darcys22/godbledger-web/backend/models/feeds"
	"github.com/gin-gonic/gin"

	pb "github.com/darcys22/godbledger/proto/transaction"
	"google.golang.org/grpc"
)

type ExternalAccountsResult struct {
	Results []Account `json:"results"`
}

func GetExternalAccountListing(c *gin.Context) {
  db := backend.GetConnection()
	queryDB := `
		select
			a.name
		from
			accounts as a
			JOIN account_tag as at on at.account_id = a.account_id
			JOIN tags as t ON t.tag_id = at.tag_id
		WHERE
			t.tag_name = "external"
	;`

	log.Debug("Querying Database")
	rows, err := db.Query(queryDB)
	if err != nil {
		log.Errorf("Could not query database (%v)", err)
	}
	defer rows.Close()

	arr := ExternalAccountsResult{}
	arr.Results = []Account{}

	index := 0

	for rows.Next() {
		//Scan one account record
		index++
		t := Account{ID: index}
		if err := rows.Scan(&t.Text); err != nil {
			log.Errorf("Could not scan rows of query (%v)", err)
		}
		arr.Results = append(arr.Results, t)
	}
	if rows.Err() != nil {
		log.Errorf("rows errored with (%v)", rows.Err())
	}
	c.JSON(200, &arr)
}

// Get Unreconciled Transactions

type UnreconciledTransactionOptions struct {
	Account   string `json:"account"`
	StartDate string `json:"startdate"`
	EndDate   string `json:"enddate"`
}

type UnreconciledTransactionsRequest struct {
	Options UnreconciledTransactionOptions `json:"options"`
	Columns []string                       `json:"columns"`
}

type UnreconciledTransactionLine struct {
	Row []string `json:"row"`
}

type ReconcileResult struct {
	Options UnreconciledTransactionOptions `json:"options"`
	Columns []string                       `json:"columns"`
	Result  []UnreconciledTransactionLine  `json:"result"`
}

type UploadCSVOptions struct {
  Columns   []string           `json:"columns" binding:"required"`
	StartRow int                 `json:"startRow"`
	EndRow   int                 `json:"endRow"`
}

type UploadCSVRequest struct {
  Account string             `json:"account" binding:"required"`
  Options UploadCSVOptions   `json:"options" binding:"required"`
  Filename string            `json:"filename" binding:"required"`
  File string                `json:"file" binding:"required"`
}

type UploadCSVResult struct {
  Something string
}

func UnreconciledTransactions(req UnreconciledTransactionsRequest) (error, *ReconcileResult) {
  db := backend.GetConnection()
	queryDB := `
	select
		s.split_date,
		s.description,
		s.amount,
    a.currency,
    c.decimals
	from
		splits as s
		join split_accounts as sa on s.split_id = sa.split_id
		join accounts as a on sa.account_id = a.account_id
		join currencies as c on a.currency = c.name
	where
		a.name = ?
		and s.split_id not in (
			select
				distinct r.split_id
			from
				reconciliations as r
		)
	;`

	log.Info("Querying Database for unreconciled transactions on account ", req.Options.Account)
	rows, err := db.Query(queryDB, req.Options.Account)
	if err != nil {
    log.Info("Could not query database ", err)
		return fmt.Errorf("Could not query database (%v)", err), nil
	}
	defer rows.Close()

	var r ReconcileResult
	r.Options = req.Options
	r.Columns = req.Columns

  var date time.time
  description := ""
  amount := 0
  currency := ""
  decimals := 0
  amount_str := ""
	for rows.Next() {
		var utl UnreconciledTransactionLine
		if err := rows.Scan(&date, &description, &amount, &currency, &decimals); err != nil {
			return fmt.Errorf("Could not scan rows of query (%v)", err), &r
		}
    amount_str = fmt.Sprintf("%.2f", float64(amount)/math.Pow(10, float64(decimals)))
		utl.Row = append(utl.Row, date)
		utl.Row = append(utl.Row, description)
		utl.Row = append(utl.Row, amount_str)
		utl.Row = append(utl.Row, currency)
		r.Result = append(r.Result, utl)
	}

	if rows.Err() != nil {
		return fmt.Errorf("rows errored with (%v)", rows.Err()), &r
	}

	return nil, &r
}

func UploadCSV(req UploadCSVRequest) (error, *UploadCSVResult) {
	var r UploadCSVResult
  rows, err := feeds.ReadCSVBankStatement(req.File, req.Options.Columns, req.Options.StartRow, req.Options.EndRow)
	if err != nil {
    log.Info(err)
		return err, nil
	}
  log.Info(rows)
  cfg := setting.GetConfig()
	address := fmt.Sprintf("%s:%s", cfg.GoDBLedgerHost, cfg.GoDBLedgerPort)
	log.WithField("address", address).Info("GRPC Dialing on port")
	opts := []grpc.DialOption{}

	if cfg.GoDBLedgerCACert != "" && cfg.GoDBLedgerCert != "" && cfg.GoDBLedgerKey != "" {
		tlsCredentials, err := loadTLSCredentials(cfg)
		if err != nil {
			return fmt.Errorf("Could not load TLS credentials (%v)", err), &r
		}
		opts = append(opts, grpc.WithTransportCredentials(tlsCredentials))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		return fmt.Errorf("Could not connect to GRPC (%v)", err), &r
	}
	defer conn.Close()
	client := pb.NewTransactorClient(conn)

	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	transactionLines := make([]*pb.TransactionFeedLine, len(rows))

	for i, feedtransaction := range rows {
		//TODO sean get this from somewhere
		decimals := float64(currenciesDecimals["USD"])
		if err != nil {
			return fmt.Errorf("Could not process the amount as a float (%v)", err), &r
		}
		amount := int64(feedtransaction.Amount * math.Pow(10, decimals))
		transactionLines[i] = &pb.TransactionFeedLine{
      Date:        feedtransaction.Date.Format("2006-01-02"),
			Description: feedtransaction.Description,
			Hash:        feedtransaction.Hash,
			Amount:      amount,
		}
	}

	grpc_req := &pb.TransactionFeedRequest{
    Account:     req.Account,
		Lines:       transactionLines,
	}
	grpc_resp, err := client.AddTransactionFeed(ctxTimeout, grpc_req)
	if err != nil {
		return fmt.Errorf("Could not call Add Transaction FeedMethod (%v)", err), &r
	}
  r.Something = grpc_resp.GetMessage()
	log.Infof("Add Transaction Feed Response: %s", r.Something)
	return nil, &r
}
