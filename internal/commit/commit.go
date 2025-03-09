package commit

import (
	"fmt"
	"regexp"

	"github.com/Supkaa/release/internal/conventional/types"
)

type Commit struct {
	Type      types.CommitType
	Scope     string
	Message   string
	IsMerge   bool
	IsInitial bool
}

func Parse(commit string) Commit {
	newCommit := Commit{
		IsMerge:   parseIsMerge(commit),
		IsInitial: parseIsInitial(commit),
		Message:   commit,
	}

	if !newCommit.IsMerge && !newCommit.IsInitial {
		newCommit.Type = parseType(commit)
		newCommit.Scope = parseScope(commit)
		newCommit.Message = parseMessage(commit)
	}

	return newCommit
}

func (c Commit) String() string {
	str := ""
	if c.Type != "" {
		str = fmt.Sprintf("%s: ", c.Type)
	}

	if c.Scope != "" && str != "" {
		str = fmt.Sprintf("%s(%s): ", c.Type, c.Scope)
	}

	if c.Message != "" {
		str = str + c.Message
	}

	return str
}

func parseType(commit string) types.CommitType {
	matches := regexp.
		MustCompile(`([\w\s]+)[:(]`).
		FindString(commit)

	if len(matches) > 2 {
		return matches[:len(matches)-1]
	}
	return matches
}

func parseIsMerge(commit string) bool {
	return regexp.
		MustCompile(`(?i)^merge`).
		MatchString(commit)
}

func parseIsInitial(commit string) bool {
	return commit == "Initial commit"
}

func parseScope(commit string) string {
	matches := regexp.
		MustCompile(`\(([^)]+)\)\:`).
		FindString(commit)

	if len(matches) > 2 {
		return matches[1 : len(matches)-2]
	}

	return ""
}

func parseMessage(commit string) string {
	matches := regexp.
		MustCompile(`:\ (.*)$`).
		FindString(commit)

	if len(matches) > 2 {
		return matches[2:]
	}

	return ""
}
