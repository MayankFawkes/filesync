package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func (stg *Settings) SendCreated(path string) {
	fp, _ := os.Open(path)
	defer fp.Close()

	(*stg).MyFiles.Add(path)

	log.Println("Created event trigger", path)

	headers := dict{
		"Path": stg.RelativePath(path),
	}

	for _, url := range (*stg).MyFriends.Url() {
		stg.SendSingleCreated(url, headers, fp)
	}
}

func (stg *Settings) SendModified(path string) {
	fp, err := os.Open(path)

	if err != nil {
		log.Println("Error mod", err)
	}

	defer fp.Close()

	stg.MyFiles.Add(path)

	log.Println("Modified event trigger", path)

	headers := dict{
		"Path": stg.RelativePath(path),
	}

	for _, url := range stg.MyFriends.Url() {
		stg.SendSingleModified(url, headers, fp)
	}
}

func (stg *Settings) SendDeleted(path string) {
	log.Println("Deleted event trigger", path)

	stg.MyFiles.Remove(path)

	headers := dict{
		"Path": stg.RelativePath(path),
	}

	for _, url := range (*stg).MyFriends.Url() {
		stg.SendSingleDelete(url, headers)
	}
}

// --------------------------------------------------

func (stg *Settings) SendSingleCreated(url string, headers dict, fp *os.File) *http.Response {
	return stg.MakeRequest(
		"POST",
		fmt.Sprintf("%s/created", url),
		headers,
		fp,
	)
}

func (stg *Settings) SendSingleModified(url string, headers dict, fp *os.File) *http.Response {
	return stg.MakeRequest(
		"POST",
		fmt.Sprintf("%s/modified", url),
		headers,
		fp,
	)
}

func (stg *Settings) SendSingleDelete(url string, headers dict) *http.Response {
	return stg.MakeRequest(
		"DELETE",
		fmt.Sprintf("%s/deleted", url),
		headers,
		nil,
	)
}
