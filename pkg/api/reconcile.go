package api

import (
	"net/http"

	m "github.com/darcys22/godbledger-web/pkg/models"
	"github.com/gin-gonic/gin"
)

func GetExternalAccountListing(c *gin.Context) {
	m.GetExternalAccountListing(c)
}

func GetUnreconciledTransactions(c *gin.Context) {
	var request m.UnreconciledTransactionsRequest

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err, unreconciledTransactionsResult := m.UnreconciledTransactions(request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	c.JSON(200, unreconciledTransactionsResult)
}
