package server

import (
	"encoding/json"
	"fmt"
	"strings"
)

func (stg *Settings) SendWelcome(host string, port int, lport int, auth string) {

	someString, _ := json.Marshal(map[string]interface{}{"ip": "", "port": lport})
	fp := strings.NewReader(string(someString))
	stg.MakeRequest(
		"POST",
		fmt.Sprintf("http://%s:%d/welcome", host, port),
		dict{},
		fp,
	)
}
