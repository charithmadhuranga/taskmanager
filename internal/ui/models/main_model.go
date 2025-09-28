package models

import (
	"fmt"

	"tappmanager/internal/services"
	"tappmanager/internal/storage"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ViewType represents different UI views
type ViewType int

const (
	ViewProcesses ViewType = iota
	ViewDetails
	ViewStats
	ViewSettings
	ViewHelp
)

// MainModel is the root model for the application
type MainModel struct {
	storage        storage.Storage
	processService *services.ProcessService
	currentView    ViewType
	processes      *ProcessesModel
	details        *DetailsModel
	stats          *StatsModel
	settings       *SettingsModel
	help           *HelpModel
	width          int
	height         int
	quitting       bool
}

// NewMainModel creates a new main model
func NewMainModel(storage storage.Storage, processService *services.ProcessService) *MainModel {
	return &MainModel{
		storage:        storage,
		processService: processService,
		currentView:    ViewProcesses,
		processes:      NewProcessesModel(processService),
		details:        NewDetailsModel(processService),
		stats:          NewStatsModel(processService),
		settings:       NewSettingsModel(storage),
		help:           NewHelpModel(),
		quitting:       false,
	}
}

// Init initializes the model
func (m MainModel) Init() tea.Cmd {
	return tea.Batch(
		m.processes.Init(),
		m.details.Init(),
		m.stats.Init(),
		m.settings.Init(),
		m.help.Init(),
	)
}

// Update handles messages and updates the model
func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Update all sub-models with new size
		*m.processes = m.processes.UpdateSize(msg.Width, msg.Height)
		*m.details = m.details.UpdateSize(msg.Width, msg.Height)
		*m.stats = m.stats.UpdateSize(msg.Width, msg.Height)
		*m.settings = m.settings.UpdateSize(msg.Width, msg.Height)
		*m.help = m.help.UpdateSize(msg.Width, msg.Height)

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "Q", "ctrl+q", "alt+f4", "cmd+q", "ctrl+d":
			m.quitting = true
			return m, tea.Quit

		case "esc":
			// ESC key - return to processes view from any other view
			if m.currentView != ViewProcesses {
				m.currentView = ViewProcesses
				cmd = m.processes.Init()
				cmds = append(cmds, cmd)
			}

		case "p", "P":
			m.currentView = ViewProcesses
			cmd = m.processes.Init()
			cmds = append(cmds, cmd)

		case "d", "D":
			m.currentView = ViewDetails
			cmd = m.details.Init()
			cmds = append(cmds, cmd)

		case "ctrl+s":
			m.currentView = ViewStats
			cmd = m.stats.Init()
			cmds = append(cmds, cmd)

		case "h", "H":
			m.currentView = ViewHelp
			cmd = m.help.Init()
			cmds = append(cmds, cmd)

		case "e", "E":
			m.currentView = ViewSettings
			cmd = m.settings.Init()
			cmds = append(cmds, cmd)

		case "cmd+w":
			// macOS specific - close current view (go back to processes)
			if m.currentView != ViewProcesses {
				m.currentView = ViewProcesses
				cmd = m.processes.Init()
				cmds = append(cmds, cmd)
			}
		}

	case SwitchViewMsg:
		// Handle view switching from sub-models
		m.currentView = msg.View
		switch msg.View {
		case ViewProcesses:
			cmd = m.processes.Init()
		case ViewDetails:
			cmd = m.details.Init()
		case ViewStats:
			cmd = m.stats.Init()
		case ViewSettings:
			cmd = m.settings.Init()
		case ViewHelp:
			cmd = m.help.Init()
		}
		cmds = append(cmds, cmd)
	}

	// Update the current view
	switch m.currentView {
	case ViewProcesses:
		*m.processes, cmd = m.processes.Update(msg)
		cmds = append(cmds, cmd)

	case ViewDetails:
		*m.details, cmd = m.details.Update(msg)
		cmds = append(cmds, cmd)

	case ViewStats:
		*m.stats, cmd = m.stats.Update(msg)
		cmds = append(cmds, cmd)

	case ViewSettings:
		*m.settings, cmd = m.settings.Update(msg)
		cmds = append(cmds, cmd)

	case ViewHelp:
		*m.help, cmd = m.help.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// View renders the current view
func (m MainModel) View() string {
	if m.quitting {
		return "Goodbye!\n"
	}

	// Check if terminal is too small
	if m.width < 80 || m.height < 20 {
		return m.renderSmallTerminalMessage()
	}

	// Create header
	header := m.renderHeader()
	
	// Create content based on current view
	var content string
	switch m.currentView {
	case ViewProcesses:
		content = m.processes.View()
	case ViewDetails:
		content = m.details.View()
	case ViewStats:
		content = m.stats.View()
	case ViewSettings:
		content = m.settings.View()
	case ViewHelp:
		content = m.help.View()
	}

	// Create footer
	footer := m.renderFooter()

	// Calculate available height for content
	headerHeight := 3
	footerHeight := 3
	availableHeight := m.height - headerHeight - footerHeight
	
	// Ensure content fits in available height
	contentStyle := lipgloss.NewStyle().
		Height(availableHeight).
		MaxHeight(availableHeight)

	content = contentStyle.Render(content)

	// Combine all parts
	return lipgloss.JoinVertical(lipgloss.Left, header, content, footer)
}

// renderHeader renders the application header
func (m MainModel) renderHeader() string {
	title := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true).
		Render("Terminal Process Manager")

	nav := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render("[P]rocesses [D]etails [S]tats [E]ettings [H]elp [Q]uit")

	header := lipgloss.JoinHorizontal(lipgloss.Center, title, "  ", nav)
	
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(0, 1).
		Render(header)
}

// renderFooter renders the application footer
func (m MainModel) renderFooter() string {
	viewNames := map[ViewType]string{
		ViewProcesses: "Processes",
		ViewDetails:   "Details", 
		ViewStats:     "Statistics",
		ViewSettings:  "Settings",
		ViewHelp:      "Help",
	}

	status := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render("View: " + viewNames[m.currentView])

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(0, 1).
		Render(status)
}

// renderSmallTerminalMessage renders a message for small terminals
func (m MainModel) renderSmallTerminalMessage() string {
	message := lipgloss.NewStyle().
		Foreground(lipgloss.Color("196")).
		Bold(true).
		Align(lipgloss.Center).
		Render("Terminal too small!\n\nPlease resize your terminal to at least 80x20 characters.\n\nCurrent size: " + 
			lipgloss.NewStyle().Foreground(lipgloss.Color("230")).Render(fmt.Sprintf("%dx%d", m.width, m.height)) + 
			"\n\nPress Ctrl+C to quit.")

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("196")).
		Padding(2, 4).
		Align(lipgloss.Center).
		Render(message)
}
