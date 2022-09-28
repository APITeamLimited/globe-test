package goja

// This is a slightly modified version of the standard Go parser to make it more compatible with ECMAScript 5.1
// Changes:
// - 6-digit extended years are supported in place of long year (2006) in the form of +123456
// - Timezone formats tolerate colons, e.g. -0700 will parse -07:00
// - Short week day will also parse long week day
// - Short month ("Jan") will also parse long month ("January")
// - Long day ("02") will also parse short day ("2").
// - Timezone in brackets, "(MST)", will match any string in brackets (e.g. "(GMT Standard Time)")
// - If offset is not set and timezone name is unknown, an error is returned
// - If offset and timezone name are both set the offset takes precedence and the resulting Location will be FixedZone("", offset)

// Original copyright message:

// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"errors"
	"time"
)

const (
	_                        = iota
	stdLongMonth             = iota + stdNeedDate  // "January"
	stdMonth                                       // "Jan"
	stdNumMonth                                    // "1"
	stdZeroMonth                                   // "01"
	stdLongWeekDay                                 // "Monday"
	stdWeekDay                                     // "Mon"
	stdDay                                         // "2"
	stdUnderDay                                    // "_2"
	stdZeroDay                                     // "02"
	stdHour                  = iota + stdNeedClock // "15"
	stdHour12                                      // "3"
	stdZeroHour12                                  // "03"
	stdMinute                                      // "4"
	stdZeroMinute                                  // "04"
	stdSecond                                      // "5"
	stdZeroSecond                                  // "05"
	stdLongYear              = iota + stdNeedDate  // "2006"
	stdYear                                        // "06"
	stdPM                    = iota + stdNeedClock // "PM"
	stdpm                                          // "pm"
	stdTZ                    = iota                // "MST"
	stdBracketTZ                                   // "(MST)"
	stdISO8601TZ                                   // "Z0700"  // prints Z for UTC
	stdISO8601SecondsTZ                            // "Z070000"
	stdISO8601ShortTZ                              // "Z07"
	stdISO8601ColonTZ                              // "Z07:00" // prints Z for UTC
	stdISO8601ColonSecondsTZ                       // "Z07:00:00"
	stdNumTZ                                       // "-0700"  // always numeric
	stdNumSecondsTz                                // "-070000"
	stdNumShortTZ                                  // "-07"    // always numeric
	stdNumColonTZ                                  // "-07:00" // always numeric
	stdNumColonSecondsTZ                           // "-07:00:00"
	stdFracSecond0                                 // ".0", ".00", ... , trailing zeros included
	stdFracSecond9                                 // ".9", ".99", ..., trailing zeros omitted

	stdNeedDate  = 1 << 8             // need month, day, year
	stdNeedClock = 2 << 8             // need hour, minute, second
	stdArgShift  = 16                 // extra argument in high bits, above low stdArgShift
	stdMask      = 1<<stdArgShift - 1 // mask out argument
)

var errBad = errors.New("bad value for field") // placeholder not passed to user

