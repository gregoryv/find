package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"
)

func NewScanner() *Scanner {
	return &Scanner{
		files:  make([]string, 0),
		Logger: log.New(ioutil.Discard, "", log.LstdFlags),
	}
}

type Scanner struct {
	files []string
	*log.Logger

	result []FileMatch
}

func (me *Scanner) Scan(expr string) error {
	re, err := regexp.Compile(expr)
	if err != nil {
		return err
	}
	me.result = make([]FileMatch, 0)

	for _, filename := range me.Files() {
		me.Println("scan:", filename)
		fh, err := os.Open(filename)
		if err != nil {
			return err
		}

		s := bufio.NewScanner(fh)
		var line int
		fm := FileMatch{
			Filename: filename,
			Result:   make([]LineMatch, 0),
		}
		for s.Scan() {
			line++ // lines start with 1
			text := s.Text()
			if re.MatchString(text) {
				fm.Result = append(fm.Result, LineMatch{
					Line: line,
					Text: text,
				})
			}
		}
		fh.Close()
		if len(fm.Result) > 0 {
			me.result = append(me.result, fm)
		}
	}
	return nil
}

func (me *Scanner) LastResult() []FileMatch {
	return me.result
}

func (me *Scanner) WriteResult(w io.Writer) {
	var i int
	for _, fm := range me.result {
		fmt.Fprintln(w, fm.Filename)
		for _, m := range fm.Result {
			fmt.Fprintln(w, i+1, m.Text)
			i++
		}
		fmt.Fprintln(w)
	}
}

// Files returns a filtered list of files
func (me *Scanner) Files() []string {
	// todo
	return me.files
}

// SetFiles to scan
func (me *Scanner) SetFiles(v []string) { me.files = v }

type FileMatch struct {
	Filename string
	Result   []LineMatch
}

type LineMatch struct {
	Line int
	Text string
}
