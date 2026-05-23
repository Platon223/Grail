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

	domainGmail "github.com/Platon223/Grail/internal/domain/email"
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
	emails []*domainGmail.Email
	cursor int
	viewport int
	scroll int
	bodyScroll int
	item *domainGmail.Email
	replying bool
	renderer *glamour.TermRenderer
	termWidth int
}


func NewInboxModel(gmailService *gmailapi.Service) InboxModel {		

	res, err := gmailService.Users.Messages.List("me").
		LabelIds("INBOX").
		MaxResults(50).
		Do()

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

		case "backspace", "esc":
			if m.replying {
				m.replying = false
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
				if m.bodyScroll > 0 {
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
				lines := strings.Split(m.item.Body, "\n")
				maxScroll := len(lines) - 50

				if maxScroll < 0 {
					maxScroll = 0
				}

				if m.bodyScroll < maxScroll {
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
			}

		case "r":

			if m.item.ID != "" {
				m.replying = true	
			}
			
		}
	}

	var cmd tea.Cmd

	return m, cmd
}



func (m InboxModel) View() string {

	if m.replying {
		s := "replying"
		title := titleStyle.Render(" Reply ")
		box := boxStyle.Render(s)
		
		return fmt.Sprintf("%s\n%s\n", title, box) 
	}

	if m.item.ID != "" {
		s := ""
		helper, err := m.renderer.Render("`j/k`: scroll - `r`: reply - `o`: browser - `esc/backspace`: exit")
		bodyRendered, err := m.renderer.Render(m.item.Body)

		if err != nil {
			s += "Error rendering markdown content."
		} else {

			lines := strings.Split(bodyRendered, "\n")
			lines[m.bodyScroll] = fmt.Sprintf("%s %s", ">", lines[m.bodyScroll])
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

