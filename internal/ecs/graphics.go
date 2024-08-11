package ecs

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Sprite struct {
	ImageSource string
	Image       *ebiten.Image
	Height      float64
	Width       float64
	IsActive    bool
}

type Color struct {
	R int
	G int
	B int
	A int
}

func (s *Sprite) Init() error {
	image, _, err := ebitenutil.NewImageFromFile(s.ImageSource)
	if err != nil {
		return err
	}
	s.Image = image
	return nil
}

func (s *Sprite) Draw(screen *ebiten.Image, X, Y float64) {
	if s.IsActive {
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Scale(s.Width, s.Height)
		opts.GeoM.Translate(X, Y)
		if s.Image != nil {
			screen.DrawImage(s.Image, opts)
		}

	}
}
