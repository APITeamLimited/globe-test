package cobra

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/spf13/pflag"
)

// Annotations for Bash completion.
const (
	BashCompFilenameExt     = "cobra_annotation_bash_completion_filename_extensions"
	BashCompCustom          = "cobra_annotation_bash_completion_custom"
	BashCompOneRequiredFlag = "cobra_annotation_bash_completion_one_required_flag"
	BashCompSubdirsInDir    = "cobra_annotation_bash_completion_subdirs_in_dir"
)

func writePreamble(buf io.StringWriter, name string) ***REMOVED***
	WriteStringAndCheck(buf, fmt.Sprintf("# bash completion for %-36s -*- shell-script -*-\n", name))
	WriteStringAndCheck(buf, fmt.Sprintf(`
__%[1]s_debug()
***REMOVED***
    if [[ -n $***REMOVED***BASH_COMP_DEBUG_FILE***REMOVED*** ]]; then
        echo "$*" >> "$***REMOVED***BASH_COMP_DEBUG_FILE***REMOVED***"
    fi
***REMOVED***

# Homebrew on Macs have version 1.3 of bash-completion which doesn't include
# _init_completion. This is a very minimal version of that function.
__%[1]s_init_completion()
***REMOVED***
    COMPREPLY=()
    _get_comp_words_by_ref "$@" cur prev words cword
***REMOVED***

__%[1]s_index_of_word()
***REMOVED***
    local w word=$1
    shift
    index=0
    for w in "$@"; do
        [[ $w = "$word" ]] && return
        index=$((index+1))
    done
    index=-1
***REMOVED***

__%[1]s_contains_word()
***REMOVED***
    local w word=$1; shift
    for w in "$@"; do
        [[ $w = "$word" ]] && return
    done
    return 1
***REMOVED***

__%[1]s_handle_go_custom_completion()
***REMOVED***
    __%[1]s_debug "$***REMOVED***FUNCNAME[0]***REMOVED***: cur is $***REMOVED***cur***REMOVED***, words[*] is $***REMOVED***words[*]***REMOVED***, #words[@] is $***REMOVED***#words[@]***REMOVED***"

    local shellCompDirectiveError=%[3]d
    local shellCompDirectiveNoSpace=%[4]d
    local shellCompDirectiveNoFileComp=%[5]d
    local shellCompDirectiveFilterFileExt=%[6]d
    local shellCompDirectiveFilterDirs=%[7]d

    local out requestComp lastParam lastChar comp directive args

    # Prepare the command to request completions for the program.
    # Calling $***REMOVED***words[0]***REMOVED*** instead of directly %[1]s allows to handle aliases
    args=("$***REMOVED***words[@]:1***REMOVED***")
    requestComp="$***REMOVED***words[0]***REMOVED*** %[2]s $***REMOVED***args[*]***REMOVED***"

    lastParam=$***REMOVED***words[$(($***REMOVED***#words[@]***REMOVED***-1))]***REMOVED***
    lastChar=$***REMOVED***lastParam:$(($***REMOVED***#lastParam***REMOVED***-1)):1***REMOVED***
    __%[1]s_debug "$***REMOVED***FUNCNAME[0]***REMOVED***: lastParam $***REMOVED***lastParam***REMOVED***, lastChar $***REMOVED***lastChar***REMOVED***"

    if [ -z "$***REMOVED***cur***REMOVED***" ] && [ "$***REMOVED***lastChar***REMOVED***" != "=" ]; then
        # If the last parameter is complete (there is a space following it)
        # We add an extra empty parameter so we can indicate this to the go method.
        __%[1]s_debug "$***REMOVED***FUNCNAME[0]***REMOVED***: Adding extra empty parameter"
        requestComp="$***REMOVED***requestComp***REMOVED*** \"\""
    fi

    __%[1]s_debug "$***REMOVED***FUNCNAME[0]***REMOVED***: calling $***REMOVED***requestComp***REMOVED***"
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
    __%[1]s_debug "$***REMOVED***FUNCNAME[0]***REMOVED***: the completion directive is: $***REMOVED***directive***REMOVED***"
    __%[1]s_debug "$***REMOVED***FUNCNAME[0]***REMOVED***: the completions are: $***REMOVED***out[*]***REMOVED***"

    if [ $((directive & shellCompDirectiveError)) -ne 0 ]; then
        # Error code.  No completion.
        __%[1]s_debug "$***REMOVED***FUNCNAME[0]***REMOVED***: received error from custom completion go code"
        return
    else
        if [ $((directive & shellCompDirectiveNoSpace)) -ne 0 ]; then
            if [[ $(type -t compopt) = "builtin" ]]; then
                __%[1]s_debug "$***REMOVED***FUNCNAME[0]***REMOVED***: activating no space"
                compopt -o nospace
            fi
        fi
        if [ $((directive & shellCompDirectiveNoFileComp)) -ne 0 ]; then
            if [[ $(type -t compopt) = "builtin" ]]; then
                __%[1]s_debug "$***REMOVED***FUNCNAME[0]***REMOVED***: activating no file completion"
                compopt +o default
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
        local subDir
        # Use printf to strip any trailing newline
        subdir=$(printf "%%s" "$***REMOVED***out[0]***REMOVED***")
        if [ -n "$subdir" ]; then
            __%[1]s_debug "Listing directories in $subdir"
            __%[1]s_handle_subdirs_in_dir_flag "$subdir"
        else
            __%[1]s_debug "Listing directories in ."
            _filedir -d
        fi
    else
        while IFS='' read -r comp; do
            COMPREPLY+=("$comp")
        done < <(compgen -W "$***REMOVED***out[*]***REMOVED***" -- "$cur")
    fi
***REMOVED***

__%[1]s_handle_reply()
***REMOVED***
    __%[1]s_debug "$***REMOVED***FUNCNAME[0]***REMOVED***"
    local comp
    case $cur in
        -*)
            if [[ $(type -t compopt) = "builtin" ]]; then
                compopt -o nospace
            fi
            local allflags
            if [ $***REMOVED***#must_have_one_flag[@]***REMOVED*** -ne 0 ]; then
                allflags=("$***REMOVED***must_have_one_flag[@]***REMOVED***")
            else
                allflags=("$***REMOVED***flags[*]***REMOVED*** $***REMOVED***two_word_flags[*]***REMOVED***")
            fi
            while IFS='' read -r comp; do
                COMPREPLY+=("$comp")
            done < <(compgen -W "$***REMOVED***allflags[*]***REMOVED***" -- "$cur")
            if [[ $(type -t compopt) = "builtin" ]]; then
                [[ "$***REMOVED***COMPREPLY[0]***REMOVED***" == *= ]] || compopt +o nospace
            fi

            # complete after --flag=abc
            if [[ $cur == *=* ]]; then
                if [[ $(type -t compopt) = "builtin" ]]; then
                    compopt +o nospace
                fi

                local index flag
                flag="$***REMOVED***cur%%=****REMOVED***"
                __%[1]s_index_of_word "$***REMOVED***flag***REMOVED***" "$***REMOVED***flags_with_completion[@]***REMOVED***"
                COMPREPLY=()
                if [[ $***REMOVED***index***REMOVED*** -ge 0 ]]; then
                    PREFIX=""
                    cur="$***REMOVED***cur#*=***REMOVED***"
                    $***REMOVED***flags_completion[$***REMOVED***index***REMOVED***]***REMOVED***
                    if [ -n "$***REMOVED***ZSH_VERSION***REMOVED***" ]; then
                        # zsh completion needs --flag= prefix
                        eval "COMPREPLY=( \"\$***REMOVED***COMPREPLY[@]/#/$***REMOVED***flag***REMOVED***=***REMOVED***\" )"
                    fi
                fi
            fi
            return 0;
            ;;
    esac

    # check if we are handling a flag with special work handling
    local index
    __%[1]s_index_of_word "$***REMOVED***prev***REMOVED***" "$***REMOVED***flags_with_completion[@]***REMOVED***"
    if [[ $***REMOVED***index***REMOVED*** -ge 0 ]]; then
        $***REMOVED***flags_completion[$***REMOVED***index***REMOVED***]***REMOVED***
        return
    fi

    # we are parsing a flag and don't have a special handler, no completion
    if [[ $***REMOVED***cur***REMOVED*** != "$***REMOVED***words[cword]***REMOVED***" ]]; then
        return
    fi

    local completions
    completions=("$***REMOVED***commands[@]***REMOVED***")
    if [[ $***REMOVED***#must_have_one_noun[@]***REMOVED*** -ne 0 ]]; then
        completions+=("$***REMOVED***must_have_one_noun[@]***REMOVED***")
    elif [[ -n "$***REMOVED***has_completion_function***REMOVED***" ]]; then
        # if a go completion function is provided, defer to that function
        __%[1]s_handle_go_custom_completion
    fi
    if [[ $***REMOVED***#must_have_one_flag[@]***REMOVED*** -ne 0 ]]; then
        completions+=("$***REMOVED***must_have_one_flag[@]***REMOVED***")
    fi
    while IFS='' read -r comp; do
        COMPREPLY+=("$comp")
    done < <(compgen -W "$***REMOVED***completions[*]***REMOVED***" -- "$cur")

    if [[ $***REMOVED***#COMPREPLY[@]***REMOVED*** -eq 0 && $***REMOVED***#noun_aliases[@]***REMOVED*** -gt 0 && $***REMOVED***#must_have_one_noun[@]***REMOVED*** -ne 0 ]]; then
        while IFS='' read -r comp; do
            COMPREPLY+=("$comp")
        done < <(compgen -W "$***REMOVED***noun_aliases[*]***REMOVED***" -- "$cur")
    fi

    if [[ $***REMOVED***#COMPREPLY[@]***REMOVED*** -eq 0 ]]; then
		if declare -F __%[1]s_custom_func >/dev/null; then
			# try command name qualified custom func
			__%[1]s_custom_func
		else
			# otherwise fall back to unqualified for compatibility
			declare -F __custom_func >/dev/null && __custom_func
		fi
    fi

    # available in bash-completion >= 2, not always present on macOS
    if declare -F __ltrim_colon_completions >/dev/null; then
        __ltrim_colon_completions "$cur"
    fi

    # If there is only 1 completion and it is a flag with an = it will be completed
    # but we don't want a space after the =
    if [[ "$***REMOVED***#COMPREPLY[@]***REMOVED***" -eq "1" ]] && [[ $(type -t compopt) = "builtin" ]] && [[ "$***REMOVED***COMPREPLY[0]***REMOVED***" == --*= ]]; then
       compopt -o nospace
    fi
***REMOVED***

# The arguments should be in the form "ext1|ext2|extn"
__%[1]s_handle_filename_extension_flag()
***REMOVED***
    local ext="$1"
    _filedir "@($***REMOVED***ext***REMOVED***)"
***REMOVED***

__%[1]s_handle_subdirs_in_dir_flag()
***REMOVED***
    local dir="$1"
    pushd "$***REMOVED***dir***REMOVED***" >/dev/null 2>&1 && _filedir -d && popd >/dev/null 2>&1 || return
***REMOVED***

__%[1]s_handle_flag()
***REMOVED***
    __%[1]s_debug "$***REMOVED***FUNCNAME[0]***REMOVED***: c is $c words[c] is $***REMOVED***words[c]***REMOVED***"

    # if a command required a flag, and we found it, unset must_have_one_flag()
    local flagname=$***REMOVED***words[c]***REMOVED***
    local flagvalue
    # if the word contained an =
    if [[ $***REMOVED***words[c]***REMOVED*** == *"="* ]]; then
        flagvalue=$***REMOVED***flagname#*=***REMOVED*** # take in as flagvalue after the =
        flagname=$***REMOVED***flagname%%=****REMOVED*** # strip everything after the =
        flagname="$***REMOVED***flagname***REMOVED***=" # but put the = back
    fi
    __%[1]s_debug "$***REMOVED***FUNCNAME[0]***REMOVED***: looking for $***REMOVED***flagname***REMOVED***"
    if __%[1]s_contains_word "$***REMOVED***flagname***REMOVED***" "$***REMOVED***must_have_one_flag[@]***REMOVED***"; then
        must_have_one_flag=()
    fi

    # if you set a flag which only applies to this command, don't show subcommands
    if __%[1]s_contains_word "$***REMOVED***flagname***REMOVED***" "$***REMOVED***local_nonpersistent_flags[@]***REMOVED***"; then
      commands=()
    fi

    # keep flag value with flagname as flaghash
    # flaghash variable is an associative array which is only supported in bash > 3.
    if [[ -z "$***REMOVED***BASH_VERSION***REMOVED***" || "$***REMOVED***BASH_VERSINFO[0]***REMOVED***" -gt 3 ]]; then
        if [ -n "$***REMOVED***flagvalue***REMOVED***" ] ; then
            flaghash[$***REMOVED***flagname***REMOVED***]=$***REMOVED***flagvalue***REMOVED***
        elif [ -n "$***REMOVED***words[ $((c+1)) ]***REMOVED***" ] ; then
            flaghash[$***REMOVED***flagname***REMOVED***]=$***REMOVED***words[ $((c+1)) ]***REMOVED***
        else
            flaghash[$***REMOVED***flagname***REMOVED***]="true" # pad "true" for bool flag
        fi
    fi

    # skip the argument to a two word flag
    if [[ $***REMOVED***words[c]***REMOVED*** != *"="* ]] && __%[1]s_contains_word "$***REMOVED***words[c]***REMOVED***" "$***REMOVED***two_word_flags[@]***REMOVED***"; then
			  __%[1]s_debug "$***REMOVED***FUNCNAME[0]***REMOVED***: found a flag $***REMOVED***words[c]***REMOVED***, skip the next argument"
        c=$((c+1))
        # if we are looking for a flags value, don't show commands
        if [[ $c -eq $cword ]]; then
            commands=()
        fi
    fi

    c=$((c+1))

***REMOVED***

__%[1]s_handle_noun()
***REMOVED***
    __%[1]s_debug "$***REMOVED***FUNCNAME[0]***REMOVED***: c is $c words[c] is $***REMOVED***words[c]***REMOVED***"

    if __%[1]s_contains_word "$***REMOVED***words[c]***REMOVED***" "$***REMOVED***must_have_one_noun[@]***REMOVED***"; then
        must_have_one_noun=()
    elif __%[1]s_contains_word "$***REMOVED***words[c]***REMOVED***" "$***REMOVED***noun_aliases[@]***REMOVED***"; then
        must_have_one_noun=()
    fi

    nouns+=("$***REMOVED***words[c]***REMOVED***")
    c=$((c+1))
***REMOVED***

__%[1]s_handle_command()
***REMOVED***
    __%[1]s_debug "$***REMOVED***FUNCNAME[0]***REMOVED***: c is $c words[c] is $***REMOVED***words[c]***REMOVED***"

    local next_command
    if [[ -n $***REMOVED***last_command***REMOVED*** ]]; then
        next_command="_$***REMOVED***last_command***REMOVED***_$***REMOVED***words[c]//:/__***REMOVED***"
    else
        if [[ $c -eq 0 ]]; then
            next_command="_%[1]s_root_command"
        else
            next_command="_$***REMOVED***words[c]//:/__***REMOVED***"
        fi
    fi
    c=$((c+1))
    __%[1]s_debug "$***REMOVED***FUNCNAME[0]***REMOVED***: looking for $***REMOVED***next_command***REMOVED***"
    declare -F "$next_command" >/dev/null && $next_command
***REMOVED***

__%[1]s_handle_word()
***REMOVED***
    if [[ $c -ge $cword ]]; then
        __%[1]s_handle_reply
        return
    fi
    __%[1]s_debug "$***REMOVED***FUNCNAME[0]***REMOVED***: c is $c words[c] is $***REMOVED***words[c]***REMOVED***"
    if [[ "$***REMOVED***words[c]***REMOVED***" == -* ]]; then
        __%[1]s_handle_flag
    elif __%[1]s_contains_word "$***REMOVED***words[c]***REMOVED***" "$***REMOVED***commands[@]***REMOVED***"; then
        __%[1]s_handle_command
    elif [[ $c -eq 0 ]]; then
        __%[1]s_handle_command
    elif __%[1]s_contains_word "$***REMOVED***words[c]***REMOVED***" "$***REMOVED***command_aliases[@]***REMOVED***"; then
        # aliashash variable is an associative array which is only supported in bash > 3.
        if [[ -z "$***REMOVED***BASH_VERSION***REMOVED***" || "$***REMOVED***BASH_VERSINFO[0]***REMOVED***" -gt 3 ]]; then
            words[c]=$***REMOVED***aliashash[$***REMOVED***words[c]***REMOVED***]***REMOVED***
            __%[1]s_handle_command
        else
            __%[1]s_handle_noun
        fi
    else
        __%[1]s_handle_noun
    fi
    __%[1]s_handle_word
***REMOVED***

`, name, ShellCompNoDescRequestCmd,
		ShellCompDirectiveError, ShellCompDirectiveNoSpace, ShellCompDirectiveNoFileComp,
		ShellCompDirectiveFilterFileExt, ShellCompDirectiveFilterDirs))
