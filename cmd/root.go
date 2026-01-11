package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/julianfbeck/overseerr-cli/internal/api"
	"github.com/julianfbeck/overseerr-cli/internal/config"
	"github.com/spf13/cobra"
)

var (
	jsonOutput bool
	quietMode  bool
	noColor    bool
	urlFlag    string
	version    = "dev"
	ctx        = context.Background()
)

var rootCmd = &cobra.Command{
	Use:   "overseerr",
	Short: "CLI for Overseerr media request management",
	Long: `A command-line interface for Overseerr - the media request management tool.

Manage media requests, search for movies and TV shows, and interact with your
Overseerr instance from the command line.`,
	Version: version,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "Output as JSON")
	rootCmd.PersistentFlags().BoolVarP(&quietMode, "quiet", "q", false, "Suppress non-essential output")
	rootCmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "Disable color output")
	rootCmd.PersistentFlags().StringVarP(&urlFlag, "url", "u", "", "Override Overseerr URL")
}

func getClient() (*api.OverseerrClient, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	if urlFlag != "" {
		cfg.URL = urlFlag
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return api.NewOverseerrClient(cfg.URL, cfg.APIKey)
}

func outputJSON(v interface{}) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(v)
}

func printInfo(format string, args ...interface{}) {
	if !quietMode {
		fmt.Printf(format, args...)
	}
}

func printError(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
}

// Helper to safely dereference string pointers
func derefStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// Helper to safely dereference float32 pointers
func derefFloat(f *float32) float32 {
	if f == nil {
		return 0
	}
	return *f
}

// Helper to safely dereference int pointers
func derefInt(i *int) int {
	if i == nil {
		return 0
	}
	return *i
}
