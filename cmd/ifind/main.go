package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gregoryv/binext"
	"github.com/gregoryv/cli"
	"github.com/gregoryv/find"
)

func main() {
	in := NewInput()
	if err := cli.Parse(in, os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}	

	if in.Help {
		WriteUsage(os.Stdout)
		return
	}

	if in.Expression == "" {
		fmt.Println("empty EXPR")
		os.Exit(1)
	}
	
	s := NewScanner()
	if in.Verbose {
		s.Logger.SetOutput(log.Writer())
	}
	filter := &smart{}
	filter.SetIncludeBinary(in.IncludeBinary)
	if err := filter.SetExclude(in.Exclude); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	s.SetFiles(
		ls(in.Files, filter),
	)

	if err := s.Scan(in.Expression); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if in.WriteAliases != "" {
		var i int
		aw, err := os.Create(in.WriteAliases)
		if err != nil {
			log.Fatal(err)
		}
		for _, fm := range s.LastResult() {
			for _, lm := range fm.Result {
				fmt.Fprintln(aw, aliasLine(i+1, in.AliasPrefix, fm, lm))
				i++
			}
		}
		aw.Close()
	}

	// results destination
	w := os.Stdout

	if in.OpenIndex == 0 { // list result
		var i int
		for _, fm := range s.LastResult() {
			fmt.Fprintln(w, fm.Filename)
			for _, m := range fm.Result {
				text := m.Text
				if in.Colors {
					colored := fmt.Sprintf("%s%s%s", green, in.Expression, reset)
					text = strings.ReplaceAll(text, in.Expression, colored)
				}
				fmt.Fprintln(w, i+1, text)
				i++
			}
			fmt.Fprintln(w)
		}
		os.Exit(0)
	}

	var i uint32
	for _, fm := range s.LastResult() {
		for _, lm := range fm.Result {
			i++
			if i != in.OpenIndex {
				continue
			}
			editor := os.Getenv("EDITOR")
			// Adapt command to open on a specific line
			var cmd *exec.Cmd
			switch editor {
			case "emacs", "emacsclient", "vi", "vim":
				cmd = exec.Command(
					editor, "-n", fmt.Sprintf("+%d", lm.Line), fm.Filename,
				)
			case "code", "Code.exe", "code.exe":
				cmd = exec.Command(
					editor, "--goto", fmt.Sprintf("%s:%d", fm.Filename, lm.Line),
				)
			default:
				cmd = exec.Command(editor, fm.Filename)
			}
			err := cmd.Start()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			return
		}
	}
	fmt.Printf("%v? there are only %v matches\n", in.OpenIndex, i)
	os.Exit(1)
}

func aliasLine(i int, prefix string, fm FileMatch, lm LineMatch) string {
	editor := os.Getenv("EDITOR")

	var cmd string
	switch editor {
	case "emacs", "emacsclient", "vi", "vim":
		cmd = fmt.Sprintf("%s -n +%d %s", editor, lm.Line, fm.Filename)
	case "code", "Code.exe", "code.exe":
		cmd = fmt.Sprintf("%s --goto %s:%d", editor, fm.Filename, lm.Line)
	}
	return fmt.Sprintf(`alias %s%v="%s"`, prefix, i, cmd)
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

type smart struct {
	includeBinary bool
	exclude       *regexp.Regexp
}

func (me *smart) SetIncludeBinary(v bool) {
	me.includeBinary = v
}

func (me *smart) SetExclude(v string) error {
	re, err := regexp.Compile(v)
	if err != nil {
		return err
	}
	me.exclude = re
	return nil
}

// Match excludes git project files, e.g. .git/
func (me *smart) Match(path string) bool {
	var executableFile bool
	if filepath.Ext(path) == "" {
		i, _ := os.Stat(path)
		executableFile = (i.Mode()&0111 != 0 && i.Mode().IsRegular())
	}

	if !me.includeBinary && (binext.IsBinary(path) || executableFile) {
		return false
	}
	return !me.exclude.MatchString(path)
}

var (
	//	red   = "\033[31m"
	green = "\033[32m"
	reset = "\033[0m"
)
