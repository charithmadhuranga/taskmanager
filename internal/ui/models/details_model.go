package models

import (
	"fmt"
	"strconv"
	"time"

	"tappmanager/internal/models"
	"tappmanager/internal/services"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// DetailsModel handles the process details view
type DetailsModel struct {
	processService *services.ProcessService
	processes      []*models.ProcessInfo
	selectedIndex  int
	width          int
	height         int
	refreshing     bool
}

// NewDetailsModel creates a new details model
func NewDetailsModel(processService *services.ProcessService) *DetailsModel {
	return &DetailsModel{
		processService: processService,
		processes:      []*models.ProcessInfo{},
		selectedIndex:  0,
		refreshing:     false,
	}
}

// Init initializes the model
func (m DetailsModel) Init() tea.Cmd {
	return tea.Batch(
		m.refreshProcesses(),
		m.startRefreshTimer(),
	)
}

// Update handles messages and updates the model
func (m DetailsModel) Update(msg tea.Msg) (DetailsModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.selectedIndex > 0 {
				m.selectedIndex--
			}

		case "down", "j":
			if m.selectedIndex < len(m.processes)-1 {
				m.selectedIndex++
			}

		case "r":
			cmd = m.refreshProcesses()

		case "ctrl+k":
			if len(m.processes) > 0 && m.selectedIndex < len(m.processes) {
				cmd = m.killProcess(m.processes[m.selectedIndex].PID)
			}

		case "f":
			cmd = m.showSearchDialog()

		case "esc":
			// Return to processes view
			cmd = func() tea.Msg { return SwitchViewMsg{View: ViewProcesses} }
		}

	case refreshProcessesMsg:
		m.processes = msg.Processes
		m.refreshing = false
		// Keep selected index within bounds
		if m.selectedIndex >= len(m.processes) {
			m.selectedIndex = len(m.processes) - 1
		}
		if m.selectedIndex < 0 {
			m.selectedIndex = 0
		}

	case refreshTimerMsg:
		cmd = m.refreshProcesses()

	case killProcessMsg:
		if msg.Success {
			// Process killed successfully, select next process
			if m.selectedIndex < len(m.processes)-1 {
				m.selectedIndex++
			} else if m.selectedIndex > 0 {
				m.selectedIndex--
			}
		}

	case SwitchViewMsg:
		// This will be handled by the main model
	}

	return m, cmd
}

// UpdateSize updates the model with new dimensions
func (m DetailsModel) UpdateSize(width, height int) DetailsModel {
	m.width = width
	m.height = height
	return m
}

// View renders the details view
func (m DetailsModel) View() string {
	if m.refreshing {
		return "Refreshing process details...\n"
	}

	if len(m.processes) == 0 {
		return "No processes available.\n"
	}

	if m.selectedIndex >= len(m.processes) {
		return "Invalid process selection.\n"
	}

	proc := m.processes[m.selectedIndex]
	
	// Create details content
	content := m.renderProcessDetails(proc)
	
	// Add navigation info
	nav := m.renderNavigation()
	
	// Combine content and navigation
	fullContent := lipgloss.JoinVertical(lipgloss.Left, content, nav)
	
	// Ensure content fits in available height
	contentStyle := lipgloss.NewStyle().
		Height(m.height - 4). // Account for borders and padding
		MaxHeight(m.height - 4).
		Width(m.width - 4). // Account for borders and padding
		MaxWidth(m.width - 4)

	styledContent := contentStyle.
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(1, 2).
		Render(fullContent)

	return styledContent
}

