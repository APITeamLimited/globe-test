# Generating Bash Completions For Your cobra.Command

Please refer to [Shell Completions](shell_completions.md) for details.

## Bash legacy dynamic completions

For backward compatibility, Cobra still supports its legacy dynamic completion solution (described below).  Unlike the `ValidArgsFunction` solution, the legacy solution will only work for Bash shell-completion and not for other shells. This legacy solution can be used along-side `ValidArgsFunction` and `RegisterFlagCompletionFunc()`, as long as both solutions are not used for the same command.  This provides a path to gradually migrate from the legacy solution to the new solution.

The legacy solution allows you to inject bash functions into the bash completion script.  Those bash functions are responsible for providing the completion choices for your own completions.

Some code that works in kubernetes:

```bash
const (
        bash_completion_func = `__kubectl_parse_get()
***REMOVED***
    local kubectl_output out
    if kubectl_output=$(kubectl get --no-headers "$1" 2>/dev/null); then
        out=($(echo "$***REMOVED***kubectl_output***REMOVED***" | awk '***REMOVED***print $1***REMOVED***'))
        COMPREPLY=( $( compgen -W "$***REMOVED***out[*]***REMOVED***" -- "$cur" ) )
    fi
***REMOVED***

__kubectl_get_resource()
***REMOVED***
    if [[ $***REMOVED***#nouns[@]***REMOVED*** -eq 0 ]]; then
        return 1
    fi
    __kubectl_parse_get $***REMOVED***nouns[$***REMOVED***#nouns[@]***REMOVED*** -1]***REMOVED***
    if [[ $? -eq 0 ]]; then
        return 0
    fi
***REMOVED***

__kubectl_custom_func() ***REMOVED***
    case $***REMOVED***last_command***REMOVED*** in
        kubectl_get | kubectl_describe | kubectl_delete | kubectl_stop)
            __kubectl_get_resource
            return
            ;;
        *)
            ;;
    esac
***REMOVED***
`)
```

And then I set that in my command definition:

```go
cmds := &cobra.Command***REMOVED***
	Use:   "kubectl",
	Short: "kubectl controls the Kubernetes cluster manager",
	Long: `kubectl controls the Kubernetes cluster manager.

Find more information at https://github.com/GoogleCloudPlatform/kubernetes.`,
	Run: runHelp,
	BashCompletionFunction: bash_completion_func,
***REMOVED***
```

The `BashCompletionFunction` option is really only valid/useful on the root command. Doing the above will cause `__kubectl_custom_func()` (`__<command-use>_custom_func()`) to be called when the built in processor was unable to find a solution. In the case of kubernetes a valid command might look something like `kubectl get pod [mypod]`. If you type `kubectl get pod [tab][tab]` the `__kubectl_customc_func()` will run because the cobra.Command only understood "kubectl" and "get." `__kubectl_custom_func()` will see that the cobra.Command is "kubectl_get" and will thus call another helper `__kubectl_get_resource()`.  `__kubectl_get_resource` will look at the 'nouns' collected. In our example the only noun will be `pod`.  So it will call `__kubectl_parse_get pod`.  `__kubectl_parse_get` will actually call out to kubernetes and get any pods.  It will then set `COMPREPLY` to valid pods!

Similarly, for flags:

```go
	annotation := make(map[string][]string)
	annotation[cobra.BashCompCustom] = []string***REMOVED***"__kubectl_get_namespaces"***REMOVED***

	flag := &pflag.Flag***REMOVED***
		Name:        "namespace",
		Usage:       usage,
		Annotations: annotation,
	***REMOVED***
	cmd.Flags().AddFlag(flag)
```

In addition add the `__kubectl_get_namespaces` implementation in the `BashCompletionFunction`
value, e.g.:

```bash
__kubectl_get_namespaces()
***REMOVED***
    local template
    template="***REMOVED******REMOVED*** range .items  ***REMOVED******REMOVED******REMOVED******REMOVED*** .metadata.name ***REMOVED******REMOVED*** ***REMOVED******REMOVED*** end ***REMOVED******REMOVED***"
    local kubectl_out
    if kubectl_out=$(kubectl get -o template --template="$***REMOVED***template***REMOVED***" namespace 2>/dev/null); then
        COMPREPLY=( $( compgen -W "$***REMOVED***kubectl_out***REMOVED***[*]" -- "$cur" ) )
    fi
***REMOVED***
```
