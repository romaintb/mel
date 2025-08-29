package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/romaintb/mel/internal/config"
	"github.com/romaintb/mel/internal/email"
	"github.com/romaintb/mel/internal/icons"
	"github.com/romaintb/mel/internal/search"
)

// UI represents the main user interface
type UI struct {
	// Configuration
	config *config.Config

	// Services
	emailManager  *email.Manager
	searchService *search.SearchService
	iconService   *icons.Service

	// Current view/mode
	currentView ViewType

	// Leader key state
	leaderPressed bool

	// Focus management
	focusedBox FocusedBox

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
	sidebar        lipgloss.Style
	sidebarFocused lipgloss.Style
	content        lipgloss.Style
	contentFocused lipgloss.Style
	statusBar      lipgloss.Style
}

// ViewType represents the current view mode
type ViewType int

const (
	ViewNormal ViewType = iota
	ViewInsert
	ViewVisual
	ViewSearch
)

// FocusedBox represents which box is currently focused
type FocusedBox int

const (
	FocusedSidebar FocusedBox = iota
	FocusedContent
)

// New creates a new UI instance
func New(cfg *config.Config, emailManager *email.Manager, searchService *search.SearchService, iconService *icons.Service) (*UI, error) {
	sidebar, err := NewSidebar(cfg, emailManager, iconService)
	if err != nil {
		return nil, fmt.Errorf("failed to create sidebar: %w", err)
	}

	threadList, err := NewThreadList(cfg, emailManager, iconService)
	if err != nil {
		return nil, fmt.Errorf("failed to create thread list: %w", err)
	}

	threadView, err := NewThreadView(cfg, emailManager, iconService)
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
			BorderForeground(lipgloss.Color("240")),
		sidebarFocused: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("117")), // Light blue (more blue, less white) for focused sidebar
		content: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")),
		contentFocused: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("117")), // Light blue (more blue, less white) for focused content
		statusBar: lipgloss.NewStyle().
			Background(lipgloss.Color("240")).
			Foreground(lipgloss.Color("255")).
			Padding(0, 1),
	}

	// Initialize status bar with default focus
	statusBar.SetFocusedBox("SIDEBAR")

	return &UI{
		config:        cfg,
		emailManager:  emailManager,
		searchService: searchService,
		iconService:   iconService,
		currentView:   ViewNormal,
		leaderPressed: false,
		focusedBox:    FocusedSidebar, // Default focus to sidebar
		sidebar:       sidebar,
		threadList:    threadList,
		threadView:    threadView,
		statusBar:     statusBar,
		styles:        styles,
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
	contentWidth := u.width - sidebarWidth - 4 // Account for borders (2) + spacing (2)
	contentHeight := u.height - 4              // Account for status bar (1) + borders (2)

	// Update component dimensions to fill their allocated space
	u.sidebar.Resize(sidebarWidth-2, contentHeight-2)    // Account for border padding (2)
	u.threadList.Resize(contentWidth-2, contentHeight-2) // Account for border padding (2)
	u.statusBar.Resize(u.width, 1)

	// Create styled components that fill their allocated space
	var sidebarStyle, contentStyle lipgloss.Style
	if u.focusedBox == FocusedSidebar {
		sidebarStyle = u.styles.sidebarFocused
		contentStyle = u.styles.content
	} else {
		sidebarStyle = u.styles.sidebar
		contentStyle = u.styles.contentFocused
	}

	sidebar := sidebarStyle.
		Width(sidebarWidth).
		Height(contentHeight).
		Render(u.sidebar.View())

	content := contentStyle.
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

	switch {
	case msg.String() == "q":
		return []tea.Cmd{tea.Quit}
	case msg.Type == tea.KeyCtrlC:
		return []tea.Cmd{tea.Quit}
	case msg.String() == "h":
		// Toggle sidebar visibility
		cmds = append(cmds, u.sidebar.Toggle())
		u.statusBar.SetMessage("Sidebar toggled")
		u.leaderPressed = false
	case msg.String() == "l":
		// Focus thread list
		cmds = append(cmds, u.threadList.Focus())
		u.leaderPressed = false
	case msg.String() == "i":
		// Enter insert mode
		u.currentView = ViewInsert
		u.statusBar.SetMode("INSERT")
	case msg.String() == "v":
		// Enter visual mode
		u.currentView = ViewVisual
		u.statusBar.SetMode("VISUAL")
	case msg.String() == "/":
		// Enter search mode
		u.currentView = ViewSearch
		u.statusBar.SetMode("SEARCH")
	case msg.String() == " ":
		// Leader key - show available commands
		u.statusBar.SetMessage("Leader key pressed - use: fg (content), fs (sender), fe (global), i (toggle icons)")
		u.leaderPressed = true
	case msg.Type == tea.KeyTab:
		// Switch focus between sidebar and content
		if u.focusedBox == FocusedSidebar {
			u.focusedBox = FocusedContent
			u.statusBar.SetFocusedBox("CONTENT")
			u.statusBar.SetMessage("Focus: Content")
		} else {
			u.focusedBox = FocusedSidebar
			u.statusBar.SetFocusedBox("SIDEBAR")
			u.statusBar.SetMessage("Focus: Sidebar")
		}
		u.leaderPressed = false
	case msg.String() == "i" && u.leaderPressed:
		// Toggle icon mode (leader+i)
		u.iconService.ToggleMode()
		u.statusBar.SetMessage(fmt.Sprintf("Icon mode: %s", u.iconService.GetModeString()))
		u.leaderPressed = false
	case msg.String() == "g":
		// Go to top of focused box
		if u.focusedBox == FocusedSidebar {
			// Go to top of sidebar
			cmds = append(cmds, u.sidebar.GoToTop())
		} else {
			// Go to top of thread list
			cmds = append(cmds, u.threadList.GoToTop())
		}
	case msg.String() == "G":
		// Go to bottom of focused box
		if u.focusedBox == FocusedSidebar {
			// Go to bottom of sidebar
			cmds = append(cmds, u.sidebar.GoToBottom())
		} else {
			// Go to bottom of thread list
			cmds = append(cmds, u.threadList.GoToBottom())
		}
	case msg.String() == "j":
		// Move down in focused box
		if u.focusedBox == FocusedSidebar {
			// Next item in sidebar
			cmds = append(cmds, u.sidebar.Next())
		} else {
			// Next thread in content
			cmds = append(cmds, u.threadList.Next())
		}
	case msg.String() == "k":
		// Move up in focused box
		if u.focusedBox == FocusedSidebar {
			// Previous item in sidebar
			cmds = append(cmds, u.sidebar.Prev())
		} else {
			// Previous thread in content
			cmds = append(cmds, u.threadList.Prev())
		}
	case msg.String() == "n":
		// Next unread thread
		cmds = append(cmds, u.threadList.NextUnread())
	case msg.String() == "p":
		// Previous unread thread
		cmds = append(cmds, u.threadList.PrevUnread())
	case msg.String() == "o":
		// Enter/select in focused box
		if u.focusedBox == FocusedSidebar {
			// TODO: Implement sidebar selection
			u.statusBar.SetMessage("Sidebar selection")
		} else {
			// Expand/collapse thread
			cmds = append(cmds, u.threadList.ToggleThread())
		}
	case msg.String() == "a":
		// Archive thread
		cmds = append(cmds, u.threadList.ArchiveCurrent())
	case msg.String() == "d":
		// Delete thread
		cmds = append(cmds, u.threadList.DeleteCurrent())
	case msg.String() == "s":
		// Star/unread thread
		cmds = append(cmds, u.threadList.ToggleStar())
	case msg.String() == "r":
		// Mark thread as read
		cmds = append(cmds, u.threadList.MarkRead())
	case msg.String() == "u":
		// Mark thread as unread
		cmds = append(cmds, u.threadList.MarkUnread())
	case msg.String() == "e":
		// Toggle sidebar (leader+e as specified in PRD)
		cmds = append(cmds, u.sidebar.Toggle())
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
	case " ":
		// Leader key in search mode
		u.statusBar.SetMessage("Search: fg (content), fs (sender), fe (global)")
	case "f":
		// Wait for next key to determine search type
		u.statusBar.SetMessage("Search type: g (content), s (sender), e (global)")
	case "g":
		// Content search (leader+fg)
		u.statusBar.SetMessage("Content search mode - type your query")
		cmds = append(cmds, u.startContentSearch())
	case "s":
		// Sender search (leader+fs)
		u.statusBar.SetMessage("Sender search mode - type the sender name")
		cmds = append(cmds, u.startSenderSearch())
	case "e":
		// Global search (leader+fe)
		u.statusBar.SetMessage("Global search mode - type your query")
		cmds = append(cmds, u.startGlobalSearch())
	}

	return cmds
}

