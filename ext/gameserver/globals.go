package gameserver

var (
	nilSpace         *GSSpace
	gameserverConfig *GameserverConfig
)

func createGlobalEntities() {

	// create global stuffs
	CreateSpaceLocally(0)
}
