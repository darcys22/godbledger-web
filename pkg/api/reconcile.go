package api

import (
	m "github.com/darcys22/godbledger-web/pkg/models"
	"github.com/gin-gonic/gin"
)

func GetExternalAccountListing(c *gin.Context) {
	m.GetExternalAccountListing(c)
}
