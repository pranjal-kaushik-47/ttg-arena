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

func (e *Environment) BuildSquareWall(X, Y, Width, Height int) {
	halfHight := Height / 2
	halfWidth := Width / 2
	square := Polygon{
		Vertices: []ebiten.Vertex{
			{
				DstX:   float32(X - halfWidth),
				DstY:   float32(Y - halfHight),
				ColorR: 1,
				ColorG: 255,
				ColorB: 1,
				ColorA: 1,
			},
			{
				DstX:   float32(X + halfWidth),
				DstY:   float32(Y - halfHight),
				ColorR: 1,
				ColorG: 255,
				ColorB: 1,
				ColorA: 1,
			},
			{
				DstX:   float32(X + halfWidth),
				DstY:   float32(Y + halfHight),
				ColorR: 1,
				ColorG: 255,
				ColorB: 1,
				ColorA: 1,
			},
			{
				DstX:   float32(X - halfWidth),
				DstY:   float32(Y + halfHight),
				ColorR: 1,
				ColorG: 255,
				ColorB: 1,
				ColorA: 1,
			},
		},
		Indices: []uint16{
			0,
			1,
			2,
			2,
			3,
			0,
		},
		Color: Color{
			R: 255,
			G: 0,
			B: 0,
			A: 1,
		},
	}
	image := "resources\\images\\wall.png"
	e.Walls = append(e.Walls, Wall{Sprite: &Sprite{BoundingBox: square, IsActive: true, ImageSource: image, Height: float64(Height), Width: float64(Width), PosX: float64(X), PosY: float64(Y)}})
}
