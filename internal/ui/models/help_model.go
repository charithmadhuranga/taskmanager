package models

import (
	"fmt"
	"runtime"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// HelpModel handles the help view
type HelpModel struct {
	width  int
	height int
}

// NewHelpModel creates a new help model
func NewHelpModel() *HelpModel {
	return &HelpModel{}
}

// Init initializes the model
func (m HelpModel) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model
func (m HelpModel) Update(msg tea.Msg) (HelpModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			// Return to processes view
			cmd = func() tea.Msg { return SwitchViewMsg{View: ViewProcesses} }
		}

	case SwitchViewMsg:
		// This will be handled by the main model
	}

	return m, cmd
}

// UpdateSize updates the model with new dimensions
func (m HelpModel) UpdateSize(width, height int) HelpModel {
	m.width = width
	m.height = height
	return m
}

// View renders the help view
func (m HelpModel) View() string {
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true)

	sectionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("62")).
		Bold(true)

	keyStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("230")).
		Bold(true)

	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))

	// Help content
	content := titleStyle.Render("Terminal Process Manager - Help") + "\n\n"
	
	// OS-specific information
	osName := runtime.GOOS
	content += sectionStyle.Render(fmt.Sprintf("Running on: %s", osName)) + "\n\n"
	
	// Navigation
	content += sectionStyle.Render("Navigation:") + "\n"
	content += keyStyle.Render("P") + " - " + descStyle.Render("Switch to Processes view") + "\n"
	content += keyStyle.Render("D") + " - " + descStyle.Render("Switch to Details view") + "\n"
	content += keyStyle.Render("Ctrl+S") + " - " + descStyle.Render("Switch to Statistics view") + "\n"
	content += keyStyle.Render("H") + " - " + descStyle.Render("Show this help") + "\n"
	content += keyStyle.Render("E") + " - " + descStyle.Render("Switch to Settings view") + "\n"
	
	// OS-specific quit shortcuts
	switch osName {
	case "windows":
		content += keyStyle.Render("Ctrl+Q") + " - " + descStyle.Render("Quit application") + "\n"
		content += keyStyle.Render("Alt+F4") + " - " + descStyle.Render("Quit application") + "\n"
	case "darwin":
		content += keyStyle.Render("Cmd+Q") + " - " + descStyle.Render("Quit application") + "\n"
		content += keyStyle.Render("Cmd+W") + " - " + descStyle.Render("Close current view") + "\n"
	case "linux":
		content += keyStyle.Render("Ctrl+D") + " - " + descStyle.Render("Quit application") + "\n"
	}
	content += keyStyle.Render("Q") + " - " + descStyle.Render("Quit application") + "\n"
	content += keyStyle.Render("Esc") + " - " + descStyle.Render("Return to processes view") + "\n\n"

	// Processes View
	content += sectionStyle.Render("Processes View:") + "\n"
	content += keyStyle.Render("↑/↓ or J/K") + " - " + descStyle.Render("Navigate up/down") + "\n"
	content += keyStyle.Render("R") + " - " + descStyle.Render("Refresh process list") + "\n"
	content += keyStyle.Render("Ctrl+K") + " - " + descStyle.Render("Kill selected process") + "\n"
	content += keyStyle.Render("F") + " - " + descStyle.Render("Toggle system processes filter") + "\n"
	content += keyStyle.Render("Ctrl+F") + " - " + descStyle.Render("Search processes (cycle through terms)") + "\n"
	content += keyStyle.Render("Ctrl+Shift+F") + " - " + descStyle.Render("Clear search filter") + "\n"
	content += keyStyle.Render("S") + " - " + descStyle.Render("Toggle system processes display") + "\n"
	content += keyStyle.Render("Ctrl+R") + " - " + descStyle.Render("Reset all filters and refresh") + "\n"
	content += keyStyle.Render("Ctrl+Shift+S") + " - " + descStyle.Render("Reset sort to default (CPU desc)") + "\n"
	content += keyStyle.Render("O") + " - " + descStyle.Render("Sort by CPU usage") + "\n"
	content += keyStyle.Render("M") + " - " + descStyle.Render("Sort by memory usage") + "\n"
	content += keyStyle.Render("Ctrl+P") + " - " + descStyle.Render("Sort by PID") + "\n"
	content += keyStyle.Render("N") + " - " + descStyle.Render("Sort by name") + "\n"
	content += keyStyle.Render("T") + " - " + descStyle.Render("Sort by status") + "\n"
	content += keyStyle.Render("U") + " - " + descStyle.Render("Sort by user") + "\n"
	content += keyStyle.Render("Ctrl+T") + " - " + descStyle.Render("Sort by threads") + "\n"
	content += keyStyle.Render("Ctrl+N") + " - " + descStyle.Render("Sort by nice value") + "\n"
	content += keyStyle.Render("Enter") + " - " + descStyle.Render("View process details") + "\n\n"

	// Details View
	content += sectionStyle.Render("Details View:") + "\n"
	content += keyStyle.Render("↑/↓") + " - " + descStyle.Render("Select previous/next process") + "\n"
	content += keyStyle.Render("Ctrl+R") + " - " + descStyle.Render("Refresh process details") + "\n"
	content += keyStyle.Render("Ctrl+K") + " - " + descStyle.Render("Kill selected process") + "\n"
	content += keyStyle.Render("Ctrl+F") + " - " + descStyle.Render("Search processes") + "\n"
	content += keyStyle.Render("Esc") + " - " + descStyle.Render("Return to processes view") + "\n\n"

	// Statistics View
	content += sectionStyle.Render("Statistics View:") + "\n"
	content += keyStyle.Render("Ctrl+R") + " - " + descStyle.Render("Refresh statistics") + "\n"
	content += keyStyle.Render("Ctrl+E") + " - " + descStyle.Render("Export statistics") + "\n"
	content += keyStyle.Render("Esc") + " - " + descStyle.Render("Return to processes view") + "\n\n"

	// Settings View
	content += sectionStyle.Render("Settings View:") + "\n"
	content += descStyle.Render("Configure refresh rate, filters, and display options") + "\n"
	content += keyStyle.Render("Esc") + " - " + descStyle.Render("Return to processes view") + "\n\n"

	// General
	content += sectionStyle.Render("General:") + "\n"
	content += keyStyle.Render("Arrow Keys") + " - " + descStyle.Render("Navigate") + "\n"
	content += keyStyle.Render("Enter") + " - " + descStyle.Render("Select/Activate") + "\n"
	content += keyStyle.Render("Tab") + " - " + descStyle.Render("Next field") + "\n"
	content += keyStyle.Render("Esc") + " - " + descStyle.Render("Cancel/Back") + "\n\n"

	// Process Management
	content += sectionStyle.Render("Process Management:") + "\n"
	content += descStyle.Render("• Real-time process monitoring") + "\n"
	content += descStyle.Render("• Advanced filtering and sorting") + "\n"
	content += descStyle.Render("• Process termination") + "\n"
	content += descStyle.Render("• Detailed process information") + "\n"
	content += descStyle.Render("• System statistics") + "\n"
	content += descStyle.Render("• Data export and backup") + "\n\n"

	// Features
	content += sectionStyle.Render("Features:") + "\n"
	content += descStyle.Render("• Cross-platform (macOS, Linux, Windows)") + "\n"
	content += descStyle.Render("• Real-time monitoring with auto-refresh") + "\n"
	content += descStyle.Render("• Advanced filtering by CPU, memory, status, user") + "\n"
	content += descStyle.Render("• Process tree visualization") + "\n"
	content += descStyle.Render("• Statistics and reporting") + "\n"
	content += descStyle.Render("• Data persistence and backup") + "\n"
	content += descStyle.Render("• Keyboard shortcuts for efficiency") + "\n\n"

	// Version
	content += sectionStyle.Render("Version:") + " " + descStyle.Render("1.0.0") + "\n"

	// Controls
	controls := "\n" + sectionStyle.Render("Controls:") + "\n"
	controls += keyStyle.Render("Esc") + " - " + descStyle.Render("Return to processes view") + "\n"

	// Combine content and controls
	fullContent := lipgloss.JoinVertical(lipgloss.Left, content, controls)
	
	// Add borders and styling
	styledContent := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(1, 2).
		Render(fullContent)

	return styledContent
}
