package index

import (
	"encoding/gob"
	"errors"
	"io"
	"math/big"

	"github.com/daviddengcn/go-villa"
)

var (
	// error of invalid doc-id (out of range)
	ErrInvalidDocID = errors.New("Invalid doc-ID")
)

// TokenSetSearcher can index documents, with which represented as a set of
// tokens. All data are stored in memory.
//
// Indexed data can be saved, and loaded again.
//
// If a customized type needs to be saved and loaded again, it must be
// registered by calling gob.Register.
type TokenSetSearcher struct {
	docs []interface{}
	// map from token to list of local IDs(indexes in docs field)
	inverted     map[string][]int32
	deleted      big.Int
	deletedCount int
}

// AddDoc indexes a document to the searcher. It returns a local doc id.
func (s *TokenSetSearcher) AddDoc(fields map[string]villa.StrSet,
	data interface{}) int32 {

	docID := int32(len(s.docs))
	s.docs = append(s.docs, data)
	if s.inverted == nil {
		s.inverted = make(map[string][]int32)
	}
	for fld, tokens := range fields {
		for token := range tokens {
			key := fld + ":" + token
			s.inverted[key] = append(s.inverted[key], docID)
		}
	}

	return docID
}

// Delete marks a specified doc as deleted.
func (s *TokenSetSearcher) Delete(docID int32) error {
	if docID < 0 || docID >= int32(len(s.docs)) {
		return ErrInvalidDocID
	}

	if s.deleted.Bit(int(docID)) != 0 {
		// already deleted
		return nil
	}
	s.deleted.SetBit(&s.deleted, int(docID), 1)
	s.deletedCount++
	return nil
}

// SingleFieldQuery returns a map[strig]villa.StrSet (same type as query int
// Search method) with a single field.
func SingleFieldQuery(field string, tokens ...string) map[string]villa.StrSet {
	return map[string]villa.StrSet{
		field: villa.NewStrSet(tokens...),
	}
}

// Search outputs all documents (docID and associated data) with all tokens
// hit, in the same order as ther were added. If output returns an nonnil error,
// the search stops, and the error is returned.
// If no tokens in query, all non-deleted documents are returned.
func (s *TokenSetSearcher) Search(query map[string]villa.StrSet,
	output func(docID int32, data interface{}) error) error {

	var tokens villa.StrSet
	for fld, tks := range query {
		for tk := range tks {
			key := fld + ":" + tk
			tokens.Put(key)
		}
	}
	if len(tokens) == 0 {
		// returns all non-deleted documents
		for docID := range s.docs {
			if s.deleted.Bit(docID) == 0 {
				err := output(int32(docID), s.docs[docID])
				if err != nil {
					return err
				}
			}
		}
		return nil
	}

	if len(tokens) == 1 {
		// for single token, iterating over the inverted list
		for token := range tokens {
			for _, docID := range s.inverted[token] {
				if s.deleted.Bit(int(docID)) == 0 {
					err := output(docID, s.docs[docID])
					if err != nil {
						return err
					}
				}
			}
		}
		return nil
	}

	N, n := len(s.docs), len(tokens)
	if N == 0 {
		return nil
	}

	mnI := 0
	invLists := make([][]int32, 0, n)
	for token := range tokens {
		list := s.inverted[token]
		if len(list) == 0 {
			// one of the inverted is empty, no results
			return nil
		}
		invLists = append(invLists, list)
		if len(list) < len(invLists[mnI]) {
			mnI = len(invLists) - 1
		}
	}
	// mnI1 is the index next to mnI
	mnI1 := (mnI + 1) % n

	// gaps is the minimum difference of docID that may cause a skip
	gaps := make([]int32, n)
	for i := range invLists {
		gaps[i] = 2 * int32(N) / int32(len(invLists[i]))
	}

	// the current indexes in inverted lists
	idxs := make([]int, n)

	docID, matched, i := invLists[mnI][0], 1, mnI1
mainloop:
	for {
		invList := invLists[i]

		if docID-invList[idxs[i]] > gaps[i] {
			// estimate skip linearly
			skip := int64(docID-invList[idxs[i]]) * int64(len(invList)) / int64(N)
			newIdx := idxs[i] + int(skip)
			if newIdx < len(invList) && invList[newIdx] <= docID {
				idxs[i] = newIdx
			}
		}
		// search for docID
		for invList[idxs[i]] < docID {
			idxs[i]++
			if idxs[i] == len(invList) {
				// no more docs in invLists[i]
				break mainloop
			}
		}
		// skip deleted docs
		for s.deleted.Bit(int(invList[idxs[i]])) != 0 {
			idxs[i]++
			if idxs[i] == len(invList) {
				// no more docs in invLists[i]
				break mainloop
			}
		}

		if invList[idxs[i]] > docID {
			// move to next docID in mnI list
			idxs[mnI]++
			if idxs[mnI] == len(invLists[mnI]) {
				break mainloop
			}
			docID, matched, i = invLists[mnI][idxs[mnI]], 1, mnI1
		} else {
			matched++
			if matched == n {
				// found a document
				err := output(docID, s.docs[docID])
				if err != nil {
					return err
				}

				// move to next docID in mnI list
				idxs[mnI]++
				if idxs[mnI] == len(invLists[mnI]) {
					break mainloop
				}
				docID, matched, i = invLists[mnI][idxs[mnI]], 1, mnI1
			} else {
				i++
				if i == n {
					i = 0
				}
			}
		}
	}
	return nil
}

