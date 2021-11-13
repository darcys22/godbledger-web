package api

import (
	m "github.com/darcys22/godbledger-web/backend/models"
	"github.com/darcys22/godbledger-web/backend/auth"
	"github.com/gin-gonic/gin"
	"net/http"
)

func NewUser(ctx *gin.Context) {
	var new_user m.PostNewUserCommand

	if err := ctx.BindJSON(&new_user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cookie, err := ctx.Request.Cookie("access_token")
	if err != nil {
		respondWithError(ctx, "Cookie required")
		return
	}
	tokenString := cookie.Value
	username, err := auth.JWTAuthService().ParseUser(tokenString)
	if err != nil {
		respondWithError(ctx, "Invalid API token")
		return
	}

	current_user, err := users.Get(username)
	if err != nil {
		respondWithError(ctx, "Could not find user")
		return
	}

	if current_user.Role != "admin" {
		respondWithError(ctx, "Unauthorised")
		return
	}

	if err := users.Insert(new_user.Username, new_user.Email, new_user.Password); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}
