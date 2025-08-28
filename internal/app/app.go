package app

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/romaintb/mel/internal/config"
	"github.com/romaintb/mel/internal/ui"
)

// App represents the main application
type App struct {
	ui     *ui.UI
	config *config.Config
}

// New creates a new application instance
func New(version string) (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	ui, err := ui.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize UI: %w", err)
	}

	return &App{
		ui:     ui,
		config: cfg,
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
