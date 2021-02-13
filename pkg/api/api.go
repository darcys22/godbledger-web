package api

import (
	"github.com/gin-gonic/gin"

	"github.com/darcys22/godbledger-web/pkg/middleware"
	"github.com/darcys22/godbledger-web/pkg/service"
)

var (
	loginService service.LoginService = service.StaticLoginService()
	jwtService   service.JWTService   = service.JWTAuthService()
	lController  LoginController      = LoginHandler(loginService, jwtService)
)

// Register adds http routes
func Register(r *gin.Engine) {

	// Unauthenticated Views
	r.GET("/logout", Logout)
	r.GET("/login", LoginView)
	r.POST("/login", Login)

	// Authenticated Views
	r.GET("/", middleware.AuthorizeJWT(), Index)
	r.GET("/reports", middleware.AuthorizeJWT(), Reports)
	r.POST("api/reports", middleware.AuthorizeJWT(), ReportsResults)

	r.GET("/api/journals", middleware.AuthorizeJWT(), GetJournals)
	r.POST("/api/journals", middleware.AuthorizeJWT(), PostJournal)
	r.DELETE("/api/journals/:id", middleware.AuthorizeJWT(), DeleteJournal)
	r.GET("/api/journals/:id", middleware.AuthorizeJWT(), GetJournal)
	r.POST("/api/journals/:id", middleware.AuthorizeJWT(), EditJournal)

	r.GET("/api/accounts/list", middleware.AuthorizeJWT(), GetAccountListing)
}
