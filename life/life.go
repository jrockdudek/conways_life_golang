// life project main.go
package life

import (
	"errors"
	"strconv"
)

type Point struct {
	x int
	y int
}

type Board [2][][]bool

func (b Board) InitBoardFromFile(filename string) error {
	points := []Point{
		Point{
			x: 55,
			y: 40,
		},
		Point{
			x: 54,
			y: 40,
		},
		Point{
			x: 53,
			y: 40,
		},
		Point{
			x: 52,
			y: 39,
		},
		Point{
			x: 50,
			y: 40,
		},
		Point{
			x: 50,
			y: 38,
		},
		Point{
			x: 49,
			y: 40,
		},
	}

	err := b.InitBoard(points)

	return err
}

func (b Board) InitBoard(startingPoints []Point) error {
	for _, cur_point := range startingPoints {
		if cur_point.x >= len(b[0]) || cur_point.y >= len(b[0][0]) {
			err := errors.New(strconv.Itoa(cur_point.x) + ", " + strconv.Itoa(cur_point.y) + " is not on the board!")
			return err
		}
		b[0][cur_point.x][cur_point.y] = true
	}

	return nil
}

func (b Board) CheckCell(curIndex int, cur_x int, cur_y int) bool {
	var numOfAliveNearby int

	for y := cur_y - 1; y <= cur_y+1; y++ {
		for x := cur_x - 1; x <= cur_x+1; x++ {
			// If we're off the board, or on the cell we're checking don't count
			// that towards our alive neighbors
			if !(x < 0 || y < 0 || x >= len(b[0]) || y >= len(b[0][0])) {
				if x != cur_x || y != cur_y {
					if b[curIndex][x][y] {
						numOfAliveNearby++
					}
				}
			}
		}
	}

	if !b[curIndex][cur_x][cur_y] && numOfAliveNearby == 3 {
		return true
	} else if b[curIndex][cur_x][cur_y] && (numOfAliveNearby == 2 || numOfAliveNearby == 3) {
		return true
	} else {
		return false
	}
}
