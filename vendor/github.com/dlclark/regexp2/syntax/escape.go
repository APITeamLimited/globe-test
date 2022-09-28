package syntax

import (
	"bytes"
	"strconv"
	"strings"
	"unicode"
)

func Escape(input string) string ***REMOVED***
	b := &bytes.Buffer***REMOVED******REMOVED***
	for _, r := range input ***REMOVED***
		escape(b, r, false)
	***REMOVED***
	return b.String()
***REMOVED***

const meta = `\.+*?()|[]***REMOVED******REMOVED***^$# `

func escape(b *bytes.Buffer, r rune, force bool) ***REMOVED***
	if unicode.IsPrint(r) ***REMOVED***
		if strings.IndexRune(meta, r) >= 0 || force ***REMOVED***
			b.WriteRune('\\')
		***REMOVED***
		b.WriteRune(r)
		return
	***REMOVED***

	switch r ***REMOVED***
	case '\a':
		b.WriteString(`\a`)
	case '\f':
		b.WriteString(`\f`)
	case '\n':
		b.WriteString(`\n`)
	case '\r':
		b.WriteString(`\r`)
	case '\t':
		b.WriteString(`\t`)
	case '\v':
		b.WriteString(`\v`)
	default:
		if r < 0x100 ***REMOVED***
			b.WriteString(`\x`)
			s := strconv.FormatInt(int64(r), 16)
			if len(s) == 1 ***REMOVED***
				b.WriteRune('0')
			***REMOVED***
			b.WriteString(s)
			break
		***REMOVED***
		b.WriteString(`\u`)
		b.WriteString(strconv.FormatInt(int64(r), 16))
	***REMOVED***
***REMOVED***

func Unescape(input string) (string, error) ***REMOVED***
	idx := strings.IndexRune(input, '\\')
	// no slashes means no unescape needed
	if idx == -1 ***REMOVED***
		return input, nil
	***REMOVED***

	buf := bytes.NewBufferString(input[:idx])
	// get the runes for the rest of the string -- we're going full parser scan on this

	p := parser***REMOVED******REMOVED***
	p.setPattern(input[idx+1:])
	for ***REMOVED***
		if p.rightMost() ***REMOVED***
			return "", p.getErr(ErrIllegalEndEscape)
		***REMOVED***
		r, err := p.scanCharEscape()
		if err != nil ***REMOVED***
			return "", err
		***REMOVED***
		buf.WriteRune(r)
		// are we done?
		if p.rightMost() ***REMOVED***
			return buf.String(), nil
		***REMOVED***

		r = p.moveRightGetChar()
		for r != '\\' ***REMOVED***
			buf.WriteRune(r)
			if p.rightMost() ***REMOVED***
				// we're done, no more slashes
				return buf.String(), nil
			***REMOVED***
			// keep scanning until we get another slash
			r = p.moveRightGetChar()
		***REMOVED***
	***REMOVED***
***REMOVED***
