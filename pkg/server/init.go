package server

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"time"
)

func (stg *Settings) InitFiles() {
	stg.LogDebug("Initiating files started")
	walk := func(path string, info os.FileInfo, merr error) error {
		fileInfo, err := os.Stat(path)
		if err != nil {
			return err
		}

		if fileInfo.IsDir() {
			return nil
		}

		stg.MyFiles.Add(path)

		return err
	}

	filepath.Walk(stg.WatchPath, walk)
	stg.LogDebug("Initiating files completed found", len(stg.MyFiles), "files")
}

func (stg *Settings) InitServer() {

	stg.LogDebug("Initiating Ack Requests")

	time.Sleep(5 * time.Second)
	for {
		for _, frnd := range stg.MyFriends {
			resp, err := stg.MakeRequest(
				&requestPayload{
					Method:  "GET",
					Friend:  frnd,
					Path:    "/ack",
					Headers: dict{"Authorization": stg.Auth},
					Body:    nil,
				},
			)

			if err != nil {
				stg.LogError(err.Error())
				continue
			}

			payload := make(fileNhash)
			bdy, err := io.ReadAll(resp.Body)

			if err != nil {
				stg.LogError(err.Error())
				continue
			}

			resp.Body.Close()

			err = json.Unmarshal(bdy, &payload)
			if err != nil {
				stg.LogError(err.Error())
				continue
			}

			go stg.Sync(
				payload.GetAllAbs(stg.WatchPath),
				frnd,
			)
		}

		time.Sleep(time.Duration(stg.SyncTime) * time.Second)
	}
}

func (stg *Settings) InitClient() {

	stg.LogDebug("Initiating Welcome Requests")
	for {
		stg.SendWelcome()
		time.Sleep(time.Duration(stg.SyncTime) * time.Second)
	}
}
