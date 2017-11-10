[find](https://godoc.org/github.com/gregoryv/find) - Go package for locating files

[![Build Status](https://travis-ci.org/gregoryv/find.svg?branch=master)](https://travis-ci.org/gregoryv/find)

## Usage

    import "github.com/gregoryv/find"

	func main() {
	    result, _ := find.File("ex*.go", ".")
		//...
	}
