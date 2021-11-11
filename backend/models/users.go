package models

import (
	"time"
)

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
	Active         bool
	Currency			 string
	DateLocale		 string
	Role           string
}

type PostCurrencyCommand struct {
	Currency string `json:"currency" binding:"required"`
}

type PostLocaleCommand struct {
	Locale string `json:"locale" binding:"required"`
}

type PostPasswordChangeCommand struct {
	Password string `json:"password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
	ConfirmNewPassword string `json:"confirm_new_password" binding:"required"`
}

type PostNewUserCommand struct {
	Username string `json:"username" binding:"required"`
	Email string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserSettingsResponse struct {
	// Simply the username/email will be displayed in client
	Name            string `json:"name"`
	// Admin or Regular user, will allow for hiding admin screens but server side will also check
	Role            string `json:"role"`
	// Used for date parsing - https://sugarjs.com/docs/#/DateLocales
	DateLocale      string `json:"datelocale"`
	// USD - will be used by client for all currency items
	DefaultCurrency string `json:"defaultcurrency"`
}

func (u *User) Settings() UserSettingsResponse {
	return UserSettingsResponse{u.Name, u.Currency, u.DateLocale, u.Role}
}
