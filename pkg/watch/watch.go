package watch

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

const (
	CREATED = 1 << iota
	DELETED
	MODIFIED

	CDM = CREATED | DELETED | MODIFIED
)

type Response struct {
	Status uint32
	Path   string
}

type typeStorage map[string]fs.FileInfo

func walkPath(storage typeStorage, ch chan Response) func(path string, f fs.FileInfo, merr error) error {

	return func(path string, f fs.FileInfo, merr error) error {
		fileInfo, err := os.Stat(path)
		if err != nil {
			return err
		}

		if fileInfo.IsDir() {
			return nil
		}

		pastFileInfo := storage[path]

		if pastFileInfo == nil {
			ch <- Response{
				Status: CREATED,
				Path:   path,
			}

			storage[path] = fileInfo
			pastFileInfo = fileInfo
			return nil
		}

		if pastFileInfo.Size() != fileInfo.Size() || pastFileInfo.ModTime() != fileInfo.ModTime() {
			ch <- Response{
				Status: MODIFIED,
				Path:   path,
			}
			storage[path] = fileInfo
		}
		return nil
	}

}

func checkRemoved(storage typeStorage, ch chan Response) func() {
	return func() {
		for path := range storage {
			if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
				delete(storage, path)
				ch <- Response{
					Status: DELETED,
					Path:   path,
				}
			}
		}
	}
}

func WatchPath(root string, ch chan Response) {

	var storage = make(typeStorage)

	fnwalk := walkPath(storage, ch)
	fnremove := checkRemoved(storage, ch)

	// Init
	filepath.Walk(root, func(path string, f fs.FileInfo, merr error) error {
		fileInfo, err := os.Stat(path)
		if err != nil {
			return err
		}

		if fileInfo.IsDir() {
			return nil
		}

		storage[path] = fileInfo
		return nil

	})

	for {
		filepath.Walk(root, fnwalk)
		fnremove()
		time.Sleep(100 * time.Millisecond)
	}

}
