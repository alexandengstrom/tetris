package main

import (
	"fmt"
	"image/color"
	"time"
	"math/rand"

	"github.com/hajimehoshi/ebiten"
	"golang.org/x/image/font"
)

const (
	BLOCKSIZE = 45
	WIDTH = BLOCKSIZE * 20
	HEIGHT = BLOCKSIZE * 20
	SECOND = 1000000000
	WAIT_X = 15
	WAIT_Y = 8
	FONT_PATH = "assets/fonts/font.ttf"
	FONT_SIZE = 50
	START_SPEED = 4
)

var (
	LightGray = color.RGBA{211, 211, 211, 255}
	DarkGray = color.RGBA{169, 169, 169, 255}
	Red = color.RGBA{255, 0, 0, 255}
	LightBlue = color.RGBA{135, 206, 235, 255}
	Blue = color.RGBA{0, 0, 255, 255}
	Yellow = color.RGBA{255, 255, 0, 255}
	Cyan = color.RGBA{0, 255, 255, 255}
	Orange = color.RGBA{255, 165, 0, 255}
	Green = color.RGBA{50, 205, 50, 255}
	Black = color.RGBA{0, 0, 0, 255}
	regularFont font.Face
)

func main() {
	rand.Seed(time.Now().UnixNano())
	ebiten.SetWindowSize(WIDTH, HEIGHT)
	ebiten.SetWindowTitle("Tetris")
	game := CreateGame()
	if err := ebiten.RunGame(&game); err != nil {
		fmt.Println(err)
	}
}
