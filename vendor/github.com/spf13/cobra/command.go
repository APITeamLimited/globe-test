// Copyright © 2013 Steve Francia <spf@spf13.com>.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package cobra is a commander providing a simple interface to create powerful modern CLI interfaces.
// In addition to providing an interface, Cobra simultaneously provides a controller to organize your application code.
package cobra

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	flag "github.com/spf13/pflag"
)

// FParseErrWhitelist configures Flag parse errors to be ignored
type FParseErrWhitelist flag.ParseErrorsWhitelist

// Command is just that, a command for your application.
// E.g.  'go run ...' - 'run' is the command. Cobra requires
// you to define the usage and description as part of your command
// definition to ensure usability.
type Command struct ***REMOVED***
	// Use is the one-line usage message.
	// Recommended syntax is as follow:
	//   [ ] identifies an optional argument. Arguments that are not enclosed in brackets are required.
	//   ... indicates that you can specify multiple values for the previous argument.
	//   |   indicates mutually exclusive information. You can use the argument to the left of the separator or the
	//       argument to the right of the separator. You cannot use both arguments in a single use of the command.
	//   ***REMOVED*** ***REMOVED*** delimits a set of mutually exclusive arguments when one of the arguments is required. If the arguments are
	//       optional, they are enclosed in brackets ([ ]).
	// Example: add [-F file | -D dir]... [-f format] profile
	Use string

	// Aliases is an array of aliases that can be used instead of the first word in Use.
	Aliases []string

	// SuggestFor is an array of command names for which this command will be suggested -
	// similar to aliases but only suggests.
	SuggestFor []string

	// Short is the short description shown in the 'help' output.
	Short string

	// Long is the long message shown in the 'help <this-command>' output.
	Long string

	// Example is examples of how to use the command.
	Example string

	// ValidArgs is list of all valid non-flag arguments that are accepted in shell completions
	ValidArgs []string
	// ValidArgsFunction is an optional function that provides valid non-flag arguments for shell completion.
	// It is a dynamic version of using ValidArgs.
	// Only one of ValidArgs and ValidArgsFunction can be used for a command.
	ValidArgsFunction func(cmd *Command, args []string, toComplete string) ([]string, ShellCompDirective)

	// Expected arguments
	Args PositionalArgs

	// ArgAliases is List of aliases for ValidArgs.
	// These are not suggested to the user in the shell completion,
	// but accepted if entered manually.
	ArgAliases []string

	// BashCompletionFunction is custom bash functions used by the legacy bash autocompletion generator.
	// For portability with other shells, it is recommended to instead use ValidArgsFunction
	BashCompletionFunction string

	// Deprecated defines, if this command is deprecated and should print this string when used.
	Deprecated string

	// Annotations are key/value pairs that can be used by applications to identify or
	// group commands.
	Annotations map[string]string

	// Version defines the version for this command. If this value is non-empty and the command does not
	// define a "version" flag, a "version" boolean flag will be added to the command and, if specified,
	// will print content of the "Version" variable. A shorthand "v" flag will also be added if the
	// command does not define one.
	Version string

	// The *Run functions are executed in the following order:
	//   * PersistentPreRun()
	//   * PreRun()
	//   * Run()
	//   * PostRun()
	//   * PersistentPostRun()
	// All functions get the same args, the arguments after the command name.
	//
	// PersistentPreRun: children of this command will inherit and execute.
	PersistentPreRun func(cmd *Command, args []string)
	// PersistentPreRunE: PersistentPreRun but returns an error.
	PersistentPreRunE func(cmd *Command, args []string) error
	// PreRun: children of this command will not inherit.
	PreRun func(cmd *Command, args []string)
	// PreRunE: PreRun but returns an error.
	PreRunE func(cmd *Command, args []string) error
	// Run: Typically the actual work function. Most commands will only implement this.
	Run func(cmd *Command, args []string)
	// RunE: Run but returns an error.
	RunE func(cmd *Command, args []string) error
	// PostRun: run after the Run command.
	PostRun func(cmd *Command, args []string)
	// PostRunE: PostRun but returns an error.
	PostRunE func(cmd *Command, args []string) error
	// PersistentPostRun: children of this command will inherit and execute after PostRun.
	PersistentPostRun func(cmd *Command, args []string)
	// PersistentPostRunE: PersistentPostRun but returns an error.
	PersistentPostRunE func(cmd *Command, args []string) error

	// args is actual args parsed from flags.
	args []string
	// flagErrorBuf contains all error messages from pflag.
	flagErrorBuf *bytes.Buffer
	// flags is full set of flags.
	flags *flag.FlagSet
	// pflags contains persistent flags.
	pflags *flag.FlagSet
	// lflags contains local flags.
	lflags *flag.FlagSet
	// iflags contains inherited flags.
	iflags *flag.FlagSet
	// parentsPflags is all persistent flags of cmd's parents.
	parentsPflags *flag.FlagSet
	// globNormFunc is the global normalization function
	// that we can use on every pflag set and children commands
	globNormFunc func(f *flag.FlagSet, name string) flag.NormalizedName

	// usageFunc is usage func defined by user.
	usageFunc func(*Command) error
	// usageTemplate is usage template defined by user.
	usageTemplate string
	// flagErrorFunc is func defined by user and it's called when the parsing of
	// flags returns an error.
	flagErrorFunc func(*Command, error) error
	// helpTemplate is help template defined by user.
	helpTemplate string
	// helpFunc is help func defined by user.
	helpFunc func(*Command, []string)
	// helpCommand is command with usage 'help'. If it's not defined by user,
	// cobra uses default help command.
	helpCommand *Command
	// versionTemplate is the version template defined by user.
	versionTemplate string

	// inReader is a reader defined by the user that replaces stdin
	inReader io.Reader
	// outWriter is a writer defined by the user that replaces stdout
	outWriter io.Writer
	// errWriter is a writer defined by the user that replaces stderr
	errWriter io.Writer

	//FParseErrWhitelist flag parse errors to be ignored
	FParseErrWhitelist FParseErrWhitelist

	// CompletionOptions is a set of options to control the handling of shell completion
	CompletionOptions CompletionOptions

	// commandsAreSorted defines, if command slice are sorted or not.
	commandsAreSorted bool
	// commandCalledAs is the name or alias value used to call this command.
	commandCalledAs struct ***REMOVED***
		name   string
		called bool
	***REMOVED***

	ctx context.Context

	// commands is the list of commands supported by this program.
	commands []*Command
	// parent is a parent command for this command.
	parent *Command
	// Max lengths of commands' string lengths for use in padding.
	commandsMaxUseLen         int
	commandsMaxCommandPathLen int
	commandsMaxNameLen        int

	// TraverseChildren parses flags on all parents before executing child command.
	TraverseChildren bool

	// Hidden defines, if this command is hidden and should NOT show up in the list of available commands.
	Hidden bool

	// SilenceErrors is an option to quiet errors down stream.
	SilenceErrors bool

	// SilenceUsage is an option to silence usage when an error occurs.
	SilenceUsage bool

	// DisableFlagParsing disables the flag parsing.
	// If this is true all flags will be passed to the command as arguments.
	DisableFlagParsing bool

	// DisableAutoGenTag defines, if gen tag ("Auto generated by spf13/cobra...")
	// will be printed by generating docs for this command.
	DisableAutoGenTag bool

	// DisableFlagsInUseLine will disable the addition of [flags] to the usage
	// line of a command when printing help or generating docs
	DisableFlagsInUseLine bool

	// DisableSuggestions disables the suggestions based on Levenshtein distance
	// that go along with 'unknown command' messages.
	DisableSuggestions bool

	// SuggestionsMinimumDistance defines minimum levenshtein distance to display suggestions.
	// Must be > 0.
	SuggestionsMinimumDistance int
***REMOVED***

// Context returns underlying command context. If command wasn't
// executed with ExecuteContext Context returns Background context.
func (c *Command) Context() context.Context ***REMOVED***
	return c.ctx
***REMOVED***

// SetArgs sets arguments for the command. It is set to os.Args[1:] by default, if desired, can be overridden
// particularly useful when testing.
func (c *Command) SetArgs(a []string) ***REMOVED***
	c.args = a
