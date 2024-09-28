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
	ScreenHeight      = 500
	ScreenWidth       = 500
	BoundryEdgeBuffer = 15
)
