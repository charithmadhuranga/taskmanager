package main

import (
	"log"
	"os"

	"tappmanager/internal/app"
	"tappmanager/internal/services"
	"tappmanager/internal/ui/models"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Create application
	application, err := app.NewApp()
	if err != nil {
		log.Fatalf("Failed to create application: %v", err)
	}

	// Create storage and process service
	storage := application.GetStorage()
	processService := services.NewProcessService(storage)

	// Create main model
	model := models.NewMainModel(storage, processService)

	// Create Bubble Tea program
	program := tea.NewProgram(model, tea.WithAltScreen())

	// Run the program
	if _, err := program.Run(); err != nil {
		log.Fatalf("Application error: %v", err)
		os.Exit(1)
	}
}
