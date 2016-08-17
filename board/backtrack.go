package board

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"
	"time"
)

var (
	Canonical   = false
	ShowResults = false
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

func Backtrack(b *Board) (solution *BoardList, history *BoardList, err error) {
	history = &BoardList{list: make([]*Board, 0)}
	solution = &BoardList{list: make([]*Board, 0)}

	s := time.Now()
	solution, history, err = backtrack(b, history, solution)
	tdelta := time.Since(s)
	log.Printf(
		"Backtrack Stats:\n%v\n%d states explored, %v per iteration",
		tdelta, history.Len(), tdelta/time.Duration(history.Len()),
	)
	if err == nil {
		log.Printf(
			"SOLVED!!!!!! in %v, %d steps and %d states explored \n",
			tdelta, solution.Len(), history.Len(),
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
	} else {
		log.Printf(
			"Not solved :( spent %v, exploring %d states\n",
			tdelta, history.Len(),
		)
	}
	return
}

func backtrack(b *Board, s, h *BoardList) (
	solution, history *BoardList, err error,
) {
	solution, history = s, h

	moves := NextMoves(b)
	for _, move := range moves {
		newB := b.Clone()
		err = ApplyMove(newB, move)
		if err != nil {
			log.Printf("Backtrack ERROR applying move: %s", err.Error())
			log.Printf("Move:%v\nBoard:\n%s", newB.String())
			return
		}
		history.Append(newB)

		if newB.Solved() {
			solution.New(newB)
			return solution, history, nil
		}

		solution, history, err = backtrack(newB, history, solution)
		if err == nil {
			solution.Prefix(newB)
			return solution, history, nil
		}
	}
	return solution, history, fmt.Errorf("No solution Found. Explored %d states", history.Len())
}

// NextMoves assumes the coordinates of the board's colors are valid
func NextMoves(b *Board) []Move {
	possibleMoves := []Move{}

	for colorIndex, color := range b.colors {
		if AreAllAjacent(color) {
			continue
		}
		lastPoint := color[len(color)-2]

		l, c := lastPoint[0], lastPoint[1]

		// go up
		if l-1 >= 0 && b.Grid[l-1][c] == 0 {
			p := Point{l - 1, c}
			if !(Canonical && len(color) > 2 && AjacentToAny(p, color[:len(color)-3])) {
				possibleMoves = append(possibleMoves, Move{Color: colorIndex, Point: p})
			}
		}

		// go down
		if l+1 < len(b.Grid) && b.Grid[l+1][c] == 0 {
			p := Point{l + 1, c}
			if !(Canonical && len(color) > 2 && AjacentToAny(p, color[:len(color)-3])) {
				possibleMoves = append(possibleMoves, Move{Color: colorIndex, Point: Point{l + 1, c}})
			}
		}

		// go left
		if c-1 >= 0 && b.Grid[l][c-1] == 0 {
			p := Point{l, c - 1}
			if !(Canonical && len(color) > 2 && AjacentToAny(p, color[:len(color)-3])) {
				possibleMoves = append(possibleMoves, Move{Color: colorIndex, Point: Point{l, c - 1}})
			}
		}

		//go right
		if c+1 < len(b.Grid[0]) && b.Grid[l][c+1] == 0 {
			p := Point{l, c + 1}
			if !(Canonical && len(color) > 2 && AjacentToAny(p, color[:len(color)-3])) {
				possibleMoves = append(possibleMoves, Move{Color: colorIndex, Point: Point{l, c + 1}})
			}
		}
		return possibleMoves
	}

	return possibleMoves
}
