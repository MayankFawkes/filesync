package server

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func (stg *Settings) GetCreated(c *gin.Context) {
	path := c.GetHeader("path")
	path = stg.AbsPath(path)

	log.Println("client created", path)

	if path == "" {
		return
	}

	fp, err := os.Create(path)
	if err != nil {
		log.Println("CLIENT fn:getCreated os.Create error:", err.Error())
	}

	defer c.Request.Body.Close()
	defer fp.Close()

	l, err := io.Copy(fp, c.Request.Body)

	if err != nil {
		log.Println("CLIENT fn:getCreated io.Copy error:", err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	} else {
		log.Println(path, "file created and writen", l, "bytes")
		(*stg).MyFiles.Add(path)
		c.AbortWithStatus(http.StatusOK)
		return
	}
}

func (stg *Settings) GetModified(c *gin.Context) {
	path := c.GetHeader("path")
	path = stg.AbsPath(path)

	if path == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	fp, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		log.Println("CLIENT fn:getModified os.OpenFile error:", err.Error())
	}

	defer fp.Close()
	defer c.Request.Body.Close()

	l, err := io.Copy(fp, c.Request.Body)

	if err != nil {
		log.Println("CLIENT fn:getModified io.Copy error:", err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	} else {
		log.Println(path, "file modified and writen", l, "bytes")
		(*stg).MyFiles.Add(path)
		c.AbortWithStatus(http.StatusOK)
		return
	}
}

func (stg *Settings) GetDeleted(c *gin.Context) {
	path := c.GetHeader("path")
	path = stg.AbsPath(path)

	if path == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err := os.Remove(path)
	if err != nil {
		log.Println(err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	} else {
		log.Println(path, "File deleted")
		(*stg).MyFiles.Remove(path)
		c.AbortWithStatus(http.StatusOK)
		return
	}

}
func (stg *Settings) GetAck(c *gin.Context) {
	c.JSON(http.StatusOK, (*stg).MyFiles.GetAllRelative((*stg).WatchPath))
}
