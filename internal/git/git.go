package git

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/Supkaa/release/internal/commit"
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

	return commit.New(string(stdout)), nil
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
