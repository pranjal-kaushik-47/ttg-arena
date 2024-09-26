package entity

import (
	"bytes"
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	ebiten "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type TextMessage struct {
	Tag       string
	Message   string
	Variables []any
	Position  *Point
	Size      float64
}

type ScreenText struct {
	Messages []*TextMessage
}

type textMap map[string]*TextMessage

var TextMap textMap = make(textMap)

func (s *ScreenText) AddText(text, tag string, X, Y float64, size float64) {
	taxtMessage := &TextMessage{
		Message:  text,
		Position: &Point{X: X, Y: Y},
		Size:     size,
	}
	if len(s.Messages) == 0 {
		s.Messages = []*TextMessage{taxtMessage}
	} else {
		s.Messages = append(s.Messages, taxtMessage)
	}
	TextMap[tag] = taxtMessage
}

func (s *ScreenText) Draw(screen *ebiten.Image) {
	for _, msg := range s.Messages {

		op := &text.DrawOptions{}
		op.GeoM.Translate(msg.Position.X, msg.Position.Y)
		op.ColorScale.ScaleWithColor(color.White)

		s, _ := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))

		text.Draw(screen, fmt.Sprintf(msg.Message, msg.Variables...), &text.GoTextFace{
			Source: s,
			Size:   msg.Size,
		}, op)
	}
}
