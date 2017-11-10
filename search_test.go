package find

import (
	"fmt"
	"testing"
)

func ExampleByName() {
	result, _ := ByName("*.go", ".")
	for _, file := range result {
		fmt.Println(file)
	}
	//output:search.go
	//search_test.go
	//shell.go
}

func TestByName(t *testing.T) {
	data := []struct {
		pattern string
		root    string
		count   int
		ok      bool
	}{
		{"search.go", ".", 1, true},
		{"search*", "", 2, true},
		{"x", ".", 0, false},
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
