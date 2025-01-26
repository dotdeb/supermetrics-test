package api

import (
	"github.com/dotdeb/supermetrics-test/internals/api/controllers"
	"github.com/dotdeb/supermetrics-test/internals/api/middlewares"
	tokenparser "github.com/dotdeb/supermetrics-test/internals/api/middlewares/tokenParser"
	"github.com/dotdeb/supermetrics-test/internals/configs"
	"github.com/dotdeb/supermetrics-test/internals/data"
	"github.com/gin-gonic/gin"
)

type RouterSettings struct {
	DB  data.Database
	Jwt configs.Jwt
}

func SetupRouter(settings RouterSettings) *gin.Engine {
	r := gin.Default()

	hmac := tokenparser.HmacParser{
		Secret: []byte(settings.Jwt.ApiSecret),
	}
	rsa := tokenparser.RsaParser{
		PublicCert: settings.Jwt.RsaPublic,
	}

	r.Use(middlewares.LoggingMiddleware())
	r.Use(middlewares.BearerAuthMiddleware(settings.Jwt, hmac, rsa))

	r.GET("/users", middlewares.UserReadMiddleware(settings.DB), controllers.GetUsers)

	return r
}
