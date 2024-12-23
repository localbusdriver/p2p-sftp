package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"os"

	"github.com/google/uuid"
)

type UserConfig struct {
	Username string `json:"username"`
	UserId   string `json:"userId"`
}

type UserConfigManager struct {
	currentUser UserConfig
}

func NewUserConfigManager() *UserConfigManager {
	return &UserConfigManager{}
}

func (m *UserConfigManager) getConfigPath() (string, error) { // get user home directory at C:\Users\<username> on Windows, /home/<username> on Linux
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		log.Printf("failed to get user home directory: %v\n", err)
		return "", err
	}
	return userHomeDir + "/.p2p-sftp", nil
}

func (m *UserConfigManager) ConfigExists(userConfigDir string) error {
	_, statError := os.Stat(userConfigDir + "/user-config.json") // check if config file exists
	if statError != nil {                                        // if config file doesn't exist
		_, createFileError := os.Create(userConfigDir + "/user-config.json")
		if createFileError != nil { // if can't create config file
			log.Printf("failed to create config file: %v\n", createFileError)
			return createFileError
		}
	}
	return nil
}

func (m *UserConfigManager) StoreUserConfig() error {
	userConfigDir, userHomeDirError := m.getConfigPath()
	if userHomeDirError != nil { // if can't find user home directory
		return userHomeDirError
	}

	configExistsError := m.ConfigExists(userConfigDir)
	if configExistsError != nil {
		return configExistsError
	}

	data, jsonError := json.MarshalIndent(m.currentUser, "", "  ")
	if jsonError != nil { // if can't marshal user config
		return jsonError
	}

	writeError := os.WriteFile(userConfigDir+"/user-config.json", data, 0644)
	if writeError != nil { // if can't write to config file
		return writeError
	}

	data, readError := os.ReadFile(userConfigDir + "/user-config.json")
	if readError != nil { // if can't read from config file
		return readError
	}

	err := json.Unmarshal(data, &m.currentUser)
	if err != nil { // if can't unmarshal user config
		return err
	}

	log.Printf("user config: %v\n", m.currentUser)

	return nil
}

func (m *UserConfigManager) FetchUserConfig() error {
	userConfigDir, userConfigDirError := m.getConfigPath()
	if userConfigDirError != nil { // if can't find user home directory
		return userConfigDirError
	}

	_, statError := os.Stat(userConfigDir + "/user-config.json") // check if config file exists
	if statError != nil {                                        // if config file doesn't exist
		log.Printf("config file doesn't exist: %v\n", statError)
		return statError
	}
	data, readError := os.ReadFile(userConfigDir + "/user-config.json")
	if readError != nil { // if can't read from config file
		return readError
	}

	// create user config struct
	err := json.Unmarshal(data, &m.currentUser) // unmarshal data into user config struct
	if err != nil {                             // if can't unmarshal user config
		return err
	}

	log.Printf("user config: %v\n", m.currentUser)
	return nil
}

func (m *UserConfigManager) SetUsername(username string) error {
	if username == "" {
		return errors.New("username cannot be empty")
	}
	m.currentUser.Username = username
	return m.StoreUserConfig()
}

func (m *UserConfigManager) SetUserID(userId string) error {
	m.currentUser.UserId = userId
	return m.StoreUserConfig()
}

func (m *UserConfigManager) SetNewUserId() error {
	newId := uuid.New().String()
	return m.SetUserID(newId)
}

func (m *UserConfigManager) Init() error {
	err := m.FetchUserConfig()
	if err != nil {
		log.Printf("Failed to get user config")
		return m.SetNewUserId()
	}
	return nil
}

func (m *UserConfigManager) GetUsername() (string, error) {
	if m.currentUser.Username == "" { // get username from currentUser, if it's empty try to get it from config file, and only if that fails return error
		err := m.FetchUserConfig()
		if err != nil {
			log.Printf("Username is empty")
			return "", err
		}
	}
	return m.currentUser.Username, nil
}

func (m *UserConfigManager) GetUserID() (string, error) {
	if m.currentUser.UserId == "" { // if user id is empty try to get it from config file, and only if that fails return error
		err := m.FetchUserConfig()
		if err != nil {
			log.Printf("UserID is empty")
			return "", err
		}
	}
	return m.currentUser.UserId, nil
}

func (m *UserConfigManager) ClearCurrentUser() {
	m.currentUser = UserConfig{}
}

func (m *UserConfigManager) GetUser() (*UserConfig, error) {
	if m.currentUser.Username == "" || m.currentUser.UserId == "" {
		err := m.FetchUserConfig()
		if err != nil {
			return nil, err
		}
	}
	return &m.currentUser, nil
}
