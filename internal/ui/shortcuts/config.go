package shortcuts

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// ShortcutConfig represents the configuration for shortcuts
type ShortcutConfig struct {
	Shortcuts map[string]ShortcutConfigItem `json:"shortcuts"`
	Presets   map[string]string             `json:"presets"`
	ActivePreset string                     `json:"active_preset"`
}

// ShortcutConfigItem represents a single shortcut configuration
type ShortcutConfigItem struct {
	Key         string `json:"key"`
	Action      string `json:"action"`
	Description string `json:"description"`
	Context     string `json:"context"`
	Enabled     bool   `json:"enabled"`
}

// ShortcutPreset represents a preset configuration
type ShortcutPreset struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Shortcuts   map[string]ShortcutConfigItem `json:"shortcuts"`
}

// Default presets
var DefaultPresets = map[string]ShortcutPreset{
	"default": {
		Name:        "Default",
		Description: "Standard shortcuts for the process manager",
		Shortcuts:   getDefaultShortcuts(),
	},
	"vim": {
		Name:        "Vim-style",
		Description: "Vim-inspired shortcuts",
		Shortcuts:   getVimShortcuts(),
	},
	"emacs": {
		Name:        "Emacs-style",
		Description: "Emacs-inspired shortcuts",
		Shortcuts:   getEmacsShortcuts(),
	},
}

// LoadConfig loads shortcut configuration from file
func LoadConfig(configPath string) (*ShortcutConfig, error) {
	// Create config directory if it doesn't exist
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}
	
	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Create default config
		config := &ShortcutConfig{
			Shortcuts:    getDefaultShortcuts(),
			Presets:      make(map[string]string),
			ActivePreset: "default",
		}
		
		// Save default config
		if err := SaveConfig(config, configPath); err != nil {
			return nil, fmt.Errorf("failed to save default config: %w", err)
		}
		
		return config, nil
	}
	
	// Load existing config
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}
	
	var config ShortcutConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}
	
	return &config, nil
}

// SaveConfig saves shortcut configuration to file
func SaveConfig(config *ShortcutConfig, configPath string) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	
	return nil
}

// ApplyConfig applies a configuration to the shortcut manager
func (m *ShortcutManager) ApplyConfig(config *ShortcutConfig) error {
	// Clear existing shortcuts
	m.registry = NewShortcutRegistry()
	
	// Apply shortcuts from config
	for _, item := range config.Shortcuts {
		shortcut := Shortcut{
			Key:         ParseKey(item.Key),
			Action:      item.Action,
			Description: item.Description,
			Context:     parseContext(item.Context),
			Enabled:     item.Enabled,
			Handler:     m.getHandlerForAction(item.Action),
		}
		m.RegisterShortcut(shortcut)
	}
	
	return nil
}

// ExportConfig exports current shortcuts to configuration
func (m *ShortcutManager) ExportConfig() *ShortcutConfig {
	config := &ShortcutConfig{
		Shortcuts:    make(map[string]ShortcutConfigItem),
		Presets:      make(map[string]string),
		ActivePreset: "custom",
	}
	
	// Export all shortcuts
	for context, shortcuts := range m.registry.shortcuts {
		for _, shortcut := range shortcuts {
			key := fmt.Sprintf("%s_%s", context.String(), shortcut.Action)
			config.Shortcuts[key] = ShortcutConfigItem{
				Key:         shortcut.Key.String(),
				Action:      shortcut.Action,
				Description: shortcut.Description,
				Context:     context.String(),
				Enabled:     shortcut.Enabled,
			}
		}
	}
	
	return config
}

// parseContext parses a context string
func parseContext(contextStr string) Context {
	switch contextStr {
	case "Global":
		return ContextGlobal
	case "Processes":
		return ContextProcesses
	case "Details":
		return ContextDetails
	case "Statistics":
		return ContextStats
	case "Settings":
		return ContextSettings
	case "Help":
		return ContextHelp
	case "Filter":
		return ContextFilter
	case "Search":
		return ContextSearch
	default:
		return ContextGlobal
	}
}

// getHandlerForAction returns a handler for a specific action
func (m *ShortcutManager) getHandlerForAction(action string) func() tea.Cmd {
	// This would be implemented with actual handlers
	// For now, return a placeholder
	return func() tea.Cmd {
		return tea.Printf("Action: %s", action)
	}
}

