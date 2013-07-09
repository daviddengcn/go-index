package index

import (
	"bytes"
	"github.com/russross/blackfriday"

//	"log"
)

// data-structure for a link
type Link struct {
	URL    string
	Anchor string
	Title  string
}

// Parsed data for a markdown text.
type MarkdownData struct {
	Text  []byte // plain text
	Links []Link // all links
}

type markdownData struct {
	*MarkdownData
}

// block-level callbacks
func (*markdownData) BlockCode(out *bytes.Buffer, text []byte, lang string) {}
func (*markdownData) BlockQuote(out *bytes.Buffer, text []byte)             {}
func (*markdownData) BlockHtml(out *bytes.Buffer, text []byte)              {}
func (md *markdownData) Header(out *bytes.Buffer, text func() bool, level int) {
	if text() {
		out.WriteRune('\n')
	}
}
func (*markdownData) HRule(out *bytes.Buffer) {}
func (*markdownData) List(out *bytes.Buffer, text func() bool, flags int) {
	if text() {
		out.WriteRune('\n')
	}
}
func (md *markdownData) ListItem(out *bytes.Buffer, text []byte, flags int) {
	out.WriteRune('\n')
	out.Write(text)
	out.WriteRune('\n')
}
func (md *markdownData) Paragraph(out *bytes.Buffer, text func() bool) {
	if text() {
		out.WriteRune('\n')
	}
}
func (*markdownData) Table(out *bytes.Buffer, header []byte, body []byte, columnData []int) {}
func (*markdownData) TableRow(out *bytes.Buffer, text []byte)                               {}
func (*markdownData) TableCell(out *bytes.Buffer, text []byte, flags int)                   {}
func (*markdownData) Footnotes(out *bytes.Buffer, text func() bool) {
	text()
}
func (*markdownData) FootnoteItem(out *bytes.Buffer, name, text []byte, flags int) {}

// Span-level callbacks
func (md *markdownData) AutoLink(out *bytes.Buffer, link []byte, kind int) {
	md.Links = append(md.Links, Link{
		URL: string(link),
	})
}
func (*markdownData) CodeSpan(out *bytes.Buffer, text []byte) {}
func (*markdownData) DoubleEmphasis(out *bytes.Buffer, text []byte) {
	out.Write(text)
}
func (*markdownData) Emphasis(out *bytes.Buffer, text []byte) {
	out.Write(text)
}
func (*markdownData) Image(out *bytes.Buffer, link []byte, title []byte, alt []byte) {}
func (*markdownData) LineBreak(out *bytes.Buffer)                                    {}
func (md *markdownData) Link(out *bytes.Buffer, link []byte, title []byte, content []byte) {
	out.Write(content)
	md.Links = append(md.Links, Link{
		URL:    string(link),
		Anchor: string(content),
		Title:  string(title),
	})
}
func (*markdownData) RawHtmlTag(out *bytes.Buffer, tag []byte) {}
func (*markdownData) TripleEmphasis(out *bytes.Buffer, text []byte) {
	out.Write(text)
}
func (*markdownData) StrikeThrough(out *bytes.Buffer, text []byte)      {}
func (*markdownData) FootnoteRef(out *bytes.Buffer, ref []byte, id int) {}

// Low-level callbacks
func (*markdownData) Entity(out *bytes.Buffer, entity []byte) {}
func (*markdownData) NormalText(out *bytes.Buffer, text []byte) {
	out.Write(text)
}

// Header and footer
func (*markdownData) DocumentHeader(out *bytes.Buffer) {}
func (*markdownData) DocumentFooter(out *bytes.Buffer) {}

// ParseMarkdown parses the markdown source and returns the plain text and link
// information.
func ParseMarkdown(src []byte) *MarkdownData {
	md := &MarkdownData{}

	md.Text = blackfriday.Markdown(src, &markdownData{md},
		blackfriday.EXTENSION_AUTOLINK)

	return md
}
