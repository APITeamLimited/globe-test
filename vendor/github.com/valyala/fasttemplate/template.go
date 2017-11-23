// Package fasttemplate implements simple and fast template library.
//
// Fasttemplate is faster than text/template, strings.Replace
// and strings.Replacer.
//
// Fasttemplate ideally fits for fast and simple placeholders' substitutions.
package fasttemplate

import (
	"bytes"
	"fmt"
	"github.com/valyala/bytebufferpool"
	"io"
)

// ExecuteFunc calls f on each template tag (placeholder) occurrence.
//
// Returns the number of bytes written to w.
//
// This function is optimized for constantly changing templates.
// Use Template.ExecuteFunc for frozen templates.
func ExecuteFunc(template, startTag, endTag string, w io.Writer, f TagFunc) (int64, error) ***REMOVED***
	s := unsafeString2Bytes(template)
	a := unsafeString2Bytes(startTag)
	b := unsafeString2Bytes(endTag)

	var nn int64
	var ni int
	var err error
	for ***REMOVED***
		n := bytes.Index(s, a)
		if n < 0 ***REMOVED***
			break
		***REMOVED***
		ni, err = w.Write(s[:n])
		nn += int64(ni)
		if err != nil ***REMOVED***
			return nn, err
		***REMOVED***

		s = s[n+len(a):]
		n = bytes.Index(s, b)
		if n < 0 ***REMOVED***
			// cannot find end tag - just write it to the output.
			ni, _ = w.Write(a)
			nn += int64(ni)
			break
		***REMOVED***

		ni, err = f(w, unsafeBytes2String(s[:n]))
		nn += int64(ni)
		s = s[n+len(b):]
	***REMOVED***
	ni, err = w.Write(s)
	nn += int64(ni)

	return nn, err
***REMOVED***

// Execute substitutes template tags (placeholders) with the corresponding
// values from the map m and writes the result to the given writer w.
//
// Substitution map m may contain values with the following types:
//   * []byte - the fastest value type
//   * string - convenient value type
//   * TagFunc - flexible value type
//
// Returns the number of bytes written to w.
//
// This function is optimized for constantly changing templates.
// Use Template.Execute for frozen templates.
func Execute(template, startTag, endTag string, w io.Writer, m map[string]interface***REMOVED******REMOVED***) (int64, error) ***REMOVED***
	return ExecuteFunc(template, startTag, endTag, w, func(w io.Writer, tag string) (int, error) ***REMOVED*** return stdTagFunc(w, tag, m) ***REMOVED***)
***REMOVED***

// ExecuteFuncString calls f on each template tag (placeholder) occurrence
// and substitutes it with the data written to TagFunc's w.
//
// Returns the resulting string.
//
// This function is optimized for constantly changing templates.
// Use Template.ExecuteFuncString for frozen templates.
func ExecuteFuncString(template, startTag, endTag string, f TagFunc) string ***REMOVED***
	tagsCount := bytes.Count(unsafeString2Bytes(template), unsafeString2Bytes(startTag))
	if tagsCount == 0 ***REMOVED***
		return template
	***REMOVED***

	bb := byteBufferPool.Get()
	if _, err := ExecuteFunc(template, startTag, endTag, bb, f); err != nil ***REMOVED***
		panic(fmt.Sprintf("unexpected error: %s", err))
	***REMOVED***
	s := string(bb.B)
	bb.Reset()
	byteBufferPool.Put(bb)
	return s
***REMOVED***

var byteBufferPool bytebufferpool.Pool

// ExecuteString substitutes template tags (placeholders) with the corresponding
// values from the map m and returns the result.
//
// Substitution map m may contain values with the following types:
//   * []byte - the fastest value type
//   * string - convenient value type
//   * TagFunc - flexible value type
//
// This function is optimized for constantly changing templates.
// Use Template.ExecuteString for frozen templates.
func ExecuteString(template, startTag, endTag string, m map[string]interface***REMOVED******REMOVED***) string ***REMOVED***
	return ExecuteFuncString(template, startTag, endTag, func(w io.Writer, tag string) (int, error) ***REMOVED*** return stdTagFunc(w, tag, m) ***REMOVED***)
***REMOVED***

// Template implements simple template engine, which can be used for fast
// tags' (aka placeholders) substitution.
type Template struct ***REMOVED***
	template string
	startTag string
	endTag   string

	texts          [][]byte
	tags           []string
	byteBufferPool bytebufferpool.Pool
***REMOVED***

// New parses the given template using the given startTag and endTag
// as tag start and tag end.
//
// The returned template can be executed by concurrently running goroutines
// using Execute* methods.
//
// New panics if the given template cannot be parsed. Use NewTemplate instead
// if template may contain errors.
func New(template, startTag, endTag string) *Template ***REMOVED***
	t, err := NewTemplate(template, startTag, endTag)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	return t
***REMOVED***

// NewTemplate parses the given template using the given startTag and endTag
// as tag start and tag end.
//
// The returned template can be executed by concurrently running goroutines
// using Execute* methods.
func NewTemplate(template, startTag, endTag string) (*Template, error) ***REMOVED***
	var t Template
	err := t.Reset(template, startTag, endTag)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &t, nil
***REMOVED***

// TagFunc can be used as a substitution value in the map passed to Execute*.
// Execute* functions pass tag (placeholder) name in 'tag' argument.
//
// TagFunc must be safe to call from concurrently running goroutines.
//
// TagFunc must write contents to w and return the number of bytes written.
type TagFunc func(w io.Writer, tag string) (int, error)

