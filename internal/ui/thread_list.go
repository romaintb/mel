package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/romaintb/mel/internal/config"
	"github.com/romaintb/mel/internal/email"
	"github.com/romaintb/mel/internal/icons"
)

// ThreadList represents the list of email threads
type ThreadList struct {
	config       *config.Config
	emailManager *email.Manager
	iconService  *icons.Service
	width        int
	height       int
	focused      bool
	selected     int
	threads      []ThreadItem
}

// Thread represents an email thread
type ThreadItem struct {
	ID      string
	Subject string
	From    string
	Date    string
	Unread  bool
	Starred bool
}

// NewThreadList creates a new thread list instance
func NewThreadList(cfg *config.Config, emailManager *email.Manager, iconService *icons.Service) (*ThreadList, error) {
	// Mock data for now
	threads := []ThreadItem{
		{ID: "1", Subject: "Welcome to Mel", From: "team@mel.com", Date: "2024-01-15", Unread: true, Starred: false},
		{ID: "2", Subject: "Project Update", From: "manager@work.com", Date: "2024-01-14", Unread: false, Starred: true},
		{ID: "3", Subject: "Meeting Tomorrow", From: "colleague@work.com", Date: "2024-01-13", Unread: true, Starred: false},
	}

	return &ThreadList{
		config:       cfg,
		emailManager: emailManager,
		iconService:  iconService,
		width:        0, // Will be set by Resize
		height:       0,
		focused:      false,
		selected:     0,
		threads:      threads,
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
	result += t.iconService.Get("email") + " Threads\n"
	result += "─────────\n"

	for i, thread := range t.threads {
		prefix := "  "
		if i == t.selected {
			prefix = t.iconService.Get("selected") + " "
		}

		unread := ""
		if thread.Unread {
			unread = t.iconService.Get("unread") + " "
		}

		starred := ""
		if thread.Starred {
			starred = t.iconService.Get("star") + " "
		}

		result += prefix + unread + starred + thread.Subject + "\n"
		result += "   " + thread.From + " • " + thread.Date + "\n"
		if i < len(t.threads)-1 {
			result += "\n"
		}
	}

	// Add more mock threads if we have more height available
	if t.height > 15 {
		result += "\n" + t.iconService.Get("email") + " More Threads\n"
		result += "──────────────\n"
		result += "  " + t.iconService.Get("email") + " Weekly Report\n"
		result += "   team@company.com • 2024-01-12\n"
		result += "\n  " + t.iconService.Get("email") + " Calendar Update\n"
		result += "   calendar@company.com • 2024-01-11\n"
		result += "\n  " + t.iconService.Get("email") + " Team Chat\n"
		result += "   chat@company.com • 2024-01-10\n"
	}

	if t.height > 25 {
		result += "\n" + t.iconService.Get("archive") + " Archive\n"
		result += "─────────\n"
		result += "  " + t.iconService.Get("email") + " Q4 Results\n"
		result += "   finance@company.com • 2024-01-09\n"
		result += "\n  " + t.iconService.Get("email") + " Holiday Party\n"
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

// GoToTop goes to the first thread
func (t *ThreadList) GoToTop() tea.Cmd {
	t.selected = 0
	return nil
}

// GoToBottom goes to the last thread
func (t *ThreadList) GoToBottom() tea.Cmd {
	t.selected = len(t.threads) - 1
	return nil
}

// Next goes to the next thread
func (t *ThreadList) Next() tea.Cmd {
	if t.selected < len(t.threads)-1 {
		t.selected++
	}
	return nil
}

// Prev goes to the previous thread
func (t *ThreadList) Prev() tea.Cmd {
	if t.selected > 0 {
		t.selected--
	}
	return nil
}

// NextUnread goes to the next unread thread
func (t *ThreadList) NextUnread() tea.Cmd {
	for i := t.selected + 1; i < len(t.threads); i++ {
		if t.threads[i].Unread {
			t.selected = i
			break
		}
	}
	return nil
}

// PrevUnread goes to the previous unread thread
func (t *ThreadList) PrevUnread() tea.Cmd {
	for i := t.selected - 1; i >= 0; i-- {
		if t.threads[i].Unread {
			t.selected = i
			break
		}
	}
	return nil
}

// ToggleThread toggles thread expansion
func (t *ThreadList) ToggleThread() tea.Cmd {
	// TODO: Implement thread expansion/collapse
	return nil
}

// ArchiveCurrent archives the current thread
func (t *ThreadList) ArchiveCurrent() tea.Cmd {
	if t.selected >= 0 && t.selected < len(t.threads) {
		// TODO: Implement actual archiving via email manager
		return nil
	}
	return nil
}

// DeleteCurrent deletes the current thread
func (t *ThreadList) DeleteCurrent() tea.Cmd {
	if t.selected >= 0 && t.selected < len(t.threads) {
		// TODO: Implement actual deletion via email manager
		return nil
	}
	return nil
}

// ToggleStar toggles star status of current thread
func (t *ThreadList) ToggleStar() tea.Cmd {
	if t.selected >= 0 && t.selected < len(t.threads) {
		t.threads[t.selected].Starred = !t.threads[t.selected].Starred
		return nil
	}
	return nil
}

// MarkRead marks the current thread as read
func (t *ThreadList) MarkRead() tea.Cmd {
	if t.selected >= 0 && t.selected < len(t.threads) {
		t.threads[t.selected].Unread = false
		return nil
	}
	return nil
}

// MarkUnread marks the current thread as unread
func (t *ThreadList) MarkUnread() tea.Cmd {
	if t.selected >= 0 && t.selected < len(t.threads) {
		t.threads[t.selected].Unread = true
		return nil
	}
	return nil
}
