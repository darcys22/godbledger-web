package middleware

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/darcys22/godbledger-web/pkg/auth"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithField("prefix", "JWTMiddleware")

func respondWithError(ctx *gin.Context, message interface{}) {
	log.Debugf("Error processing JWT: ", message)
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
				fmt.Println(claims)
			} else {
				respondWithError(ctx, "Invalid API token")
				return
			}
		}
		ctx.Next()
	}
}
