package gameserver

import (
	"fmt"

	"reflect"

	"github.com/xiaonanln/typeconv"
	"github.com/xiaonanln/vacuum/ext/entity"
	"github.com/xiaonanln/vacuum/vlog"
)

type GSEntity struct {
	entity.Entity
	aoi AOI

	space *GSSpace
	Kind  int
	Pos   Vec3
}

func (entity *GSEntity) String() string {
	return fmt.Sprintf("GSEntity|%d|%s", entity.Kind, entity.ID)
}

func (entity *GSEntity) Init() {
	args := entity.Args()

	entityKind := typeconv.Int(args[0])
	spaceID := SpaceID(typeconv.String(args[1]))

	var x, y, z len_t
	x = typeconv.Convert(args[2], reflect.TypeOf(x)).Interface().(len_t)
	y = typeconv.Convert(args[3], reflect.TypeOf(y)).Interface().(len_t)
	z = typeconv.Convert(args[4], reflect.TypeOf(z)).Interface().(len_t)

	entity.Kind = int(entityKind)
	entity.aoi.init()
	entity.Pos.Assign(x, y, z)

	space := spaceID.getLocalSpace()
	vlog.Debug("%s.Init: space=%s, pos=%s", entity, space, entity.Pos)
	if space == nil {
		// how can space be destroy
		entity.Destroy()
		return
	}

	entity.space = space
	space.onEntityCreated(entity)

	entityDelegate.OnReady(entity)
}

func (entity *GSEntity) checkAOI(other *GSEntity) {
	dist := entity.DistanceTo(other)
	if dist <= entity.space.aoiDistance {
		if !entity.aoi.InRange(other) {
			entity.onEnterAOI(other)
			vlog.Debug("%s SEE %s!", entity, other)
		}
	} else {
		if entity.aoi.InRange(other) {
			entity.onLeaveAOI(other)
			vlog.Debug("%s MISS %s.", entity, other)
		}
	}
}

func (entity *GSEntity) onEnterAOI(other *GSEntity) {
	entity.aoi.Add(other)
}

func (entity *GSEntity) onLeaveAOI(other *GSEntity) {
	entity.aoi.Remove(other)
}

func (entity *GSEntity) DistanceTo(other *GSEntity) len_t {
	return entity.Pos.DistanceTo(other.Pos)
}
