package gameserver

type AOI struct {
	sightDistance Len_t
	entities      GSEntitySet
}

func (a *AOI) init() {
	a.sightDistance = 0 // entity with initial AOI distance = 0
	a.entities = GSEntitySet{}
}

func (a *AOI) InRange(entity *GSEntity) bool {
	return a.entities.Contains(entity)
}

func (a *AOI) Add(entity *GSEntity) {
	a.entities.Add(entity)
}

func (a *AOI) Remove(entity *GSEntity) {
	a.entities.Remove(entity)
}

func (a *AOI) Clear() {
	a.entities.Clear()
}
