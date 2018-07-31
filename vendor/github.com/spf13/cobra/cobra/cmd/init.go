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
	"path"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var initCmd = &cobra.Command***REMOVED***
	Use:     "init [name]",
	Aliases: []string***REMOVED***"initialize", "initialise", "create"***REMOVED***,
	Short:   "Initialize a Cobra Application",
	Long: `Initialize (cobra init) will create a new application, with a license
and the appropriate structure for a Cobra-based CLI application.

  * If a name is provided, it will be created in the current directory;
  * If no name is provided, the current directory will be assumed;
  * If a relative path is provided, it will be created inside $GOPATH
    (e.g. github.com/spf13/hugo);
  * If an absolute path is provided, it will be created;
  * If the directory already exists but is empty, it will be used.

Init will not use an existing directory with contents.`,

	Run: func(cmd *cobra.Command, args []string) ***REMOVED***
		wd, err := os.Getwd()
		if err != nil ***REMOVED***
			er(err)
		***REMOVED***

		var project *Project
		if len(args) == 0 ***REMOVED***
			project = NewProjectFromPath(wd)
		***REMOVED*** else if len(args) == 1 ***REMOVED***
			arg := args[0]
			if arg[0] == '.' ***REMOVED***
				arg = filepath.Join(wd, arg)
			***REMOVED***
			if filepath.IsAbs(arg) ***REMOVED***
				project = NewProjectFromPath(arg)
			***REMOVED*** else ***REMOVED***
				project = NewProject(arg)
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			er("please provide only one argument")
		***REMOVED***

		initializeProject(project)

		fmt.Fprintln(cmd.OutOrStdout(), `Your Cobra application is ready at
`+project.AbsPath()+`

Give it a try by going there and running `+"`go run main.go`."+`
Add commands to it by running `+"`cobra add [cmdname]`.")
	***REMOVED***,
***REMOVED***

func initializeProject(project *Project) ***REMOVED***
	if !exists(project.AbsPath()) ***REMOVED*** // If path doesn't yet exist, create it
		err := os.MkdirAll(project.AbsPath(), os.ModePerm)
		if err != nil ***REMOVED***
			er(err)
		***REMOVED***
	***REMOVED*** else if !isEmpty(project.AbsPath()) ***REMOVED*** // If path exists and is not empty don't use it
		er("Cobra will not create a new project in a non empty directory: " + project.AbsPath())
	***REMOVED***

	// We have a directory and it's empty. Time to initialize it.
	createLicenseFile(project.License(), project.AbsPath())
	createMainFile(project)
	createRootCmdFile(project)
***REMOVED***

func createLicenseFile(license License, path string) ***REMOVED***
	data := make(map[string]interface***REMOVED******REMOVED***)
	data["copyright"] = copyrightLine()

	// Generate license template from text and data.
	text, err := executeTemplate(license.Text, data)
	if err != nil ***REMOVED***
		er(err)
	***REMOVED***

	// Write license text to LICENSE file.
	err = writeStringToFile(filepath.Join(path, "LICENSE"), text)
	if err != nil ***REMOVED***
		er(err)
	***REMOVED***
***REMOVED***

func createMainFile(project *Project) ***REMOVED***
	mainTemplate := `***REMOVED******REMOVED*** comment .copyright ***REMOVED******REMOVED***
***REMOVED******REMOVED***if .license***REMOVED******REMOVED******REMOVED******REMOVED*** comment .license ***REMOVED******REMOVED******REMOVED******REMOVED***end***REMOVED******REMOVED***

package main

import "***REMOVED******REMOVED*** .importpath ***REMOVED******REMOVED***"

func main() ***REMOVED***
	cmd.Execute()
***REMOVED***
`
	data := make(map[string]interface***REMOVED******REMOVED***)
	data["copyright"] = copyrightLine()
	data["license"] = project.License().Header
	data["importpath"] = path.Join(project.Name(), filepath.Base(project.CmdPath()))

	mainScript, err := executeTemplate(mainTemplate, data)
	if err != nil ***REMOVED***
		er(err)
	***REMOVED***

	err = writeStringToFile(filepath.Join(project.AbsPath(), "main.go"), mainScript)
	if err != nil ***REMOVED***
		er(err)
	***REMOVED***
