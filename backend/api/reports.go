package api

import (
	m "github.com/darcys22/godbledger-web/backend/models"
	"github.com/darcys22/godbledger-web/backend/models/reports"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ReportsResults(c *gin.Context) {
	var request reports.ReportsRequest

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err, reportResult := m.NewReport(request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	c.JSON(200, reportResult)
}
