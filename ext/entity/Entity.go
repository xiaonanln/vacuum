package entity

import (
	"reflect"

	"fmt"

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
)

type IEntity interface {
	//ID() EntityID
	Init()
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

func (e *Entity) String() string {
	return fmt.Sprintf("%s<%s>", e.Type, e.ID)
}

func (e *Entity) Save() {
	e.S.Save()
}

func (e *Entity) Args() []interface{} {
	return e.S.Args()[1:]
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

func CreateEntity(typeName string, args ...interface{}) EntityID {
	argscount := len(args)
	stringArgs := make([]interface{}, argscount+1, argscount+1)
	stringArgs[0] = typeName
	copy(stringArgs[1:], args)

	stringID := vacuum.CreateString(ENTITY_STRING_NAME, stringArgs...)
	return EntityID(stringID)
}

type entityString struct {
	vacuum.String

	entityPtr reflect.Value
}

func (es *entityString) Init() {
	typeName := typeconv.String(es.Args()[0]) // get entity type
	entityTyp, ok := registeredEntityTypes[typeName]
	if !ok {
		vlog.Panicf("Entity %s is not registered", typeName)
	}

	entityPtrVal := reflect.New(entityTyp) // create entity and get its pointer
	es.entityPtr = entityPtrVal

	baseEntity := reflect.Indirect(entityPtrVal).FieldByName("Entity").Addr().Interface().(*Entity)
	baseEntity.I = entityPtrVal.Interface().(IEntity)

	baseEntity.Type = typeName
	baseEntity.ID = EntityID(es.String.ID)
	baseEntity.S = &es.String

	vlog.Debug("Creating entity %s: %v %v", typeName, entityTyp, es.entityPtr)
	baseEntity.I.Init()
}

func (es *entityString) Loop(msg common.StringMessage) {
	defer func() {
		err := recover() // recover from any error during RPC call
		if err != nil {
			vlog.TraceError("RPC %s::%v paniced: %v", es.entityPtr.Type().String()[1:], msg, err)
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

	method := es.entityPtr.MethodByName(methodName)
	vlog.Debug("EntityString Loop %s(%v) => %v.%v", methodName, args, es.entityPtr, method)

	methodType := method.Type()

	in := make([]reflect.Value, len(args))

	for i, arg := range args {
		argType := methodType.In(i)
		argVal := reflect.ValueOf(arg)
		in[i] = typeconv.Convert(argVal, argType)
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