***REMOVED***

func writePostscript(buf io.StringWriter, name string) ***REMOVED***
	name = strings.Replace(name, ":", "__", -1)
	WriteStringAndCheck(buf, fmt.Sprintf("__start_%s()\n", name))
	WriteStringAndCheck(buf, fmt.Sprintf(`***REMOVED***
    local cur prev words cword
    declare -A flaghash 2>/dev/null || :
    declare -A aliashash 2>/dev/null || :
    if declare -F _init_completion >/dev/null 2>&1; then
        _init_completion -s || return
    else
        __%[1]s_init_completion -n "=" || return
    fi

    local c=0
    local flags=()
    local two_word_flags=()
    local local_nonpersistent_flags=()
    local flags_with_completion=()
    local flags_completion=()
    local commands=("%[1]s")
    local must_have_one_flag=()
    local must_have_one_noun=()
    local has_completion_function
    local last_command
    local nouns=()

    __%[1]s_handle_word
***REMOVED***

`, name))
	WriteStringAndCheck(buf, fmt.Sprintf(`if [[ $(type -t compopt) = "builtin" ]]; then
    complete -o default -F __start_%s %s
else
    complete -o default -o nospace -F __start_%s %s
fi

`, name, name, name, name))
	WriteStringAndCheck(buf, "# ex: ts=4 sw=4 et filetype=sh\n")
