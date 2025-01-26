package tokenparser

import (
	"fmt"

	"github.com/golang-jwt/jwt"
)

type HmacParser struct {
	Secret []byte
}

func (h *HmacParser) Parse(userToken string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(userToken, &ApiClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return h.Secret, nil
	})
}
