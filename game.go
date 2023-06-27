package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
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
	mainPlayer *audio.Player
	effectPlayer *audio.Player
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
		g.game_over = true
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
			g.effectPlayer.Rewind()
			g.effectPlayer.Play()
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

	switch cleared_lines {
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


func (g *Game) Update(screen *ebiten.Image) error {
	if g.game_over { return nil }
	
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
	fmt.Println(g.mainPlayer.IsPlaying())
	if !g.mainPlayer.IsPlaying() {
		g.mainPlayer.Rewind()
		g.mainPlayer.Play()
	}

	if g.game_over { g.mainPlayer.Pause() }

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

	g.next_tetramino.Draw(screen)
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
		game_over: false,
	}

	game.current_tetramino.x = 5
	game.current_tetramino.y = 0

	for i := 0; i < 20; i++ {
		for j := 0; i < 10; i++ {
			game.board[i][j].exists = false
			game.board[i][j].color = Red
		}
	}

	image, _, err := ebitenutil.NewImageFromFile("assets/graphics/background.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}
	game.background = image
	game.mainPlayer, game.effectPlayer = CreateAudioPlayer()
	game.mainPlayer.Play()

	//opts := &ebiten.DrawImageOptions{}
	//opts.GeoM.Translate(300, 300)
	
	return game
}
