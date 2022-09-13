/*
Copyright © 2022 Sabino Ramirez <sabinoramirez017@gmail.com>
*/
package test

// imports
import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sabino-ramirez/oah/data"
	"github.com/sabino-ramirez/oah/models"
	"github.com/sabino-ramirez/oah/utils"
	"github.com/spf13/cobra"
)

// tea message type for handling errors throughout program
type errMsg struct{ err error }
type statusMsg int

// app state variables will have this type
type sessionState uint

// constants for keeping track of app state
const (
	dbItemsView sessionState = iota
	resultsView
)

var netTransport = &http.Transport{
	Dial: (&net.Dialer{
		Timeout: 5 * time.Second,
	}).Dial,
	TLSHandshakeTimeout: 5 * time.Second,
}

// lipgloss styles
var (
	modelStyle = lipgloss.NewStyle().Padding(0, 0, 0, 0).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("#3C3C3C")).
			Foreground(lipgloss.Color("#3C3C3C"))

	focusedModelStyle = lipgloss.NewStyle().Padding(0, 0, 0, 0).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("69"))

	helpStyle = lipgloss.NewStyle().Align(lipgloss.Center).
			Foreground(lipgloss.Color("241"))

	choiceStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("254"))
	selectedChoiceStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("5"))

	baseStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240"))
)

type mainModel struct {
	state          sessionState
	prompt         bool
	chooseEndpoint bool
	choice         int
	dbItems        models.DbRow
	statusCode     int
	currParam      string
	table          table.Model
	textInput      textinput.Model
	width          int
	height         int
	err            error
}

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

	m := mainModel{state: resultsView, table: t, textInput: ti, chooseEndpoint: true}
	return &m
}

// calls the refreshDb command and kicks off the program
func (m *mainModel) Init() tea.Cmd {
	return m.refreshDbItems
}

// main update
func (m *mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case statusMsg:
		m.statusCode = int(msg)

	case errMsg:
		m.err = msg

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
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
					return m, addToDb(m.currParam, m.textInput.Value())
				}
			case resultsView:
				if m.chooseEndpoint {
					m.chooseEndpoint = false
					return m, m.checkStatusCode(m.choice)
				} else {
					m.chooseEndpoint = true
				}
			}

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
			} else {
				m.table, cmd = m.table.Update(msg)
				return m, cmd
			}
		} else {
			switch msg.String() {
			case "j", "down":
				if m.state == resultsView {
					m.choice += 1
					if m.choice > 1 {
						m.choice = 1
					}
				}
			case "k", "up":
				if m.state == resultsView {
					m.choice -= 1
					if m.choice < 0 {
						m.choice = 0
					}
				}
			}
		}

	case tea.WindowSizeMsg:
		return m, m.doResize(msg)
	}

	return m, m.refreshDbItems
}

// returns view for dbitems adjustment
func (m *mainModel) viewDbItems() string {
	m.table.SetRows([]table.Row{
		{"Auth", m.dbItems.Auth},
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
	var promptLabel string
	var choices string
	var s string

	c := m.choice

	promptLabel = "Select Test Operation\n\n%s\n"
	choices = lipgloss.JoinVertical(lipgloss.Left, checkbox("Get Project Templates", c == 0), checkbox("Get Requisitions", c == 1))

	if m.chooseEndpoint {
		s = fmt.Sprintf(promptLabel, choices)
	} else {
		s = fmt.Sprintf("status code is: %v", m.statusCode)
	}

	// s := focusedModelStyle.Width(m.width / 3).Height(m.height / 2).Align(lipgloss.Center).Render(fmt.Sprintf(promptLabel, choices))
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

// cmd to re-render app when window is resized
func (m *mainModel) doResize(msg tea.WindowSizeMsg) tea.Cmd {
	m.height = msg.Height
	m.width = msg.Width
	return nil
}

// cmd to refresh table with db values
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

// cmd update db value
func addToDb(key string, value string) tea.Cmd {
	return func() tea.Msg {
		if err := data.UpdateX(key, value); err != nil {
			return errMsg{err}
		}
		return nil
	}
}

// cmd for getting status code based on endpoint position in results view list
func (m *mainModel) checkStatusCode(choice int) tea.Cmd {
	return func() tea.Msg {
		client := &http.Client{Timeout: time.Second * 10, Transport: netTransport}
		ovationAPI := models.NewClient(client, 1, 1, "Bearer "+m.dbItems.Auth)

		var projectReqs models.ProjectRequisitions
		var projectIds models.ProjectTemplates

		var statusCode int

		switch choice {
		case 0:
			statusCode, _ = utils.GetProjectTemplates(ovationAPI, &projectIds)
			// if err != nil {
			// 	log.Println("error reading response:", err)
			// }

		case 1:
			statusCode, _ = utils.GetProjectRequisitions(ovationAPI, &projectReqs)
			// if err != nil {
			// log.Println("error reading response:", err)
			// }
		}

		return statusMsg(statusCode)
	}
}

// format checkboxes for prompt view
func checkbox(label string, checked bool) string {
	if checked {
		return selectedChoiceStyle.Render("[x] " + label)
	}
	return choiceStyle.Render("[ ] " + label)
}

// in order to get errMsg type to implement error interface
func (e errMsg) Error() string { return e.err.Error() }

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
