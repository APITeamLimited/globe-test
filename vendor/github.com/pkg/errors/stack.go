package errors

import (
	"fmt"
	"io"
	"path"
	"runtime"
	"strings"
)

// Frame represents a program counter inside a stack frame.
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
			pc := f.pc()
			fn := runtime.FuncForPC(pc)
			if fn == nil ***REMOVED***
				io.WriteString(s, "unknown")
			***REMOVED*** else ***REMOVED***
				file, _ := fn.FileLine(pc)
				fmt.Fprintf(s, "%s\n\t%s", fn.Name(), file)
			***REMOVED***
		default:
			io.WriteString(s, path.Base(f.file()))
		***REMOVED***
	case 'd':
		fmt.Fprintf(s, "%d", f.line())
	case 'n':
		name := runtime.FuncForPC(f.pc()).Name()
		io.WriteString(s, funcname(name))
	case 'v':
		f.Format(s, 's')
		io.WriteString(s, ":")
		f.Format(s, 'd')
	***REMOVED***
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
				fmt.Fprintf(s, "\n%+v", f)
			***REMOVED***
		case s.Flag('#'):
			fmt.Fprintf(s, "%#v", []Frame(st))
		default:
			fmt.Fprintf(s, "%v", []Frame(st))
		***REMOVED***
	case 's':
		fmt.Fprintf(s, "%s", []Frame(st))
	***REMOVED***
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
