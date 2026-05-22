package cmd

import (
	"github.com/Platon223/Grail/internal/ui/inbox"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var inboxCmd = &cobra.Command{
	Use:   "inbox",
	Short: "inbox command launches the tui for the main inbox",
	RunE: func(cmd *cobra.Command, args []string) error {
		p := tea.NewProgram(
			inbox.NewInboxModel(),
			tea.WithAltScreen(),
		)

		_, err := p.Run()
		return err
	},
}

func init() {
	rootCmd.AddCommand(inboxCmd)
}
