package api

import (
	m "github.com/darcys22/godbledger-web/backend/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ChangePassword(c *gin.Context) {
	var journal m.PostJournalCommand

	if err := c.BindJSON(&journal); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := journal.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(200, journal)
}

func DefaultCurrency(c *gin.Context) {
	var journal m.PostJournalCommand

	if err := c.BindJSON(&journal); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := journal.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(200, journal)
}
