package api

import (
	//"fmt"
	"github.com/gin-gonic/gin"
	//"net/http"
)

// Register adds http routes
func Register(r *gin.Engine) {
	//reqSignedIn := middleware.Auth(&middleware.AuthOptions{ReqSignedIn: true})
	//reqGrafanaAdmin := middleware.Auth(&middleware.AuthOptions{ReqSignedIn: true, ReqGrafanaAdmin: true})
	//reqEditorRole := middleware.RoleAuth(m.ROLE_EDITOR, m.ROLE_ADMIN)
	//reqAccountAdmin := middleware.RoleAuth(m.ROLE_ADMIN)
	//bind := binding.Bind

	r.GET("/", Index)

	//type ContactForm struct {
	//Name string `json:"name" binding:"required"`
	//}
	//r.POST("/api/journals", func(c *gin.Context) {
	//c.Bind(&ContactForm{})
	//c.String(http.StatusOK, "ok")
	//})

	r.GET("/api/journals", GetJournals)
	r.POST("/api/journals", PostJournal)

	r.GET("/api/accounts/list", GetAccountListing)
}
