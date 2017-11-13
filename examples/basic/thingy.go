package main

import (
	"fmt"
	"github.com/docopt/docopt-go"
	"github.com/ghostsquad/ronn2docopt/examples/basic/lib"
)

func main() {
	asset, err := lib.Asset("docs/docopt.txt")
	usage := fmt.Sprintf("%s", asset)
	if err != nil {
		panic("asset not found")
	}

	arguments, _ := docopt.Parse(usage, nil, true, "Naval Fate 2.0", false)
	fmt.Println(arguments)
}