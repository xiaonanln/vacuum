package vacuum

import (
	"math/rand"
	"sync"

	. "github.com/xiaonanln/vacuum/common"
	"github.com/xiaonanln/vacuum/vlog"
)

var (
	stringsLock sync.RWMutex
	strings     = map[string]*String{}

	registeredStringDelegateNewers = map[string]StringDelegateMaker{}

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

func popString(stringID string) (s *String) {
	stringsLock.Lock()
	s = strings[stringID]
	delete(strings, stringID)
	stringsLock.Unlock()
	return
}

func RegisterString(name string, newer StringDelegateMaker) {
	if registeredStringDelegateNewers[name] != nil {
		vlog.Panicf("String delegate newer of name %s is already registered", name)
	}

	registeredStringDelegateNewers[name] = newer
	vlog.Info("String delegate newer registered: %s", name)
}

func GetLocalString(stringID string) *String {
	return getString(stringID)
}

func getStringDelegateMaker(name string) StringDelegateMaker {
	return registeredStringDelegateNewers[name]
}

func declareService(stringID string, serviceName string) {
	stringIDsByServiceLock.Lock()

	stringIDs, ok := stringIDsByService[serviceName]
	if ok { // add stringID to service
		stringIDs.Add(stringID)
		sl := stringIDListByService[serviceName]
		if sl.Find(stringID) == -1 {
			sl.Append(stringID)
		}
		stringIDListByService[serviceName] = sl // maintain the stringID list

	} else { // found new service
		stringIDsByService[serviceName] = StringSet{}
		stringIDsByService[serviceName].Add(stringID)

		stringIDListByService[serviceName] = StringList{stringID}
	}

	vlog.Info("declareService: %s => %s => %v", stringID, serviceName, stringIDListByService[serviceName])
	stringIDsByServiceLock.Unlock()
}

// Undeclare all service of the specified string
func undeclareServicesOfString(stringID string) {
	stringIDsByServiceLock.Lock()

	for serviceName, stringIDs := range stringIDsByService {
		vlog.Debug("undeclareServicesOfString: checking service %s, stringIDs %v, contains %v", serviceName, stringIDs, stringIDs.Contains(stringID))
		if stringIDs.Contains(stringID) {
			// the string declared this service, remove it
			stringIDs.Remove(stringID)
			sl := stringIDListByService[serviceName]
			sl.Remove(stringID)
			stringIDListByService[serviceName] = sl

			vlog.Debug("undeclareServicesOfString %s: %s => %v", serviceName, stringID, stringIDListByService[serviceName])
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
		vlog.Panicf("chooseServiceString: get service string failed: %s", serviceName)
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
