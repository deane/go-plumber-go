package main

import (
	"flag"
	"log"
	"os"

	"github.com/deane/go-plumber-go/board"
)

var fileName = flag.String("file", "board.txt", "go-plumber-go board file to solve")
var display = flag.Bool("show-results", false, "print the solution if the solving is successful")
var canonical = flag.Bool("canonical", false, "only search for the canonical solution")
var detectDead = flag.Bool("detect-dead", false, "detect dead cells (surrounded by different colors on all sides)")
var sort = flag.Bool("sort", false, "start with the longest flows")
var sortReverse = flag.Bool("sort-reverse", false, "start with the shortest flows")

func main() {
	flag.Parse()
	f, err := os.Open(*fileName)
	if err != nil {
		log.Printf("can't open the file %s, ERROR: %s\n", *fileName, err)
		return
	}
	board.ShowResults = *display
	board.Canonical = *canonical
	board.DetectDead = *detectDead
	defer f.Close()
	b, err := board.New(f)
	if err != nil {
		log.Println(err)
	}
	if *sort {
		b = board.SortColors(b, false)
	} else if *sortReverse {
		b = board.SortColors(b, true)
    }
	log.Println(b)
	_, err = board.Backtrack(b)

	if err != nil {
		log.Fatalf("Couldn't resolve the puzzle: %s", err.Error())
	}

}
