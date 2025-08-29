package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/romaintb/mel/internal/config"
	"github.com/romaintb/mel/internal/email"
	"github.com/romaintb/mel/internal/icons"
)

// Sidebar represents the left sidebar with account/folder tree
type Sidebar struct {
	config        *config.Config
	emailManager  *email.Manager
	iconService   *icons.Service
	width         int
	height        int
	focused       bool
	collapsed     bool
	selectedIndex int // Index of selected item
}

// NewSidebar creates a new sidebar instance
func NewSidebar(cfg *config.Config, emailManager *email.Manager, iconService *icons.Service) (*Sidebar, error) {
	return &Sidebar{
		config:        cfg,
		emailManager:  emailManager,
		iconService:   iconService,
		width:         0, // Will be set by Resize
		height:        0,
		focused:       false,
		collapsed:     false,
		selectedIndex: 0, // Start with first item selected
	}, nil
}

// Init initializes the sidebar
func (s *Sidebar) Init() tea.Cmd {
	return nil
}

// Update handles sidebar updates
func (s *Sidebar) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return s.handleKeyPress(msg)
	}
	return s, nil
}

// View renders the sidebar
func (s *Sidebar) View() string {
	if s.width == 0 {
		return ""
	}

	// If collapsed, show minimal sidebar
	if s.collapsed {
		return s.iconService.Get("email")
	}

	var result string
	result += s.iconService.Get("email") + " Accounts\n"
	result += "─────────\n"

	// Accounts section (items 0-4)
	if s.selectedIndex == 0 {
		result += "├── " + s.iconService.Get("selected") + " " + s.iconService.Get("inbox") + " Inbox (3)\n"
	} else {
		result += "├── " + s.iconService.Get("inbox") + " Inbox (3)\n"
	}

	if s.selectedIndex == 1 {
		result += "├── " + s.iconService.Get("selected") + " " + s.iconService.Get("sent") + " Sent\n"
	} else {
		result += "├── " + s.iconService.Get("sent") + " Sent\n"
	}

	if s.selectedIndex == 2 {
		result += "├── " + s.iconService.Get("selected") + " " + s.iconService.Get("drafts") + " Drafts\n"
	} else {
		result += "├── " + s.iconService.Get("drafts") + " Drafts\n"
	}

	if s.selectedIndex == 3 {
		result += "├── " + s.iconService.Get("selected") + " " + s.iconService.Get("trash") + " Trash\n"
	} else {
		result += "├── " + s.iconService.Get("trash") + " Trash\n"
	}

	if s.selectedIndex == 4 {
		result += "└── " + s.iconService.Get("selected") + " " + s.iconService.Get("starred") + " Starred\n"
	} else {
		result += "└── " + s.iconService.Get("starred") + " Starred\n"
	}

	// Add more content if we have more height available
	if s.height > 10 {
		result += "\n" + s.iconService.Get("total") + " Statistics\n"
		result += "────────────\n"

		// Statistics section (items 5-7)
		if s.selectedIndex == 5 {
			result += "├── " + s.iconService.Get("selected") + " Total: 1,234\n"
		} else {
			result += "├── Total: 1,234\n"
		}

		if s.selectedIndex == 6 {
			result += "├── " + s.iconService.Get("selected") + " Unread: 42\n"
		} else {
			result += "├── Unread: 42\n"
		}

		if s.selectedIndex == 7 {
			result += "└── " + s.iconService.Get("selected") + " Starred: 15\n"
		} else {
			result += "└── Starred: 15\n"
		}
	}

	if s.height > 20 {
		result += "\n" + s.iconService.Get("settings") + " Quick Actions\n"
		result += "───────────────\n"

		// Quick Actions section (items 8-10)
		if s.selectedIndex == 8 {
			result += "├── " + s.iconService.Get("selected") + " " + s.iconService.Get("compose") + " Compose\n"
		} else {
			result += "├── " + s.iconService.Get("compose") + " Compose\n"
		}

		if s.selectedIndex == 9 {
			result += "├── " + s.iconService.Get("selected") + " " + s.iconService.Get("search") + " Search\n"
		} else {
			result += "├── " + s.iconService.Get("search") + " Search\n"
		}

		if s.selectedIndex == 10 {
			result += "└── " + s.iconService.Get("selected") + " " + s.iconService.Get("settings") + " Settings\n"
		} else {
			result += "└── " + s.iconService.Get("settings") + " Settings\n"
		}
	}

	// Let lipgloss handle all the layout and padding
	return result
}

// Focus focuses the sidebar
func (s *Sidebar) Focus() tea.Cmd {
	s.focused = true
	return nil
}

// Blur removes focus from the sidebar
func (s *Sidebar) Blur() tea.Cmd {
	s.focused = false
	return nil
}

// Resize resizes the sidebar
func (s *Sidebar) Resize(width, height int) tea.Cmd {
	s.width = width
	s.height = height
	return nil
}

// Toggle toggles the sidebar collapsed state
func (s *Sidebar) Toggle() tea.Cmd {
	s.collapsed = !s.collapsed
	return nil
}

// Next moves selection to next item
func (s *Sidebar) Next() tea.Cmd {
	s.selectedIndex++
	// Wrap around to first item
	if s.selectedIndex >= s.getItemCount() {
		s.selectedIndex = 0
	}
	return nil
}

// Prev moves selection to previous item
func (s *Sidebar) Prev() tea.Cmd {
	s.selectedIndex--
	// Wrap around to last item
	if s.selectedIndex < 0 {
		s.selectedIndex = s.getItemCount() - 1
	}
	return nil
}

// GoToTop moves selection to first item
func (s *Sidebar) GoToTop() tea.Cmd {
	s.selectedIndex = 0
	return nil
}

// GoToBottom moves selection to last item
func (s *Sidebar) GoToBottom() tea.Cmd {
	s.selectedIndex = s.getItemCount() - 1
	return nil
}

// getItemCount returns the total number of selectable items
func (s *Sidebar) getItemCount() int {
	count := 5 // Accounts section: Inbox, Sent, Drafts, Trash, Starred
	if s.height > 10 {
		count += 3 // Statistics section: Total, Unread, Starred
	}
	if s.height > 20 {
		count += 3 // Quick Actions section: Compose, Search, Settings
	}
	return count
}

// handleKeyPress handles key presses in the sidebar
func (s *Sidebar) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if !s.focused {
		return s, nil
	}

	switch msg.String() {
	case "j":
		// Move down
	case "k":
		// Move up
	case "enter":
		// Select folder
	}

	return s, nil
}
