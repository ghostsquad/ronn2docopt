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

	scanner := bufio.NewScanner(file)

	content, err := libStrings.ReadLines(scanner)
	if err != nil {
		return "", err
	}

	d := ronn.RonnToDocopt(content)

	return strings.TrimSpace(d.String()), nil
}


func main() {
	usage := `Ronn2Docopt

Usage:
  ronn2docopt RONNFILE

Options:
  -h --help     Show this screen.
  --version     Show version.`

	arguments, _ := docopt.Parse(usage, nil, true, "ronn2docopt 0.1", false)

	docoptResults, err := ronn2docopt(arguments["RONNFILE"].(string))
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(docoptResults)
}