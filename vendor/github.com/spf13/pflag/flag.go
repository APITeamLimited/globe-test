// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package pflag is a drop-in replacement for Go's flag package, implementing
POSIX/GNU-style --flags.

pflag is compatible with the GNU extensions to the POSIX recommendations
for command-line options. See
http://www.gnu.org/software/libc/manual/html_node/Argument-Syntax.html

Usage:

pflag is a drop-in replacement of Go's native flag package. If you import
pflag under the name "flag" then all code should continue to function
with no changes.

	import flag "github.com/spf13/pflag"

There is one exception to this: if you directly instantiate the Flag struct
there is one more field "Shorthand" that you will need to set.
Most code never instantiates this struct directly, and instead uses
functions such as String(), BoolVar(), and Var(), and is therefore
unaffected.

Define flags using flag.String(), Bool(), Int(), etc.

This declares an integer flag, -flagname, stored in the pointer ip, with type *int.
	var ip = flag.Int("flagname", 1234, "help message for flagname")
If you like, you can bind the flag to a variable using the Var() functions.
	var flagvar int
	func init() ***REMOVED***
		flag.IntVar(&flagvar, "flagname", 1234, "help message for flagname")
	***REMOVED***
Or you can create custom flags that satisfy the Value interface (with
pointer receivers) and couple them to flag parsing by
	flag.Var(&flagVal, "name", "help message for flagname")
For such flags, the default value is just the initial value of the variable.

After all flags are defined, call
	flag.Parse()
to parse the command line into the defined flags.

Flags may then be used directly. If you're using the flags themselves,
they are all pointers; if you bind to variables, they're values.
	fmt.Println("ip has value ", *ip)
	fmt.Println("flagvar has value ", flagvar)

After parsing, the arguments after the flag are available as the
slice flag.Args() or individually as flag.Arg(i).
The arguments are indexed from 0 through flag.NArg()-1.

The pflag package also defines some new functions that are not in flag,
that give one-letter shorthands for flags. You can use these by appending
'P' to the name of any function that defines a flag.
	var ip = flag.IntP("flagname", "f", 1234, "help message")
	var flagvar bool
	func init() ***REMOVED***
		flag.BoolVarP("boolname", "b", true, "help message")
	***REMOVED***
	flag.VarP(&flagVar, "varname", "v", 1234, "help message")
Shorthand letters can be used with single dashes on the command line.
Boolean shorthand flags can be combined with other shorthand flags.

Command line flag syntax:
	--flag    // boolean flags only
	--flag=x

Unlike the flag package, a single dash before an option means something
different than a double dash. Single dashes signify a series of shorthand
letters for flags. All but the last shorthand letter must be boolean flags.
	// boolean flags
	-f
	-abc
	// non-boolean flags
	-n 1234
	-Ifile
	// mixed
	-abcs "hello"
	-abcn1234

Flag parsing stops after the terminator "--". Unlike the flag package,
flags can be interspersed with arguments anywhere on the command line
before this terminator.

Integer flags accept 1234, 0664, 0x1234 and may be negative.
Boolean flags (in their long form) accept 1, 0, t, f, true, false,
TRUE, FALSE, True, False.
Duration flags accept any input valid for time.ParseDuration.

The default set of command-line flags is controlled by
top-level functions.  The FlagSet type allows one to define
independent sets of flags, such as to implement subcommands
in a command-line interface. The methods of FlagSet are
analogous to the top-level functions for the command-line
flag set.
*/
package pflag

