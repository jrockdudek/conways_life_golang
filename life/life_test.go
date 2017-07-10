package life

import (
	"testing"
)

func TestInitBoard(t *testing.T) {
	goodPoints := []Point{
		Point{
			x: 5,
			y: 5,
		},
		Point{
			x: 6,
			y: 5,
		},
		Point{
			x: 7,
			y: 5,
		},
	}

	badPoints := []Point{
		Point{
			x: 50,
			y: 50,
		},
		Point{
			x: 51,
			y: 50,
		},
		Point{
			x: BoardSize,
			y: 50,
		},
	}

	if err := InitBoard(goodPoints); err != nil {
		t.Error("initBoard(goodPoints) returned err: ", err)
	}

	for _, cur_point := range goodPoints {
		if Board[0][cur_point.x][cur_point.y] != true {
			t.Error("initBoard did not set all the good points to true")
		}
	}

	if err := InitBoard(badPoints); err == nil {
		t.Error("initBoard(badPoints) returned success")
	}
}
