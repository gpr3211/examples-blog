package main

import (
	"log"
	"p1-example/ui" //TODO: our own package change p1-example to ur own mod name.
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
