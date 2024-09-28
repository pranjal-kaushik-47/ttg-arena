package common

type GameMetaData struct {
	TotalEnemies      int
	CurrentEnemyCount int
	Level             int
	BoundryEdgeBuffer int

	ScreenWidth  int
	ScreenHeight int
}

const (
	ScreenHeight      = 800
	ScreenWidth       = 1700
	BoundryEdgeBuffer = 15
	FrameRate         = 3

	//Animation folders
	Walking = "walking"
	Ideal   = "ideal"
)