func parseDate(layout, value string, defaultLocation *time.Location) (time.Time, error) ***REMOVED***
	alayout, avalue := layout, value
	rangeErrString := "" // set if a value is out of range
	amSet := false       // do we need to subtract 12 from the hour for midnight?
	pmSet := false       // do we need to add 12 to the hour?

	// Time being constructed.
	var (
		year       int
		month      int = 1 // January
		day        int = 1
		hour       int
		min        int
		sec        int
		nsec       int
		z          *time.Location
		zoneOffset int = -1
		zoneName   string
	)

	// Each iteration processes one std value.
	for ***REMOVED***
		var err error
		prefix, std, suffix := nextStdChunk(layout)
		stdstr := layout[len(prefix) : len(layout)-len(suffix)]
		value, err = skip(value, prefix)
		if err != nil ***REMOVED***
			return time.Time***REMOVED******REMOVED***, &time.ParseError***REMOVED***Layout: alayout, Value: avalue, LayoutElem: prefix, ValueElem: value***REMOVED***
		***REMOVED***
		if std == 0 ***REMOVED***
			if len(value) != 0 ***REMOVED***
				return time.Time***REMOVED******REMOVED***, &time.ParseError***REMOVED***Layout: alayout, Value: avalue, ValueElem: value, Message: ": extra text: " + value***REMOVED***
			***REMOVED***
			break
		***REMOVED***
		layout = suffix
		var p string
		switch std & stdMask ***REMOVED***
		case stdYear:
			if len(value) < 2 ***REMOVED***
				err = errBad
				break
			***REMOVED***
			p, value = value[0:2], value[2:]
			year, err = atoi(p)
			if year >= 69 ***REMOVED*** // Unix time starts Dec 31 1969 in some time zones
				year += 1900
			***REMOVED*** else ***REMOVED***
				year += 2000
			***REMOVED***
		case stdLongYear:
			if len(value) >= 7 && (value[0] == '-' || value[0] == '+') ***REMOVED*** // extended year
				neg := value[0] == '-'
				p, value = value[1:7], value[7:]
				year, err = atoi(p)
				if neg ***REMOVED***
					year = -year
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				if len(value) < 4 || !isDigit(value, 0) ***REMOVED***
					err = errBad
					break
				***REMOVED***
				p, value = value[0:4], value[4:]
				year, err = atoi(p)
			***REMOVED***

		case stdMonth:
			month, value, err = lookup(longMonthNames, value)
			if err != nil ***REMOVED***
				month, value, err = lookup(shortMonthNames, value)
			***REMOVED***
			month++
		case stdLongMonth:
			month, value, err = lookup(longMonthNames, value)
			month++
		case stdNumMonth, stdZeroMonth:
			month, value, err = getnum(value, std == stdZeroMonth)
			if month <= 0 || 12 < month ***REMOVED***
				rangeErrString = "month"
			***REMOVED***
		case stdWeekDay:
			// Ignore weekday except for error checking.
			_, value, err = lookup(longDayNames, value)
			if err != nil ***REMOVED***
				_, value, err = lookup(shortDayNames, value)
			***REMOVED***
		case stdLongWeekDay:
			_, value, err = lookup(longDayNames, value)
		case stdDay, stdUnderDay, stdZeroDay:
			if std == stdUnderDay && len(value) > 0 && value[0] == ' ' ***REMOVED***
				value = value[1:]
			***REMOVED***
			day, value, err = getnum(value, false)
			if day < 0 ***REMOVED***
				// Note that we allow any one- or two-digit day here.
				rangeErrString = "day"
			***REMOVED***
		case stdHour:
			hour, value, err = getnum(value, false)
			if hour < 0 || 24 <= hour ***REMOVED***
				rangeErrString = "hour"
			***REMOVED***
		case stdHour12, stdZeroHour12:
			hour, value, err = getnum(value, std == stdZeroHour12)
			if hour < 0 || 12 < hour ***REMOVED***
				rangeErrString = "hour"
			***REMOVED***
		case stdMinute, stdZeroMinute:
			min, value, err = getnum(value, std == stdZeroMinute)
			if min < 0 || 60 <= min ***REMOVED***
				rangeErrString = "minute"
			***REMOVED***
		case stdSecond, stdZeroSecond:
			sec, value, err = getnum(value, std == stdZeroSecond)
			if sec < 0 || 60 <= sec ***REMOVED***
				rangeErrString = "second"
				break
			***REMOVED***
			// Special case: do we have a fractional second but no
			// fractional second in the format?
			if len(value) >= 2 && value[0] == '.' && isDigit(value, 1) ***REMOVED***
				_, std, _ = nextStdChunk(layout)
				std &= stdMask
				if std == stdFracSecond0 || std == stdFracSecond9 ***REMOVED***
					// Fractional second in the layout; proceed normally
					break
				***REMOVED***
				// No fractional second in the layout but we have one in the input.
				n := 2
				for ; n < len(value) && isDigit(value, n); n++ ***REMOVED***
				***REMOVED***
				nsec, rangeErrString, err = parseNanoseconds(value, n)
				value = value[n:]
			***REMOVED***
		case stdPM:
			if len(value) < 2 ***REMOVED***
				err = errBad
				break
			***REMOVED***
			p, value = value[0:2], value[2:]
			switch p ***REMOVED***
			case "PM":
				pmSet = true
			case "AM":
				amSet = true
			default:
				err = errBad
			***REMOVED***
		case stdpm:
			if len(value) < 2 ***REMOVED***
				err = errBad
				break
			***REMOVED***
			p, value = value[0:2], value[2:]
			switch p ***REMOVED***
			case "pm":
				pmSet = true
			case "am":
				amSet = true
			default:
				err = errBad
			***REMOVED***
		case stdISO8601TZ, stdISO8601ColonTZ, stdISO8601SecondsTZ, stdISO8601ShortTZ, stdISO8601ColonSecondsTZ, stdNumTZ, stdNumShortTZ, stdNumColonTZ, stdNumSecondsTz, stdNumColonSecondsTZ:
			if (std == stdISO8601TZ || std == stdISO8601ShortTZ || std == stdISO8601ColonTZ ||
				std == stdISO8601SecondsTZ || std == stdISO8601ColonSecondsTZ) && len(value) >= 1 && value[0] == 'Z' ***REMOVED***

				value = value[1:]
				z = time.UTC
				break
			***REMOVED***
			var sign, hour, min, seconds string
			if std == stdISO8601ColonTZ || std == stdNumColonTZ || std == stdNumTZ || std == stdISO8601TZ ***REMOVED***
				if len(value) < 4 ***REMOVED***
					err = errBad
					break
				***REMOVED***
				if value[3] != ':' ***REMOVED***
					if std == stdNumColonTZ || std == stdISO8601ColonTZ || len(value) < 5 ***REMOVED***
						err = errBad
						break
					***REMOVED***
					sign, hour, min, seconds, value = value[0:1], value[1:3], value[3:5], "00", value[5:]
				***REMOVED*** else ***REMOVED***
					if len(value) < 6 ***REMOVED***
						err = errBad
						break
					***REMOVED***
					sign, hour, min, seconds, value = value[0:1], value[1:3], value[4:6], "00", value[6:]
				***REMOVED***
			***REMOVED*** else if std == stdNumShortTZ || std == stdISO8601ShortTZ ***REMOVED***
				if len(value) < 3 ***REMOVED***
					err = errBad
					break
				***REMOVED***
				sign, hour, min, seconds, value = value[0:1], value[1:3], "00", "00", value[3:]
			***REMOVED*** else if std == stdISO8601ColonSecondsTZ || std == stdNumColonSecondsTZ || std == stdISO8601SecondsTZ || std == stdNumSecondsTz ***REMOVED***
				if len(value) < 7 ***REMOVED***
					err = errBad
					break
				***REMOVED***
				if value[3] != ':' || value[6] != ':' ***REMOVED***
					if std == stdISO8601ColonSecondsTZ || std == stdNumColonSecondsTZ || len(value) < 7 ***REMOVED***
						err = errBad
						break
					***REMOVED***
					sign, hour, min, seconds, value = value[0:1], value[1:3], value[3:5], value[5:7], value[7:]
				***REMOVED*** else ***REMOVED***
					if len(value) < 9 ***REMOVED***
						err = errBad
						break
					***REMOVED***
					sign, hour, min, seconds, value = value[0:1], value[1:3], value[4:6], value[7:9], value[9:]
				***REMOVED***
			***REMOVED***
			var hr, mm, ss int
			hr, err = atoi(hour)
			if err == nil ***REMOVED***
				mm, err = atoi(min)
			***REMOVED***
			if err == nil ***REMOVED***
				ss, err = atoi(seconds)
			***REMOVED***
			zoneOffset = (hr*60+mm)*60 + ss // offset is in seconds
			switch sign[0] ***REMOVED***
			case '+':
			case '-':
				zoneOffset = -zoneOffset
			default:
				err = errBad
			***REMOVED***
		case stdTZ:
			// Does it look like a time zone?
			if len(value) >= 3 && value[0:3] == "UTC" ***REMOVED***
				z = time.UTC
				value = value[3:]
				break
			***REMOVED***
			n, ok := parseTimeZone(value)
			if !ok ***REMOVED***
				err = errBad
				break
			***REMOVED***
			zoneName, value = value[:n], value[n:]
		case stdBracketTZ:
			if len(value) < 3 || value[0] != '(' ***REMOVED***
				err = errBad
				break
			***REMOVED***
			i := 1
			for ; ; i++ ***REMOVED***
				if i >= len(value) ***REMOVED***
					err = errBad
					break
				***REMOVED***
				if value[i] == ')' ***REMOVED***
					zoneName, value = value[1:i], value[i+1:]
					break
				***REMOVED***
			***REMOVED***

		case stdFracSecond0:
			// stdFracSecond0 requires the exact number of digits as specified in
			// the layout.
			ndigit := 1 + (std >> stdArgShift)
			if len(value) < ndigit ***REMOVED***
				err = errBad
				break
			***REMOVED***
			nsec, rangeErrString, err = parseNanoseconds(value, ndigit)
			value = value[ndigit:]

		case stdFracSecond9:
			if len(value) < 2 || value[0] != '.' || value[1] < '0' || '9' < value[1] ***REMOVED***
				// Fractional second omitted.
				break
			***REMOVED***
			// Take any number of digits, even more than asked for,
			// because it is what the stdSecond case would do.
			i := 0
			for i < 9 && i+1 < len(value) && '0' <= value[i+1] && value[i+1] <= '9' ***REMOVED***
				i++
			***REMOVED***
			nsec, rangeErrString, err = parseNanoseconds(value, 1+i)
			value = value[1+i:]
		***REMOVED***
		if rangeErrString != "" ***REMOVED***
			return time.Time***REMOVED******REMOVED***, &time.ParseError***REMOVED***Layout: alayout, Value: avalue, LayoutElem: stdstr, ValueElem: value, Message: ": " + rangeErrString + " out of range"***REMOVED***
		***REMOVED***
		if err != nil ***REMOVED***
			return time.Time***REMOVED******REMOVED***, &time.ParseError***REMOVED***Layout: alayout, Value: avalue, LayoutElem: stdstr, ValueElem: value***REMOVED***
		***REMOVED***
	***REMOVED***
	if pmSet && hour < 12 ***REMOVED***
		hour += 12
	***REMOVED*** else if amSet && hour == 12 ***REMOVED***
		hour = 0
	***REMOVED***

	// Validate the day of the month.
	if day < 1 || day > daysIn(time.Month(month), year) ***REMOVED***
		return time.Time***REMOVED******REMOVED***, &time.ParseError***REMOVED***Layout: alayout, Value: avalue, ValueElem: value, Message: ": day out of range"***REMOVED***
	***REMOVED***

	if z == nil ***REMOVED***
		if zoneOffset == -1 ***REMOVED***
			if zoneName != "" ***REMOVED***
				if z1, err := time.LoadLocation(zoneName); err == nil ***REMOVED***
					z = z1
				***REMOVED*** else ***REMOVED***
					return time.Time***REMOVED******REMOVED***, &time.ParseError***REMOVED***Layout: alayout, Value: avalue, ValueElem: value, Message: ": unknown timezone"***REMOVED***
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				z = defaultLocation
			***REMOVED***
		***REMOVED*** else if zoneOffset == 0 ***REMOVED***
			z = time.UTC
		***REMOVED*** else ***REMOVED***
			z = time.FixedZone("", zoneOffset)
		***REMOVED***
	***REMOVED***

	return time.Date(year, time.Month(month), day, hour, min, sec, nsec, z), nil
