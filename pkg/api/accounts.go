package api

import (
	"flag"

	"github.com/darcys22/godbledger/godbledger/cmd"
	"github.com/darcys22/godbledger/godbledger/ledger"

	"github.com/gin-gonic/gin"
	"github.com/urfave/cli/v2"
)


type Results struct {
	Results []Account `json:"results"`
}

type Account struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}

func GetAccountListing(c *gin.Context) {

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

	arr := Results{}
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
