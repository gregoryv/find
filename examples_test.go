package find_test

import (
	"fmt"
	"github.com/gregoryv/find"
)

func ExampleByName() {
	result, _ := find.ByName("*.md", ".")
	fmt.Printf("%v", result)
	//output:[README.md]
}
