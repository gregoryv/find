package find_test

import (
	"github.com/gregoryv/find"
	"io/ioutil"
	"os"
	"testing"
)

func TestGrep(t *testing.T) {
	file, _ := ioutil.TempFile(os.TempDir(), "grep_test")
	ioutil.WriteFile(file.Name(), []byte(`
a hello
a world
`), 0644)
	data := []struct {
		pattern, file string
		exp           find.Refs
		ok            bool
	}{
		{"a hello", file.Name(), find.Refs{{1, "a hello"}}, true},
		{"a*", file.Name(), find.Refs{
			{1, "a hello"},
			{2, "a world"},
		}, true},
		{"*", "nosuchfile", find.Refs{}, false},
	}

	for _, d := range data {
		res, err := find.Grep(d.pattern, d.file)
		if d.ok && (len(res) == 0 || res[0].Line != d.exp[0].Line) {
			t.Errorf("Grep(%q, %q) expected \n%v\n, got\n %v", d.pattern, d.file, d.exp, res)
		}
		if !d.ok && err == nil {
			t.Errorf("Grep(%q, %q) expected to fail", d.pattern, d.file)
		}
	}
}
