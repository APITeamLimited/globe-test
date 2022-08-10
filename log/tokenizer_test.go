package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenizer(t *testing.T) ***REMOVED***
	tokens, err := tokenize("loki=something,s.e=2231,s=12,12=3,a=[1,2,3],b=[1],s=c")
	assert.Equal(t, []token***REMOVED***
		***REMOVED***
			key:   "loki",
			value: "something",
		***REMOVED***,
		***REMOVED***
			key:   "s.e",
			value: "2231",
		***REMOVED***,
		***REMOVED***
			key:   "s",
			value: "12",
		***REMOVED***,
		***REMOVED***
			key:   "12",
			value: "3",
		***REMOVED***,
		***REMOVED***
			key:    "a",
			value:  "1,2,3",
			inside: '[',
		***REMOVED***,
		***REMOVED***
			key:    "b",
			value:  "1",
			inside: '[',
		***REMOVED***,
		***REMOVED***
			key:   "s",
			value: "c",
		***REMOVED***,
	***REMOVED***, tokens)
	assert.NoError(t, err)

	_, err = tokenize("empty=")
	assert.EqualError(t, err, "key `empty=` with no value")
***REMOVED***