// Reset resets the template t to new one defined by
// template, startTag and endTag.
//
// Reset allows Template object re-use.
//
// Reset may be called only if no other goroutines call t methods at the moment.
func (t *Template) Reset(template, startTag, endTag string) error ***REMOVED***
	// Keep these vars in t, so GC won't collect them and won't break
	// vars derived via unsafe*
	t.template = template
	t.startTag = startTag
	t.endTag = endTag
	t.texts = t.texts[:0]
	t.tags = t.tags[:0]

	if len(startTag) == 0 ***REMOVED***
		panic("startTag cannot be empty")
	***REMOVED***
	if len(endTag) == 0 ***REMOVED***
		panic("endTag cannot be empty")
	***REMOVED***

	s := unsafeString2Bytes(template)
	a := unsafeString2Bytes(startTag)
	b := unsafeString2Bytes(endTag)

	tagsCount := bytes.Count(s, a)
	if tagsCount == 0 ***REMOVED***
		return nil
	***REMOVED***

	if tagsCount+1 > cap(t.texts) ***REMOVED***
		t.texts = make([][]byte, 0, tagsCount+1)
	***REMOVED***
	if tagsCount > cap(t.tags) ***REMOVED***
		t.tags = make([]string, 0, tagsCount)
	***REMOVED***

	for ***REMOVED***
		n := bytes.Index(s, a)
		if n < 0 ***REMOVED***
			t.texts = append(t.texts, s)
			break
		***REMOVED***
		t.texts = append(t.texts, s[:n])

		s = s[n+len(a):]
		n = bytes.Index(s, b)
		if n < 0 ***REMOVED***
			return fmt.Errorf("Cannot find end tag=%q in the template=%q starting from %q", endTag, template, s)
		***REMOVED***

		t.tags = append(t.tags, unsafeBytes2String(s[:n]))
		s = s[n+len(b):]
	***REMOVED***

	return nil
***REMOVED***

// ExecuteFunc calls f on each template tag (placeholder) occurrence.
//
// Returns the number of bytes written to w.
//
// This function is optimized for frozen templates.
// Use ExecuteFunc for constantly changing templates.
func (t *Template) ExecuteFunc(w io.Writer, f TagFunc) (int64, error) ***REMOVED***
	var nn int64

	n := len(t.texts) - 1
	if n == -1 ***REMOVED***
		ni, err := w.Write(unsafeString2Bytes(t.template))
		return int64(ni), err
	***REMOVED***

	for i := 0; i < n; i++ ***REMOVED***
		ni, err := w.Write(t.texts[i])
		nn += int64(ni)
		if err != nil ***REMOVED***
			return nn, err
		***REMOVED***

		ni, err = f(w, t.tags[i])
		nn += int64(ni)
		if err != nil ***REMOVED***
			return nn, err
		***REMOVED***
	***REMOVED***
	ni, err := w.Write(t.texts[n])
	nn += int64(ni)
	return nn, err
***REMOVED***

// Execute substitutes template tags (placeholders) with the corresponding
// values from the map m and writes the result to the given writer w.
//
// Substitution map m may contain values with the following types:
//   * []byte - the fastest value type
//   * string - convenient value type
//   * TagFunc - flexible value type
//
// Returns the number of bytes written to w.
func (t *Template) Execute(w io.Writer, m map[string]interface***REMOVED******REMOVED***) (int64, error) ***REMOVED***
	return t.ExecuteFunc(w, func(w io.Writer, tag string) (int, error) ***REMOVED*** return stdTagFunc(w, tag, m) ***REMOVED***)
***REMOVED***

// ExecuteFuncString calls f on each template tag (placeholder) occurrence
// and substitutes it with the data written to TagFunc's w.
//
// Returns the resulting string.
//
// This function is optimized for frozen templates.
// Use ExecuteFuncString for constantly changing templates.
func (t *Template) ExecuteFuncString(f TagFunc) string ***REMOVED***
	bb := t.byteBufferPool.Get()
	if _, err := t.ExecuteFunc(bb, f); err != nil ***REMOVED***
		panic(fmt.Sprintf("unexpected error: %s", err))
	***REMOVED***
	s := string(bb.Bytes())
	bb.Reset()
	t.byteBufferPool.Put(bb)
	return s
***REMOVED***

// ExecuteString substitutes template tags (placeholders) with the corresponding
// values from the map m and returns the result.
//
// Substitution map m may contain values with the following types:
//   * []byte - the fastest value type
//   * string - convenient value type
//   * TagFunc - flexible value type
//
// This function is optimized for frozen templates.
// Use ExecuteString for constantly changing templates.
func (t *Template) ExecuteString(m map[string]interface***REMOVED******REMOVED***) string ***REMOVED***
	return t.ExecuteFuncString(func(w io.Writer, tag string) (int, error) ***REMOVED*** return stdTagFunc(w, tag, m) ***REMOVED***)
***REMOVED***

func stdTagFunc(w io.Writer, tag string, m map[string]interface***REMOVED******REMOVED***) (int, error) ***REMOVED***
	v := m[tag]
	if v == nil ***REMOVED***
		return 0, nil
	***REMOVED***
	switch value := v.(type) ***REMOVED***
	case []byte:
		return w.Write(value)
	case string:
		return w.Write([]byte(value))
	case TagFunc:
		return value(w, tag)
	default:
		panic(fmt.Sprintf("tag=%q contains unexpected value type=%#v. Expected []byte, string or TagFunc", tag, v))
	***REMOVED***
***REMOVED***
