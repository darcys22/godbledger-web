package api

import (
	"net/http"
	"net/url"

	"github.com/darcys22/godbledger-web/backend/auth"
	m "github.com/darcys22/godbledger-web/backend/models"

	"github.com/gin-gonic/gin"
	"github.com/dgrijalva/jwt-go"

)

type UserSettings struct {
	// Simply the username/email will be displayed in client
	Name            string `json:"name"`
	// Admin or Regular user, will allow for hiding admin screens but server side will also check
	Role            string `json:"role"`
	// Used for date parsing - https://sugarjs.com/docs/#/DateLocales
	DateLocale      string `json:"datelocale"`
	// USD - will be used by client for all currency items
	DefaultCurrency string `json:"defaultcurrency"`
}

func respondWithError(ctx *gin.Context, message interface{}) {
	log.Debugf("Error processing JWT: %v", message)
	ctx.Abort()
	location := url.URL{Path: "/login"}
	ctx.Redirect(http.StatusFound, location.RequestURI())
}

func AuthorizeJWT() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cookie, err := ctx.Request.Cookie("access_token")
		if err != nil {
			respondWithError(ctx, "Cookie required")
			return
		}
		tokenString := cookie.Value
		token, err := auth.JWTAuthService().ValidateToken(tokenString)
		if err != nil {
			respondWithError(ctx, err)
			return
		} else {
			if token.Valid {
				claims := token.Claims.(jwt.MapClaims)
				log.Println(claims)
			} else {
				respondWithError(ctx, "Invalid API token")
				return
			}
		}
		ctx.Next()
	}
}

func GetUserSettings(ctx *gin.Context) {
		settings := UserSettings{}
		cookie, err := ctx.Request.Cookie("access_token")
		if err != nil {
			respondWithError(ctx, "Cookie required")
			return
		}
		tokenString := cookie.Value
		user, err := auth.JWTAuthService().ParseUser(tokenString)
		if err != nil {
			respondWithError(ctx, "Invalid API token")
			return
		} else {
			settings.Name = user
			//TODO sean put actual currency here
			settings.DefaultCurrency = "USD"
			//TODO sean put actual currency here
			settings.DateLocale = "en-AU"
			//settings.DateLocale = "en-US"
			//TODO sean put actual role here
			settings.Role = "Admin"
		}
		ctx.JSON(200, settings)
}

func ChangePassword(c *gin.Context) {
	var journal m.PostJournalCommand

	if err := c.BindJSON(&journal); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := journal.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(200, journal)
}

func DefaultCurrency(c *gin.Context) {
	var journal m.PostJournalCommand

	if err := c.BindJSON(&journal); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := journal.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(200, journal)
}
