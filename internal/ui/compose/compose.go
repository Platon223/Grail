package compose

import (
	"fmt"
	"time"
	//	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	gmailapi "google.golang.org/api/gmail/v1"

	//	"github.com/Platon223/Grail/internal/domain/auth"
	//	"github.com/Platon223/Grail/internal/domain/email"
	domainGmail "github.com/Platon223/Grail/internal/domain/email"
	//	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
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

type clearInputMsg struct{}

type ComposeModel struct {
	svc *gmailapi.Service
	userEmail string
	composeStage int
	composeTo string
	composeSubject string
	composeBody string
	textArea textarea.Model
	textInputTo textinput.Model
	textInputSubject textinput.Model
	renderer *glamour.TermRenderer
	termWidth int
}


func NewComposeModel(gmailService *gmailapi.Service, userEmail string) ComposeModel {		


	r, _ := glamour.NewTermRenderer(
		glamour.WithStandardStyle("dark"),
		glamour.WithWordWrap(60),
	)

	ta := textarea.New()
	ta.Placeholder = "Write your message..."
	ta.SetWidth(60)
	ta.SetHeight(10)

	tito := textinput.New()

	tisub := textinput.New()


	return ComposeModel{
		svc: gmailService,
		userEmail: userEmail,
		renderer: r,
		textArea: ta,
		textInputTo: tito,
		textInputSubject: tisub,
	}
}


func (m ComposeModel) Init() tea.Cmd {
	return tea.Tick(200*time.Millisecond, func(t time.Time) tea.Msg {
        return clearInputMsg{}
    })
}


func (m ComposeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {

	
	case clearInputMsg:
		m.textInputTo.Focus()
		return m, nil

	case tea.WindowSizeMsg:
		m.termWidth = msg.Width

		m.renderer, _ = glamour.NewTermRenderer(
			glamour.WithStandardStyle("dark"),
			glamour.WithWordWrap(m.termWidth - 6),
		)
	
	case tea.KeyMsg:

		switch msg.String() {


		case "esc":

			if m.composeStage == 1 {
				m.textInputSubject.Blur()
				m.textInputTo.Focus()
				m.composeStage--
			} else if m.composeStage == 2 {
				m.textArea.Blur()
				m.textInputSubject.Focus()
				m.composeStage--
			} else if m.composeStage == 3 {
				m.textArea.Focus()
				m.composeStage--
			}

		case "enter":

			if m.composeStage == 2 {
				var cmd tea.Cmd
				m.textArea, cmd = m.textArea.Update(msg)
				return m, cmd
			}


			if m.composeStage == 0 {

				valueTo := m.textInputTo.Value()

				m.composeTo = valueTo

				m.textInputTo.Blur()
				m.textInputSubject.Focus()

			} else if m.composeStage == 1 {

				valueSubject := m.textInputSubject.Value()

				m.composeSubject = valueSubject

				m.textInputSubject.Blur()
				m.textArea.Focus()

			} else if m.composeStage == 2 {

				valueBody := m.textArea.Value()

				m.composeBody = valueBody
			}

			if m.composeStage != 2 {
				m.composeStage++
			}



		case "ctrl+d":
			if m.composeStage == 2 {
				m.composeBody = m.textArea.Value()
				m.composeStage++
			} else if m.composeStage == 3 {
				err := domainGmail.SendEmail(m.svc, m.composeTo, m.composeSubject, m.composeBody)
				if err != nil {
					fmt.Println("send error:", err)
				} else {
					fmt.Println("sent")
					return m, tea.Quit
				}
			}

		case "q":

			return m, tea.Quit

		default:

			
			var cmd tea.Cmd
			switch m.composeStage {
			case 0:
				m.textInputTo, cmd = m.textInputTo.Update(msg)
			case 1:
				m.textInputSubject, cmd = m.textInputSubject.Update(msg)
			case 2:
				m.textArea, cmd = m.textArea.Update(msg)
			}

			return m, cmd

			}
	}
	
	var cmd tea.Cmd

	return m, cmd

}

func (m ComposeModel) View() string {

	s := ""

	helper, err := m.renderer.Render("`enter`: next - `esc`: go back - `skip two lines`: make a new paragraph - `ctrl+d`: finish body and go to preview - `q`: quit")
	if err != nil {
		s += "Error rendering markdown content."
	} else {
		s += helper + "\n\n"
	}

	switch m.composeStage {
	
	case 0:
		s += fmt.Sprintf("To: %s", m.textInputTo.View())

	case 1:
		s += fmt.Sprintf("Subject: %s", m.textInputSubject.View())

	case 2:
		s += fmt.Sprintf("Body: %s", m.textArea.View())

	case 3:
		confirm, _ := m.renderer.Render("### Compose Preview")
		renderedBody, err := m.renderer.Render(m.composeBody)
		if err != nil {
			renderedBody = m.composeBody
		}
		s += fmt.Sprintf("%s\n\nTo: %s\nSubject: %s\n\n%s\n\n", confirm, m.composeTo, m.composeSubject, renderedBody)
		s += "Press ctrl+d to send, esc to go back."
	}

	title := titleStyle.Render(" Compose ")
	box := boxStyle.Render(s)
	return fmt.Sprintf("%s\n%s\n", title, box)
}
