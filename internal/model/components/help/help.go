package help

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fedeztk/sku/internal/model/components/keys"
)

type Model struct {
	help help.Model
	keys keys.KeyMap
}

func NewModel() Model {
	return Model{
		help: help.NewModel(),
		keys: keys.Keys,
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		}

	case tea.WindowSizeMsg:
		m.help.Width = msg.Width
	}

	return m, nil
}

func (m Model) View() string {
	return helpStyle.Render(m.help.View(m.keys))
}
