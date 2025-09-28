package shortcuts

import (
	"fmt"
	"sort"
	tea "github.com/charmbracelet/bubbletea"
)

// ShortcutManager manages shortcuts and handles key events
type ShortcutManager struct {
	registry *ShortcutRegistry
	context  Context
}

// NewShortcutManager creates a new shortcut manager
func NewShortcutManager() *ShortcutManager {
	registry := NewShortcutRegistry()
	manager := &ShortcutManager{
		registry: registry,
		context:  ContextGlobal,
	}
	
	// Register default shortcuts
	manager.registerDefaultShortcuts()
	
	return manager
}

// SetContext sets the current context
func (m *ShortcutManager) SetContext(context Context) {
	m.context = context
}

// GetContext returns the current context
func (m *ShortcutManager) GetContext() Context {
	return m.context
}

// HandleKey handles a key event and returns the appropriate command
func (m *ShortcutManager) HandleKey(msg tea.KeyMsg) tea.Cmd {
	// First try to find a shortcut in the current context
	shortcut := m.registry.GetShortcut(ShortcutKey{
		Key:      msg.String(),
		Modifier: m.getModifierFromMsg(msg),
	}, m.context)
	
	if shortcut != nil && shortcut.Enabled {
		return shortcut.Handler()
	}
	
	// If not found in current context, try global context
	if m.context != ContextGlobal {
		shortcut = m.registry.GetShortcut(ShortcutKey{
			Key:      msg.String(),
			Modifier: m.getModifierFromMsg(msg),
		}, ContextGlobal)
		
		if shortcut != nil && shortcut.Enabled {
			return shortcut.Handler()
		}
	}
	
	return nil
}

// getModifierFromMsg extracts modifier from tea.KeyMsg
func (m *ShortcutManager) getModifierFromMsg(msg tea.KeyMsg) Modifier {
	if msg.Ctrl && msg.Alt && msg.Shift {
		return ModCtrlAltShift
	} else if msg.Ctrl && msg.Alt {
		return ModCtrlAlt
	} else if msg.Ctrl && msg.Shift {
		return ModCtrlShift
	} else if msg.Alt && msg.Shift {
		return ModAltShift
	} else if msg.Ctrl {
		return ModCtrl
	} else if msg.Alt {
		return ModAlt
	} else if msg.Shift {
		return ModShift
	}
	return ModNone
}

// GetShortcutsForContext returns all shortcuts for a specific context
func (m *ShortcutManager) GetShortcutsForContext(context Context) []Shortcut {
	shortcuts := m.registry.GetShortcuts(context)
	
	// Sort shortcuts by key for consistent display
	sort.Slice(shortcuts, func(i, j int) bool {
		return shortcuts[i].Key.String() < shortcuts[j].Key.String()
	})
	
	return shortcuts
}

// GetShortcutsForCurrentContext returns shortcuts for the current context
func (m *ShortcutManager) GetShortcutsForCurrentContext() []Shortcut {
	return m.GetShortcutsForContext(m.context)
}

// RegisterShortcut registers a new shortcut
func (m *ShortcutManager) RegisterShortcut(shortcut Shortcut) {
	m.registry.RegisterShortcut(shortcut)
}

// RegisterShortcuts registers multiple shortcuts
func (m *ShortcutManager) RegisterShortcuts(shortcuts []Shortcut) {
	for _, shortcut := range shortcuts {
		m.registry.RegisterShortcut(shortcut)
	}
}

// EnableShortcut enables a shortcut
func (m *ShortcutManager) EnableShortcut(key ShortcutKey, context Context) {
	shortcuts := m.registry.GetShortcuts(context)
	for i, shortcut := range shortcuts {
		if shortcut.Key == key {
			shortcuts[i].Enabled = true
			break
		}
	}
}

// DisableShortcut disables a shortcut
func (m *ShortcutManager) DisableShortcut(key ShortcutKey, context Context) {
	shortcuts := m.registry.GetShortcuts(context)
	for i, shortcut := range shortcuts {
		if shortcut.Key == key {
			shortcuts[i].Enabled = false
			break
		}
	}
}

