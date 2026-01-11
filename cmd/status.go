package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show Overseerr server status",
	RunE:  runStatus,
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

func runStatus(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	resp, err := client.GetStatusWithResponse(ctx)
	if err != nil {
		return fmt.Errorf("failed to get status: %w", err)
	}

	if resp.JSON200 == nil {
		return fmt.Errorf("unexpected response: %s", resp.Status())
	}

	status := resp.JSON200

	if jsonOutput {
		outputJSON(status)
		return nil
	}

	fmt.Printf("Overseerr Status\n")
	fmt.Printf("  Version: %s\n", derefStr(status.Version))
	if status.CommitTag != nil && *status.CommitTag != "" {
		fmt.Printf("  Commit: %s\n", *status.CommitTag)
	}
	if status.UpdateAvailable != nil && *status.UpdateAvailable {
		commitsBehind := 0
		if status.CommitsBehind != nil {
			commitsBehind = int(*status.CommitsBehind)
		}
		fmt.Printf("  Update Available: Yes (%d commits behind)\n", commitsBehind)
	} else {
		fmt.Printf("  Update Available: No\n")
	}

	return nil
}