// renderProcessDetails renders detailed process information
func (m DetailsModel) renderProcessDetails(proc *models.ProcessInfo) string {
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true)

	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Bold(true)

	valueStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("230"))

	// Basic Information
	basicInfo := titleStyle.Render("Basic Information:") + "\n"
	basicInfo += labelStyle.Render("PID:") + " " + valueStyle.Render(strconv.Itoa(int(proc.PID))) + "\n"
	basicInfo += labelStyle.Render("Parent PID:") + " " + valueStyle.Render(strconv.Itoa(int(proc.PPID))) + "\n"
	basicInfo += labelStyle.Render("Name:") + " " + valueStyle.Render(proc.Name) + "\n"
	basicInfo += labelStyle.Render("Status:") + " " + valueStyle.Render(proc.Status) + "\n"
	basicInfo += labelStyle.Render("User:") + " " + valueStyle.Render(proc.Username) + "\n"

	// Resource Usage
	resourceInfo := "\n" + titleStyle.Render("Resource Usage:") + "\n"
	resourceInfo += labelStyle.Render("CPU Usage:") + " " + valueStyle.Render(fmt.Sprintf("%.2f%%", proc.CPU)) + "\n"
	resourceInfo += labelStyle.Render("Memory Usage:") + " " + valueStyle.Render(fmt.Sprintf("%.2f%%", proc.Memory)) + "\n"
	resourceInfo += labelStyle.Render("Memory (Bytes):") + " " + valueStyle.Render(strconv.FormatUint(proc.MemoryBytes, 10)) + "\n"
	resourceInfo += labelStyle.Render("Number of Threads:") + " " + valueStyle.Render(strconv.Itoa(int(proc.NumThreads))) + "\n"
	resourceInfo += labelStyle.Render("Nice Value:") + " " + valueStyle.Render(strconv.Itoa(int(proc.Nice))) + "\n"

	// Process Information
	processInfo := "\n" + titleStyle.Render("Process Information:") + "\n"
	processInfo += labelStyle.Render("Command:") + " " + valueStyle.Render(proc.Command) + "\n"
	processInfo += labelStyle.Render("Working Directory:") + " " + valueStyle.Render(proc.WorkingDir) + "\n"
	processInfo += labelStyle.Render("Create Time:") + " " + valueStyle.Render(proc.CreateTime.Format("2006-01-02 15:04:05")) + "\n"
	processInfo += labelStyle.Render("Running:") + " " + valueStyle.Render(fmt.Sprintf("%t", proc.IsRunning)) + "\n"

	// Navigation
	navigation := "\n" + titleStyle.Render("Navigation:") + "\n"
	navigation += "↑/↓ - Select previous/next process\n"
	navigation += "Ctrl+R - Refresh\n"
	navigation += "Ctrl+K - Kill selected process\n"
	navigation += "Ctrl+F - Search processes\n"
	navigation += "Esc - Return to processes view\n"

	return basicInfo + resourceInfo + processInfo + navigation
}

// renderNavigation renders navigation information
func (m DetailsModel) renderNavigation() string {
	if len(m.processes) == 0 {
		return ""
	}

	navStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Italic(true)

	return navStyle.Render(fmt.Sprintf("Process %d of %d", m.selectedIndex+1, len(m.processes)))
}

// refreshProcesses refreshes the process list
func (m DetailsModel) refreshProcesses() tea.Cmd {
	return func() tea.Msg {
		processes, err := m.processService.GetProcesses()
		if err != nil {
			return refreshProcessesMsg{Processes: []*models.ProcessInfo{}, Error: err}
		}

		return refreshProcessesMsg{Processes: processes}
	}
}

// startRefreshTimer starts the refresh timer
func (m DetailsModel) startRefreshTimer() tea.Cmd {
	return func() tea.Msg {
		time.Sleep(3 * time.Second)
		return refreshTimerMsg{}
	}
}

// killProcess kills the selected process
func (m DetailsModel) killProcess(pid int32) tea.Cmd {
	return func() tea.Msg {
		err := m.processService.KillProcess(pid)
		if err != nil {
			return killProcessMsg{Error: err}
		}
		return killProcessMsg{Success: true}
	}
}

// showSearchDialog shows the search dialog
func (m DetailsModel) showSearchDialog() tea.Cmd {
	return func() tea.Msg {
		// For now, just show a message
		// In a real implementation, this would show a search dialog
		return searchProcessMsg{Query: "search functionality"}
	}
}

// Messages
type searchProcessMsg struct {
	Query string
}
