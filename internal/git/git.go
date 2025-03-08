package git

import (
	"log"
	"os/exec"
)

func GetLatestTag() string {
	cmd := exec.Command(
		"git",
		"log",
		"--reverse",
		`--pretty=format:%s`,
	)
	stdout, err := cmd.Output()

	if err != nil {
		log.Printf("fail to execute command: %w", err)
		return ""
	}

	return string(stdout)
}
