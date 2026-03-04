package services

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/upsync/backend/internal/config"
	"github.com/upsync/backend/internal/database"
	"github.com/upsync/backend/internal/models"
)

// FileService contains all business logic for file management.
type FileService struct {
	db  *database.Client
	cfg *config.Config
}

// New creates a new FileService.
func New(db *database.Client, cfg *config.Config) *FileService {
	return &FileService{db: db, cfg: cfg}
}

// Upload stores the file in Supabase Storage and records metadata in the DB.
func (s *FileService) Upload(
	fileBytes []byte,
	originalName string,
	mimeType string,
	expiresIn string,
) (*models.UploadResponse, error) {
	// Resolve expiry duration
	secs, ok := models.ExpiryDurations[expiresIn]
	if !ok {
		secs = models.ExpiryDurations["24h"]
	}
	expiresAt := time.Now().UTC().Add(time.Duration(secs) * time.Second)

	fileID := uuid.New().String()
	storagePath := fmt.Sprintf("%s/%s", fileID, originalName)

	// Upload to Supabase Storage
	if err := s.db.UploadObject(storagePath, mimeType, fileBytes); err != nil {
		return nil, fmt.Errorf("storage upload: %w", err)
	}

	// Persist metadata
	meta := &models.File{
		ID:           fileID,
		OriginalName: originalName,
		MimeType:     mimeType,
		Size:         int64(len(fileBytes)),
		StoragePath:  storagePath,
		ExpiresAt:    expiresAt,
	}
	if err := s.db.InsertFile(meta); err != nil {
		// Best-effort cleanup of orphaned storage object
		_ = s.db.DeleteObject(storagePath)
		return nil, fmt.Errorf("db insert: %w", err)
	}

	return &models.UploadResponse{
		ID:           fileID,
		OriginalName: originalName,
		Size:         int64(len(fileBytes)),
		MimeType:     mimeType,
		ExpiresAt:    expiresAt,
		ShareURL:     fmt.Sprintf("%s/share/%s", s.cfg.FrontendURL, fileID),
	}, nil
}

// GetMeta returns public metadata for a file, checking expiry.
func (s *FileService) GetMeta(id string) (*models.FileMetaResponse, error) {
	f, err := s.db.GetFile(id)
	if err != nil {
		return nil, fmt.Errorf("db get: %w", err)
	}
	if f == nil {
		return nil, nil // not found
	}
	return &models.FileMetaResponse{
		ID:           f.ID,
		OriginalName: f.OriginalName,
		MimeType:     f.MimeType,
		Size:         f.Size,
		ExpiresAt:    f.ExpiresAt,
		CreatedAt:    f.CreatedAt,
	}, nil
}

// GetDownloadURL validates the file and returns a short-lived signed URL.
func (s *FileService) GetDownloadURL(id string) (*models.DownloadResponse, error) {
	f, err := s.db.GetFile(id)
	if err != nil {
		return nil, fmt.Errorf("db get: %w", err)
	}
	if f == nil {
		return nil, nil // not found
	}
	if f.ExpiresAt.Before(time.Now()) {
		return nil, fmt.Errorf("expired")
	}

	// 60-second signed URL
	signed, err := s.db.CreateSignedURL(f.StoragePath, 60)
	if err != nil {
		return nil, fmt.Errorf("signed url: %w", err)
	}

	return &models.DownloadResponse{
		DownloadURL: signed,
		FileName:    f.OriginalName,
	}, nil
}

// DeleteExpiredFiles cleans up all expired files from DB and storage.
func (s *FileService) DeleteExpiredFiles() (int, error) {
	files, err := s.db.GetExpiredFiles()
	if err != nil {
		return 0, fmt.Errorf("get expired: %w", err)
	}

	deleted := 0
	for _, f := range files {
		if err := s.db.DeleteObject(f.StoragePath); err != nil {
			fmt.Printf("[cleanup] storage delete error for %s: %v\n", f.ID, err)
		}
		if err := s.db.DeleteFile(f.ID); err != nil {
			fmt.Printf("[cleanup] db delete error for %s: %v\n", f.ID, err)
			continue
		}
		deleted++
	}
	return deleted, nil
}
