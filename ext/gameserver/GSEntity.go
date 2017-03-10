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

func (eid GSEntityID) GetLocalGSEntity() *GSEntity {
	entity := entity.EntityID(eid).GetLocalEntity()
	if entity != nil {
		return entity.(*GSEntity)
	} else {
		return nil
	}
}

type GSEntity struct {
	entity.Entity
	aoi      AOI
	ID       GSEntityID
	Space    *GSSpace
	KindName string
	kindVal  reflect.Value
	Kind     IGSEntityKind
	Pos      Vec3
	client   *GSClientProxy

	enteringSpaceID GSSpaceID
}

func (ge *GSEntity) String() string {
	return fmt.Sprintf("%s<%s>", ge.KindName, ge.ID)
}

func (ge *GSEntity) Init() {
	ge.ID = GSEntityID(ge.Entity.ID)

	args := ge.Args()

	entityKind := typeconv.String(args[0])
	spaceID := GSSpaceID(typeconv.String(args[1]))

	var x, y, z Len_t
	x = typeconv.Convert(args[2], reflect.TypeOf(x)).Interface().(Len_t)
	y = typeconv.Convert(args[3], reflect.TypeOf(y)).Interface().(Len_t)
	z = typeconv.Convert(args[4], reflect.TypeOf(z)).Interface().(Len_t)

	ge.KindName = entityKind
	ge.kindVal = createGSEntityKind(ge, entityKind)
	ge.Kind = ge.kindVal.Interface().(IGSEntityKind)
	ge.Kind.Init()

	ge.aoi.init()
	ge.Pos.Assign(x, y, z)

	var space *GSSpace // FIXME: should I always set space to nilSpace here ?
	if spaceID != "" {
		space = spaceID.getLocalSpace()
	}
	if space == nil {
		space = GetNilSpace()
	}

	vlog.Debug("%s.Init: spaceID=%s, space=%s, pos=%s", ge, spaceID, space, ge.Pos)
	ge.Space = space
}

func (ge *GSEntity) OnReady() {
	space := ge.Space

	space.Lock()
	space.onEntityCreated(ge)
	space.Unlock()

	ge.Kind.OnEnterSpace()
}

func (ge *GSEntity) EnterSpace(spaceID GSSpaceID) {
	// FIXME: do not migrate if target entity is local
	ge.enteringSpaceID = spaceID
	ge.Entity.MigrateTowards(entity.EntityID(spaceID))
}

func (ge *GSEntity) checkAOI(other *GSEntity) {
	if ge.aoi.sightDistance <= 0 { // AOI disabled, sees nothing
		if ge.aoi.InRange(other) {
			ge.onLeaveAOI(other)
			vlog.Debug("%s MISS %s.", ge, other)
		}
	}

	dist := ge.DistanceTo(other)
	if dist < ge.aoi.sightDistance { // use < so that if AOI sightDistance is 0, entity sees nobody
		if !ge.aoi.InRange(other) {
			ge.onEnterAOI(other)
			vlog.Debug("%s SEES %s!", ge, other)
		}
	} else {
		if ge.aoi.InRange(other) {
			ge.onLeaveAOI(other)
			vlog.Debug("%s MISS %s.", ge, other)

		}
	}
}

func (ge *GSEntity) onEnterAOI(other *GSEntity) {
	ge.aoi.Add(other)
}

func (ge *GSEntity) onLeaveAOI(other *GSEntity) {
	ge.aoi.Remove(other)
}

func (ge *GSEntity) DistanceTo(other *GSEntity) Len_t {
	return ge.Pos.DistanceTo(other.Pos)
}

func (ge *GSEntity) SetAOIDistance(dist Len_t) {
	if dist < 0 {
		vlog.Panicf("SetAOIDistance: AOI distance should be positive, not %v", dist)
	}

	ge.aoi.sightDistance = dist
	if ge.Space.Kind > 0 {
		for otherEntity, _ := range ge.Space.entities { // check all entities in space for AOI
			if otherEntity != ge {
				ge.checkAOI(otherEntity)
			}
		}
	}
}

func (ge *GSEntity) SetPos(pos Vec3) {
	vlog.Debug("%s.SetPos %s", ge, pos)
	ge.Pos = pos
	space := ge.Space

	if space.Kind > 0 {
		// position changed, recheck AOI!
		aoidist := ge.aoi.sightDistance
		for other, _ := range space.entities {
			if other != ge {
				other.checkAOI(ge)
				if aoidist > 0 {
					ge.checkAOI(other)
				}
			}
		}
	}
}

