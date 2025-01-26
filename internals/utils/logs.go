package utils

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func CallId(c *gin.Context) string {
	id_tmp, _ := c.Get(LOG_ID)
	return id_tmp.(string)
}

func InitLogger() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel) // Default log level
	log.Info().Msg("Initialize logger")

	logLevel, _ := os.LookupEnv("LOG_LEVEL")
	logLevel = strings.ToLower(logLevel)
	if logLevel == "debug" || logLevel == "0" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)

	}

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
}
