package configs

import (
	"errors"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Jwt struct {
	ApiSecret string
	Issuer    string
	Audience  string
}

type Variables struct {
	Port    string
	Timeout int
}

type Env struct {
	Jwt       Jwt
	Variables Variables
}

func (env *Env) Load() error {
	var err error
	if err = resolveDotEnv(); err != nil {
		return err
	}
	if err = env.resolveEnvs(); err != nil {
		return err
	}

	return nil
}

// Resolve environment variables required by application
func (env *Env) resolveEnvs() error {
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
	env.Variables.Port = os.Getenv("PORT")
	if env.Variables.Port == "" {
		return errors.New("Missing environment variable: PORT")
	}
	timeout, err := strconv.Atoi(os.Getenv("TIMEOUT_SEC"))
	if err != nil {
		return errors.New("Missing or invalid environment variable: TIMEOUT_SEC")
	}
	env.Variables.Timeout = timeout

	return nil
}

// Check if .env file is in use & import variables from there.
// If file doesn't exists then ignore.
// If file exists but is not readable then return error.
//
// CI/CD or docker can ignore .env file and pass envs directly.
func resolveDotEnv() error {
	if _, err := os.Stat("./.env"); errors.Is(err, os.ErrNotExist) {
		return nil
	} else if err != nil {
		return errors.New("Error while checking .env file location.")
	}

	if err := godotenv.Load(); err != nil {
		return errors.New("Error while reading .env file.")
	}

	return nil
}
