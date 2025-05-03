package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	leftMargin   = 1
	primaryColor = lipgloss.Color("#87CEEB")
	mutedColor   = lipgloss.Color("#6C6C6C")
	lightColor   = lipgloss.Color("#D0D0D0")

	SearchStyle        = lipgloss.NewStyle().Foreground(lightColor).MarginLeft(leftMargin).PaddingBottom(1)
	NumberStyle        = lipgloss.NewStyle().Foreground(mutedColor)
	BranchStyle        = lipgloss.NewStyle().Foreground(lightColor)
	CurrentBranchStyle = lipgloss.NewStyle().Foreground(lightColor)
	SelectedStyle      = lipgloss.NewStyle().Foreground(primaryColor)
	HelpTextStyle      = lipgloss.NewStyle().Foreground(mutedColor)
)

func HorizontalLayout(leftContent, rightContent string, totalWidth int) string {
	if totalWidth <= 0 {
		totalWidth = 80 // Default width
	}

	spacing := totalWidth - lipgloss.Width(leftContent) - lipgloss.Width(rightContent)
	if spacing < 1 {
		spacing = 1
	}

	return lipgloss.NewStyle().Width(totalWidth).Render(
		lipgloss.JoinHorizontal(
			lipgloss.Left,
			leftContent,
			strings.Repeat(" ", spacing),
			rightContent,
		),
	)
}
