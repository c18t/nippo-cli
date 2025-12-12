package tui

import (
	"fmt"
	"sync"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

// SpinnerModel wraps bubbles spinner for progress indication
type SpinnerModel struct {
	spinner     spinner.Model
	message     string
	done        bool
	interrupted bool
}

// DoneMsg signals that the spinner should stop
type DoneMsg struct{}

// UpdateMessageMsg updates the spinner message
type UpdateMessageMsg struct {
	Message string
}

// NewSpinnerModel creates a new spinner model with the given message
func NewSpinnerModel(message string) SpinnerModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = SpinnerStyle
	return SpinnerModel{
		spinner: s,
		message: message,
		done:    false,
	}
}

// Init initializes the spinner model
func (m SpinnerModel) Init() tea.Cmd {
	return m.spinner.Tick
}

// Update handles messages for the spinner
func (m SpinnerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case DoneMsg:
		m.done = true
		return m, tea.Quit
	case UpdateMessageMsg:
		m.message = msg.Message
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.done = true
			m.interrupted = true
			return m, tea.Quit
		}
		return m, nil
	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
}

// View renders the spinner
func (m SpinnerModel) View() string {
	if m.done {
		if m.interrupted {
			return fmt.Sprintf("%s %s %s", WarningStyle.Render("✗"), m.message, DimStyle.Render("(interrupted)"))
		}
		return fmt.Sprintf("%s %s", SuccessStyle.Render("✓"), m.message)
	}
	return fmt.Sprintf("%s %s", m.spinner.View(), m.message)
}

// RunSpinner runs a spinner while executing the given function
func RunSpinner(message string, fn func() error) error {
	m := NewSpinnerModel(message)

	p := tea.NewProgram(m)
	errCh := make(chan error, 1)

	go func() {
		err := fn()
		p.Send(DoneMsg{})
		errCh <- err
	}()

	if _, err := p.Run(); err != nil {
		return err
	}

	return <-errCh
}

// SpinnerController manages a spinner that can be started/stopped/updated
type SpinnerController struct {
	program *tea.Program
	running bool
	done    chan struct{}
	mu      sync.Mutex
}

// NewSpinnerController creates a new spinner controller
func NewSpinnerController() *SpinnerController {
	return &SpinnerController{}
}

// Start starts the spinner with the given message
func (c *SpinnerController) Start(message string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.running {
		return
	}

	m := NewSpinnerModel(message)
	c.program = tea.NewProgram(m)
	c.done = make(chan struct{})
	c.running = true

	go func() {
		c.program.Run()
		close(c.done)
	}()
}

// UpdateMessage updates the spinner message
func (c *SpinnerController) UpdateMessage(message string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.running || c.program == nil {
		return
	}

	c.program.Send(UpdateMessageMsg{Message: message})
}

// Stop stops the spinner
func (c *SpinnerController) Stop() {
	c.mu.Lock()
	if !c.running || c.program == nil {
		c.mu.Unlock()
		return
	}
	c.program.Send(DoneMsg{})
	done := c.done
	c.running = false
	c.mu.Unlock()

	// Wait for program to finish
	<-done
}