***REMOVED***

// SetOutput sets the destination for usage and error messages.
// If output is nil, os.Stderr is used.
// Deprecated: Use SetOut and/or SetErr instead
func (c *Command) SetOutput(output io.Writer) ***REMOVED***
	c.outWriter = output
	c.errWriter = output
***REMOVED***

// SetOut sets the destination for usage messages.
// If newOut is nil, os.Stdout is used.
func (c *Command) SetOut(newOut io.Writer) ***REMOVED***
	c.outWriter = newOut
***REMOVED***

// SetErr sets the destination for error messages.
// If newErr is nil, os.Stderr is used.
func (c *Command) SetErr(newErr io.Writer) ***REMOVED***
	c.errWriter = newErr
***REMOVED***

// SetIn sets the source for input data
// If newIn is nil, os.Stdin is used.
func (c *Command) SetIn(newIn io.Reader) ***REMOVED***
	c.inReader = newIn
***REMOVED***

// SetUsageFunc sets usage function. Usage can be defined by application.
func (c *Command) SetUsageFunc(f func(*Command) error) ***REMOVED***
	c.usageFunc = f
***REMOVED***

// SetUsageTemplate sets usage template. Can be defined by Application.
func (c *Command) SetUsageTemplate(s string) ***REMOVED***
	c.usageTemplate = s
***REMOVED***

// SetFlagErrorFunc sets a function to generate an error when flag parsing
// fails.
func (c *Command) SetFlagErrorFunc(f func(*Command, error) error) ***REMOVED***
	c.flagErrorFunc = f
***REMOVED***

// SetHelpFunc sets help function. Can be defined by Application.
func (c *Command) SetHelpFunc(f func(*Command, []string)) ***REMOVED***
	c.helpFunc = f
***REMOVED***

// SetHelpCommand sets help command.
func (c *Command) SetHelpCommand(cmd *Command) ***REMOVED***
	c.helpCommand = cmd
***REMOVED***

// SetHelpTemplate sets help template to be used. Application can use it to set custom template.
func (c *Command) SetHelpTemplate(s string) ***REMOVED***
	c.helpTemplate = s
***REMOVED***

// SetVersionTemplate sets version template to be used. Application can use it to set custom template.
func (c *Command) SetVersionTemplate(s string) ***REMOVED***
	c.versionTemplate = s
***REMOVED***

// SetGlobalNormalizationFunc sets a normalization function to all flag sets and also to child commands.
// The user should not have a cyclic dependency on commands.
func (c *Command) SetGlobalNormalizationFunc(n func(f *flag.FlagSet, name string) flag.NormalizedName) ***REMOVED***
	c.Flags().SetNormalizeFunc(n)
	c.PersistentFlags().SetNormalizeFunc(n)
	c.globNormFunc = n

	for _, command := range c.commands ***REMOVED***
		command.SetGlobalNormalizationFunc(n)
	***REMOVED***
***REMOVED***

// OutOrStdout returns output to stdout.
func (c *Command) OutOrStdout() io.Writer ***REMOVED***
	return c.getOut(os.Stdout)
***REMOVED***

// OutOrStderr returns output to stderr
func (c *Command) OutOrStderr() io.Writer ***REMOVED***
	return c.getOut(os.Stderr)
***REMOVED***

// ErrOrStderr returns output to stderr
func (c *Command) ErrOrStderr() io.Writer ***REMOVED***
	return c.getErr(os.Stderr)
***REMOVED***

// InOrStdin returns input to stdin
func (c *Command) InOrStdin() io.Reader ***REMOVED***
	return c.getIn(os.Stdin)
***REMOVED***

func (c *Command) getOut(def io.Writer) io.Writer ***REMOVED***
	if c.outWriter != nil ***REMOVED***
		return c.outWriter
	***REMOVED***
	if c.HasParent() ***REMOVED***
		return c.parent.getOut(def)
	***REMOVED***
	return def
***REMOVED***

func (c *Command) getErr(def io.Writer) io.Writer ***REMOVED***
	if c.errWriter != nil ***REMOVED***
		return c.errWriter
	***REMOVED***
	if c.HasParent() ***REMOVED***
		return c.parent.getErr(def)
	***REMOVED***
	return def
***REMOVED***

func (c *Command) getIn(def io.Reader) io.Reader ***REMOVED***
	if c.inReader != nil ***REMOVED***
		return c.inReader
	***REMOVED***
	if c.HasParent() ***REMOVED***
		return c.parent.getIn(def)
	***REMOVED***
	return def
***REMOVED***

// UsageFunc returns either the function set by SetUsageFunc for this command
// or a parent, or it returns a default usage function.
func (c *Command) UsageFunc() (f func(*Command) error) ***REMOVED***
	if c.usageFunc != nil ***REMOVED***
		return c.usageFunc
	***REMOVED***
	if c.HasParent() ***REMOVED***
		return c.Parent().UsageFunc()
	***REMOVED***
	return func(c *Command) error ***REMOVED***
		c.mergePersistentFlags()
		err := tmpl(c.OutOrStderr(), c.UsageTemplate(), c)
		if err != nil ***REMOVED***
			c.PrintErrln(err)
		***REMOVED***
		return err
	***REMOVED***
***REMOVED***

// Usage puts out the usage for the command.
// Used when a user provides invalid input.
// Can be defined by user by overriding UsageFunc.
func (c *Command) Usage() error ***REMOVED***
	return c.UsageFunc()(c)
***REMOVED***

// HelpFunc returns either the function set by SetHelpFunc for this command
// or a parent, or it returns a function with default help behavior.
func (c *Command) HelpFunc() func(*Command, []string) ***REMOVED***
	if c.helpFunc != nil ***REMOVED***
		return c.helpFunc
	***REMOVED***
	if c.HasParent() ***REMOVED***
		return c.Parent().HelpFunc()
	***REMOVED***
	return func(c *Command, a []string) ***REMOVED***
		c.mergePersistentFlags()
		// The help should be sent to stdout
		// See https://github.com/spf13/cobra/issues/1002
		err := tmpl(c.OutOrStdout(), c.HelpTemplate(), c)
		if err != nil ***REMOVED***
			c.PrintErrln(err)
		***REMOVED***
	***REMOVED***
***REMOVED***

// Help puts out the help for the command.
// Used when a user calls help [command].
// Can be defined by user by overriding HelpFunc.
func (c *Command) Help() error ***REMOVED***
	c.HelpFunc()(c, []string***REMOVED******REMOVED***)
	return nil
***REMOVED***

// UsageString returns usage string.
func (c *Command) UsageString() string ***REMOVED***
	// Storing normal writers
	tmpOutput := c.outWriter
	tmpErr := c.errWriter

	bb := new(bytes.Buffer)
	c.outWriter = bb
	c.errWriter = bb

	CheckErr(c.Usage())

	// Setting things back to normal
	c.outWriter = tmpOutput
	c.errWriter = tmpErr

	return bb.String()
***REMOVED***

// FlagErrorFunc returns either the function set by SetFlagErrorFunc for this
// command or a parent, or it returns a function which returns the original
// error.
func (c *Command) FlagErrorFunc() (f func(*Command, error) error) ***REMOVED***
	if c.flagErrorFunc != nil ***REMOVED***
		return c.flagErrorFunc
	***REMOVED***

	if c.HasParent() ***REMOVED***
		return c.parent.FlagErrorFunc()
	***REMOVED***
	return func(c *Command, err error) error ***REMOVED***
		return err
	***REMOVED***
***REMOVED***

var minUsagePadding = 25

// UsagePadding return padding for the usage.
func (c *Command) UsagePadding() int ***REMOVED***
	if c.parent == nil || minUsagePadding > c.parent.commandsMaxUseLen ***REMOVED***
		return minUsagePadding
	***REMOVED***
	return c.parent.commandsMaxUseLen
***REMOVED***

var minCommandPathPadding = 11

