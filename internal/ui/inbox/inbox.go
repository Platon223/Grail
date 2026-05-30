package inbox

import (
	//	"github.com/charmbracelet/bubbles/list"
	//	"github.com/charmbracelet/bubbles/viewport"
	"fmt"
	"strings"

	//	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	gmailapi "google.golang.org/api/gmail/v1"

	"github.com/Platon223/Grail/internal/domain/auth"
	domainGmail "github.com/Platon223/Grail/internal/domain/email"
	"github.com/charmbracelet/lipgloss"
	"github.com/atotto/clipboard"
)


var (
	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#6e6a86")).
			Padding(1, 2)

	titleStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#31748f")).
			Foreground(lipgloss.Color("230")).
			Padding(0, 1).
			Bold(true)

	cursorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#eb6f92")).
			Bold(true)
)

type InboxModel struct {
	emails []*domainGmail.Email
	userEmail string
	cursor int
	viewport int
	scroll int
	bodyScroll int
	bodyLine int
	item *domainGmail.Email
	replying bool
	replyStage int
	replyTo string
	replySubject string
	replyBody string
	renderer *glamour.TermRenderer
	termWidth int
}


func NewInboxModel(gmailService *gmailapi.Service) InboxModel {		

	res, err := gmailService.Users.Messages.List("me").
		LabelIds("INBOX").
		MaxResults(50).
		Do()

	profile, err := gmailService.Users.GetProfile("me").Do()

	if err != nil {
		fmt.Println(err)
		return InboxModel{}
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


	r, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(60),
	)

	return InboxModel{
		emails: userEmails,
		userEmail: profile.EmailAddress,
		cursor: 0,
		viewport: 100,
		scroll: 0,
		item: &domainGmail.Email{},
		replying: false,
		renderer: r,
	}
}

func (m InboxModel) Init() tea.Cmd {
	return nil
}

func (m InboxModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.termWidth = msg.Width

		m.renderer, _ = glamour.NewTermRenderer(
			glamour.WithAutoStyle(),
			glamour.WithWordWrap(m.termWidth - 6),
		)

	case tea.KeyMsg:
		switch msg.String() {


		case "o":
            if m.item.ID != "" && !m.replying {

                base := "https://mail.google.com/mail/u/0/"

                if m.userEmail != "" {

                    ae := strings.ReplaceAll(m.userEmail, "@", "%40")
                    ae = strings.ReplaceAll(ae, "+", "%2B")
                    ae = strings.ReplaceAll(ae, "&", "%26")
                    ae = strings.ReplaceAll(ae, "?", "%3F")
                    ae = strings.ReplaceAll(ae, "#", "%23")

                    url := base + "?authuser=" + ae + "#all/" + m.item.ID
                    auth.OpenBrowser(url)
                } else {
                    auth.OpenBrowser(base + "#all/" + m.item.ID)
                }
            }
		case "backspace", "esc":
			if m.replying {
				m.replying = false
				m.replyStage = 0
			} else {
				m.item = &domainGmail.Email{}
			}


		case "q", "ctrl+c":
			if m.item.ID == "" {
				return m, tea.Quit
			}

		case "k", "up":	
			// scrolling in the email view
			if m.item.ID != "" {
				if m.bodyLine > 0 {
					m.bodyLine--
				} else if m.bodyScroll > 0 {
					m.bodyScroll--
				}
			} else {
				// scrolling in the inbox view

				if m.cursor > 0 {
					m.cursor --
				}

				if m.cursor < m.scroll {
					m.scroll = m.cursor
				}
			}

		case "j", "down":
			// scrolling in the email view
			if m.item.ID != "" {

				visibleLines := 50

				bodyRendered, err := m.renderer.Render(m.item.Body)
				
				var lines []string

				if err != nil {
					lines = strings.Split(m.item.Body, "\n")
				} else {
					lines = strings.Split(bodyRendered, "\n")
				}

				totalLines := len(lines)

				if m.bodyLine < visibleLines-1 && m.bodyScroll+m.bodyLine < totalLines-1 {
					m.bodyLine++
				} else if m.bodyScroll+m.bodyLine < totalLines-1 {
					// cursor is at bottom of window, scroll down instead
					m.bodyScroll++
				}

			} else {
				// scrolling in the inbox view
				if m.cursor < len(m.emails)-1 {
					m.cursor++
				}

				if m.cursor >= m.scroll+m.viewport {
					m.scroll = m.cursor - m.viewport + 1
				}
			}

		case "enter":

if m.replying {
    // NOTE: to fully enable interactive text input you should:
    // 1) import the textinput package in this file (example, commented out):
    //    // import textinput "github.com/charmbracelet/bubbles/textinput"
    // 2) add a TextInput field to the InboxModel struct, for example (commented out):
    //    // TextInput textinput.Model // used to capture typed input while replying
    // 3) initialize TextInput when entering reply mode (set cursor styles, placeholders, etc.)

    // When TextInput is wired up, forward other keypresses (backspace, regular characters, etc.)
    // to the text input model from the top-level key handling. Example (commented):
    //
    // if m.replying {
    //     var cmd tea.Cmd
    //     m.TextInput, cmd = m.TextInput.Update(msg)
    //     return m, cmd
    // }

    // The block below shows how Enter should behave once the textinput is available.
    // It's commented out so you can enable it after adding the struct field + import.
    /*
    if m.replyStage < 3 {
        // accept the current value from the textinput and store it in the model
        entered := m.TextInput.Value()
        switch m.replyStage {
        case 0:
            m.replyTo = entered
        case 1:
            m.replySubject = entered
        case 2:
            m.replyBody = entered
        }

        // reset the textinput value for the next field and advance stage
        m.TextInput.SetValue("")
        m.replyStage++
        // optionally set a new placeholder or focus state on m.TextInput here
    } else {
        // replyStage == 3: this is the preview/confirm stage. Here you can send the reply.
        // e.g. go send the message, or call a method like: m.sendReply()
    }
    */

    // Fallback behavior while textinput isn't wired up: advance stages as before
    if m.replyStage < 3 {
        m.replyStage++
    } else {
        // At final stage (preview) Enter could mean "send" — implement when ready.
    }
}

			if m.item.ID == "" {
				m.item = m.emails[m.cursor]	
				m.bodyScroll = 0
				m.bodyLine = 0
			}

		case "r":

			if m.item.ID != "" {
				m.replying = true	
			}
			
		

		case "y":
				
			var lines []string

			lines = strings.Split(m.item.Body, "\n")

			currentLine := lines[m.bodyScroll + m.bodyLine]

			err := clipboard.WriteAll(currentLine)
			if err != nil {	
			}
		}

	}

	var cmd tea.Cmd

	return m, cmd
}

