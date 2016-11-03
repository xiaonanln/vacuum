package storage

type StringStorage interface {
	Write(stringID string, data interface{})
	Read(stringID string) interface{}
}
