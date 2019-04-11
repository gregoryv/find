package find_test

import (
	"github.com/gregoryv/find"
	"io/ioutil"
	"container/list"
	"os"
	"strings"
	"testing"
)

func TestInFile(t *testing.T) {
	file, _ := ioutil.TempFile(os.TempDir(), "grep_test")
	ioutil.WriteFile(file.Name(), []byte(`
a hello
a world
`), 0644)
	data := []struct {
		pattern, file string
		exp           string
		expErr            bool
	}{
		{"a hello", file.Name(), "1:a hello", true},
		{"a*", file.Name(), "1:a hello,2:a world", true},
		{"*", "nosuchfile", "", false},
	}

	for _, d := range data {
		res, err := find.InFile(d.pattern, d.file)
		result := asLine(res)
		// Assert
		if d.expErr != (err == nil) {
			t.Error(err)
		}
		if result != d.exp {
			t.Errorf("Grep(%q, %q) expected \n%v\n, got\n %v", d.pattern, d.file, d.exp, result)
		}
	}
}

func asLine(res *list.List) string {
	if res == nil {
		return ""
	}
	lines := make([]string, 0, res.Len())
	for e := res.Front(); e != nil; e = e.Next() {
		if ref, ok := e.Value.(*find.Ref); ok {
			lines = append(lines, ref.String())
		}
	}
	return strings.Join(lines, ",")
}
