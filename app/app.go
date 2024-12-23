package app

import (
	"context"
	"fmt"

	"p2p-sftp/handlers"
)

type App struct {
	ctx         context.Context
	userManager *handlers.UserConfigManager
	fileHandler *handlers.FileHandler
}

// NewApp creates a new App application struct
func NewApp() *App {
	userManager := handlers.NewUserConfigManager()
	userConfig, getUserConfigErr := userManager.GetUser()
	if getUserConfigErr != nil {
		fmt.Printf("Failed to get user config: %v\n", getUserConfigErr)
		return nil // Handle error gracefully (exit the app)
	}
	fileHandler := handlers.NewFileHandler(userConfig)

	return &App{
		userManager: userManager,
		fileHandler: fileHandler,
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) OnStartup(ctx context.Context) {
	a.ctx = ctx

	// Initializer user config
	err := a.userManager.Init()
	if err != nil {
		fmt.Printf("Failed to initialize user config: %v\n", err)
		return // Handle error gracefully (exit the app)
	}
}
