package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"

	"github.com/ddddami/git-swift/internal/git"
	"github.com/ddddami/git-swift/internal/utils"
)

type clearScreenMsg struct{}

func ClearLines(count int) {
	if count <= 0 {
		return
	}

	// Move cursor up 'count' lines
	fmt.Printf("\033[%dA", count-1)
	// Clear from cursor to the end of screen
	fmt.Print("\033[J")
}

type Model struct {
	branches         []string
	filteredBranches []string
	currentBranch    string
	textInput        textinput.Model
	cursor           int
	err              error
	lineCount        int
	switchedBranch   string
	alreadyOnBranch  bool

	showDeleteConfirm  bool
	deleteTargetBranch string
	deleteErrorMsg     string
	deleteSelectedOpt  int
	deletedBranch      string
}

func NewModel(branches []string, currentBranch string, initialQuery string) Model {
	ti := textinput.New()
	ti.Placeholder = "Search"
	ti.Focus()
	ti.Prompt = " "
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
		lineCount:        0,
		switchedBranch:   "",
		alreadyOnBranch:  false,
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

type updateLineCountMsg struct {
	count int
}

func (m Model) switchToBranch(branchIndex int) (tea.Model, tea.Cmd) {
	if branchIndex < 0 || branchIndex >= len(m.filteredBranches) {
		return m, nil
	}

	selectedBranch := m.filteredBranches[branchIndex]
	if selectedBranch == m.currentBranch {
		m.alreadyOnBranch = true
		return m, tea.Quit
	}

	err := git.SwitchBranch(selectedBranch)
	if err != nil {
		m.err = err
		return m, nil
	}

	m.switchedBranch = selectedBranch
	return m, tea.Quit
}

type errorQuitMsg struct {
	error string
}

func (m Model) deleteBranch(branchIndex int) (tea.Model, tea.Cmd) {
	if branchIndex < 0 || branchIndex >= len(m.filteredBranches) {
		return m, nil
	}

	selectedBranch := m.filteredBranches[branchIndex]

	if selectedBranch == m.currentBranch {
		return m, tea.Sequence(
			func() tea.Msg {
				return errorQuitMsg{fmt.Sprintf("Cannot delete the current branch '%s'", selectedBranch)}
			},
			tea.Quit,
		)
	}

	err := git.DeleteBranch(selectedBranch)
	if err != nil {
		if unmergedErr, ok := err.(*git.UnmergedBranchError); ok {
			// Show confirmation dialog
			m.showDeleteConfirm = true
			m.deleteTargetBranch = selectedBranch
			m.deleteErrorMsg = unmergedErr.Message
			m.deleteSelectedOpt = 1 // Default to Cancel
			return m, nil
		}
		m.err = err
		return m, tea.Sequence(
			func() tea.Msg { return errorQuitMsg{err.Error()} },
			tea.Quit,
		)
	}

	m.deletedBranch = selectedBranch
	return m, tea.Quit
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	if m.showDeleteConfirm {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyCtrlC, tea.KeyEsc:
				m.showDeleteConfirm = false
				return m, nil
			case tea.KeyLeft, tea.KeyRight:
				m.deleteSelectedOpt = 1 - m.deleteSelectedOpt
			case tea.KeyEnter:
				if m.deleteSelectedOpt == 0 {
					err := git.ForceDeleteBranch(m.deleteTargetBranch)
					if err != nil {
						m.showDeleteConfirm = false
						m.err = fmt.Errorf("failed to force delete: %s", err.Error())
						return m, nil
					}

					m.deletedBranch = m.deleteTargetBranch
					m.showDeleteConfirm = false
					return m, tea.Quit
				} else {
					// Cancel
					m.showDeleteConfirm = false
					return m, nil
				}
			}
		}
		return m, nil
	}
	switch msg := msg.(type) {
	case updateLineCountMsg:
		m.lineCount = msg.count
		return m, nil

	case tea.KeyMsg:

		if msg.Alt && msg.Type == tea.KeyRunes && len(msg.Runes) == 1 {
			r := msg.Runes[0]
			if r >= '0' && r <= '9' {
				index := int(r - '0')

				return m.switchToBranch(index)
			}
		}
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

			return m.switchToBranch(m.cursor)

		case tea.KeyDelete:

			if len(m.filteredBranches) == 0 {
				return m, nil
			}

			return m.deleteBranch(m.cursor)

		}

		m.textInput, cmd = m.textInput.Update(msg)
		m.filter()
		if m.cursor >= len(m.filteredBranches) && len(m.filteredBranches) > 0 {
			m.cursor = len(m.filteredBranches) - 1
		}

	case errorQuitMsg:
		m.err = fmt.Errorf(msg.error)
		return m, nil

	case clearScreenMsg:
		// This is handled in Run() after the program exits
	}

	viewOutput := m.View()
	lineCount := strings.Count(viewOutput, "\n") + 1

	return m, tea.Batch(
		cmd,
		func() tea.Msg {
			return updateLineCountMsg{count: lineCount}
		},
	)
}

