package gameserver

type AOI struct {
	sightDistance len_t
	entities      map[*GSEntity]bool
}

func (a *AOI) init() {
	a.sightDistance = 0 // entity with initial AOI distance = 0
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

func (a *AOI) Clear() {
	a.entities = map[*GSEntity]bool{}
}
