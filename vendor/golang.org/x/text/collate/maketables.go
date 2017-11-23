// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

// Collation table generator.
// Data read from the web.

package main

import (
	"archive/zip"
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"

	"golang.org/x/text/collate"
	"golang.org/x/text/collate/build"
	"golang.org/x/text/internal/colltab"
	"golang.org/x/text/internal/gen"
	"golang.org/x/text/language"
	"golang.org/x/text/unicode/cldr"
)

var (
	test = flag.Bool("test", false,
		"test existing tables; can be used to compare web data with package data.")
	short = flag.Bool("short", false, `Use "short" alternatives, when available.`)
	draft = flag.Bool("draft", false, `Use draft versions, when available.`)
	tags  = flag.String("tags", "", "build tags to be included after +build directive")
	pkg   = flag.String("package", "collate",
		"the name of the package in which the generated file is to be included")

	tables = flagStringSetAllowAll("tables", "collate", "collate,chars",
		"comma-spearated list of tables to generate.")
	exclude = flagStringSet("exclude", "zh2", "",
		"comma-separated list of languages to exclude.")
	include = flagStringSet("include", "", "",
		"comma-separated list of languages to include. Include trumps exclude.")
	// TODO: Not included: unihan gb2312han zhuyin big5han (for size reasons)
	// TODO: Not included: traditional (buggy for Bengali)
	types = flagStringSetAllowAll("types", "standard,phonebook,phonetic,reformed,pinyin,stroke", "",
		"comma-separated list of types that should be included.")
)

// stringSet implements an ordered set based on a list.  It implements flag.Value
// to allow a set to be specified as a comma-separated list.
type stringSet struct ***REMOVED***
	s        []string
	allowed  *stringSet
	dirty    bool // needs compaction if true
	all      bool
	allowAll bool
***REMOVED***

func flagStringSet(name, def, allowed, usage string) *stringSet ***REMOVED***
	ss := &stringSet***REMOVED******REMOVED***
	if allowed != "" ***REMOVED***
		usage += fmt.Sprintf(" (allowed values: any of %s)", allowed)
		ss.allowed = &stringSet***REMOVED******REMOVED***
		failOnError(ss.allowed.Set(allowed))
	***REMOVED***
	ss.Set(def)
	flag.Var(ss, name, usage)
	return ss
***REMOVED***

func flagStringSetAllowAll(name, def, allowed, usage string) *stringSet ***REMOVED***
	ss := &stringSet***REMOVED***allowAll: true***REMOVED***
	if allowed == "" ***REMOVED***
		flag.Var(ss, name, usage+fmt.Sprintf(` Use "all" to select all.`))
	***REMOVED*** else ***REMOVED***
		ss.allowed = &stringSet***REMOVED******REMOVED***
		failOnError(ss.allowed.Set(allowed))
		flag.Var(ss, name, usage+fmt.Sprintf(` (allowed values: "all" or any of %s)`, allowed))
	***REMOVED***
	ss.Set(def)
	return ss
***REMOVED***

func (ss stringSet) Len() int ***REMOVED***
	return len(ss.s)
***REMOVED***

func (ss stringSet) String() string ***REMOVED***
	return strings.Join(ss.s, ",")
***REMOVED***

func (ss *stringSet) Set(s string) error ***REMOVED***
	if ss.allowAll && s == "all" ***REMOVED***
		ss.s = nil
		ss.all = true
		return nil
	***REMOVED***
	ss.s = ss.s[:0]
	for _, s := range strings.Split(s, ",") ***REMOVED***
		if s := strings.TrimSpace(s); s != "" ***REMOVED***
			if ss.allowed != nil && !ss.allowed.contains(s) ***REMOVED***
				return fmt.Errorf("unsupported value %q; must be one of %s", s, ss.allowed)
			***REMOVED***
			ss.add(s)
		***REMOVED***
	***REMOVED***
	ss.compact()
	return nil
