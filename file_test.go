package find

import (
	"testing"
	"fmt"
)

func ExampleFile() {
	result, _ := File("*.go", ".")
	for _, file := range result {
		fmt.Println(file)
	}
	//output: file.go
	//file_test.go
}

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
