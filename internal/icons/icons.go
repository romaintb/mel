package icons

// IconMode represents the current icon display mode
type IconMode int

const (
	IconModeEmoji IconMode = iota
	IconModeASCII
)

// IconSet holds all the icons for a specific mode
type IconSet struct {
	// Email and communication
	Email   string
	Inbox   string
	Sent    string
	Drafts  string
	Trash   string
	Starred string
	Archive string
	Folder  string
	Spam    string

	// Actions
	Compose  string
	Search   string
	Settings string
	Reply    string
	Forward  string
	Delete   string

	// Status indicators
	Unread string
	Read   string
	Star   string
	Unstar string

	// Navigation
	Next     string
	Previous string
	Top      string
	Bottom   string

	// UI elements
	Selected   string
	Collapsed  string
	Expanded   string
	ScrollUp   string
	ScrollDown string

	// Statistics
	Total        string
	UnreadCount  string
	StarredCount string
}

// Service manages icon display modes
type Service struct {
	currentMode IconMode
	emojiSet    *IconSet
	asciiSet    *IconSet
}

// NewService creates a new icon service
func NewService(mode IconMode) *Service {
	return &Service{
		currentMode: mode,
		emojiSet:    createEmojiSet(),
		asciiSet:    createASCIISet(),
	}
}

// SetMode sets the current icon mode
func (s *Service) SetMode(mode IconMode) {
	s.currentMode = mode
}

// GetMode returns the current icon mode
func (s *Service) GetMode() IconMode {
	return s.currentMode
}

// Get retrieves an icon by name for the current mode
func (s *Service) Get(iconName string) string {
	var iconSet *IconSet

	switch s.currentMode {
	case IconModeEmoji:
		iconSet = s.emojiSet
	case IconModeASCII:
		iconSet = s.asciiSet
	default:
		iconSet = s.emojiSet
	}

	return s.getIconValue(iconSet, iconName)
}

// GetWithFallback retrieves an icon with fallback to ASCII mode
func (s *Service) GetWithFallback(iconName string) string {
	value := s.Get(iconName)
	if value == "" {
		// Fallback to ASCII mode if emoji not found
		return s.getIconValue(s.asciiSet, iconName)
	}
	return value
}

// ToggleMode switches between emoji and ASCII modes
func (s *Service) ToggleMode() {
	if s.currentMode == IconModeEmoji {
		s.currentMode = IconModeASCII
	} else {
		s.currentMode = IconModeEmoji
	}
}

// IsEmojiMode returns true if currently in emoji mode
func (s *Service) IsEmojiMode() bool {
	return s.currentMode == IconModeEmoji
}

// IsASCIIMode returns true if currently in ASCII mode
func (s *Service) IsASCIIMode() bool {
	return s.currentMode == IconModeASCII
}

// GetModeString returns a human-readable string for the current mode
func (s *Service) GetModeString() string {
	switch s.currentMode {
	case IconModeEmoji:
		return "emoji"
	case IconModeASCII:
		return "ascii"
	default:
		return "unknown"
	}
}

// SetCustomIcon allows setting a custom icon for a specific name
func (s *Service) SetCustomIcon(iconName, value string) {
	switch s.currentMode {
	case IconModeEmoji:
		s.setCustomIconValue(s.emojiSet, iconName, value)
	case IconModeASCII:
		s.setCustomIconValue(s.asciiSet, iconName, value)
	}
}

// setCustomIconValue sets a custom icon value
func (s *Service) setCustomIconValue(iconSet *IconSet, iconName, value string) {
	switch iconName {
	case "email":
		iconSet.Email = value
	case "inbox":
		iconSet.Inbox = value
	case "sent":
		iconSet.Sent = value
	case "drafts":
		iconSet.Drafts = value
	case "trash":
		iconSet.Trash = value
	case "starred":
		iconSet.Starred = value
	case "archive":
		iconSet.Archive = value
	case "folder":
		iconSet.Folder = value
	case "spam":
		iconSet.Spam = value
	case "compose":
		iconSet.Compose = value
	case "search":
		iconSet.Search = value
	case "settings":
		iconSet.Settings = value
	case "reply":
		iconSet.Reply = value
	case "forward":
		iconSet.Forward = value
	case "delete":
		iconSet.Delete = value
	case "unread":
		iconSet.Unread = value
	case "read":
		iconSet.Read = value
	case "star":
		iconSet.Star = value
	case "unstar":
		iconSet.Unstar = value
	case "next":
		iconSet.Next = value
	case "previous":
		iconSet.Previous = value
	case "top":
		iconSet.Top = value
	case "bottom":
		iconSet.Bottom = value
	case "selected":
		iconSet.Selected = value
	case "collapsed":
		iconSet.Collapsed = value
	case "expanded":
		iconSet.Expanded = value
	case "scrollUp":
		iconSet.ScrollUp = value
	case "scrollDown":
		iconSet.ScrollDown = value
	case "total":
		iconSet.Total = value
	case "unreadCount":
		iconSet.UnreadCount = value
	case "starredCount":
		iconSet.StarredCount = value
	}
}