***REMOVED***

func (ss *stringSet) add(s string) ***REMOVED***
	ss.s = append(ss.s, s)
	ss.dirty = true
***REMOVED***

func (ss *stringSet) values() []string ***REMOVED***
	ss.compact()
	return ss.s
***REMOVED***

func (ss *stringSet) contains(s string) bool ***REMOVED***
	if ss.all ***REMOVED***
		return true
	***REMOVED***
	for _, v := range ss.s ***REMOVED***
		if v == s ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

func (ss *stringSet) compact() ***REMOVED***
	if !ss.dirty ***REMOVED***
		return
	***REMOVED***
	a := ss.s
	sort.Strings(a)
	k := 0
	for i := 1; i < len(a); i++ ***REMOVED***
		if a[k] != a[i] ***REMOVED***
			a[k+1] = a[i]
			k++
		***REMOVED***
	***REMOVED***
	ss.s = a[:k+1]
	ss.dirty = false
***REMOVED***

func skipLang(l string) bool ***REMOVED***
	if include.Len() > 0 ***REMOVED***
		return !include.contains(l)
	***REMOVED***
	return exclude.contains(l)
***REMOVED***

// altInclude returns a list of alternatives (for the LDML alt attribute)
// in order of preference.  An empty string in this list indicates the
// default entry.
func altInclude() []string ***REMOVED***
	l := []string***REMOVED******REMOVED***
	if *short ***REMOVED***
		l = append(l, "short")
	***REMOVED***
	l = append(l, "")
	// TODO: handle draft using cldr.SetDraftLevel
	if *draft ***REMOVED***
		l = append(l, "proposed")
	***REMOVED***
	return l
***REMOVED***

func failOnError(e error) ***REMOVED***
	if e != nil ***REMOVED***
		log.Panic(e)
	***REMOVED***
***REMOVED***

func openArchive() *zip.Reader ***REMOVED***
	f := gen.OpenCLDRCoreZip()
	buffer, err := ioutil.ReadAll(f)
	f.Close()
	failOnError(err)
	archive, err := zip.NewReader(bytes.NewReader(buffer), int64(len(buffer)))
	failOnError(err)
	return archive
***REMOVED***

// parseUCA parses a Default Unicode Collation Element Table of the format
// specified in http://www.unicode.org/reports/tr10/#File_Format.
// It returns the variable top.
func parseUCA(builder *build.Builder) ***REMOVED***
	var r io.ReadCloser
	var err error
	for _, f := range openArchive().File ***REMOVED***
		if strings.HasSuffix(f.Name, "allkeys_CLDR.txt") ***REMOVED***
			r, err = f.Open()
		***REMOVED***
	***REMOVED***
	if r == nil ***REMOVED***
		log.Fatal("File allkeys_CLDR.txt not found in archive.")
	***REMOVED***
	failOnError(err)
	defer r.Close()
	scanner := bufio.NewScanner(r)
	colelem := regexp.MustCompile(`\[([.*])([0-9A-F.]+)\]`)
	for i := 1; scanner.Scan(); i++ ***REMOVED***
		line := scanner.Text()
		if len(line) == 0 || line[0] == '#' ***REMOVED***
			continue
		***REMOVED***
		if line[0] == '@' ***REMOVED***
			// parse properties
			switch ***REMOVED***
			case strings.HasPrefix(line[1:], "version "):
				a := strings.Split(line[1:], " ")
				if a[1] != gen.UnicodeVersion() ***REMOVED***
					log.Fatalf("incompatible version %s; want %s", a[1], gen.UnicodeVersion())
				***REMOVED***
			case strings.HasPrefix(line[1:], "backwards "):
				log.Fatalf("%d: unsupported option backwards", i)
			default:
				log.Printf("%d: unknown option %s", i, line[1:])
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			// parse entries
			part := strings.Split(line, " ; ")
			if len(part) != 2 ***REMOVED***
				log.Fatalf("%d: production rule without ';': %v", i, line)
			***REMOVED***
			lhs := []rune***REMOVED******REMOVED***
			for _, v := range strings.Split(part[0], " ") ***REMOVED***
				if v == "" ***REMOVED***
					continue
				***REMOVED***
				lhs = append(lhs, rune(convHex(i, v)))
			***REMOVED***
			var n int
			var vars []int
			rhs := [][]int***REMOVED******REMOVED***
			for i, m := range colelem.FindAllStringSubmatch(part[1], -1) ***REMOVED***
				n += len(m[0])
				elem := []int***REMOVED******REMOVED***
				for _, h := range strings.Split(m[2], ".") ***REMOVED***
					elem = append(elem, convHex(i, h))
				***REMOVED***
				if m[1] == "*" ***REMOVED***
					vars = append(vars, i)
				***REMOVED***
				rhs = append(rhs, elem)
			***REMOVED***
			if len(part[1]) < n+3 || part[1][n+1] != '#' ***REMOVED***
				log.Fatalf("%d: expected comment; found %s", i, part[1][n:])
			***REMOVED***
			if *test ***REMOVED***
				testInput.add(string(lhs))
			***REMOVED***
			failOnError(builder.Add(lhs, rhs, vars))
		***REMOVED***
	***REMOVED***
	if scanner.Err() != nil ***REMOVED***
		log.Fatal(scanner.Err())
	***REMOVED***
