package ui

import (
	"fmt"
	"io"
	"strings"

	"github.com/Rayrsn/weather-Cli/internal/api"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
)

type item struct {
	result api.GeocodingResult
}

func (i item) FilterValue() string { return i.result.Name + " " + i.result.Country }

type itemDelegate struct{}

func (d itemDelegate) Height() int                               { return 1 }
func (d itemDelegate) Spacing() int                              { return 0 }
func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s, %s (Lat: %.2f, Lon: %.2f)", index+1, i.result.Name, i.result.Country, i.result.Latitude, i.result.Longitude)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(strs ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(strs, ""))
		}
	}

	fmt.Fprint(w, fn(str))
}

type model struct {
	list     list.Model
	choice   *api.GeocodingResult
	quitting bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.choice = &i.result
			}
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.choice != nil {
		return ""
	}
	if m.quitting {
		return "Selection cancelled.\n"
	}
	return "\n" + m.list.View()
}

func SelectCity(results []api.GeocodingResult) (*api.GeocodingResult, error) {
	if len(results) == 0 {
		return nil, fmt.Errorf("no cities found")
	}
	if len(results) == 1 {
		return &results[0], nil
	}

	items := make([]list.Item, len(results))
	for i, res := range results {
		items[i] = item{result: res}
	}

	const defaultWidth = 20
	const listHeight = 14

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "Multiple cities found. Please select one:"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = lipgloss.NewStyle().Foreground(lipgloss.Color("170")).Bold(true).Padding(0, 1)

	m := model{list: l}

	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		return nil, err
	}

	m = finalModel.(model)
	if m.choice == nil {
		return nil, fmt.Errorf("no city selected")
	}

	return m.choice, nil
}
