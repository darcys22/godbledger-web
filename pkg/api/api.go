package api

import (
	//"fmt"
	"github.com/gin-gonic/gin"
	"net/http"

	"github.com/darcys22/godbledger-web/pkg/service"
)

var (
	loginService service.LoginService = service.StaticLoginService()
	jwtService   service.JWTService   = service.JWTAuthService()
	lController  LoginController      = LoginHandler(loginService, jwtService)
)

// Register adds http routes
func Register(r *gin.Engine) {

	// not logged in views
	r.GET("/logout", Logout)
	//r.Post("/login", quota("session"), bind(dtos.LoginCommand{}), routing.Wrap(hs.LoginPost))
	r.GET("/login", LoginView)

	r.POST("/login", func(ctx *gin.Context) {
		token := lController.Login(ctx)
		if token != "" {
			ctx.JSON(http.StatusOK, gin.H{
				"token": token,
			})
		} else {
			ctx.JSON(http.StatusUnauthorized, nil)
		}
	})

	r.GET("/", Index)
	r.GET("/reports", Reports)
	r.POST("api/reports", ReportsResults)

	r.GET("/api/journals", GetJournals)
	r.POST("/api/journals", PostJournal)
	r.DELETE("/api/journals/:id", DeleteJournal)
	r.GET("/api/journals/:id", GetJournal)
	r.POST("/api/journals/:id", EditJournal)

	r.GET("/api/accounts/list", GetAccountListing)
}
