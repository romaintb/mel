package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/romaintb/mel/internal/config"
)

// UI represents the main user interface
type UI struct {
	// Configuration
	config *config.Config

	// Current view/mode
	currentView ViewType

	// UI components
	sidebar    *Sidebar
	threadList *ThreadList
	threadView *ThreadView
	statusBar  *StatusBar

	// Dimensions
	width  int
	height int

	// Styles
	styles *UIStyles
}

// UIStyles holds all the styling for the UI
type UIStyles struct {
	sidebar   lipgloss.Style
	content   lipgloss.Style
	statusBar lipgloss.Style
}

// ViewType represents the current view mode
type ViewType int

const (
	ViewNormal ViewType = iota
	ViewInsert
	ViewVisual
	ViewSearch
)

// New creates a new UI instance
func New(cfg *config.Config) (*UI, error) {
	sidebar, err := NewSidebar(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create sidebar: %w", err)
	}

	threadList, err := NewThreadList(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create thread list: %w", err)
	}

	threadView, err := NewThreadView(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create thread view: %w", err)
	}

	statusBar, err := NewStatusBar(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create status bar: %w", err)
	}

	// Initialize styles
	styles := &UIStyles{
		sidebar: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(1, 1),
		content: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(1, 1),
		statusBar: lipgloss.NewStyle().
			Background(lipgloss.Color("240")).
			Foreground(lipgloss.Color("255")).
			Padding(0, 1),
	}

	return &UI{
		config:      cfg,
		currentView: ViewNormal,
		sidebar:     sidebar,
		threadList:  threadList,
		threadView:  threadView,
		statusBar:   statusBar,
		styles:      styles,
	}, nil
}

// Init initializes the UI
func (u *UI) Init() tea.Cmd {
	return tea.Batch(
		u.sidebar.Init(),
		u.threadList.Init(),
		u.threadView.Init(),
		u.statusBar.Init(),
	)
}

// Update handles UI updates
func (u *UI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		cmds = append(cmds, u.handleKeyPress(msg)...)
	case tea.WindowSizeMsg:
		u.width = msg.Width
		u.height = msg.Height
		cmds = append(cmds, u.handleResize(msg)...)
	}

	// Update child components
	if cmd := u.updateSidebar(msg); cmd != nil {
		cmds = append(cmds, cmd)
	}
	if cmd := u.updateThreadList(msg); cmd != nil {
		cmds = append(cmds, cmd)
	}
	if cmd := u.updateThreadView(msg); cmd != nil {
		cmds = append(cmds, cmd)
	}
	if cmd := u.updateStatusBar(msg); cmd != nil {
		cmds = append(cmds, cmd)
	}

	return u, tea.Batch(cmds...)
}

// View renders the UI using proper lipgloss layout
func (u *UI) View() string {
	if u.width == 0 {
		return "Initializing..."
	}

	// Fixed layout dimensions - grid-like approach
	sidebarWidth := 30                         // Fixed 30 characters
	contentWidth := u.width - sidebarWidth - 2 // Account for borders
	contentHeight := u.height - 2              // Account for status bar and borders

	// Update component dimensions to fill their allocated space
	u.sidebar.Resize(sidebarWidth-2, contentHeight-2)    // Account for border padding
	u.threadList.Resize(contentWidth-2, contentHeight-2) // Account for border padding
	u.statusBar.Resize(u.width, 1)

	// Create styled components that fill their allocated space
	sidebar := u.styles.sidebar.
		Width(sidebarWidth).
		Height(contentHeight).
		Render(u.sidebar.View())

	content := u.styles.content.
		Width(contentWidth).
		Height(contentHeight).
		Render(u.renderContent())

	// Combine sidebar and content horizontally
	mainArea := lipgloss.JoinHorizontal(lipgloss.Left, sidebar, content)

	// Status bar at bottom - full width, single line
	status := u.styles.statusBar.
		Width(u.width).
		Height(1).
		Render(u.statusBar.View())

	return lipgloss.JoinVertical(lipgloss.Top, mainArea, status)
}

// renderContent renders the main content area
func (u *UI) renderContent() string {
	// For now, just show thread list
	// TODO: Implement proper view switching
	return u.threadList.View()
}

