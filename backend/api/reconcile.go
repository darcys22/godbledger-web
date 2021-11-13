package api

import (
	"net/http"

	m "github.com/darcys22/godbledger-web/backend/models"
	"github.com/gin-gonic/gin"
)

func GetExternalAccountListing(ctx *gin.Context) {
	m.GetExternalAccountListing(ctx)
}

func GetUnreconciledTransactions(ctx *gin.Context) {
	var request m.UnreconciledTransactionsRequest

	if err := ctx.BindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err, unreconciledTransactionsResult := m.UnreconciledTransactions(request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	ctx.JSON(http.StatusOK, unreconciledTransactionsResult)
}
