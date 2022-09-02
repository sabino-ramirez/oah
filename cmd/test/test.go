/*
Copyright © 2022 Sabino Ramirez <sabinoramirez017@gmail.com>

*/
package test

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

// lipgloss styles
var (
	testStyle = lipgloss.NewStyle().Padding(2).BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("109"))
)

type mainModel struct {
	message string
}

func initialModel() *mainModel {
	m := mainModel{message: "hey it works i guess"}
	return &m
}

func (m *mainModel) Init() tea.Cmd {
	//*TODO* cannot get the cursor to blink
	return nil
}

func (m *mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// var cmd tea.Cmd
	// var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			return m, tea.Quit
		}

	}

	return m, nil
}

func (m *mainModel) View() string {
	return fmt.Sprintf("message: %s", m.message)
}

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
