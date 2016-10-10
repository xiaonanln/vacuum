package vacuum

type Vacuum struct {
	strs map[string]*String
}

var (
	vacuumInstance = Vacuum{
		strs: map[string]*String{},
	}
)

func putString(s *String) {
	vacuumInstance.strs[s.ID] = s
}

func getString(sid string) *String {
	return vacuumInstance.strs[sid]
}
