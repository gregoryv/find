package find_test

import (
	"fmt"
	"github.com/gregoryv/find"
)

func ExampleByName() {
	result, _ := find.ByName("*.go", ".")
	fmt.Printf("%v", result)
	//output:[examples_test.go search.go search_test.go shell.go]
}
