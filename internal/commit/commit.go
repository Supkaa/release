package commit

import (
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/Supkaa/release/internal/conventional/types"
)

var (
	baseFormatRegex       = regexp.MustCompile(`(?is)^(?:(?P<type>[^\(!:]+)(?:\((?P<scope>[^\)]+)\))?(?P<breaking>!)?: (?P<message>[^\n\r]+))(?P<remainder>.*)`)
	bodyFooterFormatRegex = regexp.MustCompile(`(?isU)^(?:(?P<description>.*))?(?P<footer>(?-U:(?:[\w\-]+(?:: | #).*|(?i:BREAKING CHANGE:.*)|(?i:BREAKING-CHANGE:.*))+))`)
	footerFormatRegex     = regexp.MustCompile(`(?s)^(?P<footer>(?i:(?:[\w\-]+(?:: | #).*|(?i:BREAKING CHANGE:.*)|(?i:BREAKING-CHANGE:.*))+))`)
)

func regExMapper(match []string, expectedFormatRegex *regexp.Regexp, result map[string]string) {
	for i, name := range expectedFormatRegex.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = strings.TrimSpace(match[i])
		}
	}
}

type Commit struct {
	Type             types.CommitType
	Scope            string
	Message          string
	Description      string
	Footers          map[string]string
	IsMerge          bool
	IsInitial        bool
	IsBreakingChange bool
	BreakingChange   string
}

func Parse(commit string) Commit {
	newCommit := Commit{
		IsMerge:   parseIsMerge(commit),
		IsInitial: parseIsInitial(commit),
		Message:   commit,
	}

	if newCommit.IsMerge || newCommit.IsInitial {
		return newCommit
	}

	return deepParse(commit)
}

func deepParse(commit string) Commit {
	// isBreakingChange := false
	match := baseFormatRegex.FindStringSubmatch(commit)
	if len(match) == 0 {
		return Commit{
			Message: commit,
		}
	}

	result := make(map[string]string)
	regExMapper(match, baseFormatRegex, result)

	// split the remainder into body & footer
	match = bodyFooterFormatRegex.FindStringSubmatch(result["remainder"])
	if len(match) > 0 {
		regExMapper(match, bodyFooterFormatRegex, result)
	} else {
		result["description"] = result["remainder"]
	}

	breakingChangesText, footers := extractAndRemoveBreakingChanges(parseFooters(result["footer"]))

	return Commit{
		Type:             types.CommitType(result["type"]),
		Scope:            result["scope"],
		Message:          result["message"],
		Description:      result["description"],
		Footers:          footers,
		IsBreakingChange: result["breaking"] == "!" || breakingChangesText != "",
		BreakingChange:   breakingChangesText,
	}
}

func parseFooters(rawFooter string) map[string]string {
	if rawFooter == "" {
		return nil
	}

	footers := make(map[string]string)
	currentKey := ""

	for _, line := range strings.Split(rawFooter, "\n") {
		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}

		matches := footerFormatRegex.FindStringSubmatch(line)
		if len(matches) == 0 {
			footers[currentKey] += fmt.Sprintf("\n%s", line)
			continue
		}

		l := strings.Split(matches[0], ": ")
		key := l[0]
		value := l[1]
		currentKey = key

		if existingValue, ok := footers[key]; ok {
			value = fmt.Sprintf("%s,%s", existingValue, value)
		}

		footers[key] = strings.TrimSpace(value)
	}

	if len(footers) == 0 {
		footers = nil
	}

	return footers
}

func extractAndRemoveBreakingChanges(footers map[string]string) (string, map[string]string) {
	if footers == nil {
		return "", nil
	}

	breakingChangesText := ""

	for key, value := range footers {
		if !slices.Contains([]string{"BREAKING CHANGE", "BREAKING-CHANGE"}, key) {
			continue
		}

		breakingChangesText = fmt.Sprintf("%s\n%s", breakingChangesText, value)
		delete(footers, key)
	}

	if len(footers) == 0 {
		footers = nil
	}

	return strings.TrimSpace(breakingChangesText), footers
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

func parseIsMerge(commit string) bool {
	return regexp.
		MustCompile(`(?i)^merge`).
		MatchString(commit)
}

func parseIsInitial(commit string) bool {
	return commit == "Initial commit"
}
