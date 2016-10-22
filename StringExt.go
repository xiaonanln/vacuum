package vacuum

import (
	"time"

	"typeconv"

	log "github.com/Sirupsen/logrus"
)

func WaitServiceReady(serviceName string, n int) {
	current := GetServiceProviderCount(serviceName)
	log.WithFields(log.Fields{"need": n, "current": current}).Printf("Waiting %s service to be ready", serviceName)
	for current < n {
		time.Sleep(100 * time.Millisecond)
		current = GetServiceProviderCount(serviceName)
	}
	log.WithFields(log.Fields{"num": current}).Printf("Service %s is ready now.", serviceName)
}

func (s *String) ReadInt() int64 {
	return typeconv.ToInt(s.Read())
}

func (s *String) ReadIntTuple() []int64 {
	v := s.Read()
	return typeconv.ToIntTuple(v)
}

func (s *String) ReadString() string {
	return s.Read().(string)
}
