package entity

import (
	"image/color"

	ebiten "github.com/hajimehoshi/ebiten/v2"
)

type Color struct {
	R int `json:"r"`
	G int `json:"g"`
	B int `json:"b"`
	A int `json:"a"`
}

func (c Color) ToColor() color.Color {
	return color.RGBA{R: uint8(c.R), G: uint8(c.G), B: uint8(c.B), A: uint8(c.A)}
}

type Sprite struct {
	PosX        float64
	PosY        float64
	Image       *ebiten.Image
	BoundingBox Polygon
	IsActive    bool
}

func (s *Sprite) Update() error {
	return nil
}

func (s *Sprite) Draw(screen *ebiten.Image) error {
	if s.IsActive {
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(s.PosX, s.PosY)
		screen.DrawImage(s.Image, opts)
	}
	return nil
}

func (s *Sprite) BoundingBoxShiftRight(shiftby float64) Polygon {
	VerticesCopy := make([]ebiten.Vertex, len(s.BoundingBox.Vertices))
	copy(VerticesCopy, s.BoundingBox.Vertices)
	for i := 0; i < len(VerticesCopy); i += 1 {

		VerticesCopy[i].DstX += float32(shiftby)
	}
	return Polygon{Vertices: VerticesCopy}
}

func (s *Sprite) BoundingBoxShiftLeft(shiftby float64) Polygon {
	VerticesCopy := make([]ebiten.Vertex, len(s.BoundingBox.Vertices))
	copy(VerticesCopy, s.BoundingBox.Vertices)
	for i := 0; i < len(VerticesCopy); i += 1 {

		VerticesCopy[i].DstX -= float32(shiftby)
	}
	return Polygon{Vertices: VerticesCopy}
}

func (s *Sprite) BoundingBoxShiftUp(shiftby float64) Polygon {
	VerticesCopy := make([]ebiten.Vertex, len(s.BoundingBox.Vertices))
	copy(VerticesCopy, s.BoundingBox.Vertices)
	for i := 0; i < len(VerticesCopy); i += 1 {
		VerticesCopy[i].DstY -= float32(shiftby)
	}
	return Polygon{Vertices: VerticesCopy}
}

func (s *Sprite) BoundingBoxShiftDown(shiftby float64) Polygon {
	VerticesCopy := make([]ebiten.Vertex, len(s.BoundingBox.Vertices))
	copy(VerticesCopy, s.BoundingBox.Vertices)
	for i := 0; i < len(VerticesCopy); i += 1 {

		VerticesCopy[i].DstY += float32(shiftby)
	}
	return Polygon{Vertices: VerticesCopy}
}
