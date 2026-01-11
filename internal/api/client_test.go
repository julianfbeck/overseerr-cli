package api

import (
	"testing"
)

func TestNewOverseerrClient_URLNormalization(t *testing.T) {
	tests := []struct {
		name     string
		inputURL string
		wantBase string
	}{
		{
			name:     "plain URL",
			inputURL: "https://example.com",
			wantBase: "https://example.com/api/v1",
		},
		{
			name:     "URL with trailing slash",
			inputURL: "https://example.com/",
			wantBase: "https://example.com/api/v1",
		},
		{
			name:     "URL already has api/v1",
			inputURL: "https://example.com/api/v1",
			wantBase: "https://example.com/api/v1",
		},
		{
			name:     "URL with api/v1 and trailing slash",
			inputURL: "https://example.com/api/v1/",
			wantBase: "https://example.com/api/v1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewOverseerrClient(tt.inputURL, "test-key")
			if err != nil {
				t.Fatalf("NewOverseerrClient() error = %v", err)
			}
			if client == nil {
				t.Fatal("NewOverseerrClient() returned nil client")
			}
		})
	}
}

func TestStatusString(t *testing.T) {
	tests := []struct {
		name   string
		status *float32
		want   string
	}{
		{
			name:   "nil status",
			status: nil,
			want:   "Unknown",
		},
		{
			name:   "status 1 - Unknown",
			status: Ptr(float32(1)),
			want:   "Unknown",
		},
		{
			name:   "status 2 - Pending",
			status: Ptr(float32(2)),
			want:   "Pending",
		},
		{
			name:   "status 3 - Processing",
			status: Ptr(float32(3)),
			want:   "Processing",
		},
		{
			name:   "status 4 - Partially Available",
			status: Ptr(float32(4)),
			want:   "Partially Available",
		},
		{
			name:   "status 5 - Available",
			status: Ptr(float32(5)),
			want:   "Available",
		},
		{
			name:   "status 99 - fallback",
			status: Ptr(float32(99)),
			want:   "Status(99)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := StatusString(tt.status)
			if got != tt.want {
				t.Errorf("StatusString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRequestStatusString(t *testing.T) {
	tests := []struct {
		name   string
		status *float32
		want   string
	}{
		{
			name:   "nil status",
			status: nil,
			want:   "Unknown",
		},
		{
			name:   "status 1 - Pending Approval",
			status: Ptr(float32(1)),
			want:   "Pending Approval",
		},
		{
			name:   "status 2 - Approved",
			status: Ptr(float32(2)),
			want:   "Approved",
		},
		{
			name:   "status 3 - Declined",
			status: Ptr(float32(3)),
			want:   "Declined",
		},
		{
			name:   "status 99 - fallback",
			status: Ptr(float32(99)),
			want:   "Status(99)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RequestStatusString(tt.status)
			if got != tt.want {
				t.Errorf("RequestStatusString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMediaTypeString(t *testing.T) {
	tests := []struct {
		name      string
		mediaType *string
		want      string
	}{
		{
			name:      "nil media type",
			mediaType: nil,
			want:      "Unknown",
		},
		{
			name:      "movie",
			mediaType: Ptr("movie"),
			want:      "Movie",
		},
		{
			name:      "tv",
			mediaType: Ptr("tv"),
			want:      "TV",
		},
		{
			name:      "other",
			mediaType: Ptr("other"),
			want:      "other",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MediaTypeString(tt.mediaType)
			if got != tt.want {
				t.Errorf("MediaTypeString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPtr(t *testing.T) {
	// Test with string
	strVal := "test"
	strPtr := Ptr(strVal)
	if *strPtr != strVal {
		t.Errorf("Ptr(string) = %v, want %v", *strPtr, strVal)
	}

	// Test with int
	intVal := 42
	intPtr := Ptr(intVal)
	if *intPtr != intVal {
		t.Errorf("Ptr(int) = %v, want %v", *intPtr, intVal)
	}

	// Test with float32
	floatVal := float32(3.14)
	floatPtr := Ptr(floatVal)
	if *floatPtr != floatVal {
		t.Errorf("Ptr(float32) = %v, want %v", *floatPtr, floatVal)
	}
}
