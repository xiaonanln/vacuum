package gameserver

type GSEntitySet map[*GSEntity]bool

func (es GSEntitySet) Add(entity *GSEntity) {
	es[entity] = true
}

func (es GSEntitySet) Remove(entity *GSEntity) {
	delete(es, entity)
}

func (es GSEntitySet) Copy() GSEntitySet {
	copy := GSEntitySet{}
	for ent, _ := range es {
		copy[ent] = true
	}
	return copy
}

func (es GSEntitySet) Contains(entity *GSEntity) bool {
	return es[entity]
}

func (es *GSEntitySet) Clear() {
	*es = GSEntitySet{}
}
