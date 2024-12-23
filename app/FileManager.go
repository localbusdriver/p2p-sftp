package app

import (
	"bytes"
	"p2p-sftp/handlers"
	"time"

	"github.com/google/uuid"
)

func (a *App) UploadFile(filename string, fileData []byte) error {
	userId, err := a.userManager.GetUserID()
	if err != nil {
		return err
	}

	reader := bytes.NewReader(fileData)

	file := &handlers.FileUpload{
		Id:          uuid.New().String(),
		Filename:    filename,
		StoragePath: a.fileHandler.StorageBasePath + "/" + userId + "/" + filename,
		Size:        int64(len(fileData)),
		Uploadtime:  time.Now(),
		UserID:      userId,
		Status:      "pending",
	}

	uploadError := a.fileHandler.UploadFile(file, reader)
	return uploadError
}

func (a *App) GetUploadedFiles(userId string) ([]handlers.FileUpload, error) {
	return a.fileHandler.GetUploadedFiles(userId)
}

func (a *App) DeleteUpload(uploadId string, userId string) error {
	return a.fileHandler.DeleteUpload(uploadId, userId)
}

func (a *App) ValidateFile(size int64, filename string) error {
	return a.fileHandler.ValidateFile(size, filename)
}

func (a *App) CleanupOldUploads(age time.Duration) error {
	return a.fileHandler.CleanupOldUploads(age)
}

func (a *App) GetFileUpload(uploadId string) (*handlers.FileUpload, error) {
	userId, err := a.userManager.GetUserID()
	if err != nil {
		return nil, err
	}
	return a.fileHandler.GetSelectedUpload(uploadId, userId)
}

func (a *App) GetAllUploadedFiles(uploadId string) ([]handlers.FileUpload, error) {
	userId, err := a.userManager.GetUserID()
	if err != nil {
		return nil, err
	}
	return a.fileHandler.GetUploadedFiles(userId)
}
