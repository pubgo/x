package strutil_test

import (
	"github.com/pubgo/x/pkg/strutil"
	"testing"
)

func TestMatch(t *testing.T) {
	assert := func(p, s string) {
		if !strutil.Match(p, s, 0) {
			t.Errorf("Assertion failed: Match(%#v, %#v, 0)", p, s)
		}
	}
	assert("", "")
	assert("*", "")
	assert("*", "foo")
	assert("*", "bar")
	assert("*", "*")
	assert("**", "f")
	assert("**", "foo.txt")
	assert("*.*", "foo.txt")
	assert("foo*.txt", "foobar.txt")
	assert("foo.txt", "foo.txt")
	assert("foo\\.txt", "foo.txt")
	if strutil.Match("foo\\.txt", "foo.txt", strutil.FNM_NOESCAPE) {
		t.Errorf("Assertion failed: Match(%#v, %#v, FNM_NOESCAPE) == false", "foo\\.txt", "foo.txt")
	}
}


func TestWildcard(t *testing.T) {
	// A wildcard pattern "*" should match anything
	cases := []struct {
		pattern string
		input   string
		flags   int
		want    bool
	}{
		{"*", "", 0, true},
		{"*", "foo", 0, true},
		{"*", "*", 0, true},
		{"*", "   ", 0, true},
		{"*", ".foo", 0, true},
		{"*", "わたし", 0, true},
	}

	for tc, c := range cases {
		got := strutil.Match(c.pattern, c.input, c.flags)
		if got != c.want {
			t.Errorf(
				"Testcase #%d failed: strutil.Match('%s', '%s', %d) should be %v not %v",
				tc, c.pattern, c.input, c.flags, c.want, got,
			)
		}
	}
}

func TestWildcardSlash(t *testing.T) {
	cases := []struct {
		pattern string
		input   string
		flags   int
		want    bool
	}{
		// Should match / when flags are 0
		{"*", "foo/bar", 0, true},
		{"*", "/", 0, true},
		{"*", "/foo", 0, true},
		{"*", "foo/", 0, true},
		// Shouldnt match / when flags include FNM_PATHNAME
		{"*", "foo/bar", strutil.FNM_PATHNAME, false},
		{"*", "/", strutil.FNM_PATHNAME, false},
		{"*", "/foo", strutil.FNM_PATHNAME, false},
		{"*", "foo/", strutil.FNM_PATHNAME, false},
	}

	for tc, c := range cases {
		got := strutil.Match(c.pattern, c.input, c.flags)
		if got != c.want {
			t.Errorf(
				"Testcase #%d failed: strutil.Match('%s', '%s', %d) should be %v not %v",
				tc, c.pattern, c.input, c.flags, c.want, got,
			)
		}
	}
	for _, c := range cases {
		got := strutil.Match(c.pattern, c.input, c.flags)
		if got != c.want {
			t.Errorf(
				"strutil.Match('%s', '%s', %d) should be %v not %v",
				c.pattern, c.input, c.flags, c.want, got,
			)
		}
	}
}

func TestWildcardFNMPeriod(t *testing.T) {
	// FNM_PERIOD means that . is not matched in some circumstances.
	cases := []struct {
		pattern string
		input   string
		flags   int
		want    bool
	}{
		{"*", ".foo", strutil.FNM_PERIOD, false},
		{"/*", "/.foo", strutil.FNM_PERIOD, true},
		{"/*", "/.foo", strutil.FNM_PERIOD | strutil.FNM_PATHNAME, false},
	}

	for tc, c := range cases {
		got := strutil.Match(c.pattern, c.input, c.flags)
		if got != c.want {
			t.Errorf(
				"Testcase #%d failed: strutil.Match('%s', '%s', %d) should be %v not %v",
				tc, c.pattern, c.input, c.flags, c.want, got,
			)
		}
	}
}

