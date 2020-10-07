package easyjson

import (
	jlexer "github.com/mailru/easyjson/jlexer"
	"github.com/mailru/easyjson/jwriter"
)

// UnknownFieldsProxy implemets UnknownsUnmarshaler and UnknownsMarshaler
// use it as embedded field in your structure to parse and then serialize unknown struct fields
type UnknownFieldsProxy struct ***REMOVED***
	unknownFields map[string][]byte
***REMOVED***

func (s *UnknownFieldsProxy) UnmarshalUnknown(in *jlexer.Lexer, key string) ***REMOVED***
	if s.unknownFields == nil ***REMOVED***
		s.unknownFields = make(map[string][]byte, 1)
	***REMOVED***
	s.unknownFields[key] = in.Raw()
***REMOVED***

func (s UnknownFieldsProxy) MarshalUnknowns(out *jwriter.Writer, first bool) ***REMOVED***
	for key, val := range s.unknownFields ***REMOVED***
		if first ***REMOVED***
			first = false
		***REMOVED*** else ***REMOVED***
			out.RawByte(',')
		***REMOVED***
		out.String(string(key))
		out.RawByte(':')
		out.Raw(val, nil)
	***REMOVED***
***REMOVED***
