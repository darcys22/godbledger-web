package models

import (
	"fmt"

	"github.com/darcys22/godbledger-web/backend/models/reports"
	"github.com/darcys22/godbledger-web/backend/models/reports/trialbalance"
	"github.com/darcys22/godbledger-web/backend/models/reports/generalledger"
)

func NewReport(req reports.ReportsRequest) (error, *reports.ReportResult) {
	switch req.Reports[0].Options.Title {
	case "TrialBalance":
		return trialbalance.TrialBalanceReport(req)
	case "GeneralLedger":
		return generalledger.GeneralLedgerReport(req)
	default:
		log.Errorf("Unknown Report %s", req.Reports[0].Options.Title)
	}
	return fmt.Errorf("Unknown Report %s", req.Reports[0].Options.Title), nil
}
