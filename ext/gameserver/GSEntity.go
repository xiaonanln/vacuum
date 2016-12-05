package gameserver

import (
	"github.com/xiaonanln/vacuum/ext/entity"
	. "github.com/xiaonanln/vacuum/ext/gameserver/position"
)

type GSEntity struct {
	entity.Entity

	position Pos
}
