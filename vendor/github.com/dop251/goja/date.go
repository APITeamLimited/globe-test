package goja

import (
	"time"
)

const (
	dateTimeLayout       = "Mon Jan 02 2006 15:04:05 GMT-0700 (MST)"
	isoDateTimeLayout    = "2006-01-02T15:04:05.000Z"
	dateLayout           = "Mon Jan 02 2006"
	timeLayout           = "15:04:05 GMT-0700 (MST)"
	datetimeLayout_en_GB = "01/02/2006, 15:04:05"
	dateLayout_en_GB     = "01/02/2006"
	timeLayout_en_GB     = "15:04:05"
)

type dateObject struct ***REMOVED***
	baseObject
	time  time.Time
	isSet bool
***REMOVED***

var (
	dateLayoutList = []string***REMOVED***
		"2006-01-02T15:04:05Z0700",
		"2006-01-02T15:04:05",
		"2006-01-02",
		"2006-01-02 15:04:05",
		time.RFC1123,
		time.RFC1123Z,
		dateTimeLayout,
		time.UnixDate,
		time.ANSIC,
		time.RubyDate,
		"Mon, 02 Jan 2006 15:04:05 GMT-0700 (MST)",
		"Mon, 02 Jan 2006 15:04:05 -0700 (MST)",

		"2006",
		"2006-01",

		"2006T15:04",
		"2006-01T15:04",
		"2006-01-02T15:04",

		"2006T15:04:05",
		"2006-01T15:04:05",

		"2006T15:04Z0700",
		"2006-01T15:04Z0700",
		"2006-01-02T15:04Z0700",

		"2006T15:04:05Z0700",
		"2006-01T15:04:05Z0700",
	***REMOVED***
)

func dateParse(date string) (time.Time, bool) ***REMOVED***
	var t time.Time
	var err error
	for _, layout := range dateLayoutList ***REMOVED***
		t, err = parseDate(layout, date, time.UTC)
		if err == nil ***REMOVED***
			break
		***REMOVED***
	***REMOVED***
	unix := timeToMsec(t)
	return t, err == nil && unix >= -8640000000000000 && unix <= 8640000000000000
***REMOVED***

func (r *Runtime) newDateObject(t time.Time, isSet bool) *Object ***REMOVED***
	v := &Object***REMOVED***runtime: r***REMOVED***
	d := &dateObject***REMOVED******REMOVED***
	v.self = d
	d.val = v
	d.class = classDate
	d.prototype = r.global.DatePrototype
	d.extensible = true
	d.init()
	d.time = t.In(time.Local)
	d.isSet = isSet
	return v
***REMOVED***

func dateFormat(t time.Time) string ***REMOVED***
	return t.Local().Format(dateTimeLayout)
***REMOVED***

func (d *dateObject) toPrimitive() Value ***REMOVED***
	return d.toPrimitiveString()
***REMOVED***

func (d *dateObject) export() interface***REMOVED******REMOVED*** ***REMOVED***
	if d.isSet ***REMOVED***
		return d.time
	***REMOVED***
	return nil
***REMOVED***
