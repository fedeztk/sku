# Sudoku Engine
This package is a **non** efficient implementation of a sudoku solver and generator. It is ~~stolen~~ ported from [this python implementation](https://www.101computing.net/sudoku-generator-algorithm/), all credits go to the author. It is a simple recursive and backtracking engine, easy to understand and read.

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