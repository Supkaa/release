package tag

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Tag struct {
	Major int
	Minor int
	Patch int
}

func Parse(tag string) Tag {
	rule := `^v\d+\.\d+\.\d+$`
	re := regexp.MustCompile(rule)
	tag = strings.TrimSpace(tag)

	if re.MatchString(tag) {
		return parseTagWithPreffix(tag)
	}

	return parseTagWithoutPreffix(tag)
}

func (t *Tag) String() string {
	return fmt.Sprintf("v%d.%d.%d", t.Major, t.Minor, t.Patch)
}

func (t *Tag) NextMajor() {
	t.Major++
	t.Minor = 0
	t.Patch = 0
}

func (t *Tag) NextMinor() {
	t.Minor++
	t.Patch = 0
}

func (t *Tag) NextPatch() {
	t.Patch++
}

func (t *Tag) IsZero() bool {
	return t.Major == 0 && t.Minor == 0 && t.Patch == 0
}

func parseTagWithPreffix(tag string) Tag {
	rule := `^v\d+\.\d+\.\d+$`
	re := regexp.MustCompile(rule)
	if !re.MatchString(tag) {
		return Tag{}
	}

	return parseTagWithoutPreffix(tag[1:])
}

func parseTagWithoutPreffix(tag string) Tag {
	rule := `^\d+\.\d+\.\d+$`
	re := regexp.MustCompile(rule)
	if !re.MatchString(tag) {
		return Tag{}
	}

	tagParts := strings.Split(tag, ".")
	major, _ := strconv.Atoi(tagParts[0])
	minor, _ := strconv.Atoi(tagParts[1])
	patch, _ := strconv.Atoi(tagParts[2])

	return Tag{
		Major: major,
		Minor: minor,
		Patch: patch,
	}
}
