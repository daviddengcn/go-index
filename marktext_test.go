package index

import (
	"bytes"
	"testing"
	"unicode"
	"github.com/daviddengcn/go-villa"
//	"fmt"
)

func TestMarkText(t *testing.T) {
	text := "Hello myFriend"
	
	var outBuf bytes.Buffer
	err := MarkText([]byte(text), func(last, current rune) RuneType {
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
		outBuf.Write(text)
		return nil
	}, func(text []byte) error {
		outBuf.WriteRune('<')
		outBuf.Write(text)
		outBuf.WriteRune('>')
		return nil
	})
	
	if err != nil {
		t.Errorf("MarkText failed: %v", err)
	}
	
	marked := outBuf.String()
	
	villa.AssertEquals(t, "marked", marked, "<Hello> <my><Friend>")
}