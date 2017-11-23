package escape

import "strings"

var (
	escaper   = strings.NewReplacer(`,`, `\,`, `"`, `\"`, ` `, `\ `, `=`, `\=`)
	unescaper = strings.NewReplacer(`\,`, `,`, `\"`, `"`, `\ `, ` `, `\=`, `=`)
)

// UnescapeString returns unescaped version of in.
func UnescapeString(in string) string ***REMOVED***
	if strings.IndexByte(in, '\\') == -1 ***REMOVED***
		return in
	***REMOVED***
	return unescaper.Replace(in)
***REMOVED***

// String returns the escaped version of in.
func String(in string) string ***REMOVED***
	return escaper.Replace(in)
***REMOVED***
