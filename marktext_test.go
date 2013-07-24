package index

import (
	"fmt"
	"github.com/daviddengcn/go-villa"
	"testing"
	"unicode"
)

func TestMarkText(t *testing.T) {
	text := "Hello myFriend"

	var outBuf villa.ByteSlice
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

	marked := string(outBuf)

	AssertEquals(t, "marked", marked, "<Hello> <my><Friend>")
}

/*
	AssertEquals shows error message when act and exp don't equal
*/
func AssertEquals(t *testing.T, name string, act, exp interface{}) {
	if act != exp {
		t.Errorf("%s is expected to be %v, but got %v", name, exp, act)
	}
}

/*
	AssertEquals shows error message when strings forms of act and exp don't
	equal
*/
func AssertStringEquals(t *testing.T, name string, act, exp interface{}) {
	if fmt.Sprintf("%v", act) != fmt.Sprintf("%v", exp) {
		t.Errorf("%s is expected to be %v, but got %v", name, exp, act)
	} // if
}

/*
	AssertStrSetEquals shows error message when act and exp are equal string
	sets.
*/
func AssertStrSetEquals(t *testing.T, name string, act, exp villa.StrSet) {
	if !act.Equals(exp) {
		t.Errorf("%s is expected to be %v, but got %v", name, exp, act)
	}
}