***REMOVED***

func convHex(line int, s string) int ***REMOVED***
	r, e := strconv.ParseInt(s, 16, 32)
	if e != nil ***REMOVED***
		log.Fatalf("%d: %v", line, e)
	***REMOVED***
	return int(r)
***REMOVED***

var testInput = stringSet***REMOVED******REMOVED***

var charRe = regexp.MustCompile(`&#x([0-9A-F]*);`)
var tagRe = regexp.MustCompile(`<([a-z_]*)  */>`)

var mainLocales = []string***REMOVED******REMOVED***

// charsets holds a list of exemplar characters per category.
type charSets map[string][]string

func (p charSets) fprint(w io.Writer) ***REMOVED***
	fmt.Fprintln(w, "[exN]string***REMOVED***")
	for i, k := range []string***REMOVED***"", "contractions", "punctuation", "auxiliary", "currencySymbol", "index"***REMOVED*** ***REMOVED***
		if set := p[k]; len(set) != 0 ***REMOVED***
			fmt.Fprintf(w, "\t\t%d: %q,\n", i, strings.Join(set, " "))
		***REMOVED***
	***REMOVED***
	fmt.Fprintln(w, "\t***REMOVED***,")
***REMOVED***

var localeChars = make(map[string]charSets)

const exemplarHeader = `
type exemplarType int
const (
	exCharacters exemplarType = iota
	exContractions
	exPunctuation
	exAuxiliary
	exCurrency
	exIndex
	exN
)
`

func printExemplarCharacters(w io.Writer) ***REMOVED***
	fmt.Fprintln(w, exemplarHeader)
	fmt.Fprintln(w, "var exemplarCharacters = map[string][exN]string***REMOVED***")
	for _, loc := range mainLocales ***REMOVED***
		fmt.Fprintf(w, "\t%q: ", loc)
		localeChars[loc].fprint(w)
	***REMOVED***
	fmt.Fprintln(w, "***REMOVED***")
***REMOVED***

func decodeCLDR(d *cldr.Decoder) *cldr.CLDR ***REMOVED***
	r := gen.OpenCLDRCoreZip()
	data, err := d.DecodeZip(r)
	failOnError(err)
	return data
***REMOVED***

