package find

import (
	"container/list"
	"os"
	"path/filepath"
)

// ByName returns a list of files and directories whose names match the shell like pattern
func ByName(pattern, root string) (result *list.List, err error) {
	sp := NewShellPattern(pattern)
	return By(sp, root)
}

// By returns a list of files and directories whose names match
func By(m Matcher, root string) (result *list.List, err error) {
	if root == "" {
		root = "."
	}
	result = list.New()
	visit := func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !f.IsDir() {
			matched := m.Match(f.Name())
			if matched {
				result.PushBack(path)
			}
		}
		return nil
	}
	err = filepath.Walk(root, visit)
	return
}
