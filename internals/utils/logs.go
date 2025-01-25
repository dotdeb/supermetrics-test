package utils

import (
	"github.com/gin-gonic/gin"
)

func CallId(c *gin.Context) string {
	id_tmp, _ := c.Get(LOG_ID)
	return id_tmp.(string)
}
