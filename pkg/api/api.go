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
	r.GET("/reports", Reports)

	r.GET("/api/journals", GetJournals)
	r.POST("/api/journals", PostJournal)
	r.DELETE("/api/journals/:id", DeleteJournal)
	r.GET("/api/journals/:id", GetJournal)
	r.POST("/api/journals/:id", EditJournal)

	r.GET("/api/accounts/list", GetAccountListing)
}
