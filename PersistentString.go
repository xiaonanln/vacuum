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
	GetPersistentData() map[string]interface{}
	LoadPersistentData(data map[string]interface{})
}

func (s *String) Save() {
	persistence := s.persistence
	if persistence == nil {
		vlog.TraceError("string %s is not persistent", s)
		return
	}

	logrus.Debugf("Saving %s ...", s)

	data := persistence.GetPersistentData()
	if err := stringStorage.Write(s.Name, s.ID, data); err != nil {
		vlog.TraceError("Save %s failed: %s", err.Error())
	}
}

func AssurePersistent(_ PersistentString) {

}
