package harvest

import (
	"regexp"
	"strings"
)

type EntryMatcher interface {
	Match(string) (bool, string)
}

type TagEntryMatcher struct {
	tags []string
}

func NewTagEntryMatcher(tags string) *TagEntryMatcher {
	t := strings.Split(tags, ",")
	return &TagEntryMatcher{
		tags: t,
	}
}

func (m *TagEntryMatcher) Match(text string) (bool, string) {
	for _, t := range m.tags {
		if strings.Contains(text, t) {
			return true, t
		}
	}

	return false, ""
}

type PatternEntryMatcher struct {
	re *regexp.Regexp
}

func NewPatternEntryMatcher(pattern string) (*PatternEntryMatcher, error) {
	re, err := regexp.Compile(pattern)
	return &PatternEntryMatcher{
		re: re,
	}, err
}

func (m *PatternEntryMatcher) Match(text string) (bool, string) {
	match := m.re.FindString(text)
	matched := len(match) > 0
	return matched, match
}
