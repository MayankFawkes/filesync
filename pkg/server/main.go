package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/mayankfawkes/filesync/pkg/watch"
)

func Server(stg Settings) {
	stg.MyFriends = make(friends)
	stg.MyFiles = make(fileNhash)

	r := gin.Default()

	if stg.Server {
		go stg.initServer()

		// server endpoints
		r.POST("/welcome", stg.getWelcome)

		ch := make(chan watch.Response)

		go watch.WatchPath(stg.WatchPath, ch)

		go func() {
			for d := range ch {
				if d.Status == watch.CREATED {
					go stg.sendCreated(d.Path)
				} else if d.Status == watch.MODIFIED {
					go stg.sendModified(d.Path)
				} else if d.Status == watch.DELETED {
					go stg.sendDeleted(d.Path)
				}
			}
		}()
	} else {
		fmt.Println("Client init done")
		go stg.initClient()

		// client endpoints
		r.POST("/created", stg.getCreated)
		r.POST("/modified", stg.getModified)
		r.DELETE("/deleted", stg.getDeleted)
		r.GET("/ack", stg.getAck)
	}

	stg.initFiles()

	fmt.Println(stg.MyFriends)

	// run the server
	r.Run(fmt.Sprintf(":%d", stg.Port))

}
