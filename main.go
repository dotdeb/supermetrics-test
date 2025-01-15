package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dotdeb/supermetrics-test/internals/api"
	"github.com/dotdeb/supermetrics-test/internals/configs"
	"github.com/dotdeb/supermetrics-test/internals/data"
)

func main() {
	fmt.Println("Initialize server")
	env := &configs.Env{
		Jwt:       configs.Jwt{},
		Variables: configs.Variables{},
	}
	if err := env.Load(); err != nil {
		fmt.Println("Could not read envs: ", err)
		os.Exit(1)
	}

	fmt.Println("Initialize database")
	var db data.Database = &data.HardCodedDb{}

	fmt.Println("Create server")
	server := &http.Server{
		Addr:        ":" + env.Variables.Port,
		ReadTimeout: time.Duration(env.Variables.Timeout) * time.Second,
		Handler: api.SetupRouter(api.RouterSettings{
			DB:  db,
			Jwt: env.Jwt,
		}),
	}

	fmt.Println("Start server")
	if err := server.ListenAndServe(); err != nil {
		log.Printf("Failed to start server: %v", err)
	}
}
