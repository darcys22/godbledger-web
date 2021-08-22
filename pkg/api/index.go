package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func Reports(c *gin.Context) {
	c.HTML(http.StatusOK, "reports.html", nil)
}

func Reconcile(c *gin.Context) {
	c.HTML(http.StatusOK, "reconcile.html", nil)
}

func Modules(c *gin.Context) {
	c.HTML(http.StatusOK, "modules.html", nil)
}

func User(c *gin.Context) {
	c.HTML(http.StatusOK, "user.html", nil)
}

func Admin(c *gin.Context) {
	c.HTML(http.StatusOK, "admin.html", nil)
}

func NotFound(c *gin.Context) {
	c.HTML(http.StatusNotFound, "index.html", nil)
}
