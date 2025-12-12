# Change: Add Bubbletea TUI Framework for Rich UI

## Why

The current CLI uses `fmt.Print` and simple `promptui` prompts, resulting in a
limited user experience. By introducing
[Bubbletea](https://github.com/charmbracelet/bubbletea), we can provide a
richer TUI experience with spinners, progress bars, and interactive inputs.

## What Changes

- Migrate from `promptui` to `bubbletea` + `bubbles`
- Adapt the Presenter layer's View system to Bubbletea's Model-View-Update
  (MVU) architecture
- Introduce rich UI components suited for each command:
  - Spinners (progress indication)
  - Progress bars (build/deploy progress)
  - Text inputs (configuration input)
- Preserve UI output on completion/interruption (no screen clearing)

## Impact

- Affected specs: ui-presentation (new)
- Affected code:
  - `internal/adapter/presenter/view/` - Replace with Bubbletea-based Views
  - `internal/adapter/presenter/*.go` - Update View integration
  - `go.mod` - Add dependencies (bubbletea, bubbles, lipgloss)