// CommandPathPadding return padding for the command path.
func (c *Command) CommandPathPadding() int ***REMOVED***
	if c.parent == nil || minCommandPathPadding > c.parent.commandsMaxCommandPathLen ***REMOVED***
		return minCommandPathPadding
	***REMOVED***
	return c.parent.commandsMaxCommandPathLen
***REMOVED***

var minNamePadding = 11

// NamePadding returns padding for the name.
func (c *Command) NamePadding() int ***REMOVED***
	if c.parent == nil || minNamePadding > c.parent.commandsMaxNameLen ***REMOVED***
		return minNamePadding
	***REMOVED***
	return c.parent.commandsMaxNameLen
***REMOVED***

// UsageTemplate returns usage template for the command.
func (c *Command) UsageTemplate() string ***REMOVED***
	if c.usageTemplate != "" ***REMOVED***
		return c.usageTemplate
	***REMOVED***

	if c.HasParent() ***REMOVED***
		return c.parent.UsageTemplate()
	***REMOVED***
	return `Usage:***REMOVED******REMOVED***if .Runnable***REMOVED******REMOVED***
  ***REMOVED******REMOVED***.UseLine***REMOVED******REMOVED******REMOVED******REMOVED***end***REMOVED******REMOVED******REMOVED******REMOVED***if .HasAvailableSubCommands***REMOVED******REMOVED***
  ***REMOVED******REMOVED***.CommandPath***REMOVED******REMOVED*** [command]***REMOVED******REMOVED***end***REMOVED******REMOVED******REMOVED******REMOVED***if gt (len .Aliases) 0***REMOVED******REMOVED***

Aliases:
  ***REMOVED******REMOVED***.NameAndAliases***REMOVED******REMOVED******REMOVED******REMOVED***end***REMOVED******REMOVED******REMOVED******REMOVED***if .HasExample***REMOVED******REMOVED***

Examples:
***REMOVED******REMOVED***.Example***REMOVED******REMOVED******REMOVED******REMOVED***end***REMOVED******REMOVED******REMOVED******REMOVED***if .HasAvailableSubCommands***REMOVED******REMOVED***

Available Commands:***REMOVED******REMOVED***range .Commands***REMOVED******REMOVED******REMOVED******REMOVED***if (or .IsAvailableCommand (eq .Name "help"))***REMOVED******REMOVED***
  ***REMOVED******REMOVED***rpad .Name .NamePadding ***REMOVED******REMOVED*** ***REMOVED******REMOVED***.Short***REMOVED******REMOVED******REMOVED******REMOVED***end***REMOVED******REMOVED******REMOVED******REMOVED***end***REMOVED******REMOVED******REMOVED******REMOVED***end***REMOVED******REMOVED******REMOVED******REMOVED***if .HasAvailableLocalFlags***REMOVED******REMOVED***

Flags:
***REMOVED******REMOVED***.LocalFlags.FlagUsages | trimTrailingWhitespaces***REMOVED******REMOVED******REMOVED******REMOVED***end***REMOVED******REMOVED******REMOVED******REMOVED***if .HasAvailableInheritedFlags***REMOVED******REMOVED***

Global Flags:
***REMOVED******REMOVED***.InheritedFlags.FlagUsages | trimTrailingWhitespaces***REMOVED******REMOVED******REMOVED******REMOVED***end***REMOVED******REMOVED******REMOVED******REMOVED***if .HasHelpSubCommands***REMOVED******REMOVED***

Additional help topics:***REMOVED******REMOVED***range .Commands***REMOVED******REMOVED******REMOVED******REMOVED***if .IsAdditionalHelpTopicCommand***REMOVED******REMOVED***
  ***REMOVED******REMOVED***rpad .CommandPath .CommandPathPadding***REMOVED******REMOVED*** ***REMOVED******REMOVED***.Short***REMOVED******REMOVED******REMOVED******REMOVED***end***REMOVED******REMOVED******REMOVED******REMOVED***end***REMOVED******REMOVED******REMOVED******REMOVED***end***REMOVED******REMOVED******REMOVED******REMOVED***if .HasAvailableSubCommands***REMOVED******REMOVED***

Use "***REMOVED******REMOVED***.CommandPath***REMOVED******REMOVED*** [command] --help" for more information about a command.***REMOVED******REMOVED***end***REMOVED******REMOVED***
`
***REMOVED***

// HelpTemplate return help template for the command.
func (c *Command) HelpTemplate() string ***REMOVED***
	if c.helpTemplate != "" ***REMOVED***
		return c.helpTemplate
	***REMOVED***

	if c.HasParent() ***REMOVED***
		return c.parent.HelpTemplate()
	***REMOVED***
	return `***REMOVED******REMOVED***with (or .Long .Short)***REMOVED******REMOVED******REMOVED******REMOVED***. | trimTrailingWhitespaces***REMOVED******REMOVED***

***REMOVED******REMOVED***end***REMOVED******REMOVED******REMOVED******REMOVED***if or .Runnable .HasSubCommands***REMOVED******REMOVED******REMOVED******REMOVED***.UsageString***REMOVED******REMOVED******REMOVED******REMOVED***end***REMOVED******REMOVED***`
***REMOVED***

// VersionTemplate return version template for the command.
func (c *Command) VersionTemplate() string ***REMOVED***
	if c.versionTemplate != "" ***REMOVED***
		return c.versionTemplate
	***REMOVED***

	if c.HasParent() ***REMOVED***
		return c.parent.VersionTemplate()
	***REMOVED***
	return `***REMOVED******REMOVED***with .Name***REMOVED******REMOVED******REMOVED******REMOVED***printf "%s " .***REMOVED******REMOVED******REMOVED******REMOVED***end***REMOVED******REMOVED******REMOVED******REMOVED***printf "version %s" .Version***REMOVED******REMOVED***
`
***REMOVED***

func hasNoOptDefVal(name string, fs *flag.FlagSet) bool ***REMOVED***
	flag := fs.Lookup(name)
	if flag == nil ***REMOVED***
		return false
	***REMOVED***
	return flag.NoOptDefVal != ""
***REMOVED***

func shortHasNoOptDefVal(name string, fs *flag.FlagSet) bool ***REMOVED***
	if len(name) == 0 ***REMOVED***
		return false
	***REMOVED***

	flag := fs.ShorthandLookup(name[:1])
	if flag == nil ***REMOVED***
		return false
	***REMOVED***
	return flag.NoOptDefVal != ""
***REMOVED***

func stripFlags(args []string, c *Command) []string ***REMOVED***
	if len(args) == 0 ***REMOVED***
		return args
	***REMOVED***
	c.mergePersistentFlags()

	commands := []string***REMOVED******REMOVED***
	flags := c.Flags()

Loop:
	for len(args) > 0 ***REMOVED***
		s := args[0]
		args = args[1:]
		switch ***REMOVED***
		case s == "--":
			// "--" terminates the flags
			break Loop
		case strings.HasPrefix(s, "--") && !strings.Contains(s, "=") && !hasNoOptDefVal(s[2:], flags):
			// If '--flag arg' then
			// delete arg from args.
			fallthrough // (do the same as below)
		case strings.HasPrefix(s, "-") && !strings.Contains(s, "=") && len(s) == 2 && !shortHasNoOptDefVal(s[1:], flags):
			// If '-f arg' then
			// delete 'arg' from args or break the loop if len(args) <= 1.
			if len(args) <= 1 ***REMOVED***
				break Loop
			***REMOVED*** else ***REMOVED***
				args = args[1:]
				continue
			***REMOVED***
		case s != "" && !strings.HasPrefix(s, "-"):
			commands = append(commands, s)
		***REMOVED***
	***REMOVED***

	return commands
***REMOVED***

// argsMinusFirstX removes only the first x from args.  Otherwise, commands that look like
// openshift admin policy add-role-to-user admin my-user, lose the admin argument (arg[4]).
func argsMinusFirstX(args []string, x string) []string ***REMOVED***
	for i, y := range args ***REMOVED***
		if x == y ***REMOVED***
			ret := []string***REMOVED******REMOVED***
			ret = append(ret, args[:i]...)
			ret = append(ret, args[i+1:]...)
			return ret
		***REMOVED***
	***REMOVED***
	return args
***REMOVED***

