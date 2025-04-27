package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	primaryColor = lipgloss.Color("#87CEEB")
	mutedColor   = lipgloss.Color("#6C6C6C")
	lightColor   = lipgloss.Color("#D0D0D0")

	promptStyle = lipgloss.NewStyle().Foreground(primaryColor)
	searchStyle = lipgloss.NewStyle().Foreground(primaryColor)

	numberStyle = lipgloss.NewStyle().Foreground(mutedColor)
	branchStyle = lipgloss.NewStyle().Foreground(lightColor)

	currentBranchStyle = lipgloss.NewStyle().Foreground(lightColor)
	selectedStyle      = lipgloss.NewStyle().Foreground(primaryColor)

	helpTextStyle = lipgloss.NewStyle().Foreground(mutedColor)
)

type model struct {
	branches         []string
	filteredBranches []string
	currentBranch    string
	textInput        textinput.Model
	cursor           int
	err              error
}

func fuzzyMatch(branch, query string) bool {
	if query == "" {
		return true
	}

	branch = strings.ToLower(branch)
	query = strings.ToLower(query)

	branchIdx := 0
	queryIdx := 0

	for queryIdx < len(query) && branchIdx < len(branch) {
		if query[queryIdx] == branch[branchIdx] {
			queryIdx++
		}
		branchIdx++
	}

	return queryIdx == len(query)
}

func (m *model) filter() {
	query := m.textInput.Value()

	if query == "" {
		m.filteredBranches = m.branches
		return
	}

	filtered := []string{}
	for _, branch := range m.branches {
		if fuzzyMatch(branch, query) {
			filtered = append(filtered, branch)
		}
	}
	m.filteredBranches = filtered
}

func directSwitch(branchName string) {
	cmd := exec.Command("git", "switch", branchName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error switching branch: %s\n", string(output))
		os.Exit(1)
	}
	fmt.Printf("Switched to branch '%s'\n", branchName)
	os.Exit(0)
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit

		case tea.KeyUp, tea.KeyCtrlP:
			if m.cursor > 0 {
				m.cursor--
			}

		case tea.KeyDown, tea.KeyCtrlN:
			if m.cursor < len(m.filteredBranches)-1 {
				m.cursor++
			}

		case tea.KeyEnter:

			if len(m.filteredBranches) == 0 {
				return m, nil
			}

			selectedBranch := m.filteredBranches[m.cursor]
			if selectedBranch == m.currentBranch {
				return m, tea.Quit
			}

			cmd := exec.Command("git", "switch", selectedBranch)
			err := cmd.Run()
			if err != nil {
				m.err = err
				return m, nil
			}

			fmt.Printf("  ▶ Switched to branch '%s'\n", selectedBranch)
			return m, tea.Quit
		}

		m.textInput, cmd = m.textInput.Update(msg)
		m.filter()

		// Adjust cursor if out of bounds after filtering
		if m.cursor >= len(m.filteredBranches) && len(m.filteredBranches) > 0 {
			m.cursor = len(m.filteredBranches) - 1
		}
	}

	return m, cmd
}

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v\n", m.err)
	}

	s := strings.Builder{}

	s.WriteString(promptStyle.Render("$ swift git"))
	s.WriteString("\n")
	s.WriteString(searchStyle.Render(m.textInput.View()))
	s.WriteString("\n")

	if len(m.filteredBranches) > 0 {
		s.WriteString(helpTextStyle.Render("↑↓ quick select"))
		s.WriteString("\n\n")
	}

	for i, branch := range m.filteredBranches {
		var branchText string
		num := fmt.Sprintf("%d ", i)

		if branch == m.currentBranch {
			branchText = currentBranchStyle.Render(branch)
		} else {
			branchText = branchStyle.Render(branch)
		}

		if i == m.cursor {
			branchText = selectedStyle.Render(branch)
			num = selectedStyle.Render(num)
		} else {
			num = numberStyle.Render(num)
		}
		s.WriteString(fmt.Sprintf("%s%s", num, branchText))
		s.WriteString("\n")
	}
	return s.String()
}

func getBranches() ([]string, string, error) {
	cmd := exec.Command("git", "branch")
	output, err := cmd.Output()
	if err != nil {
		return nil, "", fmt.Errorf("error getting branches: %w", err)
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	branches := make([]string, 0, len(lines))
	currentBranch := ""

	for _, line := range lines {
		line = strings.TrimSpace(line)
		isCurrent := strings.HasPrefix(line, "*")
		name := strings.TrimSpace(strings.TrimPrefix(line, "*"))

		if isCurrent {
			currentBranch = name
		}

		branches = append(branches, name)
	}

	return branches, currentBranch, nil
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Search"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	branches, currentBranch, err := getBranches()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	m := model{
		branches:         branches,
		filteredBranches: branches,
		currentBranch:    currentBranch,
		textInput:        ti,
	}
	return m
}

func main() {
	if len(os.Args) > 1 {
		directSwitch(os.Args[1])
	}

	m := initialModel()

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}
