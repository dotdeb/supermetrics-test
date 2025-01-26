package controllers

import (
	"fmt"

	"github.com/dotdeb/supermetrics-test/internals/data"
	"github.com/dotdeb/supermetrics-test/internals/utils"
	"github.com/gin-gonic/gin"

	"github.com/rs/zerolog/log"
)

func GetUsers(ctx *gin.Context) {
	id := utils.CallId(ctx)

	isReader, err := isReader(ctx)
	if err != nil {
		log.Error().Str("id", id).Msg(err.Error())
		ctx.AbortWithStatusJSON(500, gin.H{"Internal server error": err.Error()})
		return
	}
	if !isReader {
		log.Info().Str("id", id).Msg("is not reader")
		ctx.AbortWithStatusJSON(403, gin.H{"error": "Forbidden. Insufficient Permissions"})
		return
	}

	users, exists := ctx.Get(utils.GET_USERS_FUNC)
	if !exists {
		log.Error().Str("id", id).Msg("utils.GET_USERS_FUNC is empty")
		ctx.AbortWithStatusJSON(500, gin.H{"Internal server error": "Cannot query database"})
		return
	}

	ctx.JSON(200, gin.H{"users": users.(func() []data.User)()})
}

// Verify that context has passed user roles & query coming from the reader-role
func isReader(ctx *gin.Context) (bool, error) {
	isReader := false
	roles, exists := ctx.Get(utils.CALLER_USER_ROLE)
	if !exists {
		return false, fmt.Errorf("Cannot query user roles")
	}
	for _, role := range roles.([]string) {
		if role == "reader" {
			isReader = true
			break
		}
	}
	return isReader, nil
}
