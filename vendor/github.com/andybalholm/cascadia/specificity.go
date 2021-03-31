package cascadia

// Specificity is the CSS specificity as defined in
// https://www.w3.org/TR/selectors/#specificity-rules
// with the convention Specificity = [A,B,C].
type Specificity [3]int

// returns `true` if s < other (strictly), false otherwise
func (s Specificity) Less(other Specificity) bool ***REMOVED***
	for i := range s ***REMOVED***
		if s[i] < other[i] ***REMOVED***
			return true
		***REMOVED***
		if s[i] > other[i] ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

func (s Specificity) Add(other Specificity) Specificity ***REMOVED***
	for i, sp := range other ***REMOVED***
		s[i] += sp
	***REMOVED***
	return s
***REMOVED***
