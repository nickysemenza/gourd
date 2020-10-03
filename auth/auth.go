package auth

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	gauth "google.golang.org/api/oauth2/v2"
)

type Claims struct {
	User *gauth.Userinfo `json:"user_info"`
	jwt.StandardClaims
}

func New(key string) (*Auth, error) {
	if key == "" {
		return nil, fmt.Errorf("auth: key is empty")
	}
	return &Auth{jwtKey: []byte(key)}, nil
}

type Auth struct {
	jwtKey []byte
}

func (a *Auth) GetJWT(user *gauth.Userinfo) (string, error) {
	expirationTime := time.Now().Add(10 * 24 * time.Hour)
	claims := &Claims{
		User: user,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(a.jwtKey)

}
func (a *Auth) ParseJWT(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return a.jwtKey, nil
	})

	if err != nil {
		return nil, err
	}
	if !tkn.Valid {
		return nil, fmt.Errorf("failed to validate token")
	}
	return claims, nil
}
