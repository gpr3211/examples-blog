package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"image/color"
)

type FontSize int

const NormalFontSize FontSize = 12
const LargeFontSize FontSize = 16

type Text struct {
	content string
	color   *color.RGBA
	size    FontSize
	x       int
	y       int
}

var ColorWhite = color.RGBA{255, 255, 255, 0}


func NewText(str string, x, y int, size FontSize, c color.RGBA) *Text {
	return &Text{
		content: str,
		color:   &c,
		x:       x,
		y:       y,
		size:    size,
	}
}

func (g *Text) ChangeText(str string) {
	g.content = str
	g.Update()

}
func (g *Text) Resize(size FontSize) error {
	g.size = size
	g.Update()
	return nil
}

func (t *Text) Draw(screen *ebiten.Image) {
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(t.x), float64(t.y))
	op.ColorScale.ScaleWithColor(t.color)
	text.Draw(screen, t.content, &text.GoTextFace{
		Size:   float64(t.size),
		Source: currFont,
	}, op)
}
func (t *Text) Update() error { return nil }
