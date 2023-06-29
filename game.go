package main

import (
	"image/color"
	"time"
	"fmt"
	"strconv"
	"log"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/text"
)

const (
	PLAY = iota
	GAMEOVER
	RESTART
)

type Box struct {
	exists bool
	color color.RGBA
}

type Game struct {
	deltaTime        time.Time
	currentTetramino Tetramino
	nextTetramino    Tetramino
	board            [20][10]Box
	points           int
	game_over        bool
	background       *ebiten.Image
	audioMixer       AudioMixer
	playState        int
	level            int
}

func (g *Game) GameOver() {
	g.playState = GAMEOVER
	g.audioMixer.Stop()
	g.audioMixer.GameOver()
}

func (g *Game) NewTetramino() {
	g.currentTetramino = g.nextTetramino
	g.currentTetramino.x = 5
	g.currentTetramino.y = 0
	
	if g.currentTetramino.ShouldFreeze(g.board) {
		g.GameOver()
	}
	
	g.nextTetramino = createTetramino()
}

func (g *Game) FreezeTetramino() {
	for i := 0; i < 4; i++ {
		x := g.currentTetramino.x + g.currentTetramino.shape[i][0]
		y := g.currentTetramino.y + g.currentTetramino.shape[i][1]
		g.board[y][x].exists = true
		g.board[y][x].color = g.currentTetramino.color
	}
}

func (g *Game) FastForward() {
	for !g.currentTetramino.ShouldFreeze(g.board) {
		g.currentTetramino.Move(0, 1)
	}
	
	g.FreezeTetramino()
}

func (g *Game) ClearLines() int {
	clearedLines := 0
	for i := 19; i >= 0; i-- {
		cleared := true
		for j := 0; j < 10; j++ {
			if !g.board[i][j].exists {
				cleared = false
			}
		}

		if cleared {
			g.audioMixer.ClearLine()
			clearedLines++
			for j := i-1; j > 0; j-- {
				for k := 0; k < 10; k++ {
					g.board[j+1][k] = g.board[j][k]
				}
			}

			for j := 0; j < 10; j++ {
				g.board[0][j].exists = false
			}
			i++			
		}
	}

	return CalculateScore(clearedLines)
}

func (g *Game) ManageLevels() {
	if g.points / LEVEL_BOUNDARY + 1 > g.level {
		g.audioMixer.LevelUp()
		g.level++
	}
}

func (g *Game) ManageAudio() {
	if !g.audioMixer.IsPlaying() {
		g.audioMixer.Restart()
	}
}

func (g *Game) ManageInput() error {
	
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return fmt.Errorf("game is interrupted")
	} else if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		if g.currentTetramino.CanMove(g.board, [2]int{-1, 0}) {
			g.currentTetramino.Move(-1, 0)
		}
	} else if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		if g.currentTetramino.CanMove(g.board, [2]int{1, 0}) {
			g.currentTetramino.Move(1, 0)
		}
	} else if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		g.currentTetramino.Rotate()
	} else if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		g.FastForward()
	}
	
	return nil
}

func (g *Game) UpdatePlaystate() {
	g.ManageAudio()
	g.ManageInput()

	currentTime := time.Now()
	if currentTime.Sub(g.deltaTime) > time.Duration(SECOND/(START_SPEED + g.level)) {
		if g.currentTetramino.ShouldFreeze(g.board) {
			g.FreezeTetramino()
			g.NewTetramino()
		}
		g.currentTetramino.Move(0, 1)
		g.deltaTime = currentTime
	}

	g.points += g.ClearLines()
	g.ManageLevels()
}

func (g *Game) UpdateGameOverState() {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		
		if x > 5*BLOCKSIZE && x < 15*BLOCKSIZE && y > 8* BLOCKSIZE && y < 10*BLOCKSIZE {
			g.playState = RESTART
			g.points = 0
		}
	}
}

func (g *Game) Update(screen *ebiten.Image) error {
	switch g.playState {
	case PLAY:
		g.UpdatePlaystate()
	case RESTART:
		g.board = InitBoard()
		g.audioMixer.Play()
		g.playState = PLAY
		g.level = 1
	case GAMEOVER:
		g.UpdateGameOverState()
	}


	return nil
}


func (g *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(g.background, nil)
	g.currentTetramino.Draw(screen, false)
	g.DrawBoard(screen)
	g.nextTetramino.Draw(screen, true)
	score_offset := CalculateScoreOffset(g.points)
	text.Draw(screen, strconv.Itoa(g.points), regularFont, 650-score_offset*16, 170, Black)

	if g.playState == GAMEOVER {
		g.DrawGameOverBox(screen)
	}
}

func (g *Game) DrawBoard(screen *ebiten.Image) {
	for i := 0; i < 20; i++ {
		for j := 0; j < 10; j++ {
			if g.board[i][j].exists {
				ebitenutil.DrawRect(screen,
					float64(j * BLOCKSIZE),
					float64(i * BLOCKSIZE),
					BLOCKSIZE,
					BLOCKSIZE,
					Black)

				ebitenutil.DrawRect(screen,
					float64(j * BLOCKSIZE - 2),
					float64(i * BLOCKSIZE - 2),
					BLOCKSIZE - 4,
					BLOCKSIZE - 4,
					g.board[i][j].color)
			}
		}
	}	
}

func (g *Game) DrawGameOverBox(screen *ebiten.Image) {
	ebitenutil.DrawRect(screen,
		float64(5 * BLOCKSIZE),
		float64(5 * BLOCKSIZE),
		BLOCKSIZE*10,
		BLOCKSIZE*5,
		Black)

	ebitenutil.DrawRect(screen,
		float64(5 * BLOCKSIZE)+5,
		float64(5 * BLOCKSIZE)+5,
		BLOCKSIZE*10-10,
		BLOCKSIZE*5-10,
		DelftBlue)

	text.Draw(screen, "GAME OVER", regularFont, 6*BLOCKSIZE, 7*BLOCKSIZE, Plum)
	text.Draw(screen, "PLAY AGAIN", regularFont, 6*BLOCKSIZE, 9*BLOCKSIZE, Plum)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	if outsideWidth < WIDTH {
		outsideWidth = WIDTH
	}
	if outsideHeight < HEIGHT {
		outsideHeight = HEIGHT
	}
	return outsideWidth, outsideHeight
}

func CreateGame() Game {
	game := Game{
		currentTetramino: createTetramino(),
		nextTetramino: createTetramino(),
		deltaTime: time.Now(),
		board: InitBoard(),
		game_over: false,
		playState: PLAY,
		level: 1,
	}

	game.currentTetramino.x = 5
	game.currentTetramino.y = 0

	image, _, err := ebitenutil.NewImageFromFile("assets/graphics/background.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}
	game.background = image
	game.audioMixer = CreateAudioPlayer()
	game.audioMixer.Play()
	
	return game
}

func CalculateScore(clearedLines int) int {
	switch clearedLines {
	case 1:
		return 100
	case 2:
		return 300
	case 3:
		return 500
	case 4:
		return 800
	default:
		return 0
	}	
}

func CalculateScoreOffset(score int) int {
	if score > 100000 {
		return 5
	} else if score > 10000 {
		return 4
	} else if score > 1000 {
		return 3
	} else if score > 100 {
		return 2
	}

	return 0
}

func InitBoard() [20][10]Box {
	var board [20][10]Box
	
	for i := 0; i < 20; i++ {
		for j := 0; i < 10; i++ {
			board[i][j].exists = false
			board[i][j].color = Black
		}
	}

	return board
}
