# Generating Bash Completions For Your Own cobra.Command

If you are using the generator you can create a completion command by running

```bash
cobra add completion
```

Update the help text show how to install the bash_completion Linux show here [Kubectl docs show mac options](https://kubernetes.io/docs/tasks/tools/install-kubectl/#enabling-shell-autocompletion)

Writing the shell script to stdout allows the most flexible use.

```go
// completionCmd represents the completion command
var completionCmd = &cobra.Command***REMOVED***
	Use:   "completion",
	Short: "Generates bash completion scripts",
	Long: `To load completion run

. <(bitbucket completion)

To configure your bash shell to load completions for each session add to your bashrc

# ~/.bashrc or ~/.profile
. <(bitbucket completion)
`,
	Run: func(cmd *cobra.Command, args []string) ***REMOVED***
		rootCmd.GenBashCompletion(os.Stdout);
	***REMOVED***,
***REMOVED***
```

**Note:** The cobra generator may include messages printed to stdout for example if the config file is loaded, this will break the auto complete script


## Example from kubectl

Generating bash completions from a cobra command is incredibly easy. An actual program which does so for the kubernetes kubectl binary is as follows:

```go
package main

import (
	"io/ioutil"
	"os"

	"k8s.io/kubernetes/pkg/kubectl/cmd"
	"k8s.io/kubernetes/pkg/kubectl/cmd/util"
)

func main() ***REMOVED***
	kubectl := cmd.NewKubectlCommand(util.NewFactory(nil), os.Stdin, ioutil.Discard, ioutil.Discard)
	kubectl.GenBashCompletionFile("out.sh")
***REMOVED***
```

`out.sh` will get you completions of subcommands and flags. Copy it to `/etc/bash_completion.d/` as described [here](https://debian-administration.org/article/316/An_introduction_to_bash_completion_part_1) and reset your terminal to use autocompletion. If you make additional annotations to your code, you can get even more intelligent and flexible behavior.

## Creating your own custom functions

Some more actual code that works in kubernetes:

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

__custom_func() ***REMOVED***
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

The `BashCompletionFunction` option is really only valid/useful on the root command. Doing the above will cause `__custom_func()` to be called when the built in processor was unable to find a solution. In the case of kubernetes a valid command might look something like `kubectl get pod [mypod]`. If you type `kubectl get pod [tab][tab]` the `__customc_func()` will run because the cobra.Command only understood "kubectl" and "get." `__custom_func()` will see that the cobra.Command is "kubectl_get" and will thus call another helper `__kubectl_get_resource()`.  `__kubectl_get_resource` will look at the 'nouns' collected. In our example the only noun will be `pod`.  So it will call `__kubectl_parse_get pod`.  `__kubectl_parse_get` will actually call out to kubernetes and get any pods.  It will then set `COMPREPLY` to valid pods!

## Have the completions code complete your 'nouns'

In the above example "pod" was assumed to already be typed. But if you want `kubectl get [tab][tab]` to show a list of valid "nouns" you have to set them. Simplified code from `kubectl get` looks like:

```go
validArgs []string = ***REMOVED*** "pod", "node", "service", "replicationcontroller" ***REMOVED***

cmd := &cobra.Command***REMOVED***
	Use:     "get [(-o|--output=)json|yaml|template|...] (RESOURCE [NAME] | RESOURCE/NAME ...)",
	Short:   "Display one or many resources",
	Long:    get_long,
	Example: get_example,
	Run: func(cmd *cobra.Command, args []string) ***REMOVED***
		err := RunGet(f, out, cmd, args)
		util.CheckErr(err)
	***REMOVED***,
	ValidArgs: validArgs,
***REMOVED***
```

Notice we put the "ValidArgs" on the "get" subcommand. Doing so will give results like

```bash
# kubectl get [tab][tab]
node                 pod                    replicationcontroller  service
```

## Plural form and shortcuts for nouns

If your nouns have a number of aliases, you can define them alongside `ValidArgs` using `ArgAliases`:

```go
argAliases []string = ***REMOVED*** "pods", "nodes", "services", "svc", "replicationcontrollers", "rc" ***REMOVED***

cmd := &cobra.Command***REMOVED***
    ...
	ValidArgs:  validArgs,
	ArgAliases: argAliases
***REMOVED***
```

The aliases are not shown to the user on tab completion, but they are accepted as valid nouns by
the completion algorithm if entered manually, e.g. in:

```bash
# kubectl get rc [tab][tab]
backend        frontend       database 
```

Note that without declaring `rc` as an alias, the completion algorithm would show the list of nouns
in this example again instead of the replication controllers.

## Mark flags as required

Most of the time completions will only show subcommands. But if a flag is required to make a subcommand work, you probably want it to show up when the user types [tab][tab].  Marking a flag as 'Required' is incredibly easy.

```go
cmd.MarkFlagRequired("pod")
cmd.MarkFlagRequired("container")
```

and you'll get something like

```bash
# kubectl exec [tab][tab][tab]
-c            --container=  -p            --pod=  
```

# Specify valid filename extensions for flags that take a filename

In this example we use --filename= and expect to get a json or yaml file as the argument. To make this easier we annotate the --filename flag with valid filename extensions.

```go
	annotations := []string***REMOVED***"json", "yaml", "yml"***REMOVED***
	annotation := make(map[string][]string)
	annotation[cobra.BashCompFilenameExt] = annotations

	flag := &pflag.Flag***REMOVED***
		Name:        "filename",
		Shorthand:   "f",
		Usage:       usage,
		Value:       value,
		DefValue:    value.String(),
		Annotations: annotation,
	***REMOVED***
	cmd.Flags().AddFlag(flag)
```

Now when you run a command with this filename flag you'll get something like

```bash
# kubectl create -f 
test/                         example/                      rpmbuild/
hello.yml                     test.json
```

So while there are many other files in the CWD it only shows me subdirs and those with valid extensions.

# Specify custom flag completion

Similar to the filename completion and filtering using cobra.BashCompFilenameExt, you can specify
a custom flag completion function with cobra.BashCompCustom:

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

In addition add the `__handle_namespace_flag` implementation in the `BashCompletionFunction`
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
# Using bash aliases for commands

You can also configure the `bash aliases` for the commands and they will also support completions.

```bash
alias aliasname=origcommand
complete -o default -F __start_origcommand aliasname

# and now when you run `aliasname` completion will make
# suggestions as it did for `origcommand`.

$) aliasname <tab><tab>
completion     firstcommand   secondcommand
```
