package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
)

func LoginView(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "login.view.html", nil)
}

func Login(ctx *gin.Context) {
	token := lController.Login(ctx)
	if token != "" {
		ctx.SetCookie("access_token", token, 60*60*48, "/", "", false, true)
		location := url.URL{Path: "/"}
		ctx.Redirect(http.StatusFound, location.RequestURI())
	} else {
		ctx.JSON(http.StatusUnauthorized, nil)
	}
}

func Logout(ctx *gin.Context) {
	ctx.SetCookie("access_token", "", 0, "/", "", false, true)
	location := url.URL{Path: "/login"}
	ctx.Redirect(http.StatusFound, location.RequestURI())
}
