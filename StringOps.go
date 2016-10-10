package vacuum

func CreateString(routine StringRoutine) (*String, error) {
	s := newString(routine)
	putString(s)

	s.run()
	return s, nil
}
