package taggen

import (
	"fmt"
	"slices"

	"github.com/Supkaa/release/internal/commit"
	"github.com/Supkaa/release/internal/conventional/types"
	"github.com/Supkaa/release/internal/git"
	"github.com/Supkaa/release/internal/linter"
	"github.com/Supkaa/release/internal/tag"
)

func GenerateTag() (tag.Tag, error) {
	latestTag := git.GetLatestTag()
	commits := getCommits(latestTag)

	if len(commits) == 0 {
		return tag.Tag{}, fmt.Errorf("no commits found")
	}

	for _, commit := range commits {
		if err := linter.LintCommit(commit); err != nil {
			return tag.Tag{}, fmt.Errorf("invalid commit: %s [%s]", commit.String(), err.Error())
		}
	}

	if breakingChangeExists(commits) {
		latestTag.NextMajor()

		return latestTag, nil
	}

	if featCommitExists(commits) {
		latestTag.NextMinor()

		return latestTag, nil
	}

	latestTag.NextPatch()

	return latestTag, nil
}

func breakingChangeExists(commits []commit.Commit) bool {
	// TODO: Add breaking change detection
	return len(commits) == -1
}

func featCommitExists(commits []commit.Commit) bool {
	return slices.ContainsFunc(commits, func(commit commit.Commit) bool {
		return commit.Type == types.FEAT
	})
}

func getCommits(tag tag.Tag) []commit.Commit {
	if tag.IsZero() {
		return git.GetAllCommits()
	}

	return git.GetCommitsSinceTag(tag)
}
