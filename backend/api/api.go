package api

import (
	"path"

	"github.com/gin-gonic/contrib/gzip"
	"github.com/gin-gonic/gin"

	"github.com/darcys22/godbledger-web/backend/middleware"
	"github.com/darcys22/godbledger-web/backend/setting"
	
	"github.com/sirupsen/logrus"
)

var log = logrus.WithField("prefix", "API")

func mapStatic(m *gin.Engine, dir string, prefix string) {
	headers := func() gin.HandlerFunc {
		return func(c *gin.Context) {
			c.Writer.Header().Set("Cache-Control", "public, max-age=3600")
			c.Next()
		}
	}

	if setting.Env == setting.DEV {
		headers = func() gin.HandlerFunc {
			return func(c *gin.Context) {
				c.Writer.Header().Set("Cache-Control", "max-age=0, must-revalidate, no-cache")
				c.Next()
			}
		}
	}

	m.Static(prefix, path.Join(setting.StaticRootPath, dir))
	m.Use(headers())
}

// register adds http routes
func register(r *gin.Engine) {

	// ---- Unauthenticated Views -------
	r.GET("/logout", Logout)
	r.GET("/login", LoginView)
	r.POST("/login", Login)

	// ---- Authenticated Views ---------

	// Main/Journal Entry Page
	r.GET("/", middleware.AuthorizeJWT(), Index)
	r.GET("/api/journals", middleware.AuthorizeJWT(), GetJournals)
	r.POST("/api/journals", middleware.AuthorizeJWT(), PostJournal)
	r.DELETE("/api/journals/:id", middleware.AuthorizeJWT(), DeleteJournal)
	r.GET("/api/journals/:id", middleware.AuthorizeJWT(), GetJournal)
	r.POST("/api/journals/:id", middleware.AuthorizeJWT(), EditJournal)

	r.GET("/api/accounts/list", middleware.AuthorizeJWT(), GetAccountListing)

	// Reconciliation Page
	r.GET("/reconcile", middleware.AuthorizeJWT(), Reconcile)
	r.GET("/api/reconcile/listexternalaccounts", middleware.AuthorizeJWT(), GetExternalAccountListing)
	r.POST("/api/reconcile/listunreconciledtransactions", middleware.AuthorizeJWT(), GetUnreconciledTransactions)

	// Reports Page
	r.GET("/reports", middleware.AuthorizeJWT(), Reports)
	r.POST("api/reports/", middleware.AuthorizeJWT(), ReportsResults)

	// Modules Page
	r.GET("/modules", middleware.AuthorizeJWT(), Modules)

	// Users Page
	r.GET("/user", middleware.AuthorizeJWT(), User)

	// Admin Page
	r.GET("/admin", middleware.AuthorizeJWT(), Admin)

}

func NewGin() *gin.Engine {

	m := gin.Default()
	m.Use(gin.Recovery())
	if setting.EnableGzip {
		m.Use(gzip.Gzip(gzip.DefaultCompression))
	}

	mapStatic(m, "", "public")
	mapStatic(m, "app", "app")
	mapStatic(m, "css", "css")
	mapStatic(m, "img", "img")
	mapStatic(m, "fonts", "fonts")

	m.LoadHTMLGlob(path.Join(setting.StaticRootPath, "views/*.html"))

	register(m)

	return m
}


