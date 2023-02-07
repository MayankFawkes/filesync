package server

import (
	"encoding/json"
	"fmt"
	"strings"
)

func sendWelcome(host string, port int, lport int) {

	someString, _ := json.Marshal(map[string]interface{}{"ip": "", "port": lport})
	fp := strings.NewReader(string(someString))
	makeRequest(
		"POST",
		fmt.Sprintf("http://%s:%d/welcome", host, port),
		dict{},
		fp,
	)
}
