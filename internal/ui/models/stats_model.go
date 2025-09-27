package models

import (
	"fmt"
	"time"

	"tappmanager/internal/models"
	"tappmanager/internal/services"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// StatsModel handles the statistics view
type StatsModel struct {
	processService *services.ProcessService
	processes      []*models.ProcessInfo
	width          int
	height         int
	refreshing     bool
}

// NewStatsModel creates a new stats model
func NewStatsModel(processService *services.ProcessService) *StatsModel {
	return &StatsModel{
		processService: processService,
		processes:      []*models.ProcessInfo{},
		refreshing:     false,
	}
}

// Init initializes the model
func (m StatsModel) Init() tea.Cmd {
	return tea.Batch(
		m.refreshProcesses(),
		m.startRefreshTimer(),
	)
}

// Update handles messages and updates the model
func (m StatsModel) Update(msg tea.Msg) (StatsModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "r":
			cmd = m.refreshProcesses()

		case "e":
			cmd = m.exportStats()

		case "esc":
			// Return to processes view
			cmd = func() tea.Msg { return SwitchViewMsg{View: ViewProcesses} }
		}

	case refreshProcessesMsg:
		m.processes = msg.Processes
		m.refreshing = false

	case refreshTimerMsg:
		cmd = m.refreshProcesses()

	case exportStatsMsg:
		// Export completed
		cmd = tea.Printf("Statistics exported: %s", msg.Filename)

	case SwitchViewMsg:
		// This will be handled by the main model
	}

	return m, cmd
}

// UpdateSize updates the model with new dimensions
func (m StatsModel) UpdateSize(width, height int) StatsModel {
	m.width = width
	m.height = height
	return m
}

// View renders the statistics view
func (m StatsModel) View() string {
	if m.refreshing {
		return "Refreshing statistics...\n"
	}

	if len(m.processes) == 0 {
		return "No process data available.\n"
	}

	// Get process statistics
	stats := m.processService.GetProcessStats(m.processes)

	// Create statistics content
	content := m.renderStatistics(stats)
	
	// Add navigation info
	nav := m.renderNavigation()
	
	// Combine content and navigation
	fullContent := lipgloss.JoinVertical(lipgloss.Left, content, nav)
	
	// Ensure content fits in available height and width
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

// renderStatistics renders the process statistics
func (m StatsModel) renderStatistics(stats map[string]interface{}) string {
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true)

	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Bold(true)

	valueStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("230"))

	// Extract statistics
	totalProcesses := stats["total_processes"].(int)
	runningProcesses := stats["running_processes"].(int)
	totalCPU := stats["total_cpu"].(float64)
	totalMemory := stats["total_memory"].(float64)
	statusCounts := stats["status_counts"].(map[string]int)
	userCounts := stats["user_counts"].(map[string]int)

	// Overview
	overview := titleStyle.Render("Overview:") + "\n"
	overview += labelStyle.Render("Total Processes:") + " " + valueStyle.Render(fmt.Sprintf("%d", totalProcesses)) + "\n"
	overview += labelStyle.Render("Running Processes:") + " " + valueStyle.Render(fmt.Sprintf("%d", runningProcesses)) + "\n"
	overview += labelStyle.Render("Stopped Processes:") + " " + valueStyle.Render(fmt.Sprintf("%d", totalProcesses-runningProcesses)) + "\n"
	overview += labelStyle.Render("Total CPU Usage:") + " " + valueStyle.Render(fmt.Sprintf("%.2f%%", totalCPU)) + "\n"
	overview += labelStyle.Render("Total Memory Usage:") + " " + valueStyle.Render(fmt.Sprintf("%.2f%%", totalMemory)) + "\n"

	// Process Status Distribution
	statusInfo := "\n" + titleStyle.Render("Process Status Distribution:") + "\n"
	for status, count := range statusCounts {
		percentage := float64(count) / float64(totalProcesses) * 100
		statusInfo += labelStyle.Render(status) + ": " + valueStyle.Render(fmt.Sprintf("%d (%.1f%%)", count, percentage)) + "\n"
	}

	// Top Users by Process Count
	userInfo := "\n" + titleStyle.Render("Top Users by Process Count:") + "\n"
	userCount := 0
	for user, count := range userCounts {
		if userCount >= 10 {
			break
		}
		percentage := float64(count) / float64(totalProcesses) * 100
		userInfo += labelStyle.Render(user) + ": " + valueStyle.Render(fmt.Sprintf("%d (%.1f%%)", count, percentage)) + "\n"
		userCount++
	}

	// Top Processes by CPU and Memory
	topCPUProcesses := m.getTopProcesses("cpu", 5)
	topMemoryProcesses := m.getTopProcesses("memory", 5)

	cpuInfo := "\n" + titleStyle.Render("Top 5 Processes by CPU Usage:") + "\n"
	for i, proc := range topCPUProcesses {
		cpuInfo += fmt.Sprintf("%d. %s (PID: %d) - %.2f%%\n", i+1, proc.Name, proc.PID, proc.CPU)
	}

	memInfo := "\n" + titleStyle.Render("Top 5 Processes by Memory Usage:") + "\n"
	for i, proc := range topMemoryProcesses {
		memInfo += fmt.Sprintf("%d. %s (PID: %d) - %.2f%%\n", i+1, proc.Name, proc.PID, proc.Memory)
	}

	// System Information
	systemInfo := "\n" + titleStyle.Render("System Information:") + "\n"
	systemInfo += labelStyle.Render("Current Time:") + " " + valueStyle.Render(time.Now().Format("2006-01-02 15:04:05")) + "\n"
	systemInfo += labelStyle.Render("Process Count:") + " " + valueStyle.Render(fmt.Sprintf("%d", totalProcesses)) + "\n"
	systemInfo += labelStyle.Render("Average CPU per Process:") + " " + valueStyle.Render(fmt.Sprintf("%.2f%%", totalCPU/float64(totalProcesses))) + "\n"
	systemInfo += labelStyle.Render("Average Memory per Process:") + " " + valueStyle.Render(fmt.Sprintf("%.2f%%", totalMemory/float64(totalProcesses))) + "\n"

	// Controls
	controls := "\n" + titleStyle.Render("Controls:") + "\n"
	controls += "Ctrl+R - Refresh statistics\n"
	controls += "Ctrl+E - Export statistics\n"
	controls += "Esc - Return to processes view\n"

	return overview + statusInfo + userInfo + cpuInfo + memInfo + systemInfo + controls
}

