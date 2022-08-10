package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseTagKeyValue(t *testing.T) ***REMOVED***
	t.Parallel()
	testData := []struct ***REMOVED***
		input string
		name  string
		value string
		err   error
	***REMOVED******REMOVED***
		***REMOVED***
			"",
			"",
			"",
			errTagEmptyString,
		***REMOVED***,
		***REMOVED***
			"=",
			"",
			"",
			errTagEmptyName,
		***REMOVED***,
		***REMOVED***
			"=test",
			"",
			"",
			errTagEmptyName,
		***REMOVED***,
		***REMOVED***
			"test",
			"",
			"",
			errTagEmptyValue,
		***REMOVED***,
		***REMOVED***
			"test=",
			"",
			"",
			errTagEmptyValue,
		***REMOVED***,
		***REMOVED***
			"myTag=foo",
			"myTag",
			"foo",
			nil,
		***REMOVED***,
	***REMOVED***

	for _, data := range testData ***REMOVED***
		data := data
		t.Run(data.input, func(t *testing.T) ***REMOVED***
			t.Parallel()
			name, value, err := parseTagNameValue(data.input)
			assert.Equal(t, name, data.name)
			assert.Equal(t, value, data.value)
			assert.Equal(t, err, data.err)
		***REMOVED***)
	***REMOVED***
***REMOVED***
