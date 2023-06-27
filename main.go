package main

import (
	"fmt"
	"image/color"
	"time"
	"math/rand"

	"github.com/hajimehoshi/ebiten"
	//"github.com/hajimehoshi/ebiten/inpututil"
	//"github.com/hajimehoshi/ebiten/ebitenutil"
	//"github.com/hajimehoshi/ebiten/text"
	//"github.com/alexandengstrom/tetris/tetramino"
	//"github.com/alexandengstrom/tetris/game"
	//"github.com/alexandengstrom/tetris/utils"
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

var LightGray = color.RGBA{211, 211, 211, 255}
var DarkGray = color.RGBA{169, 169, 169, 255}
var Red = color.RGBA{255, 0, 0, 255}
var LightBlue = color.RGBA{135, 206, 235, 255}
var Blue = color.RGBA{0, 0, 255, 255}
var Yellow = color.RGBA{255, 255, 0, 255}
var Cyan = color.RGBA{0, 255, 255, 255}
var Orange = color.RGBA{255, 165, 0, 255}
var Green = color.RGBA{50, 205, 50, 255}
var Black = color.RGBA{0, 0, 0, 255}
var regularFont font.Face





func main() {
	rand.Seed(time.Now().UnixNano())
	ebiten.SetWindowSize(WIDTH, HEIGHT)
	ebiten.SetWindowTitle("Tetris")
	game := CreateGame()
	if err := ebiten.RunGame(&game); err != nil {
		fmt.Println(err)
	}
}
