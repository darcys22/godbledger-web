package auth

import (
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/darcys22/godbledger-web/backend/models/sqlite"
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

func (info *loginInformation) NewUser(email string, password string) bool {
	//_, err := info.users.New(email, password)
	//if err != nil {
		//return false
	//}
	return true
}

//JWT service
type JWTService interface {
	GenerateToken(email string, isUser bool) string
	ValidateToken(token string) (*jwt.Token, error)
	ParseUser(token string) (string, error)
}
type authCustomClaims struct {
	Name string `json:"name"`
	User bool   `json:"user"`
	jwt.StandardClaims
}

type jwtServices struct {
	secretKey string
	issuer    string
}

//auth-jwt
func JWTAuthService() JWTService {
	secret := os.Getenv("SECRET")
	if secret == "" {
		secret = "secret"
	}
	return &jwtServices{
		secretKey: secret,
		issuer:    "DarcyFinancial",
	}
}

func (service *jwtServices) GenerateToken(email string, isUser bool) string {
	claims := &authCustomClaims{
		email,
		isUser,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 48).Unix(),
			Issuer:    service.issuer,
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	//encoded string
	t, err := token.SignedString([]byte(service.secretKey))
	if err != nil {
		panic(err)
	}
	return t
}

func (service *jwtServices) ValidateToken(encodedToken string) (*jwt.Token, error) {
	return jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		if _, isvalid := token.Method.(*jwt.SigningMethodHMAC); !isvalid {
			return nil, fmt.Errorf("Invalid token %v", token.Header["alg"])

		}
		return []byte(service.secretKey), nil
	})

}

//type authCustomClaims struct {
func (service *jwtServices) ParseUser(encodedToken string) (string, error) {
	token, err := jwt.ParseWithClaims(encodedToken, &authCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, isvalid := token.Method.(*jwt.SigningMethodHMAC); !isvalid {
			return nil, fmt.Errorf("Invalid token %v", token.Header["alg"])

		}
		return []byte(service.secretKey), nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*authCustomClaims); ok && token.Valid {
		return claims.Name, nil
	} else {
		return "", fmt.Errorf("Invalid Token")
	}
}
