package index

import (
	"bytes"
	"github.com/russross/blackfriday"
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

// block-level callbacks
func (*MarkdownData) BlockCode(out *bytes.Buffer, text []byte, lang string) {}
func (*MarkdownData) BlockQuote(out *bytes.Buffer, text []byte)             {}
func (*MarkdownData) BlockHtml(out *bytes.Buffer, text []byte)              {}
func (md *MarkdownData) Header(out *bytes.Buffer, text func() bool, level int) {
	if text() {
		out.WriteRune('\n')
	}
}
func (*MarkdownData) HRule(out *bytes.Buffer) {}
func (*MarkdownData) List(out *bytes.Buffer, text func() bool, flags int) {
	if text() {
		out.WriteRune('\n')
	}
}
func (md *MarkdownData) ListItem(out *bytes.Buffer, text []byte, flags int) {
	out.WriteRune('\n')
	out.Write(text)
	out.WriteRune('\n')
}
func (md *MarkdownData) Paragraph(out *bytes.Buffer, text func() bool) {
	if text() {
		out.WriteRune('\n')
	}
}
func (*MarkdownData) Table(out *bytes.Buffer, header []byte, body []byte, columnData []int) {}
func (*MarkdownData) TableRow(out *bytes.Buffer, text []byte)                               {}
func (*MarkdownData) TableCell(out *bytes.Buffer, text []byte, flags int)                   {}

// Span-level callbacks
func (md *MarkdownData) AutoLink(out *bytes.Buffer, link []byte, kind int) {
	md.Links = append(md.Links, Link{
		URL: string(link),
	})
}
func (*MarkdownData) CodeSpan(out *bytes.Buffer, text []byte) {}
func (*MarkdownData) DoubleEmphasis(out *bytes.Buffer, text []byte) {
	out.Write(text)
}
func (*MarkdownData) Emphasis(out *bytes.Buffer, text []byte) {
	out.Write(text)
}
func (*MarkdownData) Image(out *bytes.Buffer, link []byte, title []byte, alt []byte) {}
func (md *MarkdownData) LineBreak(out *bytes.Buffer)                                 {}
func (md *MarkdownData) Link(out *bytes.Buffer, link []byte, title []byte, content []byte) {
	out.Write(content)
	md.Links = append(md.Links, Link{
		URL:    string(link),
		Anchor: string(content),
		Title:  string(title),
	})
}
func (*MarkdownData) RawHtmlTag(out *bytes.Buffer, tag []byte) {}
func (*MarkdownData) TripleEmphasis(out *bytes.Buffer, text []byte) {
	out.Write(text)
}
func (*MarkdownData) StrikeThrough(out *bytes.Buffer, text []byte) {}

// Low-level callbacks
func (*MarkdownData) Entity(out *bytes.Buffer, entity []byte) {}
func (md *MarkdownData) NormalText(out *bytes.Buffer, text []byte) {
	// out.Write(bytes.Replace(text, []byte("\n"), []byte(" "), -1))
	out.Write(text)
}

// Header and footer
func (*MarkdownData) DocumentHeader(out *bytes.Buffer) {}
func (*MarkdownData) DocumentFooter(out *bytes.Buffer) {}

// ParseMarkdown parses the markdown source and returns the plain text and link
// information.
func ParseMarkdown(src []byte) *MarkdownData {
	md := &MarkdownData{}

	md.Text = blackfriday.Markdown(src, md, blackfriday.EXTENSION_AUTOLINK)

	return md
}
