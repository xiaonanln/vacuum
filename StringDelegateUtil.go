package vacuum

import "github.com/xiaonanln/vacuum/common"

type _InitStringDelegate func(s *String, args ...interface{})

func (d _InitStringDelegate) Init(s *String, args ...interface{}) {
	d(s, args...)
	s.inputChan <- nil // trick: make string quit immediately
}

func (m _InitStringDelegate) Fini(s *String) {}

func (m _InitStringDelegate) Loop(s *String, msg common.StringMessage) {
	return
}

func InitStringDelegateMaker(init func(s *String, args ...interface{})) StringDelegateMaker {
	return func() StringDelegate {
		return _InitStringDelegate(init)
	}
}