***REMOVED***

var errLeadingInt = errors.New("time: bad [0-9]*") // never printed

func signedLeadingInt(s string) (x int64, rem string, err error) ***REMOVED***
	neg := false
	if s != "" && (s[0] == '-' || s[0] == '+') ***REMOVED***
		neg = s[0] == '-'
		s = s[1:]
	***REMOVED***
	x, rem, err = leadingInt(s)
	if err != nil ***REMOVED***
		return
	***REMOVED***

	if neg ***REMOVED***
		x = -x
	***REMOVED***
	return
***REMOVED***

// leadingInt consumes the leading [0-9]* from s.
func leadingInt(s string) (x int64, rem string, err error) ***REMOVED***
	i := 0
	for ; i < len(s); i++ ***REMOVED***
		c := s[i]
		if c < '0' || c > '9' ***REMOVED***
			break
		***REMOVED***
		if x > (1<<63-1)/10 ***REMOVED***
			// overflow
			return 0, "", errLeadingInt
		***REMOVED***
		x = x*10 + int64(c) - '0'
		if x < 0 ***REMOVED***
			// overflow
			return 0, "", errLeadingInt
		***REMOVED***
	***REMOVED***
	return x, s[i:], nil
***REMOVED***

// nextStdChunk finds the first occurrence of a std string in
// layout and returns the text before, the std string, and the text after.
func nextStdChunk(layout string) (prefix string, std int, suffix string) ***REMOVED***
	for i := 0; i < len(layout); i++ ***REMOVED***
		switch c := int(layout[i]); c ***REMOVED***
		case 'J': // January, Jan
			if len(layout) >= i+3 && layout[i:i+3] == "Jan" ***REMOVED***
				if len(layout) >= i+7 && layout[i:i+7] == "January" ***REMOVED***
					return layout[0:i], stdLongMonth, layout[i+7:]
				***REMOVED***
				if !startsWithLowerCase(layout[i+3:]) ***REMOVED***
					return layout[0:i], stdMonth, layout[i+3:]
				***REMOVED***
			***REMOVED***

		case 'M': // Monday, Mon, MST
			if len(layout) >= i+3 ***REMOVED***
				if layout[i:i+3] == "Mon" ***REMOVED***
					if len(layout) >= i+6 && layout[i:i+6] == "Monday" ***REMOVED***
						return layout[0:i], stdLongWeekDay, layout[i+6:]
					***REMOVED***
					if !startsWithLowerCase(layout[i+3:]) ***REMOVED***
						return layout[0:i], stdWeekDay, layout[i+3:]
					***REMOVED***
				***REMOVED***
				if layout[i:i+3] == "MST" ***REMOVED***
					return layout[0:i], stdTZ, layout[i+3:]
				***REMOVED***
			***REMOVED***

		case '0': // 01, 02, 03, 04, 05, 06
			if len(layout) >= i+2 && '1' <= layout[i+1] && layout[i+1] <= '6' ***REMOVED***
				return layout[0:i], std0x[layout[i+1]-'1'], layout[i+2:]
			***REMOVED***

		case '1': // 15, 1
			if len(layout) >= i+2 && layout[i+1] == '5' ***REMOVED***
				return layout[0:i], stdHour, layout[i+2:]
			***REMOVED***
			return layout[0:i], stdNumMonth, layout[i+1:]

		case '2': // 2006, 2
			if len(layout) >= i+4 && layout[i:i+4] == "2006" ***REMOVED***
				return layout[0:i], stdLongYear, layout[i+4:]
			***REMOVED***
			return layout[0:i], stdDay, layout[i+1:]

		case '_': // _2, _2006
			if len(layout) >= i+2 && layout[i+1] == '2' ***REMOVED***
				//_2006 is really a literal _, followed by stdLongYear
				if len(layout) >= i+5 && layout[i+1:i+5] == "2006" ***REMOVED***
					return layout[0 : i+1], stdLongYear, layout[i+5:]
				***REMOVED***
				return layout[0:i], stdUnderDay, layout[i+2:]
			***REMOVED***

		case '3':
			return layout[0:i], stdHour12, layout[i+1:]

		case '4':
			return layout[0:i], stdMinute, layout[i+1:]

		case '5':
			return layout[0:i], stdSecond, layout[i+1:]

		case 'P': // PM
			if len(layout) >= i+2 && layout[i+1] == 'M' ***REMOVED***
				return layout[0:i], stdPM, layout[i+2:]
			***REMOVED***

		case 'p': // pm
			if len(layout) >= i+2 && layout[i+1] == 'm' ***REMOVED***
				return layout[0:i], stdpm, layout[i+2:]
			***REMOVED***

		case '-': // -070000, -07:00:00, -0700, -07:00, -07
			if len(layout) >= i+7 && layout[i:i+7] == "-070000" ***REMOVED***
				return layout[0:i], stdNumSecondsTz, layout[i+7:]
			***REMOVED***
			if len(layout) >= i+9 && layout[i:i+9] == "-07:00:00" ***REMOVED***
				return layout[0:i], stdNumColonSecondsTZ, layout[i+9:]
			***REMOVED***
			if len(layout) >= i+5 && layout[i:i+5] == "-0700" ***REMOVED***
				return layout[0:i], stdNumTZ, layout[i+5:]
			***REMOVED***
			if len(layout) >= i+6 && layout[i:i+6] == "-07:00" ***REMOVED***
				return layout[0:i], stdNumColonTZ, layout[i+6:]
			***REMOVED***
			if len(layout) >= i+3 && layout[i:i+3] == "-07" ***REMOVED***
				return layout[0:i], stdNumShortTZ, layout[i+3:]
			***REMOVED***

		case 'Z': // Z070000, Z07:00:00, Z0700, Z07:00,
			if len(layout) >= i+7 && layout[i:i+7] == "Z070000" ***REMOVED***
				return layout[0:i], stdISO8601SecondsTZ, layout[i+7:]
			***REMOVED***
			if len(layout) >= i+9 && layout[i:i+9] == "Z07:00:00" ***REMOVED***
				return layout[0:i], stdISO8601ColonSecondsTZ, layout[i+9:]
			***REMOVED***
			if len(layout) >= i+5 && layout[i:i+5] == "Z0700" ***REMOVED***
				return layout[0:i], stdISO8601TZ, layout[i+5:]
			***REMOVED***
			if len(layout) >= i+6 && layout[i:i+6] == "Z07:00" ***REMOVED***
				return layout[0:i], stdISO8601ColonTZ, layout[i+6:]
			***REMOVED***
			if len(layout) >= i+3 && layout[i:i+3] == "Z07" ***REMOVED***
				return layout[0:i], stdISO8601ShortTZ, layout[i+3:]
			***REMOVED***

		case '.': // .000 or .999 - repeated digits for fractional seconds.
			if i+1 < len(layout) && (layout[i+1] == '0' || layout[i+1] == '9') ***REMOVED***
				ch := layout[i+1]
				j := i + 1
				for j < len(layout) && layout[j] == ch ***REMOVED***
					j++
				***REMOVED***
				// String of digits must end here - only fractional second is all digits.
				if !isDigit(layout, j) ***REMOVED***
					std := stdFracSecond0
					if layout[i+1] == '9' ***REMOVED***
						std = stdFracSecond9
					***REMOVED***
					std |= (j - (i + 1)) << stdArgShift
					return layout[0:i], std, layout[j:]
				***REMOVED***
			***REMOVED***
		case '(':
			if len(layout) >= i+5 && layout[i:i+5] == "(MST)" ***REMOVED***
				return layout[0:i], stdBracketTZ, layout[i+5:]
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return layout, 0, ""
***REMOVED***

