package gameserver

type AOI struct {
	entities map[*GSEntity]bool
}

func (a *AOI) init() {
	a.entities = map[*GSEntity]bool{}
}

func (a *AOI) InRange(entity *GSEntity) bool {
	return a.entities[entity]
}

func (a *AOI) Add(entity *GSEntity) {
	a.entities[entity] = true
}

func (a *AOI) Remove(entity *GSEntity) {
	delete(a.entities, entity)
}
