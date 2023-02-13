package server

import (
	"net/http"
	"os"
)

func (stg *Settings) SendCreated(path string) {
	fp, err := os.Open(path)
	if err != nil {
		stg.LogError(err)
		return
	}

	defer fp.Close()

	stg.MyFiles.Add(path)

	stg.LogDebug("Created event trigger", path)

	headers := dict{
		"Path": stg.RelativePath(path),
	}

	for _, frnd := range stg.MyFriends {
		stg.SendSingleCreated(frnd, headers, fp)
	}
}

func (stg *Settings) SendModified(path string) {
	fp, err := os.Open(path)

	if err != nil {
		stg.LogError(err)
		return
	}

	defer fp.Close()

	stg.MyFiles.Add(path)

	stg.LogDebug("Modified event trigger", path)

	headers := dict{
		"Path": stg.RelativePath(path),
	}

	for _, frnd := range stg.MyFriends {
		stg.SendSingleModified(frnd, headers, fp)
	}
}

func (stg *Settings) SendDeleted(path string) {
	stg.LogDebug("Deleted event trigger", path)

	stg.MyFiles.Remove(path)

	headers := dict{
		"Path": stg.RelativePath(path),
	}

	for _, frnd := range stg.MyFriends {
		stg.SendSingleDelete(frnd, headers)
	}
}

// --------------------------------------------------

func (stg *Settings) SendSingleCreated(frnd friend, headers dict, fp *os.File) (*http.Response, error) {
	return stg.MakeRequest(
		&requestPayload{
			Method:  "POST",
			Friend:  frnd,
			Path:    "/created",
			Headers: headers,
			Body:    fp,
		},
	)
}

func (stg *Settings) SendSingleModified(frnd friend, headers dict, fp *os.File) (*http.Response, error) {
	return stg.MakeRequest(
		&requestPayload{
			Method:  "POST",
			Friend:  frnd,
			Path:    "/modified",
			Headers: headers,
			Body:    fp,
		},
	)
}

func (stg *Settings) SendSingleDelete(frnd friend, headers dict) (*http.Response, error) {
	return stg.MakeRequest(
		&requestPayload{
			Method:  "DELETE",
			Friend:  frnd,
			Path:    "/deleted",
			Headers: headers,
			Body:    nil,
		},
	)
}