// getIconValue retrieves an icon value from an icon set
func (s *Service) getIconValue(iconSet *IconSet, iconName string) string {
	switch iconName {
	case "email":
		return iconSet.Email
	case "inbox":
		return iconSet.Inbox
	case "sent":
		return iconSet.Sent
	case "drafts":
		return iconSet.Drafts
	case "trash":
		return iconSet.Trash
	case "starred":
		return iconSet.Starred
	case "archive":
		return iconSet.Archive
	case "folder":
		return iconSet.Folder
	case "spam":
		return iconSet.Spam
	case "compose":
		return iconSet.Compose
	case "search":
		return iconSet.Search
	case "settings":
		return iconSet.Settings
	case "reply":
		return iconSet.Reply
	case "forward":
		return iconSet.Forward
	case "delete":
		return iconSet.Delete
	case "unread":
		return iconSet.Unread
	case "read":
		return iconSet.Read
	case "star":
		return iconSet.Star
	case "unstar":
		return iconSet.Unstar
	case "next":
		return iconSet.Next
	case "previous":
		return iconSet.Previous
	case "top":
		return iconSet.Top
	case "bottom":
		return iconSet.Bottom
	case "selected":
		return iconSet.Selected
	case "collapsed":
		return iconSet.Collapsed
	case "expanded":
		return iconSet.Expanded
	case "scrollUp":
		return iconSet.ScrollUp
	case "scrollDown":
		return iconSet.ScrollDown
	case "total":
		return iconSet.Total
	case "unreadCount":
		return iconSet.UnreadCount
	case "starredCount":
		return iconSet.StarredCount
	default:
		return ""
	}
}

// createEmojiSet creates the emoji icon set
func createEmojiSet() *IconSet {
	return &IconSet{
		// Email and communication
		Email:        "📧",
		Inbox:        "📥",
		Sent:         "📤",
		Drafts:       "📁",
		Trash:        "🗑️",
		Starred:      "⭐",
		Archive:      "📦",
		Folder:       "📁",
		Spam:         "🚫",
		Compose:      "📝",
		Search:       "🔍",
		Settings:     "⚙️",
		Reply:        "↩️",
		Forward:      "↪️",
		Delete:       "❌",
		Unread:       "●",
		Read:         "○",
		Star:         "⭐",
		Unstar:       "☆",
		Next:         "▶",
		Previous:     "◀",
		Top:          "⬆️",
		Bottom:       "⬇️",
		Selected:     "▶",
		Collapsed:    "▶",
		Expanded:     "▼",
		ScrollUp:     "↑",
		ScrollDown:   "↓",
		Total:        "📊",
		UnreadCount:  "●",
		StarredCount: "⭐",
	}
}

// createASCIISet creates the ASCII icon set with Neotree-style icons
func createASCIISet() *IconSet {
	return &IconSet{
		Email:   "📧",
		Inbox:   "📁",
		Sent:    "📤",
		Drafts:  "📝",
		Trash:   "🗑",
		Starred: "⭐",
		Archive: "📦",
		Folder:  "📁",
		Spam:    "🚫",

		// Actions - using Neotree-style action icons
		Compose:  "✏",
		Search:   "🔍",
		Settings: "⚙",
		Reply:    "↩",
		Forward:  "↪",
		Delete:   "✗",

		// Status indicators - using Neotree-style status icons
		Unread: "●",
		Read:   "○",
		Star:   "★",
		Unstar: "☆",

		// Navigation - using Neotree-style navigation icons
		Next:     "▶",
		Previous: "◀",
		Top:      "⬆",
		Bottom:   "⬇",

		// UI elements - using Neotree-style selection icons
		Selected:   "▶",
		Collapsed:  "▶",
		Expanded:   "▼",
		ScrollUp:   "↑",
		ScrollDown: "↓",

		// Statistics - using Neotree-style info icons
		Total:        "📊",
		UnreadCount:  "●",
		StarredCount: "★",
	}
}