var longDayNames = []string***REMOVED***
	"Sunday",
	"Monday",
	"Tuesday",
	"Wednesday",
	"Thursday",
	"Friday",
	"Saturday",
***REMOVED***

var shortDayNames = []string***REMOVED***
	"Sun",
	"Mon",
	"Tue",
	"Wed",
	"Thu",
	"Fri",
	"Sat",
***REMOVED***

var shortMonthNames = []string***REMOVED***
	"Jan",
	"Feb",
	"Mar",
	"Apr",
	"May",
	"Jun",
	"Jul",
	"Aug",
	"Sep",
	"Oct",
	"Nov",
	"Dec",
***REMOVED***

var longMonthNames = []string***REMOVED***
	"January",
	"February",
	"March",
	"April",
	"May",
	"June",
	"July",
	"August",
	"September",
	"October",
	"November",
	"December",
***REMOVED***

// isDigit reports whether s[i] is in range and is a decimal digit.
func isDigit(s string, i int) bool ***REMOVED***
	if len(s) <= i ***REMOVED***
		return false
	***REMOVED***
	c := s[i]
	return '0' <= c && c <= '9'
***REMOVED***

// getnum parses s[0:1] or s[0:2] (fixed forces the latter)
// as a decimal integer and returns the integer and the
// remainder of the string.
func getnum(s string, fixed bool) (int, string, error) ***REMOVED***
	if !isDigit(s, 0) ***REMOVED***
		return 0, s, errBad
	***REMOVED***
	if !isDigit(s, 1) ***REMOVED***
		if fixed ***REMOVED***
			return 0, s, errBad
		***REMOVED***
		return int(s[0] - '0'), s[1:], nil
	***REMOVED***
	return int(s[0]-'0')*10 + int(s[1]-'0'), s[2:], nil
