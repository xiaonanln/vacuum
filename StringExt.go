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

func (s *String) ReadInt() int {
	msg := s.Read()
	n1, ok := msg.(uint64)
	if ok {
		return int(n1)
	}

	n2, ok := msg.(int64)
	if ok {
		return int(n2)
	}

	return msg.(int)
}
