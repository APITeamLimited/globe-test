package sourcemap

import (
	"errors"
	"io"
	"strings"

	"github.com/go-sourcemap/sourcemap/internal/base64vlq"
)

type fn func(m *mappings) (fn, error)

type mapping struct ***REMOVED***
	genLine      int32
	genColumn    int32
	sourcesInd   int32
	sourceLine   int32
	sourceColumn int32
	namesInd     int32
***REMOVED***

type mappings struct ***REMOVED***
	rd  *strings.Reader
	dec base64vlq.Decoder

	hasValue bool
	hasName  bool
	value    mapping

	values []mapping
***REMOVED***

func parseMappings(s string) ([]mapping, error) ***REMOVED***
	if s == "" ***REMOVED***
		return nil, errors.New("sourcemap: mappings are empty")
	***REMOVED***

	rd := strings.NewReader(s)
	m := &mappings***REMOVED***
		rd:  rd,
		dec: base64vlq.NewDecoder(rd),

		values: make([]mapping, 0, mappingsNumber(s)),
	***REMOVED***
	m.value.genLine = 1
	m.value.sourceLine = 1

	err := m.parse()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	values := m.values
	m.values = nil
	return values, nil
***REMOVED***

func mappingsNumber(s string) int ***REMOVED***
	return strings.Count(s, ",") + strings.Count(s, ";")
***REMOVED***

func (m *mappings) parse() error ***REMOVED***
	next := parseGenCol
	for ***REMOVED***
		c, err := m.rd.ReadByte()
		if err == io.EOF ***REMOVED***
			m.pushValue()
			return nil
		***REMOVED***
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		switch c ***REMOVED***
		case ',':
			m.pushValue()
			next = parseGenCol
		case ';':
			m.pushValue()

			m.value.genLine++
			m.value.genColumn = 0

			next = parseGenCol
		default:
			err := m.rd.UnreadByte()
			if err != nil ***REMOVED***
				return err
			***REMOVED***

			next, err = next(m)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			m.hasValue = true
		***REMOVED***
	***REMOVED***
***REMOVED***

func parseGenCol(m *mappings) (fn, error) ***REMOVED***
	n, err := m.dec.Decode()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	m.value.genColumn += n
	return parseSourcesInd, nil
***REMOVED***

func parseSourcesInd(m *mappings) (fn, error) ***REMOVED***
	n, err := m.dec.Decode()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	m.value.sourcesInd += n
	return parseSourceLine, nil
***REMOVED***

func parseSourceLine(m *mappings) (fn, error) ***REMOVED***
	n, err := m.dec.Decode()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	m.value.sourceLine += n
	return parseSourceCol, nil
***REMOVED***

func parseSourceCol(m *mappings) (fn, error) ***REMOVED***
	n, err := m.dec.Decode()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	m.value.sourceColumn += n
	return parseNamesInd, nil
***REMOVED***

func parseNamesInd(m *mappings) (fn, error) ***REMOVED***
	n, err := m.dec.Decode()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	m.hasName = true
	m.value.namesInd += n
	return parseGenCol, nil
***REMOVED***

func (m *mappings) pushValue() ***REMOVED***
	if !m.hasValue ***REMOVED***
		return
	***REMOVED***
	m.hasValue = false
	if m.hasName ***REMOVED***
		m.values = append(m.values, m.value)
		m.hasName = false
	***REMOVED*** else ***REMOVED***
		m.values = append(m.values, mapping***REMOVED***
			genLine:      m.value.genLine,
			genColumn:    m.value.genColumn,
			sourcesInd:   m.value.sourcesInd,
			sourceLine:   m.value.sourceLine,
			sourceColumn: m.value.sourceColumn,
			namesInd:     -1,
		***REMOVED***)
	***REMOVED***
***REMOVED***
