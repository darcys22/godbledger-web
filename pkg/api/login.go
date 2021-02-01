package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func LoginView(c *gin.Context) {
	c.HTML(http.StatusOK, "login.view.html", nil)
}

func RegisterView(c *gin.Context) {
	c.HTML(http.StatusOK, "register.view.html", nil)
}

func Logout(c *gin.Context) {

	//err := hs.AuthTokenService.RevokeToken(c.Req.Context(), c.UserToken)
	//if err != nil && !errors.Is(err, models.ErrUserTokenNotFound) {
	//hs.log.Error("failed to revoke auth token", "error", err)
	//}

	//cookies.WriteSessionCookie(c, hs.Cfg, "", -1)

	//if setting.SignoutRedirectUrl != "" {
	//c.Redirect(setting.SignoutRedirectUrl)
	//} else {
	//hs.log.Info("Successful Logout", "User", c.Email)
	//c.Redirect(setting.AppSubUrl + "/login")
	//}
	location := url.URL{Path: "/login"}
	c.Redirect(http.StatusFound, location.RequestURI())
}
