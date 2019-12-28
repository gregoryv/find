// Package find implements search funcs for finding files by name or content
package find

import (
	"container/list"
	"os"
	"path/filepath"
)

// ByName returns a list of files whose names match the shell like pattern
func ByName(pattern, root string) (*Result, error) {
	sp := NewShellPattern(pattern)
	return By(sp, root)
}

// By returns a list of files whose names match
func By(m Matcher, root string) (*Result, error) {
	if root == "" {
		root = "."
	}
	result := &Result{list.New()}
	err := filepath.Walk(root, newVisitor(m, result))
	return result, err
}

// Returns a visitor that skips directories
func newVisitor(m Matcher, result *Result) filepath.WalkFunc {
	return func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !f.IsDir() && m.Match(f.Name()) {
			result.PushBack(path)
		}
		return nil
	}
}

type Result struct {
	*list.List
}

func (result *Result) Map(fn func(string)) {
	for e := result.Front(); e != nil; e = e.Next() {
		s, ok := e.Value.(string)
		if !ok {
			continue
		}
		fn(s)
	}
}
