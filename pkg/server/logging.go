package server

import (
	"fmt"

	"github.com/golang/glog"
)

func logPrefix(skip int, log string, server bool) string {

	// getLast := func(s string) string {
	// 	spa := strings.Split(s, "/")
	// 	return spa[len(spa)-1]
	// }

	// t := time.Now()
	// dnt := t.Format("2006-01-02 15:04:05")

	var serverType string

	if server {
		serverType = "MASTER"
	} else {
		serverType = "SLAVE"
	}

	return fmt.Sprintf("%s:%s - ", serverType, log)

	// pc, file, line, ok := runtime.Caller(skip)
	// if ok {
	// 	return fmt.Sprintf("[%s - %s/%s:#%d %s]", dnt, getLast(file), getLast(runtime.FuncForPC(pc).Name()), line, serverType)
	// } else {
	// 	return ""
	// }
}

func (stg *Settings) LogDebug(args ...any) {
	if !stg.Logging.Enable {
		return
	}

	clientType := make([]interface{}, 1)
	clientType[0] = logPrefix(2, "DEB", stg.Server)

	fargs := append(clientType, args...)
	fstr := fmt.Sprintln(fargs...)

	if glog.V(2) {
		glog.InfoDepth(1, fstr)
	}
}

func (stg *Settings) LogInfo(args ...any) {
	if !stg.Logging.Enable {
		return
	}

	clientType := make([]interface{}, 1)
	clientType[0] = logPrefix(2, "INF", stg.Server)

	fargs := append(clientType, args...)
	fstr := fmt.Sprintln(fargs...)

	glog.InfoDepth(1, fstr)
}

func (stg *Settings) LogWarning(args ...any) {
	if !stg.Logging.Enable {
		return
	}
	clientType := make([]interface{}, 1)
	clientType[0] = logPrefix(2, "WAR", stg.Server)

	fargs := append(clientType, args...)
	fstr := fmt.Sprintln(fargs...)

	glog.WarningDepth(1, fstr)
}

func (stg *Settings) LogError(args ...any) {
	if !stg.Logging.Enable {
		return
	}

	clientType := make([]interface{}, 1)
	clientType[0] = logPrefix(2, "ERR", stg.Server)

	fargs := append(clientType, args...)
	fstr := fmt.Sprintln(fargs...)

	glog.ErrorDepth(1, fstr)
}
