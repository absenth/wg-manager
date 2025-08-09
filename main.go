package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFD700"))

	infoStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#00CED1"))

	configStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FF4500"))

	currentIPStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FF6347"))

	newIPStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#00FA9A"))

	errorStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FF0000"))
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		usage()
	}

	config, state, list, check := parseArgs(args)

	if list {
		listConfigs()
	}

	if config != "" && state != "" {
		runWgQuick(config, state)
	} else if check {
		checkWG()
	} else {
		runInteractiveMode()
	}
}

func runInteractiveMode() {
	var selectedConfig, action string
	configs := getConfigList()
	if len(configs) == 0 {
		fmt.Println(errorStyle.Render("No WireGuard configuration found at /etc/wireguard"))
		os.Exit(1)
	}

	prompt := &survey.Select{
		Message: "Select a WireGuard configuration:",
		Options: configs,
	}
	if err := survey.AskOne(prompt, &selectedConfig); err != nil {
		fmt.Println(errorStyle.Render("Error selecting configuration: " + err.Error()))
		os.Exit(1)
	}

	actionPrompt := &survey.Select{
		Message: "Select an action:",
		Options: []string{"up", "down"},
	}
	if err := survey.AskOne(actionPrompt, &action); err != nil {
		fmt.Println(errorStyle.Render("Error selecting action: " + err.Error()))
		os.Exit(1)
	}

	fmt.Println(titleStyle.Render("Current IP address is:"))
	currentIP, err := getPublicIP()
	if err != nil {
		fmt.Println(errorStyle.Render("Error getting current IP address: " + err.Error()))
		os.Exit(1)
	}
	fmt.Println(currentIPStyle.Render(currentIP))
	fmt.Println(titleStyle.Render("----------------------"))

	wgQuickCmd := exec.Command("wg-quick", action, selectedConfig)
	if err := wgQuickCmd.Run(); err != nil {
		fmt.Println(errorStyle.Render("Error running wg-quick: " + err.Error()))
		os.Exit(1)
	}

	fmt.Println(titleStyle.Render("New IP address is:"))
	newIP, err := getPublicIP()
	if err != nil {
		fmt.Println(errorStyle.Render("Error getting new IP address: " + err.Error()))
		os.Exit(1)
	}
	fmt.Println(newIPStyle.Render(newIP))
	fmt.Println(titleStyle.Render("----------------------"))
}

func usage() {
	fmt.Println(titleStyle.Render("Usage: ./wg-manager --config <config_name> --state <up|down> | --list | --check"))
	fmt.Println(infoStyle.Render("   --config: Path to the WireGuard configuration file (without .conf extension)"))
	fmt.Println(infoStyle.Render("   --state: Action to perform (up or down)"))
	fmt.Println(infoStyle.Render("   --list: List available WireGuard configuration files"))
	fmt.Println(infoStyle.Render("   --check: Run the wg command and display its output with formatting"))
	os.Exit(1)
}

func listConfigs() {
	configs := getConfigList()
	if len(configs) == 0 {
		fmt.Println(errorStyle.Render("No WireGuard configuration found at /etc/wireguard"))
		os.Exit(1)
	}

	fmt.Println(titleStyle.Render("Available WireGuard configurations:"))
	for _, config := range configs {
		fmt.Println(configStyle.Render(config))
	}
	os.Exit(0)
}

func getConfigList() []string {
	configs := []string{}
	files, err := filepath.Glob("/etc/wireguard/*.conf")
	if err != nil {
		fmt.Println(errorStyle.Render("Error reading WireGuard configuration directory: " + err.Error()))
		return configs
	}

	for _, file := range files {
		configs = append(configs, strings.TrimSuffix(filepath.Base(file), ".conf"))
	}

	return configs
}

func parseArgs(args []string) (config, state string, list, check bool) {
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--config":
			if i+1 >= len(args) {
				fmt.Println(errorStyle.Render("Missing value for --config"))
				usage()
			}
			config = args[i+1]
			i++
		case "--state":
			if i+1 >= len(args) {
				fmt.Println(errorStyle.Render("Missing value for --state"))
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
			fmt.Println(errorStyle.Render("Unknown parameter passed: " + args[i]))
			usage()
		}
	}

	return config, state, list, check
}

func getPublicIP() (string, error) {
	resp, err := http.Get("https://api.ipify.org?format=text")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	ip, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(ip)), nil
}

func runWgQuick(config, state string) {
	fmt.Println(titleStyle.Render("Current IP address is:"))
	currentIP, err := getPublicIP()
	if err != nil {
		fmt.Println(errorStyle.Render("Error getting current IP address: " + err.Error()))
		os.Exit(1)
	}
	fmt.Println(currentIPStyle.Render(currentIP))
	fmt.Println(titleStyle.Render("----------------------"))

	wgQuickCmd := exec.Command("wg-quick", state, config)
	if err := wgQuickCmd.Run(); err != nil {
		fmt.Println(errorStyle.Render("Error running wg-quick: " + err.Error()))
		os.Exit(1)
	}

	fmt.Println(titleStyle.Render("New IP address is:"))
	newIP, err := getPublicIP()
	if err != nil {
		fmt.Println(errorStyle.Render("Error getting new IP address: " + err.Error()))
		os.Exit(1)
	}
	fmt.Println(newIPStyle.Render(newIP))
	fmt.Println(titleStyle.Render("----------------------"))
}

func checkWG() {
	wgOutput, err := runWgCommand()
	if err != nil {
		fmt.Println(errorStyle.Render("Error running wg command: " + err.Error()))
		os.Exit(1)
	}

	formattedOutput := formatWgOutput(wgOutput)
	fmt.Println(formattedOutput)
}

func runWgCommand() (string, error) {
	wgCmd := exec.Command("wg")
	output, err := wgCmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func formatWgOutput(output string) string {
	if output == "" {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD700")).Render("No VPN is currently active.")
	}

	scanner := bufio.NewScanner(strings.NewReader(output))
	var interfaceName string

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "interface:") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				interfaceName = strings.TrimSpace(parts[1])
				break
			}
		}
	}

	if interfaceName == "" {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD700")).Render("No VPN is currently active.")
	}

	return lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF7F")).Render(fmt.Sprintf("VPN: Active\nConnection: %s", interfaceName))
}

