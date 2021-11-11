package api

import (
	"path"

	"github.com/gin-gonic/contrib/gzip"
	"github.com/gin-gonic/gin"

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
	r.GET("/", AuthorizeJWT(), Index)
	r.GET("/api/journals", AuthorizeJWT(), GetJournals)
	r.POST("/api/journals", AuthorizeJWT(), PostJournal)
	r.GET("/api/journals/:id", AuthorizeJWT(), GetJournal)
	r.POST("/api/journals/:id", AuthorizeJWT(), EditJournal)
	r.DELETE("/api/journals/:id", AuthorizeJWT(), DeleteJournal)

	// Chart of Accounts Page
	r.GET("/accounts", AuthorizeJWT(), Accounts)
	r.GET("/api/accounts", AuthorizeJWT(), GetAccounts)
	r.POST("/api/accounts", AuthorizeJWT(), PostAccount)
	r.GET("/api/accounts/:id", AuthorizeJWT(), GetAccount)
	r.DELETE("/api/accounts/:id", AuthorizeJWT(), DeleteAccount)
	r.POST("/api/accounttags", AuthorizeJWT(), PostAccountTag)
	r.DELETE("/api/accounttags/:account/:tagid", AuthorizeJWT(), DeleteAccountTag)

	// Reconciliation Page
	r.GET("/reconcile", AuthorizeJWT(), Reconcile)
	r.GET("/api/reconcile/listexternalaccounts", AuthorizeJWT(), GetExternalAccountListing)
	r.POST("/api/reconcile/listunreconciledtransactions", AuthorizeJWT(), GetUnreconciledTransactions)

	// Reports Page
	r.GET("/reports", AuthorizeJWT(), Reports)
	r.POST("api/reports/", AuthorizeJWT(), ReportsResults)

	// Modules Page
	r.GET("/modules", AuthorizeJWT(), Modules)

	// Users Page
	r.GET("/user", AuthorizeJWT(), User)
	r.GET("/api/user/settings", AuthorizeJWT(), GetUserSettings)
	r.POST("/api/user/changepassword", AuthorizeJWT(), ChangePassword)
	r.POST("/api/user/defaultcurrency", AuthorizeJWT(), DefaultCurrency)
	r.POST("/api/user/defaultlocale", AuthorizeJWT(), DefaultLocale)

	// Admin Page
	r.GET("/admin", AuthorizeJWT(), Admin)
	r.POST("/api/admin/newuser", AuthorizeJWT(), NewUser)

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


