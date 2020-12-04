package harvest

import (
	"fmt"
	"testing"
)

func TestTagEntryMatcher(t *testing.T) {
	m := NewTagEntryMatcher("tag0,tag1,tag2")

	matched, tag := m.Match("text that includes tag1.")

	if want := true; matched != want {
		t.Errorf("Match() is %+v, want %+v", matched, want)
	}

	if want := "tag1"; tag != want {
		t.Errorf("Match() is %+v, want %+v", tag, want)
	}

	matched, tag = m.Match("text that does not include tags")
	if want := false; matched != want {
		t.Errorf("Match() is %+v, want %+v", matched, want)
	}

	if want := ""; tag != want {
		t.Errorf("Match() is %+v, want %+v", tag, want)
	}
}

func TestPatternEntryMatcher(t *testing.T) {
	var want error
	m, err := NewPatternEntryMatcher("\\[DS-\\d+\\]")
	fmt.Printf("XXX err = %v\n", err)

	if want = nil; err != want {
		t.Errorf("Match() is %+v, want %+v", err, want)
		return
	}

	matched, tag := m.Match("text that includes a pattern [DS-357] indeed")

	if want := true; matched != want {
		t.Errorf("Match() is %+v, want %+v", matched, want)
	}

	if want := "[DS-357]"; tag != want {
		t.Errorf("Match() is %+v, want %+v", tag, want)
	}
}
