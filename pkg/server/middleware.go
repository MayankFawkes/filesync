package server

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (stg *Settings) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == stg.Auth {
			log.Println("Auth matched")
			c.Next()
			return
		}

		c.AbortWithStatus(http.StatusUnauthorized)
	}
}
