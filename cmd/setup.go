/*
Copyright © 2022 Sabino Ramirez <sabinoramirez017@gmail.com>

*/
package cmd

import (
	"fmt"
	"log"
	"time"

	// "github.com/charmbracelet/bubbles/spinner"
	// "github.com/charmbracelet/bubbles/timer"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/spf13/cobra"
)

var (
	dot = " • "
)

type sessionState uint

const (
	defaultTime              = time.Minute
	inputView   sessionState = iota
	promptView
)

var (
	// Available spinners
	// spinners = []spinner.Spinner{
	// 	spinner.Line,
	// 	spinner.Dot,
	// 	spinner.MiniDot,
	// 	spinner.Jump,
	// 	spinner.Pulse,
	// 	spinner.Points,
	// 	spinner.Globe,
	// 	spinner.Moon,
	// 	spinner.Monkey,
	// }
	// unfocusedModelStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#3C3C3C"))
	// modelStyle = lipgloss.NewStyle().
	// 		Padding(2).
	// 		BorderStyle(lipgloss.NormalBorder())

	focusedModelStyle = lipgloss.NewStyle().
		// Padding(2).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("109"))
	spinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("69"))
	helpStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
)

type mainModel struct {
	state     sessionState
	TextInput textinput.Model
	// timer   timer.Model
	// spinner spinner.Model
	// index   int
	height int
	width  int
}

func newModel(timeout time.Duration) *mainModel {
	ti := textinput.New()
	ti.Placeholder = "copy/paste or type.."
	ti.Focus()
	ti.Width = 20
	m := mainModel{state: inputView, TextInput: ti}
	// m.timer = timer.New(timeout)
	// m.spinner = spinner.New()
	return &m
}

func (m *mainModel) Init() tea.Cmd {
	// start the timer and spinner on program start
	// return tea.Batch(m.timer.Init(), m.spinner.Tick)
	return textinput.Blink
}

func (m *mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			if m.state == inputView {
				m.state = promptView
			} else {
				m.state = inputView
			}
			// case "n":
			// 	if m.state == inputView {
			// 		m.timer = timer.New(defaultTime)
			// 		cmds = append(cmds, m.timer.Init())
			// 	} else {
			// 		m.Next()
			// 		m.resetSpinner()
			// 		cmds = append(cmds, spinner.Tick)
			// 	}
		}
		switch m.state {
		// update whichever model is focused
		case promptView:
			// m.spinner, cmd = m.spinner.Update(msg)
			// cmds = append(cmds, cmd)
		default:
			// m.timer, cmd = m.timer.Update(msg)
			// cmds = append(cmds, cmd)
			m.TextInput.Focus()
			m.TextInput, cmd = m.TextInput.Update(msg)
			cmds = append(cmds, cmd)
		}

	case tea.WindowSizeMsg:
		// return m, m.doResize(msg)
		cmds = append(cmds, m.doResize(msg))
		// case spinner.TickMsg:
		// 	m.spinner, cmd = m.spinner.Update(msg)
		// 	cmds = append(cmds, cmd)
		// case timer.TickMsg:
		// 	m.timer, cmd = m.timer.Update(msg)
		// 	cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func (m *mainModel) viewInput() string {
	var s string
	// s = fmt.Sprintf("%s", m.timer.View())
	s = fmt.Sprintf("Enter Auth Token\n\n%s\n\n", m.TextInput.View())

	return focusedModelStyle.Width(m.width / 2).Height(m.height / 4).Align(lipgloss.Center).Render(s)
}

func (m *mainModel) viewPrompt() string {
	var s string

	return focusedModelStyle.Width(m.width / 3).Height(m.height / 2).Align(lipgloss.Center).Render(s)
}

// mview
func (m *mainModel) View() string {
	// model := m.currentFocusedModel()
	// var footer string
	// var complete string

	inputBox := m.viewInput()
	promptBox := m.viewPrompt()
	// footer := lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Bottom, helpStyle.Render("\nesc: exit\n"))
	// _, h := lipgloss.Size(promptBox)

	if m.state == inputView {
		// complete = lipgloss.JoinVertical(lipgloss.Center, inputBox, footer)
		complete := lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, inputBox)
		return lipgloss.JoinVertical(lipgloss.Center, complete)

	}

	complete := lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, promptBox)
	return lipgloss.JoinVertical(lipgloss.Center, complete)
}

func (m *mainModel) doResize(msg tea.WindowSizeMsg) tea.Cmd {
	m.height = msg.Height
	m.width = msg.Width
	return nil
}

// func (m *mainModel) currentFocusedModel() string {
// 	if m.state == inputView {
// 		return "timer"
// 	}
// 	return "spinner"
// }

// func (m *mainModel) Next() {
// 	if m.index == len(spinners)-1 {
// 		m.index = 0
// 	} else {
// 		m.index++
// 	}
// }

func (m *mainModel) resetSpinner() {
	// m.spinner = spinner.New()
	// m.spinner.Style = spinnerStyle
	// m.spinner.Spinner = spinners[m.index]
}

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Enter Token and other parameters.",
	Run: func(cmd *cobra.Command, args []string) {
		p := tea.NewProgram(newModel(defaultTime), tea.WithAltScreen())

		if err := p.Start(); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
}
