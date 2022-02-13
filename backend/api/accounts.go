package api

import (
	"net/http"

	m "github.com/darcys22/godbledger-web/backend/models"

	"github.com/gin-gonic/gin"
)

func GetAccounts(c *gin.Context) {
	accountsModel := m.NewAccountsListing()
	err := accountsModel.SearchAccounts()
	if err != nil {
		log.Errorf("Could not get journal listing (%v)", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
    return
	}
	c.JSON(http.StatusOK, accountsModel)
}

func PostAccount(c *gin.Context) {
	var account m.PostAccountCommand

	if err := c.BindJSON(&account); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := account.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
    return
	}
	c.JSON(http.StatusOK, account)
}

func DeleteAccount(c *gin.Context) {
	id := c.Params.ByName("id")
	if err := m.DeleteAccountCommand(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
    return
	}
	c.String(http.StatusOK, "Success")
}

func GetAccount(c *gin.Context) {
	id := c.Params.ByName("id")

	account, err := m.GetAccountCommand(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
    return
	}

	c.JSON(http.StatusOK, account)
}

func PostAccountTag(c *gin.Context) {
	//var account_tag m.PostAccountTagCommand

	//if err := c.BindJSON(&account_tag); err != nil {
		//c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		//return
	//}

	//if err := account_tag.Save(); err != nil {
		//c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	//}
	//c.JSON(http.StatusOK, account_tag)
}

func DeleteAccountTag(c *gin.Context) {
  account := c.Params.ByName("account")
  tag := c.Params.ByName("tag")
  if err := m.DeleteAccountTagCommand(account, tag); err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
    return
  }
  c.String(http.StatusOK, "Success")
}

type importAccounts struct{
    Name string `json:"name" binding:"required"`
}

func PostImportAccounts(c *gin.Context) {
  var postimportaccounts importAccounts
  c.BindJSON(&postimportaccounts)

	err := m.ImportAccountsCommand(postimportaccounts.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
    return
	}
	c.JSON(http.StatusOK, "Success")
}
