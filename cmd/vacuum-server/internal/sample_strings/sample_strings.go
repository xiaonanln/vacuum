package sample_strings

import "github.com/xiaonanln/vacuum"

func Summer(s *vacuum.String) {
	var total uint64
	for {
		msg := s.Read()
		if msg != nil {
			total += msg.(uint64)
		} else {
			s.Output(total)
			break
		}
	}
}

func DoNothing(s *vacuum.String) {
	for {
		msg := s.Read()
		if msg == nil {
			break
		}
	}
}

func RedirectInputToOutput(s *vacuum.String) {
	for {
		msg := s.Read()
		if msg != nil {
			s.Output(msg)
		} else {
			break
		}
	}
}
