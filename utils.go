package main

import (
	"log"
	"io/ioutil"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

func init() {
	tt, err := ebitenutil.OpenFile(FONT_PATH)
	if err != nil {
		log.Fatal(err)
	}

	fontBytes, err := ioutil.ReadAll(tt)
	if err != nil {
		log.Fatal(err)
	}

	ttfont, err := truetype.Parse(fontBytes)
	if err != nil {
		log.Fatal(err)
	}

	regularFont = truetype.NewFace(ttfont, &truetype.Options{
		Size:    FONT_SIZE,
		DPI:     72,
		Hinting: font.HintingFull,
	})
}










