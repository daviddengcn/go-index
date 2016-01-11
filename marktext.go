package index

import (
	"unicode/utf8"
)

/*
	MarkText seperates text into separator parts and tokens, mark a token if the
	token needMark. output and mark functions are called for unmarked/marked
	texts.
*/
func MarkText(text []byte, runeType func(last, current rune) RuneType,
	needMark func([]byte) bool, output, mark func([]byte) error) error {
	if len(text) == 0 {
		return nil
	}
	r, sz := utf8.DecodeRune(text)
	tp := runeType(rune(0), r)
	for {
		// text is always non-empty here.
		// r, sz are current rune and its size. tp is r's RuneType.
		p := 0

		// seperator part, if any
		for tp == TokenSep {
			// step over this rune
			p += sz
			if p < len(text) {
				lastR := r
				r, sz = utf8.DecodeRune(text[p:])
				tp = runeType(lastR, r)
			} else {
				break
			}
		}
		// p is the first non-separator position
		if p > 0 {
			// output separator part
			if err := output(text[:p]); err != nil {
				return err
			}

			if p == len(text) {
				// text ends with a separator
				break
			}

			// skip
			text, p = text[p:], 0
		}
		// p equals 0 here, text is non-empty, tp is not TokenSep

		// word part
		for p == 0 || tp == TokenBody {
			p += sz

			if p < len(text) {
				lastR := r
				r, sz = utf8.DecodeRune(text[p:])
				tp = runeType(lastR, r)
			} else {
				break
			}
		}
		// tp == TokenStart/TokenSep/text is over

		token := text[:p]
		// output a token, mark if needed
		if needMark(token) {
			if err := mark(token); err != nil {
				return err
			}
		} else {
			if err := output(token); err != nil {
				return err
			}
		}

		if p == len(text) {
			// text ends with a token
			break
		}
		text = text[p:]
	}

	return nil
}
