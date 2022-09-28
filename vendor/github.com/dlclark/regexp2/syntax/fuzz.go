// +build gofuzz

package syntax

// Fuzz is the input point for go-fuzz
func Fuzz(data []byte) int ***REMOVED***
	sdata := string(data)
	tree, err := Parse(sdata, RegexOptions(0))
	if err != nil ***REMOVED***
		return 0
	***REMOVED***

	// translate it to code
	_, err = Write(tree)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***

	return 1
***REMOVED***
