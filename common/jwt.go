package common

import (
	"gin-use-demo/model"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var jwtKey = []byte("gin_secret_key")

type Claims struct {
	UserId uint
	jwt.StandardClaims
	UserName   string
	UserNumber string
}

// ReleaseToken 生成token
func ReleaseToken(user model.User) (string, error) {
	// token设置七天过期
	expirationTime := time.Now().Add(7 * 24 * time.Hour)
	claims := &Claims{
		UserId:     uint(user.Id),
		UserName:   user.Name,
		UserNumber: user.Number,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "02-gin-demo",
			Subject:   "user token",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ParseToken 解析token
func ParseToken(tokenString string) (*jwt.Token, *Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (i interface{}, err error) {
		return jwtKey, nil
	})

	return token, claims, err
}
