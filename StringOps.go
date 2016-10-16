package vacuum

import "github.com/xiaonanln/vacuum/vacuum_server/dispatcher_client"

func CreateString(name string) {
	dispatcher_client.CreateString(name)
}

func createString(name string) {
	routine := getStringRoutine(name)
	s := newString(routine)
	putString(s)
	go s.routine(s)
}
