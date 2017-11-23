package goja

import (
	"regexp"
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
		"2006",
		"2006-01",
		"2006-01-02",

		"2006T15:04",
		"2006-01T15:04",
		"2006-01-02T15:04",

		"2006T15:04:05",
		"2006-01T15:04:05",
		"2006-01-02T15:04:05",

		"2006T15:04:05.000",
		"2006-01T15:04:05.000",
		"2006-01-02T15:04:05.000",

		"2006T15:04-0700",
		"2006-01T15:04-0700",
		"2006-01-02T15:04-0700",

		"2006T15:04:05-0700",
		"2006-01T15:04:05-0700",
		"2006-01-02T15:04:05-0700",

		"2006T15:04:05.000-0700",
		"2006-01T15:04:05.000-0700",
		"2006-01-02T15:04:05.000-0700",

		time.RFC1123,
		dateTimeLayout,
	***REMOVED***
	matchDateTimeZone = regexp.MustCompile(`^(.*)(?:(Z)|([\+\-]\d***REMOVED***2***REMOVED***):(\d***REMOVED***2***REMOVED***))$`)
)

func dateParse(date string) (time.Time, bool) ***REMOVED***
	// YYYY-MM-DDTHH:mm:ss.sssZ
	var t time.Time
	var err error
	***REMOVED***
		date := date
		if match := matchDateTimeZone.FindStringSubmatch(date); match != nil ***REMOVED***
			if match[2] == "Z" ***REMOVED***
				date = match[1] + "+0000"
			***REMOVED*** else ***REMOVED***
				date = match[1] + match[3] + match[4]
			***REMOVED***
		***REMOVED***
		for _, layout := range dateLayoutList ***REMOVED***
			t, err = time.Parse(layout, date)
			if err == nil ***REMOVED***
				break
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return t, err == nil
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
