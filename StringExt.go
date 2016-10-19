package vacuum

import (
	"log"
	"time"
)

func WaitServiceReady(serviceName string, n int) {
	current := GetServiceProviderCount(serviceName)
	log.Printf("Waiting %s service to be ready: need %d, current %d", serviceName, n, current)
	for current < n {
		time.Sleep(100 * time.Millisecond)
		current = GetServiceProviderCount(serviceName)
	}
	log.Printf("Service %s is ready now: num=%d", serviceName, current)
}

func (s *String) ReadInt() int {
	return s.Read().(int)
}
