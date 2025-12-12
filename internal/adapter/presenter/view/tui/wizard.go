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
	DefaultValue string
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
		m.textInput.Placeholder = step.DefaultValue
		m.textInput.SetValue(step.DefaultValue)
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
			m.steps[m.currentStep].Value = m.textInput.Value()
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
