package main

import (
	. "github.com/xiaonanln/vacuum/ext/gameserver"
	"github.com/xiaonanln/vacuum/vlog"
)

var (
	spaceManager = newSpaceManager()
)

type SpaceManager struct {
	spaces map[int]GSSpaceID
}

func newSpaceManager() *SpaceManager {
	return &SpaceManager{
		spaces: map[int]GSSpaceID{},
	}
}

func (sm *SpaceManager) GetSpace(kind int) GSSpaceID {
	if kind == 0 {
		return GetNilSpace().ID
	}

	spaceID, ok := sm.spaces[kind]
	if !ok {
		spaceID = CreateSpace(kind)
		sm.spaces[kind] = spaceID
	}

	return spaceID
}

type MySpaceDelegate struct {
	SpaceDelegate
}

func (delegate *MySpaceDelegate) OnReady(space *GSSpace) {
	vlog.Debug("%s.OnReady: kind=%v", space, space.Kind)
	if space.Kind == 0 {
		delegate.onNullSpaceReady(space)
		return
	}

	//// normal space
	//for i := 0; i < NMONSTERS; i++ {
	//	space.CreateEntity(MONSTER, Vec3{100, 100, 100})
	//}
}

func (delegate *MySpaceDelegate) onNullSpaceReady(space *GSSpace) {

}

type MyEntityDelegate struct {
}
