package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func listConfigs() {
	configs := getConfigList()
	if len(configs) == 0 {
		fmt.Println(ErrorStyle.Render("No WireGuard configuration found at /etc/wireguard"))
		os.Exit(1)
	}

	fmt.Println(TitleStyle.Render("Available WireGuard configurations:"))
	for _, config := range configs {
		fmt.Println(ConfigStyle.Render(config))
	}
	os.Exit(0)
}

func getConfigList() []string {
	configs := []string{}
	files, err := filepath.Glob("/etc/wireguard/*.conf")
	if err != nil {
		fmt.Println(ErrorStyle.Render("Error reading WireGuard configuration directory: " + err.Error()))
		return configs
	}

	for _, file := range files {
		configs = append(configs, strings.TrimSuffix(filepath.Base(file), ".conf"))
	}

	return configs
}

func runWgQuick(config, state string) {
	wgQuickCmd := exec.Command("wg-quick", state, config)
	if err := wgQuickCmd.Run(); err != nil {
		fmt.Println(ErrorStyle.Render("Error running wg-quick: " + err.Error()))
		os.Exit(1)
	}
}

func checkWG() {
	wgOutput, err := runWgCommand()
	if err != nil {
		fmt.Println(ErrorStyle.Render("Error running wg command: " + err.Error()))
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
		return ErrorStyle.Render("No VPN is currently active.")
	}

	scanner := bufio.NewScanner(strings.NewReader(output))
	var formattedOutput string

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "interface:") || strings.HasPrefix(line, "peer:") {
			formattedOutput += TitleStyle.Render(line) + "\n"
		} else if strings.HasPrefix(line, "public key:") || strings.HasPrefix(line, "private key:") ||
			strings.HasPrefix(line, "listening port:") {
			formattedOutput += ConfigStyle.Render(line) + "\n"
		} else if strings.HasPrefix(line, "endpoint:") || strings.HasPrefix(line, "allowed ips:") ||
			strings.HasPrefix(line, "latest handshake:") || strings.HasPrefix(line, "transfer:") {
			formattedOutput += InfoStyle.Render(line) + "\n"
		} else {
			formattedOutput += lipgloss.NewStyle().Render(line) + "\n"
		}
	}

	if err := scanner.Err(); err != nil {
		return ErrorStyle.Render("Error reading wg output: " + err.Error())
	}

	return formattedOutput
}