// parseMain parses XML files in the main directory of the CLDR core.zip file.
func parseMain() ***REMOVED***
	d := &cldr.Decoder***REMOVED******REMOVED***
	d.SetDirFilter("main")
	d.SetSectionFilter("characters")
	data := decodeCLDR(d)
	for _, loc := range data.Locales() ***REMOVED***
		x := data.RawLDML(loc)
		if skipLang(x.Identity.Language.Type) ***REMOVED***
			continue
		***REMOVED***
		if x.Characters != nil ***REMOVED***
			x, _ = data.LDML(loc)
			loc = language.Make(loc).String()
			for _, ec := range x.Characters.ExemplarCharacters ***REMOVED***
				if ec.Draft != "" ***REMOVED***
					continue
				***REMOVED***
				if _, ok := localeChars[loc]; !ok ***REMOVED***
					mainLocales = append(mainLocales, loc)
					localeChars[loc] = make(charSets)
				***REMOVED***
				localeChars[loc][ec.Type] = parseCharacters(ec.Data())
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func parseCharacters(chars string) []string ***REMOVED***
	parseSingle := func(s string) (r rune, tail string, escaped bool) ***REMOVED***
		if s[0] == '\\' ***REMOVED***
			return rune(s[1]), s[2:], true
		***REMOVED***
		r, sz := utf8.DecodeRuneInString(s)
		return r, s[sz:], false
	***REMOVED***
	chars = strings.TrimSpace(chars)
	if n := len(chars) - 1; chars[n] == ']' && chars[0] == '[' ***REMOVED***
		chars = chars[1:n]
	***REMOVED***
	list := []string***REMOVED******REMOVED***
	var r, last, end rune
	for len(chars) > 0 ***REMOVED***
		if chars[0] == '***REMOVED***' ***REMOVED*** // character sequence
			buf := []rune***REMOVED******REMOVED***
			for chars = chars[1:]; len(chars) > 0; ***REMOVED***
				r, chars, _ = parseSingle(chars)
				if r == '***REMOVED***' ***REMOVED***
					break
				***REMOVED***
				if r == ' ' ***REMOVED***
					log.Fatalf("space not supported in sequence %q", chars)
				***REMOVED***
				buf = append(buf, r)
			***REMOVED***
			list = append(list, string(buf))
			last = 0
		***REMOVED*** else ***REMOVED*** // single character
			escaped := false
			r, chars, escaped = parseSingle(chars)
			if r != ' ' ***REMOVED***
				if r == '-' && !escaped ***REMOVED***
					if last == 0 ***REMOVED***
						log.Fatal("'-' should be preceded by a character")
					***REMOVED***
					end, chars, _ = parseSingle(chars)
					for ; last <= end; last++ ***REMOVED***
						list = append(list, string(last))
					***REMOVED***
					last = 0
				***REMOVED*** else ***REMOVED***
					list = append(list, string(r))
					last = r
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return list
***REMOVED***

var fileRe = regexp.MustCompile(`.*/collation/(.*)\.xml`)

// typeMap translates legacy type keys to their BCP47 equivalent.
var typeMap = map[string]string***REMOVED***
	"phonebook":   "phonebk",
	"traditional": "trad",
***REMOVED***

