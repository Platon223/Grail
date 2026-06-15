package cmd

import (
	"fmt"

	composeInter "github.com/Platon223/Grail/internal/commands/compose"
	"github.com/Platon223/Grail/internal/ui/compose"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)


var composeCmd = &cobra.Command{
	Use:   "compose",
	Short: "compose command launches a compose tui",
	RunE: func(cmd *cobra.Command, args []string) error {

		gmailService, err := composeInter.GetComposeSvc()
		if err != nil {
			fmt.Println(err)
			return nil
		}	


		p := tea.NewProgram(
			compose.NewComposeModel(gmailService),
			tea.WithAltScreen(),
		)

		_, err = p.Run()
		return err
	},
}

func init() {
	rootCmd.AddCommand(composeCmd)
}


