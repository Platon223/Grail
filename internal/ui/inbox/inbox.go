package inbox

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	gmailapi "google.golang.org/api/gmail/v1"

	"github.com/Platon223/Grail/internal/domain/auth"
	"github.com/Platon223/Grail/internal/domain/email"
	domainGmail "github.com/Platon223/Grail/internal/domain/email"
	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/lipgloss"
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
	svc *gmailapi.Service
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
	textInput textarea.Model
	renderer *glamour.TermRenderer
	termWidth int
	termHeight int
}


func NewInboxModel(gmailService *gmailapi.Service, inboxType string, userEmails []*domainGmail.Email, userEmail string) InboxModel {		


	r, _ := glamour.NewTermRenderer(
		glamour.WithStandardStyle("dark"),
		glamour.WithWordWrap(60),
	)

	ta := textarea.New()
	ta.Placeholder = "Write your reply..."
	ta.SetWidth(60)
	ta.SetHeight(10)


	return InboxModel{
		svc: gmailService,
		emails: userEmails,
		userEmail: userEmail,
		cursor: 0,
		viewport: 100,
		scroll: 0,
		item: &domainGmail.Email{},
		replying: false,
		renderer: r,
		textInput: ta,
	}
}

func (m InboxModel) Init() tea.Cmd {
	return nil
}

func (m InboxModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.termWidth = msg.Width
		m.termHeight = msg.Height
		m.viewport = msg.Height - 8
		if m.viewport < 1 {
			m.viewport = 1
		}

		m.renderer, _ = glamour.NewTermRenderer(
			glamour.WithStandardStyle("dark"),
			glamour.WithWordWrap(m.termWidth - 6),
		)

	case tea.KeyMsg:		

		if m.replying {
			switch msg.String() {
			case "esc":
				m.replying = false
				m.replyStage = 0
			case "enter":
				switch m.replyStage {
				case 2:

					var cmd tea.Cmd
					m.textInput, cmd = m.textInput.Update(msg)
					return m, cmd

				case 3:
					m.replyStage--

				}
				m.textInput.SetValue("")
				m.replyStage++


			case "ctrl+d":
				if m.replyStage == 2 {

					entered := m.textInput.Value()
						
					m.replyBody = entered
					m.textInput.SetValue("")
					m.replyStage++
				} else if m.replyStage == 3 {
					err := email.SendReply(
						m.svc,
						m.replyTo,
						m.replySubject,
						m.replyBody,
						m.item.ThreadID,
						m.item.MessageID,
					)

					if err != nil {
						fmt.Println("send error: ", err)
					}

					m.replying = false
					m.replyStage = 0

				}

			default:
				var cmd tea.Cmd
				m.textInput, cmd = m.textInput.Update(msg)
				return m, cmd
			}
			return m, nil
		}



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

				visibleLines := m.termHeight - 10 
				if visibleLines < 1 {
					visibleLines = 1
				}

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

			if m.item.ID == "" {
				m.item = m.emails[m.cursor]	
				m.bodyScroll = 0
				m.bodyLine = 0
			}

		case "r":

			if m.item.ID != "" {
				m.replying = true	
				m.replyStage = 0
				m.replyTo = m.item.From
				m.replySubject = "Re: " + m.item.Subject
				m.replyBody = ""
				m.textInput.SetValue(m.replyTo)
				m.textInput.Focus()
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

		helper, err := m.renderer.Render("`enter`: next - `esc`: cancel - `skip two lines`: make a new paragraph - `ctrl+d`: finish body and go to preview")
		if err != nil {
			s += "Error rendering markdown content."
		} else {
			s += helper + "\n\n"
		}
	 

		if m.replyStage == 0 {
			s += fmt.Sprintf("  To: %s\n", m.replyTo)
		} else if m.replyStage == 1 {
			s += fmt.Sprintf("  Subject: %s\n", m.replySubject)
		} 


		if m.replyStage == 2 {
			s += fmt.Sprintf("%s Body: %s\n", cursorStyle.Render(">"), m.textInput.View())
		} else if m.replyStage == 3 {
			confirm, _ := m.renderer.Render("### Reply Preview")
			renderedBody, err := m.renderer.Render(m.replyBody)
			if err != nil {
				renderedBody = m.replyBody
			}
			s += fmt.Sprintf("%s\n\nTo: %s\nSubject: %s\n\n%s\n\n", confirm, m.replyTo, m.replySubject, renderedBody)
			s += "Press ctrl+d to send, esc to go back."
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

			visibleLines := m.termHeight - 10
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

