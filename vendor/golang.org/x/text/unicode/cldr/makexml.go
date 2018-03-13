// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

// This tool generates types for the various XML formats of CLDR.
package main

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"

	"golang.org/x/text/internal/gen"
)

var outputFile = flag.String("output", "xml.go", "output file name")

func main() ***REMOVED***
	flag.Parse()

	r := gen.OpenCLDRCoreZip()
	buffer, err := ioutil.ReadAll(r)
	if err != nil ***REMOVED***
		log.Fatal("Could not read zip file")
	***REMOVED***
	r.Close()
	z, err := zip.NewReader(bytes.NewReader(buffer), int64(len(buffer)))
	if err != nil ***REMOVED***
		log.Fatalf("Could not read zip archive: %v", err)
	***REMOVED***

	var buf bytes.Buffer

	version := gen.CLDRVersion()

	for _, dtd := range files ***REMOVED***
		for _, f := range z.File ***REMOVED***
			if strings.HasSuffix(f.Name, dtd.file+".dtd") ***REMOVED***
				r, err := f.Open()
				failOnError(err)

				b := makeBuilder(&buf, dtd)
				b.parseDTD(r)
				b.resolve(b.index[dtd.top[0]])
				b.write()
				if b.version != "" && version != b.version ***REMOVED***
					println(f.Name)
					log.Fatalf("main: inconsistent versions: found %s; want %s", b.version, version)
				***REMOVED***
				break
			***REMOVED***
		***REMOVED***
	***REMOVED***
	fmt.Fprintln(&buf, "// Version is the version of CLDR from which the XML definitions are generated.")
	fmt.Fprintf(&buf, "const Version = %q\n", version)

	gen.WriteGoFile(*outputFile, "cldr", buf.Bytes())
***REMOVED***

func failOnError(err error) ***REMOVED***
	if err != nil ***REMOVED***
		log.New(os.Stderr, "", log.Lshortfile).Output(2, err.Error())
		os.Exit(1)
	***REMOVED***
***REMOVED***

// configuration data per DTD type
type dtd struct ***REMOVED***
	file string   // base file name
	root string   // Go name of the root XML element
	top  []string // create a different type for this section

	skipElem    []string // hard-coded or deprecated elements
	skipAttr    []string // attributes to exclude
	predefined  []string // hard-coded elements exist of the form <name>Elem
	forceRepeat []string // elements to make slices despite DTD
***REMOVED***

var files = []dtd***REMOVED***
	***REMOVED***
		file: "ldmlBCP47",
		root: "LDMLBCP47",
		top:  []string***REMOVED***"ldmlBCP47"***REMOVED***,
		skipElem: []string***REMOVED***
			"cldrVersion", // deprecated, not used
		***REMOVED***,
	***REMOVED***,
	***REMOVED***
		file: "ldmlSupplemental",
		root: "SupplementalData",
		top:  []string***REMOVED***"supplementalData"***REMOVED***,
		skipElem: []string***REMOVED***
			"cldrVersion", // deprecated, not used
		***REMOVED***,
		forceRepeat: []string***REMOVED***
			"plurals", // data defined in plurals.xml and ordinals.xml
		***REMOVED***,
	***REMOVED***,
	***REMOVED***
		file: "ldml",
		root: "LDML",
		top: []string***REMOVED***
			"ldml", "collation", "calendar", "timeZoneNames", "localeDisplayNames", "numbers",
		***REMOVED***,
		skipElem: []string***REMOVED***
			"cp",       // not used anywhere
			"special",  // not used anywhere
			"fallback", // deprecated, not used
			"alias",    // in Common
			"default",  // in Common
		***REMOVED***,
		skipAttr: []string***REMOVED***
			"hiraganaQuarternary", // typo in DTD, correct version included as well
		***REMOVED***,
		predefined: []string***REMOVED***"rules"***REMOVED***,
	***REMOVED***,
***REMOVED***

var comments = map[string]string***REMOVED***
	"ldmlBCP47": `
// LDMLBCP47 holds information on allowable values for various variables in LDML.
`,
	"supplementalData": `
// SupplementalData holds information relevant for internationalization
// and proper use of CLDR, but that is not contained in the locale hierarchy.
`,
	"ldml": `
// LDML is the top-level type for locale-specific data.
`,
	"collation": `
// Collation contains rules that specify a certain sort-order,
// as a tailoring of the root order. 
// The parsed rules are obtained by passing a RuleProcessor to Collation's
// Process method.
`,
	"calendar": `
// Calendar specifies the fields used for formatting and parsing dates and times.
// The month and quarter names are identified numerically, starting at 1.
// The day (of the week) names are identified with short strings, since there is
// no universally-accepted numeric designation.
`,
	"dates": `
// Dates contains information regarding the format and parsing of dates and times.
`,
	"localeDisplayNames": `
// LocaleDisplayNames specifies localized display names for for scripts, languages,
// countries, currencies, and variants.
`,
	"numbers": `
// Numbers supplies information for formatting and parsing numbers and currencies.
`,
***REMOVED***

