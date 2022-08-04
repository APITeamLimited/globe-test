package log

import "fmt"

type token struct ***REMOVED***
	key, value string
	inside     rune // shows whether it's inside a given collection, currently [ means it's an array
***REMOVED***

type tokenizer struct ***REMOVED***
	i          int
	s          string
	currentKey string
***REMOVED***

func (t *tokenizer) readKey() (string, error) ***REMOVED***
	start := t.i
	for ; t.i < len(t.s); t.i++ ***REMOVED***
		if t.s[t.i] == '=' && t.i != len(t.s)-1 ***REMOVED***
			t.i++

			return t.s[start : t.i-1], nil
		***REMOVED***
		if t.s[t.i] == ',' ***REMOVED***
			k := t.s[start:t.i]

			return k, fmt.Errorf("key `%s` with no value", k)
		***REMOVED***
	***REMOVED***

	s := t.s[start:]

	return s, fmt.Errorf("key `%s` with no value", s)
***REMOVED***

func (t *tokenizer) readValue() string ***REMOVED***
	start := t.i
	for ; t.i < len(t.s); t.i++ ***REMOVED***
		if t.s[t.i] == ',' ***REMOVED***
			t.i++

			return t.s[start : t.i-1]
		***REMOVED***
	***REMOVED***

	return t.s[start:]
***REMOVED***

func (t *tokenizer) readArray() (string, error) ***REMOVED***
	start := t.i
	for ; t.i < len(t.s); t.i++ ***REMOVED***
		if t.s[t.i] == ']' ***REMOVED***
			if t.i+1 == len(t.s) || t.s[t.i+1] == ',' ***REMOVED***
				t.i += 2

				return t.s[start : t.i-2], nil
			***REMOVED***
			t.i++

			return t.s[start : t.i-1], fmt.Errorf("there was no ',' after an array with key '%s'", t.currentKey)
		***REMOVED***
	***REMOVED***

	return t.s[start:], fmt.Errorf("array value for key `%s` didn't end", t.currentKey)
***REMOVED***

func tokenize(s string) ([]token, error) ***REMOVED***
	result := []token***REMOVED******REMOVED***
	t := &tokenizer***REMOVED***s: s***REMOVED***

	var err error
	var value string
	for t.i < len(s) ***REMOVED***
		t.currentKey, err = t.readKey()
		if err != nil ***REMOVED***
			return result, err
		***REMOVED***
		if t.s[t.i] == '[' ***REMOVED***
			t.i++
			value, err = t.readArray()

			result = append(result, token***REMOVED***
				key:    t.currentKey,
				value:  value,
				inside: '[',
			***REMOVED***)
			if err != nil ***REMOVED***
				return result, err
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			value = t.readValue()
			result = append(result, token***REMOVED***
				key:   t.currentKey,
				value: value,
			***REMOVED***)
		***REMOVED***
	***REMOVED***

	return result, nil
***REMOVED***
