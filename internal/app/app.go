package app

import (
	"os"

	"tappmanager/internal/storage"

	"github.com/rivo/tview"
)

// App represents the main application
type App struct {
	config  *Config
	storage storage.Storage
	ui      *tview.Application
}

// NewApp creates a new application instance
func NewApp() (*App, error) {
	config, err := LoadConfig()
	if err != nil {
		return nil, err
	}

	// Ensure TERM is set (tview/tcell requirement)
	if os.Getenv("TERM") == "" {
		os.Setenv("TERM", "xterm-256color")
	}

	storage := storage.NewJSONStorage(config.DataDir)
	
	// Load existing configuration
	if _, err := storage.LoadConfig(); err != nil {
		return nil, err
	}

	app := &App{
		config:  config,
		storage: storage,
		ui:      tview.NewApplication(),
	}

	return app, nil
}

// GetConfig returns the application configuration
func (a *App) GetConfig() *Config {
	return a.config
}

// GetStorage returns the storage interface
func (a *App) GetStorage() storage.Storage {
	return a.storage
}

// GetUI returns the UI application
func (a *App) GetUI() *tview.Application {
	return a.ui
}

// Run starts the application
func (a *App) Run() error {
	// This will be implemented in the UI layer
	return nil
}
