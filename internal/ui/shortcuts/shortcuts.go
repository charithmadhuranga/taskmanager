package shortcuts

import (
	"fmt"
	"os"
	"path/filepath"
	tea "github.com/charmbracelet/bubbletea"
)

// ShortcutSystem represents the complete shortcut system
type ShortcutSystem struct {
	manager      *ShortcutManager
	config       *ShortcutConfig
	configPath   string
	helpGenerator *HelpGenerator
}

// NewShortcutSystem creates a new shortcut system
func NewShortcutSystem() *ShortcutSystem {
	manager := NewShortcutManager()
	
	// Default config path
	homeDir, _ := os.UserHomeDir()
	configPath := filepath.Join(homeDir, ".tappmanager", "shortcuts.json")
	
	system := &ShortcutSystem{
		manager:      manager,
		configPath:   configPath,
		helpGenerator: NewHelpGenerator(manager),
	}
	
	// Load configuration
	config, err := LoadConfig(configPath)
	if err != nil {
		// Use default configuration
		config = &ShortcutConfig{
			Shortcuts:    getDefaultShortcuts(),
			Presets:      make(map[string]string),
			ActivePreset: "default",
		}
	}
	
	system.config = config
	system.manager.ApplyConfig(config)
	
	return system
}

// HandleKey handles a key event
func (s *ShortcutSystem) HandleKey(msg tea.KeyMsg) tea.Cmd {
	return s.manager.HandleKey(msg)
}

// SetContext sets the current context
func (s *ShortcutSystem) SetContext(context Context) {
	s.manager.SetContext(context)
}

// GetContext returns the current context
func (s *ShortcutSystem) GetContext() Context {
	return s.manager.GetContext()
}

// GetHelpText returns help text for the current context
func (s *ShortcutSystem) GetHelpText() string {
	return s.helpGenerator.GenerateContextHelp(s.manager.GetContext())
}

// GetFullHelpText returns comprehensive help text
func (s *ShortcutSystem) GetFullHelpText() string {
	return s.helpGenerator.GenerateHelp()
}

// GetQuickReference returns a quick reference card
func (s *ShortcutSystem) GetQuickReference() string {
	return s.helpGenerator.GenerateQuickReference()
}

// GetConflictsHelp returns help for conflicting shortcuts
func (s *ShortcutSystem) GetConflictsHelp() string {
	return s.helpGenerator.GenerateConflictsHelp()
}

// RegisterShortcut registers a new shortcut
func (s *ShortcutSystem) RegisterShortcut(shortcut Shortcut) {
	s.manager.RegisterShortcut(shortcut)
}

// RegisterShortcuts registers multiple shortcuts
func (s *ShortcutSystem) RegisterShortcuts(shortcuts []Shortcut) {
	s.manager.RegisterShortcuts(shortcuts)
}

// EnableShortcut enables a shortcut
func (s *ShortcutSystem) EnableShortcut(key ShortcutKey, context Context) {
	s.manager.EnableShortcut(key, context)
}

// DisableShortcut disables a shortcut
func (s *ShortcutSystem) DisableShortcut(key ShortcutKey, context Context) {
	s.manager.DisableShortcut(key, context)
}

// GetShortcutsForContext returns shortcuts for a specific context
func (s *ShortcutSystem) GetShortcutsForContext(context Context) []Shortcut {
	return s.manager.GetShortcutsForContext(context)
}

// GetShortcutsForCurrentContext returns shortcuts for the current context
func (s *ShortcutSystem) GetShortcutsForCurrentContext() []Shortcut {
	return s.manager.GetShortcutsForCurrentContext()
}

// SaveConfig saves the current configuration
func (s *ShortcutSystem) SaveConfig() error {
	config := s.manager.ExportConfig()
	return SaveConfig(config, s.configPath)
}

// LoadConfig loads configuration from file
func (s *ShortcutSystem) LoadConfig() error {
	config, err := LoadConfig(s.configPath)
	if err != nil {
		return err
	}
	
	s.config = config
	return s.manager.ApplyConfig(config)
}

// GetConfig returns the current configuration
func (s *ShortcutSystem) GetConfig() *ShortcutConfig {
	return s.config
}

// SetConfigPath sets the configuration file path
func (s *ShortcutSystem) SetConfigPath(path string) {
	s.configPath = path
}

