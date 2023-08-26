package main

import (
	"flag"
	"os"
	"strings"
	"testing"
)

func Test_main(t *testing.T) {
	os.Args = []string{"test", "Println"}
	main()

	// verify usage contains some reference to all the options
	f.VisitAll(func(f *flag.Flag) {
		opt := "-" + f.Name
		if len(f.Name) > 1 {
			opt = "--" + f.Name
		}
		if !strings.Contains(usage, opt) {
			t.Error("missing", opt, "in usage")
		}
	})
}
