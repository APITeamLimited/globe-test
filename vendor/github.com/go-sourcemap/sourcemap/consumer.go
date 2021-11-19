package sourcemap

import (
	"encoding/json"
	"fmt"
	"net/url"
	"path"
	"sort"
)

type sourceMap struct ***REMOVED***
	Version        int               `json:"version"`
	File           string            `json:"file"`
	SourceRoot     string            `json:"sourceRoot"`
	Sources        []string          `json:"sources"`
	SourcesContent []string          `json:"sourcesContent"`
	Names          []json.RawMessage `json:"names,string"`
	Mappings       string            `json:"mappings"`

	mappings []mapping
***REMOVED***

type v3 struct ***REMOVED***
	sourceMap
	Sections []section `json:"sections"`
***REMOVED***

func (m *sourceMap) parse(sourcemapURL string) error ***REMOVED***
	if err := checkVersion(m.Version); err != nil ***REMOVED***
		return err
	***REMOVED***

	var sourceRootURL *url.URL
	if m.SourceRoot != "" ***REMOVED***
		u, err := url.Parse(m.SourceRoot)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if u.IsAbs() ***REMOVED***
			sourceRootURL = u
		***REMOVED***
	***REMOVED*** else if sourcemapURL != "" ***REMOVED***
		u, err := url.Parse(sourcemapURL)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if u.IsAbs() ***REMOVED***
			u.Path = path.Dir(u.Path)
			sourceRootURL = u
		***REMOVED***
	***REMOVED***

	for i, src := range m.Sources ***REMOVED***
		m.Sources[i] = m.absSource(sourceRootURL, src)
	***REMOVED***

	mappings, err := parseMappings(m.Mappings)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	m.mappings = mappings
	// Free memory.
	m.Mappings = ""

	return nil
***REMOVED***

func (m *sourceMap) absSource(root *url.URL, source string) string ***REMOVED***
	if path.IsAbs(source) ***REMOVED***
		return source
	***REMOVED***

	if u, err := url.Parse(source); err == nil && u.IsAbs() ***REMOVED***
		return source
	***REMOVED***

	if root != nil ***REMOVED***
		u := *root
		u.Path = path.Join(u.Path, source)
		return u.String()
	***REMOVED***

	if m.SourceRoot != "" ***REMOVED***
		return path.Join(m.SourceRoot, source)
	***REMOVED***

	return source
***REMOVED***

func (m *sourceMap) name(idx int) string ***REMOVED***
	if idx >= len(m.Names) ***REMOVED***
		return ""
	***REMOVED***

	raw := m.Names[idx]
	if len(raw) == 0 ***REMOVED***
		return ""
	***REMOVED***

	if raw[0] == '"' && raw[len(raw)-1] == '"' ***REMOVED***
		var str string
		if err := json.Unmarshal(raw, &str); err == nil ***REMOVED***
			return str
		***REMOVED***
	***REMOVED***

	return string(raw)
***REMOVED***

type section struct ***REMOVED***
	Offset struct ***REMOVED***
		Line   int `json:"line"`
		Column int `json:"column"`
	***REMOVED*** `json:"offset"`
	Map *sourceMap `json:"map"`
***REMOVED***

type Consumer struct ***REMOVED***
	sourcemapURL string
	file         string
	sections     []section
***REMOVED***

func Parse(sourcemapURL string, b []byte) (*Consumer, error) ***REMOVED***
	v3 := new(v3)
	err := json.Unmarshal(b, v3)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if err := checkVersion(v3.Version); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if len(v3.Sections) == 0 ***REMOVED***
		v3.Sections = append(v3.Sections, section***REMOVED***
			Map: &v3.sourceMap,
		***REMOVED***)
	***REMOVED***

	for _, s := range v3.Sections ***REMOVED***
		err := s.Map.parse(sourcemapURL)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	reverse(v3.Sections)
	return &Consumer***REMOVED***
		sourcemapURL: sourcemapURL,
		file:         v3.File,
		sections:     v3.Sections,
	***REMOVED***, nil
***REMOVED***

func (c *Consumer) SourcemapURL() string ***REMOVED***
	return c.sourcemapURL
***REMOVED***

// File returns an optional name of the generated code
// that this source map is associated with.
func (c *Consumer) File() string ***REMOVED***
	return c.file
***REMOVED***

// Source returns the original source, name, line, and column information
// for the generated source's line and column positions.
func (c *Consumer) Source(
	genLine, genColumn int,
) (source, name string, line, column int, ok bool) ***REMOVED***
	for i := range c.sections ***REMOVED***
		s := &c.sections[i]
		if s.Offset.Line < genLine ||
			(s.Offset.Line+1 == genLine && s.Offset.Column <= genColumn) ***REMOVED***
			genLine -= s.Offset.Line
			genColumn -= s.Offset.Column
			return c.source(s.Map, genLine, genColumn)
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

func (c *Consumer) source(
	m *sourceMap, genLine, genColumn int,
) (source, name string, line, column int, ok bool) ***REMOVED***
	i := sort.Search(len(m.mappings), func(i int) bool ***REMOVED***
		m := &m.mappings[i]
		if int(m.genLine) == genLine ***REMOVED***
			return int(m.genColumn) >= genColumn
		***REMOVED***
		return int(m.genLine) >= genLine
	***REMOVED***)

	var match *mapping
	// Mapping not found
	if i == len(m.mappings) ***REMOVED***
		// lets see if the line is correct but the column is bigger
		match = &m.mappings[i-1]
		if int(match.genLine) != genLine ***REMOVED***
			return
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		match = &m.mappings[i]

		// Fuzzy match.
		if int(match.genLine) > genLine || int(match.genColumn) > genColumn ***REMOVED***
			if i == 0 ***REMOVED***
				return
			***REMOVED***
			match = &m.mappings[i-1]
		***REMOVED***
	***REMOVED***

	if match.sourcesInd >= 0 ***REMOVED***
		source = m.Sources[match.sourcesInd]
	***REMOVED***
	if match.namesInd >= 0 ***REMOVED***
		name = m.name(int(match.namesInd))
	***REMOVED***
	line = int(match.sourceLine)
	column = int(match.sourceColumn)
	ok = true
	return
***REMOVED***

// SourceContent returns the original source content for the source.
func (c *Consumer) SourceContent(source string) string ***REMOVED***
	for i := range c.sections ***REMOVED***
		s := &c.sections[i]
		for i, src := range s.Map.Sources ***REMOVED***
			if src == source ***REMOVED***
				if i < len(s.Map.SourcesContent) ***REMOVED***
					return s.Map.SourcesContent[i]
				***REMOVED***
				break
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return ""
***REMOVED***

func checkVersion(version int) error ***REMOVED***
	if version == 3 || version == 0 ***REMOVED***
		return nil
	***REMOVED***
	return fmt.Errorf(
		"sourcemap: got version=%d, but only 3rd version is supported",
		version,
	)
***REMOVED***

func reverse(ss []section) ***REMOVED***
	last := len(ss) - 1
	for i := 0; i < len(ss)/2; i++ ***REMOVED***
		ss[i], ss[last-i] = ss[last-i], ss[i]
	***REMOVED***
***REMOVED***
