package board

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Point [2]int
type Color [2]Point

type Board struct {
	Grid   [][]int
	colors []Color
}

func New(txt io.ReadCloser) (*Board, error) {
	board := &Board{}

	r := bufio.NewReader(txt)

	sizeString, err := r.ReadString('\n')
	if err != nil {
		err = fmt.Errorf("error reading input %s", err)
		return board, err
	}
	lines, cols, err := getSize(sizeString)
	if err != nil {
		return board, err
	}

	fmt.Printf("board of %d lines and %d cols\n", lines, cols)
	board.Grid = make([][]int, lines)
	for i := 0; i < cols; i++ {
		board.Grid[i] = make([]int, cols)
	}

	index := 1
	for line := ""; ; line, err = r.ReadString('\n') {
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

	return board, nil
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

func insertPoints(board *Board, line string, index int) error {
	if line == "" {
		return nil
	}

	badFormatErr := fmt.Errorf("Bad format, lines should indicate the positions of 2 points (e.g. '0,0 0,3')")
	points := strings.Split(strings.Trim(line, "\n"), " ")
	if len(points) != 2 {
		fmt.Println(points)
		return badFormatErr
	}

	c := Color{}
	for x, point := range points {
		coords := strings.Split(point, ",")
		if len(coords) != 2 {
			fmt.Println(coords)
			return badFormatErr
		}

		i, err := strconv.Atoi(coords[0])
		j, err2 := strconv.Atoi(coords[1])

		// Check points are valid coordinates whithin specified board size
		if err != nil || err2 != nil || i < 0 || i > len(board.Grid) || j < 0 || j > len(board.Grid[0]) {
			fmt.Println(err, err2, i, j)
			return badFormatErr
		}
		board.Grid[i][j] = index
		p := Point{i, j}
		c[x] = p

	}
	board.colors = append(board.colors, c)
	return nil
}

func (b *Board) Clone() *Board {
	newBoard := &Board{}

	// colors list doesn't change so we can use the same pointer
	newBoard.colors = b.colors
	lines := len(b.Grid)
	cols := len(b.Grid[0])

	newBoard.Grid = make([][]int, lines)

	for i := range newBoard.Grid {
		newBoard.Grid[i] = make([]int, cols)
	}

	for j := range b.Grid {
		for k := range b.Grid[j] {
			newBoard.Grid[j][k] = b.Grid[j][k]
		}
	}

	return newBoard
}

func (b *Board) String() string {
	return b.GridString() + b.ColorsString()
}

func (b *Board) GridString() string {
	s := ""
	printdelimiter := func() {
		for _ = range b.Grid[0] {
			s += "+---"
		}
		s += "+\n"
	}

	printdelimiter()
	for i := range b.Grid {
		for j := range b.Grid[i] {
			s += fmt.Sprintf("| %d ", b.Grid[i][j])
		}
		s += "|\n"
		printdelimiter()
	}
	return s
}

func (b *Board) ColorsString() string {
	s := ""
	for _, c := range b.colors {
		s += fmt.Sprintf("(%d,%d)-(%d,%d)\n", c[0][0], c[0][1], c[1][0], c[1][1])
	}
	return s
}
