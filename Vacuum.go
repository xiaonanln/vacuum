package vacuum

import "log"

type Vacuum struct {
	strs map[string]*String
}

var (
	vacuumInstance = Vacuum{
		strs: map[string]*String{},
	}

	registeredStringRoutines = map[string]StringRoutine{}
)

func putString(s *String) {
	vacuumInstance.strs[s.ID] = s
}

func getString(sid string) *String {
	return vacuumInstance.strs[sid]
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
