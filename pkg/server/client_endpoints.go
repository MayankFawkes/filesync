package server

import (
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func (stg *Settings) GetCreated(c *gin.Context) {
	path := c.GetHeader("path")
	path = stg.AbsPath(path)

	stg.LogInfo("Creating file", path)

	if path == "" {
		stg.LogWarning("Path is empty")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	fp, err := os.Create(path)
	if err != nil {
		stg.LogError(err.Error())
	}

	defer c.Request.Body.Close()
	defer fp.Close()

	l, err := io.Copy(fp, c.Request.Body)

	if err != nil {
		stg.LogError(err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	} else {
		stg.LogDebug("File Modified and Writen", l, "bytes loc:", path)
		(*stg).MyFiles.Add(path)
		c.AbortWithStatus(http.StatusOK)
		return
	}
}

func (stg *Settings) GetModified(c *gin.Context) {
	path := c.GetHeader("path")
	path = stg.AbsPath(path)

	stg.LogInfo("Modifying  file", path)

	if path == "" {
		stg.LogWarning("Path is empty")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	fp, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		stg.LogError(err.Error())
	}

	defer fp.Close()
	defer c.Request.Body.Close()

	l, err := io.Copy(fp, c.Request.Body)

	if err != nil {
		stg.LogError(err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	} else {
		stg.LogDebug("File Modified and Writen", l, "bytes loc:", path)
		(*stg).MyFiles.Add(path)
		c.AbortWithStatus(http.StatusOK)
		return
	}
}

func (stg *Settings) GetDeleted(c *gin.Context) {
	path := c.GetHeader("path")
	path = stg.AbsPath(path)

	stg.LogInfo("Deleting  file", path)

	if path == "" {
		stg.LogWarning("Path is empty")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err := os.Remove(path)
	if err != nil {
		stg.LogError(err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	} else {
		stg.LogDebug("File Deleted", path)
		(*stg).MyFiles.Remove(path)
		c.AbortWithStatus(http.StatusOK)
		return
	}

}
func (stg *Settings) GetAck(c *gin.Context) {
	c.JSON(http.StatusOK, stg.MyFiles.GetAllRelative(stg.WatchPath))
}
