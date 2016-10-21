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
	stringIDListByService  = map[string][]string{}
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
	defer stringIDsByServiceLock.Unlock()

	stringIDs, ok := stringIDsByService[serviceName]
	if ok { // add stringID to service
		stringIDs.Add(stringID)
		stringIDListByService[serviceName] = append(stringIDListByService[serviceName], stringID) // maintain the stringID list

	} else { // found new service
		stringIDsByService[serviceName] = StringSet{}
		stringIDsByService[serviceName].Add(stringID)

		stringIDListByService[serviceName] = []string{stringID}
	}
}

func chooseServiceString(serviceName string) string {
	stringIDsByServiceLock.RLock()
	defer stringIDsByServiceLock.RUnlock()

	stringIDs, ok := stringIDListByService[serviceName]
	if ok {
		return stringIDs[rand.Intn(len(stringIDs))]
	} else {
		log.Panicf("chooseServiceString: get service string failed: %s", serviceName)
		return ""
	}
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
