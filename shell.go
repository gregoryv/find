package find

import (
	"path/filepath"
)

type Matcher interface {
	Match(path string) (bool, error)
}

type shellPattern struct {
	pattern string
}

func NewShellPattern(pattern string) Matcher {
	return &shellPattern{pattern: pattern}
}

func (sp *shellPattern) Match(path string) (bool, error) {
	return filepath.Match(sp.pattern, path)
}
