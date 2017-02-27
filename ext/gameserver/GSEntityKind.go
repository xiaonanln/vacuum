package gameserver

import (
	"reflect"

	"github.com/xiaonanln/vacuum/ext/entity"
	"github.com/xiaonanln/vacuum/vlog"
)

var (
	registeredEntityKinds = map[string]reflect.Type{}
)

type IGSEntityKind interface {
	Init()
	Destroy()
}

type GSEntityKind struct {
}

func (kind *GSEntityKind) Init() {
	vlog.Debug("%s.Init() ...", kind)
}

func (kind *GSEntityKind) Destroy() {
	vlog.Debug("%s.Destroy() ...", kind)
}

func RegisterGSEntityKind(kindName string, entityKindPtr interface{}) {
	if _, ok := registeredEntityKinds[kindName]; ok {
		vlog.Panicf("RegisterEntity: Entity type %s already registered", kindName)
	}

	entityKindVal := reflect.Indirect(reflect.ValueOf(entityKindPtr))
	entityKindType := entityKindVal.Type()

	gsEntityKindField := entityKindVal.FieldByName("GSEntityKind")
	if !gsEntityKindField.IsValid() || gsEntityKindField.Type().Name() != "GSEntityKind" {
		vlog.Panicf("EntityKind %s is not valid, should has GSEntityKind field", entityKindType.Name())
	}

	// register the string of entity
	registeredEntityKinds[kindName] = entityKindType

	vlog.Debug(">>> RegisterGSEntityKind %s => %s <<<", kindName, entityKindType.Name())
}

func createGSEntityKind(kindName string) reflect.Value {
	kindType, ok := registeredEntityKinds[kindName]
	if !ok {
		vlog.Panicf("Entity Kind %s is not registered", kindName)
	}

	entityKindPtrVal := reflect.New(kindType) // create entity and get its pointer
	return entityKindPtrVal
}

func createGSEntity(kindName string, spaceID SpaceID, x, y, z Len_t) GSEntityID {
	entityID := entity.CreateEntityLocally("GSEntity", kindName, spaceID, x, y, z)
	return GSEntityID(entityID)
}
