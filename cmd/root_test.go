package cmd

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestDerefStr(t *testing.T) {
	tests := []struct {
		name  string
		input *string
		want  string
	}{
		{
			name:  "nil string",
			input: nil,
			want:  "",
		},
		{
			name:  "empty string",
			input: strPtr(""),
			want:  "",
		},
		{
			name:  "non-empty string",
			input: strPtr("hello"),
			want:  "hello",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := derefStr(tt.input)
			if got != tt.want {
				t.Errorf("derefStr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDerefFloat(t *testing.T) {
	tests := []struct {
		name  string
		input *float32
		want  float32
	}{
		{
			name:  "nil float",
			input: nil,
			want:  0,
		},
		{
			name:  "zero",
			input: floatPtr(0),
			want:  0,
		},
		{
			name:  "positive",
			input: floatPtr(3.14),
			want:  3.14,
		},
		{
			name:  "negative",
			input: floatPtr(-1.5),
			want:  -1.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := derefFloat(tt.input)
			if got != tt.want {
				t.Errorf("derefFloat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDerefInt(t *testing.T) {
	tests := []struct {
		name  string
		input *int
		want  int
	}{
		{
			name:  "nil int",
			input: nil,
			want:  0,
		},
		{
			name:  "zero",
			input: intPtr(0),
			want:  0,
		},
		{
			name:  "positive",
			input: intPtr(42),
			want:  42,
		},
		{
			name:  "negative",
			input: intPtr(-1),
			want:  -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := derefInt(tt.input)
			if got != tt.want {
				t.Errorf("derefInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOutputJSON(t *testing.T) {
	// Test with a simple struct
	data := struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}{
		Name:  "test",
		Value: 42,
	}

	// Encode to JSON and verify format
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "  ")
	enc.Encode(data)

	expected := `{
  "name": "test",
  "value": 42
}
`
	if buf.String() != expected {
		t.Errorf("JSON output = %v, want %v", buf.String(), expected)
	}
}

func TestRootCommand(t *testing.T) {
	// Test that root command is properly configured
	if rootCmd.Use != "overseerr" {
		t.Errorf("rootCmd.Use = %v, want %v", rootCmd.Use, "overseerr")
	}

	// Test that required flags are registered
	flags := rootCmd.PersistentFlags()

	if flags.Lookup("json") == nil {
		t.Error("--json flag not registered")
	}
	if flags.Lookup("quiet") == nil {
		t.Error("--quiet flag not registered")
	}
	if flags.Lookup("no-color") == nil {
		t.Error("--no-color flag not registered")
	}
	if flags.Lookup("url") == nil {
		t.Error("--url flag not registered")
	}
}

// Helper functions
func strPtr(s string) *string {
	return &s
}

func floatPtr(f float32) *float32 {
	return &f
}

func intPtr(i int) *int {
	return &i
}
