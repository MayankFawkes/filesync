package server

import (
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (stg *Settings) getWelcome(c *gin.Context) {
	var payload friend

	req := c.Request

	defer req.Body.Close()
	bdy, _ := io.ReadAll(req.Body)

	err := json.Unmarshal(bdy, &payload)

	if payload.Ip == nil {
		payload.Ip = net.ParseIP(c.ClientIP())
	}

	if err != nil {
		log.Println(err.Error())
		c.Data(http.StatusBadRequest, "text/plain; charset=utf-8", []byte(err.Error()))
		return
	}

	(*stg).MyFriends.Add(payload)
}
