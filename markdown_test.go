package index

import (
	"testing"
)

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
