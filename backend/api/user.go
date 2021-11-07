package api

import (
	"net/http"
	"net/url"

	"github.com/darcys22/godbledger-web/backend/auth"
	m "github.com/darcys22/godbledger-web/backend/models"
	"github.com/darcys22/godbledger-web/backend/models/sqlite"

	"github.com/gin-gonic/gin"
	"github.com/dgrijalva/jwt-go"

)

var (
	users sqlite.UserModel = sqlite.New("sqlite.db")
)

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

		ctx.JSON(200, current_user.Settings())
}

func ChangePassword(ctx *gin.Context) {
	var journal m.PostJournalCommand

	if err := ctx.BindJSON(&journal); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := journal.Save(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	ctx.JSON(200, journal)
}

func DefaultCurrency(ctx *gin.Context) {
	var currency m.PostCurrencyCommand

	if err := ctx.BindJSON(&currency); err != nil {
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

	current_user.Currency = currency.Currency

	if err := users.Save(current_user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	ctx.JSON(200, currency)
}

func DefaultLocale(ctx *gin.Context) {
	var locale m.PostLocaleCommand

	if err := ctx.BindJSON(&locale); err != nil {
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

	current_user.DateLocale = locale.Locale

	if err := users.Save(current_user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	ctx.JSON(200, locale)
}
