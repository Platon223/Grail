package inbox

import (
	//	"github.com/charmbracelet/bubbles/list"
	//	"github.com/charmbracelet/bubbles/viewport"
	"fmt"

//	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	// "github.com/charmbracelet/glamour"
	// "github.com/charmbracelet/lipgloss"
)


type state int

const (
	stateInbox state = iota
	stateReading
)

type Email struct {
	Body string
	From string
	To string
}

type InboxModel struct {
	emails []Email
	cursor int
	viewport int
	scroll int
}


func NewInboxModel() InboxModel {

	return InboxModel{
		emails: []Email{
			{
				Body: "some email",
				From: "someuser@gmail.com",
				To: "someuser@gmail.com",
			},
			{
				Body: "some email",
				From: "someuser@gmail.com",
				To: "someuser@gmail.com",
			},
			{
				Body: "some email",
				From: "someuser@gmail.com",
				To: "someuser@gmail.com",
			},
			{
				Body: "some email",
				From: "someuser@gmail.com",
				To: "someuser@gmail.com",
			},
			{
				Body: "some email",
				From: "someuser@gmail.com",
				To: "someuser@gmail.com",
			},
			{
				Body: "some email",
				From: "someuser@gmail.com",
				To: "someuser@gmail.com",
			},
		},
		cursor: 0,
		viewport: 5,
		scroll: 0,
	}
}

func (m InboxModel) Init() tea.Cmd {
	return nil
}

func (m InboxModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {

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
		}
	}

	var cmd tea.Cmd

	return m, cmd
}



func (m InboxModel) View() string {

	s := "Use up/down or j/k to scroll. Press q or ctrl+c to exit \n\n"


	end := m.scroll + m.viewport
	if end > len(m.emails) {
		end = len(m.emails)
	}

	for i := m.scroll; i < end; i++ {
		cursor := " "

		if m.cursor == i {
			cursor = ">"
		}

		s += fmt.Sprintf("%s %s\n", cursor, m.emails[i])

	}

	s += fmt.Sprintf("\n(Viewing %d-%d of %d)", m.scroll + 1, end, len(m.emails))

	return s
}

