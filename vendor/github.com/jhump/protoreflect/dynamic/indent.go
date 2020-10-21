package dynamic

import "bytes"

type indentBuffer struct ***REMOVED***
	bytes.Buffer
	indent      string
	indentCount int
	comma       bool
***REMOVED***

func (b *indentBuffer) start() error ***REMOVED***
	if b.indentCount >= 0 ***REMOVED***
		b.indentCount++
		return b.newLine(false)
	***REMOVED***
	return nil
***REMOVED***

func (b *indentBuffer) sep() error ***REMOVED***
	if b.indentCount >= 0 ***REMOVED***
		_, err := b.WriteString(": ")
		return err
	***REMOVED*** else ***REMOVED***
		return b.WriteByte(':')
	***REMOVED***
***REMOVED***

func (b *indentBuffer) end() error ***REMOVED***
	if b.indentCount >= 0 ***REMOVED***
		b.indentCount--
		return b.newLine(false)
	***REMOVED***
	return nil
***REMOVED***

func (b *indentBuffer) maybeNext(first *bool) error ***REMOVED***
	if *first ***REMOVED***
		*first = false
		return nil
	***REMOVED*** else ***REMOVED***
		return b.next()
	***REMOVED***
***REMOVED***

func (b *indentBuffer) next() error ***REMOVED***
	if b.indentCount >= 0 ***REMOVED***
		return b.newLine(b.comma)
	***REMOVED*** else if b.comma ***REMOVED***
		return b.WriteByte(',')
	***REMOVED*** else ***REMOVED***
		return b.WriteByte(' ')
	***REMOVED***
***REMOVED***

func (b *indentBuffer) newLine(comma bool) error ***REMOVED***
	if comma ***REMOVED***
		err := b.WriteByte(',')
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	err := b.WriteByte('\n')
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	for i := 0; i < b.indentCount; i++ ***REMOVED***
		_, err := b.WriteString(b.indent)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***
