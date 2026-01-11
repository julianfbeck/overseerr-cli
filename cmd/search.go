package cmd

import (
	"fmt"
	"strings"

	"github.com/julianfbeck/overseerr-cli/internal/api"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search for movies and TV shows",
	Args:  cobra.MinimumNArgs(1),
	RunE:  runSearch,
}

var searchPage int

func init() {
	rootCmd.AddCommand(searchCmd)
	searchCmd.Flags().IntVarP(&searchPage, "page", "p", 1, "Page number")
}

func runSearch(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	query := strings.Join(args, " ")
	page := float32(searchPage)

	resp, err := client.GetSearchWithResponse(ctx, &api.GetSearchParams{
		Query: query,
		Page:  &page,
	})
	if err != nil {
		return fmt.Errorf("search failed: %w", err)
	}

	if resp.JSON200 == nil {
		return fmt.Errorf("unexpected response: %s", resp.Status())
	}

	result := resp.JSON200

	if jsonOutput {
		outputJSON(result)
		return nil
	}

	if result.Results == nil || len(*result.Results) == 0 {
		fmt.Println("No results found")
		return nil
	}

	totalPages := 1
	if result.TotalPages != nil {
		totalPages = int(*result.TotalPages)
	}
	totalResults := 0
	if result.TotalResults != nil {
		totalResults = int(*result.TotalResults)
	}

	fmt.Printf("Search results for '%s' (page %d/%d, %d total)\n\n",
		query, searchPage, totalPages, totalResults)

	// Search results use union types - output JSON for detailed info
	fmt.Println("Use --json flag for detailed search results")

	return nil
}
