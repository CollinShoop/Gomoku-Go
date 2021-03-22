package main

import (
	"github.com/collinshoop/gomoku/src/board"
)

func main() {
	b := board.NewBoard(15, 15)
	b.Print()
	b.Move(1, 1, 1)
	b.Move(2, 1, 2)
	b.Print()
	b.Move(1, 1, 3)
	b.Print()
	b.Move(2, 0, 2)
	b.Move(2, 0, 0)
	b.Print()
	b.Move(2, 2, 2)
	b.Print()
	b.IsOver()
	b.Move(2, 4, 2)
	b.Print()
	b.IsOver()
	b.Move(2, 3, 2)
	b.Print()
	b.IsOver()
}
