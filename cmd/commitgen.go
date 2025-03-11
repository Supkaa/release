/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/Supkaa/release/internal/commit"
	"github.com/Supkaa/release/internal/conventional/types"
	"github.com/Supkaa/release/internal/git"
	"github.com/Supkaa/release/internal/linter"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

// commitgenCmd represents the commitgen command
var commitgenCmd = &cobra.Command{
	Use:   "commitgen",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !git.IsGitAddExecuted() {
			log.Fatal("git add not executed")
		}

		dryRun, err := cmd.Flags().GetBool("dry-run")
		if err != nil {
			log.Fatal(err)
		}

		p := tea.NewProgram(initialModel(dryRun))
		if _, err := p.Run(); err != nil {
			log.Fatalf("Error running program: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(commitgenCmd)
	commitgenCmd.Flags().BoolP("dry-run", "d", false, "Run without committing")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// commitgenCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// commitgenCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

type stage int

const (
	stageSelectType stage = iota
	stageEnterScope
	stageEnterMessage
	stageEnterDescription
	stageIsBreakingChange
	stageEnterBreakingChangeDescription
	stageConfirm
	stageError
)

type model struct {
	stage                          stage
	cursor                         int
	scopeInput                     textinput.Model
	messageInput                   textinput.Model
	descriptionInput               textinput.Model
	breakingChangeDescriptionInput textinput.Model
	errorMessage                   string
	newCommit                      commit.Commit
	isDryRun                       bool
}

func initialModel(idDryRun bool) model {
	scopeInput := textinput.New()
	scopeInput.Placeholder = "Scope"
	scopeInput.Focus()

	messageInput := textinput.New()
	messageInput.Placeholder = "Short message"

	descriptionInput := textinput.New()
	descriptionInput.Placeholder = "Long description"

	breakingChangeDescriptionInput := textinput.New()
	breakingChangeDescriptionInput.Placeholder = "Breaking change description"

	return model{
		stage:                          stageSelectType,
		cursor:                         0,
		scopeInput:                     scopeInput,
		messageInput:                   messageInput,
		descriptionInput:               descriptionInput,
		breakingChangeDescriptionInput: breakingChangeDescriptionInput,
		errorMessage:                   "",
		newCommit:                      commit.Commit{},
		isDryRun:                       idDryRun,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.stage {
		case stageSelectType:
			return m.updateSelectType(msg)
		case stageEnterScope:
			return m.updateEnterScope(msg)
		case stageEnterMessage:
			return m.updateEnterMessage(msg)
		case stageEnterDescription:
			return m.updateEnterDescription(msg)
		case stageIsBreakingChange:
			return m.updateIsBreakingChange(msg)
		case stageEnterBreakingChangeDescription:
			return m.updateEnterBreakingChangeDescription(msg)
		case stageConfirm:
			return m.updateConfirm(msg)
		case stageError:
			return m.updateError()
		}
	}

	return m, nil
}

func (m model) View() string {
	switch m.stage {
	case stageSelectType:
		return m.viewSelectType()
	case stageEnterScope:
		return m.viewEnterScope()
	case stageEnterMessage:
		return m.viewEnterMessage()
	case stageEnterDescription:
		return m.viewEnterDescription()
	case stageIsBreakingChange:
		return m.viewIsBreakingChange()
	case stageEnterBreakingChangeDescription:
		return m.viewEnterBreakingChangeDescription()
	case stageConfirm:
		return m.viewConfirm()
	case stageError:
		return m.viewError()
	}

	return "Unknown stage"
}

func (m model) updateSelectType(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < len(types.ValidCommitTypes)-1 {
			m.cursor++
		}
	case "enter":
		m.newCommit.Type = types.ValidCommitTypes[m.cursor]
		m.stage = stageEnterScope
	}
	return m, nil
}

func (m model) viewSelectType() string {
	s := "Select commit type:\n\n"
	for i, t := range types.ValidCommitTypes {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		s += fmt.Sprintf("%s %s\n", cursor, t)
	}
	s += "\nUse arrow keys or j/k to navigate, Enter to select, ctrl+c to quit"
	return s
}

func (m model) updateEnterScope(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "enter":
		m.newCommit.Scope = strings.TrimSpace(m.scopeInput.Value())
		m.stage = stageEnterMessage
		m.scopeInput.Blur()
		m.messageInput.Focus()
		return m, textinput.Blink
	}

	m.scopeInput, cmd = m.scopeInput.Update(msg)

	return m, cmd
}

func (m model) viewEnterScope() string {
	s := "Enter scope (optional):\n\n"
	s += m.scopeInput.View()
	s += "\n\nPress Enter to continue, ctrl+c to quit"
	return s
}

func (m model) updateEnterMessage(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "enter":
		m.newCommit.Message = strings.TrimSpace(m.messageInput.Value())
		m.stage = stageEnterDescription
		m.messageInput.Blur()
		m.descriptionInput.Focus()
		return m, textinput.Blink
	}

	m.messageInput, cmd = m.messageInput.Update(msg)

	return m, cmd
}

