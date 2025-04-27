package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image/color"
)

type CodeWindow struct {
	x, y          int
	width, height int
	bgColor       color.Color
	borderColor   color.Color
	text          *Text
	content       []rune
	cursorPos     int
	frame         int
}

func NewCodeWindow(x, y int, size FontSize, c color.RGBA) *CodeWindow {
	return &CodeWindow{
		text:        NewText("", x, y, size, c),
		content:     []rune{},
		cursorPos:   0,
		frame:       0,
		width:       600,
		height:      600,
		bgColor:     color.RGBA{30, 30, 30, 225},
		borderColor: ColorWhite,
	}
}

func (cw *CodeWindow) Update() error {
	cw.frame++
	cw.text.Update()

	for _, r := range ebiten.AppendInputChars(nil) {
		if r == '\b' {
			//			cw.backspace()
		} else if r == '\r' || r == '\n' {
			// TODO:
			// add lines. separate words.
		} else if r >= 32 { // Printable characters
			cw.insertRune(r)
		}
	}

	// Handle arrow keys
	if repeatingKeyPressed(ebiten.KeyBackspace) && cw.cursorPos > 0 {
		cw.backspace()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) && cw.cursorPos > 0 {
		cw.cursorPos--
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) && cw.cursorPos < len(cw.content) {
		cw.cursorPos++
	}

	// Rebuild displayed text
	cw.updateText()

	return nil
}

func (cw *CodeWindow) Draw(screen *ebiten.Image) {
	// Background
	bg := ebiten.NewImage(cw.width, cw.height)
	bg.Fill(cw.bgColor)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(cw.x), float64(cw.y))
	screen.DrawImage(bg, op)

	// Border (simple 1px border)
	borderColor := cw.borderColor
	borderW, borderH := cw.width, cw.height

	// Top
	line := ebiten.NewImage(borderW, 1)
	line.Fill(borderColor)
	opTop := &ebiten.DrawImageOptions{}
	opTop.GeoM.Translate(float64(cw.x), float64(cw.y))
	screen.DrawImage(line, opTop)

	// Bottom
	opBottom := &ebiten.DrawImageOptions{}
	opBottom.GeoM.Translate(float64(cw.x), float64(cw.y+borderH-1))
	screen.DrawImage(line, opBottom)

	// Left
	vline := ebiten.NewImage(1, borderH)
	vline.Fill(borderColor)
	opLeft := &ebiten.DrawImageOptions{}
	opLeft.GeoM.Translate(float64(cw.x), float64(cw.y))
	screen.DrawImage(vline, opLeft)

	// Right
	opRight := &ebiten.DrawImageOptions{}
	opRight.GeoM.Translate(float64(cw.x+borderW-1), float64(cw.y))
	screen.DrawImage(vline, opRight)

	// Text inside
	cw.text.Draw(screen)
}

func (cw *CodeWindow) insertRune(r rune) {
	if cw.cursorPos >= 0 && cw.cursorPos <= len(cw.content) {
		before := cw.content[:cw.cursorPos]
		after := cw.content[cw.cursorPos:]
		cw.content = append(append(before, r), after...)
		cw.cursorPos++
	}
}

func (cw *CodeWindow) backspace() {
	if cw.cursorPos > 0 {
		before := cw.content[:cw.cursorPos-1]
		after := cw.content[cw.cursorPos:]
		cw.content = append(before, after...)
		cw.cursorPos--
	}
}

func (cw *CodeWindow) updateText() {
	display := string(cw.content)
	if (cw.frame/30)%2 == 0 { // Blinking every half second (assuming 60fps)
		// Insert blinking cursor (visually a "_")
		runes := cw.content
		cursor := "_"
		before := string(runes[:cw.cursorPos])
		after := string(runes[cw.cursorPos:])
		display = before + cursor + after
	}
	cw.text.ChangeText(display)
}

// Helper functions
//

func repeatingKeyPressed(key ebiten.Key) bool {
	const (
		delay    = 30
		interval = 3
	)
	d := inpututil.KeyPressDuration(key)
	if d == 1 {
		return true
	}
	if d >= delay && (d-delay)%interval == 0 {
		return true
	}
	return false
}

func runeWidth(r rune, fontSize FontSize) int {
	if r == '\t' {
		return 4 * int(fontSize)
	}
	return int(fontSize) // Approx: 1 fontSize = 1 char width
}
