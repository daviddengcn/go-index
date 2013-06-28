package index

import (
	"bytes"
	"github.com/daviddengcn/go-villa"
	"testing"
	"unicode"

//	"fmt"
)

func TestTokenizer(t *testing.T) {
	text := "abc de'f  ghi\tjk"
	var tokens []string
	err := Tokenize(func(last, current rune) RuneType {
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
	}, bytes.NewReader([]byte(text)), func(token []byte) error {
		tokens = append(tokens, string(token))
		return nil
	})

	if err != nil {
		t.Errorf("Tokenize failed: %v", err)
	}

	//t.Print(text, "->", tokens)
	villa.AssertStringEquals(t, "tokens", tokens, "[abc de ' f ghi jk]")
}
