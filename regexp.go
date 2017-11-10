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

func (rm *reg) Match(path string) (bool, error) {
	res := rm.ex.Match([]byte(path))
	return res, nil
}
