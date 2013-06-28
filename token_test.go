package index

import (
	"bytes"
	"testing"
	"unicode"
	"github.com/daviddengcn/go-villa"
//	"fmt"
)

func TestTokenizer(t *testing.T) {
	text := "abc de'f  ghi\tjk"
	var tokens []string
	Tokenize(func(last, current rune) RuneType {
		if unicode.IsSpace(current) {
			return TokenSep
		}
		
		if current == '\'' {
			return TokenStart
		}
		
		if last == '\'' {
			return TokenStart
		}
		
		return TokenBody
	}, bytes.NewReader([]byte(text)), func(token []byte) {
		tokens = append(tokens, string(token))
	})
	
	//t.Print(text, "->", tokens)
	villa.AssertStringEquals(t, "tokens", tokens, "[abc de ' f ghi jk]")
}