package main

import (
	"log"
	"io/ioutil"
	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/audio"
	"github.com/hajimehoshi/ebiten/audio/wav"
	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
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









func CreateAudioPlayer() (*audio.Player, *audio.Player) {
	audioContext, err := audio.NewContext(44100)
	if err != nil {
		log.Fatal(err)
	}

	file, err := ebitenutil.OpenFile("assets/audio/soundtrack.wav")
	if err != nil {
		log.Fatal(err)
	}

	decodedSound, err := wav.Decode(audioContext, file)
	if err != nil {
		log.Fatal(err)
	}

	player, err := audio.NewPlayer(audioContext, decodedSound)

	file2, err := ebitenutil.OpenFile("assets/audio/clear.wav")
	if err != nil {
		log.Fatal(err)
	}

	decodedSound2, err := wav.Decode(audioContext, file2)
	if err != nil {
		log.Fatal(err)
	}

	effectplayer, err := audio.NewPlayer(audioContext, decodedSound2)
	return player, effectplayer
}
