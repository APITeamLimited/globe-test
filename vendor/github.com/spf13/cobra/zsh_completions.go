package cobra

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
)

// GenZshCompletionFile generates zsh completion file.
func (c *Command) GenZshCompletionFile(filename string) error ***REMOVED***
	outFile, err := os.Create(filename)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer outFile.Close()

	return c.GenZshCompletion(outFile)
***REMOVED***

// GenZshCompletion generates a zsh completion file and writes to the passed writer.
func (c *Command) GenZshCompletion(w io.Writer) error ***REMOVED***
	buf := new(bytes.Buffer)

	writeHeader(buf, c)
	maxDepth := maxDepth(c)
	writeLevelMapping(buf, maxDepth)
	writeLevelCases(buf, maxDepth, c)

	_, err := buf.WriteTo(w)
	return err
***REMOVED***

func writeHeader(w io.Writer, cmd *Command) ***REMOVED***
	fmt.Fprintf(w, "#compdef %s\n\n", cmd.Name())
***REMOVED***

func maxDepth(c *Command) int ***REMOVED***
	if len(c.Commands()) == 0 ***REMOVED***
		return 0
	***REMOVED***
	maxDepthSub := 0
	for _, s := range c.Commands() ***REMOVED***
		subDepth := maxDepth(s)
		if subDepth > maxDepthSub ***REMOVED***
			maxDepthSub = subDepth
		***REMOVED***
	***REMOVED***
	return 1 + maxDepthSub
***REMOVED***

func writeLevelMapping(w io.Writer, numLevels int) ***REMOVED***
	fmt.Fprintln(w, `_arguments \`)
	for i := 1; i <= numLevels; i++ ***REMOVED***
		fmt.Fprintf(w, `  '%d: :->level%d' \`, i, i)
		fmt.Fprintln(w)
	***REMOVED***
	fmt.Fprintf(w, `  '%d: :%s'`, numLevels+1, "_files")
	fmt.Fprintln(w)
***REMOVED***

func writeLevelCases(w io.Writer, maxDepth int, root *Command) ***REMOVED***
	fmt.Fprintln(w, "case $state in")
	defer fmt.Fprintln(w, "esac")

	for i := 1; i <= maxDepth; i++ ***REMOVED***
		fmt.Fprintf(w, "  level%d)\n", i)
		writeLevel(w, root, i)
		fmt.Fprintln(w, "  ;;")
	***REMOVED***
	fmt.Fprintln(w, "  *)")
	fmt.Fprintln(w, "    _arguments '*: :_files'")
	fmt.Fprintln(w, "  ;;")
***REMOVED***

func writeLevel(w io.Writer, root *Command, i int) ***REMOVED***
	fmt.Fprintf(w, "    case $words[%d] in\n", i)
	defer fmt.Fprintln(w, "    esac")

	commands := filterByLevel(root, i)
	byParent := groupByParent(commands)

	for p, c := range byParent ***REMOVED***
		names := names(c)
		fmt.Fprintf(w, "      %s)\n", p)
		fmt.Fprintf(w, "        _arguments '%d: :(%s)'\n", i, strings.Join(names, " "))
		fmt.Fprintln(w, "      ;;")
	***REMOVED***
	fmt.Fprintln(w, "      *)")
	fmt.Fprintln(w, "        _arguments '*: :_files'")
	fmt.Fprintln(w, "      ;;")

***REMOVED***

func filterByLevel(c *Command, l int) []*Command ***REMOVED***
	cs := make([]*Command, 0)
	if l == 0 ***REMOVED***
		cs = append(cs, c)
		return cs
	***REMOVED***
	for _, s := range c.Commands() ***REMOVED***
		cs = append(cs, filterByLevel(s, l-1)...)
	***REMOVED***
	return cs
***REMOVED***

func groupByParent(commands []*Command) map[string][]*Command ***REMOVED***
	m := make(map[string][]*Command)
	for _, c := range commands ***REMOVED***
		parent := c.Parent()
		if parent == nil ***REMOVED***
			continue
		***REMOVED***
		m[parent.Name()] = append(m[parent.Name()], c)
	***REMOVED***
	return m
***REMOVED***

func names(commands []*Command) []string ***REMOVED***
	ns := make([]string, len(commands))
	for i, c := range commands ***REMOVED***
		ns[i] = c.Name()
	***REMOVED***
	return ns
***REMOVED***
