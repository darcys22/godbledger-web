package api

import (
	m "github.com/darcys22/godbledger-web/pkg/models"
	"github.com/go-macaron/binding"
	"gopkg.in/macaron.v1"
)

// Register adds http routes
func Register(r *macaron.Macaron) {
	//reqSignedIn := middleware.Auth(&middleware.AuthOptions{ReqSignedIn: true})
	//reqGrafanaAdmin := middleware.Auth(&middleware.AuthOptions{ReqSignedIn: true, ReqGrafanaAdmin: true})
	//reqEditorRole := middleware.RoleAuth(m.ROLE_EDITOR, m.ROLE_ADMIN)
	//reqAccountAdmin := middleware.RoleAuth(m.ROLE_ADMIN)
	//bind := binding.Bind

	r.Get("/", Index)

	r.Get("/api/journals", GetJournals)
	r.Post("/api/journals", binding.Bind(m.PostJournalCommand{}, PostJournal))

	r.NotFound(NotFound)
}
