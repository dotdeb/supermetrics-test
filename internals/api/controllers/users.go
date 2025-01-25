package controllers

import (
	"fmt"

	"github.com/dotdeb/supermetrics-test/internals/data"
	"github.com/dotdeb/supermetrics-test/internals/utils"
	"github.com/gin-gonic/gin"
)

func GetUsers(ctx *gin.Context) {
	id_tmp, _ := ctx.Get(utils.LOG_ID)
	id := id_tmp.(string) + ":"

	isReader, err := isReader(ctx)
	if err != nil {
		fmt.Println(id, err.Error())
		ctx.AbortWithStatusJSON(500, gin.H{"Internal server error": err.Error()})
		return
	}
	if !isReader {
		fmt.Println(id, "is not reader")
		ctx.AbortWithStatusJSON(403, gin.H{"error": "Forbidden. Insufficient Permissions"})
		return
	}

	users, exists := ctx.Get(utils.GET_USERS_FUNC)
	if !exists {
		fmt.Println(id, "utils.GET_USERS_FUNC is empty")
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
