package linter

import (
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/Supkaa/release/internal/commit"
	"github.com/Supkaa/release/internal/conventional/types"
)

var (
	ErrInvalidCommitType          = errors.New("invalid commit type")
	ErrCommitTypeCannotBeEmpty    = errors.New("commit type can`t be empty")
	ErrInvalidCommitScope         = errors.New("invalid commit scope")
	ErrCommitMessageCannotBeEmpty = errors.New("commit message can`t be empty")
)

func LintCommit(commit commit.Commit) error {
	if commit.IsMerge {
		return nil
	}

	if err := lintType(commit.Type); err != nil {
		return fmt.Errorf("lint error: %w", err)
	}

	if commit.Scope != "" {
		if err := lintScope(commit.Scope); err != nil {
			return fmt.Errorf("lint error: %w", err)
		}
	}

	if err := lintMessage(commit.Message); err != nil {
		return fmt.Errorf("lint error: %w", err)
	}

	return nil
}

func lintType(commitType types.CommitType) error {
	if strings.TrimSpace(commitType) == "" {
		return ErrCommitTypeCannotBeEmpty
	}

	if !slices.Contains(types.ValidCommitTypes, commitType) {
		return fmt.Errorf("%w (%s). Valid types [%s]", ErrInvalidCommitType, commitType, strings.Join(types.ValidCommitTypes, ", "))
	}

	return nil
}

func lintScope(commitScope string) error {
	rule := `^[a-z]+$`
	re := regexp.MustCompile(rule)
	if !re.MatchString(commitScope) {
		return fmt.Errorf("%w (%s). Valid scope must satisfy the rule %s", ErrInvalidCommitScope, commitScope, rule)
	}

	return nil
}

func lintMessage(commitMessage string) error {
	if strings.TrimSpace(commitMessage) == "" {
		return ErrCommitMessageCannotBeEmpty
	}

	return nil
}
