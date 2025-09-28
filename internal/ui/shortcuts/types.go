package shortcuts

import (
	"strings"
	tea "github.com/charmbracelet/bubbletea"
)

// ShortcutKey represents a keyboard shortcut
type ShortcutKey struct {
	Key      string
	Modifier Modifier
}

// Modifier represents keyboard modifiers
type Modifier int

const (
	ModNone Modifier = iota
	ModCtrl
	ModAlt
	ModShift
	ModCtrlShift
	ModCtrlAlt
	ModAltShift
	ModCtrlAltShift
)

// String returns the string representation of a modifier
func (m Modifier) String() string {
	switch m {
	case ModNone:
		return ""
	case ModCtrl:
		return "Ctrl"
	case ModAlt:
		return "Alt"
	case ModShift:
		return "Shift"
	case ModCtrlShift:
		return "Ctrl+Shift"
	case ModCtrlAlt:
		return "Ctrl+Alt"
	case ModAltShift:
		return "Alt+Shift"
	case ModCtrlAltShift:
		return "Ctrl+Alt+Shift"
	default:
		return ""
	}
}

// Shortcut represents a keyboard shortcut with its action
type Shortcut struct {
	Key        ShortcutKey
	Action     string
	Description string
	Context    Context
	Handler    func() tea.Cmd
	Enabled    bool
}

// Context represents the context where a shortcut is active
type Context int

const (
	ContextGlobal Context = iota
	ContextProcesses
	ContextDetails
	ContextStats
	ContextSettings
	ContextHelp
	ContextFilter
	ContextSearch
)

// String returns the string representation of a context
func (c Context) String() string {
	switch c {
	case ContextGlobal:
		return "Global"
	case ContextProcesses:
		return "Processes"
	case ContextDetails:
		return "Details"
	case ContextStats:
		return "Statistics"
	case ContextSettings:
		return "Settings"
	case ContextHelp:
		return "Help"
	case ContextFilter:
		return "Filter"
	case ContextSearch:
		return "Search"
	default:
		return "Unknown"
	}
}

// ShortcutRegistry manages all shortcuts
type ShortcutRegistry struct {
	shortcuts map[Context][]Shortcut
	conflicts map[ShortcutKey][]Shortcut
}

// NewShortcutRegistry creates a new shortcut registry
func NewShortcutRegistry() *ShortcutRegistry {
	return &ShortcutRegistry{
		shortcuts: make(map[Context][]Shortcut),
		conflicts: make(map[ShortcutKey][]Shortcut),
	}
}

// RegisterShortcut registers a new shortcut
func (r *ShortcutRegistry) RegisterShortcut(shortcut Shortcut) {
	// Add to context
	r.shortcuts[shortcut.Context] = append(r.shortcuts[shortcut.Context], shortcut)
	
	// Check for conflicts
	key := shortcut.Key
	if existing, exists := r.conflicts[key]; exists {
		r.conflicts[key] = append(existing, shortcut)
	} else {
		r.conflicts[key] = []Shortcut{shortcut}
	}
}

// GetShortcuts returns shortcuts for a specific context
func (r *ShortcutRegistry) GetShortcuts(context Context) []Shortcut {
	return r.shortcuts[context]
}

// GetShortcut returns a specific shortcut by key and context
func (r *ShortcutRegistry) GetShortcut(key ShortcutKey, context Context) *Shortcut {
	shortcuts := r.shortcuts[context]
	for _, shortcut := range shortcuts {
		if shortcut.Key == key && shortcut.Enabled {
			return &shortcut
		}
	}
	return nil
}

// GetConflicts returns conflicting shortcuts for a key
func (r *ShortcutRegistry) GetConflicts(key ShortcutKey) []Shortcut {
	return r.conflicts[key]
}

// ParseKey parses a key string into a ShortcutKey
func ParseKey(keyStr string) ShortcutKey {
	parts := strings.Split(keyStr, "+")
	
	var modifier Modifier
	var key string
	
	if len(parts) == 1 {
		key = parts[0]
		modifier = ModNone
	} else {
		key = parts[len(parts)-1]
		modifiers := parts[:len(parts)-1]
		
		hasCtrl := false
		hasAlt := false
		hasShift := false
		
		for _, mod := range modifiers {
			switch strings.ToLower(mod) {
			case "ctrl":
				hasCtrl = true
			case "alt":
				hasAlt = true
			case "shift":
				hasShift = true
			}
		}
		
		if hasCtrl && hasAlt && hasShift {
			modifier = ModCtrlAltShift
		} else if hasCtrl && hasAlt {
			modifier = ModCtrlAlt
		} else if hasCtrl && hasShift {
			modifier = ModCtrlShift
		} else if hasAlt && hasShift {
			modifier = ModAltShift
		} else if hasCtrl {
			modifier = ModCtrl
		} else if hasAlt {
			modifier = ModAlt
		} else if hasShift {
			modifier = ModShift
		}
	}
	
	return ShortcutKey{
		Key:      key,
		Modifier: modifier,
	}
}

// String returns the string representation of a ShortcutKey
func (k ShortcutKey) String() string {
	if k.Modifier == ModNone {
		return k.Key
	}
	return k.Modifier.String() + "+" + k.Key
}

// Matches checks if a tea.KeyMsg matches this shortcut key
func (k ShortcutKey) Matches(msg tea.KeyMsg) bool {
	// Check if the key matches
	if strings.ToLower(msg.String()) != strings.ToLower(k.Key) {
		return false
	}
	
	// Check modifiers
	switch k.Modifier {
	case ModNone:
		return !msg.Alt && !msg.Ctrl && !msg.Shift
	case ModCtrl:
		return msg.Ctrl && !msg.Alt && !msg.Shift
	case ModAlt:
		return msg.Alt && !msg.Ctrl && !msg.Shift
	case ModShift:
		return msg.Shift && !msg.Ctrl && !msg.Alt
	case ModCtrlShift:
		return msg.Ctrl && msg.Shift && !msg.Alt
	case ModCtrlAlt:
		return msg.Ctrl && msg.Alt && !msg.Shift
	case ModAltShift:
		return msg.Alt && msg.Shift && !msg.Ctrl
	case ModCtrlAltShift:
		return msg.Ctrl && msg.Alt && msg.Shift
	}
	
	return false
}

