package tokenparser

import "github.com/golang-jwt/jwt"

type ApiClaims struct {
	jwt.StandardClaims
	Roles []string
}
