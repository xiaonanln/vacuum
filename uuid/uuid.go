package uuid

import "strconv"

var (
	nextUUID uint64 = 1
)

func GenUUID() string {
	return strconv.FormatUint(nextUUID, 36)
}
