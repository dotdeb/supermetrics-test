/*
This script creates a token for developer to test api calls.
*/
package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/dotdeb/supermetrics-test/internals/configs"
	"github.com/golang-jwt/jwt/v5"
)

func main() {
	var (
		key       []byte
		token     *jwt.Token
		signedJwt string
	)

	if _, err := os.Stat(".env"); errors.Is(err, os.ErrNotExist) {
		fmt.Println(".env file is not founded. Go to path where .env file exists and run command from there")
		os.Exit(1)
	}

	e := configs.Env{}
	e.Load()

	key = []byte(e.Jwt.ApiSecret)
	token = jwt.NewWithClaims(jwt.SigningMethodHS512,
		jwt.MapClaims{
			"alg":   "HS512",
			"iss":   e.Jwt.Issuer,
			"aud":   e.Jwt.Audience,
			"roles": []string{"reader"},
			"exp":   time.Now().Add(60 * time.Minute).Unix(),
			"nbf":   time.Now().Unix(),
		})
	signedJwt, err := token.SignedString(key)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println(signedJwt)
}
