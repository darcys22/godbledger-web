package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"

	"github.com/darcys22/godbledger-web/backend/auth"
	"github.com/darcys22/godbledger-web/backend/setting"
)

var (
	loginService auth.LoginService
	jwtService   auth.JWTService
	loginController  LoginController
)

type LoginController struct {
	loginService auth.LoginService
	jwtService   auth.JWTService
}

func InitLoginHandler() {
	loginService = auth.StaticLoginService(setting.GetConfig())
	jwtService = auth.JWTAuthService()
	loginController = LoginController{
		loginService: loginService,
		jwtService:   jwtService,
	}
}

type LoginCredentials struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}

func (controller *LoginController) Login(ctx *gin.Context) string {
	var credential LoginCredentials
	err := ctx.ShouldBind(&credential)
	if err != nil {
		return "no data found"
	}
	isUserAuthenticated := controller.loginService.LoginUser(credential.Email, credential.Password)
	if isUserAuthenticated {
		return controller.jwtService.GenerateToken(credential.Email, true)

	}
	return ""
}

func Login(ctx *gin.Context) {
	token := loginController.Login(ctx)
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
