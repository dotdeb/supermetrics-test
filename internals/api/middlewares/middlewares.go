package middlewares

import (
	"fmt"

	"github.com/dotdeb/supermetrics-test/internals/data"
	"github.com/gin-gonic/gin"
)

// Log api calls.
// Currently logs are not used.
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("Api call:", c.Request.Host+c.Request.URL.Path, " auth:", c.GetHeader("Authorization")[:10])
		c.Next()
	}
}

// Pass database GetUsers-method to target controller via context.
// This ensures that database access is limited for controllers via router.
func UserReadMiddleware(db data.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("fetch_db_users", db.GetUsers)
		c.Next()
	}
}
