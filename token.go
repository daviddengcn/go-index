package index

import (
	"bytes"
	"io"
)

// A Tokenizer interface can tokenize text into tokens.
type Tokenizer interface {
	/*
		Tokenize separates a rune sequence into some tokens. If output returns
		an non-nil error, tokenizing ceased, and the error is returned.
	*/
	Tokenize(in io.RuneReader, output func(token []byte) error) error
}

type RuneType int

const (
	/*
	    ,----> TokenBody
	   ////
	  Hello  my  friend
	  |     \___\________.> TokenSep (spaces)
	  `-> TokenStart
	*/

	TokenSep   RuneType = iota // token breaker, should ignored
	TokenStart                 // start of a new token, end current token, if any
	TokenBody                  // body of a token. It's ok for the first rune to be a TokenBody
)

/*
	Tokenize separates a rune sequence into some tokens defining a RuneType
	function.
*/
func Tokenize(runeType func(last, current rune) RuneType, in io.RuneReader,
	output func(token []byte) error) error {
	last := rune(0)
	var outBuf bytes.Buffer
	for {
		current, _, err := in.ReadRune()
		if err != nil {
			break
		}
		tp := runeType(last, current)
		if tp == TokenStart || tp == TokenSep {
			// finish current
			if outBuf.Len() > 0 {
				err = output(outBuf.Bytes())
				if err != nil {
					return err
				}
				outBuf.Reset()
			}
		}

		if tp == TokenStart || tp == TokenBody {
			outBuf.WriteRune(current)
		}
		last = current
	}

	// finish last, if any
	if outBuf.Len() > 0 {
		return output(outBuf.Bytes())
	}
	return nil
}

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
