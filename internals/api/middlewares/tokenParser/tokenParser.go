package tokenparser

import "github.com/golang-jwt/jwt"

type Parser interface {
	Parse(string) (*jwt.Token, error)
}
