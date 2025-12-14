package tui

import "fmt"

// Print outputs the given value to stdout
func Print(v interface{}) {
	fmt.Printf("%v", v)
}

// Println outputs the given value to stdout with a newline
func Println(v interface{}) {
	fmt.Printf("%v\n", v)
}

// PrintSuccess outputs a success message with styling
func PrintSuccess(message string) {
	fmt.Println(SuccessStyle.Render(message))
}

// PrintError outputs an error message with styling
func PrintError(message string) {
	fmt.Println(ErrorStyle.Render(message))
}

// PrintWarning outputs a warning message with styling
func PrintWarning(message string) {
	fmt.Println(WarningStyle.Render(message))
}

// PrintInfo outputs an informational message with styling
func PrintInfo(message string) {
	fmt.Println(InfoStyle.Render(message))
}
