package vacuum

import "os"

var (
	mainfunc func()
)

type _Main struct {
	String
}

func (s *_Main) Init() {
	mainfunc()
	//s.Send(s.ID, nil)
	os.Exit(1)
}

func RegisterMain(main func()) {
	mainfunc = main
	RegisterString("Main", &_Main{})
}
