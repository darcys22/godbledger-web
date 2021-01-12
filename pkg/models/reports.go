package models

import (
	//"context"
	//"crypto/tls"
	//"crypto/x509"
	"flag"
	"fmt"
	//"io/ioutil"
	"time"

	"github.com/darcys22/godbledger/godbledger/cmd"
	"github.com/darcys22/godbledger/godbledger/ledger"

	"github.com/urfave/cli/v2"
)

type Options struct {
	Title     string `json:"title"`
	StartDate string `json:"startdate"`
	EndDate   string `json:"enddate"`
}

type ReportRequest struct {
	Options Options  `json:"options"`
	Columns []string `json:"columns"`
}

type ReportsRequest struct {
	Reports []ReportRequest `json:"reports"`
}

type ReportLine struct {
	Styling string   `json:"styling"`
	Row     []string `json:"row"`
}
type ReportResult struct {
	Options Options      `json:"options"`
	Columns []string     `json:"columns"`
	Result  []ReportLine `json:"result"`
}

var trialBalanceColumns = map[string]string{
	"AccountName": "split_accounts.account_id",
	"Amount":      "Sum(splits.amount)",
	"Currency":    "currency.decimals",
}

func NewReport(req ReportsRequest) (error, *ReportResult) {
	set := flag.NewFlagSet("getJournalListing", 0)
	set.String("config", "", "doc")

	ctx := cli.NewContext(nil, set, nil)
	err, cfg := cmd.MakeConfig(ctx)
	if err != nil {
		return fmt.Errorf("Could not make config (%v)", err), nil
	}

	ledger, err := ledger.New(ctx, cfg)
	if err != nil {
		return fmt.Errorf("Could not make new ledger (%v)", err), nil
	}

	//queryDateStart := time.Now().Add(time.Hour * 24 * 365 * -100)
	queryDateEnd := time.Now().Add(time.Hour * 24 * 365 * 100)
	queryDB := `
			SELECT split_accounts.account_id,
						 Sum(splits.amount),
						 currency.decimals
			FROM   splits
						 JOIN split_accounts ON splits.split_id = split_accounts.split_id
						 JOIN currencies AS currency ON splits.currency = currency.NAME
			WHERE  splits.split_date <= ?
						 AND "void" NOT IN (SELECT t.tag_name
																FROM   tags AS t
																			 JOIN transaction_tag AS tt
																				 ON tt.tag_id = t.tag_id
																WHERE  tt.transaction_id = splits.transaction_id)
			GROUP  BY split_accounts.account_id, splits.currency
			;`

	log.Debug("Querying Database")
	//rows, err := ledger.LedgerDb.Query(queryDB, queryDateStart, queryDateEnd)
	rows, err := ledger.LedgerDb.Query(queryDB, queryDateEnd)

	var r ReportResult
	r.Options = req.Reports[0].Options
	r.Columns = req.Reports[0].Columns

	if err != nil {
		return fmt.Errorf("Could not query database (%v)", err), nil
	}
	defer rows.Close()

	for rows.Next() {
		t := make([]string, len(req.Reports[0].Columns))
		pointers := make([]interface{}, len(t))
		for i, _ := range pointers {
			pointers[i] = &t[i]
		}
		if err := rows.Scan(pointers...); err != nil {
			return fmt.Errorf("Could not scan rows of query (%v)", err), nil
		}
		var l ReportLine
		l.Row = t
		r.Result = append(r.Result, l)
	}
	if rows.Err() != nil {
		return fmt.Errorf("rows errored with (%v)", rows.Err()), nil
	}

	return nil, &r
}
