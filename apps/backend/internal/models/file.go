package models

import "time"

// File represents the metadata stored in the database for an uploaded file.
type File struct {
	ID           string    `json:"id"`
	OriginalName string    `json:"original_name"`
	MimeType     string    `json:"mime_type"`
	Size         int64     `json:"size"`
	StoragePath  string    `json:"storage_path"`
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
}

// UploadResponse is returned to the client after a successful upload.
type UploadResponse struct {
	ID           string    `json:"id"`
	OriginalName string    `json:"originalName"`
	Size         int64     `json:"size"`
	MimeType     string    `json:"mimeType"`
	ExpiresAt    time.Time `json:"expiresAt"`
	ShareURL     string    `json:"shareUrl"`
}

// FileMetaResponse is returned when fetching file metadata for the share page.
type FileMetaResponse struct {
	ID           string    `json:"id"`
	OriginalName string    `json:"originalName"`
	MimeType     string    `json:"mimeType"`
	Size         int64     `json:"size"`
	ExpiresAt    time.Time `json:"expiresAt"`
	CreatedAt    time.Time `json:"createdAt"`
}

// DownloadResponse is returned when the client requests a download link.
type DownloadResponse struct {
	DownloadURL string `json:"downloadUrl"`
	FileName    string `json:"fileName"`
}

// ErrorResponse is the standard JSON error shape returned by the API.
type ErrorResponse struct {
	Error string `json:"error"`
	Code  string `json:"code,omitempty"`
}

// ExpiryDuration maps human-readable expiry keys to seconds.
var ExpiryDurations = map[string]int{
	"1h":  1 * 60 * 60,
	"6h":  6 * 60 * 60,
	"12h": 12 * 60 * 60,
	"24h": 24 * 60 * 60,
	"3d":  3 * 24 * 60 * 60,
	"7d":  7 * 24 * 60 * 60,
}
