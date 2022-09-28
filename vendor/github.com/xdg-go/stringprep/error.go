package stringprep

import "fmt"

// Error describes problems encountered during stringprep, including what rune
// was problematic.
type Error struct ***REMOVED***
	Msg  string
	Rune rune
***REMOVED***

func (e Error) Error() string ***REMOVED***
	return fmt.Sprintf("%s (rune: '\\u%04x')", e.Msg, e.Rune)
***REMOVED***
