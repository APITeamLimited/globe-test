package jlexer

import "fmt"

// LexerError implements the error interface and represents all possible errors that can be
// generated during parsing the JSON data.
type LexerError struct ***REMOVED***
	Reason string
	Offset int
	Data   string
***REMOVED***

func (l *LexerError) Error() string ***REMOVED***
	return fmt.Sprintf("parse error: %s near offset %d of '%s'", l.Reason, l.Offset, l.Data)
***REMOVED***
