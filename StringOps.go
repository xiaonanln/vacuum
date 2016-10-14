package vacuum

func CreateString(name string) (string, error) {
	routine := getStringRoutine(name)
	s := newString(routine)
	putString(s)

	s.run()
	return s.ID, nil
}
