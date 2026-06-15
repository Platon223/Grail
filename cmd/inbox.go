package cmd

import (
	"fmt"

	inboxInter "github.com/Platon223/Grail/internal/commands/inbox"
	"github.com/Platon223/Grail/internal/ui/inbox"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	gmailapi "google.golang.org/api/gmail/v1"
	"time"
	"github.com/charmbracelet/bubbles/spinner"
	domainGmail "github.com/Platon223/Grail/internal/domain/email"

)

var inboxType string

var inboxCmd = &cobra.Command{
	Use:   "inbox",
	Short: "inbox command launches the tui for the desired inbox",
	RunE: func(cmd *cobra.Command, args []string) error {
		
		
		done := make(chan *gmailapi.Profile)
		var userEmailsGlobal []*domainGmail.Email
		errCh := make(chan error)
	

		gmailService, err := inboxInter.GetMainInboxEmails()
		if err != nil {
			fmt.Println(err)
			return nil
		}


		go func() {

			res, err := gmailService.Users.Messages.List("me").
				LabelIds(inboxType).
				MaxResults(50).
				Do()
	

			profile, err := gmailService.Users.GetProfile("me").Do()
			if err != nil {
				errCh <- err
				return
			}

			var userEmails []*domainGmail.Email

			for _, msg := range res.Messages {
				msg, err := gmailService.Users.Messages.Get("me", msg.Id).
					Format("full").
					Do()

				if err != nil {
					continue
				}

				email, err := domainGmail.FromMessage(msg)
				if err != nil {
					continue
				}

				userEmails = append(userEmails, email)
			}

	
			userEmailsGlobal = userEmails
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
			inbox.NewInboxModel(gmailService, inboxType, userEmailsGlobal, profile.EmailAddress),
			tea.WithAltScreen(),
		)

		_, err = p.Run()
		return err
	},
}

func init() {
	rootCmd.AddCommand(inboxCmd)

	inboxCmd.Flags().StringVarP(&inboxType, "inbox", "i", "INBOX", "Type of inbox")

}
