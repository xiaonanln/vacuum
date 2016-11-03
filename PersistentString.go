package vacuum

type PersistentString interface {
	GetPersistentData() map[string]interface{}
	LoadPersistentData(data map[string]interface{})
}
