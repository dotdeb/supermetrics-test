/*
This script creates local environment for developer usage
*/
package main

import (
	crypto "crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
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
	key := createRsaKey()

	dotenv := "../.env"
	writeDotEnv(dotenv)

	e := &configs.Env{}
	loadDotEnv(dotenv)
	resolveEnvs(e)

	token := constructHmacToken(e)
	writeTestFile("example.http", token)

	signed := createRsaToken(e, key)
	writeTestFile("example.rsa.http", signed)
}

func writeTestFile(file string, token string) {
	os.Remove(file)
	url := "GET http://localhost:8000/users\n"
	header := "Authorization: Bearer " + token
	warning := "\n\n# This file is autocreated. Run createTestToken.go to recreate"
	os.WriteFile(file, []byte(url+header+warning), 0400)
}

func createRsaToken(e *configs.Env, key *rsa.PrivateKey) string {
	claims := &jwt.MapClaims{
		"alg":   "RS256",
		"iss":   e.Jwt.Issuer,
		"aud":   e.Jwt.Audience,
		"roles": []string{"reader"},
		"exp":   time.Now().Add(60 * time.Minute).Unix(),
		"nbf":   time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	tokenString, err := token.SignedString(key)
	if err != nil {
		log.Fatal(err)
	}

	return tokenString
}

// Construct token which is used when calling /users api
func constructHmacToken(e *configs.Env) string {
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

func createRsaKey() *rsa.PrivateKey {
	privateKey, err := rsa.GenerateKey(crypto.Reader, 2048)
	if err != nil {
		log.Fatal(err)
	}

	privateKeyFile, err := os.Create("private.pem")
	if err != nil {
		log.Fatal(err)
	}
	defer privateKeyFile.Close()

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}
	if err := pem.Encode(privateKeyFile, privateKeyPEM); err != nil {
		log.Fatal(err)
	}

	publicKeyFile, err := os.Create("public.pem")
	if err != nil {
		log.Fatal(err)
	}
	defer publicKeyFile.Close()

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		log.Fatal(err)
	}
	publicKeyPEM := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	}
	if err := pem.Encode(publicKeyFile, publicKeyPEM); err != nil {
		log.Fatal(err)
	}

	return privateKey
}

func constructRsaToken(e *configs.Env) string {
	var (
		key       []byte
		token     *jwt.Token
		signedJwt string
	)
	createRsaKey()

	key = []byte(e.Jwt.ApiSecret)
	token = jwt.NewWithClaims(jwt.SigningMethodRS256,
		jwt.MapClaims{
			"alg":   "RS256",
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
	pem, _ := os.ReadFile("public.pem")

	os.Remove(file)
	api_secret := "API_SECRET=" + getRandomToken(24)
	issuer := "JWT_ISSUER=local.developer"
	audience := "JWT_AUDIENCE=local.app"
	port := "PORT=8000"
	timeout := "TIMEOUT_SEC=5"
	rsaPub := "RSA_PUB=\"" + string(pem) + "\""
	warning := "\n\n# This file is autocreated. Run init_localenv.go to recreate"
	os.WriteFile(file, []byte(
		api_secret+"\n"+
			issuer+"\n"+
			audience+"\n"+
			port+"\n"+
			timeout+"\n"+
			rsaPub+"\n"+
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
