package find

import (
	"path/filepath"
)

type Matcher interface {
	Match(path string) bool
}

type shellPattern struct {
	pattern string
}

func NewShellPattern(pattern string) Matcher {
	return &shellPattern{pattern: pattern}
}

func (sp *shellPattern) Match(path string) bool {
	res, _ := filepath.Match(sp.pattern, path)
	return res
}
