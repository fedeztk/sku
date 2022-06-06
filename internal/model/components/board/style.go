package board

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	cellStyle = func(modifiable bool) lipgloss.Style {
		if modifiable {
			return lipgloss.NewStyle().PaddingLeft(1).Background(lipgloss.Color("7")).
				PaddingRight(1).Foreground(lipgloss.Color("0"))
		} else {
			return lipgloss.NewStyle().PaddingLeft(1).Background(lipgloss.Color("8")).
				PaddingRight(1)
		}
	}

	cursorCellStyle = func(modifiable bool) lipgloss.Style {
		if modifiable {
			return lipgloss.NewStyle().PaddingLeft(1).Background(lipgloss.Color("10")).
				PaddingRight(1)
		} else {
			return lipgloss.NewStyle().PaddingLeft(1).Background(lipgloss.Color("2")).
				PaddingRight(1)
		}
	}

	errorCellStyle = func(isCursor bool) lipgloss.Style {
		if isCursor {
			return lipgloss.NewStyle().PaddingLeft(1).Background(lipgloss.Color("9")).
				PaddingRight(1)
		} else {
			return lipgloss.NewStyle().PaddingLeft(1).Background(lipgloss.Color("1")).
				PaddingRight(1)
		}
	}

	formatCell = func(isError, isCursor, modifiable bool, row, col int, c string) string {
		var s lipgloss.Style
		if isError {
			s = errorCellStyle(isCursor)
		} else if isCursor {
			s = cursorCellStyle(modifiable)
		} else {
			s = cellStyle(modifiable)
		}

		// ugly hack to get the border at the center with 1 char margin
		if col+1 == 3 || col+1 == 6 {
			return s.Render(c) + lipgloss.NewStyle().
				Border(lipgloss.NormalBorder(), false, true, false, false).Margin(0, 1).Render("")
		}

		return s.Render(c)
	}

	formatRow = func(row int, r string) string {
		// ugly hack to get the border at the center with 1 char margin
		if row+1 == 3 || row+1 == 6 {
			rSize, _ := lipgloss.Size(r)
			// there are 3 sudoku boxes hence /3, the -1 is to leave space for the crossing border
			border := strings.Repeat("─", (rSize/3)-1)

			// the extra "─" is required to center the middle box
			return r + "\n" + border + "┼" + "─" + border + "┼" + border
		}
		return r
	}

	cellsLeftStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Margin(1, 0, 0, 0)
)