func TestQuestionMark(t *testing.T) {
	//A question mark pattern "?" should match a single character
	cases := []struct {
		pattern string
		input   string
		flags   int
		want    bool
	}{
		{"?", "", 0, false},
		{"?", "f", 0, true},
		{"?", ".", 0, true},
		{"?", "?", 0, true},
		{"?", "foo", 0, false},
		{"?", "わ", 0, true},
		{"?", "わた", 0, false},
		// Even '/' when flags are 0
		{"?", "/", 0, true},
		// Except '/' when flags include FNM_PATHNAME
		{"?", "/", strutil.FNM_PATHNAME, false},
	}

	for tc, c := range cases {
		got := strutil.Match(c.pattern, c.input, c.flags)
		if got != c.want {
			t.Errorf(
				"Testcase #%d failed: strutil.Match('%s', '%s', %d) should be %v not %v",
				tc, c.pattern, c.input, c.flags, c.want, got,
			)
		}
	}
}

func TestQuestionMarkExceptions(t *testing.T) {
	//When flags include FNM_PERIOD a '?' might not match a '.' character.
	cases := []struct {
		pattern string
		input   string
		flags   int
		want    bool
	}{
		{"?", ".", strutil.FNM_PERIOD, false},
		{"foo?", "foo.", strutil.FNM_PERIOD, true},
		{"/?", "/.", strutil.FNM_PERIOD, true},
		{"/?", "/.", strutil.FNM_PERIOD | strutil.FNM_PATHNAME, false},
	}

	for tc, c := range cases {
		got := strutil.Match(c.pattern, c.input, c.flags)
		if got != c.want {
			t.Errorf(
				"Testcase #%d failed: strutil.Match('%s', '%s', %d) should be %v not %v",
				tc, c.pattern, c.input, c.flags, c.want, got,
			)
		}
	}
}

func TestRange(t *testing.T) {
	azPat := "[a-z]"
	cases := []struct {
		pattern string
		input   string
		flags   int
		want    bool
	}{
		// Should match a single character inside its range
		{azPat, "a", 0, true},
		{azPat, "q", 0, true},
		{azPat, "z", 0, true},
		{"[わ]", "わ", 0, true},

		// Should not match characters outside its range
		{azPat, "-", 0, false},
		{azPat, " ", 0, false},
		{azPat, "D", 0, false},
		{azPat, "é", 0, false},

		//Should only match one character
		{azPat, "ab", 0, false},
		{azPat, "", 0, false},

		// Should not consume more of the pattern than necessary
		{azPat + "foo", "afoo", 0, true},

		// Should match '-' if it is the first/last character or is
		// backslash escaped
		{"[-az]", "-", 0, true},
		{"[-az]", "a", 0, true},
		{"[-az]", "b", 0, false},
		{"[az-]", "-", 0, true},
		{"[a\\-z]", "-", 0, true},
		{"[a\\-z]", "b", 0, false},

		// ignore '\\' when FNM_NOESCAPE is given
		{"[a\\-z]", "\\", strutil.FNM_NOESCAPE, true},
		{"[a\\-z]", "-", strutil.FNM_NOESCAPE, false},

		// Should be negated if starting with ^ or !"
		{"[^a-z]", "a", 0, false},
		{"[!a-z]", "b", 0, false},
		{"[!a-z]", "é", 0, true},
		{"[!a-z]", "わ", 0, true},

		// Still match '-' if following the negation character
		{"[^-az]", "-", 0, false},
		{"[^-az]", "b", 0, true},

		// Should support multiple characters/ranges
		{"[abc]", "a", 0, true},
		{"[abc]", "c", 0, true},
		{"[abc]", "d", 0, false},
		{"[a-cg-z]", "c", 0, true},
		{"[a-cg-z]", "h", 0, true},
		{"[a-cg-z]", "d", 0, false},

		//Should not match '/' when flags is FNM_PATHNAME
		{"[abc/def]", "/", 0, true},
		{"[abc/def]", "/", strutil.FNM_PATHNAME, false},
		{"[.-0]", "/", 0, true}, // The range [.-0] includes /
		{"[.-0]", "/", strutil.FNM_PATHNAME, false},

		// Should normally be case-sensitive
		{"[a-z]", "A", 0, false},
		{"[A-Z]", "a", 0, false},
		//Except when FNM_CASEFOLD is given
		{"[a-z]", "A", strutil.FNM_CASEFOLD, true},
		{"[A-Z]", "a", strutil.FNM_CASEFOLD, true},
	}

	for tc, c := range cases {
		got := strutil.Match(c.pattern, c.input, c.flags)
		if got != c.want {
			t.Errorf(
				"Testcase #%d failed: strutil.Match('%s', '%s', %d) should be %v not %v",
				tc, c.pattern, c.input, c.flags, c.want, got,
			)
		}
	}
}

