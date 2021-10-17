package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gregoryv/cmdline"
	"github.com/gregoryv/find"
)

func main() {
	var (
		cli       = cmdline.NewBasicParser()
		files     = cli.Option("-f, --files").String("")
		expr      = cli.Required("EXPR").String("")
		openIndex = cli.Optional("OPEN_INDEX").String("")
	)
	cli.Parse()

	s := NewScanner()
	filter := &smart{}
	s.SetFiles(
		ls(files, filter),
	)

	if err := s.Scan(expr); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if openIndex == "" {
		s.WriteResult(os.Stdout)
		os.Exit(0)
	}

	oi, err := strconv.Atoi(openIndex)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	var i int
	for _, fm := range s.LastResult() {
		for _, lm := range fm.Result {
			i++
			if i == oi {
				editor := os.Getenv("EDITOR")
				// Adapt command to open on a specific line
				var cmd *exec.Cmd
				switch editor {
				case "emacs", "emacsclient", "vi", "vim":
					cmd = exec.Command(
						editor, "-n", fmt.Sprintf("+%d", lm.Line), fm.Filename,
					)
				default:
					cmd = exec.Command(editor, fm.Filename)
				}
				err := cmd.Start()
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				os.Exit(0)
			}
		}
	}

}

// ls returns a list of files based on the given pattern. Empty string means
// recursive from current working directory
func ls(pattern string, filter find.Matcher) []string {
	if pattern != "" {
		f, _ := filepath.Glob(pattern)
		return f
	}
	// recursive from current working directory
	result, _ := find.ByName("*", ".")
	files := make([]string, 0)
	for e := result.Front(); e != nil; e = e.Next() {
		filename := e.Value.(string)
		if filter.Match(filename) {
			files = append(files, filename)
		}
	}
	return files
}

type smart struct{}

// Match excludes git project files, e.g. .git/
func (me *smart) Match(path string) bool {
	switch {
	case strings.Index(path, ".git/") == 0:
	default:
		return true
	}
	return false
}
