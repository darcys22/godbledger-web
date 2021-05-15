package service

import (
	"github.com/darcys22/godbledger-web/pkg/models/sqlite"
)

type LoginService interface {
	LoginUser(email string, password string) bool
}
type loginInformation struct {
	users *sqlite.UserModel
}

func StaticLoginService() LoginService {
	database := sqlite.New("sqlite.db")

	return &loginInformation{users: &database}
}
func (info *loginInformation) LoginUser(email string, password string) bool {
	_, err := info.users.Authenticate(email, password)
	if err != nil {
		return false
	}
	return true
}
