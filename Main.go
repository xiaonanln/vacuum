package vacuum

var (
	mainfunc func()
)

type _Main struct {
	String
}

func (s *_Main) Init() {
	mainfunc()
	//s.Send(s.ID, nil)
}

func RegisterMain(main func()) {
	mainfunc = main
	RegisterString("Main", &_Main{})
}
