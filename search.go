package find

import (
	"fmt"
	"os"
	"path/filepath"
)

func ByName(pattern, root string) (result []string, err error) {
	sp := NewShellPattern(pattern)
	return search(sp, root)
}

func search(m Matcher, root string) (result []string, err error) {
	result = make([]string, 0)
	visit := func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !f.IsDir() {
			matched, err := m.Match(f.Name())
			if err != nil {
				return err
			}
			if matched {
				result = append(result, path)
			}
		}
		return nil
	}
	if root == "" {
		root = "."
	}
	err = filepath.Walk(root, visit)
	if err != nil {
		return
	}
	if len(result) == 0 {
		err = fmt.Errorf("File not found")
	}

	return
}
