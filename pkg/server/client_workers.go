package server

import (
	"encoding/json"
	"strings"
)

func (stg *Settings) SendWelcome() {

	someString, _ := json.Marshal(map[string]interface{}{"ip": "", "port": stg.Port})
	fp := strings.NewReader(string(someString))

	stg.LogDebug("Processing Welcome Request")

	stg.MakeRequest(
		&requestPayload{
			Method:  "POST",
			Friend:  friend{Ip: stg.MasterIp, Port: stg.MasterPort},
			Path:    "/welcome",
			Headers: dict{},
			Body:    fp,
		},
	)
	stg.LogDebug("Welcome Request Sent")
}
