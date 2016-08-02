package board

import (
	"os"
	"testing"
)

func TestColorCell(t *testing.T) {
	f, err := os.Open("../board.txt")

	if err != nil {
		t.Fatalf("can't open test board file: %s", err.Error())
	}
	defer f.Close()

	b, err := New(f)
	if err != nil {
		t.Fatalf("can't parse board file: %s", err.Error())
	}

	err = b.ColorCell(0, 1, 0)
	if err != nil {
		t.Fatalf("can't color cell: %s", err.Error())
	}
	err = b.ColorCell(0, 1, 0)
	if err == nil {
		t.Fatal("can't color cell again")
	}

}

func TestSolved(t *testing.T) {
	f, err := os.Open("../board.txt")

	if err != nil {
		t.Fatalf("can't open test board file: %s", err.Error())
	}
	defer f.Close()

	b, err := New(f)
	if err != nil {
		t.Fatalf("can't parse board file: %s", err.Error())
	}

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
