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
	"text/template"

	"github.com/gregoryv/binext"
	"github.com/gregoryv/find"
)

var usageTmpl = template.Must(template.New("").Parse(usage))

const usage = `Usage: {{.Cmd}} [OPTIONS] EXPR [FILE...] [OPEN_INDEX]

ifind - grep expression and quick open indexed result

Options
    -c, --colors
    -i, --include-binary
    -w, --write-aliases : ""
        Output file for search result aliases for shell sourcing

    -a, --alias-prefix : ""
        Use together with -w to prefix numbered aliases
        e.g -w -a t results in alias t1=...

    -e, --exclude, $IFIND_EXCLUDE_REGEXP : "^.git/|(pdf|svg)$"
        Regexp for excluding paths

    -v, --verbose
    -h, --help

Examples
    Look for EXPR in files recursively
        $ ifind EXPR

    or in specific files
        $ ifind EXPR *.txt

    Open the third match
        $ EDITOR=emacsclient ifind EXPR 3

`

// keep it outside so we can verify usage content to the actual flags
var f = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

func main() {

	var colors bool
	f.BoolVar(&colors, "color", false, "")
	f.BoolVar(&colors, "c", false, "")

	var includeBinary bool
	f.BoolVar(&includeBinary, "include-binary", false, "")
	f.BoolVar(&includeBinary, "i", false, "")

	var writeAliases string
	f.StringVar(&writeAliases, "w", "", "")
	f.StringVar(&writeAliases, "write-aliases", "", "")

	var aliasPrefix string
	f.StringVar(&aliasPrefix, "a", "", "")
	f.StringVar(&aliasPrefix, "alias-prefix", "", "")

	var exclude string
	var excludeDef = "^.git/|(pdf|svg)$"
	f.StringVar(&exclude, "e", excludeDef, "")
	f.StringVar(&exclude, "exclude", excludeDef, "")
	envs.StringVar(&exclude, excludeDef, "IFIND_EXCLUDE_REGEXP")

	var verbose bool
	f.BoolVar(&verbose, "verbose", false, "")

	f.Usage = func() {
		usageTmpl.Execute(os.Stdout, map[string]any{
			"Cmd": os.Args[0],
		})
	}

	// parse arguments
	f.Parse(os.Args[1:])

	log.SetFlags(0)

	// find expression
	rest := f.Args()
	if len(rest) == 0 {
		log.Fatal("missing expression")
	}
	expr := rest[0]
	rest = rest[1:]

	// find optional index as final argument
	var index int
	if len(rest) > 0 {
		last := len(rest) - 1
		var err error
		index, err = strconv.Atoi(rest[last])
		if err == nil {
			// remaining arguments should be a list of files
			rest = rest[:last]
		}
	}

	// setup scanner
	s := NewScanner()
	if verbose {
		s.Logger.SetOutput(log.Writer())
	}
	if len(rest) > 0 {
		s.SetFiles(rest)
	} else {
		filter := &smart{}
		filter.SetIncludeBinary(includeBinary)
		if err := filter.SetExclude(exclude); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		s.SetFiles(ls(filter))
	}

	// scan for expression
	if err := s.Scan(expr); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// write out optional aliases file
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

	// list result
	w := os.Stdout
	if index == 0 {
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

	// open selected indexed result
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
func ls(filter find.Matcher) []string {
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

// ----------------------------------------

// smart a configurable filter when searching for files
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
	green = "\033[32m"
	reset = "\033[0m"
)
