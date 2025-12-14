package cmd

import (
	"testing"
)

// Tests for commandsWithoutConfig map

func TestCommandsWithoutConfig(t *testing.T) {
	tests := []struct {
		command string
		allowed bool
	}{
		{"init", true},
		{"doctor", true},
		{"help", true},
		{"version", true},
		{"nippo", true},
		{"build", false},
		{"deploy", false},
		{"clean", false},
		{"update", false},
		{"format", false},
	}

	for _, tt := range tests {
		t.Run(tt.command, func(t *testing.T) {
			got := commandsWithoutConfig[tt.command]
			if got != tt.allowed {
				t.Errorf("commandsWithoutConfig[%q] = %v, want %v", tt.command, got, tt.allowed)
			}
		})
	}
}

// Tests for rootCmd configuration

func TestRootCmdUse(t *testing.T) {
	if rootCmd.Use != "nippo" {
		t.Errorf("rootCmd.Use = %q, want %q", rootCmd.Use, "nippo")
	}
}

func TestRootCmdShort(t *testing.T) {
	expected := "nippo - The tool to power my nippo."
	if rootCmd.Short != expected {
		t.Errorf("rootCmd.Short = %q, want %q", rootCmd.Short, expected)
	}
}

func TestRootCmdHasVersionFlag(t *testing.T) {
	flag := rootCmd.Flags().Lookup("version")
	if flag == nil {
		t.Error("rootCmd should have version flag")
	}
	if flag != nil && flag.Shorthand != "V" {
		t.Errorf("version flag shorthand = %q, want %q", flag.Shorthand, "V")
	}
}

func TestRootCmdHasConfigFlag(t *testing.T) {
	flag := rootCmd.PersistentFlags().Lookup("config")
	if flag == nil {
		t.Error("rootCmd should have config persistent flag")
	}
}

func TestRootCmdHasLicenseNoticeFlag(t *testing.T) {
	flag := rootCmd.Flags().Lookup("license-notice")
	if flag == nil {
		t.Error("rootCmd should have license-notice flag")
	}
}

// Tests for subcommand configurations

func TestAuthCmdUse(t *testing.T) {
	if authCmd.Use != "auth" {
		t.Errorf("authCmd.Use = %q, want %q", authCmd.Use, "auth")
	}
}

func TestAuthCmdShort(t *testing.T) {
	expected := "Authenticate with Google Drive"
	if authCmd.Short != expected {
		t.Errorf("authCmd.Short = %q, want %q", authCmd.Short, expected)
	}
}

func TestBuildCmdUse(t *testing.T) {
	if buildCmd.Use != "build" {
		t.Errorf("buildCmd.Use = %q, want %q", buildCmd.Use, "build")
	}
}

func TestBuildCmdShort(t *testing.T) {
	expected := "Build nippo site"
	if buildCmd.Short != expected {
		t.Errorf("buildCmd.Short = %q, want %q", buildCmd.Short, expected)
	}
}

func TestDoctorCmdUse(t *testing.T) {
	if doctorCmd.Use != "doctor" {
		t.Errorf("doctorCmd.Use = %q, want %q", doctorCmd.Use, "doctor")
	}
}

func TestDoctorCmdShort(t *testing.T) {
	expected := "Check nippo environment health"
	if doctorCmd.Short != expected {
		t.Errorf("doctorCmd.Short = %q, want %q", doctorCmd.Short, expected)
	}
}

func TestInitCmdUse(t *testing.T) {
	if initCmd.Use != "init" {
		t.Errorf("initCmd.Use = %q, want %q", initCmd.Use, "init")
	}
}

func TestInitCmdShort(t *testing.T) {
	expected := "Initialize nippo command"
	if initCmd.Short != expected {
		t.Errorf("initCmd.Short = %q, want %q", initCmd.Short, expected)
	}
}

// Tests for more subcommand configurations

func TestCleanCmdUse(t *testing.T) {
	if cleanCmd.Use != "clean" {
		t.Errorf("cleanCmd.Use = %q, want %q", cleanCmd.Use, "clean")
	}
}

func TestCleanCmdShort(t *testing.T) {
	expected := "Clean built nippo site files"
	if cleanCmd.Short != expected {
		t.Errorf("cleanCmd.Short = %q, want %q", cleanCmd.Short, expected)
	}
}

func TestDeployCmdUse(t *testing.T) {
	if deployCmd.Use != "deploy" {
		t.Errorf("deployCmd.Use = %q, want %q", deployCmd.Use, "deploy")
	}
}

func TestDeployCmdShort(t *testing.T) {
	expected := "Deploy nippo site"
	if deployCmd.Short != expected {
		t.Errorf("deployCmd.Short = %q, want %q", deployCmd.Short, expected)
	}
}

func TestFormatCmdUse(t *testing.T) {
	if formatCmd.Use != "format" {
		t.Errorf("formatCmd.Use = %q, want %q", formatCmd.Use, "format")
	}
}

func TestFormatCmdShort(t *testing.T) {
	expected := "Manage front-matter in nippo files"
	if formatCmd.Short != expected {
		t.Errorf("formatCmd.Short = %q, want %q", formatCmd.Short, expected)
	}
}

func TestUpdateCmdUse(t *testing.T) {
	if updateCmd.Use != "update" {
		t.Errorf("updateCmd.Use = %q, want %q", updateCmd.Use, "update")
	}
}

func TestUpdateCmdShort(t *testing.T) {
	expected := "Download latest nippo source project"
	if updateCmd.Short != expected {
		t.Errorf("updateCmd.Short = %q, want %q", updateCmd.Short, expected)
	}
}

// Tests for Version variable

func TestVersionVariable(t *testing.T) {
	// Version should be empty by default (set via ldflags during build)
	if Version != "" {
		t.Logf("Version = %q (expected empty or set by ldflags)", Version)
	}
}

// Tests for subcommand registration

func TestSubcommandsRegistered(t *testing.T) {
	commands := rootCmd.Commands()
	commandNames := make(map[string]bool)
	for _, cmd := range commands {
		commandNames[cmd.Name()] = true
	}

	expectedCommands := []string{"auth", "build", "clean", "deploy", "doctor", "format", "init", "update"}
	for _, name := range expectedCommands {
		if !commandNames[name] {
			t.Errorf("expected subcommand %q to be registered", name)
		}
	}
}
