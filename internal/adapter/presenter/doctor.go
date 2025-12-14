package presenter

import (
	"fmt"

	"github.com/c18t/nippo-cli/internal/usecase/port"
	"github.com/charmbracelet/lipgloss"
	"github.com/samber/do/v2"
)

type DoctorPresenter interface {
	Show(output *port.DoctorUseCaseOutputData)
}

type doctorPresenter struct{}

func NewDoctorPresenter(_ do.Injector) (DoctorPresenter, error) {
	return &doctorPresenter{}, nil
}

func (p *doctorPresenter) Show(output *port.DoctorUseCaseOutputData) {
	// Styles
	successStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
	warnStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("3"))
	categoryStyle := lipgloss.NewStyle().Bold(true).Underline(true)
	suggestionStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))

	// Group by category
	categories := make(map[string][]port.DoctorCheck)
	categoryOrder := []string{}
	for _, check := range output.Checks {
		if _, exists := categories[check.Category]; !exists {
			categoryOrder = append(categoryOrder, check.Category)
		}
		categories[check.Category] = append(categories[check.Category], check)
	}

	// Collect issues for summary
	var issues []port.DoctorCheck

	// Print results
	fmt.Println()
	for _, category := range categoryOrder {
		checks := categories[category]
		fmt.Println(categoryStyle.Render(category))

		for _, check := range checks {
			var statusIcon string
			switch check.Status {
			case port.DoctorCheckStatusPass:
				statusIcon = successStyle.Render("✓")
			case port.DoctorCheckStatusFail:
				statusIcon = errorStyle.Render("✗")
				issues = append(issues, check)
			case port.DoctorCheckStatusWarn:
				statusIcon = warnStyle.Render("!")
				issues = append(issues, check)
			}

			fmt.Printf("  %s %s: %s\n", statusIcon, check.Item, check.Message)
		}
		fmt.Println()
	}

	// Print suggestions for issues
	if len(issues) > 0 {
		fmt.Println(categoryStyle.Render("Suggestions"))
		for _, issue := range issues {
			if issue.Suggestion != "" {
				fmt.Printf("  • %s: %s\n", issue.Item, suggestionStyle.Render(issue.Suggestion))
			}
		}
		fmt.Println()
	}
}
