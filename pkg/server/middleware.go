package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (stg *Settings) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == stg.Auth {
			stg.LogDebug("Auth matched for host", c.ClientIP(), "and path", c.Request.URL.Path)
			c.Next()
			return
		}
		stg.LogDebug("Auth mismatched for host", c.ClientIP(), "and path", c.Request.URL.Path)

		c.AbortWithStatus(http.StatusUnauthorized)
	}
}
