package cmd

import (
	"fmt"
	"strconv"

	"github.com/julianfbeck/overseerr-cli/internal/api"
	"github.com/spf13/cobra"
)

var mediaCmd = &cobra.Command{
	Use:   "media",
	Short: "Get media details",
}

var movieCmd = &cobra.Command{
	Use:   "movie <tmdb-id>",
	Short: "Get movie details by TMDB ID",
	Args:  cobra.ExactArgs(1),
	RunE:  runMovie,
}

var tvCmd = &cobra.Command{
	Use:   "tv <tmdb-id>",
	Short: "Get TV show details by TMDB ID",
	Args:  cobra.ExactArgs(1),
	RunE:  runTV,
}

func init() {
	rootCmd.AddCommand(mediaCmd)
	mediaCmd.AddCommand(movieCmd)
	mediaCmd.AddCommand(tvCmd)
}

func runMovie(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid TMDB ID: %s", args[0])
	}

	resp, err := client.GetMovieMovieIdWithResponse(ctx, float32(id), nil)
	if err != nil {
		return fmt.Errorf("failed to get movie: %w", err)
	}

	if resp.JSON200 == nil {
		return fmt.Errorf("unexpected response: %s", resp.Status())
	}

	if jsonOutput {
		outputJSON(resp.JSON200)
		return nil
	}

	printMovieDetails(resp.JSON200)
	return nil
}

func printMovieDetails(m *api.MovieDetails) {
	title := derefStr(m.Title)
	date := derefStr(m.ReleaseDate)
	year := ""
	if len(date) >= 4 {
		year = date[:4]
	}

	fmt.Printf("%s (%s)\n", title, year)
	fmt.Printf("TMDB ID: %d\n", int(derefFloat(m.Id)))
	fmt.Printf("Status: %s\n", derefStr(m.Status))

	if m.Runtime != nil && *m.Runtime > 0 {
		fmt.Printf("Runtime: %d min\n", int(*m.Runtime))
	}

	fmt.Printf("Rating: %.1f/10", derefFloat(m.VoteAverage))
	if m.VoteCount != nil {
		fmt.Printf(" (%d votes)", int(*m.VoteCount))
	}
	fmt.Println()

	if m.Genres != nil && len(*m.Genres) > 0 {
		fmt.Printf("Genres: ")
		for i, g := range *m.Genres {
			if i > 0 {
				fmt.Print(", ")
			}
			fmt.Print(derefStr(g.Name))
		}
		fmt.Println()
	}

	if m.MediaInfo != nil && m.MediaInfo.Status != nil {
		fmt.Printf("Library Status: %s\n", api.StatusString(m.MediaInfo.Status))
	}

	fmt.Println()
	if m.Overview != nil && *m.Overview != "" {
		fmt.Println(*m.Overview)
	}
}

func runTV(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid TMDB ID: %s", args[0])
	}

	resp, err := client.GetTvTvIdWithResponse(ctx, float32(id), nil)
	if err != nil {
		return fmt.Errorf("failed to get TV show: %w", err)
	}

	if resp.JSON200 == nil {
		return fmt.Errorf("unexpected response: %s", resp.Status())
	}

	if jsonOutput {
		outputJSON(resp.JSON200)
		return nil
	}

	printTVDetails(resp.JSON200)
	return nil
}

func printTVDetails(t *api.TvDetails) {
	title := derefStr(t.Name)
	date := derefStr(t.FirstAirDate)
	year := ""
	if len(date) >= 4 {
		year = date[:4]
	}

	fmt.Printf("%s (%s)\n", title, year)
	fmt.Printf("TMDB ID: %d\n", int(derefFloat(t.Id)))
	fmt.Printf("Status: %s\n", derefStr(t.Status))

	seasons := 0
	episodes := 0
	if t.NumberOfSeason != nil {
		seasons = int(*t.NumberOfSeason)
	}
	if t.NumberOfEpisodes != nil {
		episodes = int(*t.NumberOfEpisodes)
	}
	fmt.Printf("Seasons: %d | Episodes: %d\n", seasons, episodes)

	fmt.Printf("Rating: %.1f/10", derefFloat(t.VoteAverage))
	if t.VoteCount != nil {
		fmt.Printf(" (%d votes)", int(*t.VoteCount))
	}
	fmt.Println()

	if t.Genres != nil && len(*t.Genres) > 0 {
		fmt.Printf("Genres: ")
		for i, g := range *t.Genres {
			if i > 0 {
				fmt.Print(", ")
			}
			fmt.Print(derefStr(g.Name))
		}
		fmt.Println()
	}

	if t.Networks != nil && len(*t.Networks) > 0 {
		fmt.Printf("Networks: ")
		for i, n := range *t.Networks {
			if i > 0 {
				fmt.Print(", ")
			}
			fmt.Print(derefStr(n.Name))
		}
		fmt.Println()
	}

	if t.MediaInfo != nil && t.MediaInfo.Status != nil {
		fmt.Printf("Library Status: %s\n", api.StatusString(t.MediaInfo.Status))
	}

	fmt.Println()
	if t.Overview != nil && *t.Overview != "" {
		fmt.Println(*t.Overview)
	}

	if t.Seasons != nil && len(*t.Seasons) > 0 {
		fmt.Println("\nSeasons:")
		for _, s := range *t.Seasons {
			seasonNum := 0
			if s.SeasonNumber != nil {
				seasonNum = int(*s.SeasonNumber)
			}
			if seasonNum == 0 {
				continue // Skip specials
			}
			airYear := ""
			airDate := derefStr(s.AirDate)
			if len(airDate) >= 4 {
				airYear = airDate[:4]
			}
			episodeCount := 0
			if s.EpisodeCount != nil {
				episodeCount = int(*s.EpisodeCount)
			}
			fmt.Printf("  Season %d: %d episodes (%s)\n", seasonNum, episodeCount, airYear)
		}
	}
}
