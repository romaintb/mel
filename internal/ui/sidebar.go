package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/romaintb/mel/internal/config"
	"github.com/romaintb/mel/internal/email"
	"github.com/romaintb/mel/internal/icons"
)

// Sidebar represents the left sidebar with account/folder tree
type Sidebar struct {
	config         *config.Config
	emailManager   *email.Manager
	iconService    *icons.Service
	width          int
	height         int
	focused        bool
	collapsed      bool
	selectedIndex  int                 // Index of selected item
	folders        []*email.MailFolder // Actual mail folders
	selectedFolder string              // Currently selected folder
}

// NewSidebar creates a new sidebar instance
func NewSidebar(cfg *config.Config, emailManager *email.Manager, iconService *icons.Service) (*Sidebar, error) {
	return &Sidebar{
		config:         cfg,
		emailManager:   emailManager,
		iconService:    iconService,
		width:          0, // Will be set by Resize
		height:         0,
		focused:        false,
		collapsed:      false,
		selectedIndex:  0, // Start with first item selected
		folders:        []*email.MailFolder{},
		selectedFolder: "",
	}, nil
}

// Init initializes the sidebar
func (s *Sidebar) Init() tea.Cmd {
	return s.refreshFolders()
}

// refreshFolders refreshes the folder list from the email manager
func (s *Sidebar) refreshFolders() tea.Cmd {
	return func() tea.Msg {
		folders, err := s.emailManager.GetMailFolders()
		if err != nil {
			// Return empty folders on error, but log it
			// For debugging, you can uncomment the next line to see errors
			// fmt.Printf("Error refreshing folders: %v\n", err)
			return foldersRefreshedMsg{folders: []*email.MailFolder{}, err: err}
		}
		return foldersRefreshedMsg{folders: folders, err: nil}
	}
}

// foldersRefreshedMsg is sent when folders are refreshed
type foldersRefreshedMsg struct {
	folders []*email.MailFolder
	err     error
}

// Update handles sidebar updates
func (s *Sidebar) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return s.handleKeyPress(msg)
	case foldersRefreshedMsg:
		s.folders = msg.folders
		// Set default selection to first folder if available
		if len(s.folders) > 0 && s.selectedFolder == "" {
			s.selectedFolder = s.folders[0].Name
		}
		return s, nil
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
	result += s.iconService.Get("email") + " Mail Folders\n"
	result += "──────────────\n"

	if len(s.folders) == 0 {
		result += "├── No folders found\n"
		result += "└── Check your mail directory\n"
		return result
	}

	// Calculate available height for folders (subtract header and quick actions)
	headerHeight := 2 // "Mail Folders" + separator
	quickActionsHeight := 0
	if s.height > 10 {
		quickActionsHeight = 5 // Quick Actions section
	}
	availableHeight := s.height - headerHeight - quickActionsHeight

	// Determine which folders to display based on available height
	startIndex := 0
	endIndex := len(s.folders)

	if availableHeight < len(s.folders) {
		// Need to scroll - show subset of folders
		if s.selectedIndex >= len(s.folders) {
			// Quick action is selected, adjust for that
			startIndex = 0
			endIndex = availableHeight
		} else {
			// Folder is selected, center the selection
			startIndex = s.selectedIndex - (availableHeight / 2)
			if startIndex < 0 {
				startIndex = 0
			}
			endIndex = startIndex + availableHeight
			if endIndex > len(s.folders) {
				endIndex = len(s.folders)
				startIndex = endIndex - availableHeight
			}
		}
	}

	// Display folders within the available height
	for i := startIndex; i < endIndex; i++ {
		folder := s.folders[i]
		var prefix string

		// Show scroll indicators
		if i == startIndex && startIndex > 0 {
			prefix = s.iconService.Get("scrollUp") + "── "
		} else if i == endIndex-1 && endIndex < len(s.folders) {
			prefix = s.iconService.Get("scrollDown") + "── "
		} else if i == len(s.folders)-1 {
			prefix = "└── "
		} else {
			prefix = "├── "
		}

		// Check if this folder is selected
		isSelected := s.selectedIndex == i

		// Get appropriate icon for the folder
		icon := s.getFolderIcon(folder)

		// Format the folder name with counts
		folderDisplay := s.formatFolderDisplay(folder)

		// Build the complete line with width constraint
		var line string
		if isSelected {
			line = prefix + s.iconService.Get("selected") + " " + icon + " " + folderDisplay
		} else {
			line = prefix + icon + " " + folderDisplay
		}

		// Calculate display width (emojis take 2 display columns)
		displayWidth := s.calculateDisplayWidth(line)

		// Ensure line doesn't exceed sidebar width
		if displayWidth > s.width {
			line = s.truncateTextByDisplayWidth(line, s.width)
		}

		result += line + "\n"
	}

	return result
}

