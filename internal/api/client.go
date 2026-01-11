package api

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

// OverseerrClient wraps the generated client with API key authentication
type OverseerrClient struct {
	*ClientWithResponses
	apiKey string
}

// NewOverseerrClient creates a new Overseerr API client
func NewOverseerrClient(baseURL, apiKey string) (*OverseerrClient, error) {
	baseURL = strings.TrimSuffix(baseURL, "/")
	if !strings.HasSuffix(baseURL, "/api/v1") {
		baseURL = baseURL + "/api/v1"
	}

	client, err := NewClientWithResponses(baseURL, WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
		req.Header.Set("X-Api-Key", apiKey)
		return nil
	}))
	if err != nil {
		return nil, err
	}

	return &OverseerrClient{
		ClientWithResponses: client,
		apiKey:              apiKey,
	}, nil
}

// Helper function to get a pointer to a value
func Ptr[T any](v T) *T {
	return &v
}

// StatusString returns a human-readable status string
func StatusString(status *float32) string {
	if status == nil {
		return "Unknown"
	}
	switch int(*status) {
	case 1:
		return "Unknown"
	case 2:
		return "Pending"
	case 3:
		return "Processing"
	case 4:
		return "Partially Available"
	case 5:
		return "Available"
	default:
		return fmt.Sprintf("Status(%d)", int(*status))
	}
}

// RequestStatusString returns a human-readable request status
func RequestStatusString(status *float32) string {
	if status == nil {
		return "Unknown"
	}
	switch int(*status) {
	case 1:
		return "Pending Approval"
	case 2:
		return "Approved"
	case 3:
		return "Declined"
	default:
		return fmt.Sprintf("Status(%d)", int(*status))
	}
}

// MediaTypeString returns a human-readable media type
func MediaTypeString(mediaType *string) string {
	if mediaType == nil {
		return "Unknown"
	}
	switch *mediaType {
	case "movie":
		return "Movie"
	case "tv":
		return "TV"
	default:
		return *mediaType
	}
}
