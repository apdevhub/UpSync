package handlers

import (
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/upsync/backend/internal/apierr"
	"github.com/upsync/backend/internal/config"
	"github.com/upsync/backend/internal/services"
)

// FileHandler handles all /api/files routes.
type FileHandler struct {
	svc *services.FileService
	cfg *config.Config
}

// NewFileHandler creates a new FileHandler.
func NewFileHandler(svc *services.FileService, cfg *config.Config) *FileHandler {
	return &FileHandler{svc: svc, cfg: cfg}
}

// Upload godoc
// POST /api/files/upload
// Accepts a multipart/form-data body with fields:
//   - file    — the file to upload (required)
//   - expiresIn — expiry key e.g. "24h" (optional, default "24h")
func (h *FileHandler) Upload(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		apierr.Respond(c, apierr.ErrNoFile)
		return
	}

	if fileHeader.Size > h.cfg.MaxFileSizeBytes {
		apierr.Respond(c, apierr.ErrFileTooLarge)
		return
	}

	f, err := fileHeader.Open()
	if err != nil {
		apierr.RespondInternal(c, err)
		return
	}
	defer f.Close()

	fileBytes, err := io.ReadAll(f)
	if err != nil {
		apierr.RespondInternal(c, err)
		return
	}

	expiresIn := c.PostForm("expiresIn")
	if expiresIn == "" {
		expiresIn = "24h"
	}

	// Detect content type if browser didn't supply one
	mimeType := fileHeader.Header.Get("Content-Type")
	if mimeType == "" || mimeType == "application/octet-stream" {
		mimeType = http.DetectContentType(fileBytes)
	}

	result, err := h.svc.Upload(fileBytes, fileHeader.Filename, mimeType, expiresIn)
	if err != nil {
		apierr.RespondInternal(c, err)
		return
	}

	c.JSON(http.StatusCreated, result)
}

// GetMeta godoc
// GET /api/files/:id
// Returns public metadata for a file and checks expiry.
func (h *FileHandler) GetMeta(c *gin.Context) {
	id := c.Param("id")

	meta, err := h.svc.GetMeta(id)
	if err != nil {
		apierr.RespondInternal(c, err)
		return
	}
	if meta == nil {
		apierr.Respond(c, apierr.ErrNotFound)
		return
	}
	if meta.ExpiresAt.Before(time.Now()) {
		apierr.Respond(c, apierr.ErrExpired)
		return
	}

	c.JSON(http.StatusOK, meta)
}

// GetDownloadURL godoc
// GET /api/files/:id/download
// Returns a 60-second signed URL for the file.
func (h *FileHandler) GetDownloadURL(c *gin.Context) {
	id := c.Param("id")

	result, err := h.svc.GetDownloadURL(id)
	if err != nil {
		if err.Error() == "expired" {
			apierr.Respond(c, apierr.ErrExpired)
			return
		}
		apierr.RespondInternal(c, err)
		return
	}
	if result == nil {
		apierr.Respond(c, apierr.ErrNotFound)
		return
	}

	c.JSON(http.StatusOK, result)
}
