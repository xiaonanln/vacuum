package filesystem

import (
	"path/filepath"

	"io/ioutil"

	"encoding/json"

	"github.com/Sirupsen/logrus"
	"github.com/xiaonanln/vacuum/storage"
	"os"
	"encoding/base64"
)

type FileSystemStringStorage struct {
	directory string
}

func encodeStringID(stringID string) string {
	return base64.URLEncoding.EncodeToString([]byte(stringID))
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
	stringSaveFile := filepath.Join(ss.directory, encodeStringID(stringID))
	dataBytes, err := ioutil.ReadFile(stringSaveFile)
	if err != nil {
		return nil, err
	}

	var data interface{}
	err = json.Unmarshal(dataBytes, &data)
	if err != nil{
		return nil, err
	}
	return data, nil
}

func newFileSystemStringStorage(directory string) *FileSystemStringStorage {
	if err := os.MkdirAll(directory, 0644); err != nil{
		logrus.Panic(err)
	}

	return &FileSystemStringStorage{
		directory: directory,
	}
}

func OpenDirectory(directory string) storage.StringStorage {
	return newFileSystemStringStorage(directory)
}
