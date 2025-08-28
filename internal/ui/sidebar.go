package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/romaintb/mel/internal/config"
)

// Sidebar represents the left sidebar with account/folder tree
type Sidebar struct {
	config  *config.Config
	width   int
	height  int
	focused bool
}

// NewSidebar creates a new sidebar instance
func NewSidebar(cfg *config.Config) (*Sidebar, error) {
	return &Sidebar{
		config:  cfg,
		width:   0, // Will be set by Resize
		height:  0,
		focused: false,
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

	var result string
	result += "📧 Accounts\n"
	result += "─────────\n"
	result += "├── 📥 Inbox (3)\n"
	result += "├── 📤 Sent\n"
	result += "├── 📁 Drafts\n"
	result += "├── 🗑️  Trash\n"
	result += "└── ⭐ Starred\n"

	// Add more content if we have more height available
	if s.height > 10 {
		result += "\n📊 Statistics\n"
		result += "────────────\n"
		result += "├── Total: 1,234\n"
		result += "├── Unread: 42\n"
		result += "└── Starred: 15\n"
	}

	if s.height > 20 {
		result += "\n🔧 Quick Actions\n"
		result += "───────────────\n"
		result += "├── 📝 Compose\n"
		result += "├── 🔍 Search\n"
		result += "└── ⚙️  Settings\n"
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
