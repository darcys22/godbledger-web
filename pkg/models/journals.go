package models

import ()

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

//type PostJournalCommand struct {
//Date          string     `json:"date" binding:"Required"`
//Narration     string     `json:"narration" binding:"Required"`
//LineItemCount string     `json:"lineItemCount" binding:"Required"`
//LineItems     []LineItem `json:"lineItems" binding:"Required"`
//}
type PostJournalCommand struct {
	Date string `json:"date" binding:"Required"`
}

//
// QUERIES
//

//type GetDashboardQuery struct {
//Slug  string
//OrgId int64

//Result *Dashboard
//}