// getDefaultShortcuts returns default shortcut configuration
func getDefaultShortcuts() map[string]ShortcutConfigItem {
	return map[string]ShortcutConfigItem{
		"global_quit": {
			Key:         "ctrl+q",
			Action:      "quit",
			Description: "Quit application",
			Context:     "Global",
			Enabled:     true,
		},
		"global_help": {
			Key:         "ctrl+h",
			Action:      "help",
			Description: "Show help",
			Context:     "Global",
			Enabled:     true,
		},
		"global_refresh": {
			Key:         "ctrl+r",
			Action:      "refresh",
			Description: "Refresh current view",
			Context:     "Global",
			Enabled:     true,
		},
		"nav_processes": {
			Key:         "ctrl+p",
			Action:      "view_processes",
			Description: "Switch to Processes view",
			Context:     "Global",
			Enabled:     true,
		},
		"nav_details": {
			Key:         "ctrl+d",
			Action:      "view_details",
			Description: "Switch to Details view",
			Context:     "Global",
			Enabled:     true,
		},
		"nav_stats": {
			Key:         "ctrl+t",
			Action:      "view_stats",
			Description: "Switch to Statistics view",
			Context:     "Global",
			Enabled:     true,
		},
		"nav_settings": {
			Key:         "ctrl+,",
			Action:      "view_settings",
			Description: "Switch to Settings view",
			Context:     "Global",
			Enabled:     true,
		},
		"process_kill": {
			Key:         "ctrl+k",
			Action:      "kill_process",
			Description: "Kill selected process",
			Context:     "Processes",
			Enabled:     true,
		},
		"process_export": {
			Key:         "ctrl+e",
			Action:      "export",
			Description: "Export process data",
			Context:     "Processes",
			Enabled:     true,
		},
		"filter_search": {
			Key:         "ctrl+f",
			Action:      "search",
			Description: "Open search dialog",
			Context:     "Processes",
			Enabled:     true,
		},
		"sort_cpu": {
			Key:         "ctrl+o",
			Action:      "sort_cpu",
			Description: "Sort by CPU usage",
			Context:     "Processes",
			Enabled:     true,
		},
		"sort_memory": {
			Key:         "ctrl+m",
			Action:      "sort_memory",
			Description: "Sort by memory usage",
			Context:     "Processes",
			Enabled:     true,
		},
	}
}

// getVimShortcuts returns vim-style shortcut configuration
func getVimShortcuts() map[string]ShortcutConfigItem {
	shortcuts := getDefaultShortcuts()
	
	// Override with vim-style shortcuts
	shortcuts["global_quit"] = ShortcutConfigItem{
		Key:         ":q",
		Action:      "quit",
		Description: "Quit application",
		Context:     "Global",
		Enabled:     true,
	}
	shortcuts["process_kill"] = ShortcutConfigItem{
		Key:         "dd",
		Action:      "kill_process",
		Description: "Kill selected process",
		Context:     "Processes",
		Enabled:     true,
	}
	shortcuts["filter_search"] = ShortcutConfigItem{
		Key:         "/",
		Action:      "search",
		Description: "Open search dialog",
		Context:     "Processes",
		Enabled:     true,
	}
	
	return shortcuts
}

// getEmacsShortcuts returns emacs-style shortcut configuration
func getEmacsShortcuts() map[string]ShortcutConfigItem {
	shortcuts := getDefaultShortcuts()
	
	// Override with emacs-style shortcuts
	shortcuts["global_quit"] = ShortcutConfigItem{
		Key:         "ctrl+x ctrl+c",
		Action:      "quit",
		Description: "Quit application",
		Context:     "Global",
		Enabled:     true,
	}
	shortcuts["global_help"] = ShortcutConfigItem{
		Key:         "ctrl+h",
		Action:      "help",
		Description: "Show help",
		Context:     "Global",
		Enabled:     true,
	}
	shortcuts["filter_search"] = ShortcutConfigItem{
		Key:         "ctrl+s",
		Action:      "search",
		Description: "Open search dialog",
		Context:     "Processes",
		Enabled:     true,
	}
	
	return shortcuts
}

