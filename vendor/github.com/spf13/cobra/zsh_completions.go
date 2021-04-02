package cobra

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

// GenZshCompletionFile generates zsh completion file including descriptions.
func (c *Command) GenZshCompletionFile(filename string) error ***REMOVED***
	return c.genZshCompletionFile(filename, true)
***REMOVED***

// GenZshCompletion generates zsh completion file including descriptions
// and writes it to the passed writer.
func (c *Command) GenZshCompletion(w io.Writer) error ***REMOVED***
	return c.genZshCompletion(w, true)
***REMOVED***

// GenZshCompletionFileNoDesc generates zsh completion file without descriptions.
func (c *Command) GenZshCompletionFileNoDesc(filename string) error ***REMOVED***
	return c.genZshCompletionFile(filename, false)
***REMOVED***

// GenZshCompletionNoDesc generates zsh completion file without descriptions
// and writes it to the passed writer.
func (c *Command) GenZshCompletionNoDesc(w io.Writer) error ***REMOVED***
	return c.genZshCompletion(w, false)
***REMOVED***

// MarkZshCompPositionalArgumentFile only worked for zsh and its behavior was
// not consistent with Bash completion. It has therefore been disabled.
// Instead, when no other completion is specified, file completion is done by
// default for every argument. One can disable file completion on a per-argument
// basis by using ValidArgsFunction and ShellCompDirectiveNoFileComp.
// To achieve file extension filtering, one can use ValidArgsFunction and
// ShellCompDirectiveFilterFileExt.
//
// Deprecated
func (c *Command) MarkZshCompPositionalArgumentFile(argPosition int, patterns ...string) error ***REMOVED***
	return nil
***REMOVED***

// MarkZshCompPositionalArgumentWords only worked for zsh. It has therefore
// been disabled.
// To achieve the same behavior across all shells, one can use
// ValidArgs (for the first argument only) or ValidArgsFunction for
// any argument (can include the first one also).
//
// Deprecated
func (c *Command) MarkZshCompPositionalArgumentWords(argPosition int, words ...string) error ***REMOVED***
	return nil
***REMOVED***

func (c *Command) genZshCompletionFile(filename string, includeDesc bool) error ***REMOVED***
	outFile, err := os.Create(filename)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer outFile.Close()

	return c.genZshCompletion(outFile, includeDesc)
***REMOVED***

func (c *Command) genZshCompletion(w io.Writer, includeDesc bool) error ***REMOVED***
	buf := new(bytes.Buffer)
	genZshComp(buf, c.Name(), includeDesc)
	_, err := buf.WriteTo(w)
	return err
***REMOVED***

