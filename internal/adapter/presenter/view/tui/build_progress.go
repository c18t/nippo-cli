package tui

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

const maxDisplayFiles = 10

type BuildProgressController struct {
	program        *tea.Program
	running        bool
	cancelled      bool
	done           chan struct{}
	mu             sync.Mutex
	totalFiles     int
	processedFiles int
}

type fileInfo struct {
	name string
	id   string
}

type buildProgressModel struct {
	spinner        spinner.Model
	progress       progress.Model
	currentFile    fileInfo
	processedFiles []fileInfo
	totalFiles     int
	processed      int
	startTime      time.Time
	elapsed        time.Duration
	quitting       bool
	cancelled      bool
}

type buildTickMsg time.Time

type buildProgressMsg struct {
	filename string
	fileId   string
	total    int
}

type buildCompleteMsg struct{}

func NewBuildProgressController() *BuildProgressController {
	return &BuildProgressController{
		done: make(chan struct{}),
	}
}

func (c *BuildProgressController) Start(totalFiles int) {
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

	m := buildProgressModel{
		spinner:        s,
		progress:       p,
		processedFiles: make([]fileInfo, 0, maxDisplayFiles),
		totalFiles:     totalFiles,
		startTime:      time.Now(),
	}

	c.program = tea.NewProgram(m)

	go func() {
		model, _ := c.program.Run()
		if bm, ok := model.(buildProgressModel); ok && bm.cancelled {
			c.mu.Lock()
			c.cancelled = true
			c.mu.Unlock()
		}
		close(c.done)
	}()
}

func (c *BuildProgressController) UpdateFile(filename string, fileId string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.running || c.program == nil {
		return
	}

	c.processedFiles++
	c.program.Send(buildProgressMsg{
		filename: filename,
		fileId:   fileId,
		total:    c.totalFiles,
	})
}

func (c *BuildProgressController) Stop() {
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
		program.Send(buildCompleteMsg{})
		<-done
	}
}

func (c *BuildProgressController) IsCancelled() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.cancelled
}

func buildTickEvery() tea.Cmd {
	return tea.Every(time.Second, func(t time.Time) tea.Msg {
		return buildTickMsg(t)
	})
}

func (m buildProgressModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, buildTickEvery())
}

func (m buildProgressModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.quitting = true
			m.cancelled = true
			return m, tea.Quit
		}

	case buildProgressMsg:
		m.currentFile = fileInfo{name: msg.filename, id: msg.fileId}
		m.processed++
		m.totalFiles = msg.total

		// Add to processed files list, keeping only last maxDisplayFiles
		m.processedFiles = append(m.processedFiles, fileInfo{name: msg.filename, id: msg.fileId})
		if len(m.processedFiles) > maxDisplayFiles {
			m.processedFiles = m.processedFiles[1:]
		}
		return m, nil

	case buildCompleteMsg:
		m.quitting = true
		return m, tea.Quit

	case buildTickMsg:
		m.elapsed = time.Since(m.startTime)
		return m, buildTickEvery()

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

func (m buildProgressModel) View() string {
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
	if !m.quitting && m.currentFile.name != "" {
		b.WriteString(fmt.Sprintf("\n  %s %s (%s)\n",
			m.spinner.View(),
			InfoStyle.Render(m.currentFile.name),
			DimStyle.Render(m.currentFile.id),
		))
	}

	// List of processed files
	if len(m.processedFiles) > 0 {
		b.WriteString("\n" + DimStyle.Render("  Recently processed:") + "\n")
		for _, file := range m.processedFiles {
			b.WriteString(fmt.Sprintf("    %s %s (%s)\n",
				SuccessStyle.Render("✓"),
				DimStyle.Render(file.name),
				DimStyle.Render(file.id),
			))
		}
	}

	return b.String()
}
