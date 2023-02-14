package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"path/filepath"
	"sync"
)

type friend struct {
	Ip   net.IP `json:"ip"`
	Port int    `json:"port"`
}

func (fs *friend) Host() string {
	return fmt.Sprintf("%s:%d", fs.Ip.String(), fs.Port)
}

// type friends map[string]friend

type friends struct {
	sync.RWMutex
	m map[string]friend
}

func (fs *friends) Add(f friend) *friends {
	fs.Lock()
	defer fs.Unlock()
	fs.m[f.Ip.String()] = f
	return fs
}

func (fs *friends) Remove(f string) *friends {
	fs.Lock()
	defer fs.Unlock()
	delete(fs.m, f)
	return fs
}

func (fs *friends) Url() []string {
	lst := []string{}
	fs.RLock()
	defer fs.RUnlock()
	for _, value := range fs.m {
		lst = append(lst, value.Host())
	}
	return lst
}

type Dict map[string]string

func (d *Dict) Json() ([]byte, error) {
	jsonString, err := json.Marshal(*d)
	return jsonString, err
}

// type fileNhash dict
type fileNhash struct {
	sync.RWMutex
	m Dict
}

func (d *fileNhash) Add(path string) {
	sum, _ := md5sum(path)
	d.Lock()
	defer d.Unlock()
	d.m[path] = sum
}

func (d *fileNhash) Remove(path string) {
	d.Lock()
	defer d.Unlock()
	delete(d.m, path)
}

func (d *fileNhash) Check(path string) bool {
	d.RLock()
	defer d.RUnlock()
	s := d.m[path]
	return len(s) != 0
}

func (d *fileNhash) Compare(key string, value string) bool {
	d.RLock()
	defer d.RUnlock()
	return d.m[key] == value
}

func (d *fileNhash) GetAllRelative(basepath string) Dict {
	p := make(Dict)

	d.RLock()
	defer d.RUnlock()

	for targpath, hash := range d.m {
		relPath, _ := filepath.Rel(basepath, targpath)
		p[relPath] = hash
	}

	return p
}

func (d *fileNhash) GetAllAbs(basepath string) Dict {
	p := make(Dict)

	d.RLock()
	defer d.RUnlock()

	for targpath, hash := range d.m {
		relPath := filepath.Join(basepath, targpath)
		p[relPath] = hash
	}

	return p
}

type Logging struct {
	Enable bool
	Debug  bool
}

type Settings struct {
	MyFriends  *friends
	MyFiles    *fileNhash
	Logging    Logging
	WatchPath  string
	Server     bool
	Port       int
	MasterIp   net.IP
	MasterPort int
	SyncTime   int
	Auth       string
}

func (res *Settings) RelativePath(targpath string) string {
	relPath, _ := filepath.Rel((*res).WatchPath, targpath)
	return relPath
}

func (res *Settings) AbsPath(targpath string) string {
	path := filepath.Join((*res).WatchPath, targpath)

	dir := filepath.Dir(path)
	EnsureDir(dir)

	return path

}

type requestPayload struct {
	Method  string
	Friend  friend
	Path    string
	Headers Dict
	Body    io.Reader
}
