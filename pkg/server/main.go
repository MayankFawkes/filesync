package server

import (
	"flag"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mayankfawkes/filesync/pkg/watch"
)

func setupLogging(lg Logging) {
	if !lg.Enable {
		return
	}
	flag.Set("logtostderr", "true")
	// flag.Set("alsologtostderr", "true")
	// flag.Set("log_dir", "filesync.log")

	if lg.Debug {
		flag.Set("v", "2")
	}

}

func Server(stg *Settings) {

	// Setup logging
	setupLogging(stg.Logging)

	stg.LogInfo("Filesync starting, Version:", os.Getenv("APP_VERSION"), "GitSHA", os.Getenv("GIT_SHA"))

	stg.MyFriends = &friends{m: make(map[string]friend)}
	stg.MyFiles = &fileNhash{m: make(Dict)}

	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	r.Use(stg.AuthMiddleware())

	if stg.Server {
		stg.LogInfo("Filesync Master initiating")
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
		stg.LogInfo("Filesync Client initiating")
		go stg.InitClient()

		// client endpoints
		r.POST("/created", stg.GetCreated)
		r.POST("/modified", stg.GetModified)
		r.DELETE("/deleted", stg.GetDeleted)
		r.GET("/ack", stg.GetAck)
	}

	stg.InitFiles()
	stg.LogInfo("Filesync initiating files")

	// run the server
	stg.LogInfo("Filesync starting on port:", stg.Port)
	r.Run(fmt.Sprintf(":%d", stg.Port))

}
