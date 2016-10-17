package vacuum

import (
	"log"
	"sync"
)

var (
	strings                  = map[string]*String{}
	registeredStringRoutines = map[string]StringRoutine{}
	declaredServices         = map[string]string{}
	stringIDsByService       = map[string]map[string]bool{}
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
	declaredServices[stringID] = serviceName

	stringIDs, ok := stringIDsByService[serviceName]
	if !ok {
		stringIDs = map[string]bool{}
		stringIDsByService[serviceName] = stringIDs
	}
	stringIDs[stringID] = true
}

func isServiceReady(serviceName string) {

}