***REMOVED***

func cutspace(s string) string ***REMOVED***
	for len(s) > 0 && s[0] == ' ' ***REMOVED***
		s = s[1:]
	***REMOVED***
	return s
***REMOVED***

// skip removes the given prefix from value,
// treating runs of space characters as equivalent.
func skip(value, prefix string) (string, error) ***REMOVED***
	for len(prefix) > 0 ***REMOVED***
		if prefix[0] == ' ' ***REMOVED***
			if len(value) > 0 && value[0] != ' ' ***REMOVED***
				return value, errBad
			***REMOVED***
			prefix = cutspace(prefix)
			value = cutspace(value)
			continue
		***REMOVED***
		if len(value) == 0 || value[0] != prefix[0] ***REMOVED***
			return value, errBad
		***REMOVED***
		prefix = prefix[1:]
		value = value[1:]
	***REMOVED***
	return value, nil
***REMOVED***

// Never printed, just needs to be non-nil for return by atoi.
var atoiError = errors.New("time: invalid number")

// Duplicates functionality in strconv, but avoids dependency.
func atoi(s string) (x int, err error) ***REMOVED***
	q, rem, err := signedLeadingInt(s)
	x = int(q)
	if err != nil || rem != "" ***REMOVED***
		return 0, atoiError
	***REMOVED***
	return x, nil
