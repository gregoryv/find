package main

import (
	"fmt"
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

	s.WriteResult(os.Stdout)
	// output:
	// ./testdata/trucks.txt
	// 1 volvo 2010 new silver
	// 2 volvo 2011 new silver
	//
	// ./testdata/cars.txt
	// 3 volvo 1999 old green
}
