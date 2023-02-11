package server

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/mayankfawkes/filesync/pkg/watch"
)

func Server(stg Settings) {
	stg.MyFriends = make(friends)
	stg.MyFiles = make(fileNhash)

	gin.SetMode(gin.ReleaseMode)

	var r *gin.Engine

	log.Println(stg.Logging)

	if stg.Logging {
		r = gin.Default()
	} else {
		r = gin.New()
	}

	r.Use(stg.AuthMiddleware())

	if stg.Server {
		go stg.InitServer()

		// server endpoints
		r.POST("/welcome", stg.GetWelcome)

		ch := make(chan watch.Response)

		go watch.WatchPath(stg.WatchPath, ch)

		go func() {
			for d := range ch {
				if d.Status == watch.CREATED {
					go stg.SendCreated(d.Path)
				} else if d.Status == watch.MODIFIED {
					go stg.SendModified(d.Path)
				} else if d.Status == watch.DELETED {
					go stg.SendDeleted(d.Path)
				}
			}
		}()
	} else {
		fmt.Println("Client init done")
		go stg.InitClient()

		// client endpoints
		r.POST("/created", stg.GetCreated)
		r.POST("/modified", stg.GetModified)
		r.DELETE("/deleted", stg.GetDeleted)
		r.GET("/ack", stg.GetAck)
	}

	stg.InitFiles()

	fmt.Println(stg.MyFriends)

	// run the server
	r.Run(fmt.Sprintf(":%d", stg.Port))

}
