package app

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/romaintb/mel/internal/config"
	"github.com/romaintb/mel/internal/email"
	"github.com/romaintb/mel/internal/icons"
	"github.com/romaintb/mel/internal/search"
	"github.com/romaintb/mel/internal/ui"
)

// App represents the main application
type App struct {
	ui            *ui.UI
	config        *config.Config
	emailManager  *email.Manager
	searchService *search.SearchService
	iconService   *icons.Service
}

// New creates a new application instance
func New(version string) (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Initialize icon service with configured mode
	var iconMode icons.IconMode
	switch strings.ToLower(strings.TrimSpace(cfg.UI.IconMode)) {
	case "", "ascii":
		iconMode = icons.IconModeASCII
	case "emoji":
		iconMode = icons.IconModeEmoji
	default:
		return nil, fmt.Errorf("invalid ui.iconMode %q; allowed: ascii, emoji", cfg.UI.IconMode)
	}
	iconService := icons.NewService(iconMode)

	// Initialize email manager with external tool paths
	if cfg.Email.Maildir == "" {
		return nil, fmt.Errorf("email.maildir is required")
	}

	emailManager := email.NewManager(
		cfg.Email.Maildir,
		cfg.ExternalTools.Notmuch,
		cfg.ExternalTools.Mbsync,
		cfg.ExternalTools.Msmtp,
	)

	// Initialize search service
	searchService := search.NewSearchService(emailManager)

	// Initialize UI with services
	ui, err := ui.New(cfg, emailManager, searchService, iconService)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize UI: %w", err)
	}

	return &App{
		ui:            ui,
		config:        cfg,
		emailManager:  emailManager,
		searchService: searchService,
		iconService:   iconService,
	}, nil
}

// Run starts the application
func Run(version string) error {
	app, err := New(version)
	if err != nil {
		return err
	}

	// Start the TUI program
	p := tea.NewProgram(
		app.ui,
		tea.WithAltScreen(),       // Use alternate screen buffer
		tea.WithMouseCellMotion(), // Turn on mouse support so we can track the mouse wheel
		tea.WithMouseAllMotion(),  // Turn on mouse support so we can track the mouse wheel
	)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v", err)
		os.Exit(1)
	}

	return nil
}
