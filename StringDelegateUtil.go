package vacuum

import (
	"os"

	"github.com/xiaonanln/vacuum/common"
)

type _FuncPtrStringDelegate struct {
	init func(s *String)
	loop func(s *String, msg common.StringMessage)
	fini func(s *String)
}

func (d *_FuncPtrStringDelegate) Init(s *String) {
	if d.init != nil {
		d.init(s)
	}
}

func (d *_FuncPtrStringDelegate) Fini(s *String) {
	if d.fini != nil {
		d.fini(s)
	}
}

func (d *_FuncPtrStringDelegate) Loop(s *String, msg common.StringMessage) {
	if d.loop != nil {
		d.loop(s, msg)
	}
}

func InitOnlyStringDelegateMaker(init func(s *String)) StringDelegateMaker {
	return func() StringDelegate {
		return &_FuncPtrStringDelegate{
			init: func(s *String) {
				init(s)
				s.Send(s.ID, nil) // trick: make string quit immediately
			},
			loop: nil,
			fini: nil,
		}
	}
}

func LoopOnlyStringDelegateMaker(loop func(s *String, msg common.StringMessage)) StringDelegateMaker {
	return func() StringDelegate {
		return &_FuncPtrStringDelegate{
			init: nil,
			loop: loop,
			fini: nil,
		}
	}
}

func RegisterMain(main func(s *String)) {
	RegisterString("Main", InitOnlyStringDelegateMaker(func(s *String) {
		main(s)
		os.Exit(0)
	}))
}
