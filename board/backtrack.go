package board

import (
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"sync"
	"time"
)

var (
	//options
	Canonical   = false
	ShowResults = false
	DetectDead  = false

	//stats
	leaves       = 0
)

type Move struct {
	Color int
	Point Point
}

// Thread Safe
type BoardList struct {
	list []*Board
	sync.Mutex
}

func (bl *BoardList) Append(b *Board) {
	bl.Lock()
	bl.list = append(bl.list, b)
	bl.Unlock()
}

func (bl *BoardList) Prefix(b *Board) {
	bl.Lock()
	bl.list = append([]*Board{b}, bl.list...)
	bl.Unlock()
}

func (bl *BoardList) Len() int {
	bl.Lock()
	defer bl.Unlock()
	return len(bl.list)
}

func (bl *BoardList) New(b *Board) {
	bl.Lock()
	bl.list = []*Board{b}
	bl.Unlock()
}

func ApplyMove(b *Board, m Move) error {
	return b.ColorCell(m.Color, m.Point[0], m.Point[1])
}

func Backtrack(b *Board) (solution *BoardList, err error) {
	history := 0
	solution = &BoardList{list: make([]*Board, 0)}

	s := time.Now()
	err = backtrack(b, solution, &history)
	tdelta := time.Since(s)
	log.Printf(
		"Backtrack Stats:\n%v\n%d states explored, %v per iteration, got to %d leaves",
		tdelta, history, tdelta/time.Duration(history), leaves,
	)
	if err != nil {
		log.Printf(
			"Not solved :(\n%s\nspent %v, exploring %d states\n",
			err.Error(), tdelta, history,
		)
		return
	}

	log.Printf(
		"SOLVED!!!!!! in %v, %d steps and %d states explored \n",
		tdelta, solution.Len(), history,
	)
	if ShowResults {
		for _, board := range solution.list {
			interval := 300 * time.Millisecond
			time.Sleep(interval)
			cmd := exec.Command("clear")
			cmd.Stdout = os.Stdout
			cmd.Run()
			log.Printf("\n%s", board.GridString())
		}
		time.Sleep(1 * time.Second)
	}

	return
}

func backtrack(b *Board, s *BoardList, h *int) error {
	solution := s

	moves := NextMoves(b)
	for _, move := range moves {
		newB := b.Clone()
		err := ApplyMove(newB, move)
		if err != nil {
			log.Printf("Backtrack ERROR applying move: %s", err.Error())
			log.Printf("Move:%v\nBoard:\n%s", newB.String())
			os.Exit(1)
		}
		*h += 1

		if newB.Solved() {
			solution.New(newB)
			return nil
		}

		err = backtrack(newB, solution, h)
		if err == nil {
			solution.Prefix(newB)
			return nil
		}
	}
	if len(moves) == 0 {
		leaves += 1
	}
	return fmt.Errorf("No solution Found. Explored %d states", *h)
}

// NextMoves assumes the coordinates of the board's colors are valid
func NextMoves(b *Board) []Move {
	possibleMoves := []Move{}

	for colorIndex, color := range b.colors {
		if AreAllAjacent(color) {
			continue
		}
		lastPoint := color[len(color)-2]

		around := surroundings(lastPoint)

		for _, p := range around {
			if !(inGrid(b, p) && b.Grid[p[0]][p[1]] == 0) {
				continue
			}
			legal := true
			if Canonical && len(color) > 2 && AjacentToAny(p, color[:len(color)-3]) {
				legal = false
			}
			if legal && DetectDead && findDeadCell(b, p, colorIndex+1) {
				legal = false
			}
			if legal {
				possibleMoves = append(possibleMoves, Move{Color: colorIndex, Point: p})
			}

		}

		return possibleMoves
	}

	return possibleMoves
}

// Check if coloring a given cell will create dead cells
// Take the surrounding cells, and check if their surroundings allow
// for a flow: a cell having n surrounding cells in the grid,
// can have a maximum of n-1 colors surrounding it
func findDeadCell(b *Board, p Point, color int) bool {
	for _, p2 := range surroundings(p) {
		if !inGrid(b, p2) || b.Grid[p2[0]][p2[1]] != 0 {
			continue
		}
		colors := [4]int{color, -1, -1, -1}
		NumAdjacentCells := 1 // p has to be ajacent to p2
		for _, p3 := range surroundings(p2) {
			if p3 == p || !inGrid(b, p3) {
				continue
			}
			NumAdjacentCells += 1
			c := b.Grid[p3[0]][p3[1]]
			if c == 0 {
				continue
			}
			for i, existingColor := range colors {
				if c == existingColor {
					break
				}
				if existingColor == -1 {
					colors[i] = c
					break
				}
			}
		}
		for i, v := range colors {
			if v == -1 {
				if i >= NumAdjacentCells {
					return true
				}
				break
			}
		}
	}
	return false
}

func inGrid(b *Board, p Point) bool {
	if p[0] < 0 || p[1] < 0 || p[0] >= len(b.Grid) || p[1] >= len(b.Grid[0]) {
		return false
	}
	return true
}

func surroundings(p Point) []Point {
	return []Point{{p[0] - 1, p[1]}, {p[0] + 1, p[1]}, {p[0], p[1] - 1}, {p[0], p[1] + 1}}
}

func SortColors(b *Board, reverse bool) *Board {
	res := b.Clone()

	l := []Color{} // new color list
	for _, c := range b.colors {
		set := false
		for i, c2 := range l {
			if Distance(c[0], c[len(c)-1]) > Distance(c2[0], c2[len(c2)-1]) {
				l = append(l[:i], append([]Color{c}, l[i:]...)...)
				set = true
				break
			}
		}
		if !set {
			l = append(l, c)
		}
	}
    if reverse {
		for i := 0; i < len(l)/2; i++ {
			l[i], l[len(l)-i-1] = l[len(l)-i-1], l[i]
		}
    }

	res.colors = l

	for j, c := range res.colors {
		for _, p := range c {
			res.Grid[p[0]][p[1]] = j + 1
		}
	}
	return res
}

func Distance(point1, point2 Point) float64 {
	dx := point2[0] - point1[0]
	dy := point2[1] - point1[1]

	return math.Sqrt(float64(dx*dx + dy*dy))
}
