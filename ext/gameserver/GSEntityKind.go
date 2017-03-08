package gameserver

import (
	"reflect"

	"fmt"

	"github.com/xiaonanln/vacuum/vlog"
)

var (
	registeredEntityKinds = map[string]reflect.Type{}
)

type IGSEntityKind interface {
	Init()
	OnDestroy()
	// Client notifications
	OnGetClient()
	OnLoseClient()

	OnEnterSpace()
	OnLeaveSpace()

	OnMigrateOut(extra map[string]interface{})
	OnMigrateIn(extra map[string]interface{})
}

type GSEntityKind struct {
	*GSEntity
	KindName string
	EntityID GSEntityID
}

func (kind *GSEntityKind) String() string {
	return fmt.Sprintf("GSEntityKind<%s><%s>", kind.KindName, kind.EntityID)
}

func (kind *GSEntityKind) Init() {
	vlog.Debug("%s.Init() ...", kind)
}

func (kind *GSEntityKind) OnDestroy() {
	vlog.Debug("%s.Destroy() ...", kind)
}

func (kind *GSEntityKind) OnGetClient() {
	vlog.Debug("%s.OnGetClient: %s", kind, kind.GSEntity.client)
}

func (kind *GSEntityKind) OnLoseClient() {
	vlog.Debug("%s.OnLoseClient ...", kind)
}

func (kind *GSEntityKind) OnEnterSpace() {
	vlog.Debug("%s.OnEnterSpace: %s", kind, kind.GSEntity.space)
}

func (kind *GSEntityKind) OnLeaveSpace() {
	vlog.Debug("%s.OnLeaveSpace: %s", kind)
}

func (kind *GSEntityKind) OnMigrateOut(extra map[string]interface{}) {

}

func (kind *GSEntityKind) OnMigrateIn(extra map[string]interface{}) {

}

// SHORTCUTS TO ENTITY FUNCTIONALITIES
func (kind *GSEntityKind) Space() *GSSpace {
	return kind.GSEntity.space
}

func RegisterGSEntityKind(kindName string, entityKindPtr IGSEntityKind) {
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

func createGSEntityKind(entity *GSEntity, kindName string) reflect.Value {
	kindType, ok := registeredEntityKinds[kindName]
	if !ok {
		vlog.Panicf("Entity Kind %s is not registered", kindName)
	}

	entityKindPtrVal := reflect.New(kindType) // create entity and get its pointer

	gsEntityKind := reflect.Indirect(entityKindPtrVal).FieldByName("GSEntityKind").Addr().Interface().(*GSEntityKind)
	gsEntityKind.GSEntity = entity
	gsEntityKind.KindName = kindName
	gsEntityKind.EntityID = entity.ID

	return entityKindPtrVal
}
