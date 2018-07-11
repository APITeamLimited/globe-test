// Copyright © 2015 Steve Francia <spf@spf13.com>.
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

package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"unicode"

	"github.com/spf13/cobra"
)

func init() ***REMOVED***
	addCmd.Flags().StringVarP(&packageName, "package", "t", "", "target package name (e.g. github.com/spf13/hugo)")
	addCmd.Flags().StringVarP(&parentName, "parent", "p", "rootCmd", "variable name of parent command for this command")
***REMOVED***

var packageName, parentName string

var addCmd = &cobra.Command***REMOVED***
	Use:     "add [command name]",
	Aliases: []string***REMOVED***"command"***REMOVED***,
	Short:   "Add a command to a Cobra Application",
	Long: `Add (cobra add) will create a new command, with a license and
the appropriate structure for a Cobra-based CLI application,
and register it to its parent (default rootCmd).

If you want your command to be public, pass in the command name
with an initial uppercase letter.

Example: cobra add server -> resulting in a new cmd/server.go`,

	Run: func(cmd *cobra.Command, args []string) ***REMOVED***
		if len(args) < 1 ***REMOVED***
			er("add needs a name for the command")
		***REMOVED***

		var project *Project
		if packageName != "" ***REMOVED***
			project = NewProject(packageName)
		***REMOVED*** else ***REMOVED***
			wd, err := os.Getwd()
			if err != nil ***REMOVED***
				er(err)
			***REMOVED***
			project = NewProjectFromPath(wd)
		***REMOVED***

		cmdName := validateCmdName(args[0])
		cmdPath := filepath.Join(project.CmdPath(), cmdName+".go")
		createCmdFile(project.License(), cmdPath, cmdName)

		fmt.Fprintln(cmd.OutOrStdout(), cmdName, "created at", cmdPath)
	***REMOVED***,
***REMOVED***

// validateCmdName returns source without any dashes and underscore.
// If there will be dash or underscore, next letter will be uppered.
// It supports only ASCII (1-byte character) strings.
// https://github.com/spf13/cobra/issues/269
func validateCmdName(source string) string ***REMOVED***
	i := 0
	l := len(source)
	// The output is initialized on demand, then first dash or underscore
	// occurs.
	var output string

	for i < l ***REMOVED***
		if source[i] == '-' || source[i] == '_' ***REMOVED***
			if output == "" ***REMOVED***
				output = source[:i]
			***REMOVED***

			// If it's last rune and it's dash or underscore,
			// don't add it output and break the loop.
			if i == l-1 ***REMOVED***
				break
			***REMOVED***

			// If next character is dash or underscore,
			// just skip the current character.
			if source[i+1] == '-' || source[i+1] == '_' ***REMOVED***
				i++
				continue
			***REMOVED***

			// If the current character is dash or underscore,
			// upper next letter and add to output.
			output += string(unicode.ToUpper(rune(source[i+1])))
			// We know, what source[i] is dash or underscore and source[i+1] is
			// uppered character, so make i = i+2.
			i += 2
			continue
		***REMOVED***

		// If the current character isn't dash or underscore,
		// just add it.
		if output != "" ***REMOVED***
			output += string(source[i])
		***REMOVED***
		i++
	***REMOVED***

	if output == "" ***REMOVED***
		return source // source is initially valid name.
	***REMOVED***
	return output
***REMOVED***

func createCmdFile(license License, path, cmdName string) ***REMOVED***
	template := `***REMOVED******REMOVED***comment .copyright***REMOVED******REMOVED***
***REMOVED******REMOVED***if .license***REMOVED******REMOVED******REMOVED******REMOVED***comment .license***REMOVED******REMOVED******REMOVED******REMOVED***end***REMOVED******REMOVED***

package ***REMOVED******REMOVED***.cmdPackage***REMOVED******REMOVED***

import (
	"fmt"

	"github.com/spf13/cobra"
)

// ***REMOVED******REMOVED***.cmdName***REMOVED******REMOVED***Cmd represents the ***REMOVED******REMOVED***.cmdName***REMOVED******REMOVED*** command
var ***REMOVED******REMOVED***.cmdName***REMOVED******REMOVED***Cmd = &cobra.Command***REMOVED***
	Use:   "***REMOVED******REMOVED***.cmdName***REMOVED******REMOVED***",
	Short: "A brief description of your command",
	Long: ` + "`" + `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.` + "`" + `,
	Run: func(cmd *cobra.Command, args []string) ***REMOVED***
		fmt.Println("***REMOVED******REMOVED***.cmdName***REMOVED******REMOVED*** called")
	***REMOVED***,
***REMOVED***

func init() ***REMOVED***
	***REMOVED******REMOVED***.parentName***REMOVED******REMOVED***.AddCommand(***REMOVED******REMOVED***.cmdName***REMOVED******REMOVED***Cmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// ***REMOVED******REMOVED***.cmdName***REMOVED******REMOVED***Cmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// ***REMOVED******REMOVED***.cmdName***REMOVED******REMOVED***Cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
***REMOVED***
`

	data := make(map[string]interface***REMOVED******REMOVED***)
	data["copyright"] = copyrightLine()
	data["license"] = license.Header
	data["cmdPackage"] = filepath.Base(filepath.Dir(path)) // last dir of path
	data["parentName"] = parentName
	data["cmdName"] = cmdName

	cmdScript, err := executeTemplate(template, data)
	if err != nil ***REMOVED***
		er(err)
	***REMOVED***
	err = writeStringToFile(path, cmdScript)
	if err != nil ***REMOVED***
		er(err)
	***REMOVED***
***REMOVED***
