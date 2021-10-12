package reports

import (
	"fmt"
	"time"
	"strings"

	"github.com/darcys22/godbledger/godbledger/ledger"

	"github.com/sirupsen/logrus"
)

var log = logrus.WithField("prefix", "Reports")

var DecimalsCache = map[string]int{
	"USD": 2,
}

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
	//Styling string   `json:"styling"`
	Row []string `json:"row"`
}
type ReportResult struct {
	Options Options      `json:"options"`
	Columns []string     `json:"columns"`
	Result  []ReportLine `json:"result"`
}

type ReportProcessor struct {
	Columns []string
	Input map[string]string
	Decimals int
}

func ProcessRows(ledger *ledger.Ledger, columns []string, inputs []string) (error, []string) {
	var rowProcessor = ReportProcessor{columns,map[string]string{},0}
	for i, column := range columns {

		log.Debug(column)
		switch column {
		case "Currency":
			if val, ok := DecimalsCache[inputs[i]]; ok {
				log.Debug("found currency ", inputs[i])
					rowProcessor.Decimals = val
			} else {
				log.Debug("not found currency ", inputs[i])
				querycurrency := "SELECT decimals FROM currencies where name = ?"
				rows, err := ledger.LedgerDb.Query(querycurrency, inputs[i])
				if err != nil {
					return fmt.Errorf("Could not query database (%v)", err), nil
				}
				defer rows.Close()
				for rows.Next() {
					if err := rows.Scan(&rowProcessor.Decimals); err != nil {
						return fmt.Errorf("Could not scan rows of query (%v)", err), nil
					}
					rowProcessor.Input[column] = inputs[i]
				}
				if rows.Err() != nil {
					return fmt.Errorf("rows errored with (%v)", rows.Err()), nil
				}
			}
		default:
			rowProcessor.Input[column] = inputs[i]
		}
	}

	var result []string

	for i, column := range columns {
		switch column {
		case "Amount":
			if rowProcessor.Decimals > 0 {
				atomicAmount := strings.TrimSpace(inputs[i])
				index := len(atomicAmount) - rowProcessor.Decimals
				decimalAmount := atomicAmount[:index] + "." + atomicAmount[index:]
				result = append(result, decimalAmount)
			} else {
				result = append(result, inputs[i])
			}

		default:
			result = append(result, inputs[i])
		}
	}

	return nil, result
}

func ProcessDate(time_str string) (time.Time, error) {
	layout := "2006-01-02"
	return time.Parse(layout, time_str)
}
