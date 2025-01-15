package api

import (
	"github.com/dotdeb/supermetrics-test/internals/api/controllers"
	"github.com/dotdeb/supermetrics-test/internals/api/middlewares"
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

	r.Use(middlewares.LoggingMiddleware())
	r.Use(middlewares.BearerAuthMiddleware(settings.Jwt))

	r.GET("/users", middlewares.UserReadMiddleware(settings.DB), controllers.GetUsers)

	return r
}
