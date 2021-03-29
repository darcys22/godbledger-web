package models

import (
	"flag"

	"github.com/darcys22/godbledger/godbledger/cmd"
	"github.com/darcys22/godbledger/godbledger/ledger"

	"github.com/gin-gonic/gin"
	"github.com/urfave/cli/v2"
)

// Get External Accounts

type Account struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}

type ExternalAccountsResult struct {
	Results []Account `json:"results"`
}

func GetExternalAccountListing(c *gin.Context) {

	set := flag.NewFlagSet("getExternalAccountListing", 0)
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
	rows, err := ledger.LedgerDb.Query(queryDB)
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

type ReconcileRequest struct {
	Account string `json:"account"`
}

type UnreconciledTransactionLine struct {
	Row string `json:"row"`
}

type ReconcileResult struct {
	Account string                        `json:"account"`
	Result  []UnreconciledTransactionLine `json:"result"`
}

func GetUnreconciledTransactions(req ReconcileRequest) (error, *ReconcileResult) {
	//set := flag.NewFlagSet("getJournalListing", 0)
	//set.String("config", "", "doc")

	//ctx := cli.NewContext(nil, set, nil)
	//err, cfg := cmd.MakeConfig(ctx)
	//if err != nil {
	//return fmt.Errorf("Could not make config (%v)", err), nil
	//}

	//ledger, err := ledger.New(ctx, cfg)
	//if err != nil {
	//return fmt.Errorf("Could not make new ledger (%v)", err), nil
	//}

	////queryDateStart := time.Now().Add(time.Hour * 24 * 365 * -100)
	//queryDateEnd := time.Now().Add(time.Hour * 24 * 365 * 100)
	//queryDB := `
	//SELECT split_accounts.account_id,
	//Sum(splits.amount),
	//currency.decimals
	//FROM   splits
	//JOIN split_accounts ON splits.split_id = split_accounts.split_id
	//JOIN currencies AS currency ON splits.currency = currency.NAME
	//WHERE  splits.split_date <= ?
	//AND "void" NOT IN (SELECT t.tag_name
	//FROM   tags AS t
	//JOIN transaction_tag AS tt
	//ON tt.tag_id = t.tag_id
	//WHERE  tt.transaction_id = splits.transaction_id)
	//GROUP  BY split_accounts.account_id, splits.currency
	//;`

	//log.Debug("Querying Database")
	////rows, err := ledger.LedgerDb.Query(queryDB, queryDateStart, queryDateEnd)
	//rows, err := ledger.LedgerDb.Query(queryDB, queryDateEnd)

	//var r ReportResult
	//r.Options = req.Reports[0].Options
	//r.Columns = req.Reports[0].Columns

	//if err != nil {
	//return fmt.Errorf("Could not query database (%v)", err), nil
	//}
	//defer rows.Close()

	//for rows.Next() {
	//t := make([]string, len(req.Reports[0].Columns))
	//pointers := make([]interface{}, len(t))
	//for i, _ := range pointers {
	//pointers[i] = &t[i]
	//}
	//if err := rows.Scan(pointers...); err != nil {
	//return fmt.Errorf("Could not scan rows of query (%v)", err), nil
	//}
	//var l ReportLine
	//l.Row = t
	//r.Result = append(r.Result, l)
	//}
	//if rows.Err() != nil {
	//return fmt.Errorf("rows errored with (%v)", rows.Err()), nil
	//}

	//return nil, &r
	return nil, nil
}
