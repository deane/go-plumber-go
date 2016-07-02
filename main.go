package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/deane/go-plumber-go/board"
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
	board, err := board.New(f)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(board)
}
