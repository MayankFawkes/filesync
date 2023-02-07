package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func (stg *Settings) sendCreated(path string) {
	fp, _ := os.Open(path)
	defer fp.Close()

	(*stg).MyFiles.Add(path)

	log.Println("Created event trigger", path)

	for _, url := range (*stg).MyFriends.Url() {
		sendSingleCreated(url, stg.RelativePath(path), fp)
	}
}

func (stg *Settings) sendModified(path string) {
	fp, err := os.Open(path)

	if err != nil {
		log.Println("Error mod", err)
	}

	defer fp.Close()

	(*stg).MyFiles.Add(path)

	log.Println("Modified event trigger", path)

	for _, url := range (*stg).MyFriends.Url() {
		sendSingleModified(url, stg.RelativePath(path), fp)
	}
}

func (stg *Settings) sendDeleted(path string) {
	log.Println("Deleted event trigger", path)

	(*stg).MyFiles.Remove(path)

	for _, url := range (*stg).MyFriends.Url() {
		sendSingleDelete(url, stg.RelativePath(path))
	}
}

// --------------------------------------------------

func sendSingleCreated(url string, path string, fp *os.File) *http.Response {
	return makeRequest(
		"POST",
		fmt.Sprintf("%s/created", url),
		dict{"path": path},
		fp,
	)
}

func sendSingleModified(url string, path string, fp *os.File) *http.Response {
	return makeRequest(
		"POST",
		fmt.Sprintf("%s/modified", url),
		dict{"path": path},
		fp,
	)
}

func sendSingleDelete(url string, path string) *http.Response {
	return makeRequest(
		"DELETE",
		fmt.Sprintf("%s/deleted", url),
		dict{"path": path},
		nil,
	)
}
