package ronn2docopt

import (
	"strings"
	"bytes"
	"regexp"
	"os"
	"bufio"
)

type DocOpt struct {
	Synopsis string
	HelpOptionSections []HelpOptionSection
}

type Synopsis struct {
	Body string
}

type HelpOptionSection struct {
	Name string
	Options []HelpOption
}

type HelpOption struct {
	Name         string
	Desc         string
	DefaultValue string
}

var brRe = regexp.MustCompile(`^(.*)\s*(<br>)\s*$`)

var sectionHeaderRe, sectionHeaderMa = RegexAndMatchNames(`^##\s+(?P<section>.*)$`)
var namedOptionRe, namedOptionMa = RegexAndMatchNames(`^ {2}\* ` + "`?" + `(?P<name>-.*):$`)
var defaultValueRe, defaultValueMa = RegexAndMatchNames(`^ {4}(?P<before>.*)(?P<default>\[default: .*\]$)`)
var shortOptionDescRe, shortOptionDescMa = RegexAndMatchNames(`^ {4}(?P<short>.*?[.!?]).*$`)

/*

1. In Options Section?
	> Yes - Continue to Step 2..
	> No - Keep looking for options section
2. At Option Declaration?
	> Yes - Begin collecting option lines until
		> Reached New Section. STOP! No other Sections to parse.
        > Reached New Option Declaration. Go to 2.
        > Reached New Option Section Description. Go to 3.
	> No, continue to step 3
3. At Option Section Description?
	> Yes - Begin collecting description lines until
		> New Section Reached? STOP! No other Sections to parse.
		> New Option Declaration Found? Go To 2.
	> No - Go to 1

 */

func (d *DocOpt) String() string {
	var buffer bytes.Buffer

	buffer.WriteString("Usage:\n")

	buffer.WriteString(d.Synopsis)
	buffer.WriteString("\n\n")

	buffer.WriteString("Options:\n")

	for i, s := range d.HelpOptionSections {
		if i > 0 && s.Name != "" {
			buffer.WriteString(s.Name)
			buffer.WriteString("\n")
		}

		longestOptionNameLen := s.longestOptionNameLen()

		for _, o := range s.Options {
			padding := 0
			if len(o.Desc) > 0 || len(o.DefaultValue) > 0 {
				// 2 spaces + name + 2 spaces
				padding = longestOptionNameLen + 4
			}

			on := PadRight("  " + o.Name, " ", padding)
			buffer.WriteString(on)
			buffer.WriteString(o.Desc)

			if len(o.DefaultValue) > 0 {
				buffer.WriteString(" ")
				buffer.WriteString(o.DefaultValue)
			}

			buffer.WriteString("\n")
		}

		buffer.WriteString("\n")
	}

	results := buffer.String()
	results = strings.TrimRight(results, "\n")

	return results
}

func  RonnToDocopt(lines []string) *DocOpt {
	var d DocOpt

	s := getSection(lines, "SYNOPSIS")
	d.Synopsis = formatSynopsis(s)

	o := getSection(lines, "OPTIONS")

	lastWasOptions := false
	var sectionLines []string

	for _, line := range o {
		if lastWasOptions && isSectionDescriptionLine(line) {
			s := newOptionSection(sectionLines)
			d.HelpOptionSections = append(d.HelpOptionSections, *s)
			sectionLines = nil
			lastWasOptions = false
		} else if f, _ := isOptionDeclaration(line); f {
			lastWasOptions = true
		}

		sectionLines = append(sectionLines, line)
	}

	// append final section
	if len(sectionLines) > 0 {
		s := newOptionSection(sectionLines)
		d.HelpOptionSections = append(d.HelpOptionSections, *s)
	}

	return &d
}

func ConvertRonnFile(ronnFile string) (string, error) {
	file, err := os.Open(ronnFile)
	defer file.Close()

	if err != nil {
		return "", err
	}

	scanner := bufio.NewScanner(file)

	content, err := ReadLines(scanner)
	if err != nil {
		return "", err
	}

	d := RonnToDocopt(content)

	return strings.TrimSpace(d.String()), nil
}

// ==================================================== //
// PRIVATE METHODS
// ---------------------------------------------------- //

func (s *HelpOptionSection) longestOptionNameLen() int {
	l := 0
	for _, o := range s.Options {
		nl := len(o.Name)
		if nl > l {
			l = nl
		}
	}

	return l
}

