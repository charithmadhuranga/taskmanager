package shortcuts

import (
	"fmt"
	"sort"
	"strings"
)

// HelpGenerator generates help text for shortcuts
type HelpGenerator struct {
	manager *ShortcutManager
}

// NewHelpGenerator creates a new help generator
func NewHelpGenerator(manager *ShortcutManager) *HelpGenerator {
	return &HelpGenerator{
		manager: manager,
	}
}

// GenerateHelp generates comprehensive help text
func (h *HelpGenerator) GenerateHelp() string {
	var help strings.Builder
	
	help.WriteString("Terminal Process Manager - Keyboard Shortcuts\n")
	help.WriteString("============================================\n\n")
	
	// Global shortcuts
	help.WriteString(h.generateContextHelp(ContextGlobal))
	
	// Context-specific shortcuts
	contexts := []Context{
		ContextProcesses,
		ContextDetails,
		ContextStats,
		ContextSettings,
		ContextHelp,
		ContextFilter,
		ContextSearch,
	}
	
	for _, context := range contexts {
		shortcuts := h.manager.GetShortcutsForContext(context)
		if len(shortcuts) > 0 {
			help.WriteString(h.generateContextHelp(context))
		}
	}
	
	// Shortcut tips
	help.WriteString(h.generateTips())
	
	return help.String()
}

// generateContextHelp generates help for a specific context
func (h *HelpGenerator) generateContextHelp(context Context) string {
	shortcuts := h.manager.GetShortcutsForContext(context)
	
	if len(shortcuts) == 0 {
		return ""
	}
	
	var help strings.Builder
	
	// Context header
	help.WriteString(fmt.Sprintf("%s Shortcuts:\n", context.String()))
	help.WriteString(strings.Repeat("-", len(context.String())+11))
	help.WriteString("\n")
	
	// Group shortcuts by category
	categories := h.groupShortcutsByCategory(shortcuts)
	
	for category, categoryShortcuts := range categories {
		if len(categoryShortcuts) > 0 {
			help.WriteString(fmt.Sprintf("\n%s:\n", category))
			
			// Sort shortcuts within category
			sort.Slice(categoryShortcuts, func(i, j int) bool {
				return categoryShortcuts[i].Key.String() < categoryShortcuts[j].Key.String()
			})
			
			for _, shortcut := range categoryShortcuts {
				if shortcut.Enabled {
					help.WriteString(fmt.Sprintf("  %-20s - %s\n", 
						shortcut.Key.String(), 
						shortcut.Description))
				}
			}
		}
	}
	
	help.WriteString("\n")
	return help.String()
}

// groupShortcutsByCategory groups shortcuts by category
func (h *HelpGenerator) groupShortcutsByCategory(shortcuts []Shortcut) map[string][]Shortcut {
	categories := make(map[string][]Shortcut)
	
	for _, shortcut := range shortcuts {
		category := h.getCategoryForAction(shortcut.Action)
		categories[category] = append(categories[category], shortcut)
	}
	
	return categories
}

// getCategoryForAction returns the category for an action
func (h *HelpGenerator) getCategoryForAction(action string) string {
	switch {
	case strings.Contains(action, "view_") || strings.Contains(action, "nav_"):
		return "Navigation"
	case strings.Contains(action, "kill") || strings.Contains(action, "process"):
		return "Process Management"
	case strings.Contains(action, "sort") || strings.Contains(action, "filter"):
		return "Sorting & Filtering"
	case strings.Contains(action, "search") || strings.Contains(action, "find"):
		return "Search"
	case strings.Contains(action, "export") || strings.Contains(action, "backup"):
		return "Data Management"
	case strings.Contains(action, "help") || strings.Contains(action, "info"):
		return "Help & Information"
	case strings.Contains(action, "refresh") || strings.Contains(action, "reload"):
		return "Refresh"
	case action == "quit" || action == "cancel":
		return "Application Control"
	default:
		return "Other"
	}
}