type element struct ***REMOVED***
	name      string // XML element name
	category  string // elements contained by this element
	signature string // category + attrKey*

	attr []*attribute // attributes supported by this element.
	sub  []struct ***REMOVED***   // parsed and evaluated sub elements of this element.
		e      *element
		repeat bool // true if the element needs to be a slice
	***REMOVED***

	resolved bool // prevent multiple resolutions of this element.
***REMOVED***

type attribute struct ***REMOVED***
	name string
	key  string
	list []string

	tag string // Go tag
***REMOVED***

var (
	reHead  = regexp.MustCompile(` *(\w+) +([\w\-]+)`)
	reAttr  = regexp.MustCompile(` *(\w+) *(?:(\w+)|\(([\w\- \|]+)\)) *(?:#([A-Z]*) *(?:\"([\.\d+])\")?)? *("[\w\-:]*")?`)
	reElem  = regexp.MustCompile(`^ *(EMPTY|ANY|\(.*\)[\*\+\?]?) *$`)
	reToken = regexp.MustCompile(`\w\-`)
)

// builder is used to read in the DTD files from CLDR and generate Go code
// to be used with the encoding/xml package.
type builder struct ***REMOVED***
	w       io.Writer
	index   map[string]*element
	elem    []*element
	info    dtd
	version string
***REMOVED***

func makeBuilder(w io.Writer, d dtd) builder ***REMOVED***
	return builder***REMOVED***
		w:     w,
		index: make(map[string]*element),
		elem:  []*element***REMOVED******REMOVED***,
		info:  d,
	***REMOVED***
***REMOVED***

// parseDTD parses a DTD file.
func (b *builder) parseDTD(r io.Reader) ***REMOVED***
	for d := xml.NewDecoder(r); ; ***REMOVED***
		t, err := d.Token()
		if t == nil ***REMOVED***
			break
		***REMOVED***
		failOnError(err)
		dir, ok := t.(xml.Directive)
		if !ok ***REMOVED***
			continue
		***REMOVED***
		m := reHead.FindSubmatch(dir)
		dir = dir[len(m[0]):]
		ename := string(m[2])
		el, elementFound := b.index[ename]
		switch string(m[1]) ***REMOVED***
		case "ELEMENT":
			if elementFound ***REMOVED***
				log.Fatal("parseDTD: duplicate entry for element %q", ename)
			***REMOVED***
			m := reElem.FindSubmatch(dir)
			if m == nil ***REMOVED***
				log.Fatalf("parseDTD: invalid element %q", string(dir))
			***REMOVED***
			if len(m[0]) != len(dir) ***REMOVED***
				log.Fatal("parseDTD: invalid element %q", string(dir), len(dir), len(m[0]), string(m[0]))
			***REMOVED***
			s := string(m[1])
			el = &element***REMOVED***
				name:     ename,
				category: s,
			***REMOVED***
			b.index[ename] = el
		case "ATTLIST":
			if !elementFound ***REMOVED***
				log.Fatalf("parseDTD: unknown element %q", ename)
			***REMOVED***
			s := string(dir)
			m := reAttr.FindStringSubmatch(s)
			if m == nil ***REMOVED***
				log.Fatal(fmt.Errorf("parseDTD: invalid attribute %q", string(dir)))
			***REMOVED***
			if m[4] == "FIXED" ***REMOVED***
				b.version = m[5]
			***REMOVED*** else ***REMOVED***
				switch m[1] ***REMOVED***
				case "draft", "references", "alt", "validSubLocales", "standard" /* in Common */ :
				case "type", "choice":
				default:
					el.attr = append(el.attr, &attribute***REMOVED***
						name: m[1],
						key:  s,
						list: reToken.FindAllString(m[3], -1),
					***REMOVED***)
					el.signature = fmt.Sprintf("%s=%s+%s", el.signature, m[1], m[2])
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

var reCat = regexp.MustCompile(`[ ,\|]*(?:(\(|\)|\#?[\w_-]+)([\*\+\?]?))?`)