func isFlagArg(arg string) bool ***REMOVED***
	return ((len(arg) >= 3 && arg[1] == '-') ||
		(len(arg) >= 2 && arg[0] == '-' && arg[1] != '-'))
***REMOVED***

// Find the target command given the args and command tree
// Meant to be run on the highest node. Only searches down.
func (c *Command) Find(args []string) (*Command, []string, error) ***REMOVED***
	var innerfind func(*Command, []string) (*Command, []string)

	innerfind = func(c *Command, innerArgs []string) (*Command, []string) ***REMOVED***
		argsWOflags := stripFlags(innerArgs, c)
		if len(argsWOflags) == 0 ***REMOVED***
			return c, innerArgs
		***REMOVED***
		nextSubCmd := argsWOflags[0]

		cmd := c.findNext(nextSubCmd)
		if cmd != nil ***REMOVED***
			return innerfind(cmd, argsMinusFirstX(innerArgs, nextSubCmd))
		***REMOVED***
		return c, innerArgs
	***REMOVED***

	commandFound, a := innerfind(c, args)
	if commandFound.Args == nil ***REMOVED***
		return commandFound, a, legacyArgs(commandFound, stripFlags(a, commandFound))
	***REMOVED***
	return commandFound, a, nil
***REMOVED***

func (c *Command) findSuggestions(arg string) string ***REMOVED***
	if c.DisableSuggestions ***REMOVED***
		return ""
	***REMOVED***
	if c.SuggestionsMinimumDistance <= 0 ***REMOVED***
		c.SuggestionsMinimumDistance = 2
	***REMOVED***
	suggestionsString := ""
	if suggestions := c.SuggestionsFor(arg); len(suggestions) > 0 ***REMOVED***
		suggestionsString += "\n\nDid you mean this?\n"
		for _, s := range suggestions ***REMOVED***
			suggestionsString += fmt.Sprintf("\t%v\n", s)
		***REMOVED***
	***REMOVED***
	return suggestionsString
***REMOVED***

func (c *Command) findNext(next string) *Command ***REMOVED***
	matches := make([]*Command, 0)
	for _, cmd := range c.commands ***REMOVED***
		if cmd.Name() == next || cmd.HasAlias(next) ***REMOVED***
			cmd.commandCalledAs.name = next
			return cmd
		***REMOVED***
		if EnablePrefixMatching && cmd.hasNameOrAliasPrefix(next) ***REMOVED***
			matches = append(matches, cmd)
		***REMOVED***
	***REMOVED***

	if len(matches) == 1 ***REMOVED***
		return matches[0]
	***REMOVED***

	return nil
***REMOVED***

// Traverse the command tree to find the command, and parse args for
// each parent.
func (c *Command) Traverse(args []string) (*Command, []string, error) ***REMOVED***
	flags := []string***REMOVED******REMOVED***
	inFlag := false

	for i, arg := range args ***REMOVED***
		switch ***REMOVED***
		// A long flag with a space separated value
		case strings.HasPrefix(arg, "--") && !strings.Contains(arg, "="):
			// TODO: this isn't quite right, we should really check ahead for 'true' or 'false'
			inFlag = !hasNoOptDefVal(arg[2:], c.Flags())
			flags = append(flags, arg)
			continue
		// A short flag with a space separated value
		case strings.HasPrefix(arg, "-") && !strings.Contains(arg, "=") && len(arg) == 2 && !shortHasNoOptDefVal(arg[1:], c.Flags()):
			inFlag = true
			flags = append(flags, arg)
			continue
		// The value for a flag
		case inFlag:
			inFlag = false
			flags = append(flags, arg)
			continue
		// A flag without a value, or with an `=` separated value
		case isFlagArg(arg):
			flags = append(flags, arg)
			continue
		***REMOVED***

		cmd := c.findNext(arg)
		if cmd == nil ***REMOVED***
			return c, args, nil
		***REMOVED***

		if err := c.ParseFlags(flags); err != nil ***REMOVED***
			return nil, args, err
		***REMOVED***
		return cmd.Traverse(args[i+1:])
	***REMOVED***
	return c, args, nil
***REMOVED***

// SuggestionsFor provides suggestions for the typedName.
func (c *Command) SuggestionsFor(typedName string) []string ***REMOVED***
	suggestions := []string***REMOVED******REMOVED***
	for _, cmd := range c.commands ***REMOVED***
		if cmd.IsAvailableCommand() ***REMOVED***
			levenshteinDistance := ld(typedName, cmd.Name(), true)
			suggestByLevenshtein := levenshteinDistance <= c.SuggestionsMinimumDistance
			suggestByPrefix := strings.HasPrefix(strings.ToLower(cmd.Name()), strings.ToLower(typedName))
			if suggestByLevenshtein || suggestByPrefix ***REMOVED***
				suggestions = append(suggestions, cmd.Name())
			***REMOVED***
			for _, explicitSuggestion := range cmd.SuggestFor ***REMOVED***
				if strings.EqualFold(typedName, explicitSuggestion) ***REMOVED***
					suggestions = append(suggestions, cmd.Name())
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return suggestions
***REMOVED***

// VisitParents visits all parents of the command and invokes fn on each parent.
func (c *Command) VisitParents(fn func(*Command)) ***REMOVED***
	if c.HasParent() ***REMOVED***
		fn(c.Parent())
		c.Parent().VisitParents(fn)
	***REMOVED***
***REMOVED***

// Root finds root command.
func (c *Command) Root() *Command ***REMOVED***
	if c.HasParent() ***REMOVED***
		return c.Parent().Root()
	***REMOVED***
	return c
***REMOVED***

// ArgsLenAtDash will return the length of c.Flags().Args at the moment
// when a -- was found during args parsing.
func (c *Command) ArgsLenAtDash() int ***REMOVED***
	return c.Flags().ArgsLenAtDash()
***REMOVED***

func (c *Command) execute(a []string) (err error) ***REMOVED***
	if c == nil ***REMOVED***
		return fmt.Errorf("Called Execute() on a nil Command")
	***REMOVED***

	if len(c.Deprecated) > 0 ***REMOVED***
		c.Printf("Command %q is deprecated, %s\n", c.Name(), c.Deprecated)
	***REMOVED***

	// initialize help and version flag at the last point possible to allow for user
	// overriding
	c.InitDefaultHelpFlag()
	c.InitDefaultVersionFlag()

	err = c.ParseFlags(a)
	if err != nil ***REMOVED***
		return c.FlagErrorFunc()(c, err)
	***REMOVED***

	// If help is called, regardless of other flags, return we want help.
	// Also say we need help if the command isn't runnable.
	helpVal, err := c.Flags().GetBool("help")
	if err != nil ***REMOVED***
		// should be impossible to get here as we always declare a help
		// flag in InitDefaultHelpFlag()
		c.Println("\"help\" flag declared as non-bool. Please correct your code")
		return err
	***REMOVED***

	if helpVal ***REMOVED***
		return flag.ErrHelp
	***REMOVED***

	// for back-compat, only add version flag behavior if version is defined
	if c.Version != "" ***REMOVED***
		versionVal, err := c.Flags().GetBool("version")
		if err != nil ***REMOVED***
			c.Println("\"version\" flag declared as non-bool. Please correct your code")
			return err
		***REMOVED***
		if versionVal ***REMOVED***
			err := tmpl(c.OutOrStdout(), c.VersionTemplate(), c)
			if err != nil ***REMOVED***
				c.Println(err)
			***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	if !c.Runnable() ***REMOVED***
		return flag.ErrHelp
	***REMOVED***

	c.preRun()

	argWoFlags := c.Flags().Args()
	if c.DisableFlagParsing ***REMOVED***
		argWoFlags = a
	***REMOVED***

	if err := c.ValidateArgs(argWoFlags); err != nil ***REMOVED***
		return err
	***REMOVED***

	for p := c; p != nil; p = p.Parent() ***REMOVED***
		if p.PersistentPreRunE != nil ***REMOVED***
			if err := p.PersistentPreRunE(c, argWoFlags); err != nil ***REMOVED***
				return err
			***REMOVED***
			break
		***REMOVED*** else if p.PersistentPreRun != nil ***REMOVED***
			p.PersistentPreRun(c, argWoFlags)
			break
		***REMOVED***
	***REMOVED***
	if c.PreRunE != nil ***REMOVED***
		if err := c.PreRunE(c, argWoFlags); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED*** else if c.PreRun != nil ***REMOVED***
		c.PreRun(c, argWoFlags)
	***REMOVED***

	if err := c.validateRequiredFlags(); err != nil ***REMOVED***
		return err
	***REMOVED***
	if c.RunE != nil ***REMOVED***
		if err := c.RunE(c, argWoFlags); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		c.Run(c, argWoFlags)
	***REMOVED***
	if c.PostRunE != nil ***REMOVED***
		if err := c.PostRunE(c, argWoFlags); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED*** else if c.PostRun != nil ***REMOVED***
		c.PostRun(c, argWoFlags)
	***REMOVED***
	for p := c; p != nil; p = p.Parent() ***REMOVED***
		if p.PersistentPostRunE != nil ***REMOVED***
			if err := p.PersistentPostRunE(c, argWoFlags); err != nil ***REMOVED***
				return err
			***REMOVED***
			break
		***REMOVED*** else if p.PersistentPostRun != nil ***REMOVED***
			p.PersistentPostRun(c, argWoFlags)
			break
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (c *Command) preRun() ***REMOVED***
	for _, x := range initializers ***REMOVED***
		x()
	***REMOVED***
