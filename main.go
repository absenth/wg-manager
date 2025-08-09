package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	args := os.Args[1:]
	config, state, list, check := parseArgs(args)

	if len(args) == 0 {
		p := tea.NewProgram(initialModel())
		if err := p.Start(); err != nil {
			fmt.Println("Error running program:", err)
			os.Exit(1)
		}
		return
	}

	if list {
		listConfigs()
	}

	if config != "" && state != "" {
		runWgQuick(config, state)
	} else if config == "" && state != "" {
		fmt.Println(ErrorStyle.Render("Configuration file name is required for --state"))
		os.Exit(1)
	} else if config != "" && state == "" {
		fmt.Println(ErrorStyle.Render("State (up/down) is required for --config"))
		os.Exit(1)
	}

	if check {
		checkWG()
	}
}

func usage() {
	fmt.Println(TitleStyle.Render(`Usage: wg-manager
--config <config_name> --state <up|down>
--list: list available WireGuard configuration files
--check: Run the wg command and display its output with formatting`))
	os.Exit(1)
}

func parseArgs(args []string) (config, state string, list, check bool) {
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--config":
			if i+1 >= len(args) {
				fmt.Println(ErrorStyle.Render("Missing value for --config"))
				usage()
			}
			config = args[i+1]
			i++
		case "--state":
			if i+1 >= len(args) {
				fmt.Println(ErrorStyle.Render("Missing value for --state"))
				usage()
			}
			state = args[i+1]
			i++
		case "--list":
			list = true
		case "--check":
			check = true
		case "--help":
			usage()
		default:
			fmt.Println(ErrorStyle.Render("Unknown parameter passed: " + args[i]))
			usage()
		}
	}

	return config, state, list, check
}

