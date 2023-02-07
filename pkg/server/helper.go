package server

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func makeRequest(method, url string, headers dict, body io.Reader) *http.Response {
	client := &http.Client{}

	req, _ := http.NewRequest(method, url, body)

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	resp, _ := client.Do(req)

	return resp
}

func md5sum(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func (stg *Settings) initFiles() {

	walk := func(path string, info os.FileInfo, merr error) error {
		fileInfo, err := os.Stat(path)
		if err != nil {
			return err
		}

		if fileInfo.IsDir() {
			return nil
		}

		(*stg).MyFiles.Add(path)

		return err
	}

	filepath.Walk((*stg).WatchPath, walk)
}

func (stg *Settings) initServer() {

	time.Sleep(5 * time.Second)
	for {
		for _, url := range (*stg).MyFriends.Url() {
			resp := makeRequest(
				"GET",
				fmt.Sprintf("%s/ack", url),
				dict{},
				nil,
			)
			payload := make(fileNhash)
			bdy, _ := io.ReadAll(resp.Body)
			resp.Body.Close()

			json.Unmarshal(bdy, &payload)

			go stg.Sync(
				payload.GetAllAbs((*stg).WatchPath),
				url,
			)
		}

		time.Sleep(time.Duration(stg.SyncTime) * time.Second)
	}
}
func (stg *Settings) initClient() {

	for {
		sendWelcome((*stg).MasterIp.String(), (*stg).MasterPort, (*stg).Port)

		time.Sleep(time.Duration(stg.SyncTime) * time.Second)
	}
}

func (stg *Settings) Sync(fnh dict, url string) {
	sfiles := stg.MyFiles
	cfiles := fnh

	for key, value := range sfiles {
		if cfiles[key] == "" {
			stg.sendCreated(key)
		} else {
			if cfiles[key] != value {
				fmt.Println("semding mod", key)
				stg.sendModified(key)
			}
		}
	}

	for key := range cfiles {
		if sfiles[key] == "" {
			stg.sendDeleted(key)
		}
	}
}

func ensureDir(dirName string) error {
	err := os.Mkdir(dirName, os.ModeDir)
	if err == nil {
		return nil
	}
	if os.IsExist(err) {
		// check that the existing path is a directory
		info, err := os.Stat(dirName)
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return errors.New("path exists but is not a directory")
		}
		return nil
	}
	return err
}
