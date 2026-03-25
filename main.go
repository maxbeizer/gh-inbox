package main

import (
	"fmt"
	"log"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/spf13/cobra"

	"github.com/maxbeizer/gh-inbox/internal/api"
	"github.com/maxbeizer/gh-inbox/internal/tui"
)

func main() {
	userMessages := log.New(os.Stderr, "", 0)

	rootCmd := &cobra.Command{
		Use:   "gh-inbox",
		Short: "A rich TUI for managing GitHub notifications",
		Long: `gh-inbox is a terminal user interface for managing your GitHub notifications.
Browse, triage, and act on notifications with keyboard shortcuts and a 
beautiful visual interface inspired by gh-dash.`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := api.NewClient()
			if err != nil {
				return fmt.Errorf("failed to initialize GitHub client: %w\nMake sure you are authenticated with 'gh auth login'", err)
			}

			app := tui.NewApp(client)
			p := tea.NewProgram(app)
			if _, err := p.Run(); err != nil {
				return fmt.Errorf("error running TUI: %w", err)
			}
			return nil
		},
	}

	if err := rootCmd.Execute(); err != nil {
		userMessages.Printf("error: %v", err)
		os.Exit(1)
	}
}
