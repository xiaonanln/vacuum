package uuid

import "strconv"

var (
	nextUUID uint64 = 0
)

func GenUUID() string {
	nextUUID = nextUUID + 1
	return strconv.FormatUint(nextUUID, 36)
}
