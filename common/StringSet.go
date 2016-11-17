package common

import "github.com/xiaonanln/vacuum/vlog"

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

type StringList []string

func (sl *StringList) Remove(elem string) {
	widx := 0
	cpsl := *sl
	for idx, _elem := range cpsl {
		if _elem == elem {
			// ignore this elem by doing nothing
		} else if idx != widx {
			cpsl[widx] = _elem
			widx += 1
		}
	}

	*sl = cpsl[:widx]
}

func (sl *StringList) Append(elem string) {
	*sl = append(*sl, elem)
}

func init() {
	var sl StringList
	sl.Append("1")
	sl.Append("2")
	sl.Append("3")
	sl.Remove("3")
	sl.Remove("4")
	vlog.Info("sl %v", sl)
}
