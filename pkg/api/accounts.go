package api

import (
	"flag"

	"github.com/darcys22/godbledger/godbledger/cmd"
	"github.com/darcys22/godbledger/godbledger/ledger"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var log = logrus.WithField("prefix", "AccountController")

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
		@rownum := @rownum +1 as idx, 
		accounts.name 
	FROM accounts 
		CROSS JOIN (SELECT @rownum := 0) r
	;`

	log.Debug("Querying Database")
	rows, err := ledger.LedgerDb.Query(queryDB)
	if err != nil {
		log.Errorf("Could not query database (%v)", err)
	}
	defer rows.Close()

	arr := Results{}
	arr.Results = []Account{}

	for rows.Next() {
		//Scan one account record
		var t Account
		if err := rows.Scan(&t.ID, &t.Text); err != nil {
			log.Errorf("Could not scan rows of query (%v)", err)
		}
		arr.Results = append(arr.Results, t)
	}
	if rows.Err() != nil {
		log.Errorf("rows errored with (%v)", rows.Err())
	}
	c.JSON(200, &arr)
}
