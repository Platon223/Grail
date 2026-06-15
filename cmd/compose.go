package cmd

import (
	"fmt"

	composeInter "github.com/Platon223/Grail/internal/commands/compose"
	"github.com/Platon223/Grail/internal/ui/compose"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	gmailapi "google.golang.org/api/gmail/v1"
	"github.com/charmbracelet/bubbles/spinner"
	"time"

)


var composeCmd = &cobra.Command{
	Use:   "compose",
	Short: "compose command launches a compose tui",
	RunE: func(cmd *cobra.Command, args []string) error {
		
		done := make(chan *gmailapi.Profile)
		errCh := make(chan error)

		gmailService, err := composeInter.GetComposeSvc()
		if err != nil {
			fmt.Println(err)
			return nil
		}	


		go func() {
			profile, err := gmailService.Users.GetProfile("me").Do()
			if err != nil {
				errCh <- err
				return
			}
			done <- profile
		}()

		s := spinner.New()
		s.Spinner = spinner.Dot

		fmt.Print("\033[?25l") // hide cursor
		defer fmt.Print("\033[?25h") // show cursor on exit

		ticker := time.NewTicker(80 * time.Millisecond)
		defer ticker.Stop()

		var profile *gmailapi.Profile

	loop:
		for {
			select {
			case p := <-done:
				profile = p
				break loop
			case err := <-errCh:
				fmt.Print("\r")
				return err
			case <-ticker.C:
				fmt.Printf("\r%s Connecting to Gmail...", s.View())
				var cmd2 tea.Cmd
				s, cmd2 = s.Update(spinner.TickMsg{})
				_ = cmd2
			}
		}

		fmt.Print("\r\033[K") // clear the line

		p := tea.NewProgram(
			compose.NewComposeModel(gmailService, profile.EmailAddress),
			tea.WithAltScreen(),
		)
		_, err = p.Run()
		return err

	},
}

func init() {
	rootCmd.AddCommand(composeCmd)
}


