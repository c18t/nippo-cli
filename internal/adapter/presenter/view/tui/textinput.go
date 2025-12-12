package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// TextInputModel wraps bubbles textinput for user input
type TextInputModel struct {
	textInput textinput.Model
	label     string
	submitted bool
	err       error
}

// NewTextInputModel creates a new text input model
func NewTextInputModel(label string, defaultValue string) TextInputModel {
	ti := textinput.New()
	ti.Placeholder = defaultValue
	ti.SetValue(defaultValue)
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = 60
	ti.PromptStyle = PromptStyle

	return TextInputModel{
		textInput: ti,
		label:     label,
		submitted: false,
	}
}

// Init initializes the text input model
func (m TextInputModel) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles messages for the text input
func (m TextInputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			m.submitted = true
			return m, tea.Quit
		case tea.KeyCtrlC, tea.KeyEsc:
			m.err = fmt.Errorf("input cancelled")
			return m, tea.Quit
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

// View renders the text input
func (m TextInputModel) View() string {
	if m.submitted {
		return fmt.Sprintf("%s\n%s %s\n", PromptStyle.Render(m.label), SuccessStyle.Render(">"), m.textInput.Value())
	}
	if m.err != nil {
		return fmt.Sprintf("%s\n%s %s\n", PromptStyle.Render(m.label), WarningStyle.Render(">"), DimStyle.Render("(cancelled)"))
	}
	return fmt.Sprintf("%s\n%s\n", PromptStyle.Render(m.label), m.textInput.View())
}

// Value returns the current input value
func (m TextInputModel) Value() string {
	return m.textInput.Value()
}

// Err returns any error that occurred
func (m TextInputModel) Err() error {
	return m.err
}

// RunTextInput runs an interactive text input and returns the result
func RunTextInput(label string, defaultValue string) (string, error) {
	m := NewTextInputModel(label, defaultValue)
	p := tea.NewProgram(m)

	finalModel, err := p.Run()
	if err != nil {
		return "", err
	}

	result := finalModel.(TextInputModel)
	if result.Err() != nil {
		return "", result.Err()
	}

	return result.Value(), nil
}
