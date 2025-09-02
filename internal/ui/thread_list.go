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
	scrollOffset int // How many items are scrolled up
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
	return &ThreadList{
		config:       cfg,
		emailManager: emailManager,
		iconService:  iconService,
		width:        0, // Will be set by Resize
		height:       0,
		focused:      false,
		selected:     0,
		threads:      []ThreadItem{}, // Start empty, will be populated by LoadThreads
	}, nil
}

// Init initializes the thread list
func (t *ThreadList) Init() tea.Cmd {
	// Load threads from INBOX by default
	return t.LoadThreads("INBOX")
}

// threadsLoadedMsg is sent when threads are loaded
type threadsLoadedMsg struct {
	threads []*email.Thread
	folder  string
	err     error
}

// Update handles thread list updates
func (t *ThreadList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return t.handleKeyPress(msg)
	case threadsLoadedMsg:
		return t.handleThreadsLoaded(msg)
	}
	return t, nil
}

// LoadThreads loads threads from a specific folder
func (t *ThreadList) LoadThreads(folderName string) tea.Cmd {
	return func() tea.Msg {
		threads, err := t.emailManager.GetThreadsFromFolder(folderName)
		if err != nil {
			return threadsLoadedMsg{threads: nil, folder: folderName, err: err}
		}

		return threadsLoadedMsg{threads: threads, folder: folderName, err: nil}
	}
}

// getPrimarySender extracts the primary sender from participants
func (t *ThreadList) getPrimarySender(participants []string) string {
	if len(participants) > 0 {
		return participants[0]
	}
	return "Unknown"
}

// handleThreadsLoaded handles when threads are loaded
func (t *ThreadList) handleThreadsLoaded(msg threadsLoadedMsg) (tea.Model, tea.Cmd) {
	if msg.err != nil {
		// On error, keep existing threads but could show error message
		return t, nil
	}

	// Convert threads to ThreadItems
	var threadItems []ThreadItem
	for _, thread := range msg.threads {
		from := t.getPrimarySender(thread.Participants)
		date := thread.Timestamp.Format("2006-01-02")
		unread := thread.UnreadCount > 0

		item := ThreadItem{
			ID:      thread.ID,
			Subject: thread.Subject,
			From:    from,
			Date:    date,
			Unread:  unread,
			Starred: false, // TODO: Check if thread is starred
		}

		threadItems = append(threadItems, item)
	}

	t.threads = threadItems
	t.selected = 0     // Reset selection to first thread
	t.scrollOffset = 0 // Reset scroll offset

	return t, nil
}

// View renders the thread list
func (t *ThreadList) View() string {
	if t.width == 0 {
		return ""
	}

	if len(t.threads) == 0 {
		return t.iconService.Get("email") + " Threads\n─────────\nNo threads"
	}

	var result string
	result += t.iconService.Get("email") + " Threads\n"
	result += "─────────\n"

	// Calculate how many thread items can fit in the available height
	// Each thread takes 1 line
	availableHeight := t.height - 2        // Subtract header height
	maxVisibleItems := availableHeight - 2 // Subtract space for scroll indicators

	// Ensure we don't try to show more items than we have
	totalItems := len(t.threads)
	if maxVisibleItems > totalItems {
		maxVisibleItems = totalItems
	}

	// Calculate the range of items to display
	startIdx := t.scrollOffset
	endIdx := startIdx + maxVisibleItems
	if endIdx > totalItems {
		endIdx = totalItems
	}

	// Display the visible threads (1 line per thread)
	for i := startIdx; i < endIdx; i++ {
		thread := t.threads[i]

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

		// Calculate available width for content (subtract prefix length)
		prefixLength := len(prefix) + len(unread) + len(starred)
		availableWidth := t.width - prefixLength - 1 // -1 for newline

		// Build the line with truncation: [subject] from [sender] • [date]
		subject := thread.Subject
		sender := thread.From
		date := thread.Date

		// Calculate the fixed parts: " from " + sender + " • " + date
		fixedParts := " from " + sender + " • " + date
		fixedLength := len(fixedParts)

		// Truncate subject if needed to fit within available width
		if fixedLength+len(subject) > availableWidth && availableWidth > 10 {
			maxSubjectLen := availableWidth - fixedLength - 3 // -3 for "..."
			if maxSubjectLen > 0 && len(subject) > maxSubjectLen {
				subject = subject[:maxSubjectLen] + "..."
			}
		}

		line := prefix + unread + starred + subject + " from " + sender + " • " + date + "\n"
		result += line
	}

	// Add scroll indicators if needed
	if startIdx > 0 {
		// Show scroll up indicator at the top
		scrollUpText := t.iconService.Get("scrollUp") + " More above..."
		if len(scrollUpText) > t.width {
			scrollUpText = scrollUpText[:t.width-3] + "..."
		}
		result = scrollUpText + "\n" + result
	}
	if endIdx < totalItems {
		// Show scroll down indicator at the bottom
		scrollDownText := t.iconService.Get("scrollDown") + " More below..."
		if len(scrollDownText) > t.width {
			scrollDownText = scrollDownText[:t.width-3] + "..."
		}
		result += "\n" + scrollDownText
	}

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
	t.scrollOffset = 0
	return nil
}