***REMOVED***

// ExecuteContext is the same as Execute(), but sets the ctx on the command.
// Retrieve ctx by calling cmd.Context() inside your *Run lifecycle or ValidArgs
// functions.
func (c *Command) ExecuteContext(ctx context.Context) error ***REMOVED***
	c.ctx = ctx
	return c.Execute()
***REMOVED***

// Execute uses the args (os.Args[1:] by default)
// and run through the command tree finding appropriate matches
// for commands and then corresponding flags.
func (c *Command) Execute() error ***REMOVED***
	_, err := c.ExecuteC()
	return err
***REMOVED***

// ExecuteContextC is the same as ExecuteC(), but sets the ctx on the command.
// Retrieve ctx by calling cmd.Context() inside your *Run lifecycle or ValidArgs
// functions.
func (c *Command) ExecuteContextC(ctx context.Context) (*Command, error) ***REMOVED***
	c.ctx = ctx
	return c.ExecuteC()
***REMOVED***

// ExecuteC executes the command.
func (c *Command) ExecuteC() (cmd *Command, err error) ***REMOVED***
	if c.ctx == nil ***REMOVED***
		c.ctx = context.Background()
	***REMOVED***

	// Regardless of what command execute is called on, run on Root only
	if c.HasParent() ***REMOVED***
		return c.Root().ExecuteC()
	***REMOVED***

	// windows hook
	if preExecHookFn != nil ***REMOVED***
		preExecHookFn(c)
	***REMOVED***

	// initialize help at the last point to allow for user overriding
	c.InitDefaultHelpCmd()
	// initialize completion at the last point to allow for user overriding
	c.initDefaultCompletionCmd()

	args := c.args

	// Workaround FAIL with "go test -v" or "cobra.test -test.v", see #155
	if c.args == nil && filepath.Base(os.Args[0]) != "cobra.test" ***REMOVED***
		args = os.Args[1:]
	***REMOVED***

	// initialize the hidden command to be used for shell completion
	c.initCompleteCmd(args)

	var flags []string
	if c.TraverseChildren ***REMOVED***
		cmd, flags, err = c.Traverse(args)
	***REMOVED*** else ***REMOVED***
		cmd, flags, err = c.Find(args)
	***REMOVED***
	if err != nil ***REMOVED***
		// If found parse to a subcommand and then failed, talk about the subcommand
		if cmd != nil ***REMOVED***
			c = cmd
		***REMOVED***
		if !c.SilenceErrors ***REMOVED***
			c.PrintErrln("Error:", err.Error())
			c.PrintErrf("Run '%v --help' for usage.\n", c.CommandPath())
		***REMOVED***
		return c, err
	***REMOVED***

	cmd.commandCalledAs.called = true
	if cmd.commandCalledAs.name == "" ***REMOVED***
		cmd.commandCalledAs.name = cmd.Name()
	***REMOVED***

	// We have to pass global context to children command
	// if context is present on the parent command.
	if cmd.ctx == nil ***REMOVED***
		cmd.ctx = c.ctx
	***REMOVED***

	err = cmd.execute(flags)
	if err != nil ***REMOVED***
		// Always show help if requested, even if SilenceErrors is in
		// effect
		if err == flag.ErrHelp ***REMOVED***
			cmd.HelpFunc()(cmd, args)
			return cmd, nil
		***REMOVED***

		// If root command has SilenceErrors flagged,
		// all subcommands should respect it
		if !cmd.SilenceErrors && !c.SilenceErrors ***REMOVED***
			c.PrintErrln("Error:", err.Error())
		***REMOVED***

		// If root command has SilenceUsage flagged,
		// all subcommands should respect it
		if !cmd.SilenceUsage && !c.SilenceUsage ***REMOVED***
			c.Println(cmd.UsageString())
		***REMOVED***
	***REMOVED***
	return cmd, err
***REMOVED***

func (c *Command) ValidateArgs(args []string) error ***REMOVED***
	if c.Args == nil ***REMOVED***
		return nil
	***REMOVED***
	return c.Args(c, args)
***REMOVED***

func (c *Command) validateRequiredFlags() error ***REMOVED***
	if c.DisableFlagParsing ***REMOVED***
		return nil
	***REMOVED***

	flags := c.Flags()
	missingFlagNames := []string***REMOVED******REMOVED***
	flags.VisitAll(func(pflag *flag.Flag) ***REMOVED***
		requiredAnnotation, found := pflag.Annotations[BashCompOneRequiredFlag]
		if !found ***REMOVED***
			return
		***REMOVED***
		if (requiredAnnotation[0] == "true") && !pflag.Changed ***REMOVED***
			missingFlagNames = append(missingFlagNames, pflag.Name)
		***REMOVED***
	***REMOVED***)

	if len(missingFlagNames) > 0 ***REMOVED***
		return fmt.Errorf(`required flag(s) "%s" not set`, strings.Join(missingFlagNames, `", "`))
	***REMOVED***
	return nil
***REMOVED***

// InitDefaultHelpFlag adds default help flag to c.
// It is called automatically by executing the c or by calling help and usage.
// If c already has help flag, it will do nothing.
func (c *Command) InitDefaultHelpFlag() ***REMOVED***
	c.mergePersistentFlags()
	if c.Flags().Lookup("help") == nil ***REMOVED***
		usage := "help for "
		if c.Name() == "" ***REMOVED***
			usage += "this command"
		***REMOVED*** else ***REMOVED***
			usage += c.Name()
		***REMOVED***
		c.Flags().BoolP("help", "h", false, usage)
	***REMOVED***
***REMOVED***

// InitDefaultVersionFlag adds default version flag to c.
// It is called automatically by executing the c.
// If c already has a version flag, it will do nothing.
// If c.Version is empty, it will do nothing.
func (c *Command) InitDefaultVersionFlag() ***REMOVED***
	if c.Version == "" ***REMOVED***
		return
	***REMOVED***

	c.mergePersistentFlags()
	if c.Flags().Lookup("version") == nil ***REMOVED***
		usage := "version for "
		if c.Name() == "" ***REMOVED***
			usage += "this command"
		***REMOVED*** else ***REMOVED***
			usage += c.Name()
		***REMOVED***
		if c.Flags().ShorthandLookup("v") == nil ***REMOVED***
			c.Flags().BoolP("version", "v", false, usage)
		***REMOVED*** else ***REMOVED***
			c.Flags().Bool("version", false, usage)
		***REMOVED***
	***REMOVED***
***REMOVED***

