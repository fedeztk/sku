package main

import (
	_ "embed"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	model "github.com/fedeztk/sku/internal/model"
)

const (
	LEVEL_EASY = iota
	LEVEL_MEDIUM
	LEVEL_HARD
	LEVEL_EXPERT
)

//go:generate ./get_version.sh
//go:embed .version
var skuVersion string

func main() {
	var mode int
	modesMap := map[string]int{
		"easy":   LEVEL_EASY,
		"medium": LEVEL_MEDIUM,
		"hard":   LEVEL_HARD,
		"expert": LEVEL_EXPERT,
	}

	if len(os.Args) < 2 {
		mode = LEVEL_EASY
	} else {
		if os.Args[1] == "-v" || os.Args[1] == "--version" {
			fmt.Println(skuVersion)
			os.Exit(0)
		}

		if os.Args[1] == "-h" || os.Args[1] == "--help" {
			fmt.Println(getHelp())
			os.Exit(0)
		}

		if _, ok := modesMap[os.Args[1]]; !ok {
			fmt.Printf("Invalid mode: %s\n", os.Args[1])
			fmt.Println(getHelp())
			os.Exit(1)
		}

		mode = modesMap[os.Args[1]]
	}

	p := tea.NewProgram(model.NewModel(mode), tea.WithAltScreen(), tea.WithMouseCellMotion())
	if err := p.Start(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func getHelp() string {
	return `sku - a simple sudoku game
    -v, --version    show version
    -h, --help       show this help
    [mode]           easy, medium, hard, expert`
}
