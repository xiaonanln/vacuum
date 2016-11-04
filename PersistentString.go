package vacuum

import "github.com/Sirupsen/logrus"

type PersistentString interface {
	GetPersistentData() map[string]interface{}
	LoadPersistentData(data map[string]interface{})
}

func (s *String) Save() error {
	persistence := s.Persistence()
	if persistence == nil {
		logrus.Panicf("string %s is not persistent", s)
	}
	logrus.Debugf("SAVING %s ...", s)
	data := persistence.GetPersistentData()
	if err := stringStorage.Write(s.ID, data); err != nil {
		logrus.Errorf("Save %s failed: %s", s, err.Error())
		return err
	}
	return nil
}