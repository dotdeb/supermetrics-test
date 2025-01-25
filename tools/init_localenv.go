/*
This script creates local environment for developer usage
*/
package main

import (
	"errors"
	"log"
	"os"
	"time"

	"math/rand"
	"strings"

	"github.com/dotdeb/supermetrics-test/internals/configs"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func main() {
	dotenv := "../.env"
	writeDotEnv(dotenv)

	e := &configs.Env{}
	loadDotEnv(dotenv)
	resolveEnvs(e)

	token := constructToken(e)
	writeTestFile("example.http", token)
}

func writeTestFile(file string, token string) {
	os.Remove(file)
	url := "GET http://localhost:8000/users\n"
	header := "Authorization: Bearer "
	warning := "\n\n# This file is autocreated. Run createTestToken.go to recreate"
	os.WriteFile(file, []byte(url+header+token+warning), 0400)
}

// Construct token which is used when calling /users api
func constructToken(e *configs.Env) string {
	var (
		key       []byte
		token     *jwt.Token
		signedJwt string
	)

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

	return signedJwt
}

// Create .env-file with environment variables required by the application
func writeDotEnv(file string) {
	os.Remove(file)
	api_secret := "API_SECRET=" + getRandomToken(24)
	issuer := "JWT_ISSUER=local.developer"
	audience := "JWT_AUDIENCE=local.app"
	port := "PORT=8000"
	timeout := "TIMEOUT_SEC=5"
	warning := "\n\n# This file is autocreated. Run init_localenv.go to recreate"
	os.WriteFile(file, []byte(
		api_secret+"\n"+
			issuer+"\n"+
			audience+"\n"+
			port+"\n"+
			timeout+"\n"+
			warning,
	), 0400)

}

func getRandomToken(size int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_!&%"
	sb := strings.Builder{}
	sb.Grow(size)
	for i := 0; i < size; i++ {
		sb.WriteByte(charset[rand.Intn(len(charset))])
	}
	return sb.String()
}

// Read .env variables to environment variables
func loadDotEnv(file string) {
	err := godotenv.Load(file)
	if err != nil {
		log.Fatal(err)
	}
}

// Read envs needed to construct test token
func resolveEnvs(env *configs.Env) error {
	env.Jwt.Issuer = os.Getenv("JWT_ISSUER")
	if env.Jwt.Issuer == "" {
		return errors.New("Missing environment variable: JWT_ISSUER")
	}
	env.Jwt.Audience = os.Getenv("JWT_AUDIENCE")
	if env.Jwt.Audience == "" {
		return errors.New("Missing environment variable: JWT_AUDIENCE")
	}
	env.Jwt.ApiSecret = os.Getenv("API_SECRET")
	if env.Jwt.ApiSecret == "" {
		return errors.New("Missing environment variable: API_SECRET")
	}
	return nil
}
