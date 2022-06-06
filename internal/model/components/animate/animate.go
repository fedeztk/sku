package animate

import (
	"fmt"
	"math"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/harmonica"
)

const (
	fps          = 60
	spriteWidth  = 12
	spriteHeight = 5
	frequency    = 5.0
	damping      = 0.15
)

type frameMsg time.Time

type Model struct {
	x       float64
	xVel    float64
	TargetX float64
	spring  harmonica.Spring
}

func animate() tea.Cmd {
	return tea.Tick(time.Second/fps, func(t time.Time) tea.Msg {
		return frameMsg(t)
	})
}

func wait(d time.Duration) tea.Cmd {
	return func() tea.Msg {
		time.Sleep(d)
		return nil
	}
}

func NewModel() Model {
	return Model{
		spring:  harmonica.NewSpring(harmonica.FPS(fps), frequency, damping),
		TargetX: 20,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Sequentially(wait(time.Second/2), animate())
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg.(type) {

	// Step foreward one frame
	case frameMsg:
		// Update x position (and velocity) with our spring.
		m.x, m.xVel = m.spring.Update(m.x, m.xVel, m.TargetX)

		// Quit when we're basically at the target position.
		if math.Abs(m.x-m.TargetX) < 0.01 {
			return m, tea.Sequentially(wait(3/4*time.Second), tea.Quit)
		}

		// Request next frame
		return m, animate()

	default:
		return m, nil
	}
}

func (m Model) View() string {
	var out strings.Builder
	fmt.Fprint(&out, "\n")

	x := int(math.Round(m.x))
	if x < 0 {
		return ""
	}

	spriteRow := spriteStyle.Render(strings.Repeat(" ", spriteWidth) +
		getWinningString(int(m.x)) +
		strings.Repeat(" ", spriteWidth))

	row := strings.Repeat(" ", x) + spriteRow + "\n"

	fmt.Fprint(&out, strings.Repeat(row, spriteHeight))

	return out.String()
}

func getWinningString(timeFrame int) string {
	won := "YOU WON!"
	res := ""

	for i, c := range won {
		if i < timeFrame {
			res += string(c)
		} else {
			res += " "
		}
	}

	return res
}
