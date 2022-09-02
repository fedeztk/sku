package board

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fedeztk/sku/internal/model/components/keys"
	"github.com/fedeztk/sku/pkg/sudoku"
)

type Model struct {
	board [sudoku_len][sudoku_len]struct {
		puzzle     int
		answer     int
		modifiable bool
	}
	KeyMap         keys.KeyMap
	cursor         coordinate
	cellsLeft      int
	errCoordinates map[coordinate]interface{}

	Err error
}

type GameWon struct{}

type gameCheck struct {
	Err    error
	result map[coordinate]interface{} // maybe not the best way to do this
}

type coordinate struct {
	row, col int
}

const (
	sudoku_len = 9
)

func NewModel(mode int) Model {
	var cellsLeft int
	var board [sudoku_len][sudoku_len]struct {
		puzzle     int
		answer     int
		modifiable bool
	}

	sudoku := sudoku.New(mode)
	puzzle, answer := sudoku.Puzzle, sudoku.Answer

	for i := 0; i < sudoku_len; i++ {
		for j := 0; j < sudoku_len; j++ {
			board[i][j].puzzle = puzzle[i*sudoku_len+j]
			board[i][j].answer = answer[i*sudoku_len+j]
			if modifiable := puzzle[i*sudoku_len+j] == 0; modifiable {
				board[i][j].modifiable = modifiable
				cellsLeft++
			}
		}
	}

	return Model{
		board:          board,
		KeyMap:         keys.Keys,
		cellsLeft:      cellsLeft,
		errCoordinates: make(map[coordinate]interface{}),
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.KeyMap.Down):
			m.cursorDown()
		case key.Matches(msg, m.KeyMap.Up):
			m.cursorUp()
		case key.Matches(msg, m.KeyMap.Left):
			m.cursorLeft()
		case key.Matches(msg, m.KeyMap.Right):
			m.cursorRight()
		case key.Matches(msg, m.KeyMap.Clear):
			m.clear(m.cursor.row, m.cursor.col)
		case key.Matches(msg, m.KeyMap.Number):
			return m, m.set(m.cursor.row, m.cursor.col, int(msg.String()[0]-'0'))
		}
	case gameCheck:
		m.Err = msg.Err
		m.errCoordinates = msg.result
		if msg.Err == nil {
			return m, m.won()
		}
	}
	return m, nil
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) View() string {
	// replace 0 with empty string
	var maybeReplace = func(v int) string {
		if v == 0 {
			return " "
		}
		return fmt.Sprintf("%d", v)
	}

	var boardView string
	for i := 0; i < sudoku_len; i++ {
		row := ""
		for j := 0; j < sudoku_len; j++ {
			_, isError := m.errCoordinates[coordinate{i, j}]

			row += formatCell(isError, m.cursor.row == i && m.cursor.col == j,
				m.board[i][j].modifiable, i, j, maybeReplace(m.board[i][j].puzzle))
		}
		boardView += formatRow(i, row) + "\n"
	}

	return boardView + cellsLeftStyle.Render(fmt.Sprintf("Cells left: %d", m.cellsLeft))
}

func (m *Model) set(row, col int, v int) tea.Cmd {
	if m.board[row][col].modifiable {
		if m.board[row][col].puzzle == 0 {
			m.cellsLeft--
		}
		m.board[row][col].puzzle = v

		delete(m.errCoordinates, coordinate{row, col})

		if m.cellsLeft == 0 {
			return m.check()
		}
	}
	return nil
}

func (m *Model) clear(row, col int) {
	if m.board[row][col].modifiable {
		if m.board[row][col].puzzle != 0 {
			m.cellsLeft++
		}
		m.board[row][col].puzzle = 0

		delete(m.errCoordinates, coordinate{row, col})
	}
}

func (m *Model) cursorDown() {
	if m.cursor.row < sudoku_len-1 {
		m.cursor.row++
	} else {
		m.cursor.row = 0
	}
}

func (m *Model) cursorUp() {
	if m.cursor.row > 0 {
		m.cursor.row--
	} else {
		m.cursor.row = sudoku_len - 1
	}
}

func (m *Model) cursorLeft() {
	if m.cursor.col > 0 {
		m.cursor.col--
	} else {
		m.cursor.col = sudoku_len - 1
	}
}

func (m *Model) cursorRight() {
	if m.cursor.col < sudoku_len-1 {
		m.cursor.col++
	} else {
		m.cursor.col = 0
	}
}

func (m *Model) check() tea.Cmd {
	return func() tea.Msg {
		result := make(map[coordinate]interface{})
		for i := 0; i < sudoku_len; i++ {
			for j := 0; j < sudoku_len; j++ {
				if m.board[i][j].puzzle != m.board[i][j].answer {
					result[coordinate{i, j}] = nil
				}
			}
		}

		if len(result) == 0 {
			return gameCheck{Err: nil, result: result}
		}
		return gameCheck{Err: fmt.Errorf("%d errors", len(result)), result: result}
	}
}

func (m *Model) won() tea.Cmd {
	return func() tea.Msg {
		return GameWon{}
	}
}
