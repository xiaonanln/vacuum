package filesystem

import (
	"path/filepath"

	"io/ioutil"

	"encoding/json"

	"github.com/Sirupsen/logrus"
	"github.com/xiaonanln/vacuum/storage"
)

type FileSystemStringStorage struct {
	directory string
}

func encodeStringID(stringID string) string {
	return stringID
}

func (ss *FileSystemStringStorage) Write(stringID string, data interface{}) error {
	stringSaveFile := filepath.Join(ss.directory, encodeStringID(stringID))
	dataBytes, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	logrus.Debugf("Saving to file %s: %s", stringSaveFile, string(dataBytes))
	return ioutil.WriteFile(stringSaveFile, dataBytes, 0644)
}

func (ss *FileSystemStringStorage) Read(stringID string) (interface{}, error) {
	return nil, nil
}

func newFileSystemStringStorage(directory string) *FileSystemStringStorage {
	return &FileSystemStringStorage{
		directory: directory,
	}
}

func OpenDirectory(directory string) storage.StringStorage {
	return newFileSystemStringStorage(directory)
}
