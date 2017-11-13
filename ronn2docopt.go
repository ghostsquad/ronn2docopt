package main

import (
	"fmt"
	"bufio"
	"os"
	"io"
	"strings"
	"github.com/docopt/docopt-go"
)

// read ronn file
// output string for use with docopt
// strip markdown specific syntax

func captureSection(reader *bufio.Reader, c func(x string) (string, bool)) ([]string, error) {
	var section []string

	var line string
	var err error

	for {
		line, err = reader.ReadString('\n')
		if strings.HasPrefix(strings.TrimSpace(line), "##") {
			break
		}

		captureLine, keep := c(line)
		if keep {
			section = append(section, captureLine)
		}

		if err != nil {
			break
		}
	}

	return section, err
}

// returns a string as well as an boolean indicator if this value should be added to capture section
func captureSynopsis(line string) (string, bool) {
	synopsisLine := strings.TrimSpace(line)
	if synopsisLine == "" {
		return synopsisLine, false
	}

	synopsisLine = strings.Replace(synopsisLine, "`", "", -1)
	synopsisLine = strings.TrimRight(synopsisLine, "<br>")
	return "  " + synopsisLine, true
}

func captureOption(line string) (string, bool) {
	optionsLine := line

	if strings.HasPrefix(optionsLine, "  * ") && strings.Contains(optionsLine, ":") {
		optionsLine = strings.TrimSpace(optionsLine)
		if optionsLine == "" {
			return optionsLine, false
		}

		optionsLine = strings.TrimLeft(optionsLine, "*")
		optionsLine = strings.TrimRight(optionsLine, ":")
		optionsLine = strings.Replace(optionsLine, "`", "", -1)
		optionsLine = strings.TrimSpace(optionsLine)
		return "  " + optionsLine, true
	}

	return optionsLine, false
}

func ronn2docopt(ronnFile string) (string, error) {
	file, err := os.Open(ronnFile)
	defer file.Close()

	if err != nil {
		return "", err
	}

	reader := bufio.NewReader(file)

	var line string
	var synopsis []string
	var options []string

	for {
		line, err = reader.ReadString('\n')
		line = strings.TrimSpace(line)

		if line == "## SYNOPSIS" {
			synopsis, err = captureSection(reader, captureSynopsis)
		}

		if line == "## OPTIONS" {
			options, err = captureSection(reader, captureOption)
		}

		if err != nil {
			break
		}
	}

	docoptUsage := fmt.Sprintf(`
Usage:
%s

Options:
%s`, strings.Join(synopsis, "\n"), strings.Join(options, "\n"))

	if err != io.EOF {
		return "", err
	}

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