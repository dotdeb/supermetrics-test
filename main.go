package main

import (
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/dotdeb/supermetrics-test/internals/api"
	"github.com/dotdeb/supermetrics-test/internals/configs"
	"github.com/dotdeb/supermetrics-test/internals/data"
	"github.com/dotdeb/supermetrics-test/internals/utils"
)

func main() {
	utils.InitLogger()

	log.Info().Msg("Initialize server")
	env := &configs.Env{
		Jwt:       configs.Jwt{},
		Variables: configs.Variables{},
	}
	if err := env.Load(); err != nil {
		log.Error().Msg("Could not read envs: " + err.Error())
		os.Exit(1)
	}

	log.Info().Msg("Initialize database")
	var db data.Database = &data.HardCodedDb{}

	log.Info().Msg("Create server")
	server := &http.Server{
		Addr:        ":" + env.Variables.Port,
		ReadTimeout: time.Duration(env.Variables.Timeout) * time.Second,
		Handler: api.SetupRouter(api.RouterSettings{
			DB:  db,
			Jwt: env.Jwt,
		}),
	}

	log.Info().Msg("Start server")
	if err := server.ListenAndServe(); err != nil {
		log.Printf("Failed to start server: %v", err)
	}
}