// GoToBottom goes to the last thread
func (t *ThreadList) GoToBottom() tea.Cmd {
	if len(t.threads) == 0 {
		return nil
	}

	t.selected = len(t.threads) - 1

	// Calculate how many items can be visible
	availableHeight := t.height - 2
	maxVisibleItems := availableHeight - 2 // Subtract space for scroll indicators
	if maxVisibleItems > len(t.threads) {
		maxVisibleItems = len(t.threads)
	}

	// Scroll so that the selected item is visible at the bottom
	if len(t.threads) > maxVisibleItems {
		t.scrollOffset = len(t.threads) - maxVisibleItems
	} else {
		t.scrollOffset = 0
	}

	return nil
}

// Next goes to the next thread
func (t *ThreadList) Next() tea.Cmd {
	if len(t.threads) == 0 {
		return nil
	}

	// Calculate how many items can be visible
	availableHeight := t.height - 2
	maxVisibleItems := availableHeight - 2 // Subtract space for scroll indicators
	if maxVisibleItems > len(t.threads) {
		maxVisibleItems = len(t.threads)
	}

	// Move selection down
	if t.selected < len(t.threads)-1 {
		t.selected++

		// Check if we need to scroll down
		visibleStart := t.scrollOffset
		visibleEnd := visibleStart + maxVisibleItems

		if t.selected >= visibleEnd {
			t.scrollOffset++
		}
	}
	return nil
}

// Prev goes to the previous thread
func (t *ThreadList) Prev() tea.Cmd {
	if len(t.threads) == 0 || t.selected <= 0 {
		return nil
	}

	// Move selection up
	t.selected--

	// Check if we need to scroll up
	if t.selected < t.scrollOffset {
		t.scrollOffset--
	}
	return nil
}

// NextUnread goes to the next unread thread
func (t *ThreadList) NextUnread() tea.Cmd {
	for i := t.selected + 1; i < len(t.threads); i++ {
		if t.threads[i].Unread {
			t.selected = i
			// Adjust scroll offset to make the selected item visible
			t.adjustScrollForSelection()
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
			// Adjust scroll offset to make the selected item visible
			t.adjustScrollForSelection()
			break
		}
	}
	return nil
}

// adjustScrollForSelection adjusts the scroll offset to make the currently selected item visible
func (t *ThreadList) adjustScrollForSelection() {
	if len(t.threads) == 0 {
		return
	}

	// Calculate how many items can be visible
	availableHeight := t.height - 2
	maxVisibleItems := availableHeight - 2 // Subtract space for scroll indicators
	if maxVisibleItems > len(t.threads) {
		maxVisibleItems = len(t.threads)
	}

	// Ensure selected item is visible
	if t.selected < t.scrollOffset {
		t.scrollOffset = t.selected
	} else if t.selected >= t.scrollOffset+maxVisibleItems {
		t.scrollOffset = t.selected - maxVisibleItems + 1
	}

	// Ensure scroll offset doesn't go negative
	if t.scrollOffset < 0 {
		t.scrollOffset = 0
	}
}

// PageDown scrolls down by one page
func (t *ThreadList) PageDown() tea.Cmd {
	if len(t.threads) == 0 {
		return nil
	}

	// Calculate how many items can be visible
	availableHeight := t.height - 2
	maxVisibleItems := availableHeight - 2 // Subtract space for scroll indicators
	if maxVisibleItems <= 0 {
		maxVisibleItems = 1
	}

	// Scroll down by one page
	newScrollOffset := t.scrollOffset + maxVisibleItems
	maxScrollOffset := len(t.threads) - maxVisibleItems
	if maxScrollOffset < 0 {
		maxScrollOffset = 0
	}

	if newScrollOffset > maxScrollOffset {
		newScrollOffset = maxScrollOffset
	}

	if newScrollOffset != t.scrollOffset {
		t.scrollOffset = newScrollOffset
		// Move selection to the first visible item
		t.selected = t.scrollOffset
	}

	return nil
}

// PageUp scrolls up by one page
func (t *ThreadList) PageUp() tea.Cmd {
	if len(t.threads) == 0 {
		return nil
	}

	// Calculate how many items can be visible
	availableHeight := t.height - 2
	maxVisibleItems := availableHeight - 2 // Subtract space for scroll indicators
	if maxVisibleItems <= 0 {
		maxVisibleItems = 1
	}

	// Scroll up by one page
	newScrollOffset := t.scrollOffset - maxVisibleItems
	if newScrollOffset < 0 {
		newScrollOffset = 0
	}

	if newScrollOffset != t.scrollOffset {
		t.scrollOffset = newScrollOffset
		// Move selection to the first visible item
		t.selected = t.scrollOffset
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
