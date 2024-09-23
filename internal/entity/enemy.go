package entity

import (
	"math/rand/v2"
	"tag-game-v2/common"

	ebiten "github.com/hajimehoshi/ebiten/v2"
)

type Enemy struct {
	Sprite   *Sprite
	Speed    float64
	Eyesight float64
}

// movement logic update for the enemy to not get stuck in corners:
// if enemy position has not changed since last frame and player is close
// move away from last point and the player

func (e *Enemy) Reset(env *Environment, metaData *common.GameMetaData) error {
	// image, _, err := ebitenutil.NewImageFromFile("resources\\images\\enemy.png")
	// if err != nil {
	// 	return err
	// }
	if e.Sprite == nil {
		e.Sprite = &Sprite{}
	}
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
		Indices: []uint16{1, 0, 3, 1, 2, 3}, // square
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
		if enemy.Sprite.PosX <= float64(screenWidth) && !env.Colliding("enemy", enemy.Sprite.BoundingBoxShiftRight(2)) {
			enemy.Sprite.BoundingBox.Vertices = enemy.Sprite.BoundingBoxShiftRight(enemy.Speed).Vertices
			x = enemy.Speed
		} else {
			x = 0
		}
	} else {
		if enemy.Sprite.PosX >= 0 && !env.Colliding("enemy", enemy.Sprite.BoundingBoxShiftLeft(2)) {
			enemy.Sprite.BoundingBox.Vertices = enemy.Sprite.BoundingBoxShiftLeft(enemy.Speed).Vertices
			x = -1 * enemy.Speed
		} else {
			x = 0
		}
	}
	if player.Sprite.PosY <= enemy.Sprite.PosY {
		if enemy.Sprite.PosY <= float64(screenHeight) && !env.Colliding("enemy", enemy.Sprite.BoundingBoxShiftDown(2)) {
			enemy.Sprite.BoundingBox.Vertices = enemy.Sprite.BoundingBoxShiftDown(enemy.Speed).Vertices
			y = enemy.Speed
		} else {
			y = 0
		}
	} else {
		if enemy.Sprite.PosY >= 0 && !env.Colliding("enemy", enemy.Sprite.BoundingBoxShiftUp(2)) {
			enemy.Sprite.BoundingBox.Vertices = enemy.Sprite.BoundingBoxShiftUp(enemy.Speed).Vertices
			y = -1 * enemy.Speed
		} else {
			y = 0
		}
	}
	if x == 0 && y == 0 {
		screenCenterX, screenCenterY := float64(screenWidth/2), float64(screenHeight/2)
		if screenCenterX > enemy.Sprite.PosX {
			x = 1
		} else {
			x = -1
		}
		if screenCenterY > enemy.Sprite.PosY {
			y = 1
		} else {
			y = -1
		}
	}
	return x, y
}

func (e *Enemy) Update(metaData *common.GameMetaData, p *Player, env *Environment) error {

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
		e.Sprite.IsActive = false
	}
	return nil
}

func (p *Enemy) Draw(screen *ebiten.Image) error {

	p.Sprite.Draw(screen)
	return nil
}