import (
	"bytes"
	"errors"
	goflag "flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

// ErrHelp is the error returned if the flag -help is invoked but no such flag is defined.
var ErrHelp = errors.New("pflag: help requested")

// ErrorHandling defines how to handle flag parsing errors.
type ErrorHandling int

const (
	// ContinueOnError will return an err from Parse() if an error is found
	ContinueOnError ErrorHandling = iota
	// ExitOnError will call os.Exit(2) if an error is found when parsing
	ExitOnError
	// PanicOnError will panic() if an error is found when parsing flags
	PanicOnError
)

// ParseErrorsWhitelist defines the parsing errors that can be ignored
type ParseErrorsWhitelist struct ***REMOVED***
	// UnknownFlags will ignore unknown flags errors and continue parsing rest of the flags
	UnknownFlags bool
***REMOVED***

// NormalizedName is a flag name that has been normalized according to rules
// for the FlagSet (e.g. making '-' and '_' equivalent).
type NormalizedName string

// A FlagSet represents a set of defined flags.
type FlagSet struct ***REMOVED***
	// Usage is the function called when an error occurs while parsing flags.
	// The field is a function (not a method) that may be changed to point to
	// a custom error handler.
	Usage func()

	// SortFlags is used to indicate, if user wants to have sorted flags in
	// help/usage messages.
	SortFlags bool

	// ParseErrorsWhitelist is used to configure a whitelist of errors
	ParseErrorsWhitelist ParseErrorsWhitelist

	name              string
	parsed            bool
	actual            map[NormalizedName]*Flag
	orderedActual     []*Flag
	sortedActual      []*Flag
	formal            map[NormalizedName]*Flag
	orderedFormal     []*Flag
	sortedFormal      []*Flag
	shorthands        map[byte]*Flag
	args              []string // arguments after flags
	argsLenAtDash     int      // len(args) when a '--' was located when parsing, or -1 if no --
	errorHandling     ErrorHandling
	output            io.Writer // nil means stderr; use out() accessor
	interspersed      bool      // allow interspersed option/non-option args
	normalizeNameFunc func(f *FlagSet, name string) NormalizedName

	addedGoFlagSets []*goflag.FlagSet
***REMOVED***

// A Flag represents the state of a flag.
type Flag struct ***REMOVED***
	Name                string              // name as it appears on command line
	Shorthand           string              // one-letter abbreviated flag
	Usage               string              // help message
	Value               Value               // value as set
	DefValue            string              // default value (as text); for usage message
	Changed             bool                // If the user set the value (or if left to default)
	NoOptDefVal         string              // default value (as text); if the flag is on the command line without any options
	Deprecated          string              // If this flag is deprecated, this string is the new or now thing to use
	Hidden              bool                // used by cobra.Command to allow flags to be hidden from help/usage text
	ShorthandDeprecated string              // If the shorthand of this flag is deprecated, this string is the new or now thing to use
	Annotations         map[string][]string // used by cobra.Command bash autocomple code
***REMOVED***

// Value is the interface to the dynamic value stored in a flag.
// (The default value is represented as a string.)
type Value interface ***REMOVED***
	String() string
	Set(string) error
	Type() string
***REMOVED***

// sortFlags returns the flags as a slice in lexicographical sorted order.
func sortFlags(flags map[NormalizedName]*Flag) []*Flag ***REMOVED***
	list := make(sort.StringSlice, len(flags))
	i := 0
	for k := range flags ***REMOVED***
		list[i] = string(k)
		i++
	***REMOVED***
	list.Sort()
	result := make([]*Flag, len(list))
	for i, name := range list ***REMOVED***
		result[i] = flags[NormalizedName(name)]
	***REMOVED***
	return result
***REMOVED***

// SetNormalizeFunc allows you to add a function which can translate flag names.
// Flags added to the FlagSet will be translated and then when anything tries to
// look up the flag that will also be translated. So it would be possible to create
// a flag named "getURL" and have it translated to "geturl".  A user could then pass
// "--getUrl" which may also be translated to "geturl" and everything will work.
func (f *FlagSet) SetNormalizeFunc(n func(f *FlagSet, name string) NormalizedName) ***REMOVED***
	f.normalizeNameFunc = n
	f.sortedFormal = f.sortedFormal[:0]
	for fname, flag := range f.formal ***REMOVED***
		nname := f.normalizeFlagName(flag.Name)
		if fname == nname ***REMOVED***
			continue
		***REMOVED***
		flag.Name = string(nname)
		delete(f.formal, fname)
		f.formal[nname] = flag
		if _, set := f.actual[fname]; set ***REMOVED***
			delete(f.actual, fname)
			f.actual[nname] = flag
		***REMOVED***
	***REMOVED***
***REMOVED***

// GetNormalizeFunc returns the previously set NormalizeFunc of a function which
// does no translation, if not set previously.
func (f *FlagSet) GetNormalizeFunc() func(f *FlagSet, name string) NormalizedName ***REMOVED***
	if f.normalizeNameFunc != nil ***REMOVED***
		return f.normalizeNameFunc
	***REMOVED***
	return func(f *FlagSet, name string) NormalizedName ***REMOVED*** return NormalizedName(name) ***REMOVED***
***REMOVED***

func (f *FlagSet) normalizeFlagName(name string) NormalizedName ***REMOVED***
	n := f.GetNormalizeFunc()
	return n(f, name)
***REMOVED***

func (f *FlagSet) out() io.Writer ***REMOVED***
	if f.output == nil ***REMOVED***
		return os.Stderr
	***REMOVED***
	return f.output
***REMOVED***

// SetOutput sets the destination for usage and error messages.
// If output is nil, os.Stderr is used.
func (f *FlagSet) SetOutput(output io.Writer) ***REMOVED***
	f.output = output
***REMOVED***

// VisitAll visits the flags in lexicographical order or
// in primordial order if f.SortFlags is false, calling fn for each.
// It visits all flags, even those not set.
func (f *FlagSet) VisitAll(fn func(*Flag)) ***REMOVED***
	if len(f.formal) == 0 ***REMOVED***
		return
	***REMOVED***

	var flags []*Flag
	if f.SortFlags ***REMOVED***
		if len(f.formal) != len(f.sortedFormal) ***REMOVED***
			f.sortedFormal = sortFlags(f.formal)
		***REMOVED***
		flags = f.sortedFormal
	***REMOVED*** else ***REMOVED***
		flags = f.orderedFormal
	***REMOVED***

	for _, flag := range flags ***REMOVED***
		fn(flag)
	***REMOVED***
***REMOVED***

// HasFlags returns a bool to indicate if the FlagSet has any flags defined.
func (f *FlagSet) HasFlags() bool ***REMOVED***
	return len(f.formal) > 0
***REMOVED***

// HasAvailableFlags returns a bool to indicate if the FlagSet has any flags
// that are not hidden.
func (f *FlagSet) HasAvailableFlags() bool ***REMOVED***
	for _, flag := range f.formal ***REMOVED***
		if !flag.Hidden ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

// VisitAll visits the command-line flags in lexicographical order or
// in primordial order if f.SortFlags is false, calling fn for each.
// It visits all flags, even those not set.
func VisitAll(fn func(*Flag)) ***REMOVED***
	CommandLine.VisitAll(fn)
***REMOVED***

// Visit visits the flags in lexicographical order or
// in primordial order if f.SortFlags is false, calling fn for each.
// It visits only those flags that have been set.
func (f *FlagSet) Visit(fn func(*Flag)) ***REMOVED***
	if len(f.actual) == 0 ***REMOVED***
		return
	***REMOVED***

	var flags []*Flag
	if f.SortFlags ***REMOVED***
		if len(f.actual) != len(f.sortedActual) ***REMOVED***
			f.sortedActual = sortFlags(f.actual)
		***REMOVED***
		flags = f.sortedActual
	***REMOVED*** else ***REMOVED***
		flags = f.orderedActual
	***REMOVED***

	for _, flag := range flags ***REMOVED***
		fn(flag)
	***REMOVED***
***REMOVED***

// Visit visits the command-line flags in lexicographical order or
// in primordial order if f.SortFlags is false, calling fn for each.
// It visits only those flags that have been set.
func Visit(fn func(*Flag)) ***REMOVED***
	CommandLine.Visit(fn)
***REMOVED***

// Lookup returns the Flag structure of the named flag, returning nil if none exists.
func (f *FlagSet) Lookup(name string) *Flag ***REMOVED***
	return f.lookup(f.normalizeFlagName(name))
***REMOVED***

// ShorthandLookup returns the Flag structure of the short handed flag,
// returning nil if none exists.
// It panics, if len(name) > 1.
func (f *FlagSet) ShorthandLookup(name string) *Flag ***REMOVED***
	if name == "" ***REMOVED***
		return nil
	***REMOVED***
	if len(name) > 1 ***REMOVED***
		msg := fmt.Sprintf("can not look up shorthand which is more than one ASCII character: %q", name)
		fmt.Fprintf(f.out(), msg)
		panic(msg)
	***REMOVED***
	c := name[0]
	return f.shorthands[c]
***REMOVED***

// lookup returns the Flag structure of the named flag, returning nil if none exists.
func (f *FlagSet) lookup(name NormalizedName) *Flag ***REMOVED***
	return f.formal[name]
***REMOVED***

// func to return a given type for a given flag name
func (f *FlagSet) getFlagType(name string, ftype string, convFunc func(sval string) (interface***REMOVED******REMOVED***, error)) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	flag := f.Lookup(name)
	if flag == nil ***REMOVED***
		err := fmt.Errorf("flag accessed but not defined: %s", name)
		return nil, err
	***REMOVED***

	if flag.Value.Type() != ftype ***REMOVED***
		err := fmt.Errorf("trying to get %s value of flag of type %s", ftype, flag.Value.Type())
		return nil, err
	***REMOVED***

	sval := flag.Value.String()
	result, err := convFunc(sval)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return result, nil
***REMOVED***

// ArgsLenAtDash will return the length of f.Args at the moment when a -- was
// found during arg parsing. This allows your program to know which args were
// before the -- and which came after.
func (f *FlagSet) ArgsLenAtDash() int ***REMOVED***
	return f.argsLenAtDash
***REMOVED***

// MarkDeprecated indicated that a flag is deprecated in your program. It will
// continue to function but will not show up in help or usage messages. Using
// this flag will also print the given usageMessage.
func (f *FlagSet) MarkDeprecated(name string, usageMessage string) error ***REMOVED***
	flag := f.Lookup(name)
	if flag == nil ***REMOVED***
		return fmt.Errorf("flag %q does not exist", name)
	***REMOVED***
	if usageMessage == "" ***REMOVED***
		return fmt.Errorf("deprecated message for flag %q must be set", name)
	***REMOVED***
	flag.Deprecated = usageMessage
	flag.Hidden = true
	return nil
***REMOVED***

// MarkShorthandDeprecated will mark the shorthand of a flag deprecated in your
// program. It will continue to function but will not show up in help or usage
// messages. Using this flag will also print the given usageMessage.
func (f *FlagSet) MarkShorthandDeprecated(name string, usageMessage string) error ***REMOVED***
	flag := f.Lookup(name)
	if flag == nil ***REMOVED***
		return fmt.Errorf("flag %q does not exist", name)
	***REMOVED***
	if usageMessage == "" ***REMOVED***
		return fmt.Errorf("deprecated message for flag %q must be set", name)
	***REMOVED***
	flag.ShorthandDeprecated = usageMessage
	return nil
***REMOVED***

// MarkHidden sets a flag to 'hidden' in your program. It will continue to
// function but will not show up in help or usage messages.
func (f *FlagSet) MarkHidden(name string) error ***REMOVED***
	flag := f.Lookup(name)
	if flag == nil ***REMOVED***
		return fmt.Errorf("flag %q does not exist", name)
	***REMOVED***
	flag.Hidden = true
	return nil
***REMOVED***

// Lookup returns the Flag structure of the named command-line flag,
// returning nil if none exists.
func Lookup(name string) *Flag ***REMOVED***
	return CommandLine.Lookup(name)
***REMOVED***

// ShorthandLookup returns the Flag structure of the short handed flag,
// returning nil if none exists.
func ShorthandLookup(name string) *Flag ***REMOVED***
	return CommandLine.ShorthandLookup(name)
***REMOVED***

// Set sets the value of the named flag.
func (f *FlagSet) Set(name, value string) error ***REMOVED***
	normalName := f.normalizeFlagName(name)
	flag, ok := f.formal[normalName]
	if !ok ***REMOVED***
		return fmt.Errorf("no such flag -%v", name)
	***REMOVED***

	err := flag.Value.Set(value)
	if err != nil ***REMOVED***
		var flagName string
		if flag.Shorthand != "" && flag.ShorthandDeprecated == "" ***REMOVED***
			flagName = fmt.Sprintf("-%s, --%s", flag.Shorthand, flag.Name)
		***REMOVED*** else ***REMOVED***
			flagName = fmt.Sprintf("--%s", flag.Name)
		***REMOVED***
		return fmt.Errorf("invalid argument %q for %q flag: %v", value, flagName, err)
	***REMOVED***

	if !flag.Changed ***REMOVED***
		if f.actual == nil ***REMOVED***
			f.actual = make(map[NormalizedName]*Flag)
		***REMOVED***
		f.actual[normalName] = flag
		f.orderedActual = append(f.orderedActual, flag)

		flag.Changed = true
	***REMOVED***

	if flag.Deprecated != "" ***REMOVED***
		fmt.Fprintf(f.out(), "Flag --%s has been deprecated, %s\n", flag.Name, flag.Deprecated)
	***REMOVED***
	return nil
***REMOVED***

// SetAnnotation allows one to set arbitrary annotations on a flag in the FlagSet.
// This is sometimes used by spf13/cobra programs which want to generate additional
// bash completion information.
func (f *FlagSet) SetAnnotation(name, key string, values []string) error ***REMOVED***
	normalName := f.normalizeFlagName(name)
	flag, ok := f.formal[normalName]
	if !ok ***REMOVED***
		return fmt.Errorf("no such flag -%v", name)
	***REMOVED***
	if flag.Annotations == nil ***REMOVED***
		flag.Annotations = map[string][]string***REMOVED******REMOVED***
	***REMOVED***
	flag.Annotations[key] = values
	return nil
***REMOVED***

// Changed returns true if the flag was explicitly set during Parse() and false
// otherwise
func (f *FlagSet) Changed(name string) bool ***REMOVED***
	flag := f.Lookup(name)
	// If a flag doesn't exist, it wasn't changed....
	if flag == nil ***REMOVED***
		return false
	***REMOVED***
	return flag.Changed
***REMOVED***

// Set sets the value of the named command-line flag.
func Set(name, value string) error ***REMOVED***
	return CommandLine.Set(name, value)
***REMOVED***

// PrintDefaults prints, to standard error unless configured
// otherwise, the default values of all defined flags in the set.
func (f *FlagSet) PrintDefaults() ***REMOVED***
	usages := f.FlagUsages()
	fmt.Fprint(f.out(), usages)
***REMOVED***

// defaultIsZeroValue returns true if the default value for this flag represents
// a zero value.
func (f *Flag) defaultIsZeroValue() bool ***REMOVED***
	switch f.Value.(type) ***REMOVED***
	case boolFlag:
		return f.DefValue == "false"
	case *durationValue:
		// Beginning in Go 1.7, duration zero values are "0s"
		return f.DefValue == "0" || f.DefValue == "0s"
	case *intValue, *int8Value, *int32Value, *int64Value, *uintValue, *uint8Value, *uint16Value, *uint32Value, *uint64Value, *countValue, *float32Value, *float64Value:
		return f.DefValue == "0"
	case *stringValue:
		return f.DefValue == ""
	case *ipValue, *ipMaskValue, *ipNetValue:
		return f.DefValue == "<nil>"
	case *intSliceValue, *stringSliceValue, *stringArrayValue:
		return f.DefValue == "[]"
	default:
		switch f.Value.String() ***REMOVED***
		case "false":
			return true
		case "<nil>":
			return true
		case "":
			return true
		case "0":
			return true
		***REMOVED***
		return false
	***REMOVED***
***REMOVED***

// UnquoteUsage extracts a back-quoted name from the usage
// string for a flag and returns it and the un-quoted usage.
// Given "a `name` to show" it returns ("name", "a name to show").
// If there are no back quotes, the name is an educated guess of the
// type of the flag's value, or the empty string if the flag is boolean.
func UnquoteUsage(flag *Flag) (name string, usage string) ***REMOVED***
	// Look for a back-quoted name, but avoid the strings package.
	usage = flag.Usage
	for i := 0; i < len(usage); i++ ***REMOVED***
		if usage[i] == '`' ***REMOVED***
			for j := i + 1; j < len(usage); j++ ***REMOVED***
				if usage[j] == '`' ***REMOVED***
					name = usage[i+1 : j]
					usage = usage[:i] + name + usage[j+1:]
					return name, usage
				***REMOVED***
			***REMOVED***
			break // Only one back quote; use type name.
		***REMOVED***
	***REMOVED***

	name = flag.Value.Type()
	switch name ***REMOVED***
	case "bool":
		name = ""
	case "float64":
		name = "float"
	case "int64":
		name = "int"
	case "uint64":
		name = "uint"
	case "stringSlice":
		name = "strings"
	case "intSlice":
		name = "ints"
	case "uintSlice":
		name = "uints"
	case "boolSlice":
		name = "bools"
	***REMOVED***

	return
***REMOVED***

// Splits the string `s` on whitespace into an initial substring up to
// `i` runes in length and the remainder. Will go `slop` over `i` if
// that encompasses the entire string (which allows the caller to
// avoid short orphan words on the final line).
func wrapN(i, slop int, s string) (string, string) ***REMOVED***
	if i+slop > len(s) ***REMOVED***
		return s, ""
	***REMOVED***

	w := strings.LastIndexAny(s[:i], " \t\n")
	if w <= 0 ***REMOVED***
		return s, ""
	***REMOVED***
	nlPos := strings.LastIndex(s[:i], "\n")
	if nlPos > 0 && nlPos < w ***REMOVED***
		return s[:nlPos], s[nlPos+1:]
	***REMOVED***
	return s[:w], s[w+1:]
***REMOVED***

// Wraps the string `s` to a maximum width `w` with leading indent
// `i`. The first line is not indented (this is assumed to be done by
// caller). Pass `w` == 0 to do no wrapping
func wrap(i, w int, s string) string ***REMOVED***
	if w == 0 ***REMOVED***
		return strings.Replace(s, "\n", "\n"+strings.Repeat(" ", i), -1)
	***REMOVED***

	// space between indent i and end of line width w into which
	// we should wrap the text.
	wrap := w - i

	var r, l string

	// Not enough space for sensible wrapping. Wrap as a block on
	// the next line instead.
	if wrap < 24 ***REMOVED***
		i = 16
		wrap = w - i
		r += "\n" + strings.Repeat(" ", i)
	***REMOVED***
	// If still not enough space then don't even try to wrap.
	if wrap < 24 ***REMOVED***
		return strings.Replace(s, "\n", r, -1)
	***REMOVED***

	// Try to avoid short orphan words on the final line, by
	// allowing wrapN to go a bit over if that would fit in the
	// remainder of the line.
	slop := 5
	wrap = wrap - slop

	// Handle first line, which is indented by the caller (or the
	// special case above)
	l, s = wrapN(wrap, slop, s)
	r = r + strings.Replace(l, "\n", "\n"+strings.Repeat(" ", i), -1)

	// Now wrap the rest
	for s != "" ***REMOVED***
		var t string

		t, s = wrapN(wrap, slop, s)
		r = r + "\n" + strings.Repeat(" ", i) + strings.Replace(t, "\n", "\n"+strings.Repeat(" ", i), -1)
	***REMOVED***

	return r

***REMOVED***

// FlagUsagesWrapped returns a string containing the usage information
// for all flags in the FlagSet. Wrapped to `cols` columns (0 for no
// wrapping)
func (f *FlagSet) FlagUsagesWrapped(cols int) string ***REMOVED***
	buf := new(bytes.Buffer)

	lines := make([]string, 0, len(f.formal))

	maxlen := 0
	f.VisitAll(func(flag *Flag) ***REMOVED***
		if flag.Hidden ***REMOVED***
			return
		***REMOVED***

		line := ""
		if flag.Shorthand != "" && flag.ShorthandDeprecated == "" ***REMOVED***
			line = fmt.Sprintf("  -%s, --%s", flag.Shorthand, flag.Name)
		***REMOVED*** else ***REMOVED***
			line = fmt.Sprintf("      --%s", flag.Name)
		***REMOVED***

		varname, usage := UnquoteUsage(flag)
		if varname != "" ***REMOVED***
			line += " " + varname
		***REMOVED***
		if flag.NoOptDefVal != "" ***REMOVED***
			switch flag.Value.Type() ***REMOVED***
			case "string":
				line += fmt.Sprintf("[=\"%s\"]", flag.NoOptDefVal)
			case "bool":
				if flag.NoOptDefVal != "true" ***REMOVED***
					line += fmt.Sprintf("[=%s]", flag.NoOptDefVal)
				***REMOVED***
			case "count":
				if flag.NoOptDefVal != "+1" ***REMOVED***
					line += fmt.Sprintf("[=%s]", flag.NoOptDefVal)
				***REMOVED***
			default:
				line += fmt.Sprintf("[=%s]", flag.NoOptDefVal)
			***REMOVED***
		***REMOVED***

		// This special character will be replaced with spacing once the
		// correct alignment is calculated
		line += "\x00"
		if len(line) > maxlen ***REMOVED***
			maxlen = len(line)
		***REMOVED***

		line += usage
		if !flag.defaultIsZeroValue() ***REMOVED***
			if flag.Value.Type() == "string" ***REMOVED***
				line += fmt.Sprintf(" (default %q)", flag.DefValue)
			***REMOVED*** else ***REMOVED***
				line += fmt.Sprintf(" (default %s)", flag.DefValue)
			***REMOVED***
		***REMOVED***
		if len(flag.Deprecated) != 0 ***REMOVED***
			line += fmt.Sprintf(" (DEPRECATED: %s)", flag.Deprecated)
		***REMOVED***

		lines = append(lines, line)
	***REMOVED***)

	for _, line := range lines ***REMOVED***
		sidx := strings.Index(line, "\x00")
		spacing := strings.Repeat(" ", maxlen-sidx)
		// maxlen + 2 comes from + 1 for the \x00 and + 1 for the (deliberate) off-by-one in maxlen-sidx
		fmt.Fprintln(buf, line[:sidx], spacing, wrap(maxlen+2, cols, line[sidx+1:]))
	***REMOVED***

	return buf.String()
***REMOVED***

// FlagUsages returns a string containing the usage information for all flags in
// the FlagSet
func (f *FlagSet) FlagUsages() string ***REMOVED***
	return f.FlagUsagesWrapped(0)
***REMOVED***

// PrintDefaults prints to standard error the default values of all defined command-line flags.
func PrintDefaults() ***REMOVED***
	CommandLine.PrintDefaults()
***REMOVED***

// defaultUsage is the default function to print a usage message.
func defaultUsage(f *FlagSet) ***REMOVED***
	fmt.Fprintf(f.out(), "Usage of %s:\n", f.name)
	f.PrintDefaults()
***REMOVED***

// NOTE: Usage is not just defaultUsage(CommandLine)
// because it serves (via godoc flag Usage) as the example
// for how to write your own usage function.

// Usage prints to standard error a usage message documenting all defined command-line flags.
// The function is a variable that may be changed to point to a custom function.
// By default it prints a simple header and calls PrintDefaults; for details about the
// format of the output and how to control it, see the documentation for PrintDefaults.
var Usage = func() ***REMOVED***
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	PrintDefaults()
***REMOVED***

// NFlag returns the number of flags that have been set.
func (f *FlagSet) NFlag() int ***REMOVED*** return len(f.actual) ***REMOVED***

// NFlag returns the number of command-line flags that have been set.
func NFlag() int ***REMOVED*** return len(CommandLine.actual) ***REMOVED***

// Arg returns the i'th argument.  Arg(0) is the first remaining argument
// after flags have been processed.
func (f *FlagSet) Arg(i int) string ***REMOVED***
	if i < 0 || i >= len(f.args) ***REMOVED***
		return ""
	***REMOVED***
	return f.args[i]
***REMOVED***

// Arg returns the i'th command-line argument.  Arg(0) is the first remaining argument
// after flags have been processed.
func Arg(i int) string ***REMOVED***
	return CommandLine.Arg(i)
***REMOVED***

// NArg is the number of arguments remaining after flags have been processed.
func (f *FlagSet) NArg() int ***REMOVED*** return len(f.args) ***REMOVED***

// NArg is the number of arguments remaining after flags have been processed.
func NArg() int ***REMOVED*** return len(CommandLine.args) ***REMOVED***

// Args returns the non-flag arguments.
func (f *FlagSet) Args() []string ***REMOVED*** return f.args ***REMOVED***

// Args returns the non-flag command-line arguments.
func Args() []string ***REMOVED*** return CommandLine.args ***REMOVED***

// Var defines a flag with the specified name and usage string. The type and
// value of the flag are represented by the first argument, of type Value, which
// typically holds a user-defined implementation of Value. For instance, the
// caller could create a flag that turns a comma-separated string into a slice
// of strings by giving the slice the methods of Value; in particular, Set would
// decompose the comma-separated string into the slice.
func (f *FlagSet) Var(value Value, name string, usage string) ***REMOVED***
	f.VarP(value, name, "", usage)
***REMOVED***

// VarPF is like VarP, but returns the flag created
func (f *FlagSet) VarPF(value Value, name, shorthand, usage string) *Flag ***REMOVED***
	// Remember the default value as a string; it won't change.
	flag := &Flag***REMOVED***
		Name:      name,
		Shorthand: shorthand,
		Usage:     usage,
		Value:     value,
		DefValue:  value.String(),
	***REMOVED***
	f.AddFlag(flag)
	return flag
***REMOVED***

// VarP is like Var, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) VarP(value Value, name, shorthand, usage string) ***REMOVED***
	f.VarPF(value, name, shorthand, usage)
***REMOVED***

// AddFlag will add the flag to the FlagSet
func (f *FlagSet) AddFlag(flag *Flag) ***REMOVED***
	normalizedFlagName := f.normalizeFlagName(flag.Name)

	_, alreadyThere := f.formal[normalizedFlagName]
	if alreadyThere ***REMOVED***
		msg := fmt.Sprintf("%s flag redefined: %s", f.name, flag.Name)
		fmt.Fprintln(f.out(), msg)
		panic(msg) // Happens only if flags are declared with identical names
	***REMOVED***
	if f.formal == nil ***REMOVED***
		f.formal = make(map[NormalizedName]*Flag)
	***REMOVED***

	flag.Name = string(normalizedFlagName)
	f.formal[normalizedFlagName] = flag
	f.orderedFormal = append(f.orderedFormal, flag)

	if flag.Shorthand == "" ***REMOVED***
		return
	***REMOVED***
	if len(flag.Shorthand) > 1 ***REMOVED***
		msg := fmt.Sprintf("%q shorthand is more than one ASCII character", flag.Shorthand)
		fmt.Fprintf(f.out(), msg)
		panic(msg)
	***REMOVED***
	if f.shorthands == nil ***REMOVED***
		f.shorthands = make(map[byte]*Flag)
	***REMOVED***
	c := flag.Shorthand[0]
	used, alreadyThere := f.shorthands[c]
	if alreadyThere ***REMOVED***
		msg := fmt.Sprintf("unable to redefine %q shorthand in %q flagset: it's already used for %q flag", c, f.name, used.Name)
		fmt.Fprintf(f.out(), msg)
		panic(msg)
	***REMOVED***
	f.shorthands[c] = flag
***REMOVED***

// AddFlagSet adds one FlagSet to another. If a flag is already present in f
// the flag from newSet will be ignored.
func (f *FlagSet) AddFlagSet(newSet *FlagSet) ***REMOVED***
	if newSet == nil ***REMOVED***
		return
	***REMOVED***
	newSet.VisitAll(func(flag *Flag) ***REMOVED***
		if f.Lookup(flag.Name) == nil ***REMOVED***
			f.AddFlag(flag)
		***REMOVED***
	***REMOVED***)
***REMOVED***

// Var defines a flag with the specified name and usage string. The type and
// value of the flag are represented by the first argument, of type Value, which
// typically holds a user-defined implementation of Value. For instance, the
// caller could create a flag that turns a comma-separated string into a slice
// of strings by giving the slice the methods of Value; in particular, Set would
// decompose the comma-separated string into the slice.
func Var(value Value, name string, usage string) ***REMOVED***
	CommandLine.VarP(value, name, "", usage)
***REMOVED***

// VarP is like Var, but accepts a shorthand letter that can be used after a single dash.
func VarP(value Value, name, shorthand, usage string) ***REMOVED***
	CommandLine.VarP(value, name, shorthand, usage)
***REMOVED***

// failf prints to standard error a formatted error and usage message and
// returns the error.
func (f *FlagSet) failf(format string, a ...interface***REMOVED******REMOVED***) error ***REMOVED***
	err := fmt.Errorf(format, a...)
	if f.errorHandling != ContinueOnError ***REMOVED***
		fmt.Fprintln(f.out(), err)
		f.usage()
	***REMOVED***
	return err
***REMOVED***

// usage calls the Usage method for the flag set, or the usage function if
// the flag set is CommandLine.
func (f *FlagSet) usage() ***REMOVED***
	if f == CommandLine ***REMOVED***
		Usage()
	***REMOVED*** else if f.Usage == nil ***REMOVED***
		defaultUsage(f)
	***REMOVED*** else ***REMOVED***
		f.Usage()
	***REMOVED***
***REMOVED***

//--unknown (args will be empty)
//--unknown --next-flag ... (args will be --next-flag ...)
//--unknown arg ... (args will be arg ...)
func stripUnknownFlagValue(args []string) []string ***REMOVED***
	if len(args) == 0 ***REMOVED***
		//--unknown
		return args
	***REMOVED***

	first := args[0]
	if first[0] == '-' ***REMOVED***
		//--unknown --next-flag ...
		return args
	***REMOVED***

	//--unknown arg ... (args will be arg ...)
	return args[1:]
***REMOVED***

func (f *FlagSet) parseLongArg(s string, args []string, fn parseFunc) (a []string, err error) ***REMOVED***
	a = args
	name := s[2:]
	if len(name) == 0 || name[0] == '-' || name[0] == '=' ***REMOVED***
		err = f.failf("bad flag syntax: %s", s)
		return
	***REMOVED***

	split := strings.SplitN(name, "=", 2)
	name = split[0]
	flag, exists := f.formal[f.normalizeFlagName(name)]

	if !exists ***REMOVED***
		switch ***REMOVED***
		case name == "help":
			f.usage()
			return a, ErrHelp
		case f.ParseErrorsWhitelist.UnknownFlags:
			// --unknown=unknownval arg ...
			// we do not want to lose arg in this case
			if len(split) >= 2 ***REMOVED***
				return a, nil
			***REMOVED***

			return stripUnknownFlagValue(a), nil
		default:
			err = f.failf("unknown flag: --%s", name)
			return
		***REMOVED***
	***REMOVED***

	var value string
	if len(split) == 2 ***REMOVED***
		// '--flag=arg'
		value = split[1]
	***REMOVED*** else if flag.NoOptDefVal != "" ***REMOVED***
		// '--flag' (arg was optional)
		value = flag.NoOptDefVal
	***REMOVED*** else if len(a) > 0 ***REMOVED***
		// '--flag arg'
		value = a[0]
		a = a[1:]
	***REMOVED*** else ***REMOVED***
		// '--flag' (arg was required)
		err = f.failf("flag needs an argument: %s", s)
		return
	***REMOVED***

	err = fn(flag, value)
	if err != nil ***REMOVED***
		f.failf(err.Error())
	***REMOVED***
	return
***REMOVED***

func (f *FlagSet) parseSingleShortArg(shorthands string, args []string, fn parseFunc) (outShorts string, outArgs []string, err error) ***REMOVED***
	if strings.HasPrefix(shorthands, "test.") ***REMOVED***
		return
	***REMOVED***

	outArgs = args
	outShorts = shorthands[1:]
	c := shorthands[0]

	flag, exists := f.shorthands[c]
	if !exists ***REMOVED***
		switch ***REMOVED***
		case c == 'h':
			f.usage()
			err = ErrHelp
			return
		case f.ParseErrorsWhitelist.UnknownFlags:
			// '-f=arg arg ...'
			// we do not want to lose arg in this case
			if len(shorthands) > 2 && shorthands[1] == '=' ***REMOVED***
				outShorts = ""
				return
			***REMOVED***

			outArgs = stripUnknownFlagValue(outArgs)
			return
		default:
			err = f.failf("unknown shorthand flag: %q in -%s", c, shorthands)
			return
		***REMOVED***
	***REMOVED***

	var value string
	if len(shorthands) > 2 && shorthands[1] == '=' ***REMOVED***
		// '-f=arg'
		value = shorthands[2:]
		outShorts = ""
	***REMOVED*** else if flag.NoOptDefVal != "" ***REMOVED***
		// '-f' (arg was optional)
		value = flag.NoOptDefVal
	***REMOVED*** else if len(shorthands) > 1 ***REMOVED***
		// '-farg'
		value = shorthands[1:]
		outShorts = ""
	***REMOVED*** else if len(args) > 0 ***REMOVED***
		// '-f arg'
		value = args[0]
		outArgs = args[1:]
	***REMOVED*** else ***REMOVED***
		// '-f' (arg was required)
		err = f.failf("flag needs an argument: %q in -%s", c, shorthands)
		return
	***REMOVED***

	if flag.ShorthandDeprecated != "" ***REMOVED***
		fmt.Fprintf(f.out(), "Flag shorthand -%s has been deprecated, %s\n", flag.Shorthand, flag.ShorthandDeprecated)
	***REMOVED***

	err = fn(flag, value)
	if err != nil ***REMOVED***
		f.failf(err.Error())
	***REMOVED***
	return
***REMOVED***

func (f *FlagSet) parseShortArg(s string, args []string, fn parseFunc) (a []string, err error) ***REMOVED***
	a = args
	shorthands := s[1:]

	// "shorthands" can be a series of shorthand letters of flags (e.g. "-vvv").
	for len(shorthands) > 0 ***REMOVED***
		shorthands, a, err = f.parseSingleShortArg(shorthands, args, fn)
		if err != nil ***REMOVED***
			return
		***REMOVED***
	***REMOVED***

	return
***REMOVED***

func (f *FlagSet) parseArgs(args []string, fn parseFunc) (err error) ***REMOVED***
	for len(args) > 0 ***REMOVED***
		s := args[0]
		args = args[1:]
		if len(s) == 0 || s[0] != '-' || len(s) == 1 ***REMOVED***
			if !f.interspersed ***REMOVED***
				f.args = append(f.args, s)
				f.args = append(f.args, args...)
				return nil
			***REMOVED***
			f.args = append(f.args, s)
			continue
		***REMOVED***

		if s[1] == '-' ***REMOVED***
			if len(s) == 2 ***REMOVED*** // "--" terminates the flags
				f.argsLenAtDash = len(f.args)
				f.args = append(f.args, args...)
				break
			***REMOVED***
			args, err = f.parseLongArg(s, args, fn)
		***REMOVED*** else ***REMOVED***
			args, err = f.parseShortArg(s, args, fn)
		***REMOVED***
		if err != nil ***REMOVED***
			return
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

// Parse parses flag definitions from the argument list, which should not
// include the command name.  Must be called after all flags in the FlagSet
// are defined and before flags are accessed by the program.
// The return value will be ErrHelp if -help was set but not defined.
func (f *FlagSet) Parse(arguments []string) error ***REMOVED***
	if f.addedGoFlagSets != nil ***REMOVED***
		for _, goFlagSet := range f.addedGoFlagSets ***REMOVED***
			goFlagSet.Parse(nil)
		***REMOVED***
	***REMOVED***
	f.parsed = true

	if len(arguments) < 0 ***REMOVED***
		return nil
	***REMOVED***

	f.args = make([]string, 0, len(arguments))

	set := func(flag *Flag, value string) error ***REMOVED***
		return f.Set(flag.Name, value)
	***REMOVED***

	err := f.parseArgs(arguments, set)
	if err != nil ***REMOVED***
		switch f.errorHandling ***REMOVED***
		case ContinueOnError:
			return err
		case ExitOnError:
			fmt.Println(err)
			os.Exit(2)
		case PanicOnError:
			panic(err)
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

type parseFunc func(flag *Flag, value string) error

// ParseAll parses flag definitions from the argument list, which should not
// include the command name. The arguments for fn are flag and value. Must be
// called after all flags in the FlagSet are defined and before flags are
// accessed by the program. The return value will be ErrHelp if -help was set
// but not defined.
func (f *FlagSet) ParseAll(arguments []string, fn func(flag *Flag, value string) error) error ***REMOVED***
	f.parsed = true
	f.args = make([]string, 0, len(arguments))

	err := f.parseArgs(arguments, fn)
	if err != nil ***REMOVED***
		switch f.errorHandling ***REMOVED***
		case ContinueOnError:
			return err
		case ExitOnError:
			os.Exit(2)
		case PanicOnError:
			panic(err)
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// Parsed reports whether f.Parse has been called.
func (f *FlagSet) Parsed() bool ***REMOVED***
	return f.parsed
***REMOVED***

// Parse parses the command-line flags from os.Args[1:].  Must be called
// after all flags are defined and before flags are accessed by the program.
func Parse() ***REMOVED***
	// Ignore errors; CommandLine is set for ExitOnError.
	CommandLine.Parse(os.Args[1:])
***REMOVED***

// ParseAll parses the command-line flags from os.Args[1:] and called fn for each.
// The arguments for fn are flag and value. Must be called after all flags are
// defined and before flags are accessed by the program.
func ParseAll(fn func(flag *Flag, value string) error) ***REMOVED***
	// Ignore errors; CommandLine is set for ExitOnError.
	CommandLine.ParseAll(os.Args[1:], fn)
***REMOVED***

// SetInterspersed sets whether to support interspersed option/non-option arguments.
func SetInterspersed(interspersed bool) ***REMOVED***
	CommandLine.SetInterspersed(interspersed)
***REMOVED***

// Parsed returns true if the command-line flags have been parsed.
func Parsed() bool ***REMOVED***
	return CommandLine.Parsed()
***REMOVED***

// CommandLine is the default set of command-line flags, parsed from os.Args.
var CommandLine = NewFlagSet(os.Args[0], ExitOnError)

// NewFlagSet returns a new, empty flag set with the specified name,
// error handling property and SortFlags set to true.
func NewFlagSet(name string, errorHandling ErrorHandling) *FlagSet ***REMOVED***
	f := &FlagSet***REMOVED***
		name:          name,
		errorHandling: errorHandling,
		argsLenAtDash: -1,
		interspersed:  true,
		SortFlags:     true,
	***REMOVED***
	return f
***REMOVED***

// SetInterspersed sets whether to support interspersed option/non-option arguments.
func (f *FlagSet) SetInterspersed(interspersed bool) ***REMOVED***
	f.interspersed = interspersed
***REMOVED***

// Init sets the name and error handling property for a flag set.
// By default, the zero FlagSet uses an empty name and the
// ContinueOnError error handling policy.
func (f *FlagSet) Init(name string, errorHandling ErrorHandling) ***REMOVED***
	f.name = name
	f.errorHandling = errorHandling
	f.argsLenAtDash = -1
***REMOVED***
