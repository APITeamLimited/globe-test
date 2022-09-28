// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package collate

import "golang.org/x/text/internal/colltab"

const blockSize = 64

func getTable(t tableIndex) *colltab.Table ***REMOVED***
	return &colltab.Table***REMOVED***
		Index: colltab.Trie***REMOVED***
			Index0:  mainLookup[:][blockSize*t.lookupOffset:],
			Values0: mainValues[:][blockSize*t.valuesOffset:],
			Index:   mainLookup[:],
			Values:  mainValues[:],
		***REMOVED***,
		ExpandElem:     mainExpandElem[:],
		ContractTries:  colltab.ContractTrieSet(mainCTEntries[:]),
		ContractElem:   mainContractElem[:],
		MaxContractLen: 18,
		VariableTop:    varTop,
	***REMOVED***
***REMOVED***

// tableIndex holds information for constructing a table
// for a certain locale based on the main table.
type tableIndex struct ***REMOVED***
	lookupOffset uint32
	valuesOffset uint32
***REMOVED***
