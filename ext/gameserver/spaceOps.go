package gameserver

import (
	"fmt"

	"github.com/xiaonanln/vacuum/ext/entity"
)

type SpaceID entity.EntityID

func (sid SpaceID) String() string {
	return fmt.Sprintf("SpaceID#%s", string(sid))
}

// Create a space
func CreateSpace(kind int) SpaceID {
	eid := entity.CreateEntity(SPACE_ENTITY_TYPE, kind)
	spaceID := SpaceID(eid)
	return spaceID
}