// parseCollation parses XML files in the collation directory of the CLDR core.zip file.
func parseCollation(b *build.Builder) ***REMOVED***
	d := &cldr.Decoder***REMOVED******REMOVED***
	d.SetDirFilter("collation")
	data := decodeCLDR(d)
	for _, loc := range data.Locales() ***REMOVED***
		x, err := data.LDML(loc)
		failOnError(err)
		if skipLang(x.Identity.Language.Type) ***REMOVED***
			continue
		***REMOVED***
		cs := x.Collations.Collation
		sl := cldr.MakeSlice(&cs)
		if len(types.s) == 0 ***REMOVED***
			sl.SelectAnyOf("type", x.Collations.Default())
		***REMOVED*** else if !types.all ***REMOVED***
			sl.SelectAnyOf("type", types.s...)
		***REMOVED***
		sl.SelectOnePerGroup("alt", altInclude())

		for _, c := range cs ***REMOVED***
			id, err := language.Parse(loc)
			if err != nil ***REMOVED***
				fmt.Fprintf(os.Stderr, "invalid locale: %q", err)
				continue
			***REMOVED***
			// Support both old- and new-style defaults.
			d := c.Type
			if x.Collations.DefaultCollation == nil ***REMOVED***
				d = x.Collations.Default()
			***REMOVED*** else ***REMOVED***
				d = x.Collations.DefaultCollation.Data()
			***REMOVED***
			// We assume tables are being built either for search or collation,
			// but not both. For search the default is always "search".
			if d != c.Type && c.Type != "search" ***REMOVED***
				typ := c.Type
				if len(c.Type) > 8 ***REMOVED***
					typ = typeMap[c.Type]
				***REMOVED***
				id, err = id.SetTypeForKey("co", typ)
				failOnError(err)
			***REMOVED***
			t := b.Tailoring(id)
			c.Process(processor***REMOVED***t***REMOVED***)
		***REMOVED***
	***REMOVED***
***REMOVED***

type processor struct ***REMOVED***
	t *build.Tailoring
***REMOVED***

func (p processor) Reset(anchor string, before int) (err error) ***REMOVED***
	if before != 0 ***REMOVED***
		err = p.t.SetAnchorBefore(anchor)
	***REMOVED*** else ***REMOVED***
		err = p.t.SetAnchor(anchor)
	***REMOVED***
	failOnError(err)
	return nil
***REMOVED***

func (p processor) Insert(level int, str, context, extend string) error ***REMOVED***
	str = context + str
	if *test ***REMOVED***
		testInput.add(str)
	***REMOVED***
	// TODO: mimic bug in old maketables: remove.
	err := p.t.Insert(colltab.Level(level-1), str, context+extend)
	failOnError(err)
	return nil
***REMOVED***

func (p processor) Index(id string) ***REMOVED***
***REMOVED***

func testCollator(c *collate.Collator) ***REMOVED***
	c0 := collate.New(language.Und)

	// iterator over all characters for all locales and check
	// whether Key is equal.
	buf := collate.Buffer***REMOVED******REMOVED***

	// Add all common and not too uncommon runes to the test set.
	for i := rune(0); i < 0x30000; i++ ***REMOVED***
		testInput.add(string(i))
	***REMOVED***
	for i := rune(0xE0000); i < 0xF0000; i++ ***REMOVED***
		testInput.add(string(i))
	***REMOVED***
	for _, str := range testInput.values() ***REMOVED***
		k0 := c0.KeyFromString(&buf, str)
		k := c.KeyFromString(&buf, str)
		if !bytes.Equal(k0, k) ***REMOVED***
			failOnError(fmt.Errorf("test:%U: keys differ (%x vs %x)", []rune(str), k0, k))
		***REMOVED***
		buf.Reset()
	***REMOVED***
	fmt.Println("PASS")
***REMOVED***

func main() ***REMOVED***
	gen.Init()
	b := build.NewBuilder()
	parseUCA(b)
	if tables.contains("chars") ***REMOVED***
		parseMain()
	***REMOVED***
	parseCollation(b)

	c, err := b.Build()
	failOnError(err)

	if *test ***REMOVED***
		testCollator(collate.NewFromTable(c))
	***REMOVED*** else ***REMOVED***
		w := &bytes.Buffer***REMOVED******REMOVED***

		gen.WriteUnicodeVersion(w)
		gen.WriteCLDRVersion(w)

		if tables.contains("collate") ***REMOVED***
			_, err = b.Print(w)
			failOnError(err)
		***REMOVED***
		if tables.contains("chars") ***REMOVED***
			printExemplarCharacters(w)
		***REMOVED***
		gen.WriteGoFile("tables.go", *pkg, w.Bytes())
	***REMOVED***
***REMOVED***
