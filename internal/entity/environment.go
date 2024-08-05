package entity

import (
	"encoding/json"
	"fmt"
	"os"
	"tag-game-v2/common"

	ebiten "github.com/hajimehoshi/ebiten/v2"
)

type Wall struct {
	Sprite *Sprite `json:"sprite"`
}

type Environment struct {
	Walls []Wall `json:"walls"`
}

// how to create square tiles in ebitan?
// find a easier way to design a map + create a level designer

func (e *Environment) Draw(screen *ebiten.Image, metaData common.GameMetaData) error {
	// screenWidth, screenHeight := metaData.ScreenWidth, metaData.ScreenHeight
	// e.buildWalls(currentLevel, screenWidth, screenHeight)
	for _, wall := range e.Walls {
		wall.Sprite.Draw(screen)
		wall.Sprite.BoundingBox.Draw(screen)

	}
	return nil
}

func (e *Environment) BuildWalls(currentLevel int, screenWidth, screenHeight int) {
	if currentLevel == 0 {
		walls := make([]Wall, 0)
		levelJson, _ := os.ReadFile("internal/gamedata/levels/level1.json")
		err := json.Unmarshal(levelJson, &walls)
		if err != nil {
			panic(err)
		}
		e.Walls = walls
		fmt.Println(e.Walls)
	}
}

func (e *Environment) Colliding(tag string, Polygon Polygon) bool {

	c := false
	for _, wall := range e.Walls {
		if wall.Sprite.BoundingBox.PolygonCollision(Polygon.Vertices) {
			c = true
		}
	}
	return c
}
