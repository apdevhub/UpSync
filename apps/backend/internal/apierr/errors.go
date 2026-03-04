package apierr

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AppError is a structured application error with an HTTP status code.
type AppError struct {
	StatusCode int
	Message    string
	Code       string
}

func (e *AppError) Error() string { return e.Message }

// Predefined common errors
var (
	ErrNotFound     = &AppError{StatusCode: http.StatusNotFound, Message: "File not found", Code: "NOT_FOUND"}
	ErrExpired      = &AppError{StatusCode: http.StatusGone, Message: "This link has expired. The file has been deleted.", Code: "EXPIRED"}
	ErrFileTooLarge = &AppError{StatusCode: http.StatusRequestEntityTooLarge, Message: "File exceeds the 50 MB limit", Code: "FILE_TOO_LARGE"}
	ErrNoFile       = &AppError{StatusCode: http.StatusBadRequest, Message: "No file provided", Code: "NO_FILE"}
	ErrInternal     = &AppError{StatusCode: http.StatusInternalServerError, Message: "Internal server error", Code: "INTERNAL"}
)

// New creates a custom AppError.
func New(status int, message, code string) *AppError {
	return &AppError{StatusCode: status, Message: message, Code: code}
}

// Respond writes an AppError as a JSON response and aborts the Gin chain.
func Respond(c *gin.Context, err *AppError) {
	c.AbortWithStatusJSON(err.StatusCode, gin.H{
		"error": err.Message,
		"code":  err.Code,
	})
}

// RespondInternal logs the real error and sends a generic 500 to the client.
func RespondInternal(c *gin.Context, internal error) {
	c.Error(internal) // picked up by the logging middleware
	Respond(c, ErrInternal)
}
