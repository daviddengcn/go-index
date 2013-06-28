package index

import (
	"bytes"
	"io"
)

type RuneType int

const (
	/*
	    ,----> TokenBody
	   ////
	  Hello my  friend
	  |        \________.> TokenSep (spaces)
	  `-> TokenStart
	*/
	
	TokenSep   RuneType = iota // token breaker, should ignored
	TokenStart                 // start of a new token, end current token, if any
	TokenBody                  // body of a token. It's ok for the first rune to be a TokenBody
)

/*
	Tokenize spearates a rune sequence into some tokens defining a RuneType
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