***REMOVED***

func writeCommands(buf io.StringWriter, cmd *Command) ***REMOVED***
	WriteStringAndCheck(buf, "    commands=()\n")
	for _, c := range cmd.Commands() ***REMOVED***
		if !c.IsAvailableCommand() && c != cmd.helpCommand ***REMOVED***
			continue
		***REMOVED***
		WriteStringAndCheck(buf, fmt.Sprintf("    commands+=(%q)\n", c.Name()))
		writeCmdAliases(buf, c)
	***REMOVED***
	WriteStringAndCheck(buf, "\n")
***REMOVED***

func writeFlagHandler(buf io.StringWriter, name string, annotations map[string][]string, cmd *Command) ***REMOVED***
	for key, value := range annotations ***REMOVED***
		switch key ***REMOVED***
		case BashCompFilenameExt:
			WriteStringAndCheck(buf, fmt.Sprintf("    flags_with_completion+=(%q)\n", name))

			var ext string
			if len(value) > 0 ***REMOVED***
				ext = fmt.Sprintf("__%s_handle_filename_extension_flag ", cmd.Root().Name()) + strings.Join(value, "|")
			***REMOVED*** else ***REMOVED***
				ext = "_filedir"
			***REMOVED***
			WriteStringAndCheck(buf, fmt.Sprintf("    flags_completion+=(%q)\n", ext))
		case BashCompCustom:
			WriteStringAndCheck(buf, fmt.Sprintf("    flags_with_completion+=(%q)\n", name))

			if len(value) > 0 ***REMOVED***
				handlers := strings.Join(value, "; ")
				WriteStringAndCheck(buf, fmt.Sprintf("    flags_completion+=(%q)\n", handlers))
			***REMOVED*** else ***REMOVED***
				WriteStringAndCheck(buf, "    flags_completion+=(:)\n")
			***REMOVED***
		case BashCompSubdirsInDir:
			WriteStringAndCheck(buf, fmt.Sprintf("    flags_with_completion+=(%q)\n", name))

			var ext string
			if len(value) == 1 ***REMOVED***
				ext = fmt.Sprintf("__%s_handle_subdirs_in_dir_flag ", cmd.Root().Name()) + value[0]
			***REMOVED*** else ***REMOVED***
				ext = "_filedir -d"
			***REMOVED***
			WriteStringAndCheck(buf, fmt.Sprintf("    flags_completion+=(%q)\n", ext))
		***REMOVED***
	***REMOVED***
