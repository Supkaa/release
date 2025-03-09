package git

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/Supkaa/release/internal/commit"
	"github.com/Supkaa/release/internal/tag"
)

func GetLatestCommit() (commit.Commit, error) {
	cmd := exec.Command(
		"git",
		"log",
		"-1",
		`--pretty=format:%s`,
	)
	stdout, err := cmd.Output()
	if err != nil {
		return commit.Commit{}, fmt.Errorf("fail to get latest commit:%w", err)
	}

	return commit.Parse(string(stdout)), nil
}

func GetLatestTag() tag.Tag {
	cmd := exec.Command(
		"git",
		"describe",
		"--tags",
		`--abbrev=0`,
	)
	stdout, err := cmd.Output()
	if err != nil {
		return tag.Tag{}
	}

	return tag.Parse(string(stdout))
}

func GetAllCommits() []commit.Commit {
	cmd := exec.Command(
		"git",
		"log",
		"--reverse",
		`--pretty=format:%s`,
	)
	stdout, err := cmd.Output()
	if err != nil {
		log.Printf("fail to execute command: %s", err.Error())
		return []commit.Commit{}
	}

	return parseCommits(string(stdout))
}

func GetCommitsSinceTag(tag tag.Tag) []commit.Commit {
	cmd := exec.Command(
		"git",
		"log",
		fmt.Sprintf("%s..HEAD", tag.String()),
		"--reverse",
		`--pretty=format:%s`,
	)
	stdout, err := cmd.Output()
	if err != nil {
		log.Printf("fail to execute command: %s", err.Error())
		return []commit.Commit{}
	}

	return parseCommits(string(stdout))
}

func parseCommits(commitsString string) []commit.Commit {
	if commitsString == "" {
		return nil
	}

	commitStrings := strings.Split(commitsString, "\n")

	var commits []commit.Commit
	for _, commitString := range commitStrings {
		commits = append(commits, commit.Parse(commitString))
	}

	return commits
}

func IsGitAddExecuted() bool {
	cmd := exec.Command(
		"git",
		"diff",
		"--cached",
		"--exit-code",
	)
	err := cmd.Run()

	return err != nil
}

func Commit(commit commit.Commit) {
	cmd := exec.Command(
		"git",
		"commit",
		"-m",
		commit.String(),
	)

	if err := cmd.Run(); err != nil {
		log.Fatalf("fail to execute command: %s", err.Error())
	}
}

func Tag(tag tag.Tag) {
	cmd := exec.Command(
		"git",
		"tag",
		tag.String(),
	)

	if err := cmd.Run(); err != nil {
		log.Fatalf("fail to execute command: %s", err.Error())
	}
}
