package index

import (
	"github.com/daviddengcn/go-algs/ed"
	"strings"
	"testing"
)

func showText(text string) string {
	return text + "."
}

func AssertTextEquals(t *testing.T, name, act, exp string) {
	if exp == act {
		return
	}

	expLines := strings.Split(exp, "\n")
	actLines := strings.Split(act, "\n")

	t.Errorf("%s unexpected(exp: %d lines, act %d lines!", name, len(expLines), len(actLines))
	t.Logf("exp ---  act +++")
	t.Logf("Difference:")
	_, matA, matB := ed.EditDistanceFFull(len(expLines), len(actLines), func(iA, iB int) int {
		sa, sb := expLines[iA], actLines[iB]
		if sa == sb {
			return 0
		}
		return ed.String(sa, sb)
	}, func(iA int) int {
		return len(expLines[iA]) + 1
	}, func(iB int) int {
		return len(actLines[iB]) + 1
	})
	for i, j := 0, 0; i < len(expLines) || j < len(actLines); {
		switch {
		case j >= len(actLines) || i < len(expLines) && matA[i] < 0:
			t.Logf("--- %3d: %s", i+1, showText(expLines[i]))
			i++
		case i >= len(expLines) || j < len(actLines) && matB[j] < 0:
			t.Logf("+++ %3d: %s", j+1, showText(actLines[j]))
			j++
		default:
			if expLines[i] != actLines[j] {
				t.Logf("--- %3d: %s", i+1, showText(expLines[i]))
				t.Logf("+++ %3d: %s", j+1, showText(actLines[j]))
			} // else
			i++
			j++
		}
	} // for i, j
}

func TestParseMarkdown(t *testing.T) {
	src :=
		`h1 text
========
_Introduction_ __to__ an [example](http://example.com/) http://www.example.com/

* L1
 * L2
Hello
Go [Go][go]

h2 text
-------
` + "```go\n" +
			"var i int\n" +
			"```" + `
[go]: http://golang.org/ "Golang"
`
	md := ParseMarkdown([]byte(src))

	t.Logf("Links:\n")
	for i, link := range md.Links {
		t.Logf("%3d: %+v\n", i, link)
	}

	MD :=
		`h1 text
Introduction to an example 

L1

L2
Hello
Go Go

h2 text

`
	AssertTextEquals(t, "Text", string(md.Text), MD)
}
