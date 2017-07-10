// _command_line_app project main.go
package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"life/life"
	"os"
	"time"
)

const (
	defaultInput     = "start.txt"
	defaultBoardSize = 80
	defaultGifOutput = "output.gif"
	defaultDelay     = 15
)

var startingInput = flag.String("i", defaultInput, "The file to read your starting points from")
var boardSize = flag.Int("b", defaultBoardSize, "How big of a baord to use(default is 80)")
var generations = flag.Int("g", 0, "Number of generations to run(default is 0, or infinite)")
var gifOutput = flag.Bool("j", false, "Output the resulting life to a gif")
var gifFilename = flag.String("o", defaultGifOutput, "Name of file to output gif to(only valid with -j)")
var gifDelay = flag.Int("d", defaultDelay, "Delay between frames of the GIF")

var palette = []color.Color{color.White, color.Black}

func main() {
	flag.Parse()
	a := new([2][][]bool)
	for i := range a {
		a[i] = make([][]bool, *boardSize)
		for j := range a[i] {
			a[i][j] = make([]bool, *boardSize)
		}
	}
	board := life.Board(*a)

	if err := board.InitBoardFromFile(*startingInput); err != nil {
		fmt.Fprintf(os.Stderr, "InitBoardFromFile: %v\n", err)
		os.Exit(1)
	}

	// Main loop
	var boardIndex int
	// We can't have infinite generations for a gif, so limit it if they put 0
	if *gifOutput && *generations == 0 {
		*generations = 500
	}
	anim := gif.GIF{LoopCount: *generations}
	if !*gifOutput {
		PrintBoardTerminal(board, boardIndex, 0)
		time.Sleep(1 * time.Second)
	} else {
		PrintBoardGif(board, boardIndex, &anim)
	}
	for cur_gen := 1; cur_gen < *generations || *generations == 0; cur_gen++ {
		newBoardIndex := (boardIndex + 1) % 2
		ch := make(chan struct{})
		for cur_y := 0; cur_y < len(board[0][0]); cur_y++ {
			for cur_x := 0; cur_x < len(board[0]); cur_x++ {
				go func(func_x, func_y int) {
					board[newBoardIndex][func_x][func_y] = board.CheckCell(boardIndex, func_x, func_y)
					ch <- struct{}{}
				}(cur_x, cur_y)
			}
		}
		for i := 0; i < (len(board[0][0]) * len(board[0])); i++ {
			<-ch
		}
		boardIndex = newBoardIndex
		if !*gifOutput {
			PrintBoardTerminal(board, boardIndex, cur_gen)
			time.Sleep(500 * time.Millisecond)
		} else {
			PrintBoardGif(board, boardIndex, &anim)
		}
	}
	if *gifOutput {
		output_file, err := os.Create(*gifFilename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "gif.EncodeAll: %v\n", err)
			os.Exit(1)
		}

		defer func() {
			if err := output_file.Close(); err != nil {
				fmt.Fprintf(os.Stderr, "gif.EncodeAll: %v\n", err)
				os.Exit(1)
			}
		}()

		if err := gif.EncodeAll(output_file, &anim); err != nil {
			fmt.Fprintf(os.Stderr, "gif.EncodeAll: %v\n", err)
			os.Exit(1)
		}
	}
}

func PrintBoardTerminal(b life.Board, boardIndex int, cur_gen int) {
	fmt.Printf("===================== Generation: %d =======================\n", cur_gen)
	for y := 0; y < len(b[0][0]); y++ {
		for x := 0; x < len(b[0]); x++ {
			cell_print := " "
			if b[boardIndex][x][y] {
				cell_print = "O"
			}
			fmt.Printf(cell_print)
		}
		fmt.Println()
	}
	fmt.Println("===========================================================")
}

func PrintBoardGif(b life.Board, boardIndex int, anim *gif.GIF) {
	pixels_per_square := 6
	rect := image.Rect(0, 0, *boardSize*pixels_per_square, *boardSize*pixels_per_square)
	img := image.NewPaletted(rect, palette)
	for y := 0; y < len(b[0][0]); y++ {
		for x := 0; x < len(b[0]); x++ {
			if b[boardIndex][x][y] {
				for i := 0; i < pixels_per_square; i++ {
					for k := 0; k < pixels_per_square; k++ {
						img.SetColorIndex((x*pixels_per_square)+i, (y*pixels_per_square)+k, 1)
					}
				}
			}
		}
	}
	anim.Delay = append(anim.Delay, *gifDelay)
	anim.Image = append(anim.Image, img)
}
