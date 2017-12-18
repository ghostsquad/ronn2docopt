package main

import (
	"fmt"
	"bufio"
	"os"
	"strings"
	"github.com/docopt/docopt-go"
	"github.com/ghostsquad/ronn2docopt/lib/ronn"
	libStrings "github.com/ghostsquad/ronn2docopt/lib/strings"
)

// read ronn file
// output string for use with docopt
// strip markdown specific syntax

func ronn2docopt(ronnFile string) (string, error) {
	file, err := os.Open(ronnFile)
	defer file.Close()

	if err != nil {
		return "", err
	}

	reader := bufio.NewReader(file)

	var docoptUsage string

	content, err := libStrings.ReadLines(reader)
	if err != nil {
		return "", err
	}

	ronn.RonnToDocopt(content)

	return strings.TrimSpace(docoptUsage), nil
}


func main() {
	usage := `Naval Fate.

Usage:
  ronn2docopt RONNFILE

Options:
  -h --help     Show this screen.
  --version     Show version.`

	arguments, _ := docopt.Parse(usage, nil, true, "ronn2docopt 0.1", false)

	usage, err := ronn2docopt(arguments["RONNFILE"].(string))
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(usage)
}