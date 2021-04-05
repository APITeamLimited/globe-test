package color

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/mattn/go-colorable"
	"github.com/mattn/go-isatty"
)

var (
	// NoColor defines if the output is colorized or not. It's dynamically set to
	// false or true based on the stdout's file descriptor referring to a terminal
	// or not. This is a global option and affects all colors. For more control
	// over each color block use the methods DisableColor() individually.
	NoColor = os.Getenv("TERM") == "dumb" ||
		(!isatty.IsTerminal(os.Stdout.Fd()) && !isatty.IsCygwinTerminal(os.Stdout.Fd()))

	// Output defines the standard output of the print functions. By default
	// os.Stdout is used.
	Output = colorable.NewColorableStdout()

	// Error defines a color supporting writer for os.Stderr.
	Error = colorable.NewColorableStderr()

	// colorsCache is used to reduce the count of created Color objects and
	// allows to reuse already created objects with required Attribute.
	colorsCache   = make(map[Attribute]*Color)
	colorsCacheMu sync.Mutex // protects colorsCache
)

// Color defines a custom color object which is defined by SGR parameters.
type Color struct ***REMOVED***
	params  []Attribute
	noColor *bool
***REMOVED***

// Attribute defines a single SGR Code
type Attribute int

const escape = "\x1b"

// Base attributes
const (
	Reset Attribute = iota
	Bold
	Faint
	Italic
	Underline
	BlinkSlow
	BlinkRapid
	ReverseVideo
	Concealed
	CrossedOut
)

// Foreground text colors
const (
	FgBlack Attribute = iota + 30
	FgRed
	FgGreen
	FgYellow
	FgBlue
	FgMagenta
	FgCyan
	FgWhite
)

// Foreground Hi-Intensity text colors
const (
	FgHiBlack Attribute = iota + 90
	FgHiRed
	FgHiGreen
	FgHiYellow
	FgHiBlue
	FgHiMagenta
	FgHiCyan
	FgHiWhite
)

// Background text colors
const (
	BgBlack Attribute = iota + 40
	BgRed
	BgGreen
	BgYellow
	BgBlue
	BgMagenta
	BgCyan
	BgWhite
)

// Background Hi-Intensity text colors
const (
	BgHiBlack Attribute = iota + 100
	BgHiRed
	BgHiGreen
	BgHiYellow
	BgHiBlue
	BgHiMagenta
	BgHiCyan
	BgHiWhite
)

// New returns a newly created color object.
func New(value ...Attribute) *Color ***REMOVED***
	c := &Color***REMOVED***params: make([]Attribute, 0)***REMOVED***
	c.Add(value...)
	return c
***REMOVED***

// Set sets the given parameters immediately. It will change the color of
// output with the given SGR parameters until color.Unset() is called.
func Set(p ...Attribute) *Color ***REMOVED***
	c := New(p...)
	c.Set()
	return c
***REMOVED***

// Unset resets all escape attributes and clears the output. Usually should
// be called after Set().
func Unset() ***REMOVED***
	if NoColor ***REMOVED***
		return
	***REMOVED***

	fmt.Fprintf(Output, "%s[%dm", escape, Reset)
***REMOVED***

// Set sets the SGR sequence.
func (c *Color) Set() *Color ***REMOVED***
	if c.isNoColorSet() ***REMOVED***
		return c
	***REMOVED***

	fmt.Fprintf(Output, c.format())
	return c
***REMOVED***

func (c *Color) unset() ***REMOVED***
	if c.isNoColorSet() ***REMOVED***
		return
	***REMOVED***

	Unset()
***REMOVED***

func (c *Color) setWriter(w io.Writer) *Color ***REMOVED***
	if c.isNoColorSet() ***REMOVED***
		return c
	***REMOVED***

	fmt.Fprintf(w, c.format())
	return c
***REMOVED***

