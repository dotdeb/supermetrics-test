package middlewares

import (
	"fmt"
	"strings"
	"time"

	"github.com/dotdeb/supermetrics-test/internals/configs"
	"github.com/dotdeb/supermetrics-test/internals/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type ApiClaims struct {
	jwt.StandardClaims
	Roles []string
}

// Verify that api-call has sufficient token. This doesn't check user roles
// but pass them to controllers via context.
func BearerAuthMiddleware(system configs.Jwt) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := utils.CallId(c)

		userToken, err := parseBearer(c.GetHeader("Authorization"))
		if err != nil {
			fmt.Println(id+":", err.Error())
			c.AbortWithStatusJSON(401, gin.H{"message": "Invalid API key"})
			return
		}

		token, err := parseToken(userToken, []byte(system.ApiSecret))
		if err != nil {
			fmt.Println(id+":", err.Error())
			c.AbortWithStatusJSON(401, gin.H{"message": "Invalid API key"})
			return
		}

		claims, ok := token.Claims.(*ApiClaims)
		if !ok || !token.Valid {
			c.AbortWithStatusJSON(401, gin.H{"message": "Invalid API key"})
			return
		}

		// Require jwt token has attributes listed here. Otherwise reject token.
		unixNow := time.Now().Unix()
		if claims.Issuer != system.Issuer ||
			claims.Audience != system.Audience ||
			claims.NotBefore > unixNow ||
			claims.ExpiresAt < unixNow {
			c.AbortWithStatusJSON(401, gin.H{"message": "Invalid API key"})
			return
		}

		c.Set("user_roles", claims.Roles)
		c.Next()
	}
}

func parseToken(userToken string, apiToken []byte) (*jwt.Token, error) {
	return jwt.ParseWithClaims(userToken, &ApiClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return apiToken, nil
	})
}

// Verify that query has bearer prefix and remove this prefix
func parseBearer(queryToken string) (string, error) {
	if !strings.Contains(queryToken, "Bearer") {
		return "", fmt.Errorf("Missing bearer token")
	}
	queryToken = strings.TrimSpace(strings.Replace(queryToken, "Bearer", "", 1))
	return queryToken, nil
}
