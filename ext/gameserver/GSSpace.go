package gameserver

import (
	"fmt"

	"github.com/xiaonanln/typeconv"
	"github.com/xiaonanln/vacuum/ext/entity"
	"github.com/xiaonanln/vacuum/vlog"
)

const (
	DEFAULT_AOI_DISTANCE = 100
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
	Kind        int
	aoiDistance len_t

	entities map[*GSEntity]bool
}

func (space *GSSpace) Init() {
	args := space.Args()
	spaceKind := typeconv.Int(args[0])

	space.Kind = int(spaceKind)
	space.entities = map[*GSEntity]bool{}
	space.aoiDistance = DEFAULT_AOI_DISTANCE

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
	for other, _ := range space.entities {
		if other != entity {
			other.checkAOI(entity)
			entity.checkAOI(other)
		}
	}
}

func (space *GSSpace) SetAOIDistance(dist len_t) {
	if dist <= 0 {
		vlog.Panicf("SetAOIDistance: AOI distance should be positive, not %v", dist)
	}

	space.aoiDistance = dist
}