***REMOVED***

func createRootCmdFile(project *Project) ***REMOVED***
	template := `***REMOVED******REMOVED***comment .copyright***REMOVED******REMOVED***
***REMOVED******REMOVED***if .license***REMOVED******REMOVED******REMOVED******REMOVED***comment .license***REMOVED******REMOVED******REMOVED******REMOVED***end***REMOVED******REMOVED***

package cmd

import (
	"fmt"
	"os"
***REMOVED******REMOVED***if .viper***REMOVED******REMOVED***
	homedir "github.com/mitchellh/go-homedir"***REMOVED******REMOVED***end***REMOVED******REMOVED***
	"github.com/spf13/cobra"***REMOVED******REMOVED***if .viper***REMOVED******REMOVED***
	"github.com/spf13/viper"***REMOVED******REMOVED***end***REMOVED******REMOVED***
)***REMOVED******REMOVED***if .viper***REMOVED******REMOVED***

var cfgFile string***REMOVED******REMOVED***end***REMOVED******REMOVED***

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command***REMOVED***
	Use:   "***REMOVED******REMOVED***.appName***REMOVED******REMOVED***",
	Short: "A brief description of your application",
	Long: ` + "`" + `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.` + "`" + `,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) ***REMOVED*** ***REMOVED***,
***REMOVED***

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() ***REMOVED***
	if err := rootCmd.Execute(); err != nil ***REMOVED***
		fmt.Println(err)
		os.Exit(1)
	***REMOVED***
***REMOVED***

func init() ***REMOVED*** ***REMOVED******REMOVED***- if .viper***REMOVED******REMOVED***
	cobra.OnInitialize(initConfig)
***REMOVED******REMOVED***end***REMOVED******REMOVED***
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.***REMOVED******REMOVED*** if .viper ***REMOVED******REMOVED***
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.***REMOVED******REMOVED*** .appName ***REMOVED******REMOVED***.yaml)")***REMOVED******REMOVED*** else ***REMOVED******REMOVED***
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.***REMOVED******REMOVED*** .appName ***REMOVED******REMOVED***.yaml)")***REMOVED******REMOVED*** end ***REMOVED******REMOVED***

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
***REMOVED******REMOVED******REMOVED*** if .viper ***REMOVED******REMOVED***

// initConfig reads in config file and ENV variables if set.
func initConfig() ***REMOVED***
	if cfgFile != "" ***REMOVED***
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	***REMOVED*** else ***REMOVED***
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil ***REMOVED***
			fmt.Println(err)
			os.Exit(1)
		***REMOVED***

		// Search config in home directory with name ".***REMOVED******REMOVED*** .appName ***REMOVED******REMOVED***" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".***REMOVED******REMOVED*** .appName ***REMOVED******REMOVED***")
	***REMOVED***

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil ***REMOVED***
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	***REMOVED***
***REMOVED******REMOVED******REMOVED*** end ***REMOVED******REMOVED***
`

	data := make(map[string]interface***REMOVED******REMOVED***)
	data["copyright"] = copyrightLine()
	data["viper"] = viper.GetBool("useViper")
	data["license"] = project.License().Header
	data["appName"] = path.Base(project.Name())

	rootCmdScript, err := executeTemplate(template, data)
	if err != nil ***REMOVED***
		er(err)
	***REMOVED***

	err = writeStringToFile(filepath.Join(project.CmdPath(), "root.go"), rootCmdScript)
	if err != nil ***REMOVED***
		er(err)
	***REMOVED***

***REMOVED***