// Save serializes the searcher data to a Writer with the gob encoder.
func (s *TokenSetSearcher) Save(w io.Writer) error {
	enc := gob.NewEncoder(w)

	if err := enc.Encode(len(s.docs)); err != nil {
		return err
	}

	for i := range s.docs {
		if err := enc.Encode(&s.docs[i]); err != nil {
			return err
		}
	}

	if err := enc.Encode(len(s.inverted)); err != nil {
		return err
	}
	for token, ids := range s.inverted {
		if err := enc.Encode(token); err != nil {
			return err
		}
		if err := enc.Encode(ids); err != nil {
			return err
		}
	}

	if err := enc.Encode(s.deleted.Bytes()); err != nil {
		return err
	}
	if err := enc.Encode(s.deletedCount); err != nil {
		return err
	}
	return nil
}

// Load restores the searcher data from a Reader with the gob decoder.
func (s *TokenSetSearcher) Load(r io.Reader) error {
	*s = TokenSetSearcher{}

	dec := gob.NewDecoder(r)

	var docsLen int
	if err := dec.Decode(&docsLen); err != nil {
		return err
	}

	s.docs = make([]interface{}, docsLen)
	for i := 0; i < docsLen; i++ {
		if err := dec.Decode(&s.docs[i]); err != nil {
			return err
		}
	}

	var invLen int
	if err := dec.Decode(&invLen); err != nil {
		return err
	}
	if invLen > 0 {
		s.inverted = make(map[string][]int32)
		for i := 0; i < invLen; i++ {
			var token string
			var ids []int32
			if err := dec.Decode(&token); err != nil {
				return err
			}
			if err := dec.Decode(&ids); err != nil {
				return err
			}
			s.inverted[token] = ids
		}
	}

	var bytes []byte
	if err := dec.Decode(&bytes); err != nil {
		return err
	}
	s.deleted.SetBytes(bytes)

	if err := dec.Decode(&(s.deletedCount)); err != nil {
		return err
	}
	return nil
}

// DocInfo returns the doc-info of specified doc
func (s *TokenSetSearcher) DocInfo(docID int32) interface{} {
	if docID < 0 || docID >= int32(len(s.docs)) {
		return ErrInvalidDocID
	}

	return s.docs[docID]
}

// DocCount returns the number of (non-deleted) docs.
func (s *TokenSetSearcher) DocCount() int {
	return len(s.docs) - s.deletedCount
}

// DeletedCount returns the number of documents marked as deleted
func (s *TokenSetSearcher) DeletedCount() int {
	return s.deletedCount
}

// Returns the docIDs of a speicified token.
func (s *TokenSetSearcher) TokenDocList(field, token string) []int32 {
	return s.inverted[field+":"+token]
}
