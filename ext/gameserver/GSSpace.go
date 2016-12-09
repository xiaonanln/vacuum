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
	args := space.Args()
	spaceKind := typeconv.Int(args[0])

	space.Kind = int(spaceKind)
	space.entities = map[*GSEntity]bool{}

	spaceDelegate.OnReady(space)
}

func (space *GSSpace) String() string {
	return fmt.Sprintf("GSSpace|%d|%s", space.Kind, space.ID)
}

// Create entity in space
func (space *GSSpace) CreateEntity(kind int, pos Vec3) {
	entity.CreateEntityLocally(ENTITY_TYPE, kind, space.ID, pos.X, pos.Y, pos.Z)
}

func (space *GSSpace) onEntityCreated(entity *GSEntity) {
	space.entities[entity] = true
	aoidist := entity.aoi.sightDistance
	for other, _ := range space.entities {
		if other != entity {
			other.checkAOI(entity)
			if aoidist > 0 {
				entity.checkAOI(other)
			}
		}
	}
}

func (space *GSSpace) GetEntityCount() int {
	return len(space.entities)
}
