package entity

import (
	"math/rand/v2"
	"tag-game-v2/common"

	"github.com/google/uuid"
	ebiten "github.com/hajimehoshi/ebiten/v2"
)

const (
	Runner int = iota
	Chaser int = iota
)

type Enemy struct {
	id       string
	Sprite   *Sprite
	Speed    float64
	Eyesight float64
	Type     int
}

func (e *Enemy) Reset(env *Environment, metaData *common.GameMetaData) error {
	if e.Sprite == nil {
		e.Sprite = &Sprite{}
	}
	e.id = uuid.New().String()
	e.Type = 0  //rand.IntN(2)
	e.Speed = 1 //float64(rand.IntN(3))
	e.Sprite.ImageSource = "resources\\images\\enemy.png"
	e.Sprite.PosX = float64(rand.IntN(metaData.ScreenWidth - 15))
	e.Sprite.PosY = float64(rand.IntN(metaData.ScreenHeight - 15))
	e.Sprite.IsActive = true
	e.Eyesight = 100 //float64(rand.IntN(100))
	e.Sprite.BoundingBox = Polygon{
		Vertices: []ebiten.Vertex{
			{DstX: float32(e.Sprite.PosX), DstY: float32(e.Sprite.PosY), ColorR: 0, ColorG: 255, ColorB: 0, ColorA: 1},
			{DstX: float32(e.Sprite.PosX) + 15, DstY: float32(e.Sprite.PosY), ColorR: 0, ColorG: 255, ColorB: 0, ColorA: 1},
			{DstX: float32(e.Sprite.PosX) + 15, DstY: float32(e.Sprite.PosY) + 15, ColorR: 0, ColorG: 255, ColorB: 0, ColorA: 1},
			{DstX: float32(e.Sprite.PosX), DstY: float32(e.Sprite.PosY) + 15, ColorR: 0, ColorG: 255, ColorB: 0, ColorA: 1},
		},
		Indices: SquareIndex,
		Color:   Color{R: 255, G: 0, B: 0, A: 1},
	}
	if env.Colliding("enemy", e.Sprite.BoundingBox) {
		e.Reset(env, metaData)
	}
	return nil
}

