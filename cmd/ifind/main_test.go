package main

import (
	"os"
	"testing"
)

func Test_main(t *testing.T) {
	os.Args = []string{"test", "Println"}
	main()
}
