package generalledger

import (
	//"context"
	//"crypto/tls"
	//"crypto/x509"
	"flag"
	"fmt"
	"strings"
	//"io/ioutil"
	"time"

	"github.com/darcys22/godbledger-web/backend/models/reports"
	"github.com/darcys22/godbledger/godbledger/cmd"
	"github.com/darcys22/godbledger/godbledger/ledger"

	"github.com/urfave/cli/v2"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithField("prefix", "Reports-General-Ledger")

// General Ledger
var GeneralLedgerColumns = map[string]string{
	"ID":           "transactions.transaction_id",
	"Date":         "splits.split_date",
	"Description":  "splits.description",
	"Currency":     "splits.currency",
	"Decimals":     "currency.decimals",
	"Amount":       "splits.amount",
	"AtomicAmount": "splits.amount",
	"Account":      "split_accounts.account_id",
}

func GeneralLedgerReport(req reports.ReportsRequest) (error, *reports.ReportResult) {
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

	queryDateStart := time.Now().Add(time.Hour * 24 * 365 * -100)
	queryDateEnd := time.Now().Add(time.Hour * 24 * 365 * 100)

	queryDB := strings.Builder{}
	queryDB.WriteString("SELECT\n")

	for i, col := range req.Reports[0].Columns {
		queryDB.WriteString(GeneralLedgerColumns[col])
		if i != len(req.Reports[0].Columns)-1 {
			queryDB.WriteString(",")
		}
		queryDB.WriteString("\n")
	}

	queryDB.WriteString(`FROM splits
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
					tt.transaction_id = splits.transaction_id)
			AND "main" IN (
				SELECT
					t.tag_name
				FROM
					tags AS t
					JOIN account_tag AS at ON at.tag_id = t.tag_id
				WHERE
					at.account_id = split_accounts.account_id)
	;`)

	log.Debug("Querying Database")
	rows, err := ledger.LedgerDb.Query(queryDB.String(), queryDateStart, queryDateEnd)

	if err != nil {
		return fmt.Errorf("Could not query database (%v)", err), nil
	}
	defer rows.Close()

	var r reports.ReportResult
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
		err, processedRow := reports.ProcessRows(ledger, req.Reports[0].Columns, t)
		if err != nil {
			return fmt.Errorf("Could not process rows of query (%v)", err), nil
		}
		r.Result = append(r.Result, reports.ReportLine{processedRow})
	}
	if rows.Err() != nil {
		return fmt.Errorf("rows errored with (%v)", rows.Err()), nil
	}

	return nil, &r
}
