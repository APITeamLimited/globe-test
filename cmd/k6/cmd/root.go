/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2016 Load Impact
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package cmd

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var Version = "0.17.1"
var Banner = `
          /\      |‾‾|  /‾‾/  /‾/   
     /\  /  \     |  |_/  /  / /   
    /  \/    \    |      |  /  ‾‾\  
   /          \   |  |‾\  \ | (_) | 
  / __________ \  |__|  \__\ \___/ .io`

var cfgFile string

// RootCmd represents the base command when called without any subcommands.
var RootCmd = &cobra.Command***REMOVED***
	Use:   "k6",
	Short: "a next-generation load generator",
	Long:  Banner,
	PersistentPreRun: func(cmd *cobra.Command, args []string) ***REMOVED***
		if viper.GetBool("verbose") ***REMOVED***
			log.SetLevel(log.DebugLevel)
		***REMOVED***
	***REMOVED***,
***REMOVED***

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() ***REMOVED***
	if err := RootCmd.Execute(); err != nil ***REMOVED***
		fmt.Println(err)
		os.Exit(-1)
	***REMOVED***
***REMOVED***

func init() ***REMOVED***
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().BoolP("verbose", "v", false, "enable debug logging")
	if err := viper.BindPFlags(RootCmd.PersistentFlags()); err != nil ***REMOVED***
		panic(err)
	***REMOVED***

	// It makes no sense to bind this to viper, so register it afterwards.
	RootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default $HOME/.k6.yaml)")
***REMOVED***

// initConfig reads in config file and ENV variables if set.
func initConfig() ***REMOVED***
	// Enable ability to specify config file via flag.
	if cfgFile != "" ***REMOVED***
		viper.SetConfigFile(cfgFile)
	***REMOVED***

	viper.SetConfigName(".k6")   // Name of config file (without extension).
	viper.AddConfigPath("$HOME") // Adding home directory as first search path.
	viper.SetEnvPrefix("k6")     // Accept environment variables starting with K6_.
	viper.AutomaticEnv()         // Read in environment variables that match.

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil ***REMOVED***
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok ***REMOVED***
			log.WithError(err).Error("Couldn't read global config")
		***REMOVED***
	***REMOVED***
***REMOVED***
