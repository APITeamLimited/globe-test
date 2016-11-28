package ui

import (
	"fmt"
	"strings"
)

type ProgressBar struct ***REMOVED***
	Width    int
	Progress float64
***REMOVED***

func (b *ProgressBar) String() string ***REMOVED***
	space := b.Width - 2
	filled := int(float64(space) * b.Progress)

	filling := ""
	caret := ""
	if filled > 0 ***REMOVED***
		if filled < space ***REMOVED***
			filling = strings.Repeat("=", filled-1)
			caret = ">"
		***REMOVED*** else ***REMOVED***
			filling = strings.Repeat("=", filled)
		***REMOVED***
	***REMOVED***

	padding := ""
	if space > filled ***REMOVED***
		padding = strings.Repeat(" ", space-filled)
	***REMOVED***

	return fmt.Sprintf("[%s%s%s]", filling, caret, padding)
***REMOVED***