// GetConfigPath returns the configuration file path
func (s *ShortcutSystem) GetConfigPath() string {
	return s.configPath
}

// GetConflicts returns conflicting shortcuts
func (s *ShortcutSystem) GetConflicts() map[ShortcutKey][]Shortcut {
	return s.manager.GetConflicts()
}

// ResolveConflict resolves a shortcut conflict by disabling one of the conflicting shortcuts
func (s *ShortcutSystem) ResolveConflict(key ShortcutKey, keepContext Context) {
	conflicts := s.manager.GetConflicts()
	if shortcutList, exists := conflicts[key]; exists {
		for _, shortcut := range shortcutList {
			if shortcut.Context != keepContext {
				s.DisableShortcut(shortcut.Key, shortcut.Context)
			}
		}
	}
}

// GetShortcutHelp returns help for a specific shortcut
func (s *ShortcutSystem) GetShortcutHelp(key ShortcutKey) string {
	return s.manager.GetShortcutHelp(key)
}

// GetShortcutList returns a list of all shortcuts
func (s *ShortcutSystem) GetShortcutList() []string {
	return s.helpGenerator.GenerateShortcutList()
}

// ValidateShortcut validates a shortcut key
func (s *ShortcutSystem) ValidateShortcut(key ShortcutKey) error {
	// Basic validation
	if key.Key == "" {
		return fmt.Errorf("shortcut key cannot be empty")
	}
	
	// Check for conflicts
	conflicts := s.GetConflicts()
	if shortcutList, exists := conflicts[key]; exists && len(shortcutList) > 1 {
		return fmt.Errorf("shortcut %s conflicts with other shortcuts", key.String())
	}
	
	return nil
}

// CreateShortcut creates a new shortcut with validation
func (s *ShortcutSystem) CreateShortcut(keyStr, action, description string, context Context, handler func() tea.Cmd) error {
	key := ParseKey(keyStr)
	
	// Validate the shortcut
	if err := s.ValidateShortcut(key); err != nil {
		return err
	}
	
	shortcut := Shortcut{
		Key:         key,
		Action:      action,
		Description: description,
		Context:     context,
		Handler:     handler,
		Enabled:     true,
	}
	
	s.RegisterShortcut(shortcut)
	return nil
}

// UpdateShortcut updates an existing shortcut
func (s *ShortcutSystem) UpdateShortcut(key ShortcutKey, context Context, updates Shortcut) error {
	// Find and update the shortcut
	shortcuts := s.manager.GetShortcutsForContext(context)
	for i, shortcut := range shortcuts {
		if shortcut.Key == key {
			// Update the shortcut
			if updates.Action != "" {
				shortcuts[i].Action = updates.Action
			}
			if updates.Description != "" {
				shortcuts[i].Description = updates.Description
			}
			if updates.Handler != nil {
				shortcuts[i].Handler = updates.Handler
			}
			shortcuts[i].Enabled = updates.Enabled
			return nil
		}
	}
	
	return fmt.Errorf("shortcut not found")
}

// DeleteShortcut deletes a shortcut
func (s *ShortcutSystem) DeleteShortcut(key ShortcutKey, context Context) error {
	// This would require modifying the registry to support deletion
	// For now, just disable it
	s.DisableShortcut(key, context)
	return nil
}

// GetPresets returns available presets
func (s *ShortcutSystem) GetPresets() map[string]ShortcutPreset {
	return DefaultPresets
}

// ApplyPreset applies a preset configuration
func (s *ShortcutSystem) ApplyPreset(presetName string) error {
	presets := s.GetPresets()
	preset, exists := presets[presetName]
	if !exists {
		return fmt.Errorf("preset %s not found", presetName)
	}
	
	// Convert preset to config
	config := &ShortcutConfig{
		Shortcuts:    preset.Shortcuts,
		Presets:      make(map[string]string),
		ActivePreset: presetName,
	}
	
	// Apply the preset
	s.config = config
	return s.manager.ApplyConfig(config)
}

// GetActivePreset returns the active preset name
func (s *ShortcutSystem) GetActivePreset() string {
	return s.config.ActivePreset
}

// SetActivePreset sets the active preset
func (s *ShortcutSystem) SetActivePreset(presetName string) error {
	return s.ApplyPreset(presetName)
}
