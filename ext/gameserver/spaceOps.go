package gameserver

import (
	"fmt"

	"github.com/xiaonanln/vacuum/ext/entity"
)

type GSSpaceID entity.EntityID

func (sid GSSpaceID) String() string {
	return fmt.Sprintf("SpaceID#%s", string(sid))
}

func (spaceID GSSpaceID) getLocalSpace() *GSSpace {
	localEntity := entity.EntityID(spaceID).GetLocalEntity()
	if localEntity == nil {
		return nil
	}
	return localEntity.(*GSSpace)
}

// Create a space
func CreateSpace(kind int) GSSpaceID {
	eid := entity.CreateEntity("GSSpace", kind)
	return GSSpaceID(eid)
}

func CreateSpaceLocally(kind int) GSSpaceID {
	eid := entity.CreateEntityLocally("GSSpace", kind)
	return GSSpaceID(eid)
}

func GetNilSpace() *GSSpace {
	return nilSpace
}
