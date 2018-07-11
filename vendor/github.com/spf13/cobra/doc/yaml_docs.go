// Copyright 2016 French Ben. All rights reserved.
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

package doc

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v2"
)

type cmdOption struct ***REMOVED***
	Name         string
	Shorthand    string `yaml:",omitempty"`
	DefaultValue string `yaml:"default_value,omitempty"`
	Usage        string `yaml:",omitempty"`
***REMOVED***

type cmdDoc struct ***REMOVED***
	Name             string
	Synopsis         string      `yaml:",omitempty"`
	Description      string      `yaml:",omitempty"`
	Options          []cmdOption `yaml:",omitempty"`
	InheritedOptions []cmdOption `yaml:"inherited_options,omitempty"`
	Example          string      `yaml:",omitempty"`
	SeeAlso          []string    `yaml:"see_also,omitempty"`
***REMOVED***

// GenYamlTree creates yaml structured ref files for this command and all descendants
// in the directory given. This function may not work
// correctly if your command names have `-` in them. If you have `cmd` with two
// subcmds, `sub` and `sub-third`, and `sub` has a subcommand called `third`
// it is undefined which help output will be in the file `cmd-sub-third.1`.
func GenYamlTree(cmd *cobra.Command, dir string) error ***REMOVED***
	identity := func(s string) string ***REMOVED*** return s ***REMOVED***
	emptyStr := func(s string) string ***REMOVED*** return "" ***REMOVED***
	return GenYamlTreeCustom(cmd, dir, emptyStr, identity)
***REMOVED***

// GenYamlTreeCustom creates yaml structured ref files.
func GenYamlTreeCustom(cmd *cobra.Command, dir string, filePrepender, linkHandler func(string) string) error ***REMOVED***
	for _, c := range cmd.Commands() ***REMOVED***
		if !c.IsAvailableCommand() || c.IsAdditionalHelpTopicCommand() ***REMOVED***
			continue
		***REMOVED***
		if err := GenYamlTreeCustom(c, dir, filePrepender, linkHandler); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	basename := strings.Replace(cmd.CommandPath(), " ", "_", -1) + ".yaml"
	filename := filepath.Join(dir, basename)
	f, err := os.Create(filename)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer f.Close()

	if _, err := io.WriteString(f, filePrepender(filename)); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := GenYamlCustom(cmd, f, linkHandler); err != nil ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

// GenYaml creates yaml output.
func GenYaml(cmd *cobra.Command, w io.Writer) error ***REMOVED***
	return GenYamlCustom(cmd, w, func(s string) string ***REMOVED*** return s ***REMOVED***)
***REMOVED***

// GenYamlCustom creates custom yaml output.
func GenYamlCustom(cmd *cobra.Command, w io.Writer, linkHandler func(string) string) error ***REMOVED***
	cmd.InitDefaultHelpCmd()
	cmd.InitDefaultHelpFlag()

	yamlDoc := cmdDoc***REMOVED******REMOVED***
	yamlDoc.Name = cmd.CommandPath()

	yamlDoc.Synopsis = forceMultiLine(cmd.Short)
	yamlDoc.Description = forceMultiLine(cmd.Long)

	if len(cmd.Example) > 0 ***REMOVED***
		yamlDoc.Example = cmd.Example
	***REMOVED***

	flags := cmd.NonInheritedFlags()
	if flags.HasFlags() ***REMOVED***
		yamlDoc.Options = genFlagResult(flags)
	***REMOVED***
	flags = cmd.InheritedFlags()
	if flags.HasFlags() ***REMOVED***
		yamlDoc.InheritedOptions = genFlagResult(flags)
	***REMOVED***

	if hasSeeAlso(cmd) ***REMOVED***
		result := []string***REMOVED******REMOVED***
		if cmd.HasParent() ***REMOVED***
			parent := cmd.Parent()
			result = append(result, parent.CommandPath()+" - "+parent.Short)
		***REMOVED***
		children := cmd.Commands()
		sort.Sort(byName(children))
		for _, child := range children ***REMOVED***
			if !child.IsAvailableCommand() || child.IsAdditionalHelpTopicCommand() ***REMOVED***
				continue
			***REMOVED***
			result = append(result, child.Name()+" - "+child.Short)
		***REMOVED***
		yamlDoc.SeeAlso = result
	***REMOVED***

	final, err := yaml.Marshal(&yamlDoc)
	if err != nil ***REMOVED***
		fmt.Println(err)
		os.Exit(1)
	***REMOVED***

	if _, err := w.Write(final); err != nil ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

func genFlagResult(flags *pflag.FlagSet) []cmdOption ***REMOVED***
	var result []cmdOption

	flags.VisitAll(func(flag *pflag.Flag) ***REMOVED***
		// Todo, when we mark a shorthand is deprecated, but specify an empty message.
		// The flag.ShorthandDeprecated is empty as the shorthand is deprecated.
		// Using len(flag.ShorthandDeprecated) > 0 can't handle this, others are ok.
		if !(len(flag.ShorthandDeprecated) > 0) && len(flag.Shorthand) > 0 ***REMOVED***
			opt := cmdOption***REMOVED***
				flag.Name,
				flag.Shorthand,
				flag.DefValue,
				forceMultiLine(flag.Usage),
			***REMOVED***
			result = append(result, opt)
		***REMOVED*** else ***REMOVED***
			opt := cmdOption***REMOVED***
				Name:         flag.Name,
				DefaultValue: forceMultiLine(flag.DefValue),
				Usage:        forceMultiLine(flag.Usage),
			***REMOVED***
			result = append(result, opt)
		***REMOVED***
	***REMOVED***)

	return result
***REMOVED***
