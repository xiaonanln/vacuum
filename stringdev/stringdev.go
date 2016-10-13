package stringdev

type StringRoutine func(StringInt)
type StringMessage interface{}

type StringInt interface {
	ID() string
	Read() StringMessage
	Output(msg StringMessage)
}
