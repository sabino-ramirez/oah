/*
Copyright © 2022 Sabino Ramirez <sabinoramirez017@gmail.com>
*/
package setup

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sabino-ramirez/oah/data"

	"github.com/spf13/cobra"
)

// app state variables will have this type
type sessionState uint

// constants for keeping track of app state
const (
	inputView sessionState = iota
	promptView
)

// lipgloss styles
var (
	choiceStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("254"))
	selectedChoiceStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("5"))
	focusedModelStyle   = lipgloss.NewStyle().Padding(2).BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("69"))
	helpStyle           = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
)

type mainModel struct {
	state     sessionState
	TextInput textinput.Model

	params    []string
	currParam int

	choice int

	height int
	width  int

	err error
}

// tea message type for handling errors throughout program
type errMsg struct{ err error }

// in order to get errMsg type to implement error interface
func (e errMsg) Error() string { return e.err.Error() }

// returns what the initial model state will be
func initialModel() *mainModel {
	ti := textinput.New()
	ti.Placeholder = "copy/paste or type.."
	ti.Focus()
	ti.Width = 20

	params := []string{"auth", "orgId", "projTempId"}
	m := mainModel{state: inputView, TextInput: ti, params: params, currParam: 0, err: nil}
	return &m
}

// tea init function
func (m *mainModel) Init() tea.Cmd {
	//*TODO* cannot get the cursor to blink
	return tea.Batch(textinput.Blink, checkDatabase)
}

func (m *mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case errMsg:
		m.err = msg

	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			if m.state == promptView {
				m.choice += 1
				if m.choice > 1 {
					m.choice = 1
				}
			}
		case "k", "up":
			if m.state == promptView {
				m.choice -= 1
				if m.choice < 0 {
					m.choice = 0
				}
			}
		}

		switch msg.Type {
		case tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			if m.state == inputView {
				cmds = append(cmds, addToDb(m.params[m.currParam], m.TextInput.Value()))
				m.currParam++
				m.TextInput.Reset()
				m.state = promptView
			} else {
				if m.choice == 0 && m.currParam < 3 {
					m.state = inputView
				} else {
					if m.currParam > 2 {
						return m, tea.Quit
					}
					m.currParam++
					m.choice = 0
				}
			}
		}

		// update whichever model is focused
		switch m.state {
		case promptView:
			//
		default:
			m.TextInput, cmd = m.TextInput.Update(msg)
			cmds = append(cmds, cmd)
		}

	case tea.WindowSizeMsg:
		cmds = append(cmds, m.doResize(msg))
	}

	// if m.currParam > 3 {
	// 	return m, tea.Quit
	// }
	return m, tea.Batch(cmds...)
}

// input view
func (m *mainModel) viewInput() string {
	var s string
	var param string

	switch m.currParam {
	case 0:
		param = "Auth Token"
	case 1:
		param = "Organization Id"
	case 2:
		param = "Project Template Id"
	}

	s = fmt.Sprintf("Enter %s\n\n%s\n\n", param, m.TextInput.View())

	return focusedModelStyle.Width(m.width / 2).Height(m.height / 4).Align(lipgloss.Center).Render(s)
}

// prompt view
func (m *mainModel) viewPrompt() string {
	var promptLabel string
	var choices string

	c := m.choice

	switch m.currParam {
	case 1:
		promptLabel = "Do you have Organization Id?\n\n"
		choices = lipgloss.JoinVertical(lipgloss.Left, checkbox("yes", c == 0), checkbox("no", c == 1))
	case 2:
		promptLabel = "Do you have Project Template Id?\n\n"
		choices = lipgloss.JoinVertical(lipgloss.Left, checkbox("yes", c == 0), checkbox("no", c == 1))
	case 3:
		promptLabel = "Great! Run 'oah test' to test some endpoints.\n\n"
		choices = lipgloss.JoinVertical(lipgloss.Left, checkbox("got it!", true))
	}

	promptLabel += "%s\n"

	return focusedModelStyle.Width(m.width / 3).Height(m.height / 2).Align(lipgloss.Center).Render(fmt.Sprintf(promptLabel, choices))
}

// main view
func (m *mainModel) View() string {
	if m.err != nil {
		return fmt.Sprintf("Encountered Error: %v", m.err)
	}

	inputBox := m.viewInput()
	promptBox := m.viewPrompt()
	footer := helpStyle.Render("\n↑/↓, j/k: navigate • ↵: enter/select • esc: exit\n")

	if m.state == inputView {
		complete := lipgloss.JoinVertical(lipgloss.Center, inputBox, footer)
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, complete)

	}

	complete := lipgloss.JoinVertical(lipgloss.Center, promptBox, footer)
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, complete)
}

// tea command to add value to db
func addToDb(key string, value string) tea.Cmd {
	return func() tea.Msg {
		if err := data.UpdateX(key, value); err != nil {
			return errMsg{err}
		}
		return nil
	}
}

// tea command to create table and do default insert
func checkDatabase() tea.Msg {
	if err := data.CreateTable(); err != nil {
		return errMsg{err}
	}
	return nil
}

// tea command to re-render app when window is resized
func (m *mainModel) doResize(msg tea.WindowSizeMsg) tea.Cmd {
	m.height = msg.Height
	m.width = msg.Width
	return nil
}

// format checkboxes for prompt view
func checkbox(label string, checked bool) string {
	if checked {
		return selectedChoiceStyle.Render("[x] " + label)
	}
	return choiceStyle.Render("[ ] " + label)
}

// cobra setup
var SetupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Enter Token and other parameters.",
	Run: func(cmd *cobra.Command, args []string) {
		p := tea.NewProgram(initialModel(), tea.WithAltScreen())

		if err := p.Start(); err != nil {
			log.Fatal(err)
		}
	},
}
