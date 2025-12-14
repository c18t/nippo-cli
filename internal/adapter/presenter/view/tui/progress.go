package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

// ProgressModel wraps bubbles progress bar for trackable operations
type ProgressModel struct {
	progress    progress.Model
	message     string
	percent     float64
	done        bool
	interrupted bool
}

// ProgressMsg updates the progress percentage
type ProgressMsg struct {
	Percent float64
}

// ProgressDoneMsg signals that the progress is complete
type ProgressDoneMsg struct{}

// NewProgressModel creates a new progress bar model
func NewProgressModel(message string) ProgressModel {
	p := progress.New(progress.WithSolidFill("#6B8E23"))
	return ProgressModel{
		progress: p,
		message:  message,
		percent:  0,
		done:     false,
	}
}

// Init initializes the progress model
func (m ProgressModel) Init() tea.Cmd {
	return nil
}

// Update handles messages for the progress bar
func (m ProgressModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case ProgressMsg:
		m.percent = msg.Percent
		return m, nil
	case ProgressDoneMsg:
		m.done = true
		m.percent = 1.0
		return m, tea.Quit
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.done = true
			m.interrupted = true
			return m, tea.Quit
		}
		return m, nil
	case tea.WindowSizeMsg:
		m.progress.Width = msg.Width - 4
		if m.progress.Width > 80 {
			m.progress.Width = 80
		}
		return m, nil
	default:
		var cmd tea.Cmd
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd
	}
}

// View renders the progress bar
func (m ProgressModel) View() string {
	suffix := ""
	if m.done && m.interrupted {
		suffix = " " + DimStyle.Render("(interrupted)")
	}
	return fmt.Sprintf("%s\n%s%s", m.message, m.progress.ViewAs(m.percent), suffix)
}
