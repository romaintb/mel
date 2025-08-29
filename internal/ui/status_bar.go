package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/romaintb/mel/internal/config"
)

// StatusBar represents the bottom status bar
type StatusBar struct {
	config     *config.Config
	width      int
	height     int
	message    string
	mode       string
	focusedBox string
}

// NewStatusBar creates a new status bar instance
func NewStatusBar(cfg *config.Config) (*StatusBar, error) {
	return &StatusBar{
		config:  cfg,
		width:   0,
		height:  1,
		message: "Ready",
		mode:    "NORMAL",
	}, nil
}

// Init initializes the status bar
func (s *StatusBar) Init() tea.Cmd {
	return nil
}

// Update handles status bar updates
func (s *StatusBar) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return s.handleKeyPress(msg)
	}
	return s, nil
}

// View renders the status bar
func (s *StatusBar) View() string {
	if s.width == 0 {
		return ""
	}

	// Left side: mode, focused box, and message
	left := "[" + s.mode + "][" + s.focusedBox + "] " + s.message

	// Right side: shortcuts
	right := "q:quit h:sidebar l:list i:insert v:visual /:search"

	// Calculate spacing
	spacing := s.width - len(left) - len(right)
	if spacing < 0 {
		spacing = 0
	}

	// Build the status bar
	result := left
	for i := 0; i < spacing; i++ {
		result += " "
	}
	result += right

	return result
}

// SetMessage sets the status message
func (s *StatusBar) SetMessage(msg string) {
	s.message = msg
}

// SetMode sets the current mode
func (s *StatusBar) SetMode(mode string) {
	s.mode = mode
}

// SetFocusedBox sets the currently focused box
func (s *StatusBar) SetFocusedBox(box string) {
	s.focusedBox = box
}

// Resize resizes the status bar
func (s *StatusBar) Resize(width, height int) tea.Cmd {
	s.width = width
	s.height = height
	return nil
}

// handleKeyPress handles key presses in the status bar
func (s *StatusBar) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Status bar doesn't handle key presses directly
	return s, nil
}
