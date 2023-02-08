package main

import (
	"net"
	"os"
	"strconv"

	"github.com/mayankfawkes/filesync/pkg/server"
)

func getenv(key string, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func main() {

	watchPath := getenv("WATCH_PATH", "")

	if watchPath == "" {
		panic("WATCH_PATH not found.")
	}

	var nodeType bool

	if getenv("NODE", "MASTER") == "MASTER" {
		nodeType = true
	}

	port, err := strconv.Atoi(getenv("PORT", "8000"))
	if err != nil {
		panic(err.Error())
	}

	mip := getenv("MASTER_IP", "")
	if (!nodeType) && (mip == "") {
		panic("Master server IP not found")
	}

	mport, err := strconv.Atoi(getenv("MASTER_PORT", "8000"))
	if err != nil {
		panic(err.Error())
	}

	syncTime, err := strconv.Atoi(getenv("SYNC_TIME", "300"))
	if err != nil {
		panic(err.Error())
	}

	stg := server.Settings{
		WatchPath:  watchPath,
		Server:     nodeType,
		Port:       port,
		MasterIp:   net.ParseIP(mip),
		MasterPort: mport,
		SyncTime:   syncTime,
	}
	server.Server(stg)

}