// GetConflicts returns conflicting shortcuts
func (m *ShortcutManager) GetConflicts() map[ShortcutKey][]Shortcut {
	conflicts := make(map[ShortcutKey][]Shortcut)
	
	for key, shortcuts := range m.registry.conflicts {
		if len(shortcuts) > 1 {
			conflicts[key] = shortcuts
		}
	}
	
	return conflicts
}

// registerDefaultShortcuts registers all default shortcuts
func (m *ShortcutManager) registerDefaultShortcuts() {
	// Global shortcuts
	globalShortcuts := []Shortcut{
		{
			Key:         ParseKey("ctrl+q"),
			Action:      "quit",
			Description: "Quit application",
			Context:     ContextGlobal,
			Handler:     func() tea.Cmd { return tea.Quit },
			Enabled:     true,
		},
		{
			Key:         ParseKey("q"),
			Action:      "quit",
			Description: "Quit application",
			Context:     ContextGlobal,
			Handler:     func() tea.Cmd { return tea.Quit },
			Enabled:     true,
		},
		{
			Key:         ParseKey("ctrl+h"),
			Action:      "help",
			Description: "Show help",
			Context:     ContextGlobal,
			Handler:     func() tea.Cmd { return tea.Printf("Help requested") },
			Enabled:     true,
		},
		{
			Key:         ParseKey("f1"),
			Action:      "help",
			Description: "Show help",
			Context:     ContextGlobal,
			Handler:     func() tea.Cmd { return tea.Printf("Help requested") },
			Enabled:     true,
		},
		{
			Key:         ParseKey("esc"),
			Action:      "cancel",
			Description: "Cancel current operation",
			Context:     ContextGlobal,
			Handler:     func() tea.Cmd { return tea.Printf("Operation cancelled") },
			Enabled:     true,
		},
		{
			Key:         ParseKey("ctrl+r"),
			Action:      "refresh",
			Description: "Refresh current view",
			Context:     ContextGlobal,
			Handler:     func() tea.Cmd { return tea.Printf("Refreshing...") },
			Enabled:     true,
		},
	}
	
	// Navigation shortcuts
	navShortcuts := []Shortcut{
		{
			Key:         ParseKey("ctrl+p"),
			Action:      "view_processes",
			Description: "Switch to Processes view",
			Context:     ContextGlobal,
			Handler:     func() tea.Cmd { return tea.Printf("Switching to Processes view") },
			Enabled:     true,
		},
		{
			Key:         ParseKey("ctrl+d"),
			Action:      "view_details",
			Description: "Switch to Details view",
			Context:     ContextGlobal,
			Handler:     func() tea.Cmd { return tea.Printf("Switching to Details view") },
			Enabled:     true,
		},
		{
			Key:         ParseKey("ctrl+t"),
			Action:      "view_stats",
			Description: "Switch to Statistics view",
			Context:     ContextGlobal,
			Handler:     func() tea.Cmd { return tea.Printf("Switching to Statistics view") },
			Enabled:     true,
		},
		{
			Key:         ParseKey("ctrl+,"),
			Action:      "view_settings",
			Description: "Switch to Settings view",
			Context:     ContextGlobal,
			Handler:     func() tea.Cmd { return tea.Printf("Switching to Settings view") },
			Enabled:     true,
		},
	}
	
	// Process management shortcuts
	processShortcuts := []Shortcut{
		{
			Key:         ParseKey("ctrl+k"),
			Action:      "kill_process",
			Description: "Kill selected process",
			Context:     ContextProcesses,
			Handler:     func() tea.Cmd { return tea.Printf("Killing process...") },
			Enabled:     true,
		},
		{
			Key:         ParseKey("ctrl+shift+k"),
			Action:      "force_kill_process",
			Description: "Force kill selected process",
			Context:     ContextProcesses,
			Handler:     func() tea.Cmd { return tea.Printf("Force killing process...") },
			Enabled:     true,
		},
		{
			Key:         ParseKey("ctrl+e"),
			Action:      "export",
			Description: "Export process data",
			Context:     ContextProcesses,
			Handler:     func() tea.Cmd { return tea.Printf("Exporting data...") },
			Enabled:     true,
		},
		{
			Key:         ParseKey("ctrl+b"),
			Action:      "backup",
			Description: "Create backup",
			Context:     ContextProcesses,
			Handler:     func() tea.Cmd { return tea.Printf("Creating backup...") },
			Enabled:     true,
		},
	}
	
	// Filtering shortcuts
	filterShortcuts := []Shortcut{
		{
			Key:         ParseKey("ctrl+f"),
			Action:      "search",
			Description: "Open search dialog",
			Context:     ContextProcesses,
			Handler:     func() tea.Cmd { return tea.Printf("Opening search...") },
			Enabled:     true,
		},
		{
			Key:         ParseKey("ctrl+shift+f"),
			Action:      "advanced_filter",
			Description: "Open advanced filter dialog",
			Context:     ContextProcesses,
			Handler:     func() tea.Cmd { return tea.Printf("Opening advanced filter...") },
			Enabled:     true,
		},
		{
			Key:         ParseKey("ctrl+l"),
			Action:      "clear_filters",
			Description: "Clear all filters",
			Context:     ContextProcesses,
			Handler:     func() tea.Cmd { return tea.Printf("Clearing filters...") },
			Enabled:     true,
		},
	}
	
	// Sorting shortcuts
	sortShortcuts := []Shortcut{
		{
			Key:         ParseKey("ctrl+o"),
			Action:      "sort_cpu",
			Description: "Sort by CPU usage",
			Context:     ContextProcesses,
			Handler:     func() tea.Cmd { return tea.Printf("Sorting by CPU...") },
			Enabled:     true,
		},
		{
			Key:         ParseKey("ctrl+m"),
			Action:      "sort_memory",
			Description: "Sort by memory usage",
			Context:     ContextProcesses,
			Handler:     func() tea.Cmd { return tea.Printf("Sorting by memory...") },
			Enabled:     true,
		},
		{
			Key:         ParseKey("ctrl+n"),
			Action:      "sort_name",
			Description: "Sort by name",
			Context:     ContextProcesses,
			Handler:     func() tea.Cmd { return tea.Printf("Sorting by name...") },
			Enabled:     true,
		},
		{
			Key:         ParseKey("ctrl+s"),
			Action:      "sort_status",
			Description: "Sort by status",
			Context:     ContextProcesses,
			Handler:     func() tea.Cmd { return tea.Printf("Sorting by status...") },
			Enabled:     true,
		},
	}
	
	// Register all shortcuts
	m.RegisterShortcuts(globalShortcuts)
	m.RegisterShortcuts(navShortcuts)
	m.RegisterShortcuts(processShortcuts)
	m.RegisterShortcuts(filterShortcuts)
	m.RegisterShortcuts(sortShortcuts)
}

// GetHelpText returns formatted help text for shortcuts
func (m *ShortcutManager) GetHelpText(context Context) string {
	shortcuts := m.GetShortcutsForContext(context)
	
	if len(shortcuts) == 0 {
		return "No shortcuts available for this context."
	}
	
	help := fmt.Sprintf("Shortcuts for %s context:\n\n", context.String())
	
	for _, shortcut := range shortcuts {
		if shortcut.Enabled {
			help += fmt.Sprintf("%-20s - %s\n", shortcut.Key.String(), shortcut.Description)
		}
	}
	
	return help
}

// GetShortcutHelp returns help for a specific shortcut
func (m *ShortcutManager) GetShortcutHelp(key ShortcutKey) string {
	// Try current context first
	shortcut := m.registry.GetShortcut(key, m.context)
	if shortcut != nil {
		return fmt.Sprintf("%s: %s", shortcut.Key.String(), shortcut.Description)
	}
	
	// Try global context
	shortcut = m.registry.GetShortcut(key, ContextGlobal)
	if shortcut != nil {
		return fmt.Sprintf("%s: %s", shortcut.Key.String(), shortcut.Description)
	}
	
	return fmt.Sprintf("No shortcut found for %s", key.String())
}

