package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"image/color"
)

type FontSize int

const NormalFontSize FontSize = 12
const LargeFontSize FontSize = 16

// Add these fields to your Text struct
type Text struct {
	content       string // Full content
	displayedText string // Currently displayed text (grows gradually)
	color         *color.RGBA
	x, y          int
	size          FontSize

	// Animation properties
	IsAnimating bool
	AnimIndex   int // Current rune index
	CharDelay   int // Frames between each character
	FrameCount  int // Frame counter
}

var ColorWhite = color.RGBA{255, 255, 255, 0}

// Modified NewText constructor to initialize displayedText
func NewText(str string, x, y int, size FontSize, c color.RGBA) *Text {
	return &Text{
		content:       str,
		displayedText: str, // Start with full text displayed
		color:         &c,
		x:             x,
		y:             y,
		size:          size,
		IsAnimating:   false,
	}
}

// Update ChangeText to work with animation
func (g *Text) ChangeText(str string) {
	g.content = str
	if !g.IsAnimating {
		g.displayedText = str
	}
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

	// Use displayedText during animation, otherwise use full content
	displayText := t.content
	if t.IsAnimating {
		displayText = t.displayedText
	}

	text.Draw(screen, displayText, &text.GoTextFace{
		Size:   float64(t.size),
		Source: currFont,
	}, op)
}
func (t *Text) AnimateText(delay int) {
	t.IsAnimating = true
	t.AnimIndex = 0
	t.CharDelay = delay
	t.FrameCount = 0
	t.displayedText = ""
}

// Updated Update method to handle animation
func (t *Text) Update() error {
	if t.IsAnimating {
		t.FrameCount++

		// Time to show next character
		if t.FrameCount >= t.CharDelay {
			t.FrameCount = 0

			// If we haven't shown all characters yet
			if t.AnimIndex < len([]rune(t.content)) {
				// Get runes to handle multi-byte characters properly
				runes := []rune(t.content)
				// Add the next character
				t.displayedText = string(runes[:t.AnimIndex+1])
				t.AnimIndex++
			} else {
				// Animation complete
				t.IsAnimating = false
			}
		}
	}

	return nil
}
