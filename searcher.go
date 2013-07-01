package index

import(
	"github.com/daviddengcn/go-villa"
	"encoding/gob"
	"io"
	"math/big"
	"errors"
	
//	"fmt"
)

var (
	ErrInvalidDocID = errors.New("Invalid doc-ID")
)

// TokenSetSearcher can index documents, with which represented as a set of 
// tokens. All data are stored in memory.
//
// Indexed data can be saved, and loaded again.
//
// If a customized type needs to be saved and loaded again, it must be
// registered by gob.Register.
type TokenSetSearcher struct {
	docs []interface{}
	inverted map[string][]int32
	deleted big.Int
}

// IndexDoc indexes a document to the searcher. It returns a local doc id.
func (s *TokenSetSearcher) IndexDoc(fields map[string]villa.StrSet,
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
	s.deleted.SetBit(&s.deleted, int(docID), 1)
	return nil
}


// SingleFieldQuery returns a map[strig]villa.StrSet (same type as query int
// Search method) with a single field.
func SingleFieldQuery(field string, tokens []string) map[string]villa.StrSet {
	return map[string]villa.StrSet {
		field: villa.NewStrSet(tokens...),
	}
}

// Search ouputs all documents (docID and associated data) with all tokens
// hit, in the same order as ther were added. If output returns an nonnil error,
// the search stops, and the error is returned.
// If no tokens in query, all non-deleted documents are returned.
func (s *TokenSetSearcher) Search(query map[string]villa.StrSet, 
		output func(docID int32, data interface{})error) error {

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
			list := s.inverted[token]
			if len(list) == 0 {
				return nil
			}
			for _, docID := range list {
				if s.deleted.Bit(int(docID)) == 0 {
					err := output(docID, s.docs[docID])
					if err != nil {
						return err
					}
				}
			}
			break
		}
		return nil
	}
	
	invLists := make([][]int32, 0, len(tokens))
	for token := range tokens {
		list := s.inverted[token]
		if len(list) == 0 {
			// one of the inverted is empty, no results
			return nil
		}
		invLists = append(invLists, list)
	}
	
	
	N, n := len(s.docs), len(tokens)
	if N == 0 {
		return nil
	}
	
	gaps := make([]int32, n)
	mnI := 0
	for i := range invLists {
		gaps[i] = 2*int32(N) / int32(len(invLists[i]))
		if len(invLists[i]) < len(invLists[mnI]) {
			mnI = i
		}
	}
	mnI1 := mnI + 1
	if mnI1 == n {
		mnI1 = 0
	}
	
	// the current position in inverted lists
	hds := make([]int, len(tokens))
	docID, matched, i := invLists[mnI][0], 1, mnI1
mainloop:
	for {
		invList := invLists[i]

		if docID - invList[hds[i]] > gaps[i] {
			// estimate skip linearly
			skip := int(docID - invList[hds[i]]) * len(invList) / N
			newHd := hds[i] + skip
			if newHd >= len(invList) || invList[newHd] > docID {
				break
			}
			hds[i] = newHd
		}
		// search for docID
		for invList[hds[i]] < docID {
			hds[i] ++
			if hds[i] == len(invList) {
				// no more docs in invLists[i]
				break mainloop
			}
		}
		// skip deleted docs
		for s.deleted.Bit(int(invList[hds[i]])) != 0 {
			hds[i] ++
			if hds[i] == len(invList) {
				// no more docs in invLists[i]
				break mainloop
			}
		}
		
		if invList[hds[i]] > docID {
			// new docID
			hds[mnI] ++
			if hds[mnI] == len(invLists[mnI]) {
				break mainloop
			}
			docID, matched, i = invLists[mnI][hds[mnI]], 1, mnI1
		} else {
			matched ++
			if matched == n {
				// found a document
				err := output(docID, s.docs[docID])
				if err != nil {
					return err
				}

				/*	
				matched = 0
				docID ++
				/*/
				/*
				hds[i] ++
				if hds[i] == len(invList) {
					break mainloop
				}
				docID, matched = invList[hds[i]], 1
				//*/
				hds[mnI] ++
				if hds[mnI] == len(invLists[mnI]) {
					break mainloop
				}
				docID, matched, i = invLists[mnI][hds[mnI]], 1, mnI1
			} else {
				mnI ++
				if mnI == n {
					mnI = 0
				}
			}
		}
	}
	return nil
}

// Saves serializes the searcher data to a Writer with the gob encoder.
func (s *TokenSetSearcher) Save(w io.Writer) error {
	enc := gob.NewEncoder(w)
	err := enc.Encode(s.docs)
	if err != nil {
		return err
	}
	err = enc.Encode(s.inverted)
	if err != nil {
		return err
	}
	err = enc.Encode(s.deleted.Bytes())
	if err != nil {
		return err
	}
	return nil
}

// Load restores the searcher data from a Reader with the gob decoder.
func (s *TokenSetSearcher) Load(r io.Reader) error {
	s.docs = nil
	s.inverted = nil
	
	dec := gob.NewDecoder(r)
	err := dec.Decode(&(s.docs))
	if err != nil {
		return err
	}
	err = dec.Decode(&(s.inverted))
	if err != nil {
		return err
	}
	var bytes []byte
	err = dec.Decode(&bytes)
	if err != nil {
		return err
	}
	s.deleted.SetBytes(bytes)
	return nil
}

// DocCount returns the number of docs (included deleted).
func (s *TokenSetSearcher) DocCount() int {
	return len(s.docs)
}
