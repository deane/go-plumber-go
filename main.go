package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

var fileName = flag.String("file", "board.txt", "go-plumber-go board file to solve")

func main() {
	flag.Parse()
	f, err := os.Open(*fileName)
	if err != nil {
		fmt.Printf("can't open the file %s, ERROR: %s\n", *fileName, err)
		return
	}
	defer f.Close()
	file := bufio.NewReader(f)

	sizeString, err := file.ReadString('\n')
	if err != nil {
		fmt.Printf("error reading the file %s, ERROR: %s\n", *fileName, err)
		return
	}
	lines, cols, err := getSize(sizeString)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("board of %d lines and %d cols\n", lines, cols)
	board := initBoard(lines, cols)

	index := 1
	for line := ""; ; line, err = file.ReadString('\n') {
		readErr := insertPoints(board, line, index)
		if err != nil || readErr != nil {
			err = readErr
			break
		}
		index++
	}
	if err != io.EOF && err != nil {
		fmt.Println("Error reading file:", err)
	}

	printBoard(board)
}

func getSize(s string) (int, int, error) {
	badFormatErr := fmt.Errorf("Bad format, first line should indicate the size of the board (e.g. '5,5')")
	split := strings.Split(strings.Trim(s, "\n"), ",")
	if len(split) != 2 {
		fmt.Println(split)
		return 0, 0, badFormatErr
	}

	lines, err := strconv.Atoi(split[0])
	cols, err2 := strconv.Atoi(split[1])
	if err != nil || err2 != nil {
		fmt.Println(err, err2)
		return 0, 0, badFormatErr
	}

	return lines, cols, nil

}

func initBoard(lines, cols int) [][]int {
	board := make([][]int, lines)
	for i := 0; i < cols; i++ {
		board[i] = make([]int, cols)
	}
	return board
}

func printBoard(board [][]int) {
	printdelimiter := func() {
		for _ = range board[0] {
			fmt.Print("+---")
		}
		fmt.Print("+\n")
	}

	printdelimiter()
	for i := range board {
		for j := range board[i] {
			fmt.Printf("| %d ", board[i][j])
		}
		fmt.Print("|\n")
		printdelimiter()

	}

}

func insertPoints(board [][]int, line string, index int) error {
	if line == "" {
		return nil
	}

	badFormatErr := fmt.Errorf("Bad format, lines should indicate the positions of 2 points (e.g. '0,0 0,3')")
	points := strings.Split(strings.Trim(line, "\n"), " ")
	if len(points) != 2 {
		fmt.Println(points)
		return badFormatErr
	}

	for _, point := range points {
		coords := strings.Split(point, ",")
		if len(coords) != 2 {
			fmt.Println(coords)
			return badFormatErr
		}

		i, err := strconv.Atoi(coords[0])
		j, err2 := strconv.Atoi(coords[1])

		// Check points are valid coordinates whithin specified board size
		if err != nil || err2 != nil || i < 0 || i > len(board) || j < 0 || j > len(board[0]) {
			fmt.Println(err, err2, coords)
			return badFormatErr
		}
		board[i][j] = index

	}
	return nil
}