// InitDefaultHelpCmd adds default help command to c.
// It is called automatically by executing the c or by calling help and usage.
// If c already has help command or c has no subcommands, it will do nothing.
func (c *Command) InitDefaultHelpCmd() ***REMOVED***
	if !c.HasSubCommands() ***REMOVED***
		return
	***REMOVED***

	if c.helpCommand == nil ***REMOVED***
		c.helpCommand = &Command***REMOVED***
			Use:   "help [command]",
			Short: "Help about any command",
			Long: `Help provides help for any command in the application.
Simply type ` + c.Name() + ` help [path to command] for full details.`,
			ValidArgsFunction: func(c *Command, args []string, toComplete string) ([]string, ShellCompDirective) ***REMOVED***
				var completions []string
				cmd, _, e := c.Root().Find(args)
				if e != nil ***REMOVED***
					return nil, ShellCompDirectiveNoFileComp
				***REMOVED***
				if cmd == nil ***REMOVED***
					// Root help command.
					cmd = c.Root()
				***REMOVED***
				for _, subCmd := range cmd.Commands() ***REMOVED***
					if subCmd.IsAvailableCommand() || subCmd == cmd.helpCommand ***REMOVED***
						if strings.HasPrefix(subCmd.Name(), toComplete) ***REMOVED***
							completions = append(completions, fmt.Sprintf("%s\t%s", subCmd.Name(), subCmd.Short))
						***REMOVED***
					***REMOVED***
				***REMOVED***
				return completions, ShellCompDirectiveNoFileComp
			***REMOVED***,
			Run: func(c *Command, args []string) ***REMOVED***
				cmd, _, e := c.Root().Find(args)
				if cmd == nil || e != nil ***REMOVED***
					c.Printf("Unknown help topic %#q\n", args)
					CheckErr(c.Root().Usage())
				***REMOVED*** else ***REMOVED***
					cmd.InitDefaultHelpFlag() // make possible 'help' flag to be shown
					CheckErr(cmd.Help())
				***REMOVED***
			***REMOVED***,
		***REMOVED***
	***REMOVED***
	c.RemoveCommand(c.helpCommand)
	c.AddCommand(c.helpCommand)
***REMOVED***

// ResetCommands delete parent, subcommand and help command from c.
func (c *Command) ResetCommands() ***REMOVED***
	c.parent = nil
	c.commands = nil
	c.helpCommand = nil
	c.parentsPflags = nil
***REMOVED***

// Sorts commands by their names.
type commandSorterByName []*Command

func (c commandSorterByName) Len() int           ***REMOVED*** return len(c) ***REMOVED***
func (c commandSorterByName) Swap(i, j int)      ***REMOVED*** c[i], c[j] = c[j], c[i] ***REMOVED***
func (c commandSorterByName) Less(i, j int) bool ***REMOVED*** return c[i].Name() < c[j].Name() ***REMOVED***

// Commands returns a sorted slice of child commands.
func (c *Command) Commands() []*Command ***REMOVED***
	// do not sort commands if it already sorted or sorting was disabled
	if EnableCommandSorting && !c.commandsAreSorted ***REMOVED***
		sort.Sort(commandSorterByName(c.commands))
		c.commandsAreSorted = true
	***REMOVED***
	return c.commands
***REMOVED***

// AddCommand adds one or more commands to this parent command.
func (c *Command) AddCommand(cmds ...*Command) ***REMOVED***
	for i, x := range cmds ***REMOVED***
		if cmds[i] == c ***REMOVED***
			panic("Command can't be a child of itself")
		***REMOVED***
		cmds[i].parent = c
		// update max lengths
		usageLen := len(x.Use)
		if usageLen > c.commandsMaxUseLen ***REMOVED***
			c.commandsMaxUseLen = usageLen
		***REMOVED***
		commandPathLen := len(x.CommandPath())
		if commandPathLen > c.commandsMaxCommandPathLen ***REMOVED***
			c.commandsMaxCommandPathLen = commandPathLen
		***REMOVED***
		nameLen := len(x.Name())
		if nameLen > c.commandsMaxNameLen ***REMOVED***
			c.commandsMaxNameLen = nameLen
		***REMOVED***
		// If global normalization function exists, update all children
		if c.globNormFunc != nil ***REMOVED***
			x.SetGlobalNormalizationFunc(c.globNormFunc)
		***REMOVED***
		c.commands = append(c.commands, x)
		c.commandsAreSorted = false
	***REMOVED***
***REMOVED***

// RemoveCommand removes one or more commands from a parent command.
func (c *Command) RemoveCommand(cmds ...*Command) ***REMOVED***
	commands := []*Command***REMOVED******REMOVED***
main:
	for _, command := range c.commands ***REMOVED***
		for _, cmd := range cmds ***REMOVED***
			if command == cmd ***REMOVED***
				command.parent = nil
				continue main
			***REMOVED***
		***REMOVED***
		commands = append(commands, command)
	***REMOVED***
	c.commands = commands
	// recompute all lengths
	c.commandsMaxUseLen = 0
	c.commandsMaxCommandPathLen = 0
	c.commandsMaxNameLen = 0
	for _, command := range c.commands ***REMOVED***
		usageLen := len(command.Use)
		if usageLen > c.commandsMaxUseLen ***REMOVED***
			c.commandsMaxUseLen = usageLen
		***REMOVED***
		commandPathLen := len(command.CommandPath())
		if commandPathLen > c.commandsMaxCommandPathLen ***REMOVED***
			c.commandsMaxCommandPathLen = commandPathLen
		***REMOVED***
		nameLen := len(command.Name())
		if nameLen > c.commandsMaxNameLen ***REMOVED***
			c.commandsMaxNameLen = nameLen
		***REMOVED***
	***REMOVED***
***REMOVED***

// Print is a convenience method to Print to the defined output, fallback to Stderr if not set.
func (c *Command) Print(i ...interface***REMOVED******REMOVED***) ***REMOVED***
	fmt.Fprint(c.OutOrStderr(), i...)
***REMOVED***

// Println is a convenience method to Println to the defined output, fallback to Stderr if not set.
func (c *Command) Println(i ...interface***REMOVED******REMOVED***) ***REMOVED***
	c.Print(fmt.Sprintln(i...))
***REMOVED***

// Printf is a convenience method to Printf to the defined output, fallback to Stderr if not set.
func (c *Command) Printf(format string, i ...interface***REMOVED******REMOVED***) ***REMOVED***
	c.Print(fmt.Sprintf(format, i...))
***REMOVED***

// PrintErr is a convenience method to Print to the defined Err output, fallback to Stderr if not set.
func (c *Command) PrintErr(i ...interface***REMOVED******REMOVED***) ***REMOVED***
	fmt.Fprint(c.ErrOrStderr(), i...)
***REMOVED***

// PrintErrln is a convenience method to Println to the defined Err output, fallback to Stderr if not set.
func (c *Command) PrintErrln(i ...interface***REMOVED******REMOVED***) ***REMOVED***
	c.PrintErr(fmt.Sprintln(i...))
***REMOVED***

// PrintErrf is a convenience method to Printf to the defined Err output, fallback to Stderr if not set.
func (c *Command) PrintErrf(format string, i ...interface***REMOVED******REMOVED***) ***REMOVED***
	c.PrintErr(fmt.Sprintf(format, i...))
***REMOVED***

// CommandPath returns the full path to this command.
func (c *Command) CommandPath() string ***REMOVED***
	if c.HasParent() ***REMOVED***
		return c.Parent().CommandPath() + " " + c.Name()
	***REMOVED***
	return c.Name()
***REMOVED***

// UseLine puts out the full usage for a given command (including parents).
func (c *Command) UseLine() string ***REMOVED***
	var useline string
	if c.HasParent() ***REMOVED***
		useline = c.parent.CommandPath() + " " + c.Use
	***REMOVED*** else ***REMOVED***
		useline = c.Use
	***REMOVED***
	if c.DisableFlagsInUseLine ***REMOVED***
		return useline
	***REMOVED***
	if c.HasAvailableFlags() && !strings.Contains(useline, "[flags]") ***REMOVED***
		useline += " [flags]"
	***REMOVED***
	return useline
***REMOVED***

