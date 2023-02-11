package server

import (
	"encoding/json"
	"fmt"
	"net"
	"path/filepath"
)

type friend struct {
	Ip   net.IP `json:"ip"`
	Port int    `json:"port"`
}

type friends map[string]friend

func (fs *friends) Add(f friend) *friends {
	(*fs)[f.Ip.String()] = f
	return fs
}

func (fs *friends) Remove(f string) *friends {
	delete((*fs), f)
	return fs
}

func (fs *friends) Url() []string {
	lst := []string{}
	for _, value := range *fs {
		lst = append(lst, fmt.Sprintf("http://%s:%d", value.Ip.String(), value.Port))
	}
	return lst
}

type dict map[string]string

func (d *dict) Json() ([]byte, error) {
	jsonString, err := json.Marshal(*d)
	return jsonString, err
}

type fileNhash dict

func (d *fileNhash) Add(path string) {
	sum, _ := md5sum(path)
	(*d)[path] = sum
}

func (d *fileNhash) Remove(path string) {
	delete((*d), path)
}

func (d *fileNhash) GetAllRelative(basepath string) dict {
	p := make(dict)

	for targpath, hash := range *d {
		relPath, _ := filepath.Rel(basepath, targpath)
		p[relPath] = hash
	}

	return p
}

func (d *fileNhash) GetAllAbs(basepath string) dict {
	p := make(dict)

	for targpath, hash := range *d {
		relPath := filepath.Join(basepath, targpath)
		p[relPath] = hash
	}

	return p
}

type Settings struct {
	MyFriends  friends
	MyFiles    fileNhash
	Logging    bool
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
	ensureDir(dir)

	return path

}
