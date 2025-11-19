package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var signSecret []byte

type Claims struct {
	Sub string `json:"sub"`
	Exp int64  `json:"exp"`
}

func Init(secret string) {
	signSecret = []byte(secret)
}

func Sign(c Claims) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": c.Sub,
		"exp": c.Exp,
	})
	return t.SignedString(signSecret)
}

func Parse(tokenStr string) (*Claims, error) {
	t, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return signSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if !t.Valid {
		return nil, jwt.ErrSignatureInvalid
	}
	mc := t.Claims.(jwt.MapClaims)
	return &Claims{
		Sub: mc["sub"].(string),
		Exp: int64(mc["exp"].(float64)),
	}, nil
}
