package sudoku

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	EASY   = 45
	MEDIUM = 36
	HARD   = 27
	EXPERT = 13
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
	erased := 0
	const SUDOKU_SIZE = 81
	for SUDOKU_SIZE-erased > difficulty {
		row, col := 0, 0
		for s.Puzzle[row*9+col] == 0 {
			rand.Seed(time.Now().UnixNano())
			row, col = rand.Intn(9), rand.Intn(9)
		}

		backup := s.Puzzle[row*9+col]
		s.Puzzle[row*9+col] = 0

		copyGrid := s.Puzzle

		count := 0
		solve(&copyGrid, &count)
		if count != 1 {
			s.Puzzle[row*9+col] = backup
		} else {
			erased++
		}
	}
}

func solve(grid *[81]int, count *int) bool {
	var row, col int
	for i := 0; i < 81; i++ {
		row = i / 9
		col = i % 9
		if grid[row*9+col] == 0 {
			for n := 1; n <= 9; n++ {
				if isValid(grid, row, col, n) {
					grid[row*9+col] = n
					if checkFull(grid) {
						*count++
						break
					} else {
						if solve(grid, count) {
							return true
						}
					}
				}
			}
			break
		}
	}
	grid[row*9+col] = 0
	return false
}

func fill(grid *[81]int) bool {
	numberList := [9]int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	var row, col int
	for i := 0; i < 81; i++ {
		row = i / 9
		col = i % 9
		if grid[row*9+col] == 0 {
			rand.Seed(time.Now().UnixNano())
			rand.Shuffle(len(numberList), func(i, j int) {
				numberList[i], numberList[j] = numberList[j], numberList[i]
			})
			for _, n := range numberList {
				if isValid(grid, row, col, n) {
					grid[row*9+col] = n
					if checkFull(grid) {
						return true
					} else {
						if fill(grid) {
							return true
						}
					}
				}
			}
			break
		}
	}
	grid[row*9+col] = 0
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

func isValid(grid *[81]int, row, col, n int) bool {
	for i := 0; i < 9; i++ {
		if grid[row*9+i] == n || grid[i*9+col] == n {
			return false
		}
	}
	// check numbers on the same 3x3 block
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if grid[(row/3*3+i)*9+(col/3*3+j)] == n {
				return false
			}
		}
	}
	return true
}

func printSudoku(s *Sudoku) {
	for i := 0; i < 81; i = i + 9 {
		fmt.Println(s.Puzzle[i : i+9])
	}
	fmt.Println()
	for i := 0; i < 81; i = i + 9 {
		fmt.Println(s.Answer[i : i+9])
	}
	erased := 0
	for i := 0; i < 81; i++ {
		if s.Puzzle[i] == 0 {
			fmt.Printf(".")
			erased++
		} else {
			fmt.Printf("%d", s.Puzzle[i])
		}
	}
	fmt.Println()
	fmt.Println(erased)
}
