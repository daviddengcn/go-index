package index

import (
	"io"
	"testing"
	"unicode"

	"github.com/golangplus/bytes"
	"github.com/golangplus/testing/assert"
)

func TestMarkText_Empty(t *testing.T) {
	assert.NoError(t, MarkText(nil, nil, nil, nil, nil))
}

func TestMarkText_AllSeparators(t *testing.T) {
	var out bytesp.Slice
	assert.NoError(t, MarkText([]byte("Hello"), func(last, current rune) RuneType {
		return TokenSep
	}, func(token []byte) bool {
		return true
	}, func(text []byte) error {
		out.Write(text)
		return nil
	}, func(text []byte) error {
		out.WriteRune('<')
		out.Write(text)
		out.WriteRune('>')
		return nil
	}))
	assert.Equal(t, "out", string(out), "Hello")
}

func TestMarkText_OutputError(t *testing.T) {
	assert.Equal(t, "MarkText", MarkText([]byte("H"), func(last, current rune) RuneType {
		return TokenSep
	}, func(token []byte) bool {
		return true
	}, func(text []byte) error {
		return io.EOF
	}, func(text []byte) error {
		return nil
	}), io.EOF)
	assert.Equal(t, "MarkText", MarkText([]byte("H"), func(last, current rune) RuneType {
		return TokenBody
	}, func(token []byte) bool {
		return false
	}, func(text []byte) error {
		return io.EOF
	}, func(text []byte) error {
		return nil
	}), io.EOF)
}

func TestMarkText_MarkError(t *testing.T) {
	assert.Equal(t, "MarkText", MarkText([]byte("H"), func(last, current rune) RuneType {
		return TokenBody
	}, func(token []byte) bool {
		return true
	}, func(text []byte) error {
		return nil
	}, func(text []byte) error {
		return io.EOF
	}), io.EOF)
}

func TestMarkText(t *testing.T) {
	var out bytesp.Slice
	assert.NoError(t, MarkText([]byte("Hello myFriend"), func(last, current rune) RuneType {
		if unicode.IsSpace(current) {
			return TokenSep
		}
		if current >= 'A' && current <= 'Z' {
			return TokenStart
		}
		return TokenBody
	}, func(token []byte) bool {
		return true
	}, func(text []byte) error {
		out.Write(text)
		return nil
	}, func(text []byte) error {
		out.WriteRune('<')
		out.Write(text)
		out.WriteRune('>')
		return nil
	}))
	assert.Equal(t, "out", string(out), "<Hello> <my><Friend>")
}