// renderNavigation renders navigation information
func (m StatsModel) renderNavigation() string {
	navStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Italic(true)

	return navStyle.Render("Statistics updated every 5 seconds")
}

// getTopProcesses returns the top N processes by the specified field
func (m StatsModel) getTopProcesses(field string, n int) []*models.ProcessInfo {
	// Create a copy of processes for sorting
	processes := make([]*models.ProcessInfo, len(m.processes))
	copy(processes, m.processes)

	// Sort by field
	switch field {
	case "cpu":
		for i := 0; i < len(processes)-1; i++ {
			for j := i + 1; j < len(processes); j++ {
				if processes[i].CPU < processes[j].CPU {
					processes[i], processes[j] = processes[j], processes[i]
				}
			}
		}
	case "memory":
		for i := 0; i < len(processes)-1; i++ {
			for j := i + 1; j < len(processes); j++ {
				if processes[i].Memory < processes[j].Memory {
					processes[i], processes[j] = processes[j], processes[i]
				}
			}
		}
	}

	// Return top N
	if n > len(processes) {
		n = len(processes)
	}
	return processes[:n]
}

// refreshProcesses refreshes the process list
func (m StatsModel) refreshProcesses() tea.Cmd {
	return func() tea.Msg {
		processes, err := m.processService.GetProcesses()
		if err != nil {
			return refreshProcessesMsg{Processes: []*models.ProcessInfo{}, Error: err}
		}

		return refreshProcessesMsg{Processes: processes}
	}
}

// startRefreshTimer starts the refresh timer
func (m StatsModel) startRefreshTimer() tea.Cmd {
	return func() tea.Msg {
		time.Sleep(5 * time.Second)
		return refreshTimerMsg{}
	}
}

// exportStats exports the current statistics
func (m StatsModel) exportStats() tea.Cmd {
	return func() tea.Msg {
		// This would integrate with the storage service to export statistics
		filename := fmt.Sprintf("process_stats_%s.txt", time.Now().Format("20060102_150405"))
		return exportStatsMsg{Filename: filename}
	}
}

// Messages
type exportStatsMsg struct {
	Filename string
}