// generateTips generates helpful tips
func (h *HelpGenerator) generateTips() string {
	var tips strings.Builder
	
	tips.WriteString("Tips:\n")
	tips.WriteString("-----\n")
	tips.WriteString("• Shortcuts are context-sensitive - different views have different shortcuts\n")
	tips.WriteString("• Use Ctrl+H or F1 to show help for the current view\n")
	tips.WriteString("• Use Esc to cancel current operation or go back\n")
	tips.WriteString("• Use Tab to cycle through focusable elements\n")
	tips.WriteString("• Use Ctrl+R to refresh the current view\n")
	tips.WriteString("• Use Ctrl+Q or 'q' to quit the application\n")
	tips.WriteString("• Shortcuts can be customized in the settings\n")
	tips.WriteString("• Use Ctrl+Shift+F for advanced filtering options\n")
	tips.WriteString("• Use Ctrl+E to export data from any view\n")
	tips.WriteString("• Use Ctrl+B to create backups\n")
	
	return tips.String()
}

// GenerateContextHelp generates help for a specific context
func (h *HelpGenerator) GenerateContextHelp(context Context) string {
	shortcuts := h.manager.GetShortcutsForContext(context)
	
	if len(shortcuts) == 0 {
		return fmt.Sprintf("No shortcuts available for %s context.", context.String())
	}
	
	var help strings.Builder
	
	help.WriteString(fmt.Sprintf("Shortcuts for %s View:\n", context.String()))
	help.WriteString(strings.Repeat("=", len(context.String())+20))
	help.WriteString("\n\n")
	
	// Group by category
	categories := h.groupShortcutsByCategory(shortcuts)
	
	for category, categoryShortcuts := range categories {
		if len(categoryShortcuts) > 0 {
			help.WriteString(fmt.Sprintf("%s:\n", category))
			help.WriteString(strings.Repeat("-", len(category)+1))
			help.WriteString("\n")
			
			for _, shortcut := range categoryShortcuts {
				if shortcut.Enabled {
					help.WriteString(fmt.Sprintf("  %-20s - %s\n", 
						shortcut.Key.String(), 
						shortcut.Description))
				}
			}
			help.WriteString("\n")
		}
	}
	
	return help.String()
}

// GenerateQuickReference generates a quick reference card
func (h *HelpGenerator) GenerateQuickReference() string {
	var ref strings.Builder
	
	ref.WriteString("Quick Reference Card\n")
	ref.WriteString("===================\n\n")
	
	// Most commonly used shortcuts
	commonShortcuts := []string{
		"ctrl+q", "ctrl+h", "ctrl+r", "esc",
		"ctrl+p", "ctrl+d", "ctrl+t", "ctrl+,",
		"ctrl+k", "ctrl+e", "ctrl+f", "ctrl+o",
	}
	
	ref.WriteString("Most Common Shortcuts:\n")
	ref.WriteString("---------------------\n")
	
	for _, keyStr := range commonShortcuts {
		key := ParseKey(keyStr)
		help := h.manager.GetShortcutHelp(key)
		ref.WriteString(fmt.Sprintf("%s\n", help))
	}
	
	return ref.String()
}

// GenerateShortcutList generates a simple list of all shortcuts
func (h *HelpGenerator) GenerateShortcutList() []string {
	var shortcuts []string
	
	contexts := []Context{
		ContextGlobal,
		ContextProcesses,
		ContextDetails,
		ContextStats,
		ContextSettings,
		ContextHelp,
		ContextFilter,
		ContextSearch,
	}
	
	for _, context := range contexts {
		contextShortcuts := h.manager.GetShortcutsForContext(context)
		for _, shortcut := range contextShortcuts {
			if shortcut.Enabled {
				shortcuts = append(shortcuts, fmt.Sprintf("%s: %s (%s)", 
					shortcut.Key.String(), 
					shortcut.Description,
					context.String()))
			}
		}
	}
	
	// Sort shortcuts
	sort.Strings(shortcuts)
	
	return shortcuts
}

// GenerateConflictsHelp generates help for conflicting shortcuts
func (h *HelpGenerator) GenerateConflictsHelp() string {
	conflicts := h.manager.GetConflicts()
	
	if len(conflicts) == 0 {
		return "No conflicting shortcuts found."
	}
	
	var help strings.Builder
	
	help.WriteString("Conflicting Shortcuts:\n")
	help.WriteString("=====================\n\n")
	
	for key, shortcutList := range conflicts {
		help.WriteString(fmt.Sprintf("Key: %s\n", key.String()))
		help.WriteString("Conflicts:\n")
		
		for _, shortcut := range shortcutList {
			help.WriteString(fmt.Sprintf("  - %s (%s): %s\n", 
				shortcut.Action,
				shortcut.Context.String(),
				shortcut.Description))
		}
		help.WriteString("\n")
	}
	
	return help.String()
}

