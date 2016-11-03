package vacuum

import "github.com/xiaonanln/vacuum/storage"

var (
	serverID      int
	stringStorage storage.StringStorage
)

func Setup(_serverID int, _storage storage.StringStorage) {
	serverID = _serverID
	stringStorage = _storage
}
