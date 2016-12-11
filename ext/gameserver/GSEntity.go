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

	var x, y, z Len_t
	x = typeconv.Convert(args[2], reflect.TypeOf(x)).Interface().(Len_t)
	y = typeconv.Convert(args[3], reflect.TypeOf(y)).Interface().(Len_t)
	z = typeconv.Convert(args[4], reflect.TypeOf(z)).Interface().(Len_t)

	entity.Kind = int(entityKind)
	entity.aoi.init()
	entity.Pos.Assign(x, y, z)

	space := spaceID.getLocalSpace()
	vlog.Debug("%s.Init: space=%s, pos=%s", entity, space, entity.Pos)
	if space == nil {
		// how can space be destroy
		entity.I.Destroy()
		return
	}

	entity.space = space

	space.Lock()

	space.onEntityCreated(entity)
	entityDelegate.OnReady(entity)
	entityDelegate.OnEnterSpace(entity, space)

	space.Unlock()
}

func (entity *GSEntity) checkAOI(other *GSEntity) {
	if entity.aoi.sightDistance <= 0 { // AOI disabled, sees nothing
		if entity.aoi.InRange(other) {
			entity.onLeaveAOI(other)
			vlog.Debug("%s MISS %s.", entity, other)
		}
	}

	dist := entity.DistanceTo(other)
	if dist < entity.aoi.sightDistance { // use < so that if AOI sightDistance is 0, entity sees nobody
		if !entity.aoi.InRange(other) {
			entity.onEnterAOI(other)
			vlog.Debug("%s SEES %s!", entity, other)
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
	entityDelegate.OnEnterAOI(entity, other)
}

func (entity *GSEntity) onLeaveAOI(other *GSEntity) {
	entity.aoi.Remove(other)
	entityDelegate.OnLeaveAOI(entity, other)
}

func (entity *GSEntity) DistanceTo(other *GSEntity) Len_t {
	return entity.Pos.DistanceTo(other.Pos)
}

func (entity *GSEntity) SetAOIDistance(dist Len_t) {
	if dist < 0 {
		vlog.Panicf("SetAOIDistance: AOI distance should be positive, not %v", dist)
	}

	entity.aoi.sightDistance = dist
	for otherEntity, _ := range entity.space.entities { // check all entities in space for AOI
		if otherEntity != entity {
			entity.checkAOI(otherEntity)
		}
	}
}

func (entity *GSEntity) SetPos(pos Vec3) {
	vlog.Debug("%s.SetPos %s", entity, pos)
	entity.Pos = pos
	// position changed, recheck AOI!

	aoidist := entity.aoi.sightDistance
	space := entity.space

	for other, _ := range space.entities {
		if other != entity {
			other.checkAOI(entity)
			if aoidist > 0 {
				entity.checkAOI(other)
			}
		}
	}
}

func (entity *GSEntity) AOIEntities() GSEntitySet {
	return entity.aoi.entities
}