func (c *Color) unsetWriter(w io.Writer) ***REMOVED***
	if c.isNoColorSet() ***REMOVED***
		return
	***REMOVED***

	if NoColor ***REMOVED***
		return
	***REMOVED***

	fmt.Fprintf(w, "%s[%dm", escape, Reset)
***REMOVED***

// Add is used to chain SGR parameters. Use as many as parameters to combine
// and create custom color objects. Example: Add(color.FgRed, color.Underline).
func (c *Color) Add(value ...Attribute) *Color ***REMOVED***
	c.params = append(c.params, value...)
	return c
***REMOVED***

func (c *Color) prepend(value Attribute) ***REMOVED***
	c.params = append(c.params, 0)
	copy(c.params[1:], c.params[0:])
	c.params[0] = value
***REMOVED***

// Fprint formats using the default formats for its operands and writes to w.
// Spaces are added between operands when neither is a string.
// It returns the number of bytes written and any write error encountered.
// On Windows, users should wrap w with colorable.NewColorable() if w is of
// type *os.File.
func (c *Color) Fprint(w io.Writer, a ...interface***REMOVED******REMOVED***) (n int, err error) ***REMOVED***
	c.setWriter(w)
	defer c.unsetWriter(w)

	return fmt.Fprint(w, a...)
***REMOVED***

// Print formats using the default formats for its operands and writes to
// standard output. Spaces are added between operands when neither is a
// string. It returns the number of bytes written and any write error
// encountered. This is the standard fmt.Print() method wrapped with the given
// color.
func (c *Color) Print(a ...interface***REMOVED******REMOVED***) (n int, err error) ***REMOVED***
	c.Set()
	defer c.unset()

	return fmt.Fprint(Output, a...)
***REMOVED***

// Fprintf formats according to a format specifier and writes to w.
// It returns the number of bytes written and any write error encountered.
// On Windows, users should wrap w with colorable.NewColorable() if w is of
// type *os.File.
func (c *Color) Fprintf(w io.Writer, format string, a ...interface***REMOVED******REMOVED***) (n int, err error) ***REMOVED***
	c.setWriter(w)
	defer c.unsetWriter(w)

	return fmt.Fprintf(w, format, a...)
***REMOVED***

// Printf formats according to a format specifier and writes to standard output.
// It returns the number of bytes written and any write error encountered.
// This is the standard fmt.Printf() method wrapped with the given color.
func (c *Color) Printf(format string, a ...interface***REMOVED******REMOVED***) (n int, err error) ***REMOVED***
	c.Set()
	defer c.unset()

	return fmt.Fprintf(Output, format, a...)
***REMOVED***

// Fprintln formats using the default formats for its operands and writes to w.
// Spaces are always added between operands and a newline is appended.
// On Windows, users should wrap w with colorable.NewColorable() if w is of
// type *os.File.
func (c *Color) Fprintln(w io.Writer, a ...interface***REMOVED******REMOVED***) (n int, err error) ***REMOVED***
	c.setWriter(w)
	defer c.unsetWriter(w)

	return fmt.Fprintln(w, a...)
***REMOVED***

// Println formats using the default formats for its operands and writes to
// standard output. Spaces are always added between operands and a newline is
// appended. It returns the number of bytes written and any write error
// encountered. This is the standard fmt.Print() method wrapped with the given
// color.
func (c *Color) Println(a ...interface***REMOVED******REMOVED***) (n int, err error) ***REMOVED***
	c.Set()
	defer c.unset()

	return fmt.Fprintln(Output, a...)
***REMOVED***

// Sprint is just like Print, but returns a string instead of printing it.
func (c *Color) Sprint(a ...interface***REMOVED******REMOVED***) string ***REMOVED***
	return c.wrap(fmt.Sprint(a...))
***REMOVED***

// Sprintln is just like Println, but returns a string instead of printing it.
func (c *Color) Sprintln(a ...interface***REMOVED******REMOVED***) string ***REMOVED***
	return c.wrap(fmt.Sprintln(a...))