***REMOVED***

// match reports whether s1 and s2 match ignoring case.
// It is assumed s1 and s2 are the same length.
func match(s1, s2 string) bool ***REMOVED***
	for i := 0; i < len(s1); i++ ***REMOVED***
		c1 := s1[i]
		c2 := s2[i]
		if c1 != c2 ***REMOVED***
			// Switch to lower-case; 'a'-'A' is known to be a single bit.
			c1 |= 'a' - 'A'
			c2 |= 'a' - 'A'
			if c1 != c2 || c1 < 'a' || c1 > 'z' ***REMOVED***
				return false
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

func lookup(tab []string, val string) (int, string, error) ***REMOVED***
	for i, v := range tab ***REMOVED***
		if len(val) >= len(v) && match(val[0:len(v)], v) ***REMOVED***
			return i, val[len(v):], nil
		***REMOVED***
	***REMOVED***
	return -1, val, errBad
***REMOVED***

// daysBefore[m] counts the number of days in a non-leap year
// before month m begins. There is an entry for m=12, counting
// the number of days before January of next year (365).
var daysBefore = [...]int32***REMOVED***
	0,
	31,
	31 + 28,
	31 + 28 + 31,
	31 + 28 + 31 + 30,
	31 + 28 + 31 + 30 + 31,
	31 + 28 + 31 + 30 + 31 + 30,
	31 + 28 + 31 + 30 + 31 + 30 + 31,
	31 + 28 + 31 + 30 + 31 + 30 + 31 + 31,
	31 + 28 + 31 + 30 + 31 + 30 + 31 + 31 + 30,
	31 + 28 + 31 + 30 + 31 + 30 + 31 + 31 + 30 + 31,
	31 + 28 + 31 + 30 + 31 + 30 + 31 + 31 + 30 + 31 + 30,
	31 + 28 + 31 + 30 + 31 + 30 + 31 + 31 + 30 + 31 + 30 + 31,
