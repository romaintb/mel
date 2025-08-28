package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/romaintb/mel/internal/config"
)

// ThreadList represents the list of email threads
type ThreadList struct {
	config   *config.Config
	width    int
	height   int
	focused  bool
	selected int
	threads  []Thread
}

// Thread represents an email thread
type Thread struct {
	ID      string
	Subject string
	From    string
	Date    string
	Unread  bool
	Starred bool
}

// NewThreadList creates a new thread list instance
func NewThreadList(cfg *config.Config) (*ThreadList, error) {
	// Mock data for now
	threads := []Thread{
		{ID: "1", Subject: "Welcome to Mel", From: "team@mel.com", Date: "2024-01-15", Unread: true, Starred: false},
		{ID: "2", Subject: "Project Update", From: "manager@work.com", Date: "2024-01-14", Unread: false, Starred: true},
		{ID: "3", Subject: "Meeting Tomorrow", From: "colleague@work.com", Date: "2024-01-13", Unread: true, Starred: false},
	}

	return &ThreadList{
		config:   cfg,
		width:    0, // Will be set by Resize
		height:   0,
		focused:  false,
		selected: 0,
		threads:  threads,
	}, nil
}

// Init initializes the thread list
func (t *ThreadList) Init() tea.Cmd {
	return nil
}

// Update handles thread list updates
func (t *ThreadList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return t.handleKeyPress(msg)
	}
	return t, nil
}

// View renders the thread list
func (t *ThreadList) View() string {
	if t.width == 0 || len(t.threads) == 0 {
		return "No threads"
	}

	var result string
	result += "📧 Threads\n"
	result += "─────────\n"

	for i, thread := range t.threads {
		prefix := "  "
		if i == t.selected {
			prefix = "▶ "
		}

		unread := ""
		if thread.Unread {
			unread = "● "
		}

		starred := ""
		if thread.Starred {
			starred = "⭐ "
		}

		result += prefix + unread + starred + thread.Subject + "\n"
		result += "   " + thread.From + " • " + thread.Date + "\n"
		if i < len(t.threads)-1 {
			result += "\n"
		}
	}

	// Add more mock threads if we have more height available
	if t.height > 15 {
		result += "\n📧 More Threads\n"
		result += "──────────────\n"
		result += "  📋 Weekly Report\n"
		result += "   team@company.com • 2024-01-12\n"
		result += "\n  📅 Calendar Update\n"
		result += "   calendar@company.com • 2024-01-11\n"
		result += "\n  💬 Team Chat\n"
		result += "   chat@company.com • 2024-01-10\n"
	}

	if t.height > 25 {
		result += "\n📧 Archive\n"
		result += "─────────\n"
		result += "  📊 Q4 Results\n"
		result += "   finance@company.com • 2024-01-09\n"
		result += "\n  🎉 Holiday Party\n"
		result += "   hr@company.com • 2024-01-08\n"
	}

	// Let lipgloss handle all the layout and padding
	return result
}

// Focus focuses the thread list
func (t *ThreadList) Focus() tea.Cmd {
	t.focused = true
	return nil
}

// Blur removes focus from the thread list
func (t *ThreadList) Blur() tea.Cmd {
	t.focused = false
	return nil
}

// Resize resizes the thread list
func (t *ThreadList) Resize(width, height int) tea.Cmd {
	t.width = width
	t.height = height
	return nil
}

// handleKeyPress handles key presses in the thread list
func (t *ThreadList) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if !t.focused {
		return t, nil
	}

	switch msg.String() {
	case "j":
		if t.selected < len(t.threads)-1 {
			t.selected++
		}
	case "k":
		if t.selected > 0 {
			t.selected--
		}
	case "enter":
		// Select thread
	case "gg":
		t.selected = 0
	case "G":
		t.selected = len(t.threads) - 1
	}

	return t, nil
}
