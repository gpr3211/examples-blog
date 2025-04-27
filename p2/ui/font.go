package ui

import (
	"embed"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"log"
)

//go:embed assets/*
var fontAssets embed.FS

var (
	DejaVuSans *text.GoTextFaceSource
	ThaleahFat *text.GoTextFaceSource
	currFont   *text.GoTextFaceSource
)

// LoadFont laods a font to be used in text assets.
func LoadFont(name string, assets embed.FS) (*text.GoTextFaceSource, error) {
	font, err := assets.Open(fmt.Sprintf("assets/%s.ttf", name))
	if err != nil {
		return nil, err
	}
	defer font.Close()
	s, err := text.NewGoTextFaceSource(font)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func init() {

	fmt.Println("Loading fonts")
	tf, err := LoadFont("ThaleahFat", fontAssets)
	if err != nil {
		log.Println("Failed to load font")
		panic(err)
	}
	ThaleahFat = tf

	currFont = ThaleahFat
	fmt.Println("Font loaded")
}
