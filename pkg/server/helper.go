package server

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
)

func checkError(err error) string {
	var errStr string
	if err, ok := err.(*url.Error); ok {
		if err, ok := err.Err.(*net.OpError); ok {
			if _, ok := err.Err.(*net.DNSError); ok {
				errStr = "Timeout Error"
			}
		}
	}

	if e, ok := err.(*url.Error); ok && e.Timeout() {
		errStr = "DNS Error"
	}

	return errStr
}

func (stg *Settings) MakeRequest(rp *requestPayload) (*http.Response, error) {
	client := &http.Client{}

	stg.LogDebug("Sending", rp.Method, rp.Path)

	req, _ := http.NewRequest(rp.Method, fmt.Sprintf("http://%s%s", rp.Friend.Host(), rp.Path), rp.Body)

	req.Header.Add("Authorization", stg.Auth)
	for key, value := range rp.Headers {
		req.Header.Add(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		errStr := checkError(err)

		stg.LogError(rp.Method, fmt.Sprintf("%s%s", rp.Friend.Host(), rp.Path), err.Error(), errStr)
		stg.LogInfo("Removing", rp.Friend.Host())

		stg.MyFriends.Remove(rp.Friend.Ip.String())
		return nil, err
	}

	stg.LogInfo(resp.StatusCode, rp.Method, fmt.Sprintf("%s%s", rp.Friend.Host(), rp.Path))

	return resp, nil
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

func (stg *Settings) Sync(fnh Dict, frnd friend) {
	sfiles := stg.MyFiles
	cfiles := fnh

	stg.LogInfo("Sync Started for", frnd.Host())

	for key, value := range sfiles.m {
		if cfiles[key] == "" {
			stg.SendCreated(key)
		} else {
			if cfiles[key] != value {
				stg.LogDebug("Sync Sending Modification", frnd.Host())
				stg.SendModified(key)
			}
		}
	}

	for key := range cfiles {
		if !sfiles.Check(key) {
			stg.LogDebug("Sync Sending Deletion", frnd.Host())
			stg.SendDeleted(key)
		}
	}
}

func EnsureDir(dirName string) error {
	err := os.MkdirAll(dirName, os.ModeDir)
	if err == nil {
		return err
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
