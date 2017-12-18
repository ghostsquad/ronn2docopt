package strings

import (
	"testing"
)

func TestNamedMatches(t *testing.T) {
	str := "Hello World"

	t.Run("when string matches", func(t *testing.T) {
		re, ma := RegexAndMatchNames("(?P<first_char>.)")
		na := NamedMatches(re, ma, str)

		want := "H"
		got := na["first_char"]

		if got != want {
			t.Errorf("got = %s, want %s", got, want)
		}
	})

	t.Run("when string does not match", func(t *testing.T) {
		re, ma := RegexAndMatchNames("(?P<foo>foo)")
		na := NamedMatches(re, ma, str)

		want := 0
		got := len(na)

		if got != want {
			t.Errorf("got = %s, want %s", got, want)
		}
	})

	t.Run("with multiple named matches", func(t *testing.T) {
		re, ma := RegexAndMatchNames("^(?P<first>.).*(?P<last>.)$")
		na := NamedMatches(re, ma, str)

		want := "H"
		got := na["first"]

		if got != want {
			t.Errorf("got = %s, want %s", got, want)
		}

		want = "d"
		got = na["last"]

		if got != want {
			t.Errorf("got = %s, want %s", got, want)
		}
	})
}

func TestPadRight(t *testing.T) {
	str := "foo"

	t.Run("when string is smaller than padding", func(t *testing.T) {
		got := PadRight(str, "*", 5)

		want := "foo**"

		if got != want {
			t.Errorf("got = %s, want %s", got, want)
		}
	})

	t.Run("when string is same length as padding", func(t *testing.T) {
		got := PadRight(str, "*", 3)

		want := str

		if got != want {
			t.Errorf("got = %s, want %s", got, want)
		}
	})

	t.Run("when string is longer than padding", func(t *testing.T) {
		got := PadRight(str, "*", 1)

		want := str

		if got != want {
			t.Errorf("got = %s, want %s", got, want)
		}
	})

	t.Run("when string is empty", func(t *testing.T) {
		got := PadRight("", "*", 3)

		want := "***"

		if got != want {
			t.Errorf("got = %s, want %s", got, want)
		}
	})
}