// startContentSearch starts content search
func (u *UI) startContentSearch() tea.Cmd {
	return func() tea.Msg {
		// TODO: Implement actual search input handling
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f', 'g'}}
	}
}

// startSenderSearch starts sender search
func (u *UI) startSenderSearch() tea.Cmd {
	return func() tea.Msg {
		// TODO: Implement actual search input handling
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f', 's'}}
	}
}

// startGlobalSearch starts global search
func (u *UI) startGlobalSearch() tea.Cmd {
	return func() tea.Msg {
		// TODO: Implement actual search input handling
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f', 'e'}}
	}
}

// handleResize handles window resize events
func (u *UI) handleResize(msg tea.WindowSizeMsg) []tea.Cmd {
	var cmds []tea.Cmd

	// Update UI dimensions first
	u.width = msg.Width
	u.height = msg.Height

	// Calculate component dimensions
	sidebarWidth := 30                         // Fixed 30 characters
	contentWidth := u.width - sidebarWidth - 4 // Account for borders (2) + spacing (2)
	contentHeight := u.height - 4              // Account for status bar (1) + borders (2)

	// Resize child components with their allocated dimensions
	if cmd := u.sidebar.Resize(sidebarWidth-4, contentHeight-4); cmd != nil {
		cmds = append(cmds, cmd)
	}
	if cmd := u.threadList.Resize(contentWidth-4, contentHeight-4); cmd != nil {
		cmds = append(cmds, cmd)
	}
	if cmd := u.threadView.Resize(contentWidth-4, contentHeight-4); cmd != nil {
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
