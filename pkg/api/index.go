package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func NotFound(c *gin.Context) {

	c.HTML(http.StatusNotFound, "index.html", nil)
}
