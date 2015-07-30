package index

import (
	"encoding/gob"
	"io"
	"sort"

	"github.com/golangplus/sort"
	"github.com/golangplus/strings"

	"github.com/daviddengcn/go-villa"
)

// TokenIndexer is main used to compute outlinks from inlinks.
type TokenIndexer struct {
	// key -> sorted tokens
	idTokens map[string][]string
	// token -> sorted keys
	tokenIds map[string][]string
}

func removeFromSortedString(l []string, el string) []string {
	i := sort.SearchStrings(l, el)
	if i >= len(l) || l[i] != el {
		return l
	}

	return stringsp.SliceRemove(l, i)
}

func addToSortedString(l []string, el string) []string {
	i := sort.SearchStrings(l, el)
	if i < len(l) && l[i] == el {
		return l
	}
	return stringsp.SliceInsert(l, i, el)
}

// Put sets the tokens for a specified id. If the id was put before, the tokens
// are updated.
func (ti *TokenIndexer) Put(id string, tokens villa.StrSet) {
	oldTokens := ti.idTokens[id]

	newTokens := tokens.Elements()
	sort.Strings(newTokens)

	if ti.tokenIds == nil {
		ti.tokenIds = make(map[string][]string)
	}

	sortp.DiffSortedList(len(oldTokens), len(newTokens), func(o, n int) int {
		return stringsp.Compare(oldTokens[o], newTokens[n])
	}, func(o int) {
		// To remove
		token := oldTokens[o]
		ti.tokenIds[token] = removeFromSortedString(ti.tokenIds[token], id)
	}, func(n int) {
		// To add
		token := newTokens[n]
		ti.tokenIds[token] = addToSortedString(ti.tokenIds[token], id)
	})

	if ti.idTokens == nil {
		ti.idTokens = make(map[string][]string)
	}

	ti.idTokens[id] = newTokens
}

// IdsOfToken returns a sorted slice of ids for a specified token.
//
// NOTE Do NOT change the elements of the returned slice
func (ti *TokenIndexer) IdsOfToken(token string) []string {
	return ti.tokenIds[token]
}

// TokensOfId returns a sorted slice of tokens for a specified id.
//
// NOTE Do NOT change the elements of the returned slice
func (ti *TokenIndexer) TokensOfId(id string) []string {
	return ti.idTokens[id]
}

// Saves serializes the TokenIndexer data to a Writer with the gob encoder.
func (ti *TokenIndexer) Save(w io.Writer) error {
	enc := gob.NewEncoder(w)
	if err := enc.Encode(ti.idTokens); err != nil {
		return err
	}
	if err := enc.Encode(ti.tokenIds); err != nil {
		return err
	}
	return nil
}

// Load restores the TokenIndexer data from a Reader with the gob decoder.
func (ti *TokenIndexer) Load(r io.Reader) error {
	*ti = TokenIndexer{}

	dec := gob.NewDecoder(r)
	if err := dec.Decode(&(ti.idTokens)); err != nil {
		return err
	}
	if err := dec.Decode(&(ti.tokenIds)); err != nil {
		return err
	}

	return nil
}