***REMOVED***

// Sprintf is just like Printf, but returns a string instead of printing it.
func (c *Color) Sprintf(format string, a ...interface***REMOVED******REMOVED***) string ***REMOVED***
	return c.wrap(fmt.Sprintf(format, a...))
***REMOVED***

// FprintFunc returns a new function that prints the passed arguments as
// colorized with color.Fprint().
func (c *Color) FprintFunc() func(w io.Writer, a ...interface***REMOVED******REMOVED***) ***REMOVED***
	return func(w io.Writer, a ...interface***REMOVED******REMOVED***) ***REMOVED***
		c.Fprint(w, a...)
	***REMOVED***
***REMOVED***

// PrintFunc returns a new function that prints the passed arguments as
// colorized with color.Print().
func (c *Color) PrintFunc() func(a ...interface***REMOVED******REMOVED***) ***REMOVED***
	return func(a ...interface***REMOVED******REMOVED***) ***REMOVED***
		c.Print(a...)
	***REMOVED***
***REMOVED***

// FprintfFunc returns a new function that prints the passed arguments as
// colorized with color.Fprintf().
func (c *Color) FprintfFunc() func(w io.Writer, format string, a ...interface***REMOVED******REMOVED***) ***REMOVED***
	return func(w io.Writer, format string, a ...interface***REMOVED******REMOVED***) ***REMOVED***
		c.Fprintf(w, format, a...)
	***REMOVED***
***REMOVED***

// PrintfFunc returns a new function that prints the passed arguments as
// colorized with color.Printf().
func (c *Color) PrintfFunc() func(format string, a ...interface***REMOVED******REMOVED***) ***REMOVED***
	return func(format string, a ...interface***REMOVED******REMOVED***) ***REMOVED***
		c.Printf(format, a...)
	***REMOVED***
***REMOVED***

// FprintlnFunc returns a new function that prints the passed arguments as
// colorized with color.Fprintln().
func (c *Color) FprintlnFunc() func(w io.Writer, a ...interface***REMOVED******REMOVED***) ***REMOVED***
	return func(w io.Writer, a ...interface***REMOVED******REMOVED***) ***REMOVED***
		c.Fprintln(w, a...)
	***REMOVED***
***REMOVED***

// PrintlnFunc returns a new function that prints the passed arguments as
// colorized with color.Println().
func (c *Color) PrintlnFunc() func(a ...interface***REMOVED******REMOVED***) ***REMOVED***
	return func(a ...interface***REMOVED******REMOVED***) ***REMOVED***
		c.Println(a...)
	***REMOVED***
***REMOVED***

// SprintFunc returns a new function that returns colorized strings for the
// given arguments with fmt.Sprint(). Useful to put into or mix into other
// string. Windows users should use this in conjunction with color.Output, example:
//
//	put := New(FgYellow).SprintFunc()
//	fmt.Fprintf(color.Output, "This is a %s", put("warning"))
func (c *Color) SprintFunc() func(a ...interface***REMOVED******REMOVED***) string ***REMOVED***
	return func(a ...interface***REMOVED******REMOVED***) string ***REMOVED***
		return c.wrap(fmt.Sprint(a...))
	***REMOVED***
***REMOVED***

// SprintfFunc returns a new function that returns colorized strings for the
// given arguments with fmt.Sprintf(). Useful to put into or mix into other
// string. Windows users should use this in conjunction with color.Output.
func (c *Color) SprintfFunc() func(format string, a ...interface***REMOVED******REMOVED***) string ***REMOVED***
	return func(format string, a ...interface***REMOVED******REMOVED***) string ***REMOVED***
		return c.wrap(fmt.Sprintf(format, a...))
	***REMOVED***
***REMOVED***

