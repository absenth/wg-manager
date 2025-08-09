package main

import "github.com/charmbracelet/lipgloss"

var (
	TitleStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("14"))
	ConfigStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	InfoStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("14"))
	ErrorStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
)

