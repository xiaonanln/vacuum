package gameserver

import (
	"fmt"

	"reflect"

	"github.com/xiaonanln/typeconv"
	"github.com/xiaonanln/vacuum/ext/entity"
	"github.com/xiaonanln/vacuum/vlog"
)

type GSEntityID entity.EntityID

// RPC call from client
func (eid GSEntityID) callGSRPC_OwnClient(method string, args []interface{}) {
	entity.EntityID(eid).Call("CallGSRPC_OwnClient", method, args)
}

// Notify the GSEntity for its own client
func (eid GSEntityID) notifyGetClient(gateID GSGateID, clientID GSClientID) {
	entity.EntityID(eid).Call("NotifyGetClient", gateID, clientID)
}

func (eid GSEntityID) notifyLoseClient(gateID GSGateID, clientID GSClientID) {
	entity.EntityID(eid).Call("NotifyLoseClient", gateID, clientID)
}

type GSEntity struct {
	entity.Entity
	aoi      AOI
	ID       GSEntityID
	space    *GSSpace
	KindName string
	kindVal  reflect.Value
	Kind     IGSEntityKind
	Pos      Vec3
	client   *GSClientProxy
}

func (entity *GSEntity) String() string {
	return fmt.Sprintf("%s<%s>", entity.KindName, entity.ID)
}

func (entity *GSEntity) Init() {
	entity.ID = GSEntityID(entity.Entity.ID)

	args := entity.Args()

	entityKind := typeconv.String(args[0])
	spaceID := GSSpaceID(typeconv.String(args[1]))

	var x, y, z Len_t
	x = typeconv.Convert(args[2], reflect.TypeOf(x)).Interface().(Len_t)
	y = typeconv.Convert(args[3], reflect.TypeOf(y)).Interface().(Len_t)
	z = typeconv.Convert(args[4], reflect.TypeOf(z)).Interface().(Len_t)

	entity.KindName = entityKind
	entity.kindVal = createGSEntityKind(entity, entityKind)
	entity.Kind = entity.kindVal.Interface().(IGSEntityKind)
	entity.Kind.Init()

	entity.aoi.init()
	entity.Pos.Assign(x, y, z)

	var space *GSSpace
	if spaceID != "" {
		space = spaceID.getLocalSpace()
	} else {
		space = GetNilSpace()
	}

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

func (entity *GSEntity) EnterSpace(spaceID GSSpaceID) {

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
	if entity.space.Kind > 0 {
		for otherEntity, _ := range entity.space.entities { // check all entities in space for AOI
			if otherEntity != entity {
				entity.checkAOI(otherEntity)
			}
		}
	}
}

func (entity *GSEntity) SetPos(pos Vec3) {
	vlog.Debug("%s.SetPos %s", entity, pos)
	entity.Pos = pos
	space := entity.space

	if space.Kind > 0 {
		// position changed, recheck AOI!
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
}

func (entity *GSEntity) AOIEntities() GSEntitySet {
	return entity.aoi.entities
}

//func CreateGSEntity(kind int, spaceID SpaceID, pos Vec3) GSEntityID {
//	eid := entity.CreateEntity("GSEntity", kind, spaceID, pos.X, pos.Y, pos.Z)
//	return GSEntityID(eid)
//}
//
//func CreateGSEntityLocally(kind int, spaceID SpaceID, pos Vec3) GSEntityID {
//	eid := entity.CreateEntityLocally("GSEntity", kind, spaceID, pos.X, pos.Y, pos.Z)
//	return GSEntityID(eid)
//}

func (entity *GSEntity) Destroy() {
	entity.Kind.Destroy() // destroy kind before destroy entity
	entity.Entity.Destroy()
}

// Give client to another entity
func (entity *GSEntity) GiveClientTo(otherID GSEntityID) {
	client := entity.client
	if client == nil {
		vlog.Warn("%s.GiveClientTo %s: has no client", entity, otherID)
		return
	}

	entity.client = nil
	// Tell the client to change owner
	client.notifyChangeOwner(entity.ID, otherID)

	entity.Kind.OnLoseClient()
}

func (entity *GSEntity) CallClient(methodName string, args ...interface{}) {
	if entity.client == nil {
		vlog.Debug("%s.CallClient: %s: client is nil", entity, methodName)
		return
	}

	entity.client.callClient(entity.ID, methodName, args)
}

func (entity *GSEntity) CallGSRPC_OwnClient(methodName string, args []interface{}) {
	methodName = methodName + "_OwnClient"
	method := entity.kindVal.MethodByName(methodName)
	vlog.Debug("CallGSRPC_OwnClient: method=%s(%v), args=%v", methodName, method, args)
	methodType := method.Type()

	in := make([]reflect.Value, len(args))

	for i, arg := range args {
		argType := methodType.In(i)
		in[i] = typeconv.Convert(arg, argType)
	}
	method.Call(in)
}

func (entity *GSEntity) NotifyGetClient(gateID GSGateID, clientID GSClientID) {
	client := newGSClientProxy(gateID, clientID)
	vlog.Debug("%s.NotifyGetClient: %s", entity, client)

	if entity.client != nil {
		// entity already has client, fail
		vlog.Panicf("%s.NotifyGetClient: new client %s, already has client %s", entity, client, entity.client)
	}

	entity.client = client
	entity.Kind.OnGetClient()
}

func (entity *GSEntity) NotifyLoseClient(gateID GSGateID, clientID GSClientID) {
	vlog.Debug("%s.NotifyLoseClient: lose client %s@%s", entity, clientID, gateID)

	if entity.client == nil || entity.client.ClientID != clientID {
		vlog.Warn("%s.NotifyLoseClient: has client %s, but lose client %s", entity, entity.client, clientID)
		return
	}

	entity.client = nil
	entity.Kind.OnLoseClient()
}

func (entity *GSEntity) OnMigrateOut(extra map[string]interface{}) {
	extra["C"] = entity.client.getClientProxyData()
	kindExtra := map[string]interface{}{}
	entity.Kind.OnMigrateOut(kindExtra)
	extra["K"] = kindExtra
}

func (entity *GSEntity) OnMigrateIn(extra map[string]interface{}) {
	clientProxyData := extra["C"]
	if clientProxyData != nil {
		client := &GSClientProxy{}
		client.setClientProxyData(clientProxyData)
		entity.client = client // just store client, do not call OnGetClient
	}

	kindExtra := extra["K"]
	entity.Kind.OnMigrateIn(kindExtra)
}

func CreateGSEntityAnywhere(kindName string) GSEntityID {
	entityID := entity.CreateEntityLocally("GSEntity", kindName, "", 0, 0, 0)
	return GSEntityID(entityID)
}

func createGSEntity(kindName string, spaceID GSSpaceID, pos Vec3) GSEntityID {
	entityID := entity.CreateEntityLocally("GSEntity", kindName, spaceID, pos.X, pos.Y, pos.Z)
	return GSEntityID(entityID)
}
