package tui

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const maxFormatDisplayFiles = 10

// FormatFileStatus represents the result status of processing a file
type FormatFileStatus int

const (
	FormatFileStatusSuccess  FormatFileStatus = iota // Successfully updated
	FormatFileStatusNoChange                         // No changes needed
	FormatFileStatusFailed                           // Failed to process/upload
)

type FormatProgressController struct {
	program        *tea.Program
	running        bool
	cancelled      bool
	done           chan struct{}
	mu             sync.Mutex
	totalFiles     int
	processedFiles int
}

type formatFileInfo struct {
	name    string
	id      string
	status  FormatFileStatus
	message string
}

type formatProgressModel struct {
	spinner        spinner.Model
	progress       progress.Model
	currentFile    string
	currentFileId  string
	processedFiles []formatFileInfo
	totalFiles     int
	processed      int
	startTime      time.Time
	elapsed        time.Duration
	quitting       bool
	cancelled      bool
}

type formatTickMsg time.Time

type formatProgressMsg struct {
	filename string
	fileId   string
	status   FormatFileStatus
	message  string
	total    int
}

type formatCompleteMsg struct{}

func NewFormatProgressController() *FormatProgressController {
	return &FormatProgressController{
		done: make(chan struct{}),
	}
}

func (c *FormatProgressController) Start(totalFiles int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.running {
		return
	}

	c.totalFiles = totalFiles
	c.processedFiles = 0
	c.cancelled = false
	c.done = make(chan struct{})
	c.running = true

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = SpinnerStyle

	p := progress.New(progress.WithSolidFill("#6B8E23"))
	p.Width = 40

	m := formatProgressModel{
		spinner:        s,
		progress:       p,
		processedFiles: make([]formatFileInfo, 0, maxFormatDisplayFiles),
		totalFiles:     totalFiles,
		startTime:      time.Now(),
	}

	c.program = tea.NewProgram(m)

	go func() {
		model, _ := c.program.Run()
		if fm, ok := model.(formatProgressModel); ok && fm.cancelled {
			c.mu.Lock()
			c.cancelled = true
			c.mu.Unlock()
		}
		close(c.done)
	}()
}

func (c *FormatProgressController) UpdateFile(filename string, fileId string, status FormatFileStatus, message string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.running || c.program == nil {
		return
	}

	c.processedFiles++
	c.program.Send(formatProgressMsg{
		filename: filename,
		fileId:   fileId,
		status:   status,
		message:  message,
		total:    c.totalFiles,
	})
}

func (c *FormatProgressController) Stop() {
	c.mu.Lock()
	if !c.running {
		c.mu.Unlock()
		return
	}
	c.running = false
	program := c.program
	done := c.done
	c.mu.Unlock()

	if program != nil {
		program.Send(formatCompleteMsg{})
		<-done
	}
}

func (c *FormatProgressController) IsCancelled() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.cancelled
}

func formatTickEvery() tea.Cmd {
	return tea.Every(time.Second, func(t time.Time) tea.Msg {
		return formatTickMsg(t)
	})
}

func (m formatProgressModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, formatTickEvery())
}

func (m formatProgressModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.quitting = true
			m.cancelled = true
			return m, tea.Quit
		}

	case formatProgressMsg:
		m.currentFile = msg.filename
		m.currentFileId = msg.fileId
		m.processed++
		m.totalFiles = msg.total

		// Add to processed files list, keeping only last maxFormatDisplayFiles
		m.processedFiles = append(m.processedFiles, formatFileInfo{
			name:    msg.filename,
			id:      msg.fileId,
			status:  msg.status,
			message: msg.message,
		})
		if len(m.processedFiles) > maxFormatDisplayFiles {
			m.processedFiles = m.processedFiles[1:]
		}
		return m, nil

	case formatCompleteMsg:
		m.quitting = true
		return m, tea.Quit

	case formatTickMsg:
		m.elapsed = time.Since(m.startTime)
		return m, formatTickEvery()

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd
	}

	return m, nil
}

func (m formatProgressModel) View() string {
	var b strings.Builder

	// Progress bar
	var percent float64
	if m.totalFiles > 0 {
		percent = float64(m.processed) / float64(m.totalFiles)
	}

	// Status indicator for the progress line
	statusSuffix := ""
	if m.quitting {
		if m.cancelled {
			statusSuffix = " " + DimStyle.Render("(interrupted)")
		}
	}

	elapsed := m.elapsed.Truncate(time.Second)
	b.WriteString(fmt.Sprintf("\n  %s %d/%d files %s%s\n",
		m.progress.ViewAs(percent),
		m.processed,
		m.totalFiles,
		DimStyle.Render(fmt.Sprintf("(%s)", elapsed)),
		statusSuffix,
	))

	// Current file with spinner (only when not quitting)
	if !m.quitting && m.currentFile != "" {
		b.WriteString(fmt.Sprintf("\n  %s %s (%s)\n",
			m.spinner.View(),
			InfoStyle.Render(m.currentFile),
			DimStyle.Render(m.currentFileId),
		))
	}

	// List of processed files
	if len(m.processedFiles) > 0 {
		b.WriteString("\n" + DimStyle.Render("  Recently processed:") + "\n")
		for _, file := range m.processedFiles {
			icon, style := getFormatStatusStyle(file.status)
			line := fmt.Sprintf("    %s %s (%s)", style.Render(icon), DimStyle.Render(file.name), DimStyle.Render(file.id))
			if file.status == FormatFileStatusFailed && file.message != "" {
				line += fmt.Sprintf(" - %s", ErrorStyle.Render(file.message))
			}
			b.WriteString(line + "\n")
		}
	}

	return b.String()
}

func getFormatStatusStyle(status FormatFileStatus) (string, lipgloss.Style) {
	switch status {
	case FormatFileStatusSuccess:
		return "✓", SuccessStyle
	case FormatFileStatusNoChange:
		return "○", DimStyle
	case FormatFileStatusFailed:
		return "✗", ErrorStyle
	default:
		return "?", DimStyle
	}
}
