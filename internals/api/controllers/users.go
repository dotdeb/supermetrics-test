package controllers

import (
	"fmt"

	"github.com/dotdeb/supermetrics-test/internals/data"
	"github.com/gin-gonic/gin"
)

func GetUsers(ctx *gin.Context) {
	isReader, err := isReader(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(500, gin.H{"Internal server error": err.Error()})
		return
	}
	if !isReader {
		ctx.AbortWithStatusJSON(403, gin.H{"error": "Forbidden. Insufficient Permissions"})
		return
	}

	users, exists := ctx.Get("fetch_db_users")
	if !exists {
		ctx.AbortWithStatusJSON(500, gin.H{"Internal server error": "Cannot query database"})
		return
	}

	ctx.JSON(200, gin.H{"users": users.(func() []data.User)()})
}

// Verify that context has passed user roles & query coming from the reader-role
func isReader(ctx *gin.Context) (bool, error) {
	isReader := false
	roles, exists := ctx.Get("user_roles")
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
