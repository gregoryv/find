package find

import (
	"bufio"
	"io"
	"os"
)

// Ref points to a line in a file starting from 0
type Ref struct {
	LineNo int
	Line   string
}

type Refs []Ref

// Grep finds each line that matches the given pattern.
func Grep(pattern, file string) (res Refs, err error) {
	var stream *os.File
	stream, err = os.Open(file)
	if err != nil {
		return
	}
	defer stream.Close()
	return grep(pattern, stream), nil
}

func grep(pattern string, stream io.Reader) Refs {
	res := make([]Ref, 0)
	scanner := bufio.NewScanner(stream)
	sp := NewShellPattern("*" + pattern + "*")
	var line string
	var no int
	for scanner.Scan() {
		line = scanner.Text()
		if sp.Match(line) {
			res = append(res, Ref{no, line})
		}
		no++
	}

	return res
}
