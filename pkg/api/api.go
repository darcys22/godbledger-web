package api

import (
	//"fmt"
	//m "github.com/darcys22/godbledger-web/pkg/models"
	//"github.com/go-macaron/binding"
	//"gopkg.in/macaron.v1"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Register adds http routes
func Register(r *gin.Engine) {
	//reqSignedIn := middleware.Auth(&middleware.AuthOptions{ReqSignedIn: true})
	//reqGrafanaAdmin := middleware.Auth(&middleware.AuthOptions{ReqSignedIn: true, ReqGrafanaAdmin: true})
	//reqEditorRole := middleware.RoleAuth(m.ROLE_EDITOR, m.ROLE_ADMIN)
	//reqAccountAdmin := middleware.RoleAuth(m.ROLE_ADMIN)
	//bind := binding.Bind

	r.GET("/", Index)

	type ContactForm struct {
		Name string `json:"name" binding:"Required"`
	}

	r.GET("/api/journals", GetJournals)
	r.POST("/api/journals", func(c *gin.Context) {
		c.Bind(&ContactForm{})
		c.String(http.StatusOK, "ok")
	})

	r.GET("/api/accounts/list", GetAccountListing)

	//r.NotFound(NotFound)
}
