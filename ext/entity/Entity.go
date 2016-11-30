package entity

import (
	"reflect"

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

type Entity interface {
}

type BaseEntity struct {
}

func RegisterEntity(typeName string, entityVal interface{}) {
	if !isEntityStringRegistered {
		registerEntityString()
	}

	if _, ok := registeredEntityTypes[typeName]; ok {
		vlog.Panicf("RegisterEntity: Entity type %s already registered", typeName)
	}
	vlog.Debug(">>> RegisterEntity: %s <<<", typeName)
	entityVal = reflect.Indirect(reflect.ValueOf(entityVal))
	entityType := reflect.TypeOf(entityVal)

	// register the string of entity
	registeredEntityTypes[typeName] = entityType
}

func registerEntityString() {
	vacuum.RegisterString(ENTITY_STRING_NAME, func() vacuum.StringDelegate {
		return &entityString{}
	})
}

func CreateEntity(typeName string) EntityID {
	stringID := vacuum.CreateString(ENTITY_STRING_NAME, typeName)
	return EntityID(stringID)
}

type entityString struct {
	entity Entity
}

func (es *entityString) Init(s *vacuum.String) {
	typeName := typeconv.String(s.Args()[0]) // get entity type
	vlog.Debug("Creating entity %s ...", typeName)
	entityTyp := registeredEntityTypes[typeName]
	entityPtrVal := reflect.New(entityTyp) // create entity and get its pointer
	es.entity = entityPtrVal.Interface().(Entity)
}

func (es *entityString) Loop(s *vacuum.String, msg common.StringMessage) {
	methodNameAndArgs := msg.([]interface{})
	methodName := typeconv.String(methodNameAndArgs[0])

	var args []interface{}
	if methodNameAndArgs[1] == nil {
		args = []interface{}{}
	} else {
		args = methodNameAndArgs[1].([]interface{})
	}

	vlog.Debug("EntityString Loop %s(%v)", methodName, args)

}

func (es *entityString) Fini(s *vacuum.String) {

}

func (es *entityString) GetPersistentData() map[string]interface{} {
	return nil
}

func (es *entityString) LoadPersistentData(data map[string]interface{}) {

}

func (eid EntityID) Call(methodName string, args ...interface{}) {
	vacuum.Send(string(eid), []interface{}{
		methodName,
		args,
	})
}