***REMOVED***

const cbn = "\")\n"

func writeShortFlag(buf io.StringWriter, flag *pflag.Flag, cmd *Command) ***REMOVED***
	name := flag.Shorthand
	format := "    "
	if len(flag.NoOptDefVal) == 0 ***REMOVED***
		format += "two_word_"
	***REMOVED***
	format += "flags+=(\"-%s" + cbn
	WriteStringAndCheck(buf, fmt.Sprintf(format, name))
	writeFlagHandler(buf, "-"+name, flag.Annotations, cmd)
***REMOVED***

func writeFlag(buf io.StringWriter, flag *pflag.Flag, cmd *Command) ***REMOVED***
	name := flag.Name
	format := "    flags+=(\"--%s"
	if len(flag.NoOptDefVal) == 0 ***REMOVED***
		format += "="
	***REMOVED***
	format += cbn
	WriteStringAndCheck(buf, fmt.Sprintf(format, name))
	if len(flag.NoOptDefVal) == 0 ***REMOVED***
		format = "    two_word_flags+=(\"--%s" + cbn
		WriteStringAndCheck(buf, fmt.Sprintf(format, name))
	***REMOVED***
	writeFlagHandler(buf, "--"+name, flag.Annotations, cmd)
***REMOVED***

func writeLocalNonPersistentFlag(buf io.StringWriter, flag *pflag.Flag) ***REMOVED***
	name := flag.Name
	format := "    local_nonpersistent_flags+=(\"--%[1]s" + cbn
	if len(flag.NoOptDefVal) == 0 ***REMOVED***
		format += "    local_nonpersistent_flags+=(\"--%[1]s=" + cbn
	***REMOVED***
	WriteStringAndCheck(buf, fmt.Sprintf(format, name))
	if len(flag.Shorthand) > 0 ***REMOVED***
		WriteStringAndCheck(buf, fmt.Sprintf("    local_nonpersistent_flags+=(\"-%s\")\n", flag.Shorthand))
	***REMOVED***
