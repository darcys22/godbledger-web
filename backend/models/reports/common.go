package models

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

var log = logrus.WithField("prefix", "Reports")

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

func NewReport(req ReportsRequest) (error, *ReportResult) {
	switch req.Reports[0].Options.Title {
	case "TrialBalance":
		return TrialBalanceReport(req)
	case "GeneralLedger":
		return GeneralLedgerReport(req)
	default:
		log.Errorf("Unknown Report %s", req.Reports[0].Options.Title)
	}
	return fmt.Errorf("Unknown Report %s", req.Reports[0].Options.Title), nil
}
