package sudoku

import (
	"math/rand"
	"time"
)

const (
	// available sudoku difficulties
	EASY   = 42
	MEDIUM = 36
	HARD   = 27
	EXPERT = 25
	// time constraints
	MAX_ATTEMPTS  = 5
	MAX_EXEC_TIME = 3
	// sudoku board size
	SUDOKU_LENGTH = 9
	SUDOKU_SIZE   = SUDOKU_LENGTH * SUDOKU_LENGTH
)

type Sudoku struct {
	Puzzle [SUDOKU_SIZE]int
	Answer [SUDOKU_SIZE]int
}

func New(difficulty int) *Sudoku {
	for i := 0; i < MAX_ATTEMPTS; i++ {
		if s, ok := newWithTimer(difficulty); ok {
			return s
		}
	}
	return nil
}

func newWithTimer(difficulty int) (*Sudoku, bool) {
	s := &Sudoku{}

	done := make(chan struct{})
	var ok bool

	go func() {
		s = &Sudoku{}
		fill(&s.Puzzle)
		s.Answer = s.Puzzle
		s.eraseSome(difficulty)
		close(done)
	}()

	select {
	case <-done:
		ok = true
	case <-time.After(time.Second * MAX_EXEC_TIME):
	}

	return s, ok
}

func (s *Sudoku) eraseSome(difficulty int) {
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

func solve(grid *[SUDOKU_SIZE]int, count *int) {
	if *count > 1 { // no need to go further
		return
	}

	for idx := 0; idx < SUDOKU_SIZE; idx++ {
		if grid[idx] == 0 {
			for n := 1; n <= SUDOKU_LENGTH; n++ {
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

func fill(grid *[SUDOKU_SIZE]int) bool {
	numberList := [SUDOKU_LENGTH]int{1, 2, 3, 4, 5, 6, 7, 8, 9}

	for idx := 0; idx < SUDOKU_SIZE; idx++ {
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

func checkFull(grid *[SUDOKU_SIZE]int) bool {
	for i := 0; i < SUDOKU_SIZE; i++ {
		if grid[i] == 0 {
			return false
		}
	}
	return true
}

func isValid(grid *[SUDOKU_SIZE]int, idx, n int) bool {
	// check if num is valid in row, col, and 3x3 box
	row := idx / SUDOKU_LENGTH
	col := idx % SUDOKU_LENGTH
	for i := 0; i < SUDOKU_LENGTH; i++ {
		if grid[row*SUDOKU_LENGTH+i] == n ||
			grid[i*SUDOKU_LENGTH+col] == n ||
			grid[((row/3)*3+i/3)*SUDOKU_LENGTH+((col/3)*3+i%3)] == n {
			return false
		}
	}
	return true
}
