package cmd

import (
	"fmt"

	"github.com/julianfbeck/overseerr-cli/internal/api"
	"github.com/spf13/cobra"
)

// Unused import guard
var _ = api.StatusString

var discoverCmd = &cobra.Command{
	Use:   "discover",
	Short: "Discover movies and TV shows",
}

var discoverMoviesCmd = &cobra.Command{
	Use:   "movies",
	Short: "Discover popular movies",
	RunE:  runDiscoverMovies,
}

var discoverTVCmd = &cobra.Command{
	Use:   "tv",
	Short: "Discover popular TV shows",
	RunE:  runDiscoverTV,
}

var discoverTrendingCmd = &cobra.Command{
	Use:   "trending",
	Short: "Show trending movies and TV shows",
	RunE:  runDiscoverTrending,
}

var discoverPage int

func init() {
	rootCmd.AddCommand(discoverCmd)
	discoverCmd.AddCommand(discoverMoviesCmd)
	discoverCmd.AddCommand(discoverTVCmd)
	discoverCmd.AddCommand(discoverTrendingCmd)

	discoverCmd.PersistentFlags().IntVarP(&discoverPage, "page", "p", 1, "Page number")
}

func runDiscoverMovies(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	page := float32(discoverPage)
	resp, err := client.GetDiscoverMoviesWithResponse(ctx, &api.GetDiscoverMoviesParams{
		Page: &page,
	})
	if err != nil {
		return fmt.Errorf("failed to discover movies: %w", err)
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
		fmt.Println("No movies found")
		return nil
	}

	totalPages := 1
	if result.TotalPages != nil {
		totalPages = int(*result.TotalPages)
	}

	fmt.Printf("Popular Movies (page %d/%d)\n\n", discoverPage, totalPages)

	for _, item := range *result.Results {
		printMovieResult(&item)
	}

	return nil
}

func printMovieResult(m *api.MovieResult) {
	// MovieResult has non-pointer Id and Title
	title := m.Title
	date := derefStr(m.ReleaseDate)
	if len(date) >= 4 {
		date = date[:4]
	}
	status := ""
	if m.MediaInfo != nil && m.MediaInfo.Status != nil {
		status = fmt.Sprintf(" [%s]", api.StatusString(m.MediaInfo.Status))
	}
	fmt.Printf("[Movie] %s (%s) - TMDB ID: %d%s\n", title, date, int(m.Id), status)
	if m.Overview != nil && *m.Overview != "" {
		overview := *m.Overview
		if len(overview) > 150 {
			overview = overview[:147] + "..."
		}
		fmt.Printf("  %s\n", overview)
	}
	fmt.Println()
}

func runDiscoverTV(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	page := float32(discoverPage)
	resp, err := client.GetDiscoverTvWithResponse(ctx, &api.GetDiscoverTvParams{
		Page: &page,
	})
	if err != nil {
		return fmt.Errorf("failed to discover TV shows: %w", err)
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
		fmt.Println("No TV shows found")
		return nil
	}

	totalPages := 1
	if result.TotalPages != nil {
		totalPages = int(*result.TotalPages)
	}

	fmt.Printf("Popular TV Shows (page %d/%d)\n\n", discoverPage, totalPages)

	for _, item := range *result.Results {
		printTVResult(&item)
	}

	return nil
}

func printTVResult(t *api.TvResult) {
	title := derefStr(t.Name)
	date := derefStr(t.FirstAirDate)
	if len(date) >= 4 {
		date = date[:4]
	}
	status := ""
	if t.MediaInfo != nil && t.MediaInfo.Status != nil {
		status = fmt.Sprintf(" [%s]", api.StatusString(t.MediaInfo.Status))
	}
	fmt.Printf("[TV] %s (%s) - TMDB ID: %d%s\n", title, date, int(derefFloat(t.Id)), status)
	if t.Overview != nil && *t.Overview != "" {
		overview := *t.Overview
		if len(overview) > 150 {
			overview = overview[:147] + "..."
		}
		fmt.Printf("  %s\n", overview)
	}
	fmt.Println()
}

func runDiscoverTrending(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	page := float32(discoverPage)
	resp, err := client.GetDiscoverTrendingWithResponse(ctx, &api.GetDiscoverTrendingParams{
		Page: &page,
	})
	if err != nil {
		return fmt.Errorf("failed to get trending: %w", err)
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
		fmt.Println("No trending content found")
		return nil
	}

	totalPages := 1
	if result.TotalPages != nil {
		totalPages = int(*result.TotalPages)
	}

	fmt.Printf("Trending (page %d/%d)\n\n", discoverPage, totalPages)

	// Trending uses a complex union type - output JSON for now
	fmt.Println("Use --json flag for detailed trending results")

	return nil
}
