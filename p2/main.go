package main

import (
	"log"

	"github.com/gpr3211/examples-blog/p2/ui"
	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	UI               *ui.Text
	Window           *ui.CodeWindow
	animationStarted bool
}

func (g *Game) Update() error {
	// Start animation once when the game begins
	//	if !g.animationStarted {
	//		g.UI.AnimateText(5) // 5 frames per character
	//		g.animationStarted = true
	//	}

	//	g.UI.Update()
	g.Window.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	//	g.UI.Draw(screen)
	g.Window.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640 * 2, 480 * 2
}

func main() {
	ebiten.SetWindowSize(640*2, 480*2)
	ebiten.SetWindowTitle("Hello, World!")
	//	hello := ui.NewText("Hello, World!", 100, 100, ui.LargeFontSize, ui.ColorWhite)
	codeWin := ui.NewCodeWindow(10, 10, ui.LargeFontSize, ui.ColorWhite)

	game := &Game{
		//		UI:               hello,
		Window:           codeWin,
		animationStarted: false,
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