func (m Model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v\n", m.err)
	}

	if m.showDeleteConfirm {
		return m.renderDeleteConfirm()
	}

	s := strings.Builder{}

	inputView := SearchStyle.Render(m.textInput.View())
	helpText := ""
	if len(m.filteredBranches) > 0 {
		helpText = HelpTextStyle.Render("↑↓ quick select • Alt+n quick switch ")
	}
	termWidth, _, _ := term.GetSize(0)
	if termWidth == 0 {
		termWidth = 80
	}
	row := HorizontalLayout(inputView, helpText, termWidth)
	s.WriteString(row)
	s.WriteString("\n")

	if len(m.filteredBranches) == 0 {
		s.WriteString(" No matching branches\n")
	} else {
		for i, branch := range m.filteredBranches {
			branchText := branch
			num := fmt.Sprintf("%d ", i)
			if i > 9 {
				num = "  "
			}

			if i == m.cursor {
				num = SelectedStyle.Render(num)

				if branch == m.currentBranch {
					branchText = SelectedStyle.Render(fmt.Sprintf("%s *", branch))
				} else {
					branchText = SelectedStyle.Render(branchText)
				}
			} else {
				num = NumberStyle.Render(num)

				if branch == m.currentBranch {
					branchText = CurrentBranchStyle.Render(fmt.Sprintf("%s *", branch))
				} else {
					branchText = BranchStyle.Render(branchText)
				}
			}

			s.WriteString(fmt.Sprintf(" %s%s", num, branchText))
			s.WriteString("\n")
		}
	}

	return s.String()
}

func (m Model) renderDeleteConfirm() string {
	var s strings.Builder

	s.WriteString(fmt.Sprintf("  Delete Branch: %s\n\n", m.deleteTargetBranch))

	s.WriteString(fmt.Sprintf("\033[31m%s\033[0m\n\n", m.deleteErrorMsg))

	s.WriteString("\033[33mThis branch has unmerged changes!\033[0m\n")
	s.WriteString("Force delete will permanently remove all unmerged commits.\n\n")

	forceText := "Force Delete"
	cancelText := "Cancel"

	if m.deleteSelectedOpt == 0 {
		forceText = fmt.Sprintf("\033[41m %s \033[0m", forceText) // Red bg
	} else {
		cancelText = fmt.Sprintf("\033[46m %s \033[0m", cancelText) // Cyan bg
	}

	s.WriteString(fmt.Sprintf("Do you want to force delete this branch?\n\n%s  %s\n\n", forceText, cancelText))
	s.WriteString("← → navigate • Enter select • Esc cancel")

	return s.String()
}

func Run(branches []string, currentBranch string, initialQuery string) error {
	m := NewModel(branches, currentBranch, initialQuery)
	p := tea.NewProgram(m)

	var switchedBranch string
	var alreadyOnBranch bool
	var deletedBranch string
	var errorMsg string

	finalModel, err := p.Run()
	if err != nil {
		return err
	}

	if fm, ok := finalModel.(Model); ok {
		if fm.lineCount > 0 {
			ClearLines(fm.lineCount)
		}
		switchedBranch = fm.switchedBranch
		alreadyOnBranch = fm.alreadyOnBranch
		deletedBranch = fm.deletedBranch
		if fm.err != nil {
			errorMsg = fm.err.Error()
		}
	}

	if errorMsg != "" {
		fmt.Printf("\n  Error: %s\n", errorMsg)
	} else if switchedBranch != "" {
		fmt.Printf("\n  ▶ Switched to branch '%s'\n", switchedBranch)
	} else if alreadyOnBranch {
		fmt.Printf("\n  ▶ Already on branch '%s'\n", currentBranch)
	} else if deletedBranch != "" {
		fmt.Printf("\n  Branch '%s' deleted successfully\n", deletedBranch)
	}

	return nil
}