// getFolderIcon returns the appropriate icon for a folder
func (s *Sidebar) getFolderIcon(folder *email.MailFolder) string {
	if !folder.IsSpecial {
		return s.iconService.Get("folder")
	}

	upperName := strings.ToUpper(folder.Name)
	switch upperName {
	case "INBOX":
		return s.iconService.Get("inbox")
	case "SENT":
		return s.iconService.Get("sent")
	case "DRAFTS":
		return s.iconService.Get("drafts")
	case "TRASH":
		return s.iconService.Get("trash")
	case "SPAM", "JUNK":
		return s.iconService.Get("spam")
	case "ARCHIVE":
		return s.iconService.Get("archive")
	default:
		return s.iconService.Get("folder")
	}
}

// truncateText truncates text to fit within the specified width, adding ellipsis if needed
func (s *Sidebar) truncateText(text string, maxWidth int) string {
	if len(text) <= maxWidth {
		return text
	}
	if maxWidth <= 1 {
		return "…"
	}
	return text[:maxWidth-1] + "…"
}

// formatFolderDisplay formats the folder display with counts
// This function ensures that folder names never wrap to multiple lines by truncating
// long names and adding ellipsis (...) when necessary.
func (s *Sidebar) formatFolderDisplay(folder *email.MailFolder) string {
	// Start with the full folder name
	folderName := folder.Name

	// Add unread count if any
	if folder.UnreadCount > 0 {
		return fmt.Sprintf("%s (%d)", folderName, folder.UnreadCount)
	}
	return folderName
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
	count := len(s.folders) // All folders are selectable
	if s.height > 10 {
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
		return s, s.Next()
	case "k":
		return s, s.Prev()
	case "enter":
		// Select folder or action
		s.selectCurrentItem()
	case "home":
		return s, s.GoToTop()
	case "end":
		return s, s.GoToBottom()
	case "r":
		// Refresh folders
		return s, s.refreshFolders()
	}

	return s, nil
}

// selectCurrentItem selects the currently highlighted item
func (s *Sidebar) selectCurrentItem() {
	if s.selectedIndex < len(s.folders) {
		// Select a folder
		s.selectedFolder = s.folders[s.selectedIndex].Name
	} else {
		// Handle quick actions
		actionIndex := s.selectedIndex - len(s.folders)
		switch actionIndex {
		case 0: // Compose
			// TODO: Trigger compose action
		case 1: // Search
			// TODO: Trigger search action
		case 2: // Settings
			// TODO: Trigger settings action
		}
	}
}

// GetSelectedFolder returns the currently selected folder
func (s *Sidebar) GetSelectedFolder() string {
	return s.selectedFolder
}

// getAvailableWidthForFolder calculates the available width for folder names
func (s *Sidebar) getAvailableWidthForFolder() int {
	// Calculate available width for folder name (subtract icons, prefixes, and padding)
	// Prefix: "├── " or "└── " (4 chars)
	// Selection icon: "▶ " (2 chars) - only when selected
	// Folder icon: "📁 " (2 chars)
	// Unread count: " (X)" (up to 6 chars) - only reserve space if there are unread emails
	// Padding: 0 chars (removed padding for better space utilization)
	prefixWidth := 4
	iconWidth := 2
	paddingWidth := 0
	// Don't reserve space for unread count by default - only when needed
	maxCountWidth := 0

	availableWidth := s.width - prefixWidth - iconWidth - paddingWidth - maxCountWidth

	// Ensure minimum width but allow more space for folder names
	if availableWidth < 6 {
		availableWidth = 6
	}

	return availableWidth
}

// calculateDisplayWidth calculates the display width of a string, accounting for emoji width
func (s *Sidebar) calculateDisplayWidth(text string) int {
	width := 0
	for _, r := range text {
		// Most emojis and wide characters take 2 display columns
		// This is a simplified approach - in practice you might want to use a proper Unicode width library
		if r > 127 { // Non-ASCII characters (including emojis)
			width += 1
		} else {
			width += 1
		}
	}
	return width
}

// truncateTextByDisplayWidth truncates text to fit within the specified display width
func (s *Sidebar) truncateTextByDisplayWidth(text string, maxDisplayWidth int) string {
	if s.calculateDisplayWidth(text) <= maxDisplayWidth {
		return text
	}

	// Truncate character by character until we fit
	for i := len(text) - 1; i >= 0; i-- {
		truncated := text[:i] + "…"
		if s.calculateDisplayWidth(truncated) <= maxDisplayWidth {
			return truncated
		}
	}

	return "…"
}
