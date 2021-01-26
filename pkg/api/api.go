package api

import (
	//"fmt"
	"github.com/gin-gonic/gin"
	//"net/http"

	"github.com/darcys22/godbledger-web/pkg/service"
)

// Register adds http routes
func Register(r *gin.Engine) {
	var loginService service.LoginService = service.StaticLoginService()
	var jwtService service.JWTService = service.JWTAuthService()
	var loginController LoginController = LoginHandler(loginService, jwtService)

	//reqSignedIn := middleware.Auth(&middleware.AuthOptions{ReqSignedIn: true})
	//reqGrafanaAdmin := middleware.Auth(&middleware.AuthOptions{ReqSignedIn: true, ReqGrafanaAdmin: true})
	//reqEditorRole := middleware.RoleAuth(m.ROLE_EDITOR, m.ROLE_ADMIN)
	//reqAccountAdmin := middleware.RoleAuth(m.ROLE_ADMIN)
	//bind := binding.Bind

	server.POST("/login", func(ctx *gin.Context) {
		token := loginController.Login(ctx)
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
