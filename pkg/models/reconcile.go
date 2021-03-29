package models

import (
	"flag"
	"fmt"

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
type UnreconciledTransactionsRequest struct {
	Account string `json:"account"`
}

type UnreconciledTransactionLine struct {
	Description string `json:"description"`
}

type ReconcileResult struct {
	Account string                        `json:"account"`
	Result  []UnreconciledTransactionLine `json:"result"`
}

func UnreconciledTransactions(req UnreconciledTransactionsRequest) (error, *ReconcileResult) {
	set := flag.NewFlagSet("UnreconciledTransactionsRequest", 0)
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

	queryDB := `
	select
		s.description
	from
		splits as s
		join split_accounts as sa on s.split_id = sa.split_id
		join accounts as a on sa.account_id = a.account_id
	where
		a.name = ?
		and s.split_id not in (
			select
				distinct r.split_id
			from
				reconciliat ions as r
		)
	;`

	log.Debug("Querying Database")
	rows, err := ledger.LedgerDb.Query(queryDB, req.Account)
	if err != nil {
		return fmt.Errorf("Could not query database (%v)", err), nil
	}
	defer rows.Close()

	var r ReconcileResult
	r.Account = req.Account

	for rows.Next() {
		var utl UnreconciledTransactionLine
		if err := rows.Scan(&utl.Description); err != nil {
			return fmt.Errorf("Could not scan rows of query (%v)", err), &r
		}
		r.Result = append(r.Result, utl)
	}
	if rows.Err() != nil {
		return fmt.Errorf("rows errored with (%v)", rows.Err()), &r
	}

	return nil, &r
}
