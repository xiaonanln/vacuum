package vacuum

import (
	"time"

	"github.com/xiaonanln/typeconv"
	"github.com/xiaonanln/vacuum/vlog"
)

func WaitServiceReady(serviceName string, n int) {
	current := GetServiceProviderCount(serviceName)
	vlog.Debugf("Waiting service %s to be ready: need %d, current %d", serviceName, n, current)
	for current < n {
		time.Sleep(100 * time.Millisecond)
		current = GetServiceProviderCount(serviceName)
		vlog.Debugf("WaitServiceReady %s: current %d, need %d", serviceName, current, n)
	}
	vlog.Debugf("Service %s is ready now: num=%d", serviceName, current)
}

func WaitServiceGone(serviceName string) {
	current := GetServiceProviderCount(serviceName)
	vlog.Debugf("Waiting service %s to be gone: current %d", serviceName, current)
	for current > 0 {
		time.Sleep(100 * time.Millisecond)
		current = GetServiceProviderCount(serviceName)
	}
	vlog.Debugf("Service %s is gone now", serviceName)
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
