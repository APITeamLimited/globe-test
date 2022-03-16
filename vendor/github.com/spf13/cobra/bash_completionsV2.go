package cobra

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

func (c *Command) genBashCompletion(w io.Writer, includeDesc bool) error ***REMOVED***
	buf := new(bytes.Buffer)
	genBashComp(buf, c.Name(), includeDesc)
	_, err := buf.WriteTo(w)
	return err
***REMOVED***

func genBashComp(buf io.StringWriter, name string, includeDesc bool) ***REMOVED***
	compCmd := ShellCompRequestCmd
	if !includeDesc ***REMOVED***
		compCmd = ShellCompNoDescRequestCmd
	***REMOVED***

	WriteStringAndCheck(buf, fmt.Sprintf(`# bash completion V2 for %-36[1]s -*- shell-script -*-

__%[1]s_debug()
***REMOVED***
    if [[ -n $***REMOVED***BASH_COMP_DEBUG_FILE:-***REMOVED*** ]]; then
        echo "$*" >> "$***REMOVED***BASH_COMP_DEBUG_FILE***REMOVED***"
    fi
***REMOVED***

# Macs have bash3 for which the bash-completion package doesn't include
# _init_completion. This is a minimal version of that function.
__%[1]s_init_completion()
***REMOVED***
    COMPREPLY=()
    _get_comp_words_by_ref "$@" cur prev words cword
***REMOVED***

# This function calls the %[1]s program to obtain the completion
# results and the directive.  It fills the 'out' and 'directive' vars.
__%[1]s_get_completion_results() ***REMOVED***
    local requestComp lastParam lastChar args

    # Prepare the command to request completions for the program.
    # Calling $***REMOVED***words[0]***REMOVED*** instead of directly %[1]s allows to handle aliases
    args=("$***REMOVED***words[@]:1***REMOVED***")
    requestComp="$***REMOVED***words[0]***REMOVED*** %[2]s $***REMOVED***args[*]***REMOVED***"

    lastParam=$***REMOVED***words[$(($***REMOVED***#words[@]***REMOVED***-1))]***REMOVED***
    lastChar=$***REMOVED***lastParam:$(($***REMOVED***#lastParam***REMOVED***-1)):1***REMOVED***
    __%[1]s_debug "lastParam $***REMOVED***lastParam***REMOVED***, lastChar $***REMOVED***lastChar***REMOVED***"

    if [ -z "$***REMOVED***cur***REMOVED***" ] && [ "$***REMOVED***lastChar***REMOVED***" != "=" ]; then
        # If the last parameter is complete (there is a space following it)
        # We add an extra empty parameter so we can indicate this to the go method.
        __%[1]s_debug "Adding extra empty parameter"
        requestComp="$***REMOVED***requestComp***REMOVED*** ''"
    fi

    # When completing a flag with an = (e.g., %[1]s -n=<TAB>)
    # bash focuses on the part after the =, so we need to remove
    # the flag part from $cur
    if [[ "$***REMOVED***cur***REMOVED***" == -*=* ]]; then
        cur="$***REMOVED***cur#*=***REMOVED***"
    fi

    __%[1]s_debug "Calling $***REMOVED***requestComp***REMOVED***"
    # Use eval to handle any environment variables and such
    out=$(eval "$***REMOVED***requestComp***REMOVED***" 2>/dev/null)

    # Extract the directive integer at the very end of the output following a colon (:)
    directive=$***REMOVED***out##*:***REMOVED***
    # Remove the directive
    out=$***REMOVED***out%%:****REMOVED***
    if [ "$***REMOVED***directive***REMOVED***" = "$***REMOVED***out***REMOVED***" ]; then
        # There is not directive specified
        directive=0
    fi
    __%[1]s_debug "The completion directive is: $***REMOVED***directive***REMOVED***"
    __%[1]s_debug "The completions are: $***REMOVED***out[*]***REMOVED***"
***REMOVED***

__%[1]s_process_completion_results() ***REMOVED***
    local shellCompDirectiveError=%[3]d
    local shellCompDirectiveNoSpace=%[4]d
    local shellCompDirectiveNoFileComp=%[5]d
    local shellCompDirectiveFilterFileExt=%[6]d
    local shellCompDirectiveFilterDirs=%[7]d

    if [ $((directive & shellCompDirectiveError)) -ne 0 ]; then
        # Error code.  No completion.
        __%[1]s_debug "Received error from custom completion go code"
        return
    else
        if [ $((directive & shellCompDirectiveNoSpace)) -ne 0 ]; then
            if [[ $(type -t compopt) = "builtin" ]]; then
                __%[1]s_debug "Activating no space"
                compopt -o nospace
            else
                __%[1]s_debug "No space directive not supported in this version of bash"
            fi
        fi
        if [ $((directive & shellCompDirectiveNoFileComp)) -ne 0 ]; then
            if [[ $(type -t compopt) = "builtin" ]]; then
                __%[1]s_debug "Activating no file completion"
                compopt +o default
            else
                __%[1]s_debug "No file completion directive not supported in this version of bash"
            fi
        fi
    fi

    if [ $((directive & shellCompDirectiveFilterFileExt)) -ne 0 ]; then
        # File extension filtering
        local fullFilter filter filteringCmd

        # Do not use quotes around the $out variable or else newline
        # characters will be kept.
        for filter in $***REMOVED***out[*]***REMOVED***; do
            fullFilter+="$filter|"
        done

        filteringCmd="_filedir $fullFilter"
        __%[1]s_debug "File filtering command: $filteringCmd"
        $filteringCmd
    elif [ $((directive & shellCompDirectiveFilterDirs)) -ne 0 ]; then
        # File completion for directories only

        # Use printf to strip any trailing newline
        local subdir
        subdir=$(printf "%%s" "$***REMOVED***out[0]***REMOVED***")
        if [ -n "$subdir" ]; then
            __%[1]s_debug "Listing directories in $subdir"
            pushd "$subdir" >/dev/null 2>&1 && _filedir -d && popd >/dev/null 2>&1 || return
        else
            __%[1]s_debug "Listing directories in ."
            _filedir -d
        fi
    else
        __%[1]s_handle_completion_types
    fi

    __%[1]s_handle_special_char "$cur" :
    __%[1]s_handle_special_char "$cur" =
***REMOVED***

__%[1]s_handle_completion_types() ***REMOVED***
    __%[1]s_debug "__%[1]s_handle_completion_types: COMP_TYPE is $COMP_TYPE"

    case $COMP_TYPE in
    37|42)
        # Type: menu-complete/menu-complete-backward and insert-completions
        # If the user requested inserting one completion at a time, or all
        # completions at once on the command-line we must remove the descriptions.
        # https://github.com/spf13/cobra/issues/1508
        local tab comp
        tab=$(printf '\t')
        while IFS='' read -r comp; do
            # Strip any description
            comp=$***REMOVED***comp%%%%$tab****REMOVED***
            # Only consider the completions that match
            comp=$(compgen -W "$comp" -- "$cur")
            if [ -n "$comp" ]; then
                COMPREPLY+=("$comp")
            fi
        done < <(printf "%%s\n" "$***REMOVED***out[@]***REMOVED***")
        ;;

    *)
        # Type: complete (normal completion)
        __%[1]s_handle_standard_completion_case
        ;;
    esac
***REMOVED***

__%[1]s_handle_standard_completion_case() ***REMOVED***
    local tab comp
    tab=$(printf '\t')

    local longest=0
    # Look for the longest completion so that we can format things nicely
    while IFS='' read -r comp; do
        # Strip any description before checking the length
        comp=$***REMOVED***comp%%%%$tab****REMOVED***
        # Only consider the completions that match
        comp=$(compgen -W "$comp" -- "$cur")
        if (($***REMOVED***#comp***REMOVED***>longest)); then
            longest=$***REMOVED***#comp***REMOVED***
        fi
    done < <(printf "%%s\n" "$***REMOVED***out[@]***REMOVED***")

    local completions=()
    while IFS='' read -r comp; do
        if [ -z "$comp" ]; then
            continue
        fi

        __%[1]s_debug "Original comp: $comp"
        comp="$(__%[1]s_format_comp_descriptions "$comp" "$longest")"
        __%[1]s_debug "Final comp: $comp"
        completions+=("$comp")
    done < <(printf "%%s\n" "$***REMOVED***out[@]***REMOVED***")

    while IFS='' read -r comp; do
        COMPREPLY+=("$comp")
    done < <(compgen -W "$***REMOVED***completions[*]***REMOVED***" -- "$cur")

    # If there is a single completion left, remove the description text
    if [ $***REMOVED***#COMPREPLY[*]***REMOVED*** -eq 1 ]; then
        __%[1]s_debug "COMPREPLY[0]: $***REMOVED***COMPREPLY[0]***REMOVED***"
        comp="$***REMOVED***COMPREPLY[0]%%%% ****REMOVED***"
        __%[1]s_debug "Removed description from single completion, which is now: $***REMOVED***comp***REMOVED***"
        COMPREPLY=()
        COMPREPLY+=("$comp")
    fi
***REMOVED***

__%[1]s_handle_special_char()
***REMOVED***
    local comp="$1"
    local char=$2
    if [[ "$comp" == *$***REMOVED***char***REMOVED**** && "$COMP_WORDBREAKS" == *$***REMOVED***char***REMOVED**** ]]; then
        local word=$***REMOVED***comp%%"$***REMOVED***comp##*$***REMOVED***char***REMOVED******REMOVED***"***REMOVED***
        local idx=$***REMOVED***#COMPREPLY[*]***REMOVED***
        while [[ $((--idx)) -ge 0 ]]; do
            COMPREPLY[$idx]=$***REMOVED***COMPREPLY[$idx]#"$word"***REMOVED***
        done
    fi
***REMOVED***

__%[1]s_format_comp_descriptions()
***REMOVED***
    local tab
    tab=$(printf '\t')
    local comp="$1"
    local longest=$2

    # Properly format the description string which follows a tab character if there is one
    if [[ "$comp" == *$tab* ]]; then
        desc=$***REMOVED***comp#*$tab***REMOVED***
        comp=$***REMOVED***comp%%%%$tab****REMOVED***

        # $COLUMNS stores the current shell width.
        # Remove an extra 4 because we add 2 spaces and 2 parentheses.
        maxdesclength=$(( COLUMNS - longest - 4 ))

        # Make sure we can fit a description of at least 8 characters
        # if we are to align the descriptions.
        if [[ $maxdesclength -gt 8 ]]; then
            # Add the proper number of spaces to align the descriptions
            for ((i = $***REMOVED***#comp***REMOVED*** ; i < longest ; i++)); do
                comp+=" "
            done
        else
            # Don't pad the descriptions so we can fit more text after the completion
            maxdesclength=$(( COLUMNS - $***REMOVED***#comp***REMOVED*** - 4 ))
        fi

        # If there is enough space for any description text,
        # truncate the descriptions that are too long for the shell width
        if [ $maxdesclength -gt 0 ]; then
            if [ $***REMOVED***#desc***REMOVED*** -gt $maxdesclength ]; then
                desc=$***REMOVED***desc:0:$(( maxdesclength - 1 ))***REMOVED***
                desc+="…"
            fi
            comp+="  ($desc)"
        fi
    fi

    # Must use printf to escape all special characters
    printf "%%q" "$***REMOVED***comp***REMOVED***"
***REMOVED***

__start_%[1]s()
***REMOVED***
    local cur prev words cword split

    COMPREPLY=()

    # Call _init_completion from the bash-completion package
    # to prepare the arguments properly
    if declare -F _init_completion >/dev/null 2>&1; then
        _init_completion -n "=:" || return
    else
        __%[1]s_init_completion -n "=:" || return
    fi

    __%[1]s_debug
    __%[1]s_debug "========= starting completion logic =========="
    __%[1]s_debug "cur is $***REMOVED***cur***REMOVED***, words[*] is $***REMOVED***words[*]***REMOVED***, #words[@] is $***REMOVED***#words[@]***REMOVED***, cword is $cword"

    # The user could have moved the cursor backwards on the command-line.
    # We need to trigger completion from the $cword location, so we need
    # to truncate the command-line ($words) up to the $cword location.
    words=("$***REMOVED***words[@]:0:$cword+1***REMOVED***")
    __%[1]s_debug "Truncated words[*]: $***REMOVED***words[*]***REMOVED***,"

    local out directive
    __%[1]s_get_completion_results
    __%[1]s_process_completion_results
***REMOVED***

if [[ $(type -t compopt) = "builtin" ]]; then
    complete -o default -F __start_%[1]s %[1]s
else
    complete -o default -o nospace -F __start_%[1]s %[1]s
fi

# ex: ts=4 sw=4 et filetype=sh
`, name, compCmd,
		ShellCompDirectiveError, ShellCompDirectiveNoSpace, ShellCompDirectiveNoFileComp,
		ShellCompDirectiveFilterFileExt, ShellCompDirectiveFilterDirs))
***REMOVED***

// GenBashCompletionFileV2 generates Bash completion version 2.
func (c *Command) GenBashCompletionFileV2(filename string, includeDesc bool) error ***REMOVED***
	outFile, err := os.Create(filename)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer outFile.Close()

	return c.GenBashCompletionV2(outFile, includeDesc)
***REMOVED***

// GenBashCompletionV2 generates Bash completion file version 2
// and writes it to the passed writer.
func (c *Command) GenBashCompletionV2(w io.Writer, includeDesc bool) error ***REMOVED***
	return c.genBashCompletion(w, includeDesc)
***REMOVED***
