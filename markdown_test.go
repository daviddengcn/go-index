package index

import (
	"testing"
	
	"github.com/daviddengcn/go-assert"
	"github.com/russross/blackfriday"
)

func TestParseMarkdown_bug(t *testing.T) {
	t.Logf("%s", blackfriday.MarkdownCommon([]byte("[[t]](/t)")))
	t.Logf("%s", blackfriday.MarkdownCommon([]byte(
		"[![Build Status](https://secure.travis-ci.org/daaku/go.pqueue.png)](http://travis-ci.org/daaku/go.pqueue)")))
	
	ParseMarkdown([]byte("[[t]](/t)"))
	
	md := string(ParseMarkdown([]byte(
		"[![Build Status](https://secure.travis-ci.org/daaku/go.pqueue.png)](http://travis-ci.org/daaku/go.pqueue)")).Text)
	assert.StringEquals(t, "md", md, "[]()")
}

func TestParseMarkdown(t *testing.T) {
	src :=
		`h1 text
========
_Introduction_ __to__ an [example](http://example.com/) http://www.example.com/

* L1
 * L2

go.pqueue [![Build Status](https://secure.travis-ci.org/daaku/go.pqueue.png)](http://travis-ci.org/daaku/go.pqueue)


Hello
Go [Go][go]

h2 text
-------

` + "```go" + `
var i int
package main
` + "```\n" + `
[go]: http://golang.org/ "Golang"
`
	md := ParseMarkdown([]byte(src))

	t.Logf("Links:\n")
	for i, link := range md.Links {
		t.Logf("%3d: %+v\n", i, link)
	}
	t.Logf("act:\n%s", string(md.Text))

	MD :=
		`h1 text

Introduction to an example 

L1

L2

go.pqueue []()

Hello
Go Go

h2 text

`
	assert.TextEquals(t, "Text", string(md.Text), MD)
}
