package ronn2docopt

import (
	"regexp"
	"strings"
	"bufio"
)

func RegexAndMatchNames(pattern string) (*regexp.Regexp, []string) {
	var re = regexp.MustCompile(pattern)
	var ma = re.SubexpNames()

	return re, ma
}

func NamedMatches(re *regexp.Regexp, matchNames []string, str string) map[string]string {
	namedMatches := map[string]string{}

	r2 := re.FindAllStringSubmatch(str, -1)

	if len(r2) > 0 {
		r3 := r2[0]

		for i, n := range r3 {
			namedMatches[matchNames[i]] = n
		}
	}

	return namedMatches
}

func PadRight(str string, pad string, length int) string {
	strLen := len(str)

	if strLen >= length {
		return str
	}

	padding := strings.Repeat(pad, length - strLen)

	return str + padding
}

func ReadLines(scanner *bufio.Scanner) ([]string, error) {
	var lines []string

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, nil
}
