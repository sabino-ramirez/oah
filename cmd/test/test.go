/*
Copyright © 2022 Sabino Ramirez <sabinoramirez017@gmail.com>
*/
package test

import (
	"fmt"
	"log"
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sabino-ramirez/oah/data"
	"github.com/sabino-ramirez/oah/models"
	"github.com/spf13/cobra"
)

// tea message type for handling errors throughout program
type errMsg struct{ err error }

// type dbMessage models.DbRow

// app state variables will have this type
type sessionState uint

// constants for keeping track of app state
const (
	dbItemsView sessionState = iota
	resultsView
)

// lipgloss styles
var (
	modelStyle        = lipgloss.NewStyle().Padding(0, 0, 0, 0).BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("#3C3C3C")).Foreground(lipgloss.Color("#3C3C3C"))
	focusedModelStyle = lipgloss.NewStyle().Padding(0, 0, 0, 0).BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("69"))
	helpStyle         = lipgloss.NewStyle().Align(lipgloss.Center).Foreground(lipgloss.Color("241"))
	baseStyle         = lipgloss.NewStyle().
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("240"))
)

// in order to get errMsg type to implement error interface
func (e errMsg) Error() string { return e.err.Error() }

type mainModel struct {
	state     sessionState
	prompt    bool
	dbItems   models.DbRow
	currParam string
	table     table.Model
	textInput textinput.Model
	width     int
	height    int
	err       error
}

// initial command to get variables from db
// func getDbInfo() tea.Msg {
// 	params, err := data.GetValues()
// 	if err != nil {
// 		return errMsg{err}
// 	}
//
// 	return dbMessage(params)
// }

// function returns initial state
func initialModel() *mainModel {
	ti := textinput.New()
	ti.Placeholder = "copy/paste or type.."
	ti.Focus()
	ti.Width = 20

	columns := []table.Column{
		{Title: "Key", Width: 8},
		{Title: "Value", Width: 10},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
	)

	m := mainModel{state: resultsView, table: t, textInput: ti}
	return &m
}

// calls the getDbInfo command and kicks off the program
func (m *mainModel) Init() tea.Cmd {
	return m.refreshDbItems //getDbInfo
}

func (m *mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	// var cmds []tea.Cmd

	switch msg := msg.(type) {
	case errMsg:
		m.err = msg

	// case dbMessage:
	// 	m.dbItems = models.DbRow(msg)

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			// cmds = append(cmds, tea.Quit)
			return m, tea.Quit

		case tea.KeyEnter:
			switch m.state {
			case dbItemsView:
				if m.prompt == false {
					switch m.table.SelectedRow()[0] {
					case "Auth":
						m.currParam = "auth"
					case "Org Id":
						m.currParam = "orgId"
					case "Proj. Temp. Id":
						m.currParam = "projTempId"
					}
					m.prompt = true
				} else {
					m.prompt = false
					// m.textInput.Reset()
					return m, addToDb(m.currParam, m.textInput.Value())
				}
			}
			// cmds = append(cmds, m.refreshDbItems)

		case tea.KeyTab:
			if m.state == dbItemsView {
				m.state = resultsView
			} else {
				m.state = dbItemsView
			}
		}

		if m.state == dbItemsView {
			if m.prompt {
				m.textInput, cmd = m.textInput.Update(msg)
				return m, cmd
				// cmds = append(cmds, cmd)
			} else {
				m.table, cmd = m.table.Update(msg)
				return m, cmd
				// cmds = append(cmds, cmd)
			}
		}

	case tea.WindowSizeMsg:
		// cmds = append(cmds, m.doResize(msg))
		return m, m.doResize(msg)
	}

	// m.table, cmd = m.table.Update(msg)
	// cmds = append(cmds, cmd)

	return m, m.refreshDbItems
	// return m, tea.Batch(cmds...)
}

// returns view for dbitems adjustment
func (m *mainModel) viewDbItems() string {
	// tableContent := "Org Id  |  Proj. Temp. Id\n%v  |  %v"
	// tableContent = fmt.Sprintf(tableContent, m.dbItems.OrgId, m.dbItems.ProjTempId)
	// table := focusedModelStyle.Width(m.width / 2).Height(m.height / 8).Align(lipgloss.Center).Render(tableContent)
	// return tableContent

	m.table.SetRows([]table.Row{
		{"Auth", m.dbItems.Auth},
		// {"Org Id", m.dbItems.OrgId},
		// {"Proj. Temp. Id", m.dbItems.ProjTempId},
		{"Org Id", strconv.Itoa(m.dbItems.OrgId)},
		{"Proj. Temp. Id", strconv.Itoa(m.dbItems.ProjTempId)},
	})

	s := table.DefaultStyles()
	s.Cell.Align(lipgloss.Center)
	s.Header = s.Header.Align(lipgloss.Center).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)

	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(true)

	m.table.SetStyles(s)

	m.table.SetHeight(m.height / 12)
	m.table.SetWidth(m.width / 2)

	if m.prompt {
		return fmt.Sprintf("Enter %s\n\n%s\n\n", m.currParam, m.textInput.View())
	}

	return baseStyle.Render(m.table.View()) + "\n\nMake a selection to edit value."
}

// returns view for endpoint results
func (m *mainModel) viewResults() string {
	s := fmt.Sprintf("db items: %v", m.dbItems)
	return s
}

// main view
func (m *mainModel) View() string {
	var complete string
	dbItemsBox := m.viewDbItems()
	resultsBox := m.viewResults()
	footer := helpStyle.Render("\n↑/↓, j/k: navigate • ↵: enter/select\ntab: switch view • esc: exit\n")

	if m.state == dbItemsView {
		complete = lipgloss.JoinVertical(lipgloss.Center, focusedModelStyle.Width(m.width/2).Height(m.height/4).Align(lipgloss.Center).Render("DB Items\n"+dbItemsBox), modelStyle.Width(m.width/2).Height(m.height/10).Align(lipgloss.Center).Render("Results"), footer)
	} else {
		complete = lipgloss.JoinVertical(lipgloss.Center, modelStyle.Width(m.width/2).Height(m.height/10).Align(lipgloss.Center).Render("DB Items"), focusedModelStyle.Width(m.width/2).Height(m.height/2).Align(lipgloss.Center).Render("Results\n\n"+resultsBox), footer)
	}
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, complete)
}

// tea command to re-render app when window is resized
func (m *mainModel) doResize(msg tea.WindowSizeMsg) tea.Cmd {
	m.height = msg.Height
	m.width = msg.Width
	return nil
}

func (m *mainModel) refreshDbItems() tea.Msg {
	params, err := data.GetValues()

	m.dbItems.Auth = params.Auth
	m.dbItems.OrgId = params.OrgId
	m.dbItems.ProjTempId = params.ProjTempId

	if err != nil {
		return errMsg{err}
	}

	return nil
}

func addToDb(key string, value string) tea.Cmd {
	return func() tea.Msg {
		if err := data.UpdateX(key, value); err != nil {
			return errMsg{err}
		}
		return nil
	}
}

// cobra stuff
var TestCmd = &cobra.Command{
	Use:   "test",
	Short: "Test the endpoints with parameters entered in setup",
	Run: func(cmd *cobra.Command, args []string) {
		p := tea.NewProgram(initialModel(), tea.WithAltScreen())

		if err := p.Start(); err != nil {
			log.Fatal(err)
		}
	},
}