// handleKeyPress handles keyboard input
func (u *UI) handleKeyPress(msg tea.KeyMsg) []tea.Cmd {
	var cmds []tea.Cmd

	switch u.currentView {
	case ViewNormal:
		cmds = append(cmds, u.handleNormalMode(msg)...)
	case ViewInsert:
		cmds = append(cmds, u.handleInsertMode(msg)...)
	case ViewVisual:
		cmds = append(cmds, u.handleVisualMode(msg)...)
	case ViewSearch:
		cmds = append(cmds, u.handleSearchMode(msg)...)
	}

	return cmds
}

// handleNormalMode handles key presses in normal mode
func (u *UI) handleNormalMode(msg tea.KeyMsg) []tea.Cmd {
	var cmds []tea.Cmd

	switch msg.String() {
	case "q", "ctrl+c":
		return []tea.Cmd{tea.Quit}
	case "h":
		// Focus sidebar
		cmds = append(cmds, u.sidebar.Focus())
	case "l":
		// Focus thread list
		cmds = append(cmds, u.threadList.Focus())
	case "i":
		// Enter insert mode
		u.currentView = ViewInsert
		u.statusBar.SetMode("INSERT")
	case "v":
		// Enter visual mode
		u.currentView = ViewVisual
		u.statusBar.SetMode("VISUAL")
	case "/":
		// Enter search mode
		u.currentView = ViewSearch
		u.statusBar.SetMode("SEARCH")
	}

	return cmds
}

// handleInsertMode handles key presses in insert mode
func (u *UI) handleInsertMode(msg tea.KeyMsg) []tea.Cmd {
	var cmds []tea.Cmd

	switch msg.String() {
	case "esc":
		// Exit insert mode
		u.currentView = ViewNormal
		u.statusBar.SetMode("NORMAL")
	}

	return cmds
}

// handleVisualMode handles key presses in visual mode
func (u *UI) handleVisualMode(msg tea.KeyMsg) []tea.Cmd {
	var cmds []tea.Cmd

	switch msg.String() {
	case "esc":
		// Exit visual mode
		u.currentView = ViewNormal
		u.statusBar.SetMode("NORMAL")
	}

	return cmds
}

// handleSearchMode handles key presses in search mode
func (u *UI) handleSearchMode(msg tea.KeyMsg) []tea.Cmd {
	var cmds []tea.Cmd

	switch msg.String() {
	case "esc":
		// Exit search mode
		u.currentView = ViewNormal
		u.statusBar.SetMode("NORMAL")
	}

	return cmds
}

// handleResize handles window resize events
func (u *UI) handleResize(msg tea.WindowSizeMsg) []tea.Cmd {
	var cmds []tea.Cmd

	// Update UI dimensions first
	u.width = msg.Width
	u.height = msg.Height

	// Calculate component dimensions
	sidebarWidth := 30                         // Fixed 30 characters
	contentWidth := u.width - sidebarWidth - 2 // Account for borders
	contentHeight := u.height - 2              // Account for status bar and borders

	// Resize child components with their allocated dimensions
	if cmd := u.sidebar.Resize(sidebarWidth-2, contentHeight-2); cmd != nil {
		cmds = append(cmds, cmd)
	}
	if cmd := u.threadList.Resize(contentWidth-2, contentHeight-2); cmd != nil {
		cmds = append(cmds, cmd)
	}
	if cmd := u.threadView.Resize(contentWidth-2, contentHeight-2); cmd != nil {
		cmds = append(cmds, cmd)
	}
	if cmd := u.statusBar.Resize(u.width, 1); cmd != nil {
		cmds = append(cmds, cmd)
	}

	return cmds
}

// Helper methods for updating child components
func (u *UI) updateSidebar(msg tea.Msg) tea.Cmd {
	_, cmd := u.sidebar.Update(msg)
	return cmd
}

func (u *UI) updateThreadList(msg tea.Msg) tea.Cmd {
	_, cmd := u.threadList.Update(msg)
	return cmd
}

func (u *UI) updateThreadView(msg tea.Msg) tea.Cmd {
	_, cmd := u.threadView.Update(msg)
	return cmd
}

func (u *UI) updateStatusBar(msg tea.Msg) tea.Cmd {
	_, cmd := u.statusBar.Update(msg)
	return cmd
}
