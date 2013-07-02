package index

import (
	"encoding/gob"
	"github.com/daviddengcn/go-villa"
	"io"
	"sort"
)


// TokenIndexer is main used to compute outlinks from inlinks.
type TokenIndexer struct {
	// key -> sorted tokens
	idTokens map[string][]string
	// token -> sorted keys
	tokenIds map[string][]string
}

func removeFromSortedString(l []string, el string) []string {
	pos, found := villa.StrValueCompare.BinarySearch(l, el)
	if found {
		(*villa.StringSlice)(&l).Remove(pos)
	}
	return l
}

func addToSortedString(l []string, el string) []string {
	pos, found := villa.StrValueCompare.BinarySearch(l, el)
	if !found {
		(*villa.StringSlice)(&l).Insert(pos, el)
	}
	return l
}

// Put set the tokens for a specified id. If the id was put before, the tokens
// are updated.
func (ti *TokenIndexer) Put(id string, tokens villa.StrSet) {
	oldTokens, _ := ti.idTokens[id]
	
	newTokens := tokens.Elements()
	sort.Strings(newTokens)
	
	toRemove, toAdd := villa.StrValueCompare.DiffSlicePair(oldTokens, newTokens)
	
	if ti.tokenIds == nil {
		ti.tokenIds = make(map[string][]string)
	}
	
	for _, token := range toRemove {
		ti.tokenIds[token] = removeFromSortedString(ti.tokenIds[token], id)
	}
	for _, token := range toAdd {
		ti.tokenIds[token] = addToSortedString(ti.tokenIds[token], id)
	}
	
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
