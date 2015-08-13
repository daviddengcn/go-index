package index

import (
	"io"

	"github.com/golangplus/bytes"
)

// A Tokenizer interface can tokenize text into tokens.
type Tokenizer interface {
	/*
		Tokenize separates a rune sequence into some tokens. If output returns
		a non-nil error, tokenizing stops and the error is returned.
	*/
	Tokenize(in io.RuneReader, output func(token []byte) error) error
}

/*
     ,----> TokenBody
    ////
  Hello  my  friend
  |     \___\________.> TokenSep (spaces)
  `-> TokenStart
*/
type RuneType int

const (
	TokenSep   RuneType = iota // token breaker, should ignored
	TokenStart                 // start of a new token, end current token, if any
	TokenBody                  // body of a token. It's ok for the first rune to be a TokenBody
)

// the type of func for determine RuneType give last and current runes.
type RuneTypeFunc func(last, current rune) RuneType

/*
	Tokenize separates a rune sequence into some tokens defining a RuneType
	function.
*/
func Tokenize(runeType RuneTypeFunc, in io.RuneReader,
	output func(token []byte) error) error {
	last := rune(0)
	var outBuf bytesp.ByteSlice
	for {
		current, _, err := in.ReadRune()
		if err != nil {
			break
		}
		tp := runeType(last, current)
		if tp == TokenStart || tp == TokenSep {
			// finish current
			if len(outBuf) > 0 {
				err = output([]byte(outBuf))
				if err != nil {
					return err
				}
				outBuf = outBuf[:0]
			}
		}

		if tp == TokenStart || tp == TokenBody {
			outBuf.WriteRune(current)
		}
		last = current
	}

	// finish last, if any
	if len(outBuf) > 0 {
		return output([]byte(outBuf))
	}
	return nil
}

// TokenizeBySeparators uses the runes of seps as seprators to
// tokenize in.
func TokenizeBySeparators(seps string, in io.RuneReader,
	output func(token []byte) error) error {
	isSap := make(map[rune]bool)
	for _, r := range seps {
		isSap[r] = true
	}

	return Tokenize(func(last, current rune) RuneType {
		if isSap[current] {
			return TokenSep
		}

		return TokenBody
	}, in, output)
}

/*
	SeparatorFRuneTypeF returns a rune-type function (used in func Tokenize)
	which can splits text by separators defined by func IsSeparator.
*/
func SeparatorFRuneTypeFunc(IsSeparator func(r rune) bool) RuneTypeFunc {
	return func(last, current rune) RuneType {
		if IsSeparator(current) {
			return TokenSep
		}

		return TokenBody
	}
}
