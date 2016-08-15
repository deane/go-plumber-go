package board

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

var (
	ColorFuncs = []func(string, ...interface{}) string{
		colorwrapper(color.New(color.BgBlack, color.FgWhite)),
		colorwrapper(color.New(color.BgRed, color.FgWhite)),
		colorwrapper(color.New(color.BgGreen, color.FgBlack)),
		colorwrapper(color.New(color.BgYellow, color.FgBlack)),
		colorwrapper(color.New(color.BgBlue, color.FgWhite)),
		colorwrapper(color.New(color.BgMagenta, color.FgWhite)),
		colorwrapper(color.New(color.BgCyan, color.FgBlack)),
		colorwrapper(color.New(color.BgWhite, color.FgBlack)),
	}
)

type Point [2]int
type Color []Point

type Board struct {
	Grid   [][]int
	colors []Color
}

func colorwrapper(c *color.Color) func(string, ...interface{}) string {
	return func(s string, args ...interface{}) string {
		return c.SprintFunc()(fmt.Sprintf(s, args...))
	}
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

	index := 0
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
	for _, point := range points {
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
		c = append(c, p)

	}
	board.colors = append(board.colors, c)
	return nil
}

func (b *Board) Clone() *Board {
	newBoard := &Board{}

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
	for _, color := range b.colors {
		newBoard.colors = append(newBoard.colors, color)
	}

	return newBoard
}

func (b *Board) ColorCell(colorIndex, line, col int) error {
	if colorIndex < 0 || colorIndex >= len(b.colors) {
		return errors.New("color index out of range")
	}
	if line < 0 || line >= len(b.Grid) {
		return errors.New("X out of range")
	}
	if col < 0 || col >= len(b.Grid[0]) {
		return errors.New("Y out of range")
	}

	if b.Grid[line][col] != 0 {
		return errors.New("Cell already occupied")
	}

	sliceLen := len(b.colors[colorIndex])
	c := make([]Point, sliceLen, sliceLen+1)
	for i, p := range b.colors[colorIndex] {
		c[i] = p
	}
	updatedC := append(c[:len(c)-1], Point{line, col}, c[len(c)-1])
	if !AreAllAjacent(updatedC[:len(c)]) {
		return fmt.Errorf("Cells are not ajacent: %v", updatedC[:len(c)])
	}
	b.Grid[line][col] = colorIndex + 1
	b.colors[colorIndex] = updatedC

	return nil
}

func (b *Board) Solved() bool {
	//check the grid is full
	for i := 0; i < len(b.Grid); i++ {
		for j := 0; j < len(b.Grid[0]); j++ {
			if b.Grid[i][j] == 0 {
				return false
			}
		}
	}
	for _, c := range b.colors {
		if !AreAllAjacent(c) {
			fmt.Println(c)
			return false
		}
	}

	return true
}

func AreAjacent(point, nextPoint Point) bool {

	dx := point[0] - nextPoint[0]
	dy := point[1] - nextPoint[1]
	if dx*dx > 1 || dy*dy > 1 || dx*dx+dy*dy != 1 {
		return false
	}
	return true
}

func AreAllAjacent(c Color) bool {
	for i, point := range c {
		if i == len(c)-1 {
			break
		}
		nextPoint := c[i+1]
		if !AreAjacent(point, nextPoint) {
			return false
		}
	}
	return true
}

func AjacentToAny(point Point, c Color) bool {
	for _, p := range c {
		if AreAjacent(p, point) {
			return true
		}
	}
	return false
}

func (b *Board) String() string {
	return b.GridString() + b.ColorsString()
}

func (b *Board) GridString() string {
	s := "\n"
	printdelimiter := func() {
		for _ = range b.Grid[0] {
			s += "+---"
		}
		s += "+\n"
	}

	printdelimiter()
	for i := range b.Grid {
		for j := range b.Grid[i] {
			val := b.Grid[i][j]
			s += "|"
			s += ColorFuncs[val%8](" %d ", val)
		}
		s += "|\n"
		printdelimiter()
	}
	return s
}

func (b *Board) ColorsString() string {
	s := ""
	for _, c := range b.colors {
		for i, point := range c {
			if i == len(c)-1 && !AreAllAjacent(c) {
				s += "[???]->"
			}
			s += fmt.Sprintf("(%d,%d)", point[0], point[1])
			if i != len(c)-1 {
				s += "->"
			}
		}
		s += "\n"
	}
	return s
}
