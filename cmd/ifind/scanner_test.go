package main

import (
	"fmt"
	"io"
	"os"
)

func Example() {
	s := NewScanner()
	s.SetFiles([]string{
		"./testdata/trucks.txt",
		"./testdata/cars.txt",
	})

	if err := s.Scan("volvo"); err != nil {
		fmt.Println(err)
	}

	WriteResult(os.Stdout, s.LastResult())
	// output:
	// ./testdata/trucks.txt
	// 1 volvo 2010 new silver
	// 2 volvo 2011 new silver
	//
	// ./testdata/cars.txt
	// 3 volvo 1999 old green
}

func WriteResult(w io.Writer, result []FileMatch) {
	var i int
	for _, fm := range result {
		fmt.Fprintln(w, fm.Filename)
		for _, m := range fm.Result {
			fmt.Fprintln(w, i+1, m.Text)
			i++
		}
		fmt.Fprintln(w)
	}
}