// SprintlnFunc returns a new function that returns colorized strings for the
// given arguments with fmt.Sprintln(). Useful to put into or mix into other
// string. Windows users should use this in conjunction with color.Output.
func (c *Color) SprintlnFunc() func(a ...interface***REMOVED******REMOVED***) string ***REMOVED***
	return func(a ...interface***REMOVED******REMOVED***) string ***REMOVED***
		return c.wrap(fmt.Sprintln(a...))
	***REMOVED***
***REMOVED***

// sequence returns a formatted SGR sequence to be plugged into a "\x1b[...m"
// an example output might be: "1;36" -> bold cyan
func (c *Color) sequence() string ***REMOVED***
	format := make([]string, len(c.params))
	for i, v := range c.params ***REMOVED***
		format[i] = strconv.Itoa(int(v))
	***REMOVED***

	return strings.Join(format, ";")
***REMOVED***

// wrap wraps the s string with the colors attributes. The string is ready to
// be printed.
func (c *Color) wrap(s string) string ***REMOVED***
	if c.isNoColorSet() ***REMOVED***
		return s
	***REMOVED***

	return c.format() + s + c.unformat()
***REMOVED***

func (c *Color) format() string ***REMOVED***
	return fmt.Sprintf("%s[%sm", escape, c.sequence())
***REMOVED***

func (c *Color) unformat() string ***REMOVED***
	return fmt.Sprintf("%s[%dm", escape, Reset)
***REMOVED***

// DisableColor disables the color output. Useful to not change any existing
// code and still being able to output. Can be used for flags like
// "--no-color". To enable back use EnableColor() method.
func (c *Color) DisableColor() ***REMOVED***
	c.noColor = boolPtr(true)
***REMOVED***

// EnableColor enables the color output. Use it in conjunction with
// DisableColor(). Otherwise this method has no side effects.
func (c *Color) EnableColor() ***REMOVED***
	c.noColor = boolPtr(false)
***REMOVED***

func (c *Color) isNoColorSet() bool ***REMOVED***
	// check first if we have user setted action
	if c.noColor != nil ***REMOVED***
		return *c.noColor
	***REMOVED***

	// if not return the global option, which is disabled by default
	return NoColor
***REMOVED***

// Equals returns a boolean value indicating whether two colors are equal.
func (c *Color) Equals(c2 *Color) bool ***REMOVED***
	if len(c.params) != len(c2.params) ***REMOVED***
		return false
	***REMOVED***

	for _, attr := range c.params ***REMOVED***
		if !c2.attrExists(attr) ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***

	return true
***REMOVED***

func (c *Color) attrExists(a Attribute) bool ***REMOVED***
	for _, attr := range c.params ***REMOVED***
		if attr == a ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***

	return false
***REMOVED***

func boolPtr(v bool) *bool ***REMOVED***
	return &v
***REMOVED***

func getCachedColor(p Attribute) *Color ***REMOVED***
	colorsCacheMu.Lock()
	defer colorsCacheMu.Unlock()

	c, ok := colorsCache[p]
	if !ok ***REMOVED***
		c = New(p)
		colorsCache[p] = c
	***REMOVED***

	return c
***REMOVED***

func colorPrint(format string, p Attribute, a ...interface***REMOVED******REMOVED***) ***REMOVED***
	c := getCachedColor(p)

	if !strings.HasSuffix(format, "\n") ***REMOVED***
		format += "\n"
	***REMOVED***

	if len(a) == 0 ***REMOVED***
		c.Print(format)
	***REMOVED*** else ***REMOVED***
		c.Printf(format, a...)
	***REMOVED***
***REMOVED***

func colorString(format string, p Attribute, a ...interface***REMOVED******REMOVED***) string ***REMOVED***
	c := getCachedColor(p)

	if len(a) == 0 ***REMOVED***
		return c.SprintFunc()(format)
	***REMOVED***

	return c.SprintfFunc()(format, a...)
***REMOVED***

// Black is a convenient helper function to print with black foreground. A
// newline is appended to format by default.
func Black(format string, a ...interface***REMOVED******REMOVED***) ***REMOVED*** colorPrint(format, FgBlack, a...) ***REMOVED***