// A section is delimited by section headers
func getSection(lines []string, sectionName string) []string {
	var section []string

	sectionFound := false

	for _, line := range lines {
		r, s := isSectionHeader(line)

		// if we've reached a new section,
		// but we were already in the desired section,
		// quit reading
		if r && sectionFound == true {
			break
		}

		// if we've reached a new nection
		// and it's the desired section, skip this line
		// then indicate we should start recording the lines
		if r && s == sectionName {
			sectionFound = true
			continue
		}

		if sectionFound {
			section = append(section, line)
		}
	}

	return section
}

func newOption(name string, lines []string) *HelpOption {
	h := &HelpOption{
		Name: name,
	}

	for _, line := range lines {
		h.updateWithLine(line)
	}

	return h
}

func newOptionSection(lines []string) *HelpOptionSection {
	h := &HelpOptionSection{}

	var sectionDescriptionLines []string
	var prevOptionName string
	var optionLines []string
	inSectionDesc := true

	for _, line := range lines {
		if f, newOptionName := isOptionDeclaration(line); f {

			// the first option found basically kicks off the line "collection"
			// subsequent option declarations kick off the creation of a new option
			// from previously collected lines
			if len(optionLines) > 0 {
				o := newOption(prevOptionName, optionLines)
				h.Options = append(h.Options, *o)
				optionLines = nil
			}

			prevOptionName = newOptionName
			inSectionDesc = false
			continue
		}

		if inSectionDesc {
			sectionDescriptionLines = append(sectionDescriptionLines, line)
			continue
		}

		optionLines = append(optionLines, line)
	}

	// append final option to list
	if len(prevOptionName) > 0 {
		o := newOption(prevOptionName, optionLines)
		h.Options = append(h.Options, *o)
	}

	// finalize section description
	if len(sectionDescriptionLines) > 0 {
		h.Name = strings.TrimSpace(sectionDescriptionLines[0])
	}

	return h
}

// A Section Header is a markdown H2 e.g.
// ## Foo
func isSectionHeader(line string) (bool, string) {
	ma := NamedMatches(sectionHeaderRe, sectionHeaderMa, line)
	if len(ma) > 0 {
		return true, ma["section"]
	}

	return false, ""
}

// A Section Description Line is a line that does not begin with any spaces (and is not a section header)
func isSectionDescriptionLine(line string) bool {
	if strings.HasPrefix(line, " ") {
		return false
	}

	if sH, _ := isSectionHeader(line); sH {
		return false
	}

	if strings.TrimSpace(line) == "" {
		return false
	}

	return true
}

// An Option Declaration looks like this:
//   * `--help`:
func isOptionDeclaration(line string) (bool, string) {
	ma := NamedMatches(namedOptionRe, namedOptionMa, line)

	if len(ma) > 0 {
		name := strings.Replace(ma["name"], "`", "", -1)
		name = strings.Replace(name, ",", "", -1)

		return true, name
	}

	return false, ""
}


// The Synopsis should be stripped of specific markdown/html syntax:
// * backticks (`)
// * html line break (<br>)
// It also should be formatted such that each line begins with 2 spaces.
func formatSynopsis(lines []string) string {
	var buffer bytes.Buffer
	for _, line := range lines {
		line = strings.Replace(line, "`", "", -1)
		line = brRe.ReplaceAllString(line, "$1")
		line = strings.TrimSpace(line)

		if line != "" {
			buffer.WriteString("  " + line)
		}

		buffer.WriteString("\n")
	}

	synopsis := buffer.String()
	synopsis = strings.TrimLeft(synopsis, "\n")
	synopsis = strings.TrimRight(synopsis, "\n")

	return synopsis
}

func (option *HelpOption) updateWithLine(line string) {
	if option.Desc == "" {
		if ma := NamedMatches(shortOptionDescRe, shortOptionDescMa, line); len(ma) > 0 {
			option.Desc = strings.TrimSpace(ma["short"])
		} else if ma := NamedMatches(defaultValueRe, defaultValueMa, line); len(ma) > 0 {
			option.Desc = strings.TrimSpace(ma["before"])
			option.DefaultValue = ma["default"]
			return
		} else {
			option.Desc = strings.TrimSpace(line)
		}
	}

	if option.DefaultValue == "" {
		ma := NamedMatches(defaultValueRe, defaultValueMa, line)
		if len(ma) > 0 {
			option.DefaultValue = ma["default"]
		}
	}
}
