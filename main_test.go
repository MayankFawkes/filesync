package main_test

import (
	"net"
	"testing"
	"time"

	"github.com/mayankfawkes/filesync/pkg/server"
)

func TestMain(t *testing.T) {

	stg := server.Settings{
		WatchPath: "test/loc1",
		Server:    true,
		Port:      8000,
		SyncTime:  30,
		Auth:      "123456789",
		Logging:   server.Logging{Enable: true, Debug: true},
	}

	cstg := server.Settings{
		WatchPath:  "test/loc2",
		Server:     false,
		Port:       8001,
		MasterIp:   net.ParseIP("127.0.0.1"),
		MasterPort: 8000,
		SyncTime:   30,
		Auth:       "123456789",
		Logging:    server.Logging{Enable: true, Debug: true},
	}

	go server.Server(stg)

	time.Sleep(5 * time.Second)

	server.Server(cstg)

}
