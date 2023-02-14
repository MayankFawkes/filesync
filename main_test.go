package main_test

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"path"
	"path/filepath"
	"testing"
	"time"

	"github.com/mayankfawkes/filesync/pkg/server"
)

func TestMain(t *testing.T) {

	log := false

	masterDIR := "test/loc1"
	files := []string{
		fmt.Sprintf("%s/file1.txt", masterDIR),
		fmt.Sprintf("%s/file2.txt", masterDIR),
		fmt.Sprintf("%s/file3.txt", masterDIR),
		fmt.Sprintf("%s/file4.txt", masterDIR),
		fmt.Sprintf("%s/file5.txt", masterDIR),
	}

	makeFiles(files, t)

	stg := server.Settings{
		WatchPath: masterDIR,
		Server:    true,
		Port:      8000,
		SyncTime:  30,
		Auth:      "123456789",
		Logging:   server.Logging{Enable: log, Debug: log},
	}

	cstg := server.Settings{
		WatchPath:  "test/loc2",
		Server:     false,
		Port:       8001,
		MasterIp:   net.ParseIP("127.0.0.1"),
		MasterPort: 8000,
		SyncTime:   30,
		Auth:       "123456789",
		Logging:    server.Logging{Enable: log, Debug: log},
	}

	go server.Server(&stg)
	go server.Server(&cstg)

	// Let both server do the init things
	time.Sleep(2 * time.Second)

	// Max wait time for testing - 2min
	MAX_WAIT_TIME := time.Now().Unix() + 60 + 60

	for MAX_WAIT_TIME > time.Now().Unix() {
		masterFiles := stg.MyFiles.GetAllRelative(stg.WatchPath)
		slaveFiles := cstg.MyFiles.GetAllRelative(cstg.WatchPath)

		if checkDict(masterFiles, slaveFiles) {
			break
		}

		time.Sleep(1000 * time.Millisecond)
	}
	remFiles()
}

func makeFiles(f []string, t *testing.T) {
	for _, name := range f {
		dir := filepath.Dir(name)
		err := server.EnsureDir(dir)
		if err != nil {
			t.Errorf(err.Error())
		}

		token := make([]byte, 500)
		rand.Read(token)
		r := bytes.NewReader(token)

		fs, err := os.Create(name)
		if err != nil {
			t.Errorf(err.Error())
		}

		io.Copy(fs, r)
	}
}

func remFiles() {
	dir, _ := ioutil.ReadDir("test")
	for _, d := range dir {
		os.RemoveAll(path.Join([]string{"test", d.Name()}...))
	}

}
func checkDict(s server.Dict, c server.Dict) bool {
	if len(s) != len(c) {
		return false
	}

	for k, v := range s {
		if c[k] != v {
			return false
		}
	}
	for k, v := range c {
		if s[k] != v {
			return false
		}
	}

	return true

}
