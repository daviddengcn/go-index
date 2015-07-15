package index

import (
	"fmt"
	"testing"

	"github.com/golangplus/testing/assert"

	"github.com/russross/blackfriday"
)

func TestParseMarkdown_bug(t *testing.T) {
	t.Logf("%s", blackfriday.MarkdownCommon([]byte("[[t]](/t)")))
	t.Logf("%s", blackfriday.MarkdownCommon([]byte(
		"[![Build Status](https://secure.travis-ci.org/daaku/go.pqueue.png)](http://travis-ci.org/daaku/go.pqueue)")))

	ParseMarkdown([]byte("[[t]](/t)"))

	psd := ParseMarkdown([]byte(
		"[![Build Status](https://secure.travis-ci.org/daaku/go.pqueue.png)](http://travis-ci.org/daaku/go.pqueue)"))
	t.Logf("%+v", psd)
	md := string(psd.Text)
	assert.StringEqual(t, "md", md, " \n\n")
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
	var links []string
	for i, link := range md.Links {
		t.Logf("%3d: %+v\n", i, link)
		links = append(links, fmt.Sprintf("%+v", link))
	}
	t.Logf("act:\n%s", string(md.Text))
	assert.StringEqual(t, "links", links, []string{
		"{URL:http://example.com/ Anchor:example Title:}",
		"{URL:http://www.example.com/ Anchor: Title:}",
		"{URL:http://travis-ci.org/daaku/go.pqueue Anchor:  Title:}",
		"{URL:http://golang.org/ Anchor:Go Title:Golang}",
	})

	MD :=
		`h1 text

Introduction to an example 

L1

L2

go.pqueue  

Hello
Go Go

h2 text

`
	assert.StringEqual(t, "Text", string(md.Text), MD)
}
