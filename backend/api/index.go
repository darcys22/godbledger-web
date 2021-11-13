package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Index(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "index.html", nil)
}

func Reports(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "reports.html", nil)
}

func Reconcile(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "reconcile.html", nil)
}

func Accounts(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "accounts.html", nil)
}

func Modules(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "modules.html", nil)
}

func User(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "user.html", nil)
}

func Admin(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "admin.html", nil)
}

func LoginView(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "login.view.html", nil)
}
