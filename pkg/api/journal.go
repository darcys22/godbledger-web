package api

import (
	m "github.com/darcys22/godbledger-web/pkg/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

//const journallisting = `[{"account":"Expenses:Groceries","id":"bt0a8pn7f64relhj2l00","date":"2011-03-15T00:00:00Z","desc":"Whole Food Market",     "amount":"7500","currency":"USD"},{"account":"Assets:Checking","id":"bt0a8pn7f64relhj2l00","date":"2011-03-15T00:00:00Z",  "desc":"Whole Food Market","amount":"-7500","currency":"USD"},{"account":"Expenses:Groceries","id":"bt0a8un7f64rb8lrumt0", "date":"2011-03-15T00:00:00Z","desc":"Whole Food Market","amount":"7500","currency":"USD"},{"account":"Assets:Checking",   "id":"bt0a8un7f64rb8lrumt0","date":"2011-03-15T00:00:00Z","desc":"Whole Food Market","amount":"-7500","currency":"USD"},   {"account":"Expenses:Groceries","id":"bt0a9b77f64r8fsahvag","date":"2011-03-15T00:00:00Z","desc":"Whole Food Market",      "amount":"7500","currency":"USD"},{"account":"Assets:Checking","id":"bt0a9b77f64r8fsahvag","date":"2011-03-15T00:00:00Z",  "desc":"Whole Food Market","amount":"-7500","currency":"USD"},{"account":"1","id":"bt14fb77f64ta6jaa50g","date":"2020-08-  23T10:40:44.704396691Z","desc":"Cash Income","amount":"10","currency":"AUD"},{"account":"2","id":"bt14fb77f64ta6jaa50g",   "date":"2020-08-23T10:40:44.704399619Z","desc":"Cash Income","amount":"-10","currency":"AUD"},{"account":"1","id":         "bt14n4v7f64u1b0k67gg","date":"2020-08-23T10:57:23.272618803Z","desc":"Cash Income","amount":"10","currency":"AUD"},       {"account":"2","id":"bt14n4v7f64u1b0k67gg","date":"2020-08-23T10:57:23.272620034Z","desc":"Cash Income","amount":"-10",    "currency":"AUD"}]`

func GetJournals(c *gin.Context) {
	journalsModel := m.NewJournalsListing()
	err := journalsModel.SearchJournals()
	if err != nil {
		log.Errorf("Could not get journal listing (%v)", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(200, journalsModel)
}

func PostJournal(c *gin.Context) {
	var journal m.PostJournalCommand

	if err := c.BindJSON(&journal); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, journal)
}
