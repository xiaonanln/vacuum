package vacuum

import (
	"time"

	"github.com/xiaonanln/typeconv"
	"github.com/xiaonanln/vacuum/vlog"
)

func WaitServiceReady(serviceName string, n int) {
	current := GetServiceProviderCount(serviceName)
	vlog.Debug("Waiting service %s to be ready: need %d, current %d", serviceName, n, current)
	for current < n {
		time.Sleep(WAIT_OPS_LOOP_INTERVAL)
		current = GetServiceProviderCount(serviceName)
		//vlog.Debug("WaitServiceReady %s: current %d, need %d", serviceName, current, n)
	}
	vlog.Debug("Service %s is ready now: num=%d", serviceName, current)
}

func WaitServiceGone(serviceName string) {
	current := GetServiceProviderCount(serviceName)
	vlog.Debug("Waiting service %s to be gone: current %d", serviceName, current)
	for current > 0 {
		time.Sleep(WAIT_OPS_LOOP_INTERVAL)
		current = GetServiceProviderCount(serviceName)
	}
	vlog.Debug("Service %s is gone now", serviceName)
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
