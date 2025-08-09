package main

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

func main() {
	args := os.Args[1:]
	config, state, list, check := parseArgs(args)

	if list {
		listConfigs()
	}

	if config != "" && state != "" {
		runWgQuick(config, state)
	} else if config == "" && state != "" {
		fmt.Println(errorStyle.Render("Configuration file name is required for --state"))
		os.Exit(1)
	} else if config != "" && state == "" {
		fmt.Println(errorStyle.Render("State (up/down) is required for --config"))
		os.Exit(1)
	}

	if check {
		checkWG()
	}
}

func usage() {
	fmt.Println(lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("14")).Render(`Usage: wg-manager
--config <config_name> --state <up|down>
--list: list available WireGuard configuration files
--check: Run the wg command and display its output with formatting`))
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

func getPublicIPWithTimeout(timeout time.Duration) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.ipify.org?format=text", nil)
	if err != nil {
		return "", err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
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
	wgQuickCmd := exec.Command("wg-quick", state, config)
	if err := wgQuickCmd.Run(); err != nil {
		fmt.Println(errorStyle.Render("Error running wg-quick: " + err.Error()))
		os.Exit(1)
	}
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
		return errorStyle.Render("No VPN is currently active.")
	}

	scanner := bufio.NewScanner(strings.NewReader(output))
	var formattedOutput string

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "interface:") || strings.HasPrefix(line, "peer:") {
			formattedOutput += titleStyle.Render(line) + "\n"
		} else if strings.HasPrefix(line, "public key:") || strings.HasPrefix(line, "private key:") ||
			strings.HasPrefix(line, "listening port:") {
			formattedOutput += configStyle.Render(line) + "\n"
		} else if strings.HasPrefix(line, "endpoint:") || strings.HasPrefix(line, "allowed ips:") ||
			strings.HasPrefix(line, "latest handshake:") || strings.HasPrefix(line, "transfer:") {
			formattedOutput += infoStyle.Render(line) + "\n"
		} else {
			formattedOutput += lipgloss.NewStyle().Render(line) + "\n"
		}
	}

	if err := scanner.Err(); err != nil {
		return errorStyle.Render("Error reading wg output: " + err.Error())
	}

	return formattedOutput
}

var (
	titleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("14"))
	configStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	infoStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("14"))
	errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
)

