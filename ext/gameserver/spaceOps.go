package gameserver

import (
	"fmt"

	"github.com/xiaonanln/vacuum/ext/entity"
)

type SpaceID entity.EntityID

func (sid SpaceID) String() string {
	return fmt.Sprintf("SpaceID#%s", string(sid))
}

func (spaceID SpaceID) getLocalSpace() *GSSpace {
	localEntity := entity.EntityID(spaceID).GetLocalEntity()
	if localEntity == nil {
		return nil
	}
	return localEntity.(*GSSpace)
}

// Create a space
func CreateSpace(kind int) SpaceID {
	eid := entity.CreateEntity("GSSpace", kind)
	spaceID := SpaceID(eid)
	return spaceID
}
