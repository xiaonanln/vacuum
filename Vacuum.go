package vacuum

import (
	"log"
	"math/rand"
	"sync"

	. "github.com/xiaonanln/vacuum/common"
)

var (
	strings                  = map[string]*String{}
	registeredStringRoutines = map[string]StringRoutine{}
	stringIDListByService    = map[string][]string{}
	stringIDsByService       = map[string]StringSet{}
	lock                     sync.RWMutex
)

func putString(s *String) {
	strings[s.ID] = s
}

func getString(sid string) *String {
	return strings[sid]
}

func RegisterString(name string, routine StringRoutine) {
	if registeredStringRoutines[name] != nil {
		log.Panicf("String routine of name %s is already registered", name)
	}

	registeredStringRoutines[name] = routine
	log.Printf("String routine registered: %s", name)
}

func getStringRoutine(name string) StringRoutine {
	return registeredStringRoutines[name]
}

func declareService(stringID string, serviceName string) {
	lock.Lock()
	defer lock.Unlock()

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
	lock.RLock()
	defer lock.RUnlock()

	stringIDs, ok := stringIDListByService[serviceName]
	if ok {
		return stringIDs[rand.Intn(len(stringIDs))]
	} else {
		return ""
	}
}

func GetServiceProviderCount(serviceName string) int {
	lock.RLock()
	defer lock.RUnlock()

	strs, ok := stringIDsByService[serviceName]
	if ok {
		return len(strs)
	} else {
		return 0
	}
}