// Red is a convenient helper function to print with red foreground. A
// newline is appended to format by default.
func Red(format string, a ...interface***REMOVED******REMOVED***) ***REMOVED*** colorPrint(format, FgRed, a...) ***REMOVED***

// Green is a convenient helper function to print with green foreground. A
// newline is appended to format by default.
func Green(format string, a ...interface***REMOVED******REMOVED***) ***REMOVED*** colorPrint(format, FgGreen, a...) ***REMOVED***

// Yellow is a convenient helper function to print with yellow foreground.
// A newline is appended to format by default.
func Yellow(format string, a ...interface***REMOVED******REMOVED***) ***REMOVED*** colorPrint(format, FgYellow, a...) ***REMOVED***

// Blue is a convenient helper function to print with blue foreground. A
// newline is appended to format by default.
func Blue(format string, a ...interface***REMOVED******REMOVED***) ***REMOVED*** colorPrint(format, FgBlue, a...) ***REMOVED***

// Magenta is a convenient helper function to print with magenta foreground.
// A newline is appended to format by default.
func Magenta(format string, a ...interface***REMOVED******REMOVED***) ***REMOVED*** colorPrint(format, FgMagenta, a...) ***REMOVED***

// Cyan is a convenient helper function to print with cyan foreground. A
// newline is appended to format by default.
func Cyan(format string, a ...interface***REMOVED******REMOVED***) ***REMOVED*** colorPrint(format, FgCyan, a...) ***REMOVED***

// White is a convenient helper function to print with white foreground. A
// newline is appended to format by default.
func White(format string, a ...interface***REMOVED******REMOVED***) ***REMOVED*** colorPrint(format, FgWhite, a...) ***REMOVED***

// BlackString is a convenient helper function to return a string with black
// foreground.
func BlackString(format string, a ...interface***REMOVED******REMOVED***) string ***REMOVED*** return colorString(format, FgBlack, a...) ***REMOVED***

// RedString is a convenient helper function to return a string with red
// foreground.
func RedString(format string, a ...interface***REMOVED******REMOVED***) string ***REMOVED*** return colorString(format, FgRed, a...) ***REMOVED***

// GreenString is a convenient helper function to return a string with green
// foreground.
func GreenString(format string, a ...interface***REMOVED******REMOVED***) string ***REMOVED*** return colorString(format, FgGreen, a...) ***REMOVED***

// YellowString is a convenient helper function to return a string with yellow
// foreground.
func YellowString(format string, a ...interface***REMOVED******REMOVED***) string ***REMOVED*** return colorString(format, FgYellow, a...) ***REMOVED***

// BlueString is a convenient helper function to return a string with blue
// foreground.
func BlueString(format string, a ...interface***REMOVED******REMOVED***) string ***REMOVED*** return colorString(format, FgBlue, a...) ***REMOVED***

// MagentaString is a convenient helper function to return a string with magenta
// foreground.
func MagentaString(format string, a ...interface***REMOVED******REMOVED***) string ***REMOVED***
	return colorString(format, FgMagenta, a...)
***REMOVED***

// CyanString is a convenient helper function to return a string with cyan
// foreground.
func CyanString(format string, a ...interface***REMOVED******REMOVED***) string ***REMOVED*** return colorString(format, FgCyan, a...) ***REMOVED***

// WhiteString is a convenient helper function to return a string with white
// foreground.
func WhiteString(format string, a ...interface***REMOVED******REMOVED***) string ***REMOVED*** return colorString(format, FgWhite, a...) ***REMOVED***

// HiBlack is a convenient helper function to print with hi-intensity black foreground. A
// newline is appended to format by default.
func HiBlack(format string, a ...interface***REMOVED******REMOVED***) ***REMOVED*** colorPrint(format, FgHiBlack, a...) ***REMOVED***

