package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#7851A9")).
			Padding(0, 1)

	inputStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#333333")).
			Padding(0, 1)

	currentBranchStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#00FF00")).
				Bold(true)

	matchStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5"))

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#7851A9")).
			Bold(true)

	statusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#FF6B6B"))
)

type Branch struct {
	Name      string
	IsCurrent bool
}

type model struct {
	branches      []Branch
	currentBranch string
	textInput     textinput.Model
	selected      int
	initialized   bool
	ready         bool
	windowSize    tea.WindowSizeMsg
	viewport      viewport.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// cmd  tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.windowSize = msg

		if !m.initialized {
			m.initialized = true

			m.viewport = viewport.New(msg.Width, msg.Height-3)
			m.viewport.SetContent("")
			m.textInput = textinput.New()
			m.textInput.Placeholder = "Search"
			// m.textInput.Focus()
			m.textInput.Width = msg.Width
			// m.textInput.Prompt = "branch › "
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - 3
		}
		m.ready = true

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	viewportContent := strings.Builder{}
	for i, branch := range m.branches {
		branchName := branch.Name

		if branch.IsCurrent {
			branchName = "* " + currentBranchStyle.Render(branchName)
		} else {
			branchName = "  " + branchName
		}

		if i == m.selected {
			viewportContent.WriteString(selectedStyle.Render("› "+branchName) + "\n")
		} else {
			viewportContent.WriteString(matchStyle.Render("  "+branchName) + "\n")
		}

		m.viewport.SetContent(viewportContent.String())

		_, viewportCmd := m.viewport.Update(msg)
		cmds = append(cmds, viewportCmd)
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if !m.ready {
		return "Loading branches..."
	}
	s := strings.Builder{}
	s.WriteString(inputStyle.Render(m.textInput.View()))
	s.WriteString("\n")
	s.WriteString(m.viewport.View())
	s.WriteString("\n")
	s.WriteString(" [↑/↓] Navigate [Enter] Switch [Esc] Quit")
	return s.String()
}

func getBranches() ([]Branch, string, error) {
	cmd := exec.Command("git", "branch")
	output, err := cmd.Output()
	if err != nil {
		return nil, "", fmt.Errorf("error getting branches: %w", err)
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	branches := make([]Branch, 0, len(lines))
	currentBranch := ""

	for _, line := range lines {
		line = strings.TrimSpace(line)
		isCurrent := strings.HasPrefix(line, "*")
		name := strings.TrimSpace(strings.TrimPrefix(line, "*"))

		if isCurrent {
			currentBranch = name
		}

		branches = append(branches, Branch{
			Name:      name,
			IsCurrent: isCurrent,
		})
	}

	return branches, currentBranch, nil
}

func main() {
	branches, currentBranch, err := getBranches()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	m := model{
		branches:      branches,
		currentBranch: currentBranch,
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}
