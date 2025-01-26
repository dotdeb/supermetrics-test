package tokenparser

import (
	"crypto/rsa"
	"fmt"

	"github.com/golang-jwt/jwt"
)

type RsaParser struct {
	PublicCert *rsa.PublicKey
}

func (r *RsaParser) Parse(userToken string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(userToken, &ApiClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return r.PublicCert, nil
	})

	return token, err
}
