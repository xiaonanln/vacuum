package storage

type StringStorage interface {
	Write(name string, stringID string, data interface{}) error
	Read(name string, stringID string) (interface{}, error)
}



