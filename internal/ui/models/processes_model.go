package models

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"tappmanager/internal/models"
	"tappmanager/internal/services"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ProcessesModel handles the processes view
type ProcessesModel struct {
	processService *services.ProcessService
	processes      []*models.ProcessInfo
	filter         *models.ProcessFilter
	sort           *models.ProcessSort
	selectedIndex  int
	width          int
	height         int
	showSystem     bool
	refreshing     bool
}

// NewProcessesModel creates a new processes model
func NewProcessesModel(processService *services.ProcessService) *ProcessesModel {
	return &ProcessesModel{
		processService: processService,
		processes:      []*models.ProcessInfo{},
		filter:         &models.ProcessFilter{},
		sort:           &models.ProcessSort{Field: "cpu", Order: "desc"},
		selectedIndex:  0,
		showSystem:     false,
		refreshing:     false,
	}
}

// Init initializes the model
func (m ProcessesModel) Init() tea.Cmd {
	return tea.Batch(
		m.refreshProcesses(),
		m.startRefreshTimer(),
	)
}

// Update handles messages and updates the model
func (m ProcessesModel) Update(msg tea.Msg) (ProcessesModel, tea.Cmd) {
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
			cmd = m.showFilterDialog()

		case "s":
			m.showSystem = !m.showSystem
			m.filter.ShowSystem = m.showSystem
			cmd = m.refreshProcesses()

		case "o":
			m.sortByField("cpu")
			cmd = m.refreshProcesses()

		case "m":
			m.sortByField("memory")
			cmd = m.refreshProcesses()

		case "p":
			m.sortByField("pid")
			cmd = m.refreshProcesses()

		case "n":
			m.sortByField("name")
			cmd = m.refreshProcesses()

		case "t":
			m.sortByField("status")
			cmd = m.refreshProcesses()

		case "enter":
			if len(m.processes) > 0 && m.selectedIndex < len(m.processes) {
				// Switch to details view
				cmd = tea.Sequence(
					tea.Printf("Switching to details view for process %d", m.processes[m.selectedIndex].PID),
					func() tea.Msg { return SwitchViewMsg{View: ViewDetails} },
				)
			}
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

	case SwitchViewMsg:
		// This will be handled by the main model
	}

	return m, cmd
}

// UpdateSize updates the model with new dimensions
func (m ProcessesModel) UpdateSize(width, height int) ProcessesModel {
	m.width = width
	m.height = height
	return m
}

// View renders the processes view
func (m ProcessesModel) View() string {
	if m.refreshing {
		return "Refreshing processes...\n"
	}

	if len(m.processes) == 0 {
		return "No processes found.\n"
	}

	// Create table header
	header := m.renderTableHeader()
	
	// Create table rows
	rows := m.renderTableRows()
	
	// Create separator line
	colWidths := m.calculateColumnWidths()
	separator := m.renderSeparator(colWidths)
	
	// Create table
	table := lipgloss.JoinVertical(lipgloss.Left, header, separator, rows)
	
	// Ensure table fits in available height and width
	tableStyle := lipgloss.NewStyle().
		Height(m.height - 4). // Account for borders and padding
		MaxHeight(m.height - 4).
		Width(m.width - 4). // Account for borders and padding
		MaxWidth(m.width - 4)

	styledTable := tableStyle.
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(0, 1).
		Render(table)

	return styledTable
}

// renderTableHeader renders the table header
func (m ProcessesModel) renderTableHeader() string {
	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true).
		Align(lipgloss.Center)

	// Calculate column widths based on terminal width
	colWidths := m.calculateColumnWidths()
	
	headers := []string{"PID", "Name", "Status", "CPU%", "Memory%", "User", "Threads", "Nice"}
	
	var headerCells []string
	for i, header := range headers {
		width := colWidths[i]
		cell := headerStyle.Width(width).Align(lipgloss.Center).Render(header)
		headerCells = append(headerCells, cell)
	}

	return lipgloss.JoinHorizontal(lipgloss.Left, headerCells...)
}

// renderTableRows renders the table rows
func (m ProcessesModel) renderTableRows() string {
	var rows []string
	
	// Calculate column widths
	colWidths := m.calculateColumnWidths()
	
	for i, proc := range m.processes {
		rowStyle := lipgloss.NewStyle()
		if i == m.selectedIndex {
			rowStyle = rowStyle.
				Background(lipgloss.Color("62")).
				Foreground(lipgloss.Color("230"))
		}

		// Color coding for CPU usage
		cpuColor := "white"
		if proc.CPU > 50 {
			cpuColor = "red"
		} else if proc.CPU > 20 {
			cpuColor = "yellow"
		} else if proc.CPU > 5 {
			cpuColor = "green"
		}

		// Color coding for memory usage
		memColor := "white"
		if proc.Memory > 50 {
			memColor = "red"
		} else if proc.Memory > 20 {
			memColor = "yellow"
		} else if proc.Memory > 5 {
			memColor = "green"
		}

		// Color coding for status
		statusColor := "white"
		switch proc.Status {
		case "running", "R":
			statusColor = "green"
		case "sleeping", "S":
			statusColor = "blue"
		case "zombie", "Z":
			statusColor = "red"
		case "stopped", "T":
			statusColor = "yellow"
		}

		// Truncate and format data based on column widths
		pidStr := strconv.Itoa(int(proc.PID))
		name := m.truncateString(proc.Name, colWidths[1]-2)
		status := m.truncateString(proc.Status, colWidths[2]-2)
		cpuStr := fmt.Sprintf("%.2f", proc.CPU)
		memStr := fmt.Sprintf("%.2f", proc.Memory)
		user := m.truncateString(proc.Username, colWidths[5]-2)
		threadsStr := strconv.Itoa(int(proc.NumThreads))
		niceStr := strconv.Itoa(int(proc.Nice))

		cells := []string{
			rowStyle.Width(colWidths[0]).Align(lipgloss.Right).Render(pidStr),
			rowStyle.Width(colWidths[1]).Align(lipgloss.Left).Render(name),
			rowStyle.Width(colWidths[2]).Align(lipgloss.Center).Foreground(lipgloss.Color(statusColor)).Render(status),
			rowStyle.Width(colWidths[3]).Align(lipgloss.Right).Foreground(lipgloss.Color(cpuColor)).Render(cpuStr),
			rowStyle.Width(colWidths[4]).Align(lipgloss.Right).Foreground(lipgloss.Color(memColor)).Render(memStr),
			rowStyle.Width(colWidths[5]).Align(lipgloss.Center).Render(user),
			rowStyle.Width(colWidths[6]).Align(lipgloss.Right).Render(threadsStr),
			rowStyle.Width(colWidths[7]).Align(lipgloss.Right).Render(niceStr),
		}

		rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Left, cells...))
	}

	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}