***REMOVED***

func isLeap(year int) bool ***REMOVED***
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
***REMOVED***

func daysIn(m time.Month, year int) int ***REMOVED***
	if m == time.February && isLeap(year) ***REMOVED***
		return 29
	***REMOVED***
	return int(daysBefore[m] - daysBefore[m-1])
***REMOVED***

// parseTimeZone parses a time zone string and returns its length. Time zones
// are human-generated and unpredictable. We can't do precise error checking.
// On the other hand, for a correct parse there must be a time zone at the
// beginning of the string, so it's almost always true that there's one
// there. We look at the beginning of the string for a run of upper-case letters.
// If there are more than 5, it's an error.
// If there are 4 or 5 and the last is a T, it's a time zone.
// If there are 3, it's a time zone.
// Otherwise, other than special cases, it's not a time zone.
// GMT is special because it can have an hour offset.
func parseTimeZone(value string) (length int, ok bool) ***REMOVED***
	if len(value) < 3 ***REMOVED***
		return 0, false
	***REMOVED***
	// Special case 1: ChST and MeST are the only zones with a lower-case letter.
	if len(value) >= 4 && (value[:4] == "ChST" || value[:4] == "MeST") ***REMOVED***
		return 4, true
	***REMOVED***
	// Special case 2: GMT may have an hour offset; treat it specially.
	if value[:3] == "GMT" ***REMOVED***
		length = parseGMT(value)
		return length, true
	***REMOVED***
	// Special Case 3: Some time zones are not named, but have +/-00 format
	if value[0] == '+' || value[0] == '-' ***REMOVED***
		length = parseSignedOffset(value)
		return length, true
	***REMOVED***
	// How many upper-case letters are there? Need at least three, at most five.
	var nUpper int
	for nUpper = 0; nUpper < 6; nUpper++ ***REMOVED***
		if nUpper >= len(value) ***REMOVED***
			break
		***REMOVED***
		if c := value[nUpper]; c < 'A' || 'Z' < c ***REMOVED***
			break
		***REMOVED***
	***REMOVED***
	switch nUpper ***REMOVED***
	case 0, 1, 2, 6:
		return 0, false
	case 5: // Must end in T to match.
		if value[4] == 'T' ***REMOVED***
			return 5, true
		***REMOVED***
	case 4:
		// Must end in T, except one special case.
		if value[3] == 'T' || value[:4] == "WITA" ***REMOVED***
			return 4, true
		***REMOVED***
	case 3:
		return 3, true
	***REMOVED***
	return 0, false
