package gameserver

import (
	"fmt"

	"sync"

	"time"

	"github.com/xiaonanln/goTimer"
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
	sync.RWMutex
	Kind int

	entities GSEntitySet
	timers   map[*timer.Timer]bool
}

func (space *GSSpace) Init() {
	args := space.Args()
	spaceKind := typeconv.Int(args[0])

	space.Kind = int(spaceKind)
	space.entities = GSEntitySet{}

	spaceDelegate.OnReady(space)
}

func (space *GSSpace) String() string {
	return fmt.Sprintf("GSSpace|%d|%s", space.Kind, space.ID)
}

// Create entity in space
func (space *GSSpace) CreateEntity(kind int, pos Vec3) {
	entity.CreateEntityLocally("GSEntity", kind, space.ID, pos.X, pos.Y, pos.Z)
}

func (space *GSSpace) onEntityCreated(entity *GSEntity) {
	space.entities.Add(entity)
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

func (space *GSSpace) AddCallback(d time.Duration, callback func()) {
	timer.AddCallback(d, func() {
		space.Lock()
		callback()
		space.Unlock()
	})
}

func (space *GSSpace) AddTimer(d time.Duration, callback func()) {
	timer.AddTimer(d, func() {
		space.Lock()
		callback()
		space.Unlock()
	})
}

func (space *GSSpace) Destroy() {
	space.Entity.Destroy() // super.Destroy
}

func (space *GSSpace) Entities() GSEntitySet {
	return space.entities
}
