package cobra

import (
	"github.com/spf13/pflag"
)

// MarkFlagRequired instructs the various shell completion implementations to
// prioritize the named flag when performing completion,
// and causes your command to report an error if invoked without the flag.
func (c *Command) MarkFlagRequired(name string) error ***REMOVED***
	return MarkFlagRequired(c.Flags(), name)
***REMOVED***

// MarkPersistentFlagRequired instructs the various shell completion implementations to
// prioritize the named persistent flag when performing completion,
// and causes your command to report an error if invoked without the flag.
func (c *Command) MarkPersistentFlagRequired(name string) error ***REMOVED***
	return MarkFlagRequired(c.PersistentFlags(), name)
***REMOVED***

// MarkFlagRequired instructs the various shell completion implementations to
// prioritize the named flag when performing completion,
// and causes your command to report an error if invoked without the flag.
func MarkFlagRequired(flags *pflag.FlagSet, name string) error ***REMOVED***
	return flags.SetAnnotation(name, BashCompOneRequiredFlag, []string***REMOVED***"true"***REMOVED***)
***REMOVED***

// MarkFlagFilename instructs the various shell completion implementations to
// limit completions for the named flag to the specified file extensions.
func (c *Command) MarkFlagFilename(name string, extensions ...string) error ***REMOVED***
	return MarkFlagFilename(c.Flags(), name, extensions...)
***REMOVED***

// MarkFlagCustom adds the BashCompCustom annotation to the named flag, if it exists.
// The bash completion script will call the bash function f for the flag.
//
// This will only work for bash completion.
// It is recommended to instead use c.RegisterFlagCompletionFunc(...) which allows
// to register a Go function which will work across all shells.
func (c *Command) MarkFlagCustom(name string, f string) error ***REMOVED***
	return MarkFlagCustom(c.Flags(), name, f)
***REMOVED***

// MarkPersistentFlagFilename instructs the various shell completion
// implementations to limit completions for the named persistent flag to the
// specified file extensions.
func (c *Command) MarkPersistentFlagFilename(name string, extensions ...string) error ***REMOVED***
	return MarkFlagFilename(c.PersistentFlags(), name, extensions...)
***REMOVED***

// MarkFlagFilename instructs the various shell completion implementations to
// limit completions for the named flag to the specified file extensions.
func MarkFlagFilename(flags *pflag.FlagSet, name string, extensions ...string) error ***REMOVED***
	return flags.SetAnnotation(name, BashCompFilenameExt, extensions)
***REMOVED***

// MarkFlagCustom adds the BashCompCustom annotation to the named flag, if it exists.
// The bash completion script will call the bash function f for the flag.
//
// This will only work for bash completion.
// It is recommended to instead use c.RegisterFlagCompletionFunc(...) which allows
// to register a Go function which will work across all shells.
func MarkFlagCustom(flags *pflag.FlagSet, name string, f string) error ***REMOVED***
	return flags.SetAnnotation(name, BashCompCustom, []string***REMOVED***f***REMOVED***)
***REMOVED***

// MarkFlagDirname instructs the various shell completion implementations to
// limit completions for the named flag to directory names.
func (c *Command) MarkFlagDirname(name string) error ***REMOVED***
	return MarkFlagDirname(c.Flags(), name)
***REMOVED***

// MarkPersistentFlagDirname instructs the various shell completion
// implementations to limit completions for the named persistent flag to
// directory names.
func (c *Command) MarkPersistentFlagDirname(name string) error ***REMOVED***
	return MarkFlagDirname(c.PersistentFlags(), name)
***REMOVED***

// MarkFlagDirname instructs the various shell completion implementations to
// limit completions for the named flag to directory names.
func MarkFlagDirname(flags *pflag.FlagSet, name string) error ***REMOVED***
	return flags.SetAnnotation(name, BashCompSubdirsInDir, []string***REMOVED******REMOVED***)
***REMOVED***