// DebugFlags used to determine which flags have been assigned to which commands
// and which persist.
func (c *Command) DebugFlags() ***REMOVED***
	c.Println("DebugFlags called on", c.Name())
	var debugflags func(*Command)

	debugflags = func(x *Command) ***REMOVED***
		if x.HasFlags() || x.HasPersistentFlags() ***REMOVED***
			c.Println(x.Name())
		***REMOVED***
		if x.HasFlags() ***REMOVED***
			x.flags.VisitAll(func(f *flag.Flag) ***REMOVED***
				if x.HasPersistentFlags() && x.persistentFlag(f.Name) != nil ***REMOVED***
					c.Println("  -"+f.Shorthand+",", "--"+f.Name, "["+f.DefValue+"]", "", f.Value, "  [LP]")
				***REMOVED*** else ***REMOVED***
					c.Println("  -"+f.Shorthand+",", "--"+f.Name, "["+f.DefValue+"]", "", f.Value, "  [L]")
				***REMOVED***
			***REMOVED***)
		***REMOVED***
		if x.HasPersistentFlags() ***REMOVED***
			x.pflags.VisitAll(func(f *flag.Flag) ***REMOVED***
				if x.HasFlags() ***REMOVED***
					if x.flags.Lookup(f.Name) == nil ***REMOVED***
						c.Println("  -"+f.Shorthand+",", "--"+f.Name, "["+f.DefValue+"]", "", f.Value, "  [P]")
					***REMOVED***
				***REMOVED*** else ***REMOVED***
					c.Println("  -"+f.Shorthand+",", "--"+f.Name, "["+f.DefValue+"]", "", f.Value, "  [P]")
				***REMOVED***
			***REMOVED***)
		***REMOVED***
		c.Println(x.flagErrorBuf)
		if x.HasSubCommands() ***REMOVED***
			for _, y := range x.commands ***REMOVED***
				debugflags(y)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	debugflags(c)
***REMOVED***

// Name returns the command's name: the first word in the use line.
func (c *Command) Name() string ***REMOVED***
	name := c.Use
	i := strings.Index(name, " ")
	if i >= 0 ***REMOVED***
		name = name[:i]
	***REMOVED***
	return name
***REMOVED***

// HasAlias determines if a given string is an alias of the command.
func (c *Command) HasAlias(s string) bool ***REMOVED***
	for _, a := range c.Aliases ***REMOVED***
		if a == s ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

// CalledAs returns the command name or alias that was used to invoke
// this command or an empty string if the command has not been called.
func (c *Command) CalledAs() string ***REMOVED***
	if c.commandCalledAs.called ***REMOVED***
		return c.commandCalledAs.name
	***REMOVED***
	return ""
***REMOVED***

// hasNameOrAliasPrefix returns true if the Name or any of aliases start
// with prefix
func (c *Command) hasNameOrAliasPrefix(prefix string) bool ***REMOVED***
	if strings.HasPrefix(c.Name(), prefix) ***REMOVED***
		c.commandCalledAs.name = c.Name()
		return true
	***REMOVED***
	for _, alias := range c.Aliases ***REMOVED***
		if strings.HasPrefix(alias, prefix) ***REMOVED***
			c.commandCalledAs.name = alias
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

// NameAndAliases returns a list of the command name and all aliases
func (c *Command) NameAndAliases() string ***REMOVED***
	return strings.Join(append([]string***REMOVED***c.Name()***REMOVED***, c.Aliases...), ", ")
***REMOVED***

// HasExample determines if the command has example.
func (c *Command) HasExample() bool ***REMOVED***
	return len(c.Example) > 0
***REMOVED***

// Runnable determines if the command is itself runnable.
func (c *Command) Runnable() bool ***REMOVED***
	return c.Run != nil || c.RunE != nil
***REMOVED***

// HasSubCommands determines if the command has children commands.
func (c *Command) HasSubCommands() bool ***REMOVED***
	return len(c.commands) > 0
***REMOVED***

// IsAvailableCommand determines if a command is available as a non-help command
// (this includes all non deprecated/hidden commands).
func (c *Command) IsAvailableCommand() bool ***REMOVED***
	if len(c.Deprecated) != 0 || c.Hidden ***REMOVED***
		return false
	***REMOVED***

	if c.HasParent() && c.Parent().helpCommand == c ***REMOVED***
		return false
	***REMOVED***

	if c.Runnable() || c.HasAvailableSubCommands() ***REMOVED***
		return true
	***REMOVED***

	return false
***REMOVED***

// IsAdditionalHelpTopicCommand determines if a command is an additional
// help topic command; additional help topic command is determined by the
// fact that it is NOT runnable/hidden/deprecated, and has no sub commands that
// are runnable/hidden/deprecated.
// Concrete example: https://github.com/spf13/cobra/issues/393#issuecomment-282741924.
func (c *Command) IsAdditionalHelpTopicCommand() bool ***REMOVED***
	// if a command is runnable, deprecated, or hidden it is not a 'help' command
	if c.Runnable() || len(c.Deprecated) != 0 || c.Hidden ***REMOVED***
		return false
	***REMOVED***

	// if any non-help sub commands are found, the command is not a 'help' command
	for _, sub := range c.commands ***REMOVED***
		if !sub.IsAdditionalHelpTopicCommand() ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***

	// the command either has no sub commands, or no non-help sub commands
	return true
***REMOVED***

// HasHelpSubCommands determines if a command has any available 'help' sub commands
// that need to be shown in the usage/help default template under 'additional help
// topics'.
func (c *Command) HasHelpSubCommands() bool ***REMOVED***
	// return true on the first found available 'help' sub command
	for _, sub := range c.commands ***REMOVED***
		if sub.IsAdditionalHelpTopicCommand() ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***

	// the command either has no sub commands, or no available 'help' sub commands
	return false
***REMOVED***

// HasAvailableSubCommands determines if a command has available sub commands that
// need to be shown in the usage/help default template under 'available commands'.
func (c *Command) HasAvailableSubCommands() bool ***REMOVED***
	// return true on the first found available (non deprecated/help/hidden)
	// sub command
	for _, sub := range c.commands ***REMOVED***
		if sub.IsAvailableCommand() ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***

	// the command either has no sub commands, or no available (non deprecated/help/hidden)
	// sub commands
	return false
***REMOVED***

// HasParent determines if the command is a child command.
func (c *Command) HasParent() bool ***REMOVED***
	return c.parent != nil
***REMOVED***

// GlobalNormalizationFunc returns the global normalization function or nil if it doesn't exist.
func (c *Command) GlobalNormalizationFunc() func(f *flag.FlagSet, name string) flag.NormalizedName ***REMOVED***
	return c.globNormFunc
***REMOVED***

// Flags returns the complete FlagSet that applies
// to this command (local and persistent declared here and by all parents).
func (c *Command) Flags() *flag.FlagSet ***REMOVED***
	if c.flags == nil ***REMOVED***
		c.flags = flag.NewFlagSet(c.Name(), flag.ContinueOnError)
		if c.flagErrorBuf == nil ***REMOVED***
			c.flagErrorBuf = new(bytes.Buffer)
		***REMOVED***
		c.flags.SetOutput(c.flagErrorBuf)
	***REMOVED***

	return c.flags
***REMOVED***

// LocalNonPersistentFlags are flags specific to this command which will NOT persist to subcommands.
func (c *Command) LocalNonPersistentFlags() *flag.FlagSet ***REMOVED***
	persistentFlags := c.PersistentFlags()

	out := flag.NewFlagSet(c.Name(), flag.ContinueOnError)
	c.LocalFlags().VisitAll(func(f *flag.Flag) ***REMOVED***
		if persistentFlags.Lookup(f.Name) == nil ***REMOVED***
			out.AddFlag(f)
		***REMOVED***
	***REMOVED***)
	return out
***REMOVED***

