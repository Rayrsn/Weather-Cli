package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type errMsg error

type LoadingModel struct {
	spinner  spinner.Model
	quitting bool
	err      error
	message  string
	Done     chan bool
}

func NewLoadingModel(message string) LoadingModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return LoadingModel{
		spinner: s,
		message: message,
		Done:    make(chan bool),
	}
}

func (m LoadingModel) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m LoadingModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		default:
			return m, nil
		}

	case bool:
		m.quitting = true
		return m, tea.Quit

	case errMsg:
		m.err = msg
		return m, tea.Quit

	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
}

func (m LoadingModel) View() string {
	if m.err != nil {
		return fmt.Sprintf("\nError: %v\n\n", m.err)
	}
	if m.quitting {
		return ""
	}
	return fmt.Sprintf("\n %s %s...\n\n", m.spinner.View(), m.message)
}
