package index

import (
	"bytes"
	"testing"
	"unicode"

	"github.com/golangplus/testing/assert"
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

	assert.StringEqual(t, "tokens", tokens, "[abc de ' f ghi jk]")
}

func TestSeparatorFRuneTypeFunc(t *testing.T) {
	f := SeparatorFRuneTypeFunc(unicode.IsSpace)
	assert.Equal(t, "f", f('a', ' '), TokenSep)
	assert.Equal(t, "f", f('a', 'a'), TokenBody)
}
