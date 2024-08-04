package entity

import (
	"tag-game-v2/common"

	ebiten "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Player struct {
	Sprite *Sprite
	Speed  float64
}

func (p *Player) Reset() error {
	image, _, err := ebitenutil.NewImageFromFile("resources\\images\\player.png")
	if err != nil {
		return err
	}
	if p.Sprite == nil {
		p.Sprite = &Sprite{}
	}
	p.Speed = 2
	p.Sprite.Image = image
	p.Sprite.PosX = 100
	p.Sprite.PosY = 100
	p.Sprite.IsActive = true
	p.Sprite.BoundingBox = Polygon{
		Vertices: []ebiten.Vertex{
			{DstX: float32(p.Sprite.PosX), DstY: float32(p.Sprite.PosY), ColorR: 0, ColorG: 255, ColorB: 0, ColorA: 1},
			{DstX: float32(p.Sprite.PosX) + 15, DstY: float32(p.Sprite.PosY), ColorR: 0, ColorG: 255, ColorB: 0, ColorA: 1},
			{DstX: float32(p.Sprite.PosX) + 15, DstY: float32(p.Sprite.PosY) + 15, ColorR: 0, ColorG: 255, ColorB: 0, ColorA: 1},
			{DstX: float32(p.Sprite.PosX), DstY: float32(p.Sprite.PosY) + 15, ColorR: 0, ColorG: 255, ColorB: 0, ColorA: 1},
		},
		Indices: []uint16{1, 0, 3, 1, 2, 3}, // square
		Color:   Color{R: 255, G: 0, B: 0, A: 1},
	}
	return nil
}

func (p *Player) Update(metaData *common.GameMetaData, env *Environment) error {

	screenWidth, screenHeight := metaData.ScreenWidth-metaData.BoundryEdgeBuffer, metaData.ScreenHeight-metaData.BoundryEdgeBuffer

	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) && !env.Colliding("player", p.Sprite.BoundingBoxShiftUp(2)) {
		if p.Sprite.PosY >= 0 {
			p.Sprite.BoundingBox.Vertices = p.Sprite.BoundingBoxShiftUp(p.Speed).Vertices
			p.Sprite.PosY -= p.Speed
		}

	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		if p.Sprite.PosY <= float64(screenHeight) && !env.Colliding("player", p.Sprite.BoundingBoxShiftDown(2)) {
			p.Sprite.BoundingBox.Vertices = p.Sprite.BoundingBoxShiftDown(p.Speed).Vertices
			p.Sprite.PosY += p.Speed

		}

	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		if p.Sprite.PosX >= 0 && !env.Colliding("player", p.Sprite.BoundingBoxShiftLeft(2)) {
			p.Sprite.BoundingBox.Vertices = p.Sprite.BoundingBoxShiftLeft(p.Speed).Vertices
			p.Sprite.PosX -= p.Speed

		}

	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		if p.Sprite.PosX <= float64(screenWidth) && !env.Colliding("player", p.Sprite.BoundingBoxShiftRight(2)) {
			p.Sprite.BoundingBox.Vertices = p.Sprite.BoundingBoxShiftRight(p.Speed).Vertices
			p.Sprite.PosX += p.Speed

		}
	}

	return nil
}

func (p *Player) Draw(screen *ebiten.Image) error {
	p.Sprite.Draw(screen)
	return nil
}
