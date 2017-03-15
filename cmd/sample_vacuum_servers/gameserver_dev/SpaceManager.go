package main

import (
	"sync"

	. "github.com/xiaonanln/vacuum/ext/gameserver"
	"github.com/xiaonanln/vacuum/vlog"
)

var (
	spaceManager = newSpaceManager()
)

type SpaceManager struct {
	sync.Mutex
	spaces map[int]GSSpaceID
}

func newSpaceManager() *SpaceManager {
	return &SpaceManager{
		spaces: map[int]GSSpaceID{},
	}
}

func (sm *SpaceManager) GetSpace(kind int) GSSpaceID {
	sm.Lock()
	defer sm.Unlock()
	if kind == 0 {
		return GetNilSpace().ID
	}
	return sm.spaces[kind]
}

func (sm *SpaceManager) LoadSpace(kind int) {
	sm.Lock()
	defer sm.Unlock()
	if sm.GetSpace(kind) != "" {
		return // space already exists
	}

	CreateSpace(kind)
}

func (sm *SpaceManager) onSpaceReady(space *GSSpace) {
	sm.Lock()
	defer sm.Unlock()

	kind := space.Kind
	if sm.spaces[kind] != "" {
		vlog.Warn("%s.onSpaceReady: duplicate space of kind %v", sm, kind)
		return
	}
	sm.spaces[kind] = space.ID
}

type MySpaceDelegate struct {
	SpaceDelegate
}

func (delegate *MySpaceDelegate) OnReady(space *GSSpace) {
	vlog.Debug("%s.OnReady: kind=%v, existing=%v", space, space.Kind, spaceManager.GetSpace(space.Kind))
	if space.Kind == 0 {
		delegate.onNullSpaceReady(space)
		return
	}

	if spaceManager.GetSpace(space.Kind) != "" {
		return // duplicate space ?
	}

	spaceManager.onSpaceReady(space)

	//// normal space
	//for i := 0; i < NMONSTERS; i++ {
	//	space.CreateEntity(MONSTER, Vec3{100, 100, 100})
	//}
}

func (delegate *MySpaceDelegate) onNullSpaceReady(space *GSSpace) {

}

type MyEntityDelegate struct {
}