// HiRed is a convenient helper function to print with hi-intensity red foreground. A
// newline is appended to format by default.
func HiRed(format string, a ...interface***REMOVED******REMOVED***) ***REMOVED*** colorPrint(format, FgHiRed, a...) ***REMOVED***

// HiGreen is a convenient helper function to print with hi-intensity green foreground. A
// newline is appended to format by default.
func HiGreen(format string, a ...interface***REMOVED******REMOVED***) ***REMOVED*** colorPrint(format, FgHiGreen, a...) ***REMOVED***

// HiYellow is a convenient helper function to print with hi-intensity yellow foreground.
// A newline is appended to format by default.
func HiYellow(format string, a ...interface***REMOVED******REMOVED***) ***REMOVED*** colorPrint(format, FgHiYellow, a...) ***REMOVED***

// HiBlue is a convenient helper function to print with hi-intensity blue foreground. A
// newline is appended to format by default.
func HiBlue(format string, a ...interface***REMOVED******REMOVED***) ***REMOVED*** colorPrint(format, FgHiBlue, a...) ***REMOVED***

// HiMagenta is a convenient helper function to print with hi-intensity magenta foreground.
// A newline is appended to format by default.
func HiMagenta(format string, a ...interface***REMOVED******REMOVED***) ***REMOVED*** colorPrint(format, FgHiMagenta, a...) ***REMOVED***

// HiCyan is a convenient helper function to print with hi-intensity cyan foreground. A
// newline is appended to format by default.
func HiCyan(format string, a ...interface***REMOVED******REMOVED***) ***REMOVED*** colorPrint(format, FgHiCyan, a...) ***REMOVED***

// HiWhite is a convenient helper function to print with hi-intensity white foreground. A
// newline is appended to format by default.
func HiWhite(format string, a ...interface***REMOVED******REMOVED***) ***REMOVED*** colorPrint(format, FgHiWhite, a...) ***REMOVED***

// HiBlackString is a convenient helper function to return a string with hi-intensity black
// foreground.
func HiBlackString(format string, a ...interface***REMOVED******REMOVED***) string ***REMOVED***
	return colorString(format, FgHiBlack, a...)
***REMOVED***

// HiRedString is a convenient helper function to return a string with hi-intensity red
// foreground.
func HiRedString(format string, a ...interface***REMOVED******REMOVED***) string ***REMOVED*** return colorString(format, FgHiRed, a...) ***REMOVED***

// HiGreenString is a convenient helper function to return a string with hi-intensity green
// foreground.
func HiGreenString(format string, a ...interface***REMOVED******REMOVED***) string ***REMOVED***
	return colorString(format, FgHiGreen, a...)
***REMOVED***

// HiYellowString is a convenient helper function to return a string with hi-intensity yellow
// foreground.
func HiYellowString(format string, a ...interface***REMOVED******REMOVED***) string ***REMOVED***
	return colorString(format, FgHiYellow, a...)
***REMOVED***

// HiBlueString is a convenient helper function to return a string with hi-intensity blue
// foreground.
func HiBlueString(format string, a ...interface***REMOVED******REMOVED***) string ***REMOVED*** return colorString(format, FgHiBlue, a...) ***REMOVED***

// HiMagentaString is a convenient helper function to return a string with hi-intensity magenta
// foreground.
func HiMagentaString(format string, a ...interface***REMOVED******REMOVED***) string ***REMOVED***
	return colorString(format, FgHiMagenta, a...)
***REMOVED***

// HiCyanString is a convenient helper function to return a string with hi-intensity cyan
// foreground.
func HiCyanString(format string, a ...interface***REMOVED******REMOVED***) string ***REMOVED*** return colorString(format, FgHiCyan, a...) ***REMOVED***

// HiWhiteString is a convenient helper function to return a string with hi-intensity white
// foreground.
func HiWhiteString(format string, a ...interface***REMOVED******REMOVED***) string ***REMOVED***
	return colorString(format, FgHiWhite, a...)
***REMOVED***
