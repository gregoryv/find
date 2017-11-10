package find

import (
	"testing"
)

func TestFile(t *testing.T) {
	data := []struct {
		pattern string
		root    string
		count int
		ok      bool

	}{
		{"file.go", ".", 1, true},
		{"file*", "", 2, true},
		{"x", ".", 0, false},
	}
	for _, d := range data {
		result, err := File(d.pattern, d.root)
		if d.ok && err != nil {
			t.Errorf("File(%q, %q): %s", d.pattern, d.root, err)
		}
		if len(result) != d.count {
			t.Errorf("File(%q, %q) expected to find %v files, found %v", d.pattern, d.root, d.count, len(result))
		}
	}

}
