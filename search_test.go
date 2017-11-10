package find

import (
	"testing"
)

func TestByName(t *testing.T) {
	data := []struct {
		pattern string
		root    string
		count   int
		ok      bool
	}{
		{"search.go", ".", 1, true},
		{"search*", "", 2, true}, // no directory means "."
		{"x", ".", 0, false},
		{"whatever", "nosuchdir", 0, false},
		{"", ".", 0, false},
	}
	for _, d := range data {
		result, err := ByName(d.pattern, d.root)
		if d.ok && err != nil {
			t.Errorf("ByName(%q, %q): %s", d.pattern, d.root, err)
		}
		if len(result) != d.count {
			t.Errorf("ByName(%q, %q) expected to find %v files, found %v", d.pattern, d.root, d.count, len(result))
		}
	}
}
