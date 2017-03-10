package sync_util

import "sync"

type PersistentCond struct {
	L        sync.Locker
	cond     *sync.Cond
	signaled bool
}

func NewPersistentCond(l sync.Locker) *PersistentCond {
	return &PersistentCond{
		L:    l,
		cond: sync.NewCond(l),
	}
}

func (p *PersistentCond) Wait() {
	p.L.Lock()
	for !p.signaled {
		p.cond.Wait()
	}
	p.L.Unlock()
}

func (p *PersistentCond) Signal() {
	p.L.Lock()
	p.signaled = true
	p.cond.Broadcast()
	p.L.Unlock()
}
