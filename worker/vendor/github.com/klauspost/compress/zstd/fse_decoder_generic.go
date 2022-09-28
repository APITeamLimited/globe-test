//go:build !amd64 || appengine || !gc || noasm
// +build !amd64 appengine !gc noasm

package zstd

import (
	"errors"
	"fmt"
)

// buildDtable will build the decoding table.
func (s *fseDecoder) buildDtable() error ***REMOVED***
	tableSize := uint32(1 << s.actualTableLog)
	highThreshold := tableSize - 1
	symbolNext := s.stateTable[:256]

	// Init, lay down lowprob symbols
	***REMOVED***
		for i, v := range s.norm[:s.symbolLen] ***REMOVED***
			if v == -1 ***REMOVED***
				s.dt[highThreshold].setAddBits(uint8(i))
				highThreshold--
				symbolNext[i] = 1
			***REMOVED*** else ***REMOVED***
				symbolNext[i] = uint16(v)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// Spread symbols
	***REMOVED***
		tableMask := tableSize - 1
		step := tableStep(tableSize)
		position := uint32(0)
		for ss, v := range s.norm[:s.symbolLen] ***REMOVED***
			for i := 0; i < int(v); i++ ***REMOVED***
				s.dt[position].setAddBits(uint8(ss))
				position = (position + step) & tableMask
				for position > highThreshold ***REMOVED***
					// lowprob area
					position = (position + step) & tableMask
				***REMOVED***
			***REMOVED***
		***REMOVED***
		if position != 0 ***REMOVED***
			// position must reach all cells once, otherwise normalizedCounter is incorrect
			return errors.New("corrupted input (position != 0)")
		***REMOVED***
	***REMOVED***

	// Build Decoding table
	***REMOVED***
		tableSize := uint16(1 << s.actualTableLog)
		for u, v := range s.dt[:tableSize] ***REMOVED***
			symbol := v.addBits()
			nextState := symbolNext[symbol]
			symbolNext[symbol] = nextState + 1
			nBits := s.actualTableLog - byte(highBits(uint32(nextState)))
			s.dt[u&maxTableMask].setNBits(nBits)
			newState := (nextState << nBits) - tableSize
			if newState > tableSize ***REMOVED***
				return fmt.Errorf("newState (%d) outside table size (%d)", newState, tableSize)
			***REMOVED***
			if newState == uint16(u) && nBits == 0 ***REMOVED***
				// Seems weird that this is possible with nbits > 0.
				return fmt.Errorf("newState (%d) == oldState (%d) and no bits", newState, u)
			***REMOVED***
			s.dt[u&maxTableMask].setNewState(newState)
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***
