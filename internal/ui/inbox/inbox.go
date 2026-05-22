package inbox

import (
	//	"github.com/charmbracelet/bubbles/list"
	//	"github.com/charmbracelet/bubbles/viewport"
	"fmt"

	//	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"

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


type Email struct {
	Header string
	Name string
	Body string
	From string
	To string
}

type InboxModel struct {
	emails []Email
	cursor int
	viewport int
	scroll int
	item Email
	renderer *glamour.TermRenderer
	termWidth int
}


func NewInboxModel() InboxModel {

	r, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(60),
	)

	return InboxModel{
		emails: []Email{
			{
				Header: "Header123",
				Name: "Company",
				Body: "some email 1",
				From: "someuser@gmail.comhfffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffjfjfjfjfjfjfjjfjfjfjfjffjj",
				To: "someuser@gmail.com",
			},
			{
				Header: "Header123",
				Name: "Company",
				Body: "some email 2",
				From: "someuser@gmail.com",
				To: "someuser@gmail.com",
			},
			{
				Header: "Header123",
				Name: "Company",
				Body: "some email 3",
				From: "someuser@gmail.com",
				To: "someuser@gmail.com",
			},
			{
				Header: "Header123",
				Name: "Company",
				Body: "some email 4",
				From: "someuser@gmail.com",
				To: "someuser@gmail.com",
			},
			{
				Header: "Header123",
				Name: "Company",
				Body: "some email 5",
				From: "someuser@gmail.com",
				To: "someuser@gmail.com",
			},
			{
				Header: "Header123",
				Name: "Company",
				Body: "some email 6",
				From: "someuser@gmail.com",
				To: "someuser@gmail.com",
			},
		},
		cursor: 0,
		viewport: 10,
		scroll: 0,
		item: Email{},
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
			m.item = Email{}


		case "q", "ctrl+c":
			return m, tea.Quit

		case "k", "up":	
			if m.cursor > 0 {
				m.cursor --
			}

			if m.cursor < m.scroll {
				m.scroll = m.cursor
			}
		case "j", "down":
			if m.cursor < len(m.emails)-1 {
				m.cursor++
			}

			if m.cursor >= m.scroll+m.viewport {
				m.scroll = m.cursor - m.viewport + 1
			}

		case "enter":

			if m.item == (Email{}) {
				m.item = m.emails[m.cursor]	
			}
			
		}
	}

	var cmd tea.Cmd

	return m, cmd
}



func (m InboxModel) View() string {

	if m.item != (Email{}) {
		s := fmt.Sprintf("j/k scroll - esc/backspace exit - r reply - o browser \n\nHeader: %s \nFrom: %s \n\nBody: %s", m.item.Header, m.item.From, m.item.Body)
 
		title := titleStyle.Render(" Email View ")
		box := boxStyle.Render(s)
		return fmt.Sprintf("%s\n%s\n", title, box) 
	}

	s := "j/k scroll - enter pick - q/ctrl+c exit \n\n"


	end := m.scroll + m.viewport
	if end > len(m.emails) {
		end = len(m.emails)
	}

	for i := m.scroll; i < end; i++ {
		cursor := " "

		if m.cursor == i {
			cursor = cursorStyle.Render(">")
		}

		s += fmt.Sprintf("%s %s - %s - %s\n", cursor, m.emails[i].Name, m.emails[i].Header, m.emails[i].From)

	}

	title := titleStyle.Render(" Inbox ")
	box := boxStyle.Render(s)

	return fmt.Sprintf("%s\n%s\n", title, box)
}

