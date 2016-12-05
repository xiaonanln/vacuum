package vacuum

import (
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/xiaonanln/vacuum/vlog"
)

const (
	SAVE_ERROR_RETRY_DELAY = 3 * time.Second
)

type PersistentString interface {
	IsPersistent() bool
	GetPersistentData() map[string]interface{}
	LoadPersistentData(data map[string]interface{})
}

func (s *String) Save() {
	is := s.I

	if !is.IsPersistent() {
		vlog.Debug("String %s is NOT persistent, save ignored.", s)
		return
	}

	logrus.Debugf("Saving %s ...", s)

	data := is.GetPersistentData()
	if err := stringStorage.Write(s.Name, s.ID, data); err != nil {
		vlog.TraceError("Save %s failed: %s", err.Error())
	}
}

func AssurePersistent(_ PersistentString) {

}
