package database

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/upsync/backend/internal/config"
	"github.com/upsync/backend/internal/models"
)

// Client wraps Supabase REST and Storage APIs.
type Client struct {
	cfg        *config.Config
	httpClient *http.Client
}

// New creates a new Supabase client.
func New(cfg *config.Config) *Client {
	return &Client{
		cfg:        cfg,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// ─── Helpers ──────────────────────────────────────────────────

func (c *Client) restURL(path string) string {
	return fmt.Sprintf("%s/rest/v1/%s", c.cfg.SupabaseURL, path)
}

func (c *Client) storageURL(path string) string {
	return fmt.Sprintf("%s/storage/v1/%s", c.cfg.SupabaseURL, path)
}

func (c *Client) authHeaders(req *http.Request) {
	req.Header.Set("apikey", c.cfg.SupabaseServiceKey)
	req.Header.Set("Authorization", "Bearer "+c.cfg.SupabaseServiceKey)
}

// ─── Database Operations ──────────────────────────────────────

// InsertFile saves file metadata to the files table.
func (c *Client) InsertFile(f *models.File) error {
	body, _ := json.Marshal(map[string]any{
		"id":            f.ID,
		"original_name": f.OriginalName,
		"mime_type":     f.MimeType,
		"size":          f.Size,
		"storage_path":  f.StoragePath,
		"expires_at":    f.ExpiresAt.UTC().Format(time.RFC3339),
	})

	req, _ := http.NewRequest(http.MethodPost, c.restURL("files"), bytes.NewReader(body))
	c.authHeaders(req)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Prefer", "return=minimal")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("insert file request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("insert file: status %d: %s", resp.StatusCode, string(b))
	}
	return nil
}

// GetFile retrieves file metadata by ID.
func (c *Client) GetFile(id string) (*models.File, error) {
	url := c.restURL("files") + "?id=eq." + id + "&limit=1"
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	c.authHeaders(req)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("get file request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get file: status %d", resp.StatusCode)
	}

	var rows []models.File
	if err := json.NewDecoder(resp.Body).Decode(&rows); err != nil {
		return nil, fmt.Errorf("decode file: %w", err)
	}
	if len(rows) == 0 {
		return nil, nil // not found
	}
	return &rows[0], nil
}

// GetExpiredFiles returns all files whose expires_at is in the past.
func (c *Client) GetExpiredFiles() ([]models.File, error) {
	now := time.Now().UTC().Format(time.RFC3339)
	url := c.restURL("files") + "?expires_at=lt." + now + "&select=id,storage_path"

	req, _ := http.NewRequest(http.MethodGet, url, nil)
	c.authHeaders(req)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("get expired files: %w", err)
	}
	defer resp.Body.Close()

	var rows []models.File
	if err := json.NewDecoder(resp.Body).Decode(&rows); err != nil {
		return nil, fmt.Errorf("decode expired files: %w", err)
	}
	return rows, nil
}

// DeleteFile removes a file record from the database.
func (c *Client) DeleteFile(id string) error {
	url := c.restURL("files") + "?id=eq." + id
	req, _ := http.NewRequest(http.MethodDelete, url, nil)
	c.authHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("delete file request: %w", err)
	}
	defer resp.Body.Close()
	return nil
}

// ─── Storage Operations ───────────────────────────────────────

// UploadObject uploads file bytes to Supabase Storage.
func (c *Client) UploadObject(storagePath, mimeType string, data []byte) error {
	url := c.storageURL(fmt.Sprintf("object/%s/%s", c.cfg.SupabaseBucket, storagePath))

	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	c.authHeaders(req)
	req.Header.Set("Content-Type", mimeType)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("upload object request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("upload object: status %d: %s", resp.StatusCode, string(b))
	}
	return nil
}

// DeleteObject removes a file from Supabase Storage.
func (c *Client) DeleteObject(storagePath string) error {
	url := c.storageURL(fmt.Sprintf("object/%s", c.cfg.SupabaseBucket))

	body, _ := json.Marshal(map[string]any{"prefixes": []string{storagePath}})
	req, _ := http.NewRequest(http.MethodDelete, url, bytes.NewReader(body))
	c.authHeaders(req)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("delete object request: %w", err)
	}
	defer resp.Body.Close()
	return nil
}

// CreateSignedURL generates a short-lived signed download URL (expiresInSec seconds).
func (c *Client) CreateSignedURL(storagePath string, expiresInSec int) (string, error) {
	url := c.storageURL(fmt.Sprintf("object/sign/%s/%s", c.cfg.SupabaseBucket, storagePath))

	body, _ := json.Marshal(map[string]any{"expiresIn": expiresInSec})
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	c.authHeaders(req)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("create signed url request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("create signed url: status %d: %s", resp.StatusCode, string(b))
	}

	var result struct {
		SignedURL string `json:"signedURL"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decode signed url: %w", err)
	}

	// Signed URL from Supabase is a relative path starting with /object/sign...
	// We must prepend the base storage URL: https://[project].supabase.co/storage/v1
	basePath := c.cfg.SupabaseURL + "/storage/v1"

	if strings.HasPrefix(result.SignedURL, "/") {
		result.SignedURL = basePath + result.SignedURL
	} else {
		result.SignedURL = basePath + "/" + result.SignedURL
	}
	return result.SignedURL, nil
}
