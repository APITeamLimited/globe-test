package cobra

import (
	"fmt"
)

type PositionalArgs func(cmd *Command, args []string) error

// Legacy arg validation has the following behaviour:
// - root commands with no subcommands can take arbitrary arguments
// - root commands with subcommands will do subcommand validity checking
// - subcommands will always accept arbitrary arguments
func legacyArgs(cmd *Command, args []string) error ***REMOVED***
	// no subcommand, always take args
	if !cmd.HasSubCommands() ***REMOVED***
		return nil
	***REMOVED***

	// root command with subcommands, do subcommand checking.
	if !cmd.HasParent() && len(args) > 0 ***REMOVED***
		return fmt.Errorf("unknown command %q for %q%s", args[0], cmd.CommandPath(), cmd.findSuggestions(args[0]))
	***REMOVED***
	return nil
***REMOVED***

// NoArgs returns an error if any args are included.
func NoArgs(cmd *Command, args []string) error ***REMOVED***
	if len(args) > 0 ***REMOVED***
		return fmt.Errorf("unknown command %q for %q", args[0], cmd.CommandPath())
	***REMOVED***
	return nil
***REMOVED***

// OnlyValidArgs returns an error if any args are not in the list of ValidArgs.
func OnlyValidArgs(cmd *Command, args []string) error ***REMOVED***
	if len(cmd.ValidArgs) > 0 ***REMOVED***
		for _, v := range args ***REMOVED***
			if !stringInSlice(v, cmd.ValidArgs) ***REMOVED***
				return fmt.Errorf("invalid argument %q for %q%s", v, cmd.CommandPath(), cmd.findSuggestions(args[0]))
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// ArbitraryArgs never returns an error.
func ArbitraryArgs(cmd *Command, args []string) error ***REMOVED***
	return nil
***REMOVED***

// MinimumNArgs returns an error if there is not at least N args.
func MinimumNArgs(n int) PositionalArgs ***REMOVED***
	return func(cmd *Command, args []string) error ***REMOVED***
		if len(args) < n ***REMOVED***
			return fmt.Errorf("requires at least %d arg(s), only received %d", n, len(args))
		***REMOVED***
		return nil
	***REMOVED***
***REMOVED***

// MaximumNArgs returns an error if there are more than N args.
func MaximumNArgs(n int) PositionalArgs ***REMOVED***
	return func(cmd *Command, args []string) error ***REMOVED***
		if len(args) > n ***REMOVED***
			return fmt.Errorf("accepts at most %d arg(s), received %d", n, len(args))
		***REMOVED***
		return nil
	***REMOVED***
***REMOVED***

// ExactArgs returns an error if there are not exactly n args.
func ExactArgs(n int) PositionalArgs ***REMOVED***
	return func(cmd *Command, args []string) error ***REMOVED***
		if len(args) != n ***REMOVED***
			return fmt.Errorf("accepts %d arg(s), received %d", n, len(args))
		***REMOVED***
		return nil
	***REMOVED***
***REMOVED***

// RangeArgs returns an error if the number of args is not within the expected range.
func RangeArgs(min int, max int) PositionalArgs ***REMOVED***
	return func(cmd *Command, args []string) error ***REMOVED***
		if len(args) < min || len(args) > max ***REMOVED***
			return fmt.Errorf("accepts between %d and %d arg(s), received %d", min, max, len(args))
		***REMOVED***
		return nil
	***REMOVED***
***REMOVED***
