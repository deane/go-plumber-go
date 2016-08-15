package board

import (
	"os"
	"testing"
)

func getBoard(name string, t *testing.T) *Board {

	f, err := os.Open("../" + name)

	if err != nil {
		t.Fatalf("can't open test board file: %s", err.Error())
	}
	defer f.Close()

	b, err := New(f)
	if err != nil {
		t.Fatalf("can't parse board file: %s", err.Error())
	}
	return b
}

func TestColorCell(t *testing.T) {
	b := getBoard("board.txt", t)
	err := b.ColorCell(0, 1, 0)
	if err != nil {
		t.Fatalf("can't color cell: %s", err.Error())
	}
	err = b.ColorCell(0, 1, 0)
	if err == nil {
		t.Fatal("can't color cell again")
	}

}

func TestSolved(t *testing.T) {
	b := getBoard("board.txt", t)

	if b.Solved() {
		t.Fatalf("The board:\n%s\nShouldn't be considered solved", b.String())
	}
	b.ColorCell(0, 1, 0)
	b.ColorCell(0, 2, 0)
	b.ColorCell(0, 3, 0)
	b.ColorCell(0, 4, 0)

	b.ColorCell(1, 0, 1)
	b.ColorCell(1, 1, 1)
	b.ColorCell(1, 2, 1)
	b.ColorCell(1, 4, 0)

	b.ColorCell(2, 0, 3)
	b.ColorCell(2, 1, 3)
	b.ColorCell(2, 2, 3)

	b.ColorCell(3, 2, 2)
	b.ColorCell(3, 3, 2)

	b.ColorCell(4, 2, 4)
	b.ColorCell(4, 3, 4)
	b.ColorCell(4, 4, 4)

	if !b.Solved() {
		t.Fatalf("The board:\n%s\nShould be considered solved", b.String())
	}
	t.Logf("\n%s", b.String())
}

func TestNextMoves(t *testing.T) {
	b := getBoard("board.txt", t)

	moves := NextMoves(b)

	for _, m := range moves {
		newB := b.Clone()
		err := newB.ColorCell(m.Color, m.Point[0], m.Point[1])
		if err != nil {
			t.Fatalf("\n%sShouldn't have %v in next moves\n%s", b.String(), m, err.Error())
		}
	}
}

func TestBacktrack(t *testing.T) {
	b := getBoard("bigger-board.txt", t)

	solution, history, err := Backtrack(b)

	if err != nil {
		t.Fatalf("Backtrack failed: %s\n%v", err.Error())
	}

	t.Logf("Backtrack SUCCESS!!!! solution in %d steps, explored %d board states", len(solution), len(history))
	t.Logf("finalBoard:\n%s", solution[len(solution)-1])
}

func TestParseBigBoard(t *testing.T) {
	b := getBoard("big-board.txt", t)
	t.Logf("Big Board:\n%s", b)
}
