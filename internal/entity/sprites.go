package entity

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"tag-game-v2/common"

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

type CachedData struct {
	PosX float64 `json:"posx"`
	PosY float64 `json:"posy"`
}

type Sprite struct {
	PosX         float64       `json:"posx"`
	PosY         float64       `json:"posy"`
	Image        *ebiten.Image `json:"image"`
	ImageSource  string        `json:"imagesource"`
	Type         string        `json:"type"`
	Verb         string        `json:"verb"`
	CurrentFrame int           `json:"currentframe"`
	BoundingBox  Polygon       `json:"boundingbox"`
	IsActive     bool          `json:"isactive"`
	Height       float64       `json:"height"`
	Width        float64       `json:"width"`

	CachedData CachedData
}

var RenderSpriteBoundingBox bool = false

func (s *Sprite) Update() error {
	return nil
}

func (s *Sprite) Draw(screen *ebiten.Image) error {
	var f *os.File
	var err error
	if s.IsActive {
		opts := &ebiten.DrawImageOptions{}
		if s.Height != 0.0 && s.Width != 0.0 {
			opts.GeoM.Scale(s.Height*0.063, s.Width*0.063)
		}
		halfHight := s.Height / 2
		halfWidth := s.Width / 2
		opts.GeoM.Translate(s.PosX-halfWidth, s.PosY-halfHight)
		if s.ImageSource != "" {
			if s.Verb == "" {
				f, err = os.Open(s.ImageSource)
				if err != nil {
					panic(err)
				}
			} else {
				f, err = os.Open(fmt.Sprintf("resources\\animations\\%v\\%v\\%v.png", s.Type, s.Verb, (s.CurrentFrame%(common.FrameRate*10))-1))
				if err != nil {
					panic(err)
				}

			}
			defer f.Close()
			img, _, err := image.Decode(f)
			if err != nil {
				panic(err)
			}
			s.Image = ebiten.NewImageFromImage(img)
			screen.DrawImage(s.Image, opts)
		}

		if RenderSpriteBoundingBox {
			s.BoundingBox.Draw(screen)
		}

	}
	return nil
}

func (s *Sprite) ApplyBoundingBox(polygon Polygon) {
	s.BoundingBox = polygon
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
