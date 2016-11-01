package vacuum

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/xiaonanln/typeconv"
)

func WaitServiceReady(serviceName string, n int) {
	current := GetServiceProviderCount(serviceName)
	log.WithFields(log.Fields{"need": n, "current": current}).Printf("Waiting %s service to be ready", serviceName)
	for current < n {
		time.Sleep(100 * time.Millisecond)
		current = GetServiceProviderCount(serviceName)
		log.Printf("current %d, need %d", current, n)
	}
	log.WithFields(log.Fields{"num": current}).Printf("Service %s is ready now.", serviceName)
}

func WaitServiceGone(serviceName string) {
	current := GetServiceProviderCount(serviceName)
	log.WithFields(log.Fields{"current": current}).Printf("Waiting %s service to be gone", serviceName)
	for current > 0 {
		time.Sleep(100 * time.Millisecond)
		current = GetServiceProviderCount(serviceName)
	}
	log.WithFields(log.Fields{"num": current}).Printf("Service %s is gone now.", serviceName)
}

func (s *String) ReadInt() int64 {
	return typeconv.Int(s.Read())
}

func (s *String) ReadIntTuple() []int64 {
	v := s.Read()
	return typeconv.IntTuple(v)
}

func (s *String) ReadString() string {
	return s.Read().(string)
}
