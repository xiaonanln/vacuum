package vacuum

import (
	"time"

	log "github.com/Sirupsen/logrus"
)

func WaitServiceReady(serviceName string, n int) {
	current := GetServiceProviderCount(serviceName)
	log.Printf("Waiting %s service to be ready: need %d, current %d", serviceName, n, current)
	for current < n {
		time.Sleep(100 * time.Millisecond)
		current = GetServiceProviderCount(serviceName)
	}
	log.WithFields(log.Fields{
		"num": current,
	}).Printf("Service %s is ready now.", serviceName)
}

func interfaceToInt(v interface{}) int {
	n1, ok := v.(uint64)
	if ok {
		return int(n1)
	}

	n2, ok := v.(int64)
	if ok {
		return int(n2)
	}

	return v.(int)
}

func (s *String) ReadInt() int {
	return interfaceToInt(s.Read())
}

func (s *String) ReadIntTuple() []int {
	v := s.Read()
	intTuple, ok := v.([]int)
	if ok {
		return intTuple
	}

	tuple := v.([]interface{})
	ret := make([]int, len(tuple))
	for i, v := range tuple {
		ret[i] = interfaceToInt(v)
	}
	return ret
}