func genZshComp(buf io.StringWriter, name string, includeDesc bool) ***REMOVED***
	compCmd := ShellCompRequestCmd
	if !includeDesc ***REMOVED***
		compCmd = ShellCompNoDescRequestCmd
	***REMOVED***
	WriteStringAndCheck(buf, fmt.Sprintf(`#compdef _%[1]s %[1]s

# zsh completion for %-36[1]s -*- shell-script -*-

__%[1]s_debug()
***REMOVED***
    local file="$BASH_COMP_DEBUG_FILE"
    if [[ -n $***REMOVED***file***REMOVED*** ]]; then
        echo "$*" >> "$***REMOVED***file***REMOVED***"
    fi
***REMOVED***

_%[1]s()
***REMOVED***
    local shellCompDirectiveError=%[3]d
    local shellCompDirectiveNoSpace=%[4]d
    local shellCompDirectiveNoFileComp=%[5]d
    local shellCompDirectiveFilterFileExt=%[6]d
    local shellCompDirectiveFilterDirs=%[7]d

    local lastParam lastChar flagPrefix requestComp out directive compCount comp lastComp
    local -a completions

    __%[1]s_debug "\n========= starting completion logic =========="
    __%[1]s_debug "CURRENT: $***REMOVED***CURRENT***REMOVED***, words[*]: $***REMOVED***words[*]***REMOVED***"

    # The user could have moved the cursor backwards on the command-line.
    # We need to trigger completion from the $CURRENT location, so we need
    # to truncate the command-line ($words) up to the $CURRENT location.
    # (We cannot use $CURSOR as its value does not work when a command is an alias.)
    words=("$***REMOVED***=words[1,CURRENT]***REMOVED***")
    __%[1]s_debug "Truncated words[*]: $***REMOVED***words[*]***REMOVED***,"

    lastParam=$***REMOVED***words[-1]***REMOVED***
    lastChar=$***REMOVED***lastParam[-1]***REMOVED***
    __%[1]s_debug "lastParam: $***REMOVED***lastParam***REMOVED***, lastChar: $***REMOVED***lastChar***REMOVED***"

    # For zsh, when completing a flag with an = (e.g., %[1]s -n=<TAB>)
    # completions must be prefixed with the flag
    setopt local_options BASH_REMATCH
    if [[ "$***REMOVED***lastParam***REMOVED***" =~ '-.*=' ]]; then
        # We are dealing with a flag with an =
        flagPrefix="-P $***REMOVED***BASH_REMATCH***REMOVED***"
    fi

    # Prepare the command to obtain completions
    requestComp="$***REMOVED***words[1]***REMOVED*** %[2]s $***REMOVED***words[2,-1]***REMOVED***"
    if [ "$***REMOVED***lastChar***REMOVED***" = "" ]; then
        # If the last parameter is complete (there is a space following it)
        # We add an extra empty parameter so we can indicate this to the go completion code.
        __%[1]s_debug "Adding extra empty parameter"
        requestComp="$***REMOVED***requestComp***REMOVED*** \"\""
    fi

    __%[1]s_debug "About to call: eval $***REMOVED***requestComp***REMOVED***"

    # Use eval to handle any environment variables and such
    out=$(eval $***REMOVED***requestComp***REMOVED*** 2>/dev/null)
    __%[1]s_debug "completion output: $***REMOVED***out***REMOVED***"

    # Extract the directive integer following a : from the last line
    local lastLine
    while IFS='\n' read -r line; do
        lastLine=$***REMOVED***line***REMOVED***
    done < <(printf "%%s\n" "$***REMOVED***out[@]***REMOVED***")
    __%[1]s_debug "last line: $***REMOVED***lastLine***REMOVED***"

    if [ "$***REMOVED***lastLine[1]***REMOVED***" = : ]; then
        directive=$***REMOVED***lastLine[2,-1]***REMOVED***
        # Remove the directive including the : and the newline
        local suffix
        (( suffix=$***REMOVED***#lastLine***REMOVED***+2))
        out=$***REMOVED***out[1,-$suffix]***REMOVED***
    else
        # There is no directive specified.  Leave $out as is.
        __%[1]s_debug "No directive found.  Setting do default"
        directive=0
    fi

    __%[1]s_debug "directive: $***REMOVED***directive***REMOVED***"
    __%[1]s_debug "completions: $***REMOVED***out***REMOVED***"
    __%[1]s_debug "flagPrefix: $***REMOVED***flagPrefix***REMOVED***"

    if [ $((directive & shellCompDirectiveError)) -ne 0 ]; then
        __%[1]s_debug "Completion received error. Ignoring completions."
        return
    fi

    compCount=0
    while IFS='\n' read -r comp; do
        if [ -n "$comp" ]; then
            # If requested, completions are returned with a description.
            # The description is preceded by a TAB character.
            # For zsh's _describe, we need to use a : instead of a TAB.
            # We first need to escape any : as part of the completion itself.
            comp=$***REMOVED***comp//:/\\:***REMOVED***

            local tab=$(printf '\t')
            comp=$***REMOVED***comp//$tab/:***REMOVED***

            ((compCount++))
            __%[1]s_debug "Adding completion: $***REMOVED***comp***REMOVED***"
            completions+=$***REMOVED***comp***REMOVED***
            lastComp=$comp
        fi
    done < <(printf "%%s\n" "$***REMOVED***out[@]***REMOVED***")

    if [ $((directive & shellCompDirectiveFilterFileExt)) -ne 0 ]; then
        # File extension filtering
        local filteringCmd
        filteringCmd='_files'
        for filter in $***REMOVED***completions[@]***REMOVED***; do
            if [ $***REMOVED***filter[1]***REMOVED*** != '*' ]; then
                # zsh requires a glob pattern to do file filtering
                filter="\*.$filter"
            fi
            filteringCmd+=" -g $filter"
        done
        filteringCmd+=" $***REMOVED***flagPrefix***REMOVED***"

        __%[1]s_debug "File filtering command: $filteringCmd"
        _arguments '*:filename:'"$filteringCmd"
    elif [ $((directive & shellCompDirectiveFilterDirs)) -ne 0 ]; then
        # File completion for directories only
        local subDir
        subdir="$***REMOVED***completions[1]***REMOVED***"
        if [ -n "$subdir" ]; then
            __%[1]s_debug "Listing directories in $subdir"
            pushd "$***REMOVED***subdir***REMOVED***" >/dev/null 2>&1
        else
            __%[1]s_debug "Listing directories in ."
        fi

        _arguments '*:dirname:_files -/'" $***REMOVED***flagPrefix***REMOVED***"
        if [ -n "$subdir" ]; then
            popd >/dev/null 2>&1
        fi
    elif [ $((directive & shellCompDirectiveNoSpace)) -ne 0 ] && [ $***REMOVED***compCount***REMOVED*** -eq 1 ]; then
        __%[1]s_debug "Activating nospace."
        # We can use compadd here as there is no description when
        # there is only one completion.
        compadd -S '' "$***REMOVED***lastComp***REMOVED***"
    elif [ $***REMOVED***compCount***REMOVED*** -eq 0 ]; then
        if [ $((directive & shellCompDirectiveNoFileComp)) -ne 0 ]; then
            __%[1]s_debug "deactivating file completion"
        else
            # Perform file completion
            __%[1]s_debug "activating file completion"
            _arguments '*:filename:_files'" $***REMOVED***flagPrefix***REMOVED***"
        fi
    else
        _describe "completions" completions $(echo $flagPrefix)
    fi
***REMOVED***

# don't run the completion function when being source-ed or eval-ed
if [ "$funcstack[1]" = "_%[1]s" ]; then
	_%[1]s
fi
`, name, compCmd,
		ShellCompDirectiveError, ShellCompDirectiveNoSpace, ShellCompDirectiveNoFileComp,
		ShellCompDirectiveFilterFileExt, ShellCompDirectiveFilterDirs))
***REMOVED***
