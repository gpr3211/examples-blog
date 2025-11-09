package main

import (
	"log"

	"github.com/gpr3211/examples-blog/p1/ui"
	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	UI *ui.Text
}

func (g *Game) Update() error {
	g.UI.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.UI.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 480
}

func main() {
	ebiten.SetWindowSize(640*2, 480*2)
	ebiten.SetWindowTitle("Hello, World!")

	hello := ui.NewText("Hello,World!", 100, 100, ui.LargeFontSize, ui.ColorWhite)

	game := &Game{
		UI: hello,
	}
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
