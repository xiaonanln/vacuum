package gameserver

import (
	"fmt"

	"github.com/xiaonanln/typeconv"
	"github.com/xiaonanln/vacuum/ext/entity"
)

//var (
//	localSpacesLock sync.RWMutex
//	_localSpaces    = map[SpaceID]*GSSpace{}
//)
//
//func setLocalSpace(spaceID SpaceID, space *GSSpace) {
//	localSpacesLock.Lock()
//	_localSpaces[spaceID] = space
//	localSpacesLock.Unlock()
//}
//
//func getLocalSpace(spaceID SpaceID) (ret *GSSpace) {
//	localSpacesLock.RLock()
//	ret = _localSpaces[spaceID]
//	localSpacesLock.RUnlock()
//	return
//}

type GSSpace struct {
	entity.Entity
	Kind int

	entities map[*GSEntity]bool
}

func (space *GSSpace) Init() {
	spaceKind := typeconv.Int(space.Args()[0])
	space.Kind = int(spaceKind)

	spaceDelegate.OnReady(space)
}

func (space *GSSpace) String() string {
	return fmt.Sprintf("GSSpace|%d|%s", space.Kind, space.ID)
}

// Create entity in space
func (space *GSSpace) CreateEntity(kind int) {
	entity.CreateEntityLocally(ENTITY_TYPE, kind, space.ID)
}

func (space *GSSpace) onEntityCreated(entity *GSEntity) {

}