***REMOVED***

// Setup annotations for go completions for registered flags
func prepareCustomAnnotationsForFlags(cmd *Command) ***REMOVED***
	for flag := range flagCompletionFunctions ***REMOVED***
		// Make sure the completion script calls the __*_go_custom_completion function for
		// every registered flag.  We need to do this here (and not when the flag was registered
		// for completion) so that we can know the root command name for the prefix
		// of __<prefix>_go_custom_completion
		if flag.Annotations == nil ***REMOVED***
			flag.Annotations = map[string][]string***REMOVED******REMOVED***
		***REMOVED***
		flag.Annotations[BashCompCustom] = []string***REMOVED***fmt.Sprintf("__%[1]s_handle_go_custom_completion", cmd.Root().Name())***REMOVED***
	***REMOVED***
***REMOVED***

func writeFlags(buf io.StringWriter, cmd *Command) ***REMOVED***
	prepareCustomAnnotationsForFlags(cmd)
	WriteStringAndCheck(buf, `    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

`)
	localNonPersistentFlags := cmd.LocalNonPersistentFlags()
	cmd.NonInheritedFlags().VisitAll(func(flag *pflag.Flag) ***REMOVED***
		if nonCompletableFlag(flag) ***REMOVED***
			return
		***REMOVED***
		writeFlag(buf, flag, cmd)
		if len(flag.Shorthand) > 0 ***REMOVED***
			writeShortFlag(buf, flag, cmd)
		***REMOVED***
		// localNonPersistentFlags are used to stop the completion of subcommands when one is set
		// if TraverseChildren is true we should allow to complete subcommands
		if localNonPersistentFlags.Lookup(flag.Name) != nil && !cmd.Root().TraverseChildren ***REMOVED***
			writeLocalNonPersistentFlag(buf, flag)
		***REMOVED***
	***REMOVED***)
	cmd.InheritedFlags().VisitAll(func(flag *pflag.Flag) ***REMOVED***
		if nonCompletableFlag(flag) ***REMOVED***
			return
		***REMOVED***
		writeFlag(buf, flag, cmd)
		if len(flag.Shorthand) > 0 ***REMOVED***
			writeShortFlag(buf, flag, cmd)
		***REMOVED***
	***REMOVED***)

	WriteStringAndCheck(buf, "\n")
