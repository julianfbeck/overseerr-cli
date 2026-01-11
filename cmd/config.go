package cmd

import (
	"fmt"

	"github.com/julianfbeck/overseerr-cli/internal/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage CLI configuration",
}

var setURLCmd = &cobra.Command{
	Use:   "set-url <url>",
	Short: "Set Overseerr server URL",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := config.SetURL(args[0]); err != nil {
			return fmt.Errorf("failed to save URL: %w", err)
		}
		if !quietMode {
			fmt.Printf("URL set to: %s\n", args[0])
		}
		return nil
	},
}

var setKeyCmd = &cobra.Command{
	Use:   "set-key <api-key>",
	Short: "Set API key",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := config.SetAPIKey(args[0]); err != nil {
			return fmt.Errorf("failed to save API key: %w", err)
		}
		if !quietMode {
			fmt.Println("API key saved")
		}
		return nil
	},
}

var showConfigCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}

		if jsonOutput {
			outputJSON(cfg)
			return nil
		}

		fmt.Printf("URL: %s\n", cfg.URL)
		if cfg.APIKey != "" {
			fmt.Printf("API Key: %s...%s\n", cfg.APIKey[:8], cfg.APIKey[len(cfg.APIKey)-4:])
		} else {
			fmt.Println("API Key: (not set)")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(setURLCmd)
	configCmd.AddCommand(setKeyCmd)
	configCmd.AddCommand(showConfigCmd)
}
