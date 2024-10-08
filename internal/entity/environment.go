package entity

import (
	"encoding/json"
	"os"
	"tag-game-v2/common"

	ebiten "github.com/hajimehoshi/ebiten/v2"
)

type Wall struct {
	Sprite *Sprite `json:"sprite"`
}

type Environment struct {
	Walls []Wall `json:"walls"`

	// Temp Var for level editing : Vertex Points
	TempPoints []Point
}

func (e *Environment) Draw(screen *ebiten.Image, metaData common.GameMetaData, playerPosition Point) error {
	for _, wall := range e.Walls {
		if Distance(&playerPosition, &Point{X: wall.Sprite.PosX, Y: wall.Sprite.PosY}) < 100 {
			wall.Sprite.Draw(screen)
		}
	}
	return nil
}

func (e *Environment) BuildWalls(currentLevel int, screenWidth, screenHeight int) {
	if currentLevel == 0 {
		walls := make([]Wall, 0)
		levelJson, _ := os.ReadFile("internal/gamedata/levels/newlevel.json")
		err := json.Unmarshal(levelJson, &walls)
		if err != nil {
			panic(err)
		}
		e.Walls = walls
	}
}

func (e *Environment) Colliding(tag string, Polygon Polygon) bool {

	c := false
	for _, wall := range e.Walls {
		if wall.Sprite.BoundingBox.IsActive {
			if wall.Sprite.BoundingBox.PolygonCollision(Polygon.Vertices) {
				c = true
			}
		}
	}
	return c
}

func (e *Environment) BuildSquareWall(X, Y, Width, Height int, enableColliders bool) {
	halfHight := Height / 2
	halfWidth := Width / 2

	var square Polygon
	if enableColliders {
		square = Polygon{
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
			Indices: SquareIndex,
			Color: Color{
				R: 255,
				G: 0,
				B: 0,
				A: 1,
			},
			IsActive: true,
		}
	} else {
		square = Polygon{
			IsActive: false,
		}
	}
	image := "resources\\currentimage\\sprite.png"
	e.Walls = append(e.Walls, Wall{Sprite: &Sprite{BoundingBox: square, IsActive: true, ImageSource: image, Height: float64(Height), Width: float64(Width), PosX: float64(X), PosY: float64(Y)}})
}

func (e *Environment) DrawCollider(X, Y int) {
	if len(e.TempPoints) == 3 {
		v := make([]ebiten.Vertex, 0)
		for _, i := range e.TempPoints {
			v = append(v, ebiten.Vertex{
				DstX: float32(i.X),
				DstY: float32(i.Y)})
		}
		v = append(v, ebiten.Vertex{
			DstX: float32(X),
			DstY: float32(Y),
		})

		in := SquareIndex
		p := Polygon{
			Vertices: v,
			Indices:  in,
			IsActive: true,
		}
		e.TempPoints = []Point{}
		e.Walls = append(e.Walls, Wall{Sprite: &Sprite{BoundingBox: p, IsActive: true}})
	} else {
		e.TempPoints = append(e.TempPoints, Point{X: float64(X), Y: float64(Y)})
	}
}