***REMOVED***

func writeRequiredFlag(buf io.StringWriter, cmd *Command) ***REMOVED***
	WriteStringAndCheck(buf, "    must_have_one_flag=()\n")
	flags := cmd.NonInheritedFlags()
	flags.VisitAll(func(flag *pflag.Flag) ***REMOVED***
		if nonCompletableFlag(flag) ***REMOVED***
			return
		***REMOVED***
		for key := range flag.Annotations ***REMOVED***
			switch key ***REMOVED***
			case BashCompOneRequiredFlag:
				format := "    must_have_one_flag+=(\"--%s"
				if flag.Value.Type() != "bool" ***REMOVED***
					format += "="
				***REMOVED***
				format += cbn
				WriteStringAndCheck(buf, fmt.Sprintf(format, flag.Name))

				if len(flag.Shorthand) > 0 ***REMOVED***
					WriteStringAndCheck(buf, fmt.Sprintf("    must_have_one_flag+=(\"-%s"+cbn, flag.Shorthand))
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***)
***REMOVED***

func writeRequiredNouns(buf io.StringWriter, cmd *Command) ***REMOVED***
	WriteStringAndCheck(buf, "    must_have_one_noun=()\n")
	sort.Strings(cmd.ValidArgs)
	for _, value := range cmd.ValidArgs ***REMOVED***
		// Remove any description that may be included following a tab character.
		// Descriptions are not supported by bash completion.
		value = strings.Split(value, "\t")[0]
		WriteStringAndCheck(buf, fmt.Sprintf("    must_have_one_noun+=(%q)\n", value))
	***REMOVED***
	if cmd.ValidArgsFunction != nil ***REMOVED***
		WriteStringAndCheck(buf, "    has_completion_function=1\n")
	***REMOVED***
