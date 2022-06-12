package sudoku

import "testing"

func TestSudoku(t *testing.T) {
	t.Log("TestSudoku")
	var sudoku Sudoku
	modes := map[int]string{
		EASY:   "easy",
		MEDIUM: "medium",
		HARD:   "hard",
		EXPERT: "expert",
	}

	for d, mode := range modes {
		t.Log("Mode:", mode)
		sudoku = *New(d)
		printSudoku(&sudoku)
	}
}
