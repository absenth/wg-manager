package main

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	configs []string
	cursor  int
	state   string
}

func initialModel() model {
	return model{
		configs: getConfigList(),
		cursor:  0,
		state:   "",
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

type tickMsg struct{}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.configs)-1 {
				m.cursor++
			}
		case "enter":
			if m.state == "" {
				return m, nil
			}
			runWgQuick(m.configs[m.cursor], m.state)
			return m, tea.Quit
		case "u":
			if m.state != "up" {
				m.state = "up"
			} else {
				m.state = ""
			}
		case "d":
			if m.state != "down" {
				m.state = "down"
			} else {
				m.state = ""
			}
		}
	case tickMsg:
		return m, nil
	}

	return m, tickCmd()
}

func (m model) View() string {
	s := TitleStyle.Render("WireGuard Configuration Manager") + "\n\n"

	for i, config := range m.configs {
		cursor := "  • "
		if i == m.cursor {
			cursor = "->"
		}
		s += fmt.Sprintf("%s %s\n", cursor, ConfigStyle.Render(config))
	}

	stateInfo := ""
	if m.state == "up" {
		stateInfo = InfoStyle.Render(" [UP]")
	} else if m.state == "down" {
		stateInfo = InfoStyle.Render(" [DOWN]")
	}
	s += fmt.Sprintf("\nState: %s\n", stateInfo)

	s += "\nPress 'u' to set state to up, 'd' to set state to down, and 'enter' to apply.\n"
	s += "Press 'q' or 'ctrl+c' to quit."

	return s
}

func tickCmd() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(time.Time) tea.Msg {
		return tickMsg{}
	})
}
