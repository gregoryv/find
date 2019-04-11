package find_test

import (
	"container/list"
	"github.com/gregoryv/find"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/gregoryv/asserter"
)

func TestInFile(t *testing.T) {
	file, _ := ioutil.TempFile(os.TempDir(), "grep_test")
	ioutil.WriteFile(file.Name(), []byte(`
a hello
a world
`), 0644)
	cases := []struct {
		pattern, file string
		exp           string
		expErr        bool
	}{
		{"a hello", file.Name(), "1:a hello", true},
		{"a*", file.Name(), "1:a hello,2:a world", true},
		{"*", "nosuchfile", "", false},
	}

	assert := asserter.New(t)
	for _, c := range cases {
		res, err := find.InFile(c.pattern, c.file)
		result := asLine(res)
		failed := err == nil
		assert(c.expErr == failed).Error(err)
		assert().Equals(result, c.exp)
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
