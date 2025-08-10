package main

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	configs []string
	cursor  int
	state   string
}

func initialModel() model {
	// Fetch available configurations
	configs := getConfigList()

	// Check if a VPN is active and get its state
	activeState, activeConfig := getVPNState()

	// Set default state to UP if no VPN is active
	if activeState == "" {
		activeState = "up"
	}

	// Auto-select the active configuration if applicable
	cursor := 0
	if activeConfig != "" {
		for i, config := range configs {
			if config == activeConfig {
				cursor = i
				break
			}
		}
	}

	return model{
		configs: configs,
		cursor:  cursor,
		state:   activeState,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

type tickMsg struct{}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
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
			// Apply the configuration state
			runWgQuick(m.configs[m.cursor], m.state)
			return m, tea.Quit
		case "u":
			m.state = "up"
		case "d":
			m.state = "down"
		}
	case tickMsg:
		return m, nil
	}

	return m, tea.Batch(tickCmd(), tea.Blink)
}

func (m model) View() string {
	// Render the Interface
	s := TitleStyle.Render("WireGuard Configuration Manager") + "\n\n"

	for i, config := range m.configs {
		cursor := "  • "
		if i == m.cursor {
			cursor = "->"
		}
		s += fmt.Sprintf("%s %s\n", cursor, ConfigStyle.Render(config))
	}

	stateInfo := InfoStyle.Render(fmt.Sprintf(" [%s]", strings.ToUpper(m.state)))
	s += fmt.Sprintf("\nState: %s\n", stateInfo)

	s += "\nPress 'u' to set state to UP, 'd' to set state to DOWN, and 'enter' to apply.\n"
	s += "Press 'q' or 'ctrl+c' to quit."

	return s
}

func tickCmd() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(time.Time) tea.Msg {
		return tickMsg{}
	})
}

func getVPNState() (string, string) {
	// Execute `wg` command to evaluate VPN status
	out, err := exec.Command("wg").Output()
	if err != nil || len(out) == 0 {
		// No VPN is active
		return "", ""
	}

	// Parse output for active configuration and state
	output := strings.TrimSpace(string(out))
	lines := strings.Split(output, "\n")

	if len(lines) > 0 {
		// Derive active configuration name (first line after "interface:")
		activeConfig := strings.TrimSpace(strings.TrimPrefix(lines[0], "interface:"))
		return "down", activeConfig
	}

	return "", ""
}

func runWgQuick(config string, state string) {
	// Execute `wg-quick up/down <configuration>` based on the state
	cmd := exec.Command("wg-quick", state, config)
	cmd.Run()
}
