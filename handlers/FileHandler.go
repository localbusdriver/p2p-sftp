package handlers

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

// FileHandler is a struct that handles file operations
// This handler is responsible for storing files, managing files, tracking uploaded files for later transfer, and retrieving files

// Key concepts to implement:
// 1. File handler struct
// 2. Methods:
//    - NewFileHandler: Initialize the file handler
//    - UploadFile: Handle the file upload process
//    - GetUploadedFiles: List currently uploaded files
//    - ValidateFile: Check file size, type, etc.

type FileUpload struct {
	Id          string
	Filename    string // Original filename
	StoragePath string
	Size        int64     // File size in bytes
	Uploadtime  time.Time // Time of upload
	UserID      string
	Status      string // e.g. Completed, Failed, pending, deleting ...
}

type FileHandler struct {
	// Base configuration
	StorageBasePath string // Base path for all uploads (e.g. /storage/uploads...)
	MaxFileSize     int64
	AllowedTypes    []string // Allowed file extensions/types (e.g. ["jpg", "png", "pdf"])

	// State tracking
	ActiveUploads map[string]*FileUpload // Track current active uploads

	// Synchronisation
	mutex sync.Mutex // For thread-safe operations

	// Dependencies
	userConfig UserConfig // Reference to the user config
}

func NewFileHandler(userConfig *UserConfig) *FileHandler {
	// Initialize the file handler with default values
	userHomeDir, userHomeDirErr := os.UserHomeDir()
	if userHomeDirErr != nil {
		log.Printf("Failed to get user home directory: %v\n", userHomeDirErr)
		return nil
	}
	return &FileHandler{
		StorageBasePath: userHomeDir + "/.p2p-sftp/storage/uploads/" + userConfig.UserId,
		MaxFileSize:     1024 * 1024 * 10, // 10MB ... for now
		AllowedTypes:    []string{"jpg", "png", "pdf"},
		ActiveUploads:   map[string]*FileUpload{}, // map with [string] as key and [FileUpload] as value
		userConfig:      *userConfig,
	}
}

func (h *FileHandler) UploadFile(file *FileUpload, fileData io.Reader) error {
	// Handle file upload process

	// Set variables for uploading file
	file.Id = uuid.New().String()
	file.UserID = h.userConfig.UserId
	file.Uploadtime = time.Now()
	file.Status = "pending"

	// Create file at the storage path
	file.StoragePath = filepath.Join(h.StorageBasePath, file.Filename)

	dst, err := os.Create(file.StoragePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close() // close thread after function right before function ends

	written, err := io.Copy(dst, fileData) // write file data to storage path
	if err != nil {
		os.Remove(file.StoragePath) // clean partially uploaded file
		return fmt.Errorf("failed to write file: %w", err)
	}

	file.Size = written
	file.Status = "completed" // set status to completed

	h.mutex.Lock()                  // lock thread
	defer h.mutex.Unlock()          // unlock thread after function right before function ends
	h.ActiveUploads[file.Id] = file // add file to active uploads
	if err := h.ValidateFile(file.Size, file.Filename); err != nil {
		os.Remove(file.StoragePath)
		file.Status = "failed"
		return fmt.Errorf("failed to validate file: %w", err)
	}

	return nil
}

func (h *FileHandler) GetUploadedFiles(userId string) ([]FileUpload, error) {
	// List files uploaded by user
	var files []fs.DirEntry
	files, err := os.ReadDir(h.StorageBasePath + "/" + userId)
	if err != nil {
		return nil, fmt.Errorf("directory not found: %w", err)
	}

	var uploads []FileUpload
	for _, file := range files {
		if file.IsDir() { // skip directories
			continue
		}

		fileInfo, err := file.Info()
		if err != nil {
			fmt.Printf("failed to get file info: %v", err)
			continue
		}

		uploadedFile := FileUpload{
			Id:          h.userConfig.UserId,
			Filename:    file.Name(),
			StoragePath: filepath.Join(h.StorageBasePath, file.Name()),
			Size:        fileInfo.Size(),
			Uploadtime:  fileInfo.ModTime(),
			UserID:      userId,
			Status:      "completed",
		}

		uploads = append(uploads, uploadedFile)

	}
	return uploads, nil
}

func (h *FileHandler) DeleteUpload(uploadId string, userId string) error {
	// Remove selected uploaded file

	fullPath := filepath.Join(h.StorageBasePath, userId, uploadId)

	h.ActiveUploads[uploadId].Status = "deleting"
	h.mutex.Lock()
	defer h.mutex.Unlock()

	delete(h.ActiveUploads, uploadId)

	if err := os.Remove(fullPath); err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}

func (h *FileHandler) ValidateFile(size int64, filename string) error {
	// Check file size and type
	if size > h.MaxFileSize {
		return fmt.Errorf("file size exceeds maximum allowed size: %d", h.MaxFileSize)
	}

	for _, allowedType := range h.AllowedTypes {
		if strings.HasSuffix(filename, allowedType) {
			return nil
		}
	}
	return fmt.Errorf("file type not allowed: %s", filename)
}

func (h *FileHandler) CleanupOldUploads(age time.Duration) error {
	// Remove old uploads
	h.mutex.Lock()
	defer h.mutex.Unlock()
	files, err := os.ReadDir(h.StorageBasePath)
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	for _, file := range files {
		if file.IsDir() {
			err := h.CleanupOldUploads(age)
			if err != nil {
				log.Printf("failed to cleanup old uploads: %v\n%v", err, file)
			}
		}
		fileInfo, err := file.Info()
		if err != nil {
			log.Printf("failed to get file info: %v\n%v", err, fileInfo)
			continue
		}
		if time.Since(fileInfo.ModTime()) > age {
			if err := os.Remove(file.Name()); err != nil {
				return fmt.Errorf("failed to remove file: %w", err)
			}
		}
	}

	return nil
}

func (h *FileHandler) GetSelectedUpload(uploadId string, userId string) (*FileUpload, error) {
	// Get selected uploaded file
	fullPath := filepath.Join(h.StorageBasePath, userId, uploadId)
	file, err := os.Open(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}
	uploadedFile := FileUpload{
		Id:          uploadId,
		Filename:    fileInfo.Name(),
		StoragePath: fullPath,
		Size:        fileInfo.Size(),
		Uploadtime:  fileInfo.ModTime(),
		UserID:      userId,
		Status:      "completed",
	}

	return &uploadedFile, nil
}
