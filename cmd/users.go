package cmd

import (
	"fmt"

	"github.com/julianfbeck/overseerr-cli/internal/api"
	"github.com/spf13/cobra"
)

var usersCmd = &cobra.Command{
	Use:     "users",
	Aliases: []string{"user"},
	Short:   "Manage users",
}

var usersListCmd = &cobra.Command{
	Use:   "list",
	Short: "List users",
	RunE:  runUsersList,
}

var usersMeCmd = &cobra.Command{
	Use:   "me",
	Short: "Get current user",
	RunE:  runUsersMe,
}

var (
	usersLimit int
	usersSkip  int
)

func init() {
	rootCmd.AddCommand(usersCmd)
	usersCmd.AddCommand(usersListCmd)
	usersCmd.AddCommand(usersMeCmd)

	usersListCmd.Flags().IntVarP(&usersLimit, "limit", "l", 20, "Number of users to show")
	usersListCmd.Flags().IntVarP(&usersSkip, "skip", "s", 0, "Number of users to skip")
}

func runUsersList(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	take := float32(usersLimit)
	skip := float32(usersSkip)

	resp, err := client.GetUserWithResponse(ctx, &api.GetUserParams{
		Take: &take,
		Skip: &skip,
	})
	if err != nil {
		return fmt.Errorf("failed to list users: %w", err)
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
		fmt.Println("No users found")
		return nil
	}

	total := 0
	if result.PageInfo != nil && result.PageInfo.Results != nil {
		total = int(*result.PageInfo.Results)
	}

	fmt.Printf("Users (showing %d of %d)\n\n", len(*result.Results), total)

	for _, user := range *result.Results {
		printUser(&user)
	}

	return nil
}

func printUser(u *api.User) {
	name := derefStr(u.Username)
	if name == "" {
		name = derefStr(u.Email)
	}

	id := 0
	if u.Id != nil {
		id = *u.Id
	}

	fmt.Printf("[%d] %s\n", id, name)
	if u.Email != nil && *u.Email != "" {
		fmt.Printf("  Email: %s\n", *u.Email)
	}
	if u.PlexUsername != nil && *u.PlexUsername != "" {
		fmt.Printf("  Plex: %s\n", *u.PlexUsername)
	}
	if u.RequestCount != nil {
		fmt.Printf("  Requests: %d\n", int(*u.RequestCount))
	}
	if u.CreatedAt != nil {
		// CreatedAt is a string in the API
		created := *u.CreatedAt
		if len(created) >= 10 {
			created = created[:10]
		}
		fmt.Printf("  Created: %s\n", created)
	}
	fmt.Println()
}

func runUsersMe(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	resp, err := client.GetAuthMeWithResponse(ctx)
	if err != nil {
		return fmt.Errorf("failed to get current user: %w", err)
	}

	if resp.JSON200 == nil {
		return fmt.Errorf("unexpected response: %s", resp.Status())
	}

	if jsonOutput {
		outputJSON(resp.JSON200)
		return nil
	}

	printUser(resp.JSON200)
	return nil
}
