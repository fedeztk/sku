package ui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fedeztk/sku/internal/model/components/animate"
	"github.com/fedeztk/sku/internal/model/components/board"
	"github.com/fedeztk/sku/internal/model/components/help"
	"github.com/fedeztk/sku/internal/model/components/keys"
	"github.com/fedeztk/sku/internal/model/components/stopwatch"
)

type Model struct {
	help      help.Model
	stopwatch stopwatch.Model
	board     board.Model
	animate   animate.Model

	width, height int
	gameEnded     bool

	err error
}

func (m Model) Init() tea.Cmd {
	return m.stopwatch.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd, helpCmd, boardCmd, stopwatchCmd, animCmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Keys.Quit):
			cmd = tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case board.GameWon:
		m.gameEnded = true
		cmd = m.animate.Init()
	}

	m.help, helpCmd = m.help.Update(msg)
	m.board, boardCmd = m.board.Update(msg)
	m.stopwatch, stopwatchCmd = m.stopwatch.Update(msg)
	m.animate, animCmd = m.animate.Update(msg)

	cmds = append(cmds, cmd, helpCmd, boardCmd, stopwatchCmd, animCmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.err != nil {
		return m.err.Error()
	}

	if m.gameEnded {
		return lipgloss.Place(m.width-int(m.animate.TargetX),
			m.height, lipgloss.Center, lipgloss.Center, m.animate.View())
	}

	// spaces on board and help view are needed to make the view centered
	mainView := lipgloss.JoinVertical(lipgloss.Center,
		m.board.View()+" ", m.stopwatch.View(), " "+m.help.View())

	mainView = lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, mainView)

	return mainView
}

func NewModel(mode int) Model {
	return Model{
		help:      help.NewModel(),
		stopwatch: stopwatch.NewModel(),
		board:     board.NewModel(mode),
		animate:   animate.NewModel(),
	}
}
