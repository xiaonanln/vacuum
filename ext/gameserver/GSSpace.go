package gameserver

import (
	"fmt"

	"github.com/xiaonanln/typeconv"
	"github.com/xiaonanln/vacuum/ext/entity"
)

type GSSpace struct {
	entity.Entity
	Kind int

	entities map[*GSEntity]bool
}

func (space *GSSpace) Init() {
	spaceKind := typeconv.Int(space.Args()[0])
	space.Kind = int(spaceKind)
	spaceDelegate.OnLoaded(space)
}

func (space *GSSpace) String() string {
	return fmt.Sprintf("GSSpace|%d|%s", space.Kind, space.ID)
}

// Create entity in space
func (space *GSSpace) CreateEntity(kind int) {
	entity.CreateEntityLocally(ENTITY_TYPE, kind)
}
