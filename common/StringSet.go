package common

type StringSet map[string]bool

func (ss StringSet) Contains(elem string) bool {
	return ss[elem]
}

func (ss StringSet) Add(elem string) {
	ss[elem] = true
}

func (ss StringSet) Remove(elem string) {
	delete(ss, elem)
}
