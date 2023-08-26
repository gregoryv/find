package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/gregoryv/binext"
	"github.com/gregoryv/find"
)

func main() {

	flag := flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	var colors bool
	flag.BoolVar(&colors, "c", false, "")
	flag.BoolVar(&colors, "color", false, "")

	var includeBinary bool
	flag.BoolVar(&includeBinary, "i", false, "")
	flag.BoolVar(&includeBinary, "include-binary", false, "")

	var writeAliases string
	usage := "Output file for search result aliases for shell sourcing"
	flag.StringVar(&writeAliases, "w", "", usage)
	flag.StringVar(&writeAliases, "write-aliases", "", usage)

	var aliasPrefix string
	usage = `Use together with -w to prefix numbered aliases
	e.g -w -a t results in alias t1=...`
	flag.StringVar(&aliasPrefix, "a", "", usage)
	flag.StringVar(&aliasPrefix, "alias-prefix", "", usage)

	var exclude string = "^.git/|(pdf|svg)$"
	usage = "Regexp for excluding paths"
	flag.StringVar(&exclude, "e", exclude, usage)
	flag.StringVar(&exclude, "exclude", exclude, usage)

	var verbose bool
	flag.BoolVar(&verbose, "verbose", false, "")

	flag.Parse(os.Args[1:])
	rest := flag.Args()

	log.SetFlags(0)
	if len(rest) == 0 {
		log.Fatal("missing expression")
	}

	expr := rest[0]
	rest = rest[1:]

	var index int
	if len(rest) > 0 {
		last := len(rest) - 1
		var err error
		index, err = strconv.Atoi(rest[last])
		if err == nil {
			// last is index
			rest = rest[:last]
		}
	}

	// ----------------------------------------

	if expr == "" {
		fmt.Println("empty EXPR")
		os.Exit(1)
	}

	s := NewScanner()
	if verbose {
		s.Logger.SetOutput(log.Writer())
	}
	filter := &smart{}
	filter.SetIncludeBinary(includeBinary)
	if err := filter.SetExclude(exclude); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if len(rest) > 0 {
		s.SetFiles(rest)
	} else {
		s.SetFiles(
			ls("", filter),
		)
	}

	if err := s.Scan(expr); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if writeAliases != "" {
		var i int
		aw, err := os.Create(writeAliases)
		if err != nil {
			log.Fatal(err)
		}
		for _, fm := range s.LastResult() {
			for _, lm := range fm.Result {
				fmt.Fprintln(aw, aliasLine(i+1, aliasPrefix, fm, lm))
				i++
			}
		}
		aw.Close()
	}

	// results destination
	w := os.Stdout

	if index == 0 { // list result
		var i int
		for _, fm := range s.LastResult() {
			fmt.Fprintln(w, fm.Filename)
			for _, m := range fm.Result {
				text := m.Text
				if colors {
					colored := fmt.Sprintf("%s%s%s", green, expr, reset)
					text = strings.ReplaceAll(text, expr, colored)
				}
				fmt.Fprintln(w, i+1, text)
				i++
			}
			fmt.Fprintln(w)
		}
		return
	}

	var i int
	for _, fm := range s.LastResult() {
		for _, lm := range fm.Result {
			i++
			if i != index {
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
	fmt.Printf("%v? there are only %v matches\n", index, i)
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
	fmt.Println("ls", files)
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
