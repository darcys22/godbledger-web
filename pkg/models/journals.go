package models

import ()

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
	Date        string `json:"_date" binding:"required"`
	Description string `json:"_description"`
	Account     string `json:"_account" binding:"required"`
	Amount      int    `json:"_amount" binding:"required"`
}

type PostJournalCommand struct {
	Date          string     `json:"date" binding:"required"`
	Narration     string     `json:"narration"`
	LineItemCount int        `json:"_lineItemCount" binding:"required"`
	LineItems     []LineItem `json:"lineItems" binding:"required"`
}