func (m model) viewEnterMessage() string {
	s := "Enter short description:\n\n"
	s += m.messageInput.View()
	s += "\n\nPress Enter to continue, ctrl+c to quit"
	return s
}

func (m model) updateEnterDescription(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "enter":
		m.newCommit.Description = strings.TrimSpace(m.descriptionInput.Value())
		m.stage = stageIsBreakingChange
		m.descriptionInput.Blur()
		m.breakingChangeDescriptionInput.Focus()
		return m, textinput.Blink
	}

	m.descriptionInput, cmd = m.descriptionInput.Update(msg)

	return m, cmd
}

func (m model) viewEnterDescription() string {
	s := "Enter long description (optional):\n\n"
	s += m.descriptionInput.View()
	s += "\n\nPress Enter to continue, ctrl+c to quit"
	return s
}

func (m model) updateIsBreakingChange(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "y", "Y":
		m.newCommit.IsBreakingChange = true
		m.stage = stageEnterBreakingChangeDescription
	case "n", "N", "enter":
		m.newCommit.IsBreakingChange = false
		m.stage = stageConfirm
	}

	return m, nil
}

func (m model) viewIsBreakingChange() string {
	s := "Is this a breaking change? (y/N):\n\n"
	s += fmt.Sprintf("Current: %v\n", m.newCommit.IsBreakingChange)
	s += "\nPress y/n to toggle, Enter to continue, ctrl+c to quit"
	return s
}

func (m model) updateEnterBreakingChangeDescription(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "enter":
		m.newCommit.BreakingChange = strings.TrimSpace(m.breakingChangeDescriptionInput.Value())
		m.stage = stageConfirm
		m.breakingChangeDescriptionInput.Blur()
		return m, textinput.Blink
	}

	m.breakingChangeDescriptionInput, cmd = m.breakingChangeDescriptionInput.Update(msg)

	return m, cmd
}

func (m model) viewEnterBreakingChangeDescription() string {
	s := "Enter breaking change description:\n\n"
	s += m.breakingChangeDescriptionInput.View()
	s += "\n\nPress Enter to continue, ctrl+c to quit"
	return s
}

func (m model) updateConfirm(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "y", "Y", "enter":
		if err := linter.LintCommit(m.newCommit); err != nil {
			m.errorMessage = fmt.Sprintf("Linting failed: %v", err)
			m.stage = stageError

			return m, nil
		}

		if !m.isDryRun {
			git.Commit(m.newCommit)
		}

		return m, tea.Quit
	case "n", "N":
		return m, tea.Quit
	}

	return m, nil
}

func (m model) viewConfirm() string {
	s := "Commit Preview:\n\n"
	s += m.newCommit.String()
	s += "\n\nProceed with commit? (y/N)"
	return s
}

func (m model) updateError() (tea.Model, tea.Cmd) {
	return m, tea.Quit
}

func (m model) viewError() string {
	s := "Commit Preview:\n\n"
	s += m.newCommit.String()
	s += "\n\n" + m.errorMessage

	return s
}
