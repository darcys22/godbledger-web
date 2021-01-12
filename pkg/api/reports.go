package api

import (
	m "github.com/darcys22/godbledger-web/pkg/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ReportsResults(c *gin.Context) {
	var request m.ReportsRequest

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
