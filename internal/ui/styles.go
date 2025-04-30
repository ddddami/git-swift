package ui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	primaryColor = lipgloss.Color("#87CEEB")
	mutedColor   = lipgloss.Color("#6C6C6C")
	lightColor   = lipgloss.Color("#D0D0D0")

	PromptStyle        = lipgloss.NewStyle().Foreground(primaryColor)
	SearchStyle        = lipgloss.NewStyle().Foreground(primaryColor)
	NumberStyle        = lipgloss.NewStyle().Foreground(mutedColor)
	BranchStyle        = lipgloss.NewStyle().Foreground(lightColor)
	CurrentBranchStyle = lipgloss.NewStyle().Foreground(lightColor)
	SelectedStyle      = lipgloss.NewStyle().Foreground(primaryColor)
	HelpTextStyle      = lipgloss.NewStyle().Foreground(mutedColor)
)