func (m InboxModel) View() string {

	if m.replying {
		s := ""

		helper, err := m.renderer.Render("`enter`: next - `esc/backspace`: cancel")
		if err != nil {
			s += "Error rendering markdown content."
		} else {
			s += helper + "\n\n"
		}
	 

		toCursor := " "
		subjCursor := " "
		bodyCursor := " "
		if m.replyStage == 0 {
			toCursor = cursorStyle.Render(">")
		} else if m.replyStage == 1 {
			subjCursor = cursorStyle.Render(">")
		} else if m.replyStage == 2 {
			bodyCursor = cursorStyle.Render(">")
		}

		s += fmt.Sprintf("%s To: %s\n", toCursor, m.replyTo)
		s += fmt.Sprintf("%s Subject: %s\n", subjCursor, m.replySubject)
		s += fmt.Sprintf("%s Body: %s\n\n", bodyCursor, m.replyBody)

		if m.replyStage == 3 {
			confirm, _ := m.renderer.Render("`Reply Preview`")
			s += fmt.Sprintf("%s\n\nTo: %s\nSubject: %s\n\n%s\n", confirm, m.replyTo, m.replySubject, m.replyBody)
		}

		title := titleStyle.Render(" Reply ")
		box := boxStyle.Render(s)
		return fmt.Sprintf("%s\n%s\n", title, box)
	}

	if m.item.ID != "" {
		s := ""
		bodyRendered, err := m.renderer.Render(m.item.Body)
		helper, err := m.renderer.Render("`j/k`: scroll - `r`: reply - `o`: browser - `esc/backspace`: exit - `y`: Copy the current line")

		if err != nil {
			s += "Error rendering markdown content."
		} else {

			lines := strings.Split(bodyRendered, "\n")

			selectedIdx := m.bodyScroll + m.bodyLine
			if selectedIdx < 0 {
				selectedIdx = 0
			}
			if selectedIdx >= len(lines) {
				selectedIdx = len(lines) - 1
			}

			if m.bodyScroll + m.bodyLine < len(lines) {	
				lines[m.bodyScroll + m.bodyLine] = fmt.Sprintf("%s %s", ">", lines[m.bodyScroll + m.bodyLine])
			}

			visibleLines := 50
			end := m.bodyScroll + visibleLines

			if end > len(lines) {
				end = len(lines)
			}

			start := m.bodyScroll
			if start > len(lines) {
				end = len(lines)
			}

			visibleBody := strings.Join(lines[start:end], "\n")

			s += fmt.Sprintf("%s \n\nSubject: %s \nFrom: %s \n\n %s",
				helper,
				m.item.Subject,
				m.item.From,
				visibleBody,
			)
		}
 
	 
		title := titleStyle.Render(" Email View ")
		box := boxStyle.Render(s)
		return fmt.Sprintf("%s\n%s\n", title, box) 
	}

	s := ""
	helper , err := m.renderer.Render("`j/k`: scroll - `enter`: pick - `q/ctrl+c`: exit")
	if err != nil {
		s += "Error rendering markdown content."
	} else {
		s += helper
		end := m.scroll + m.viewport
		if end > len(m.emails) {
			end = len(m.emails)
		}

		for i := m.scroll; i < end; i++ {
			cursor := " "

			if m.cursor == i {
				cursor = cursorStyle.Render(">")
			}

			s += fmt.Sprintf("%s %s - %s\n", cursor, m.emails[i].From, m.emails[i].Subject)

		}
	

	}

	title := titleStyle.Render(" Inbox ")
	box := boxStyle.Render(s)

	return fmt.Sprintf("%s\n%s\n", title, box)
}

