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
var font *text.GoTextFaceSource

func init() {
	font, _ = text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
}

func DrawAllText(screen *ebiten.Image) {
	for _, msg := range TextMap {
		op := &text.DrawOptions{}
		op.GeoM.Translate(msg.Position.X, msg.Position.Y)
		op.ColorScale.ScaleWithColor(color.White)

		text.Draw(screen, fmt.Sprintf(msg.Message, msg.Variables...), &text.GoTextFace{
			Source: font,
			Size:   msg.Size,
		}, op)
	}
}
