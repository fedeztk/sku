# Sudoku Engine
This package is a **non** efficient implementation of a sudoku solver and generator. The core logic is ~~stolen~~ inspired by [this python implementation](https://www.101computing.net/sudoku-generator-algorithm/), all credits go to the author. It is a simple recursive and backtracking engine, easy to understand and read.

The API should be considered stable, however there is some work in progress about the internal implementation due to performance experimentation, see [here](#performance-analysis).

## Motivation
This minimal implementation was done due to the lack of a package that simply has the two sensible features needed by a sudoku engine: 
1. can both solve & generate 
2. do generate a one-solution puzzle  

without useless features or generally cumbersome development experience (e.g. overly complicated structures, interfaces, initialization methods, etc).

If you need the aforementioned features this package is not the right one, since it only does point 1 and 2.

If you need a super efficient implementation this package is not the right one since it is aimed to be simple at the expense of speed; if you need something like that you can use/write one that leverages a more performant algorithm such as dancing links (dlx).

## Usage
```go
sudoku := sudoku.New(sudoku.EASY) // also available: MEDIUM, HARD, EXPERT

fmt.Println(sudoku.Puzzle) // unsolved sudoku [81]int
fmt.Println(sudoku.Answer) // solved sudoku   [81]int
```

## Performance analysis

> This section is not relevant to the package or its usage, but rather it is an insight on the performance of various strategies for speeding up the sudoku-generation problem.
Since I'm new to Go I wanted to experiment a bit more on concurrency, channel, context, etc., and I thought that leaving here some notes would be cool for other newcomers like me.

> For the experienced Go programmers reading this: I'm sorry if you find awful code and bad practices, I'm open to PRs and ideas on more idiomatic ways to solve the problem.

> For the non-experienced Go programmers reading this: please note that the code below is written by a novice, take everything with a grain of salt.
Basic knwoledge of channel, goroutines and context is required.

#### The Problem

Obtain a sudoku that has one **and only one** solution in a reasonable amount of time (e.g. ~3s max for harder puzzles).

### First approach: no concurrency

Having `MAX_ATTEMPTS` and `MAX_EXEC_TIME` defined, we try `#MAX_ATTEMPTS`
```go
func New(difficulty int) *Sudoku {
	for i := 0; i < MAX_ATTEMPTS; i++ {
		if s, ok := newWithTimer(difficulty); ok {
			return s
		}
	}
	return nil
}
```
...to find a solution in `#MAX_EXEC_TIME`' 
```go
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
```

This solution is suitable for almost all puzzles, but seems to be just sufficient for complex ones (e.g. expert puzzles take 2-3s).

### Second approach: concurrency with context.Context and channels

This approach takes advantage of concurrency by starting one process on each CPUs and waiting for the first one to finish.
The quickest goroutine that find a valid sudoku, sends it through the channel `sudokuCh`, where the `New` function was hanging. After that, the context on all suboutines is cancelled with `cancelFunc`.

```go
func NewParallel(difficulty int) *Sudoku {
	CPUs := runtime.NumCPU()
	runtime.GOMAXPROCS(CPUs)

	ctx, cancelFunc := context.WithCancel(context.Background())

	sudokuCh := make(chan *Sudoku)
	for c := 0; c < CPUs; c++ {
		go newWithTimerParallel(difficulty, ctx, sudokuCh)
	}

	sudoku := <-sudokuCh
	cancelFunc()
	return sudoku
}
```

Here it's the same as before except for the periodic check of the context.

```go
func newWithTimerParallel(difficulty int, ctx context.Context, sudokuCh chan *Sudoku) {
	s := &Sudoku{}

	done := make(chan struct{})

	go func() {
		s = &Sudoku{}
		fill(&s.Puzzle)
        // eraseSome can be skipped
		select {
		case <-ctx.Done():
			return
		default:
		}
		s.Answer = s.Puzzle
		s.eraseSome(difficulty)
		close(done)
	}()

	select {
	// channel is closed, a valid sudoku was generated
	case <-done:
		sudokuCh <- s

	// timeout expired
	case <-time.After(time.Second * MAX_EXEC_TIME):

	// there is another process that found a valid sudoku quicker
	case <-ctx.Done():
	}
}
```

### Third approach: concurrency with just channels

This is a simplification of [the previous solution](#second-approach%3A-concurrency-with-context.context-and-channels) that does not use contexts. This cause the goroutines spawned to run for more time than what is really needed for the sake of simplicity (the whole point of this package after all).

```go
func NewParallelNoCtx(difficulty int) *Sudoku {
	CPUs := runtime.NumCPU()
	runtime.GOMAXPROCS(CPUs)

	sudokuCh := make(chan *Sudoku)
	for c := 0; c < CPUs; c++ {
		go newWithTimerParallelNoCtx(difficulty, sudokuCh)
	}

	sudoku := <-sudokuCh
	return sudoku
}
```

```go
func newWithTimerParallelNoCtx(difficulty int, sudokuCh chan *Sudoku) {
	s := &Sudoku{}

	done := make(chan struct{})

	go func() {
		s = &Sudoku{}
		fill(&s.Puzzle)
		s.Answer = s.Puzzle
		s.eraseSome(difficulty)
		close(done)
	}()

	select {
	// channel is closed, a valid sudoku was generated
	case <-done:
		sudokuCh <- s

	// timeout expired
	case <-time.After(time.Second * MAX_EXEC_TIME):
	}
}
```

### Results

I tested the three solutions with a simple benchmark:
```go
# first approach
func BenchmarkSudoku(b *testing.B) {
	for i := 0; i < b.N; i++ {
		New(EXPERT)
	}
}

# second approach
func BenchmarkSudokuParallel(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewParallel(EXPERT)
	}
}

# third approach
func BenchmarkSudokuParallelNoCtx(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewParallelNoCtx(EXPERT)
	}
}
```

That's the result of the benchmark:

```
$ go test -bench=. -count=2
goos: linux
goarch: amd64
pkg: github.com/fedeztk/sku/pkg/sudoku
cpu: AMD Ryzen 7 PRO 4750U with Radeon Graphics
# first approach
BenchmarkSudoku-16                 	      55	 619782643 ns/op
BenchmarkSudoku-16                 	       6	 868804023 ns/op
# second approach
BenchmarkSudokuParallel-16         	      13	 346577499 ns/op
BenchmarkSudokuParallel-16         	       3	 446622813 ns/op
# third approach
BenchmarkSudokuParallelNoCtx-16    	       8	 936224567 ns/op
BenchmarkSudokuParallelNoCtx-16    	       1	2072859339 ns/op
PASS
ok  	github.com/fedeztk/sku/pkg/sudoku	62.719s
```

#### Things learned

The third option is inadequate performance wise. When doing concurrent programming, a reasonable amount of "complexity" is required, otherwise it is better to stick with simple sequential processes.

Since I'm not convinced on the second approach, despite it being faster, I will stick with the easier to read and comprehend: the first one. To me the performance gained with the second one over the first is still too low (going from 1 CPU to 16 CPUs to find a proper solution, I expected an improvement of roughly one order of magnitude). It probably is my fault since, as stated before, I'm new to Go. I will update here once I found out more.
