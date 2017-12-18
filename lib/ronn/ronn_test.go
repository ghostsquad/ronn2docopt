package ronn

import (
	"fmt"
	"github.com/pmezard/go-difflib/difflib"
	"testing"
	"strings"
)

var optionSection = []string{
	"## OPTIONS",
	"These options control whether output is written to file(s), standard output, or",
	"directly to a man pager.",
	"",
	"  * `-h`, `--help`:",
	"    Show this screen.",
	"  * `--version`:",
	"    Show version.",
	"  * `--speed=<kn>`:",
	"    Speed in knots. [default: 10]",
	"    The server respects the `--style` and document attribute options",
	"    (`--manual`, `--date`, etc.). These same options can be varied at request",
	"",
	"    *NOTE: This is a note",
	"",
	"  * `--foo`:",
	"    Multiline description",
	"    Only captures first line",
	"",
	"Other Options",
	"These options are for other use #cases",
	"",
	"  * `-b`:",
	"    Thingy",
	"    [default: baz]",
	"",
}

var exampleFile = append(append([]string{
	"ronn(1) -- convert markdown files to manpages",
	"=============================================",
	"",
	"## SYNOPSIS",
	"",
	"`naval_fate` `ship new <name>...`<br>",
	"`naval_fate` `ship <name> move <x> <y> [--speed=<kn>]`<br>",
	"`naval_fate` `ship shoot <x> <y>`<br>",
	"`naval_fate` `mine (set|remove) <x> <y> [--moored|--drifting]`<br>",
	"`naval_fate` `-h | --help`<br>",
	"`naval_fate` `--version`<br>",
	"",
	"## DESCRIPTION",
	"",
	"**Ronn** converts textfiles to standard roff-formatted UNIX manpages or HTML.",
	"",
}, optionSection...), []string{"## Another Section"}...)

func ExampleIsSectionHeader() {
	for _, line := range exampleFile {
		matched, section := isSectionHeader(line)
		if matched {
			fmt.Println(section)
		}
	}

	// Output:
	// SYNOPSIS
	// DESCRIPTION
	// OPTIONS
	// Another Section
}

func ExampleIsOptionDeclaration() {
	for _, line := range exampleFile {
		matched, name := isOptionDeclaration(line)
		if matched {
			fmt.Println(name)
		}
	}

	// Output:
	// -h --help
	// --version
	// --speed=<kn>
	// --foo
	// -b
}

func TestIsSectionDescriptionLine(t *testing.T) {
	got := ""

	for _, line := range append(optionSection, "## Something Else") {
		matched := isSectionDescriptionLine(line)
		if matched {
			got += line + "\n"
		}
	}

	got = strings.TrimRight(got, "\n")

	want :=
		"These options control whether output is written to file(s), standard output, or\n" +
		"directly to a man pager.\n" +
		"Other Options\n" +
		"These options are for other use #cases"

	if got != want {
		diff := difflib.UnifiedDiff{
			A:       difflib.SplitLines(want),
			B:       difflib.SplitLines(got),
			Context: 1,
		}
		text, _ := difflib.GetUnifiedDiffString(diff)

		fmt.Println(text)
		t.Error()
	}
}

func TestSanitizeSynopsis(t *testing.T) {
	got := formatSynopsis(exampleFile[4:12])

	want :=
		"  naval_fate ship new <name>...\n" +
		"  naval_fate ship <name> move <x> <y> [--speed=<kn>]\n" +
		"  naval_fate ship shoot <x> <y>\n" +
		"  naval_fate mine (set|remove) <x> <y> [--moored|--drifting]\n" +
		"  naval_fate -h | --help\n" +
		"  naval_fate --version"

	if got != want {
		diff := difflib.UnifiedDiff{
			A:       difflib.SplitLines(want),
			B:       difflib.SplitLines(got),
			Context: 1,
		}
		text, _ := difflib.GetUnifiedDiffString(diff)

		fmt.Println(text)
		t.Error()
	}
}

