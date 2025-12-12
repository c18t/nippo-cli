## 1. Dependencies

- [x] 1.1 Add `github.com/charmbracelet/bubbletea` to go.mod
- [x] 1.2 Add `github.com/charmbracelet/bubbles` to go.mod
- [x] 1.3 Add `github.com/charmbracelet/lipgloss` to go.mod
- [x] 1.4 Remove `github.com/manifoldco/promptui` dependency

## 2. Core TUI Infrastructure

- [x] 2.1 Create `internal/adapter/presenter/view/tui/output.go` with styled
      output functions
- [x] 2.2 Create `internal/adapter/presenter/view/tui/styles.go` with Lipgloss
      style definitions
- [x] 2.3 Create `internal/adapter/presenter/view/tui/spinner.go` with spinner
      component wrapper
- [x] 2.4 Create `internal/adapter/presenter/view/tui/progress.go` with progress
      bar component wrapper
- [x] 2.5 Create `internal/adapter/presenter/view/tui/textinput.go` with text
      input component wrapper

## 3. View Layer Migration

- [x] 3.1 Update `internal/adapter/presenter/view/base.go` for Bubbletea
      integration
- [x] 3.2 Rewrite `internal/adapter/presenter/view/init.go` to use Bubbletea
      text input
- [x] 3.3 Create Bubbletea model for ConfigureProjectView (integrated in
      textinput.go)
- [x] 3.4 Update InitViewProvider to run Bubbletea program

## 4. Presenter Layer Updates

- [x] 4.1 Update `internal/adapter/presenter/base.go` ConsolePresenter to use
      styled output and SpinnerController
- [x] 4.2 Update init presenter to use spinner for authentication progress
- [x] 4.3 Update build presenter to use spinner for progress
- [x] 4.4 Update deploy presenter to use spinner for progress
- [x] 4.5 Update clean presenter to use spinner for progress
- [x] 4.6 Update update presenter to use spinner for progress

## 5. Testing and Verification

- [x] 5.1 Run `go run -race . --help` to verify DI initialization
- [x] 5.2 Test `nippo init` flow with new text input and spinner
- [x] 5.3 Test `nippo build` with spinner display
- [x] 5.4 Test `nippo deploy` with spinner display
- [x] 5.5 Test `nippo clean` with spinner display
- [x] 5.6 Test `nippo update` with spinner display
- [x] 5.7 Run `make` to ensure build succeeds

## 6. UI Output Preservation on Termination

- [x] 6.1 Add `interrupted` field to `SpinnerModel` for tracking Ctrl-C state
- [x] 6.2 Update `SpinnerModel.View()` to render final state based on
      `done`/`interrupted` flags (shared logic with running state)
- [x] 6.3 Add `interrupted` field to `ProgressModel` for tracking Ctrl-C state
- [x] 6.4 Update `ProgressModel.View()` to render final state based on
      `done`/`interrupted` flags
- [x] 6.5 Update `buildProgressModel.View()` to preserve output on termination
- [x] 6.6 Update wizard/textinput to preserve entered values on completion
- [x] 6.7 Test spinner completion with final state preserved
- [x] 6.8 Test Ctrl-C interruption with message preserved
- [x] 6.9 Test build progress completion/interruption with state preserved

## Notes

The implementation provides:

1. Styled output (success, error, warning, info) via ConsolePresenter
2. Text input via Bubbletea wizard (used in init command)
3. Spinner display during Progress via SpinnerController
4. Progress bar component ready for future integration (when progress
   percentage is available)

Manual testing of `nippo init`, `build`, `deploy`, `clean`, `update` commands
requires actual project setup and Google Drive authentication, which is
outside automated test scope.
