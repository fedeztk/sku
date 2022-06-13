package sudoku

import (
	"fmt"
	"testing"
)

func TestSudoku(t *testing.T) {
	t.Log("TestSudoku")
	var sudoku Sudoku
	modes := map[int]string{
		EASY:   "easy",
		MEDIUM: "medium",
		HARD:   "hard",
		// EXPERT: "expert",
	}

	for d, mode := range modes {
		t.Log("Mode:", mode)
		sudoku = *New(d)
		t.Log(printSudoku(&sudoku))
	}
}

func printSudoku(s *Sudoku) string {
	board := "\n"
	for i := 0; i < 81; i = i + 9 {
		board += fmt.Sprintf("%d\n", s.Puzzle[i:i+9])
	}
	board += "\n"
	for i := 0; i < 81; i = i + 9 {
		board += fmt.Sprintf("%d\n", s.Answer[i:i+9])
	}
	notErased := 0
	for i := 0; i < 81; i++ {
		if s.Puzzle[i] == 0 {
			board += "."
		} else {
			board += fmt.Sprintf("%d", s.Puzzle[i])
			notErased++
		}
	}
	board += fmt.Sprintf("\nNot erased: %d", notErased)
	return board
}
