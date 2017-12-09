package find

import (
	"regexp"
)

type reg struct {
	ex *regexp.Regexp
}

func NewRegexp(ex *regexp.Regexp) Matcher {
	return &reg{ex: ex}
}

func (rm *reg) Match(path string) bool {
	return rm.ex.Match([]byte(path))
}