// resolve takes a parsed element and converts it into structured data
// that can be used to generate the XML code.
func (b *builder) resolve(e *element) ***REMOVED***
	if e.resolved ***REMOVED***
		return
	***REMOVED***
	b.elem = append(b.elem, e)
	e.resolved = true
	s := e.category
	found := make(map[string]bool)
	sequenceStart := []int***REMOVED******REMOVED***
	for len(s) > 0 ***REMOVED***
		m := reCat.FindStringSubmatch(s)
		if m == nil ***REMOVED***
			log.Fatalf("%s: invalid category string %q", e.name, s)
		***REMOVED***
		repeat := m[2] == "*" || m[2] == "+" || in(b.info.forceRepeat, m[1])
		switch m[1] ***REMOVED***
		case "":
		case "(":
			sequenceStart = append(sequenceStart, len(e.sub))
		case ")":
			if len(sequenceStart) == 0 ***REMOVED***
				log.Fatalf("%s: unmatched closing parenthesis", e.name)
			***REMOVED***
			for i := sequenceStart[len(sequenceStart)-1]; i < len(e.sub); i++ ***REMOVED***
				e.sub[i].repeat = e.sub[i].repeat || repeat
			***REMOVED***
			sequenceStart = sequenceStart[:len(sequenceStart)-1]
		default:
			if in(b.info.skipElem, m[1]) ***REMOVED***
			***REMOVED*** else if sub, ok := b.index[m[1]]; ok ***REMOVED***
				if !found[sub.name] ***REMOVED***
					e.sub = append(e.sub, struct ***REMOVED***
						e      *element
						repeat bool
					***REMOVED******REMOVED***sub, repeat***REMOVED***)
					found[sub.name] = true
					b.resolve(sub)
				***REMOVED***
			***REMOVED*** else if m[1] == "#PCDATA" || m[1] == "ANY" ***REMOVED***
			***REMOVED*** else if m[1] != "EMPTY" ***REMOVED***
				log.Fatalf("resolve:%s: element %q not found", e.name, m[1])
			***REMOVED***
		***REMOVED***
		s = s[len(m[0]):]
	***REMOVED***
***REMOVED***

// return true if s is contained in set.
func in(set []string, s string) bool ***REMOVED***
	for _, v := range set ***REMOVED***
		if v == s ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

var repl = strings.NewReplacer("-", " ", "_", " ")

// title puts the first character or each character following '_' in title case and
// removes all occurrences of '_'.
func title(s string) string ***REMOVED***
	return strings.Replace(strings.Title(repl.Replace(s)), " ", "", -1)
***REMOVED***

// writeElem generates Go code for a single element, recursively.
func (b *builder) writeElem(tab int, e *element) ***REMOVED***
	p := func(f string, x ...interface***REMOVED******REMOVED***) ***REMOVED***
		f = strings.Replace(f, "\n", "\n"+strings.Repeat("\t", tab), -1)
		fmt.Fprintf(b.w, f, x...)
	***REMOVED***
	if len(e.sub) == 0 && len(e.attr) == 0 ***REMOVED***
		p("Common")
		return
	***REMOVED***
	p("struct ***REMOVED***")
	tab++
	p("\nCommon")
	for _, attr := range e.attr ***REMOVED***
		if !in(b.info.skipAttr, attr.name) ***REMOVED***
			p("\n%s string `xml:\"%s,attr\"`", title(attr.name), attr.name)
		***REMOVED***
	***REMOVED***
	for _, sub := range e.sub ***REMOVED***
		if in(b.info.predefined, sub.e.name) ***REMOVED***
			p("\n%sElem", sub.e.name)
			continue
		***REMOVED***
		if in(b.info.skipElem, sub.e.name) ***REMOVED***
			continue
		***REMOVED***
		p("\n%s ", title(sub.e.name))
		if sub.repeat ***REMOVED***
			p("[]")
		***REMOVED***
		p("*")
		if in(b.info.top, sub.e.name) ***REMOVED***
			p(title(sub.e.name))
		***REMOVED*** else ***REMOVED***
			b.writeElem(tab, sub.e)
		***REMOVED***
		p(" `xml:\"%s\"`", sub.e.name)
	***REMOVED***
	tab--
	p("\n***REMOVED***")
***REMOVED***

// write generates the Go XML code.
func (b *builder) write() ***REMOVED***
	for i, name := range b.info.top ***REMOVED***
		e := b.index[name]
		if e != nil ***REMOVED***
			fmt.Fprintf(b.w, comments[name])
			name := title(e.name)
			if i == 0 ***REMOVED***
				name = b.info.root
			***REMOVED***
			fmt.Fprintf(b.w, "type %s ", name)
			b.writeElem(0, e)
			fmt.Fprint(b.w, "\n")
		***REMOVED***
	***REMOVED***
***REMOVED***
