package vacuum

import (
	"math/rand"
	"sync"

	log "github.com/Sirupsen/logrus"

	. "github.com/xiaonanln/vacuum/common"
)

var (
	stringsLock sync.RWMutex
	strings     = map[string]*String{}

	registeredStringRoutines = map[string]StringRoutine{}

	stringIDsByServiceLock sync.RWMutex
	stringIDListByService  = map[string]StringList{}
	stringIDsByService     = map[string]StringSet{}
)

func putString(s *String) {
	stringsLock.Lock()
	strings[s.ID] = s
	stringsLock.Unlock()
}

func getString(stringID string) (s *String) {
	stringsLock.RLock()
	s = strings[stringID]
	stringsLock.RUnlock()
	return
}

func delString(stringID string) {
	stringsLock.Lock()
	delete(strings, stringID)
	stringsLock.Unlock()
	return
}

func RegisterString(name string, routine StringRoutine) {
	if registeredStringRoutines[name] != nil {
		log.Panicf("String routine of name %s is already registered", name)
	}

	registeredStringRoutines[name] = routine
	log.Infof("String routine registered: %s", name)
}

func getStringRoutine(name string) StringRoutine {
	return registeredStringRoutines[name]
}

func declareService(stringID string, serviceName string) {
	stringIDsByServiceLock.Lock()

	stringIDs, ok := stringIDsByService[serviceName]
	if ok { // add stringID to service
		stringIDs.Add(stringID)
		sl := stringIDListByService[serviceName]
		sl.Append(stringID)
		stringIDListByService[serviceName] = sl // maintain the stringID list

	} else { // found new service
		stringIDsByService[serviceName] = StringSet{}
		stringIDsByService[serviceName].Add(stringID)

		stringIDListByService[serviceName] = StringList{stringID}
	}

	stringIDsByServiceLock.Unlock()
}

// Undeclare all service of the specified string
func undeclareServicesOfString(stringID string) {
	stringIDsByServiceLock.Lock()

	for serviceName, stringIDs := range stringIDsByService {
		log.Debugf("undeclareServicesOfString: checking service %s, stringIDs %v, contains %v", serviceName, stringIDs, stringIDs.Contains(stringID))
		if stringIDs.Contains(stringID) {
			// the string declared this service, remove it
			log.Debugf("Undeclaring service %s of String %s!", serviceName, stringID)
			stringIDs.Remove(stringID)
			sl := stringIDListByService[serviceName]
			sl.Remove(stringID)
			stringIDListByService[serviceName] = sl
		}
	}

	stringIDsByServiceLock.Unlock()
}

func chooseServiceString(serviceName string) (stringID string) {
	stringIDsByServiceLock.RLock()

	stringIDs, ok := stringIDListByService[serviceName]
	if ok {
		stringID = stringIDs[rand.Intn(len(stringIDs))]
	} else {
		log.Panicf("chooseServiceString: get service string failed: %s", serviceName)
		stringID = ""
	}

	stringIDsByServiceLock.RUnlock()
	return
}

func getAllServiceStringIDs(serviceName string) (stringIDs []string) {
	stringIDsByServiceLock.RLock()

	_stringIDs := stringIDListByService[serviceName]
	stringIDs = make([]string, len(_stringIDs))
	copy(stringIDs, _stringIDs)

	stringIDsByServiceLock.RUnlock()
	return
}

func GetServiceProviderCount(serviceName string) int {
	stringIDsByServiceLock.RLock()
	defer stringIDsByServiceLock.RUnlock()

	strs, ok := stringIDsByService[serviceName]
	if ok {
		return len(strs)
	} else {
		return 0
	}
}
