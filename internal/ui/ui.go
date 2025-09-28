package ui

import (
	"tappmanager/internal/app"
	"tappmanager/internal/services"
	"tappmanager/internal/ui/models"

	tea "github.com/charmbracelet/bubbletea"
)

// UIApp represents the UI application
type UIApp struct {
	app           *app.App
	processService *services.ProcessService
	program       *tea.Program
}

// NewUIApp creates a new UI application
func NewUIApp(app *app.App) *UIApp {
	// Get storage from app
	storage := app.GetStorage()
	
	// Create process service
	processService := services.NewProcessService(storage)
	
	// Create main model
	model := models.NewMainModel(storage, processService)
	
	// Create Bubble Tea program
	program := tea.NewProgram(model, tea.WithAltScreen())
	
	return &UIApp{
		app:           app,
		processService: processService,
		program:       program,
	}
}

// Run starts the UI application
func (u *UIApp) Run() error {
	// Run the Bubble Tea program
	if _, err := u.program.Run(); err != nil {
		return err
	}
	return nil
}

// GetProcessService returns the process service
func (u *UIApp) GetProcessService() *services.ProcessService {
	return u.processService
}

// Stop stops the UI application
func (u *UIApp) Stop() {
	if u.program != nil {
		u.program.Quit()
	}
}

// SendMessage sends a message to the UI program
func (u *UIApp) SendMessage(msg tea.Msg) {
	if u.program != nil {
		u.program.Send(msg)
	}
}
