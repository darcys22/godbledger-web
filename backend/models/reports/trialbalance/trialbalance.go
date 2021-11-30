package trialbalance

import (
	"fmt"
	"strings"

	"github.com/darcys22/godbledger-web/backend/models/reports"
	"github.com/darcys22/godbledger-web/backend/models/backend"

	"github.com/sirupsen/logrus"
)

var log = logrus.WithField("prefix", "Reports-Trial-Balance")

// Trial Balance
var trialBalanceColumns = map[string]string{
	"Accountname":  "split_accounts.account_id",
	"Amount":       "Sum(splits.amount)",
	"AtomicAmount": "Sum(splits.amount)",
	"Currency":     "currency.name",
}

func TrialBalanceReport(req reports.ReportsRequest) (error, *reports.ReportResult) {
  db := backend.GetConnection()
	queryDateStart, err := reports.ProcessDate(req.Reports[0].Options.StartDate)
	if err != nil {
		return fmt.Errorf("Could not process start date (%v)", err), nil
	}
	queryDateEnd, err := reports.ProcessDate(req.Reports[0].Options.EndDate)
	if err != nil {
		return fmt.Errorf("Could not process start date (%v)", err), nil
	}

	queryDB := strings.Builder{}
	queryDB.WriteString("SELECT\n")

	for i, col := range req.Reports[0].Columns {
		queryDB.WriteString(trialBalanceColumns[col])
		if i != len(req.Reports[0].Columns)-1 {
			queryDB.WriteString(",")
		}
		queryDB.WriteString("\n")
	}

	queryDB.WriteString(`FROM splits
	JOIN split_accounts ON splits.split_id = split_accounts.split_id
	JOIN currencies AS currency ON splits.currency = currency.NAME
WHERE  splits.split_date >= ?
	AND splits.split_date <= ?
	AND "void" NOT IN (SELECT t.tag_name
		FROM   tags AS t
			JOIN transaction_tag AS tt
				ON tt.tag_id = t.tag_id
		WHERE  tt.transaction_id = splits.transaction_id)
	AND "main" IN (SELECT t.tag_name
		FROM   tags AS t
			JOIN account_tag AS at
				ON at.tag_id = t.tag_id
		WHERE  at.account_id = split_accounts.account_id)
GROUP  BY split_accounts.account_id, splits.currency
			;`)

	log.Debug("Querying Database")
	log.Trace(queryDB.String())
	rows, err := db.Query(queryDB.String(), queryDateStart, queryDateEnd)

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
		err, processedRow := reports.ProcessRows(db, req.Reports[0].Columns, t)
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