func (ge *GSEntity) AOIEntities() GSEntitySet {
	return ge.aoi.entities
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
	entity.Kind.OnDestroy() // destroy kind before destroy entity
	entity.Entity.Destroy()
}

// TODO: think how should GiveClientTo works, should it only support local entity ?
// Give client to another entity
func (ge *GSEntity) GiveClientTo(otherID GSEntityID) {
	client := ge.client
	if client == nil {
		vlog.Warn("%s.GiveClientTo %s: has no client", ge, otherID)
		return
	}

	ge.client = nil
	// Tell the client to change owner
	client.notifyChangeOwner(ge.ID, otherID, "Avatar")

	ge.Kind.OnLoseClient()
}

func (ge *GSEntity) CallClient(methodName string, args ...interface{}) {
	if ge.client == nil {
		vlog.Debug("%s.CallClient: %s: client is nil", ge, methodName)
		return
	}

	ge.client.callClient(ge.ID, methodName, args)
}

func (ge *GSEntity) CallGSRPC_OwnClient(methodName string, args []interface{}) {
	methodName = methodName + "_OwnClient"
	method := ge.kindVal.MethodByName(methodName)
	vlog.Debug("CallGSRPC_OwnClient: method=%s(%v), args=%v", methodName, method, args)
	methodType := method.Type()

	in := make([]reflect.Value, len(args))

	for i, arg := range args {
		argType := methodType.In(i)
		in[i] = typeconv.Convert(arg, argType)
	}
	method.Call(in)
}

func (ge *GSEntity) NotifyGetClient(gateID GSGateID, clientID GSClientID) {
	client := newGSClientProxy(gateID, clientID)
	vlog.Debug("%s.NotifyGetClient: %s", ge, client)

	if ge.client != nil {
		// entity already has client, fail
		vlog.Panicf("%s.NotifyGetClient: new client %s, already has client %s", ge, client, ge.client)
	}

	ge.client = client
	ge.Kind.OnGetClient()
}

func (ge *GSEntity) NotifyLoseClient(gateID GSGateID, clientID GSClientID) {
	vlog.Debug("%s.NotifyLoseClient: lose client %s@%s", ge, clientID, gateID)

	if ge.client == nil || ge.client.ClientID != clientID {
		vlog.Warn("%s.NotifyLoseClient: has client %s, but lose client %s", ge, ge.client, clientID)
		return
	}

	ge.client = nil
	ge.Kind.OnLoseClient()
}

func (ge *GSEntity) OnMigrateOut(extra map[string]interface{}) {
	extra["C"] = ge.client.getClientProxyData()
	extra["ES"] = ge.enteringSpaceID
	kindExtra := map[string]interface{}{}
	ge.Kind.OnMigrateOut(kindExtra)
	extra["K"] = kindExtra
}

func (ge *GSEntity) OnMigrateIn(extra map[string]interface{}) {
	clientProxyData := extra["C"]
	if clientProxyData != nil {
		client := &GSClientProxy{}
		client.setClientProxyData(clientProxyData)
		ge.client = client // just store client, do not call OnGetClient
	}
	enteringSpaceID := GSSpaceID(extra["ES"].(string))
	if enteringSpaceID != "" {
		// entity entering space
		space := enteringSpaceID.getLocalSpace()
		if space != nil {
			ge.Space = space
		} else {
			// space not found ?
			vlog.Warn("%s.OnMigrateIn: entering space %s, but not found on local server", ge, enteringSpaceID)
		}
	}

	kindExtra := extra["K"]
	ge.Kind.OnMigrateIn(typeconv.MapStringAnything(kindExtra))
}

func (ge *GSEntity) MigrateTowards(otherID GSEntityID) {
	vlog.Debug("MigrateTowards %s", otherID)
	ge.Entity.MigrateTowards(entity.EntityID(otherID))
}

func CreateGSEntity(kindName string) GSEntityID {
	entityID := entity.CreateEntity("GSEntity", kindName, "", 0, 0, 0)
	return GSEntityID(entityID)
}

func CreateGSEntityLocally(kindName string) GSEntityID {
	entityID := entity.CreateEntityLocally("GSEntity", kindName, "", 0, 0, 0)
	return GSEntityID(entityID)
}

func createGSEntity(kindName string, spaceID GSSpaceID, pos Vec3) GSEntityID {
	entityID := entity.CreateEntityLocally("GSEntity", kindName, spaceID, pos.X, pos.Y, pos.Z)
	return GSEntityID(entityID)
}

func LoadGSEntity(kindName string, entityID GSEntityID) {
	entity.LoadEntity("GSEntity", entity.EntityID(entityID), kindName, "", 0, 0, 0)
}
