package main

import (
	"github.com/hajimehoshi/ebiten"
	"image/color"
)

type Tetramino struct {
	shape [4][2]int
	x int
	y int
	color color.RGBA
}

func (t *Tetramino) Draw(screen *ebiten.Image) {
	for i := 0; i < 4; i++ {
		ebitenutil.DrawRect(screen,
			float64(t.x * BLOCKSIZE + t.shape[i][0] * BLOCKSIZE),
			float64(t.y * BLOCKSIZE + t.shape[i][1] * BLOCKSIZE),
			BLOCKSIZE,
			BLOCKSIZE,
			Black,
		)
		ebitenutil.DrawRect(screen,
			float64((t.x * BLOCKSIZE + t.shape[i][0] * BLOCKSIZE)-2),
			float64((t.y * BLOCKSIZE + t.shape[i][1] * BLOCKSIZE))-2,
			BLOCKSIZE-4,
			BLOCKSIZE-4,
			t.color,
		)
	}
}

func (t *Tetramino) Move(dx int, dy int) {
	t.x += dx
	t.y += dy
}

func (t *Tetramino) Rotate() {
	for i := 1; i < 4; i++ {
		new_x := t.shape[i][1] * -1
		new_y := t.shape[i][0]
		t.shape[i][0] = new_x
		t.shape[i][1] = new_y
	}

	out_of_bounds := 0
	for i := 0; i < 4; i++ {
		if t.x + t.shape[i][0] > 9 && t.x + t.shape[i][0] > out_of_bounds {
			out_of_bounds = (t.x + t.shape[i][0]) - 9
		} else if t.x + t.shape[i][0] < 0 && t.x + t.shape[i][0] < out_of_bounds {
			out_of_bounds = (0 - (t.x + t.shape[i][0])) * -1
		}
	}

	if out_of_bounds != 0 {
		t.x -= out_of_bounds
	}
}

func (t *Tetramino) CanRotate(board [20][10]Box) bool {
	for i := 1; i < 4; i++ {
		new_x := t.shape[i][1] * -1
		new_y := t.shape[i][0]
		if board[t.y + new_y][t.x + new_x].exists {
			return false
		}
	}
	return true
}

func (t *Tetramino) CanMove(board [20][10]Box, delta_pos [2]int) bool {
	for i := 0; i < 4; i++ {
		if t.x + t.shape[i][0] + delta_pos[0] > 9 {
			return false
		} else if t.x + t.shape[i][0] + delta_pos[0] < 0 {
			return false
		} else if board[t.y + t.shape[i][1]][t.x + t.shape[i][0] + delta_pos[0]].exists {
			return false
		}
		
	}
	return true
}

func (t *Tetramino) ShouldFreeze(board [20][10]Box) bool {
	for i := 0; i < 4; i++ {
		if t.y + t.shape[i][1] >= 19 {
			return true
		} else if board[t.y + t.shape[i][1] + 1][t.x + t.shape[i][0]].exists {
			return true
		}
	}
	return false
}

func createTetramino() Tetramino {
	switch rand.Intn(7) + 1 {
	case 1:
		return Tetramino{
			shape: [4][2]int{{0,0},{-1, 0}, {1,0}, {2,0}},
			color: LightBlue,
			x: WAIT_X-1,
			y: WAIT_Y,
		}
	case 2:
		return Tetramino{
			shape: [4][2]int{{0,0},{0, 1}, {1,0}, {1,1}},
			color: Yellow,
			x: WAIT_X,
			y: WAIT_Y,
		}
	case 3:
		return Tetramino{
			shape: [4][2]int{{0,0},{-1, 0}, {1,0}, {0,1}},
			color: Cyan,
			x: WAIT_X,
			y: WAIT_Y,
		}
	case 4:
		return Tetramino{
			shape: [4][2]int{{0,0},{0, -1}, {0,1}, {1,1}},
			color: Orange,
			x: WAIT_X,
			y: WAIT_Y,
		}
	case 5:
		return Tetramino{
			shape: [4][2]int{{0,0},{-1, 0}, {0,1}, {1,1}},
			color: Green,
			x: WAIT_X,
			y: WAIT_Y,
		}
	case 6:
		return Tetramino{
			shape: [4][2]int{{0,0},{0, 1}, {0,-1}, {-1,-1}},
			color: Blue,
			x: WAIT_X,
			y: WAIT_Y,
		}
	case 7:
		return Tetramino{
			shape: [4][2]int{{0,0},{-1, 0}, {0,-1}, {1,-1}},
			color: Red,
			x: WAIT_X,
			y: WAIT_Y,
		}
	default:
		return Tetramino{
			shape: [4][2]int{{0,0},{-1, 0}, {0,-1}, {1,-1}},
			color: Blue,
			x: WAIT_X,
			y: WAIT_Y,
		}
	}
}