***REMOVED***

func writeCmdAliases(buf io.StringWriter, cmd *Command) ***REMOVED***
	if len(cmd.Aliases) == 0 ***REMOVED***
		return
	***REMOVED***

	sort.Strings(cmd.Aliases)

	WriteStringAndCheck(buf, fmt.Sprint(`    if [[ -z "$***REMOVED***BASH_VERSION***REMOVED***" || "$***REMOVED***BASH_VERSINFO[0]***REMOVED***" -gt 3 ]]; then`, "\n"))
	for _, value := range cmd.Aliases ***REMOVED***
		WriteStringAndCheck(buf, fmt.Sprintf("        command_aliases+=(%q)\n", value))
		WriteStringAndCheck(buf, fmt.Sprintf("        aliashash[%q]=%q\n", value, cmd.Name()))
	***REMOVED***
	WriteStringAndCheck(buf, `    fi`)
	WriteStringAndCheck(buf, "\n")
***REMOVED***
func writeArgAliases(buf io.StringWriter, cmd *Command) ***REMOVED***
	WriteStringAndCheck(buf, "    noun_aliases=()\n")
	sort.Strings(cmd.ArgAliases)
	for _, value := range cmd.ArgAliases ***REMOVED***
		WriteStringAndCheck(buf, fmt.Sprintf("    noun_aliases+=(%q)\n", value))
	***REMOVED***
