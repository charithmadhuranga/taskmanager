package main

import (
	"log"

	"tappmanager/internal/app"
	"tappmanager/internal/ui"
)

func main() {
	// Create application
	application, err := app.NewApp()
	if err != nil {
		log.Fatalf("Failed to create application: %v", err)
	}

	// Create UI
	uiApp := ui.NewUIApp(application)

	// Run application
	if err := uiApp.Run(); err != nil {
		log.Fatalf("Application error: %v", err)
	}
}
