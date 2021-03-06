package entity

import (
	"reflect"

	"fmt"

	"sync"

	"time"

	"github.com/xiaonanln/typeconv"
	"github.com/xiaonanln/vacuum"
	"github.com/xiaonanln/vacuum/common"
	"github.com/xiaonanln/vacuum/vlog"
)

type EntityID string

const (
	ENTITY_STRING_NAME = "__entity_string__"
)

var (
	isEntityStringRegistered = false
	registeredEntityTypes    = map[string]reflect.Type{}

	entitiesLock sync.RWMutex
	entities     = map[EntityID]IEntity{} // all entities
)

func putEntity(id EntityID, entity IEntity) {
	entitiesLock.Lock()
	entities[id] = entity
	entitiesLock.Unlock()
}

func delEntity(id EntityID) {
	entitiesLock.Lock()
	delete(entities, id)
	entitiesLock.Unlock()
}

func getEntity(id EntityID) (ret IEntity) {
	entitiesLock.RLock()
	ret = entities[id]
	entitiesLock.RUnlock()
	return
}

type IEntity interface {
	//ID() EntityID
	Init()
	OnReady() // called when ready
	Destroy() // destroy the entity

	OnMigrateIn(extra map[string]interface{})
	OnMigrateOut(extra map[string]interface{})
}

type Entity struct {
	I    IEntity
	ID   EntityID
	Type string
	S    *vacuum.String
}

func (e *Entity) Init() {
	vlog.Debug("%s.Init: Args=%v", e, e.Args())
}

func (e *Entity) Migrate(serverID int) {
	e.S.Migrate(serverID)
}

func (e *Entity) MigrateTowards(otherID EntityID) {
	vlog.Debug("MigrateTowards %s", otherID)
	e.S.MigrateTowards(string(otherID))
}

func (e *Entity) OnMigrateOut(extra map[string]interface{}) {
	vlog.Debug("%s.OnMigrateOut: %s", e, extra)
}

func (e *Entity) OnMigrateIn(extra map[string]interface{}) {
	vlog.Debug("%s.OnMigrateIn: %s", e, extra)
}

// Destroy entity
func (e *Entity) Destroy() {
	vlog.Debug("%s.Destroy ...", e)
	delEntity(e.ID)
	e.S.Send(e.S.ID, nil) // send nil to self to terminate the string
}

func (e *Entity) String() string {
	return fmt.Sprintf("%s<%s>", e.Type, e.ID)
}

func (e *Entity) Save() {
	e.S.Save()
}

func (e *Entity) Args() []interface{} {
	return e.S.Args()[1:]
}

func (e *Entity) AddCallback(d time.Duration, callback func()) {
	e.S.AddCallback(d, callback)
}

//
//func (e *BaseEntity) ID() EntityID {
//	return
//}

func RegisterEntity(typeName string, entityPtr IEntity) {
	if !isEntityStringRegistered {
		isEntityStringRegistered = true
		registerEntityString()
	}

	if _, ok := registeredEntityTypes[typeName]; ok {
		vlog.Panicf("RegisterEntity: Entity type %s already registered", typeName)
	}
	entityVal := reflect.Indirect(reflect.ValueOf(entityPtr))
	entityType := entityVal.Type()

	// register the string of entity
	registeredEntityTypes[typeName] = entityType

	vlog.Debug(">>> RegisterEntity %s => %s <<<", typeName, entityType.Name())
}

func registerEntityString() {
	vacuum.RegisterString(ENTITY_STRING_NAME, &entityString{})
}

func prependTypeNameToArgs(typeName string, args []interface{}) []interface{} {
	argscount := len(args)
	newArgs := make([]interface{}, argscount+1, argscount+1)
	newArgs[0] = typeName
	copy(newArgs[1:], args)
	return newArgs
}

func CreateEntity(typeName string, args ...interface{}) EntityID {
	stringArgs := prependTypeNameToArgs(typeName, args)
	stringID := vacuum.CreateString(ENTITY_STRING_NAME, stringArgs...)
	return EntityID(stringID)
}

func CreateEntityLocally(typeName string, args ...interface{}) EntityID {
	stringArgs := prependTypeNameToArgs(typeName, args)
	stringID := vacuum.CreateStringLocally(ENTITY_STRING_NAME, stringArgs...)
	return EntityID(stringID)
}

// TODO: load entity fail if entity was loaded and destroyed
func LoadEntity(typeName string, entityID EntityID, args ...interface{}) {
	stringArgs := prependTypeNameToArgs(typeName, args)
	vacuum.LoadString(ENTITY_STRING_NAME, string(entityID), stringArgs...)
}

type entityString struct {
	vacuum.String
	entityID     EntityID
	entity       IEntity
	entityPtrVal reflect.Value
}

func (es *entityString) Init() {
	typeName := typeconv.String(es.Args()[0]) // get entity type
	entityTyp, ok := registeredEntityTypes[typeName]
	if !ok {
		vlog.Panicf("Entity %s is not registered", typeName)
	}

	entityPtrVal := reflect.New(entityTyp) // create entity and get its pointer
	es.entityPtrVal = entityPtrVal
	es.entity = entityPtrVal.Interface().(IEntity)
	es.entityID = EntityID(es.String.ID)

	baseEntity := reflect.Indirect(entityPtrVal).FieldByName("Entity").Addr().Interface().(*Entity)
	baseEntity.I = es.entity

	baseEntity.Type = typeName
	baseEntity.ID = es.entityID
	baseEntity.S = &es.String

	putEntity(baseEntity.ID, baseEntity.I)
	vlog.Debug("Creating entity %s: %v %v", typeName, entityTyp, es.entityPtrVal)

	baseEntity.I.Init()
}

func (es *entityString) OnReady() {
	es.entity.OnReady()
}

func (es *entityString) OnMigrateIn(extra map[string]interface{}) {
	es.entity.OnMigrateIn(extra)
}

func (es *entityString) OnMigrateOut(extra map[string]interface{}) {
	es.entity.OnMigrateOut(extra)
}

func (es *entityString) OnMigratedAway() {
	delEntity(es.entityID)
}

func (es *entityString) Fini() {
	// entity ID should already be removed from entities map
}

func (es *entityString) Loop(msg common.StringMessage) {
	defer func() {
		err := recover() // recover from any error during RPC call
		if err != nil {
			vlog.TraceError("RPC %s::%v paniced: %v", es.entityPtrVal.Type().String()[1:], msg, err)
		}
	}()
	methodNameAndArgs := msg.([]interface{})
	methodName := typeconv.String(methodNameAndArgs[0])

	var args []interface{}
	if methodNameAndArgs[1] == nil {
		args = []interface{}{}
	} else {
		args = methodNameAndArgs[1].([]interface{})
	}

	method := es.entityPtrVal.MethodByName(methodName)
	vlog.Debug("EntityString Loop %s(%v) => %v.%v", methodName, args, es.entityPtrVal, method)

	methodType := method.Type()

	in := make([]reflect.Value, len(args))

	for i, arg := range args {
		argType := methodType.In(i)
		in[i] = typeconv.Convert(arg, argType)
		// log.Printf("Arg %d is %T %v value %v => %v", i, arg, arg, argVal, in[i])
	}
	// log.Printf("arguments: %v", in)
	method.Call(in)
}

func (eid EntityID) Call(methodName string, args ...interface{}) {
	vacuum.Send(string(eid), []interface{}{
		methodName,
		args,
	})
}

func (eid EntityID) GetLocalEntity() IEntity {
	return getEntity(eid)
}
