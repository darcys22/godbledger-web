package models

import (
	"errors"
	"time"
)

// Typed errors
var (
	ErrDashboardNotFound = errors.New("Account not found")
)

// Dashboard model
type Dashboard struct {
	Id      int64
	Slug    string
	OrgId   int64
	Version int

	Created time.Time
	Updated time.Time

	Title string
	Data  map[string]interface{}
}

// NewDashboard creates a new dashboard
func NewDashboard(title string) *Dashboard {
	dash := &Dashboard{}
	dash.Data = make(map[string]interface{})
	dash.Data["title"] = title
	dash.Title = title
	return dash
}

// GetTags turns the tags in data json into go string array
func (dash *Dashboard) GetTags() []string {
	jsonTags := dash.Data["tags"]
	if jsonTags == nil {
		return []string{}
	}

	arr := jsonTags.([]interface{})
	b := make([]string, len(arr))
	for i := range arr {
		b[i] = arr[i].(string)
	}
	return b
}

//
// COMMANDS
//

//class LineItem {
//constructor() {
////this.date = new Date();
//this._date = "";
//this._description = "";
//this._account = "";
//this._amount = 0;
//}

//class Journal {
//constructor() {
//this.date = new Date();
//this.narration = "Display Me";
//this.lineitems = [];
//this._lineItemCount = 0;

type LineItem struct {
	Date        string `json:"date" binding:"Required"`
	Description string `json:"description" binding:"Required"`
	Account     string `json:"Account" binding:"Required"`
	Amount      string `json:"Amount" binding:"Required"`
}

type PostJournalCommand struct {
	Date          string     `json:"date" binding:"Required"`
	Narration     string     `json:"narration" binding:"Required"`
	LineItemCount string     `json:"lineItemCount" binding:"Required"`
	LineItems     []LineItem `json:"lineItems" binding:"Required"`
}

//
// QUERIES
//

type GetDashboardQuery struct {
	Slug  string
	OrgId int64

	Result *Dashboard
}
