package storage

type StringStorage interface {
	Write(stringID string, data interface{}) error
	Read(stringID string) (interface{}, error)
}
