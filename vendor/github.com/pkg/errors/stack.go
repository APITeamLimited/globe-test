package errors

import (
	"fmt"
	"io"
	"path"
	"runtime"
	"strconv"
	"strings"
)

// Frame represents a program counter inside a stack frame.
// For historical reasons if Frame is interpreted as a uintptr
// its value represents the program counter + 1.
type Frame uintptr

// pc returns the program counter for this frame;
// multiple frames may have the same PC value.
func (f Frame) pc() uintptr ***REMOVED*** return uintptr(f) - 1 ***REMOVED***

// file returns the full path to the file that contains the
// function for this Frame's pc.
func (f Frame) file() string ***REMOVED***
	fn := runtime.FuncForPC(f.pc())
	if fn == nil ***REMOVED***
		return "unknown"
	***REMOVED***
	file, _ := fn.FileLine(f.pc())
	return file
***REMOVED***

// line returns the line number of source code of the
// function for this Frame's pc.
func (f Frame) line() int ***REMOVED***
	fn := runtime.FuncForPC(f.pc())
	if fn == nil ***REMOVED***
		return 0
	***REMOVED***
	_, line := fn.FileLine(f.pc())
	return line
***REMOVED***

// name returns the name of this function, if known.
func (f Frame) name() string ***REMOVED***
	fn := runtime.FuncForPC(f.pc())
	if fn == nil ***REMOVED***
		return "unknown"
	***REMOVED***
	return fn.Name()
***REMOVED***

// Format formats the frame according to the fmt.Formatter interface.
//
//    %s    source file
//    %d    source line
//    %n    function name
//    %v    equivalent to %s:%d
//
// Format accepts flags that alter the printing of some verbs, as follows:
//
//    %+s   function name and path of source file relative to the compile time
//          GOPATH separated by \n\t (<funcname>\n\t<path>)
//    %+v   equivalent to %+s:%d
func (f Frame) Format(s fmt.State, verb rune) ***REMOVED***
	switch verb ***REMOVED***
	case 's':
		switch ***REMOVED***
		case s.Flag('+'):
			io.WriteString(s, f.name())
			io.WriteString(s, "\n\t")
			io.WriteString(s, f.file())
		default:
			io.WriteString(s, path.Base(f.file()))
		***REMOVED***
	case 'd':
		io.WriteString(s, strconv.Itoa(f.line()))
	case 'n':
		io.WriteString(s, funcname(f.name()))
	case 'v':
		f.Format(s, 's')
		io.WriteString(s, ":")
		f.Format(s, 'd')
	***REMOVED***
***REMOVED***

// MarshalText formats a stacktrace Frame as a text string. The output is the
// same as that of fmt.Sprintf("%+v", f), but without newlines or tabs.
func (f Frame) MarshalText() ([]byte, error) ***REMOVED***
	name := f.name()
	if name == "unknown" ***REMOVED***
		return []byte(name), nil
	***REMOVED***
	return []byte(fmt.Sprintf("%s %s:%d", name, f.file(), f.line())), nil
***REMOVED***

// StackTrace is stack of Frames from innermost (newest) to outermost (oldest).
type StackTrace []Frame

// Format formats the stack of Frames according to the fmt.Formatter interface.
//
//    %s	lists source files for each Frame in the stack
//    %v	lists the source file and line number for each Frame in the stack
//
// Format accepts flags that alter the printing of some verbs, as follows:
//
//    %+v   Prints filename, function, and line number for each Frame in the stack.
func (st StackTrace) Format(s fmt.State, verb rune) ***REMOVED***
	switch verb ***REMOVED***
	case 'v':
		switch ***REMOVED***
		case s.Flag('+'):
			for _, f := range st ***REMOVED***
				io.WriteString(s, "\n")
				f.Format(s, verb)
			***REMOVED***
		case s.Flag('#'):
			fmt.Fprintf(s, "%#v", []Frame(st))
		default:
			st.formatSlice(s, verb)
		***REMOVED***
	case 's':
		st.formatSlice(s, verb)
	***REMOVED***
***REMOVED***

// formatSlice will format this StackTrace into the given buffer as a slice of
// Frame, only valid when called with '%s' or '%v'.
func (st StackTrace) formatSlice(s fmt.State, verb rune) ***REMOVED***
	io.WriteString(s, "[")
	for i, f := range st ***REMOVED***
		if i > 0 ***REMOVED***
			io.WriteString(s, " ")
		***REMOVED***
		f.Format(s, verb)
	***REMOVED***
	io.WriteString(s, "]")
***REMOVED***

// stack represents a stack of program counters.
type stack []uintptr

func (s *stack) Format(st fmt.State, verb rune) ***REMOVED***
	switch verb ***REMOVED***
	case 'v':
		switch ***REMOVED***
		case st.Flag('+'):
			for _, pc := range *s ***REMOVED***
				f := Frame(pc)
				fmt.Fprintf(st, "\n%+v", f)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func (s *stack) StackTrace() StackTrace ***REMOVED***
	f := make([]Frame, len(*s))
	for i := 0; i < len(f); i++ ***REMOVED***
		f[i] = Frame((*s)[i])
	***REMOVED***
	return f
***REMOVED***

func callers() *stack ***REMOVED***
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(3, pcs[:])
	var st stack = pcs[0:n]
	return &st
***REMOVED***

// funcname removes the path prefix component of a function's name reported by func.Name().
func funcname(name string) string ***REMOVED***
	i := strings.LastIndex(name, "/")
	name = name[i+1:]
	i = strings.Index(name, ".")
	return name[i+1:]
***REMOVED***
