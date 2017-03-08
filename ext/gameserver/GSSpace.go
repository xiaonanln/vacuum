package gameserver

import (
	"fmt"

	"sync"

	"time"

	"github.com/xiaonanln/goTimer"
	"github.com/xiaonanln/typeconv"
	"github.com/xiaonanln/vacuum/ext/entity"
	"github.com/xiaonanln/vacuum/vlog"
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
	ID   GSSpaceID
	Kind int

	entities GSEntitySet
	timers   map[*timer.Timer]bool
}

func (space *GSSpace) Init() {
	space.ID = GSSpaceID(space.Entity.ID)
	args := space.Args()
	spaceKind := typeconv.Int(args[0])

	space.Kind = int(spaceKind)
	space.entities = GSEntitySet{}

	if space.Kind == 0 { // nil space Init
		nilSpace = space
		vlog.Info("Nil space is set to: %s", nilSpace)
	}
}

func (space *GSSpace) OnReady() {
	spaceDelegate.OnReady(space)
}

func (space *GSSpace) String() string {
	return fmt.Sprintf("GSSpace|%d|%s", space.Kind, space.ID)
}

// Create entity in space
func (space *GSSpace) CreateEntity(kindName string, pos Vec3) GSEntityID {
	return createGSEntity(kindName, space.ID, pos)
}

func (space *GSSpace) onEntityCreated(entity *GSEntity) {
	space.entities.Add(entity)
	aoidist := entity.aoi.sightDistance
	if space.Kind > 0 {
		for other, _ := range space.entities {
			if other != entity {
				other.checkAOI(entity)
				if aoidist > 0 {
					entity.checkAOI(other)
				}
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

func (space *GSSpace) IsNil() bool {
	return space.Kind == 0
}
