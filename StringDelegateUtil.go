package vacuum

import "github.com/xiaonanln/vacuum/common"

type _FuncPtrStringDelegate struct {
	init func(s *String, args ...interface{})
	loop func(s *String, msg common.StringMessage)
	fini func(s *String)
}

func (d *_FuncPtrStringDelegate) Init(s *String, args ...interface{}) {
	if d.init != nil {
		d.init(s, args...)
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

//type _InitStringDelegate func(s *String, args ...interface{})
//
//func (d _InitStringDelegate) Init(s *String, args ...interface{}) {
//	d(s, args...)
//}
//
//func (m _InitStringDelegate) Fini(s *String) {}
//
//func (m _InitStringDelegate) Loop(s *String, msg common.StringMessage) {
//	return
//}

func InitStringDelegateMaker(init func(s *String, args ...interface{})) StringDelegateMaker {
	//return func() StringDelegate {
	//	return _InitStringDelegate(init)
	//}

	return func() StringDelegate {
		return &_FuncPtrStringDelegate{
			init: func(s *String, args ...interface{}) {
				init(s, args...)
				s.Send(s.ID, nil) // trick: make string quit immediately
			},
			loop: nil,
			fini: nil,
		}
	}
}
