package linter

import (
	"testing"

	"github.com/Supkaa/release/internal/commit"
	"github.com/Supkaa/release/internal/conventional/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLintCommit_MergeCommit(t *testing.T) {
	commit := commit.Commit{
		IsMerge: true,
	}

	err := LintCommit(commit)

	require.NoError(t, err)
}

func TestLintCommit_EmptyCommitType(t *testing.T) {
	commit := commit.Commit{
		Type: "",
	}

	err := LintCommit(commit)

	assert.ErrorIs(t, err, ErrCommitTypeCannotBeEmpty)
}

func TestLintCommit_InvalidCommitType(t *testing.T) {
	commit := commit.Commit{
		Type: types.CommitType(" invalid"),
	}

	err := LintCommit(commit)

	assert.ErrorIs(t, err, ErrInvalidCommitType)
}

func TestLintCommit_InvalidCommitScope(t *testing.T) {
	commit := commit.Commit{
		Type:  types.TEST,
		Scope: "InvalidScope",
	}

	err := LintCommit(commit)

	assert.ErrorIs(t, err, ErrInvalidCommitScope)
}

func TestLintCommit_EmptyCommitMessage(t *testing.T) {
	commit := commit.Commit{
		Type:    types.TEST,
		Scope:   "test",
		Message: "",
	}

	err := LintCommit(commit)

	assert.ErrorIs(t, err, ErrCommitMessageCannotBeEmpty)
}

func TestLintCommit_ValidCommit(t *testing.T) {
	commit := commit.Commit{
		Type:    types.TEST,
		Scope:   "test",
		Message: "valid message",
	}

	err := LintCommit(commit)

	assert.NoError(t, err)
}
