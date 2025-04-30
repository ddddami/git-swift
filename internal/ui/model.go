package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/ddddami/swift-git/internal/git"
	"github.com/ddddami/swift-git/internal/utils"
)

type Model struct {
	branches         []string
	filteredBranches []string
	currentBranch    string
	textInput        textinput.Model
	cursor           int
	err              error
}

func NewModel(branches []string, currentBranch string, initialQuery string) Model {
	ti := textinput.New()
	ti.Placeholder = "Search"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20
	ti.PromptStyle = lipgloss.NewStyle()

	if initialQuery != "" {
		ti.SetValue(initialQuery)
	}

	m := Model{
		branches:         branches,
		filteredBranches: branches,
		currentBranch:    currentBranch,
		textInput:        ti,
	}

	if initialQuery != "" {
		m.filter()
	}

	return m
}

func (m *Model) filter() {
	query := m.textInput.Value()

	if query == "" {
		m.filteredBranches = m.branches
		return
	}

	filtered := []string{}
	for _, branch := range m.branches {
		if utils.FuzzyMatch(branch, query) {
			filtered = append(filtered, branch)
		}
	}

	m.filteredBranches = filtered
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
				fmt.Printf("Already on branch '%s'\n", selectedBranch)
				return m, tea.Quit
			}

			err := git.SwitchBranch(selectedBranch)
			if err != nil {
				m.err = err
				return m, nil
			}

			fmt.Printf("  ▶ Switched to branch '%s'\n", selectedBranch)
			return m, tea.Quit
		}

		m.textInput, cmd = m.textInput.Update(msg)
		m.filter()

		if m.cursor >= len(m.filteredBranches) && len(m.filteredBranches) > 0 {
			m.cursor = len(m.filteredBranches) - 1
		}
	}

	return m, cmd
}

func (m Model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v\n", m.err)
	}

	s := strings.Builder{}

	s.WriteString(PromptStyle.Render("$ swift git"))
	s.WriteString("\n")
	s.WriteString(SearchStyle.Render(m.textInput.View()))
	s.WriteString("\n")

	if len(m.filteredBranches) > 0 {
		s.WriteString(HelpTextStyle.Render("↑↓ quick select"))
		s.WriteString("\n\n")
	}

	for i, branch := range m.filteredBranches {
		var branchText string
		num := fmt.Sprintf("%d ", i)

		if branch == m.currentBranch {
			branchText = CurrentBranchStyle.Render(branch)
		} else {
			branchText = BranchStyle.Render(branch)
		}

		if i == m.cursor {
			branchText = SelectedStyle.Render(branch)
			num = SelectedStyle.Render(num)
		} else {
			num = NumberStyle.Render(num)
		}
		s.WriteString(fmt.Sprintf("%s%s", num, branchText))
		s.WriteString("\n")
	}

	return s.String()
}

func Run(branches []string, currentBranch string, initialQuery string) error {
	m := NewModel(branches, currentBranch, initialQuery)
	p := tea.NewProgram(m)
	_, err := p.Run()
	return err
}