func TestBackSlash(t *testing.T) {
	cases := []struct {
		pattern string
		input   string
		flags   int
		want    bool
	}{
		//A backslash should escape the following characters
		{"\\\\", "\\", 0, true},
		{"\\*", "*", 0, true},
		{"\\*", "foo", 0, false},
		{"\\?", "?", 0, true},
		{"\\?", "f", 0, false},
		{"\\[a-z]", "[a-z]", 0, true},
		{"\\[a-z]", "a", 0, false},
		{"\\foo", "foo", 0, true},
		{"\\わ", "わ", 0, true},

		// Unless FNM_NOESCAPE is given
		{"\\\\", "\\", strutil.FNM_NOESCAPE, false},
		{"\\\\", "\\\\", strutil.FNM_NOESCAPE, true},
		{"\\*", "foo", strutil.FNM_NOESCAPE, false},
		{"\\*", "\\*", strutil.FNM_NOESCAPE, true},
	}

	for tc, c := range cases {
		got := strutil.Match(c.pattern, c.input, c.flags)
		if got != c.want {
			t.Errorf(
				"Testcase #%d failed: strutil.Match('%s', '%s', %d) should be %v not %v",
				tc, c.pattern, c.input, c.flags, c.want, got,
			)
		}
	}
}

func TestLiteral(t *testing.T) {
	cases := []struct {
		pattern string
		input   string
		flags   int
		want    bool
	}{
		//Literal characters should match themselves
		{"foo", "foo", 0, true},
		{"foo", "foobar", 0, false},
		{"foobar", "foo", 0, false},
		{"foo", "Foo", 0, false},
		{"わたし", "わたし", 0, true},
		// And perform case-folding when FNM_CASEFOLD is given
		{"foo", "FOO", strutil.FNM_CASEFOLD, true},
		{"FoO", "fOo", strutil.FNM_CASEFOLD, true},
	}

	for tc, c := range cases {
		got := strutil.Match(c.pattern, c.input, c.flags)
		if got != c.want {
			t.Errorf(
				"Testcase #%d failed: strutil.Match('%s', '%s', %d) should be %v not %v",
				tc, c.pattern, c.input, c.flags, c.want, got,
			)
		}
	}
}

func TestFNMLeadingDir(t *testing.T) {
	cases := []struct {
		pattern string
		input   string
		flags   int
		want    bool
	}{
		// FNM_LEADING_DIR should ignore trailing '/*'
		{"foo", "foo/bar", 0, false},
		{"foo", "foo/bar", strutil.FNM_LEADING_DIR, true},
		{"*", "foo/bar", strutil.FNM_PATHNAME, false},
		{"*", "foo/bar", strutil.FNM_PATHNAME | strutil.FNM_LEADING_DIR, true},
	}

	for tc, c := range cases {
		got := strutil.Match(c.pattern, c.input, c.flags)
		if got != c.want {
			t.Errorf(
				"Testcase #%d failed: strutil.Match('%s', '%s', %d) should be %v not %v",
				tc, c.pattern, c.input, c.flags, c.want, got,
			)
		}
	}
}
