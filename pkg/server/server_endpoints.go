package server

import (
	"encoding/json"
	"io"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (stg *Settings) GetWelcome(c *gin.Context) {
	var payload friend

	req := c.Request

	defer req.Body.Close()
	bdy, err := io.ReadAll(req.Body)

	if err != nil {
		stg.LogError(err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(bdy, &payload)

	if err != nil {
		stg.LogError(err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if payload.Ip == nil {
		payload.Ip = net.ParseIP(c.ClientIP())
	}

	stg.MyFriends.Add(payload)
}