***REMOVED***

// parseGMT parses a GMT time zone. The input string is known to start "GMT".
// The function checks whether that is followed by a sign and a number in the
// range -14 through 12 excluding zero.
func parseGMT(value string) int ***REMOVED***
	value = value[3:]
	if len(value) == 0 ***REMOVED***
		return 3
	***REMOVED***

	return 3 + parseSignedOffset(value)
***REMOVED***

// parseSignedOffset parses a signed timezone offset (e.g. "+03" or "-04").
// The function checks for a signed number in the range -14 through +12 excluding zero.
// Returns length of the found offset string or 0 otherwise
func parseSignedOffset(value string) int ***REMOVED***
	sign := value[0]
	if sign != '-' && sign != '+' ***REMOVED***
		return 0
	***REMOVED***
	x, rem, err := leadingInt(value[1:])
	if err != nil ***REMOVED***
		return 0
	***REMOVED***
	if sign == '-' ***REMOVED***
		x = -x
	***REMOVED***
	if x == 0 || x < -14 || 12 < x ***REMOVED***
		return 0
	***REMOVED***
	return len(value) - len(rem)
***REMOVED***

func parseNanoseconds(value string, nbytes int) (ns int, rangeErrString string, err error) ***REMOVED***
	if value[0] != '.' ***REMOVED***
		err = errBad
		return
	***REMOVED***
	if ns, err = atoi(value[1:nbytes]); err != nil ***REMOVED***
		return
	***REMOVED***
	if ns < 0 || 1e9 <= ns ***REMOVED***
		rangeErrString = "fractional second"
		return
	***REMOVED***
	// We need nanoseconds, which means scaling by the number
	// of missing digits in the format, maximum length 10. If it's
	// longer than 10, we won't scale.
	scaleDigits := 10 - nbytes
	for i := 0; i < scaleDigits; i++ ***REMOVED***
		ns *= 10
	***REMOVED***
	return
***REMOVED***

// std0x records the std values for "01", "02", ..., "06".
var std0x = [...]int***REMOVED***stdZeroMonth, stdZeroDay, stdZeroHour12, stdZeroMinute, stdZeroSecond, stdYear***REMOVED***

// startsWithLowerCase reports whether the string has a lower-case letter at the beginning.
// Its purpose is to prevent matching strings like "Month" when looking for "Mon".
func startsWithLowerCase(str string) bool ***REMOVED***
	if len(str) == 0 ***REMOVED***
		return false
	***REMOVED***
	c := str[0]
	return 'a' <= c && c <= 'z'
***REMOVED***
