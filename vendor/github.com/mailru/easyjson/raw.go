package easyjson

import (
	"github.com/mailru/easyjson/jlexer"
	"github.com/mailru/easyjson/jwriter"
)

// RawMessage is a raw piece of JSON (number, string, bool, object, array or
// null) that is extracted without parsing and output as is during marshaling.
type RawMessage []byte

// MarshalEasyJSON does JSON marshaling using easyjson interface.
func (v *RawMessage) MarshalEasyJSON(w *jwriter.Writer) ***REMOVED***
	if len(*v) == 0 ***REMOVED***
		w.RawString("null")
	***REMOVED*** else ***REMOVED***
		w.Raw(*v, nil)
	***REMOVED***
***REMOVED***

// UnmarshalEasyJSON does JSON unmarshaling using easyjson interface.
func (v *RawMessage) UnmarshalEasyJSON(l *jlexer.Lexer) ***REMOVED***
	*v = RawMessage(l.Raw())
***REMOVED***

// UnmarshalJSON implements encoding/json.Unmarshaler interface.
func (v *RawMessage) UnmarshalJSON(data []byte) error ***REMOVED***
	*v = data
	return nil
***REMOVED***

var nullBytes = []byte("null")

// MarshalJSON implements encoding/json.Marshaler interface.
func (v RawMessage) MarshalJSON() ([]byte, error) ***REMOVED***
	if len(v) == 0 ***REMOVED***
		return nullBytes, nil
	***REMOVED***
	return v, nil
***REMOVED***

// IsDefined is required for integration with omitempty easyjson logic.
func (v *RawMessage) IsDefined() bool ***REMOVED***
	return len(*v) > 0
***REMOVED***