func TestUpdateWithLine(t *testing.T) {
	t.Run("when has trailing period, no default", func(t *testing.T){
		o := HelpOption{}
		o.updateWithLine("    Show this help screen.")

		got := o.Desc

		want := "Show this help screen."

		if got != want {
			t.Errorf("option.Desc got = %s, want %s", got, want)
		}

		got = o.DefaultValue
		if got != "" {
			t.Errorf("option.DefaultValue got = %s, want <empty>", got)
		}
	})

	t.Run("when help option uses multiple sentences, no default", func(t *testing.T){
		o := HelpOption{}
		o.updateWithLine("    Show this help screen. Use it to get help.")

		got := o.Desc

		want := "Show this help screen."

		if got != want {
			t.Errorf("option.Desc got = %s, want %s", got, want)
		}

		got = o.DefaultValue
		if got != "" {
			t.Errorf("option.DefaultValue got = %s, want <empty>", got)
		}
	})

	t.Run("when desc has no trailing period, no default", func(t *testing.T){
		o := HelpOption{}
		o.updateWithLine("    Show this help screen")

		got := o.Desc
		want := "Show this help screen"

		if got != want {
			t.Errorf("option.Desc got = %s, want %s", got, want)
		}

		got = o.DefaultValue
		if got != "" {
			t.Errorf("option.DefaultValue got = %s, want <empty>", got)
		}
	})

	t.Run("when desc has no trailing period, includes default", func(t *testing.T){
		o := HelpOption{}
		o.updateWithLine("    Show this help screen [default: foo]")

		got := o.Desc
		want := "Show this help screen"
		if got != want {
			t.Errorf("option.Desc got = %s, want %s", got, want)
		}

		got = o.DefaultValue
		want = "[default: foo]"
		if got != want {
			t.Errorf("option.DefaultValue got = %s, want %s", got, want)
		}
	})

	t.Run("when desc has trailing period, includes default", func(t *testing.T){
		o := HelpOption{}
		o.updateWithLine("    Speed in knots. [default: 10]")

		got := o.Desc
		want := "Speed in knots."
		if got != want {
			t.Errorf("option.Desc got = %s, want %s", got, want)
		}

		got = o.DefaultValue
		want = "[default: 10]"
		if got != want {
			t.Errorf("option.DefaultValue got = %s, want %s", got, want)
		}
	})

	t.Run("when desc has multiple sentences, includes default", func(t *testing.T){
		o := HelpOption{}
		o.updateWithLine("    Show this help screen. Use it to get help. [default: foo]")

		got := o.Desc
		want := "Show this help screen."
		if got != want {
			t.Errorf("option.Desc got = %s, want %s", got, want)
		}

		got = o.DefaultValue
		want = "[default: foo]"
		if got != want {
			t.Errorf("option.DefaultValue got = %s, want %s", got, want)
		}
	})

	t.Run("when desc has multiple lines, includes default", func(t *testing.T){
		o := HelpOption{}
		o.updateWithLine("    Show this help screen.")
		o.updateWithLine("    Use it to get help. [default: foo]")

		got := o.Desc
		want := "Show this help screen."
		if got != want {
			t.Errorf("option.Desc got = %s, want %s", got, want)
		}

		got = o.DefaultValue
		want = "[default: foo]"
		if got != want {
			t.Errorf("option.DefaultValue got = %s, want %s", got, want)
		}
	})
}

func TestRonnToDocopt(t *testing.T) {
	t.Run("example file has 2 option sub-sections", func(t *testing.T) {
		d := RonnToDocopt(exampleFile)

		got := len(d.HelpOptionSections)
		want := 2
		if got != want {
			t.Errorf("number of option subsections got = %d, want %d", got, want)
		}
	})

	t.Run("first option sub-section has 4 options", func(t *testing.T) {
		d := RonnToDocopt(exampleFile)

		got := len(d.HelpOptionSections)
		want := 2
		if got != want {
			t.Errorf("number of help option sections got = %d, want %d", got, want)
			return
		}

		got = len(d.HelpOptionSections[0].Options)
		want = 4
		if got != want {
			t.Errorf("number of options got = %d, want %d", got, want)
		}
	})

	t.Run("section option sub-section has 1 option", func(t *testing.T) {
		d := RonnToDocopt(exampleFile)

		got := len(d.HelpOptionSections)
		want := 2
		if got != want {
			t.Errorf("number of help option sections got = %d, want %d", got, want)
			return
		}

		got = len(d.HelpOptionSections[1].Options)
		want = 1
		if got != want {
			t.Errorf("number of options got = %d, want %d", got, want)
		}
	})
	
	t.Run("options section", func(t *testing.T) {
		t.Run("when starts with an option", func(t *testing.T) {
			t.Skip()
		})

		t.Run("when starts with section description", func(t *testing.T) {
			t.Skip()
		})

		t.Run("when has multiple sub-sections", func(t *testing.T) {
			t.Skip()
		})

		t.Run("when has multiple sections each with multiple options", func(t *testing.T) {
			t.Skip()
		})

		t.Run("when sub-section description is multiple paragraphs", func(t *testing.T) {
			t.Skip()
		})

		t.Run("when option description is multiple paragraphs", func(t *testing.T) {
			t.Skip()
		})
	})
}

func TestDocOpt_String(t *testing.T) {
	t.Run("end 2 end", func(t *testing.T) {
		d := RonnToDocopt(exampleFile)
		got := d.String()

		want := "Usage:\n" +
				"  naval_fate ship new <name>...\n" +
				"  naval_fate ship <name> move <x> <y> [--speed=<kn>]\n" +
				"  naval_fate ship shoot <x> <y>\n" +
				"  naval_fate mine (set|remove) <x> <y> [--moored|--drifting]\n" +
				"  naval_fate -h | --help\n" +
				"  naval_fate --version\n" +
				"\n" +
				"Options:\n" +
				"These options control whether output is written to file(s), standard output, or\n" +
				"directly to a man pager.\n" +
				"\n" +
				"  -h --help     Show this screen.\n" +
				"  --version     Show version.\n" +
				"  --speed=<kn>  Speed in knots. [default: 10]\n" +
				"  --foo         Multiline description\n" +
				"\n" +
				"Other Options\n" +
				"These options are for other use #cases\n" +
				"\n" +
				"  -b  Thingy [default: baz]"

		if got != want {
			diff := difflib.UnifiedDiff{
				A:       difflib.SplitLines(want),
				B:       difflib.SplitLines(got),
				Context: 3,
			}
			text, _ := difflib.GetUnifiedDiffString(diff)

			fmt.Println(text)
			t.Error()
		}
	})

	t.Run("when starts with an option", func(t *testing.T) {
		t.Skip()
	})

	t.Run("when starts with section description", func(t *testing.T) {
		t.Skip()
	})

	t.Run("when has multiple sections each with multiple options", func(t *testing.T) {
		t.Skip()
	})

}
