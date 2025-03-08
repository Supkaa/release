package git

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
	"slices"
	"strings"

	"github.com/Supkaa/release/internal/commit"
	"github.com/Supkaa/release/internal/conventional/types"
)

var (
	ErrCommitTypeIsNotProvided = errors.New("commit type is not provided")
	ErrInvalidCommitType       = errors.New("invalid commit type")
	ErrCommitTypeCannotBeEmpty = errors.New("commit type can`t be empty")
)

func LintCommit(commit commit.Commit) error {
	log.Printf("%#v", commit)
	if commit.IsMerge {
		return nil
	}

	if err := lintType(commit.Type); err != nil {
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

func GetLatestCommit() string {
	cmd := exec.Command(
		"git",
		"log",
		"-1",
		`--pretty=format:%s`,
	)
	stdout, err := cmd.Output()
	if err != nil {
		log.Printf("fail to execute command: %s", err.Error())
		return ""
	}

	return string(stdout)
}

func GetLatestTag() string {
	cmd := exec.Command(
		"git",
		"log",
		"--reverse",
		`--pretty=format:%s`,
	)
	stdout, err := cmd.Output()
	if err != nil {
		log.Printf("fail to execute command: %s", err.Error())
		return ""
	}

	return string(stdout)
}
