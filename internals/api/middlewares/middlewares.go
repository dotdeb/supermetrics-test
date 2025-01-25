package middlewares

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/dotdeb/supermetrics-test/internals/data"
	"github.com/dotdeb/supermetrics-test/internals/utils"
	"github.com/gin-gonic/gin"
)

// Log api calls & create identifier for call to be identified in every step
// Currently logs are not used.
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := uuid.New().String()
		c.Set(utils.LOG_ID, id)
		c.Writer.Header().Add("trace", id) // Trace is used to trace logs for single call

		fmt.Println(id+":", c.Request.Host+c.Request.URL.Path)
		c.Next()
	}
}

// Pass database GetUsers-method to target controller via context.
// This ensures that database access is limited for controllers via router.
func UserReadMiddleware(db data.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := c.Get(utils.LOG_ID)
		if db == nil {
			fmt.Println(id.(string)+":", "Database is nil")
		}

		c.Set(utils.GET_USERS_FUNC, db.GetUsers)
		c.Next()
	}
}
