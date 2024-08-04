package entity

import (
	"encoding/json"
	"os"
	"tag-game-v2/common"

	ebiten "github.com/hajimehoshi/ebiten/v2"
)

type Environment struct {
	Walls        []Polygon
	ScreenWidth  int
	ScreenHeight int
}

func (e *Environment) Draw(screen *ebiten.Image, metaData common.GameMetaData) error {
	currentLevel := metaData.Level
	// screenWidth, screenHeight := metaData.ScreenWidth, metaData.ScreenHeight
	// e.buildWalls(currentLevel, screenWidth, screenHeight)
	for _, wall := range e.Walls {
		wall.Draw(screen, currentLevel)
	}
	return nil
}

func (e *Environment) BuildWalls(currentLevel int, screenWidth, screenHeight int) {
	if currentLevel == 0 {
		walls := make([]Polygon, 0)
		if len(e.Walls) == 0 {
			plan, _ := os.ReadFile("internal/gamedata/levels/level1.json")
			json.Unmarshal(plan, &walls)
			e.Walls = walls
		}
	}
}

func (e *Environment) Colliding(tag string, Polygon Polygon) bool {

	c := false
	for _, wall := range e.Walls {
		if wall.PolygonCollision(Polygon.Vertices) {
			c = true
		}
	}
	return c
}
