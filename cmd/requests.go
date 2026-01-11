package cmd

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/julianfbeck/overseerr-cli/internal/api"
	"github.com/spf13/cobra"
)

var requestsCmd = &cobra.Command{
	Use:     "requests",
	Aliases: []string{"request", "req"},
	Short:   "Manage media requests",
}

var requestsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List media requests",
	RunE:  runRequestsList,
}

var requestsGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get request details",
	Args:  cobra.ExactArgs(1),
	RunE:  runRequestsGet,
}

var requestsApproveCmd = &cobra.Command{
	Use:   "approve <id>",
	Short: "Approve a pending request",
	Args:  cobra.ExactArgs(1),
	RunE:  runRequestsApprove,
}

var requestsDeclineCmd = &cobra.Command{
	Use:   "decline <id>",
	Short: "Decline a pending request",
	Args:  cobra.ExactArgs(1),
	RunE:  runRequestsDecline,
}

var requestsDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a request",
	Args:  cobra.ExactArgs(1),
	RunE:  runRequestsDelete,
}

var requestsMovieCmd = &cobra.Command{
	Use:   "movie <tmdb-id>",
	Short: "Request a movie by TMDB ID",
	Args:  cobra.ExactArgs(1),
	RunE:  runRequestsMovie,
}

var requestsTVCmd = &cobra.Command{
	Use:   "tv <tmdb-id>",
	Short: "Request a TV show by TMDB ID",
	Args:  cobra.ExactArgs(1),
	RunE:  runRequestsTV,
}

var (
	requestsLimit  int
	requestsSkip   int
	requestsFilter string
	requestsSort   string
	tvSeasons      []int
	forceDelete    bool
)

func init() {
	rootCmd.AddCommand(requestsCmd)
	requestsCmd.AddCommand(requestsListCmd)
	requestsCmd.AddCommand(requestsGetCmd)
	requestsCmd.AddCommand(requestsApproveCmd)
	requestsCmd.AddCommand(requestsDeclineCmd)
	requestsCmd.AddCommand(requestsDeleteCmd)
	requestsCmd.AddCommand(requestsMovieCmd)
	requestsCmd.AddCommand(requestsTVCmd)

	requestsListCmd.Flags().IntVarP(&requestsLimit, "limit", "l", 20, "Number of requests to show")
	requestsListCmd.Flags().IntVarP(&requestsSkip, "skip", "s", 0, "Number of requests to skip")
	requestsListCmd.Flags().StringVarP(&requestsFilter, "filter", "f", "", "Filter: all, pending, approved, available, processing, unavailable")
	requestsListCmd.Flags().StringVar(&requestsSort, "sort", "", "Sort: added, modified")

	requestsTVCmd.Flags().IntSliceVar(&tvSeasons, "seasons", nil, "Specific seasons to request (default: all)")

	requestsDeleteCmd.Flags().BoolVar(&forceDelete, "force", false, "Skip confirmation")
}

func runRequestsList(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	take := float32(requestsLimit)
	skip := float32(requestsSkip)

	params := &api.GetRequestParams{
		Take: &take,
		Skip: &skip,
	}

	if requestsFilter != "" {
		filter := api.GetRequestParamsFilter(requestsFilter)
		params.Filter = &filter
	}
	if requestsSort != "" {
		sort := api.GetRequestParamsSort(requestsSort)
		params.Sort = &sort
	}

	resp, err := client.GetRequestWithResponse(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to list requests: %w", err)
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
		fmt.Println("No requests found")
		return nil
	}

	total := 0
	if result.PageInfo != nil && result.PageInfo.Results != nil {
		total = int(*result.PageInfo.Results)
	}

	fmt.Printf("Requests (showing %d of %d)\n\n", len(*result.Results), total)

	for _, req := range *result.Results {
		printRequest(&req)
	}

	return nil
}

func printRequest(req *api.MediaRequest) {
	status := api.RequestStatusString(req.Status)
	tmdbID := 0
	if req.Media != nil && req.Media.TmdbId != nil {
		tmdbID = int(*req.Media.TmdbId)
	}

	fmt.Printf("[%d] TMDB: %d - %s\n",
		int(derefFloat(req.Id)), tmdbID, status)

	if req.RequestedBy != nil {
		name := derefStr(req.RequestedBy.Username)
		if name == "" {
			name = derefStr(req.RequestedBy.Email)
		}
		fmt.Printf("  Requested by: %s\n", name)
	}

	if req.CreatedAt != nil {
		created := *req.CreatedAt
		if len(created) >= 16 {
			created = created[:16]
		}
		fmt.Printf("  Created: %s\n", created)
	}

	fmt.Println()
}

