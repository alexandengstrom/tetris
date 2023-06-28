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
	//"github.com/hajimehoshi/ebiten/audio"
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
	delta_time time.Time
	current_tetramino Tetramino
	next_tetramino Tetramino
	board [20][10]Box
	points int
	game_over bool
	background *ebiten.Image
	audioMixer AudioMixer
	playState int
}

func (g *Game) NewTetramino() {
	for i := 0; i < 4; i++ {
		x := g.current_tetramino.x + g.current_tetramino.shape[i][0]
		y := g.current_tetramino.y + g.current_tetramino.shape[i][1]
		g.board[y][x].exists = true
		g.board[y][x].color = g.current_tetramino.color
	}
	g.current_tetramino = g.next_tetramino
	g.current_tetramino.x = 5
	g.current_tetramino.y = 0
	if g.current_tetramino.ShouldFreeze(g.board) {
		g.playState = GAMEOVER
	}
	g.next_tetramino = createTetramino()
}

func (g *Game) FastForward() {
	for !g.current_tetramino.ShouldFreeze(g.board) {
		g.current_tetramino.Move(0, 1)
	}
}

func (g *Game) ClearLines() int {
	cleared_lines := 0
	for i := 19; i >= 0; i-- {
		cleared := true
		for j := 0; j < 10; j++ {
			if !g.board[i][j].exists {
				cleared = false
			}
		}

		if cleared {
			g.audioMixer.ClearLine()
			cleared_lines++
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

	return CalculateScore(cleared_lines)
}

func (g *Game) ManageAudio() {
	if g.game_over {
		g.audioMixer.Stop()
	} else if !g.audioMixer.IsPlaying() {
		g.audioMixer.Restart()
	}
}

func (g *Game) Update(screen *ebiten.Image) error {
	switch g.playState {
	case PLAY:
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			return fmt.Errorf("game is interrupted")
		} else if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
			if g.current_tetramino.CanMove(g.board, [2]int{-1, 0}) {
				g.current_tetramino.Move(-1, 0)
			}
		} else if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
			if g.current_tetramino.CanMove(g.board, [2]int{1, 0}) {
				g.current_tetramino.Move(1, 0)
			}
		} else if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
			g.current_tetramino.Rotate()
		} else if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
			g.FastForward()
		}

		current_time := time.Now()
		if current_time.Sub(g.delta_time) > SECOND/START_SPEED {
			if g.current_tetramino.ShouldFreeze(g.board) {
				g.NewTetramino()
			}
			g.current_tetramino.Move(0, 1)
			g.delta_time = current_time
		}

		g.points += g.ClearLines()
		g.ManageAudio()
	case RESTART:
		g.board = InitBoard()
		g.playState = PLAY
	case GAMEOVER:
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			x, y := ebiten.CursorPosition()
			
			
			fmt.Println(x, y)
			
			if x > 5*BLOCKSIZE && x < 15*BLOCKSIZE && y > 8* BLOCKSIZE && y < 10*BLOCKSIZE {
				g.playState = RESTART
				g.points = 0
			}
		}
	}


	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(g.background, nil)
	g.current_tetramino.Draw(screen)

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

	g.next_tetramino.DrawQueue(screen)
	score_offset := 0
	if g.points > 100000 {
		score_offset = 5
	} else if g.points > 10000 {
		score_offset = 4
	} else if g.points > 1000 {
		score_offset = 3
	} else if g.points > 100 {
		score_offset = 2
	}
	text.Draw(screen, strconv.Itoa(g.points), regularFont, 650-score_offset*16, 170, Black)

	if g.playState == GAMEOVER {
		ebitenutil.DrawRect(screen,
			float64(5 * BLOCKSIZE),
			float64(5 * BLOCKSIZE),
			BLOCKSIZE*10,
			BLOCKSIZE*5,
			Black)

		text.Draw(screen, "GAME OVER", regularFont, 6*BLOCKSIZE, 7*BLOCKSIZE, Red)
		text.Draw(screen, "PLAY AGAIN", regularFont, 6*BLOCKSIZE, 9*BLOCKSIZE, Red)
	}
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
		current_tetramino: createTetramino(),
		next_tetramino: createTetramino(),
		delta_time: time.Now(),
		board: InitBoard(),
		game_over: false,
		playState: PLAY,
	}

	game.current_tetramino.x = 5
	game.current_tetramino.y = 0

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
