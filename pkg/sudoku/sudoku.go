package sudoku

import (
	"math/rand"
	"time"
)

const (
	EASY   = 45
	MEDIUM = 36
	HARD   = 27
	// EXPERT = 13
)

type Sudoku struct {
	Puzzle [81]int
	Answer [81]int
}

func New(difficulty int) *Sudoku {
	s := &Sudoku{}
	fill(&s.Puzzle)
	s.Answer = s.Puzzle
	s.eraseSome(difficulty)

	return s
}

func (s *Sudoku) eraseSome(difficulty int) {
	const SUDOKU_SIZE = 81
	for erased := 0; SUDOKU_SIZE-erased > difficulty; {
		idx := 0
		for s.Puzzle[idx] == 0 {
			rand.Seed(time.Now().UnixNano())
			idx = rand.Intn(SUDOKU_SIZE)
		}

		copyGrid := s.Puzzle
		copyGrid[idx] = 0

		count := 0
		solve(&copyGrid, &count)
		if count == 1 {
			s.Puzzle[idx] = 0
			erased++
		}
	}
}

func solve(grid *[81]int, count *int) {
	if *count > 1 { // no need to go further
		return
	}

	var idx int
	for idx = 0; idx < 81; idx++ {
		if grid[idx] == 0 {
			for n := 1; n <= 9; n++ {
				if isValid(grid, idx, n) {
					grid[idx] = n
					if checkFull(grid) {
						*count++
					}
					solve(grid, count)
					grid[idx] = 0
				}
			}
			break
		}
	}
}

func fill(grid *[81]int) bool {
	numberList := [9]int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	var idx int
	for idx = 0; idx < 81; idx++ {
		if grid[idx] == 0 {
			rand.Seed(time.Now().UnixNano())
			rand.Shuffle(len(numberList), func(i, j int) {
				numberList[i], numberList[j] = numberList[j], numberList[i]
			})

			for _, n := range numberList {
				if isValid(grid, idx, n) {
					grid[idx] = n
					if checkFull(grid) || fill(grid) {
						return true
					}
					grid[idx] = 0
				}
			}
			break
		}
	}
	return false
}

func checkFull(grid *[81]int) bool {
	for i := 0; i < 81; i++ {
		if grid[i] == 0 {
			return false
		}
	}
	return true
}

func isValid(grid *[81]int, idx, n int) bool {
	// check if num is valid in row, col, and 3x3 box
	row := idx / 9
	col := idx % 9
	for i := 0; i < 9; i++ {
		if grid[row*9+i] == n ||
			grid[i*9+col] == n ||
			grid[((row/3)*3+i/3)*9+((col/3)*3+i%3)] == n {
			return false
		}
	}
	return true
}
