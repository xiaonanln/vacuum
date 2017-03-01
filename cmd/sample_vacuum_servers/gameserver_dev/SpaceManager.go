package main

import . "github.com/xiaonanln/vacuum/ext/gameserver"

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
