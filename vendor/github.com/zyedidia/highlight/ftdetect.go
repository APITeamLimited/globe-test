package highlight

// DetectFiletype will use the list of syntax definitions provided and the filename and first line of the file
// to determine the filetype of the file
// It will return the corresponding syntax definition for the filetype
func DetectFiletype(defs []*Def, filename string, firstLine []byte) *Def ***REMOVED***
	for _, d := range defs ***REMOVED***
		if d.ftdetect[0].MatchString(filename) ***REMOVED***
			return d
		***REMOVED***
		if len(d.ftdetect) > 1 ***REMOVED***
			if d.ftdetect[1].MatchString(string(firstLine)) ***REMOVED***
				return d
			***REMOVED***
		***REMOVED***
	***REMOVED***

	emptyDef := new(Def)
	emptyDef.FileType = "Unknown"
	emptyDef.rules = new(rules)
	return emptyDef
***REMOVED***
