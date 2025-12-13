package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// WizardStep represents a single step in the wizard
type WizardStep struct {
	Label        string
	Placeholder  string // Hint text shown when input is empty
	InitialValue string // Pre-filled value (from existing config)
	Value        string
	Completed    bool
}

// WizardModel manages multiple input steps in sequence
type WizardModel struct {
	steps       []WizardStep
	currentStep int
	textInput   textinput.Model
	cancelled   bool
	err         error
}

// NewWizardModel creates a new wizard with the given steps
func NewWizardModel(steps []WizardStep) WizardModel {
	ti := textinput.New()
	ti.CharLimit = 256
	ti.Width = 60
	ti.PromptStyle = PromptStyle

	m := WizardModel{
		steps:       steps,
		currentStep: 0,
		textInput:   ti,
	}
	m.initCurrentStep()
	return m
}

func (m *WizardModel) initCurrentStep() {
	if m.currentStep < len(m.steps) {
		step := m.steps[m.currentStep]
		m.textInput.Placeholder = step.Placeholder
		m.textInput.SetValue(step.InitialValue)
		m.textInput.Focus()
	}
}

// Init initializes the wizard model
func (m WizardModel) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles messages for the wizard
func (m WizardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			// Save current value and move to next step
			// Use placeholder as fallback if input is empty
			value := m.textInput.Value()
			if value == "" {
				value = m.steps[m.currentStep].Placeholder
			}
			m.steps[m.currentStep].Value = value
			m.steps[m.currentStep].Completed = true
			m.currentStep++

			if m.currentStep >= len(m.steps) {
				// All steps completed
				return m, tea.Quit
			}

			// Initialize next step
			m.initCurrentStep()
			return m, nil

		case tea.KeyCtrlC, tea.KeyEsc:
			m.cancelled = true
			m.err = fmt.Errorf("input cancelled")
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

// View renders the wizard
func (m WizardModel) View() string {
	var b strings.Builder

	// Show completed steps
	for i := 0; i < len(m.steps); i++ {
		step := m.steps[i]
		if step.Completed {
			b.WriteString(PromptStyle.Render(step.Label))
			b.WriteString("\n")
			b.WriteString(SuccessStyle.Render(fmt.Sprintf("> %s ✓", step.Value)))
			b.WriteString("\n\n")
		} else if i == m.currentStep && !m.cancelled {
			// Show current step (if not cancelled)
			b.WriteString(PromptStyle.Render(step.Label))
			b.WriteString("\n")
			b.WriteString(m.textInput.View())
			b.WriteString("\n")
		} else if m.cancelled && i == m.currentStep {
			// Show cancelled step
			b.WriteString(PromptStyle.Render(step.Label))
			b.WriteString("\n")
			b.WriteString(WarningStyle.Render("> ") + DimStyle.Render("(cancelled)"))
			b.WriteString("\n")
		}
	}

	return b.String()
}

// Values returns all entered values in order
func (m WizardModel) Values() []string {
	values := make([]string, len(m.steps))
	for i, step := range m.steps {
		values[i] = step.Value
	}
	return values
}

// Cancelled returns true if the wizard was cancelled
func (m WizardModel) Cancelled() bool {
	return m.cancelled
}

// Err returns any error that occurred
func (m WizardModel) Err() error {
	return m.err
}

// RunWizard runs an interactive wizard and returns all values
func RunWizard(steps []WizardStep) ([]string, error) {
	m := NewWizardModel(steps)
	p := tea.NewProgram(m)

	finalModel, err := p.Run()
	if err != nil {
		return nil, err
	}

	result := finalModel.(WizardModel)
	if result.Err() != nil {
		return nil, result.Err()
	}

	return result.Values(), nil
}

// ConfirmModel is a simple yes/no confirmation prompt
type ConfirmModel struct {
	prompt    string
	confirmed bool
	answered  bool
	cancelled bool
	err       error
}

// NewConfirmModel creates a new confirmation prompt
func NewConfirmModel(prompt string, defaultYes bool) ConfirmModel {
	return ConfirmModel{
		prompt:    prompt,
		confirmed: defaultYes,
	}
}

// Init initializes the confirm model
func (m ConfirmModel) Init() tea.Cmd {
	return nil
}

// Update handles messages for the confirm model
func (m ConfirmModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "y", "Y":
			m.confirmed = true
			m.answered = true
			return m, tea.Quit
		case "n", "N", "enter":
			m.confirmed = false
			m.answered = true
			return m, tea.Quit
		case "ctrl+c", "esc":
			m.cancelled = true
			m.err = fmt.Errorf("input cancelled")
			return m, tea.Quit
		}
	}
	return m, nil
}

// View renders the confirmation prompt
func (m ConfirmModel) View() string {
	if m.answered {
		answer := "No"
		if m.confirmed {
			answer = "Yes"
		}
		return PromptStyle.Render(m.prompt) + " " + DimStyle.Render("[y/N]") + " " + SuccessStyle.Render(answer+"\n")
	}
	if m.cancelled {
		return PromptStyle.Render(m.prompt) + " " + DimStyle.Render("[y/N]") + " " + WarningStyle.Render("(cancelled)\n")
	}
	return PromptStyle.Render(m.prompt) + " " + DimStyle.Render("[y/N]") + " "
}

// RunConfirm runs an interactive yes/no confirmation and returns the result
func RunConfirm(prompt string, defaultYes bool) (bool, error) {
	m := NewConfirmModel(prompt, defaultYes)
	p := tea.NewProgram(m)

	finalModel, err := p.Run()
	if err != nil {
		return false, err
	}

	result := finalModel.(ConfirmModel)
	if result.err != nil {
		return false, result.err
	}

	return result.confirmed, nil
}