***REMOVED***

func gen(buf io.StringWriter, cmd *Command) ***REMOVED***
	for _, c := range cmd.Commands() ***REMOVED***
		if !c.IsAvailableCommand() && c != cmd.helpCommand ***REMOVED***
			continue
		***REMOVED***
		gen(buf, c)
	***REMOVED***
	commandName := cmd.CommandPath()
	commandName = strings.Replace(commandName, " ", "_", -1)
	commandName = strings.Replace(commandName, ":", "__", -1)

	if cmd.Root() == cmd ***REMOVED***
		WriteStringAndCheck(buf, fmt.Sprintf("_%s_root_command()\n***REMOVED***\n", commandName))
	***REMOVED*** else ***REMOVED***
		WriteStringAndCheck(buf, fmt.Sprintf("_%s()\n***REMOVED***\n", commandName))
	***REMOVED***

	WriteStringAndCheck(buf, fmt.Sprintf("    last_command=%q\n", commandName))
	WriteStringAndCheck(buf, "\n")
	WriteStringAndCheck(buf, "    command_aliases=()\n")
	WriteStringAndCheck(buf, "\n")

	writeCommands(buf, cmd)
	writeFlags(buf, cmd)
	writeRequiredFlag(buf, cmd)
	writeRequiredNouns(buf, cmd)
	writeArgAliases(buf, cmd)
	WriteStringAndCheck(buf, "***REMOVED***\n\n")
***REMOVED***

// GenBashCompletion generates bash completion file and writes to the passed writer.
func (c *Command) GenBashCompletion(w io.Writer) error ***REMOVED***
	buf := new(bytes.Buffer)
	writePreamble(buf, c.Name())
	if len(c.BashCompletionFunction) > 0 ***REMOVED***
		buf.WriteString(c.BashCompletionFunction + "\n")
	***REMOVED***
	gen(buf, c)
	writePostscript(buf, c.Name())

	_, err := buf.WriteTo(w)
	return err
***REMOVED***

func nonCompletableFlag(flag *pflag.Flag) bool ***REMOVED***
	return flag.Hidden || len(flag.Deprecated) > 0
***REMOVED***

// GenBashCompletionFile generates bash completion file.
func (c *Command) GenBashCompletionFile(filename string) error ***REMOVED***
	outFile, err := os.Create(filename)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer outFile.Close()

	return c.GenBashCompletion(outFile)
***REMOVED***
