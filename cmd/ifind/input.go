package main

import (
	"errors"
	"io"
	"os"
	"strconv"
	"text/template"

	"github.com/gregoryv/cli"
)

func NewInput() *Input {
	in := Input{
		Exclude: "^.git/|(pdf|svg)$",
	}
	if v := os.Getenv("IFIND_EXCLUDE_REGEXP"); v != "" {
		in.Exclude = v
	}
	return &in
}

type Input struct {
	Help          bool
	Colors        bool
	Verbose       bool
	IncludeBinary bool

	Expression string
	OpenIndex  uint32

	Glob    string // glob expression
	Exclude string
	Files []string

	WriteAliases string
	AliasPrefix  string
}

func (in *Input) SetArg(option, value string) (err error) {
	switch option {
	case "-h", "--help":
		err = parseBool(&in.Help, value)

	case "-c", "--colors":
		err = parseBool(&in.Colors, value)

	case "-i", "--include-binary":
		err = parseBool(&in.IncludeBinary, value)

	case "-v", "--verbose":
		err = parseBool(&in.Verbose, value)

	case "-f", "--files":
		in.Glob = value

	case "-w", "--write-alias":
		in.WriteAliases = value

	case "-a", "--alias-prefix":
		in.AliasPrefix = value

	case "-e", "--exclude":
		in.Exclude = value

	case "":
		if in.Expression == "" {
			in.Expression = value
			return nil
		}
		// if it's not a number assume it's a file.  Limitation files
		// named as numbers only cannot be specified like this. Use -f
		// for that.
		if err := parseUint32(&in.OpenIndex, value); err != nil {
			in.Files = append(in.Files, value)
		}

	default:
		return cli.ErrOption
	}
	return
}

const usage = `Usage: {{.Cmd}} [OPTIONS] EXPR [FILE...] [OPEN_INDEX]

ifind - grep expression and quick open indexed result

Options
    -f, --files : ""
        Empty means current working directory and recursive.
        The pattern is a glob format like *.go or *.txt
        Alternatively specify files after the expression.

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
    Look for EXPR in all text files
        $ ifind -f *.txt EXPR
    or
        $ ifind EXPR *.txt

    Open the third match
        $ EDITOR=emacsclient ifind EXPR 3

`

func WriteUsage(w io.Writer) {
	usageTmpl.Execute(w, map[string]any{
		"Cmd": os.Args[0],
	})
}

var usageTmpl = template.Must(template.New("").Parse(usage))

// empty value means true, otherwise strconv.ParseBool is used
func parseBool(dst *bool, value string) error {
	if value == "" {
		*dst = true
		return nil
	}
	tmp, err := strconv.ParseBool(value)
	if err != nil {
		return errors.Unwrap(err)
	}
	*dst = tmp
	return nil
}

func parseUint32(dst *uint32, value string) error {
	tmp, err := strconv.ParseUint(value, 10, 32)
	if err != nil {
		return errors.Unwrap(err)
	}
	*dst = uint32(tmp)
	return nil
}
