package middlewares

import (
	"fmt"
	"strings"
	"time"

	tokenparser "github.com/dotdeb/supermetrics-test/internals/api/middlewares/tokenParser"
	"github.com/dotdeb/supermetrics-test/internals/configs"
	"github.com/dotdeb/supermetrics-test/internals/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/rs/zerolog/log"
)

// Verify that api-call has sufficient token. This doesn't check user roles
// but pass them to controllers via context.
func BearerAuthMiddleware(system configs.Jwt, hmac tokenparser.HmacParser, rsa tokenparser.RsaParser) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := utils.CallId(c)
		userToken, err := parseBearer(c.GetHeader("Authorization"))
		if err != nil {
			log.Info().Str("id", id).Msg(err.Error())
			c.AbortWithStatusJSON(401, gin.H{"message": "Invalid API key"})
			return
		}

		token, err := resolveToken(userToken, hmac, rsa)
		if err != nil {
			log.Info().Str("id", id).Msg(err.Error())
			c.AbortWithStatusJSON(401, gin.H{"message": "Invalid API key"})
			return
		}

		tokenRoles, err := validateToken(token, system)
		if err != nil {
			log.Info().Str("id", id).Msg(err.Error())
			c.AbortWithStatusJSON(401, gin.H{"message": "Invalid API key"})
			return
		}

		log.Info().Str("id", id).Msg("Authenticated with roles " + strings.Join(tokenRoles, " "))
		c.Set(utils.CALLER_USER_ROLE, tokenRoles)
		c.Next()
	}
}

// Validate token has correct attributes. Does not check roles as this middleware is multipurpose.
func validateToken(token *jwt.Token, system configs.Jwt) ([]string, error) {
	claims, ok := token.Claims.(*tokenparser.ApiClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("Token is not valid")
	}

	// Require jwt token has attributes listed here. Otherwise reject token.
	unixNow := time.Now().Unix()
	if claims.Issuer != system.Issuer ||
		claims.Audience != system.Audience ||
		claims.NotBefore > unixNow ||
		claims.ExpiresAt < unixNow {
		return nil, fmt.Errorf("token claims are not valid")
	}

	return claims.Roles, nil
}

// Verify that query has bearer prefix and remove this prefix
func parseBearer(queryToken string) (string, error) {
	if !strings.Contains(queryToken, "Bearer") {
		return "", fmt.Errorf("Missing bearer token")
	}
	queryToken = strings.TrimSpace(strings.Replace(queryToken, "Bearer", "", 1))
	return queryToken, nil
}

// Check token algorithm and parse token by given algorithms.
func resolveToken(userToken string, hmac tokenparser.HmacParser, rsa tokenparser.RsaParser) (*jwt.Token, error) {
	var token *jwt.Token

	alg, err := getAlg(userToken)
	if err != nil {
		return nil, err
	}

	if strings.HasPrefix(alg, "RS") { // RSA
		token, err = rsa.Parse(userToken)
	} else if strings.HasPrefix(alg, "HS") { // HMAC
		token, err = hmac.Parse(userToken)
	} else {
		return nil, fmt.Errorf("Unsupported algorithm %s", alg)
	}

	return token, err
}

// Read algorithm used in the token. Does not verify token
func getAlg(userToken string) (string, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(userToken, jwt.MapClaims{})
	if err != nil {
		return "", fmt.Errorf("Error parsing token: %v", err)
	}

	signingMethod := token.Method.Alg()
	return signingMethod, nil
}
