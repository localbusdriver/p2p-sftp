package app

// StoreUserConfig stores the user configuration
func (a *App) SetUsername(username string) error {
	return a.userManager.SetUsername(username)
}

// GetUserConfig retrieves the user configuration
func (a *App) GetUsername() (string, error) {
	return a.userManager.GetUsername()
}

func (a *App) GetUserId() (string, error) {
	return a.userManager.GetUserID()
}
