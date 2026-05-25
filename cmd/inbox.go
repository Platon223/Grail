package cmd

import (
	"fmt"

	inboxInter "github.com/Platon223/Grail/internal/commands/inbox"
	"github.com/Platon223/Grail/internal/ui/inbox"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var inboxCmd = &cobra.Command{
	Use:   "inbox",
	Short: "inbox command launches the tui for the main inbox",
	RunE: func(cmd *cobra.Command, args []string) error {

		gmailService, err := inboxInter.GetMainInboxEmails()
		if err != nil {
			fmt.Println(err)
			return nil
		}	


		p := tea.NewProgram(
			inbox.NewInboxModel(gmailService),
			tea.WithAltScreen(),
		)

		_, err = p.Run()
		return err
	},
}

func init() {
	rootCmd.AddCommand(inboxCmd)
}
