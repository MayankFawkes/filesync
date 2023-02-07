package main

import (
	"net"
	"time"

	"github.com/mayankfawkes/filesync/pkg/server"
)

// func getenv(key string, fallback string) string {
// 	value := os.Getenv(key)
// 	if len(value) == 0 {
// 		return fallback
// 	}
// 	return value
// }

func main() {

	// watchPath := getenv("WATCH_PATH", "")

	// if watchPath == "" {
	// 	panic("WATCH_PATH not found.")
	// }

	// var nodeType bool

	// if getenv("NODE", "MASTER") == "MASTER" {
	// 	nodeType = true
	// }

	// port, err := strconv.Atoi(getenv("PORT", "8000"))
	// if err != nil {
	// 	panic(err.Error())
	// }

	// mip := getenv("MASTER_IP", "")
	// if (!nodeType) && (mip == "") {
	// 	panic("Master server IP not found")
	// }

	// mport, err := strconv.Atoi(getenv("MASTER_PORT", "8000"))
	// if err != nil {
	// 	panic(err.Error())
	// }

	// ackSyncTime, err := strconv.Atoi(getenv("ACK_SYNC_TIME", "300"))
	// if err != nil {
	// 	panic(err.Error())
	// }

	// stg := server.Settings{
	// 	WatchPath:   watchPath,
	// 	Server:      nodeType,
	// 	Port:        port,
	// 	MasterIp:    net.ParseIP(mip),
	// 	MasterPort:  mport,
	// 	AckSyncTime: ackSyncTime,
	// }
	// server.Server(stg)

	stg := server.Settings{
		WatchPath: "test/loc1/",
		Server:    true,
		Port:      8000,
		SyncTime:  30,
	}

	cstg := server.Settings{
		WatchPath:  "test/loc2",
		Server:     false,
		Port:       8001,
		MasterIp:   net.ParseIP("127.0.0.1"),
		MasterPort: 8000,
		SyncTime:   30,
	}

	go server.Server(stg)

	time.Sleep(1 * time.Second)

	server.Server(cstg)

}
