package models

import (
	"fmt"
	"strconv"

	"tappmanager/internal/models"
	"tappmanager/internal/storage"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// SettingsModel handles the settings view
type SettingsModel struct {
	storage storage.Storage
	config  *AppConfig
	width   int
	height  int
}

// NewSettingsModel creates a new settings model
func NewSettingsModel(storage storage.Storage) *SettingsModel {
	return &SettingsModel{
		storage: storage,
		config:  NewAppConfig(),
	}
}

// Init initializes the model
func (m SettingsModel) Init() tea.Cmd {
	return func() tea.Msg {
		config, err := m.storage.LoadConfig()
		if err != nil {
			return loadConfigMsg{Error: err}
		}
		return loadConfigMsg{Config: config}
	}
}

// Update handles messages and updates the model
func (m SettingsModel) Update(msg tea.Msg) (SettingsModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			// Return to processes view
			cmd = func() tea.Msg { return SwitchViewMsg{View: ViewProcesses} }
		}

	case loadConfigMsg:
		if msg.Error == nil {
			// Convert from internal models to UI models
			m.config = &AppConfig{
				RefreshRate: msg.Config.RefreshRate,
				ShowSystem:  msg.Config.ShowSystem,
				DefaultSort: ProcessSort{
					Field: msg.Config.DefaultSort.Field,
					Order: msg.Config.DefaultSort.Order,
				},
				DefaultFilter: ProcessFilter{
					SearchTerm: msg.Config.DefaultFilter.SearchTerm,
					MinCPU:     msg.Config.DefaultFilter.MinCPU,
					MaxCPU:     msg.Config.DefaultFilter.MaxCPU,
					MinMemory:  msg.Config.DefaultFilter.MinMemory,
					MaxMemory:  msg.Config.DefaultFilter.MaxMemory,
					Status:     msg.Config.DefaultFilter.Status,
					Username:   msg.Config.DefaultFilter.Username,
					ShowSystem: msg.Config.DefaultFilter.ShowSystem,
				},
				AutoRefresh: msg.Config.AutoRefresh,
				Theme:       msg.Config.Theme,
				DataDir:     msg.Config.DataDir,
				Version:     msg.Config.Version,
				CreatedAt:   msg.Config.CreatedAt,
				UpdatedAt:   msg.Config.UpdatedAt,
			}
		}

	case SwitchViewMsg:
		// This will be handled by the main model
	}

	return m, cmd
}

// UpdateSize updates the model with new dimensions
func (m SettingsModel) UpdateSize(width, height int) SettingsModel {
	m.width = width
	m.height = height
	return m
}

// View renders the settings view
func (m SettingsModel) View() string {
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true)

	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Bold(true)

	valueStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("230"))

	// Settings content
	content := titleStyle.Render("Process Manager Settings") + "\n\n"
	
	// Refresh Rate
	content += labelStyle.Render("Refresh Rate (seconds):") + " " + valueStyle.Render(strconv.Itoa(m.config.RefreshRate)) + "\n"
	
	// Show System Processes
	content += labelStyle.Render("Show System Processes:") + " " + valueStyle.Render(fmt.Sprintf("%t", m.config.ShowSystem)) + "\n"
	
	// Default Sort Field
	content += labelStyle.Render("Default Sort Field:") + " " + valueStyle.Render(m.config.DefaultSort.Field) + "\n"
	
	// Default Sort Order
	content += labelStyle.Render("Default Sort Order:") + " " + valueStyle.Render(m.config.DefaultSort.Order) + "\n"
	
	// Min CPU Filter
	content += labelStyle.Render("Min CPU Filter:") + " " + valueStyle.Render(fmt.Sprintf("%.2f", m.config.DefaultFilter.MinCPU)) + "\n"
	
	// Max CPU Filter
	content += labelStyle.Render("Max CPU Filter:") + " " + valueStyle.Render(fmt.Sprintf("%.2f", m.config.DefaultFilter.MaxCPU)) + "\n"
	
	// Min Memory Filter
	content += labelStyle.Render("Min Memory Filter:") + " " + valueStyle.Render(fmt.Sprintf("%.2f", m.config.DefaultFilter.MinMemory)) + "\n"
	
	// Max Memory Filter
	content += labelStyle.Render("Max Memory Filter:") + " " + valueStyle.Render(fmt.Sprintf("%.2f", m.config.DefaultFilter.MaxMemory)) + "\n"
	
	// Auto Refresh
	content += labelStyle.Render("Auto Refresh:") + " " + valueStyle.Render(fmt.Sprintf("%t", m.config.AutoRefresh)) + "\n"
	
	// Theme
	content += labelStyle.Render("Theme:") + " " + valueStyle.Render(m.config.Theme) + "\n"
	
	// Data Directory
	content += labelStyle.Render("Data Directory:") + " " + valueStyle.Render(m.config.DataDir) + "\n"

	// Controls
	controls := "\n" + titleStyle.Render("Controls:") + "\n"
	controls += "Esc - Return to processes view\n"
	controls += "Note: Settings are read-only in this demo\n"

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

// Messages
type loadConfigMsg struct {
	Config *models.AppConfig
	Error  error
}