// LocalFlags returns the local FlagSet specifically set in the current command.
func (c *Command) LocalFlags() *flag.FlagSet ***REMOVED***
	c.mergePersistentFlags()

	if c.lflags == nil ***REMOVED***
		c.lflags = flag.NewFlagSet(c.Name(), flag.ContinueOnError)
		if c.flagErrorBuf == nil ***REMOVED***
			c.flagErrorBuf = new(bytes.Buffer)
		***REMOVED***
		c.lflags.SetOutput(c.flagErrorBuf)
	***REMOVED***
	c.lflags.SortFlags = c.Flags().SortFlags
	if c.globNormFunc != nil ***REMOVED***
		c.lflags.SetNormalizeFunc(c.globNormFunc)
	***REMOVED***

	addToLocal := func(f *flag.Flag) ***REMOVED***
		if c.lflags.Lookup(f.Name) == nil && c.parentsPflags.Lookup(f.Name) == nil ***REMOVED***
			c.lflags.AddFlag(f)
		***REMOVED***
	***REMOVED***
	c.Flags().VisitAll(addToLocal)
	c.PersistentFlags().VisitAll(addToLocal)
	return c.lflags
***REMOVED***

// InheritedFlags returns all flags which were inherited from parent commands.
func (c *Command) InheritedFlags() *flag.FlagSet ***REMOVED***
	c.mergePersistentFlags()

	if c.iflags == nil ***REMOVED***
		c.iflags = flag.NewFlagSet(c.Name(), flag.ContinueOnError)
		if c.flagErrorBuf == nil ***REMOVED***
			c.flagErrorBuf = new(bytes.Buffer)
		***REMOVED***
		c.iflags.SetOutput(c.flagErrorBuf)
	***REMOVED***

	local := c.LocalFlags()
	if c.globNormFunc != nil ***REMOVED***
		c.iflags.SetNormalizeFunc(c.globNormFunc)
	***REMOVED***

	c.parentsPflags.VisitAll(func(f *flag.Flag) ***REMOVED***
		if c.iflags.Lookup(f.Name) == nil && local.Lookup(f.Name) == nil ***REMOVED***
			c.iflags.AddFlag(f)
		***REMOVED***
	***REMOVED***)
	return c.iflags
***REMOVED***

// NonInheritedFlags returns all flags which were not inherited from parent commands.
func (c *Command) NonInheritedFlags() *flag.FlagSet ***REMOVED***
	return c.LocalFlags()
***REMOVED***

// PersistentFlags returns the persistent FlagSet specifically set in the current command.
func (c *Command) PersistentFlags() *flag.FlagSet ***REMOVED***
	if c.pflags == nil ***REMOVED***
		c.pflags = flag.NewFlagSet(c.Name(), flag.ContinueOnError)
		if c.flagErrorBuf == nil ***REMOVED***
			c.flagErrorBuf = new(bytes.Buffer)
		***REMOVED***
		c.pflags.SetOutput(c.flagErrorBuf)
	***REMOVED***
	return c.pflags
***REMOVED***

// ResetFlags deletes all flags from command.
func (c *Command) ResetFlags() ***REMOVED***
	c.flagErrorBuf = new(bytes.Buffer)
	c.flagErrorBuf.Reset()
	c.flags = flag.NewFlagSet(c.Name(), flag.ContinueOnError)
	c.flags.SetOutput(c.flagErrorBuf)
	c.pflags = flag.NewFlagSet(c.Name(), flag.ContinueOnError)
	c.pflags.SetOutput(c.flagErrorBuf)

	c.lflags = nil
	c.iflags = nil
	c.parentsPflags = nil
***REMOVED***

// HasFlags checks if the command contains any flags (local plus persistent from the entire structure).
func (c *Command) HasFlags() bool ***REMOVED***
	return c.Flags().HasFlags()
***REMOVED***

// HasPersistentFlags checks if the command contains persistent flags.
func (c *Command) HasPersistentFlags() bool ***REMOVED***
	return c.PersistentFlags().HasFlags()
***REMOVED***

// HasLocalFlags checks if the command has flags specifically declared locally.
func (c *Command) HasLocalFlags() bool ***REMOVED***
	return c.LocalFlags().HasFlags()
***REMOVED***

// HasInheritedFlags checks if the command has flags inherited from its parent command.
func (c *Command) HasInheritedFlags() bool ***REMOVED***
	return c.InheritedFlags().HasFlags()
***REMOVED***

// HasAvailableFlags checks if the command contains any flags (local plus persistent from the entire
// structure) which are not hidden or deprecated.
func (c *Command) HasAvailableFlags() bool ***REMOVED***
	return c.Flags().HasAvailableFlags()
***REMOVED***

// HasAvailablePersistentFlags checks if the command contains persistent flags which are not hidden or deprecated.
func (c *Command) HasAvailablePersistentFlags() bool ***REMOVED***
	return c.PersistentFlags().HasAvailableFlags()
***REMOVED***

// HasAvailableLocalFlags checks if the command has flags specifically declared locally which are not hidden
// or deprecated.
func (c *Command) HasAvailableLocalFlags() bool ***REMOVED***
	return c.LocalFlags().HasAvailableFlags()
***REMOVED***

// HasAvailableInheritedFlags checks if the command has flags inherited from its parent command which are
// not hidden or deprecated.
func (c *Command) HasAvailableInheritedFlags() bool ***REMOVED***
	return c.InheritedFlags().HasAvailableFlags()
***REMOVED***

// Flag climbs up the command tree looking for matching flag.
func (c *Command) Flag(name string) (flag *flag.Flag) ***REMOVED***
	flag = c.Flags().Lookup(name)

	if flag == nil ***REMOVED***
		flag = c.persistentFlag(name)
	***REMOVED***

	return
***REMOVED***

// Recursively find matching persistent flag.
func (c *Command) persistentFlag(name string) (flag *flag.Flag) ***REMOVED***
	if c.HasPersistentFlags() ***REMOVED***
		flag = c.PersistentFlags().Lookup(name)
	***REMOVED***

	if flag == nil ***REMOVED***
		c.updateParentsPflags()
		flag = c.parentsPflags.Lookup(name)
	***REMOVED***
	return
***REMOVED***

// ParseFlags parses persistent flag tree and local flags.
func (c *Command) ParseFlags(args []string) error ***REMOVED***
	if c.DisableFlagParsing ***REMOVED***
		return nil
	***REMOVED***

	if c.flagErrorBuf == nil ***REMOVED***
		c.flagErrorBuf = new(bytes.Buffer)
	***REMOVED***
	beforeErrorBufLen := c.flagErrorBuf.Len()
	c.mergePersistentFlags()

	// do it here after merging all flags and just before parse
	c.Flags().ParseErrorsWhitelist = flag.ParseErrorsWhitelist(c.FParseErrWhitelist)

	err := c.Flags().Parse(args)
	// Print warnings if they occurred (e.g. deprecated flag messages).
	if c.flagErrorBuf.Len()-beforeErrorBufLen > 0 && err == nil ***REMOVED***
		c.Print(c.flagErrorBuf.String())
	***REMOVED***

	return err
***REMOVED***

// Parent returns a commands parent command.
func (c *Command) Parent() *Command ***REMOVED***
	return c.parent
***REMOVED***

// mergePersistentFlags merges c.PersistentFlags() to c.Flags()
// and adds missing persistent flags of all parents.
func (c *Command) mergePersistentFlags() ***REMOVED***
	c.updateParentsPflags()
	c.Flags().AddFlagSet(c.PersistentFlags())
	c.Flags().AddFlagSet(c.parentsPflags)
***REMOVED***

// updateParentsPflags updates c.parentsPflags by adding
// new persistent flags of all parents.
// If c.parentsPflags == nil, it makes new.
func (c *Command) updateParentsPflags() ***REMOVED***
	if c.parentsPflags == nil ***REMOVED***
		c.parentsPflags = flag.NewFlagSet(c.Name(), flag.ContinueOnError)
		c.parentsPflags.SetOutput(c.flagErrorBuf)
		c.parentsPflags.SortFlags = false
	***REMOVED***

	if c.globNormFunc != nil ***REMOVED***
		c.parentsPflags.SetNormalizeFunc(c.globNormFunc)
	***REMOVED***

	c.Root().PersistentFlags().AddFlagSet(flag.CommandLine)

	c.VisitParents(func(parent *Command) ***REMOVED***
		c.parentsPflags.AddFlagSet(parent.PersistentFlags())
	***REMOVED***)
***REMOVED***