func runRequestsGet(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	resp, err := client.GetRequestRequestIdWithResponse(ctx, args[0])
	if err != nil {
		return fmt.Errorf("failed to get request: %w", err)
	}

	if resp.JSON200 == nil {
		return fmt.Errorf("unexpected response: %s", resp.Status())
	}

	if jsonOutput {
		outputJSON(resp.JSON200)
		return nil
	}

	printRequest(resp.JSON200)
	return nil
}

func runRequestsApprove(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	resp, err := client.PostRequestRequestIdStatusWithResponse(ctx, args[0], api.Approve)
	if err != nil {
		return fmt.Errorf("failed to approve request: %w", err)
	}

	if resp.JSON200 == nil {
		return fmt.Errorf("unexpected response: %s", resp.Status())
	}

	if jsonOutput {
		outputJSON(resp.JSON200)
		return nil
	}

	fmt.Printf("Request %s approved\n", args[0])
	return nil
}

func runRequestsDecline(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	resp, err := client.PostRequestRequestIdStatusWithResponse(ctx, args[0], api.Decline)
	if err != nil {
		return fmt.Errorf("failed to decline request: %w", err)
	}

	if resp.JSON200 == nil {
		return fmt.Errorf("unexpected response: %s", resp.Status())
	}

	if jsonOutput {
		outputJSON(resp.JSON200)
		return nil
	}

	fmt.Printf("Request %s declined\n", args[0])
	return nil
}

func runRequestsDelete(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	if !forceDelete && !quietMode {
		fmt.Printf("Delete request %s? This cannot be undone. Use --force to skip confirmation.\n", args[0])
		return nil
	}

	resp, err := client.DeleteRequestRequestIdWithResponse(ctx, args[0])
	if err != nil {
		return fmt.Errorf("failed to delete request: %w", err)
	}

	if resp.StatusCode() >= 400 {
		return fmt.Errorf("failed to delete: %s", resp.Status())
	}

	if !quietMode {
		fmt.Printf("Request %s deleted\n", args[0])
	}
	return nil
}

func runRequestsMovie(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	tmdbID, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid TMDB ID: %s", args[0])
	}

	mediaType := api.PostRequestJSONBodyMediaTypeMovie
	mediaID := float32(tmdbID)

	resp, err := client.PostRequestWithResponse(ctx, api.PostRequestJSONRequestBody{
		MediaType: mediaType,
		MediaId:   mediaID,
	})
	if err != nil {
		return fmt.Errorf("failed to request movie: %w", err)
	}

	if resp.JSON201 == nil {
		return fmt.Errorf("unexpected response: %s", resp.Status())
	}

	if jsonOutput {
		outputJSON(resp.JSON201)
		return nil
	}

	fmt.Printf("Movie requested successfully (Request ID: %d)\n", int(derefFloat(resp.JSON201.Id)))
	return nil
}

func runRequestsTV(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	tmdbID, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid TMDB ID: %s", args[0])
	}

	mediaType := api.PostRequestJSONBodyMediaTypeTv
	mediaID := float32(tmdbID)

	body := api.PostRequestJSONRequestBody{
		MediaType: mediaType,
		MediaId:   mediaID,
	}

	// Handle seasons - set via raw JSON if specified
	if len(tvSeasons) > 0 {
		seasons := make([]float32, len(tvSeasons))
		for i, s := range tvSeasons {
			seasons[i] = float32(s)
		}
		seasonsJSON, _ := json.Marshal(seasons)
		var seasonsUnion api.PostRequestJSONBody_Seasons
		_ = json.Unmarshal(seasonsJSON, &seasonsUnion)
		body.Seasons = &seasonsUnion
	}

	resp, err := client.PostRequestWithResponse(ctx, body)
	if err != nil {
		return fmt.Errorf("failed to request TV show: %w", err)
	}

	if resp.JSON201 == nil {
		return fmt.Errorf("unexpected response: %s", resp.Status())
	}

	if jsonOutput {
		outputJSON(resp.JSON201)
		return nil
	}

	fmt.Printf("TV show requested successfully (Request ID: %d)\n", int(derefFloat(resp.JSON201.Id)))
	return nil
}