// refreshProcesses refreshes the process list
func (m ProcessesModel) refreshProcesses() tea.Cmd {
	return func() tea.Msg {
		processes, err := m.processService.GetProcesses()
		if err != nil {
			return refreshProcessesMsg{Processes: []*models.ProcessInfo{}, Error: err}
		}

		// Apply filters
		filteredProcesses := m.processService.FilterProcesses(processes, m.filter)
		
		// Apply sorting
		m.processService.SortProcesses(filteredProcesses, m.sort)

		return refreshProcessesMsg{Processes: filteredProcesses}
	}
}

// startRefreshTimer starts the refresh timer
func (m ProcessesModel) startRefreshTimer() tea.Cmd {
	return func() tea.Msg {
		time.Sleep(2 * time.Second)
		return refreshTimerMsg{}
	}
}

// killProcess kills the selected process
func (m ProcessesModel) killProcess(pid int32) tea.Cmd {
	return func() tea.Msg {
		err := m.processService.KillProcess(pid)
		if err != nil {
			return killProcessMsg{Error: err}
		}
		return killProcessMsg{Success: true}
	}
}

// showFilterDialog shows the filter dialog
func (m ProcessesModel) showFilterDialog() tea.Cmd {
	return func() tea.Msg {
		// For now, just toggle system processes
		// In a real implementation, this would show a filter dialog
		return filterProcessesMsg{Filter: m.filter}
	}
}

// sortByField sorts processes by the specified field
func (m ProcessesModel) sortByField(field string) {
	if m.sort.Field == field {
		// Toggle sort order
		if m.sort.Order == "asc" {
			m.sort.Order = "desc"
		} else {
			m.sort.Order = "asc"
		}
	} else {
		m.sort.Field = field
		m.sort.Order = "desc"
	}
}

// calculateColumnWidths calculates appropriate column widths based on terminal width
func (m ProcessesModel) calculateColumnWidths() []int {
	// Minimum column widths
	minWidths := []int{8, 20, 10, 8, 8, 12, 8, 6} // PID, Name, Status, CPU%, Memory%, User, Threads, Nice
	
	// Available width (account for borders and padding)
	availableWidth := m.width - 4 // Account for borders
	
	// Calculate total minimum width
	totalMinWidth := 0
	for _, w := range minWidths {
		totalMinWidth += w
	}
	
	// If terminal is too narrow, use minimum widths
	if availableWidth < totalMinWidth {
		return minWidths
	}
	
	// Calculate extra width to distribute
	extraWidth := availableWidth - totalMinWidth
	
	// Distribute extra width proportionally, with Name getting the most
	colWidths := make([]int, len(minWidths))
	copy(colWidths, minWidths)
	
	// Give extra space to Name column (index 1) and User column (index 5)
	nameExtra := extraWidth * 3 / 5  // 60% of extra width
	userExtra := extraWidth * 1 / 5  // 20% of extra width
	otherExtra := extraWidth * 1 / 5 // 20% of extra width
	
	colWidths[1] += nameExtra  // Name
	colWidths[5] += userExtra  // User
	
	// Distribute remaining extra width to other columns
	remainingExtra := otherExtra
	for i := range colWidths {
		if i != 1 && i != 5 && remainingExtra > 0 {
			colWidths[i] += 1
			remainingExtra--
		}
	}
	
	return colWidths
}

// truncateString truncates a string to fit within the specified width
func (m ProcessesModel) truncateString(s string, maxWidth int) string {
	if maxWidth <= 0 {
		return ""
	}
	
	if len(s) <= maxWidth {
		return s
	}
	
	if maxWidth <= 3 {
		return "..."
	}
	
	return s[:maxWidth-3] + "..."
}

// renderSeparator renders a separator line between header and rows
func (m ProcessesModel) renderSeparator(colWidths []int) string {
	var separatorCells []string
	
	for _, width := range colWidths {
		separator := lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Width(width).
			Render(strings.Repeat("─", width))
		separatorCells = append(separatorCells, separator)
	}
	
	return lipgloss.JoinHorizontal(lipgloss.Left, separatorCells...)
}

// Messages
type refreshProcessesMsg struct {
	Processes []*models.ProcessInfo
	Error     error
}

type refreshTimerMsg struct{}

type killProcessMsg struct {
	Success bool
	Error   error
}

type filterProcessesMsg struct {
	Filter *models.ProcessFilter
}

type SwitchViewMsg struct {
	View ViewType
}