func GetMovementDirection(enemy *Enemy, player *Player, env *Environment, screenWidth, screenHeight int) (float64, float64) {
	var x, y float64
	if player.Sprite.PosX <= enemy.Sprite.PosX {
		if enemy.Type == 0 {
			if enemy.Sprite.PosX <= float64(screenWidth) && !env.Colliding("enemy", enemy.Sprite.BoundingBoxShiftRight(2)) {
				enemy.Sprite.BoundingBox.Vertices = enemy.Sprite.BoundingBoxShiftRight(enemy.Speed).Vertices
				x = enemy.Speed
			} else {
				x = 0
			}
		} else {
			if enemy.Sprite.PosX <= float64(screenWidth) && !env.Colliding("enemy", enemy.Sprite.BoundingBoxShiftRight(2)) {
				enemy.Sprite.BoundingBox.Vertices = enemy.Sprite.BoundingBoxShiftLeft(enemy.Speed).Vertices
				x = -1 * enemy.Speed
			} else {
				x = 0
			}
		}
	} else {
		if enemy.Type == 0 {
			if enemy.Sprite.PosX >= 0 && !env.Colliding("enemy", enemy.Sprite.BoundingBoxShiftLeft(2)) {
				enemy.Sprite.BoundingBox.Vertices = enemy.Sprite.BoundingBoxShiftLeft(enemy.Speed).Vertices
				x = -1 * enemy.Speed
			} else {
				x = 0
			}
		} else {
			if enemy.Sprite.PosX >= 0 && !env.Colliding("enemy", enemy.Sprite.BoundingBoxShiftLeft(2)) {
				enemy.Sprite.BoundingBox.Vertices = enemy.Sprite.BoundingBoxShiftRight(enemy.Speed).Vertices
				x = enemy.Speed
			} else {
				x = 0
			}
		}
	}
	if player.Sprite.PosY <= enemy.Sprite.PosY {
		if enemy.Type == 0 {
			if enemy.Sprite.PosY <= float64(screenHeight) && !env.Colliding("enemy", enemy.Sprite.BoundingBoxShiftDown(2)) {
				enemy.Sprite.BoundingBox.Vertices = enemy.Sprite.BoundingBoxShiftDown(enemy.Speed).Vertices
				y = enemy.Speed
			} else {
				y = 0
			}
		} else {
			if enemy.Sprite.PosY <= float64(screenHeight) && !env.Colliding("enemy", enemy.Sprite.BoundingBoxShiftDown(2)) {
				enemy.Sprite.BoundingBox.Vertices = enemy.Sprite.BoundingBoxShiftUp(enemy.Speed).Vertices
				y = -1 * enemy.Speed
			} else {
				y = 0
			}
		}
	} else {
		if enemy.Type == 0 {
			if enemy.Sprite.PosY >= 0 && !env.Colliding("enemy", enemy.Sprite.BoundingBoxShiftUp(2)) {
				enemy.Sprite.BoundingBox.Vertices = enemy.Sprite.BoundingBoxShiftUp(enemy.Speed).Vertices
				y = -1 * enemy.Speed
			} else {
				y = 0
			}
		} else {
			if enemy.Sprite.PosY >= 0 && !env.Colliding("enemy", enemy.Sprite.BoundingBoxShiftUp(2)) {
				enemy.Sprite.BoundingBox.Vertices = enemy.Sprite.BoundingBoxShiftDown(enemy.Speed).Vertices
				y = enemy.Speed
			} else {
				y = 0
			}
		}
	}
	if x == 0 && y == 0 {
		screenCenterX, screenCenterY := float64(screenWidth/2), float64(screenHeight/2)
		if screenCenterX > enemy.Sprite.PosX {
			if enemy.Type == 0 {
				x = 1
			} else {
				x = -1
			}

		} else {
			if enemy.Type == 0 {
				x = -1
			} else {
				x = 1
			}
		}
		if screenCenterY > enemy.Sprite.PosY {
			if enemy.Type == 0 {
				y = 1
			} else {
				y = -1
			}
		} else {
			if enemy.Type == 0 {
				y = -1
			} else {
				y = 1
			}
		}
	}
	return x, y
}

func (e *Enemy) MoveEnemy(metaData *common.GameMetaData, p *Player, env *Environment) error {

	screenWidth, screenHeight := metaData.ScreenWidth-metaData.BoundryEdgeBuffer, metaData.ScreenHeight-metaData.BoundryEdgeBuffer

	distanceFromPlayer := Distance(&Point{X: e.Sprite.PosX, Y: e.Sprite.PosY}, &Point{X: p.Sprite.PosX, Y: p.Sprite.PosY})
	if distanceFromPlayer <= e.Eyesight {
		x, y := GetMovementDirection(e, p, env, screenWidth, screenHeight)

		e.Sprite.PosX += x
		e.Sprite.PosY += y

		e.Sprite.CachedData = CachedData{
			PosX: e.Sprite.PosX,
			PosY: e.Sprite.PosY,
		}
	}
	if distanceFromPlayer <= 5 {
		if e.Type == 0 {
			e.Sprite.IsActive = false
			delete(TextMap, e.id)
		} else if e.Type == 1 {
			p.Sprite.IsActive = false
		}

	}
	return nil
}

func (e *Enemy) Update(metaData *common.GameMetaData, p *Player, env *Environment) error {
	e.MoveEnemy(metaData, p, env)
	// reset collider if misplaced
	return nil
}

func (e *Enemy) Draw(screen *ebiten.Image) error {

	if e.Sprite.IsActive {
		e.Sprite.Draw(screen)
		TextMap[e.id] = &TextMessage{
			Message:   "Type: %v, Pos: (%v, %v)",
			Variables: []any{e.Type, e.Sprite.PosX, e.Sprite.PosY},
			Position:  &Point{X: e.Sprite.PosX + 10, Y: e.Sprite.PosY},
			Size:      10,
		}
	}

	return nil
}
