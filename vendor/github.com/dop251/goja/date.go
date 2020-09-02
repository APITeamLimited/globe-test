package goja

import (
	"math"
	"time"
)

const (
	dateTimeLayout       = "Mon Jan 02 2006 15:04:05 GMT-0700 (MST)"
	utcDateTimeLayout    = "Mon, 02 Jan 2006 15:04:05 GMT"
	isoDateTimeLayout    = "2006-01-02T15:04:05.000Z"
	dateLayout           = "Mon Jan 02 2006"
	timeLayout           = "15:04:05 GMT-0700 (MST)"
	datetimeLayout_en_GB = "01/02/2006, 15:04:05"
	dateLayout_en_GB     = "01/02/2006"
	timeLayout_en_GB     = "15:04:05"

	maxTime   = 8.64e15
	timeUnset = math.MinInt64
)

type dateObject struct ***REMOVED***
	baseObject
	msec int64
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
	return t, err == nil && unix >= -maxTime && unix <= maxTime
***REMOVED***

func (r *Runtime) newDateObject(t time.Time, isSet bool, proto *Object) *Object ***REMOVED***
	v := &Object***REMOVED***runtime: r***REMOVED***
	d := &dateObject***REMOVED******REMOVED***
	v.self = d
	d.val = v
	d.class = classDate
	d.prototype = proto
	d.extensible = true
	d.init()
	if isSet ***REMOVED***
		d.msec = timeToMsec(t)
	***REMOVED*** else ***REMOVED***
		d.msec = timeUnset
	***REMOVED***
	return v
***REMOVED***

func dateFormat(t time.Time) string ***REMOVED***
	return t.Local().Format(dateTimeLayout)
***REMOVED***

func timeFromMsec(msec int64) time.Time ***REMOVED***
	sec := msec / 1000
	nsec := (msec % 1000) * 1e6
	return time.Unix(sec, nsec)
***REMOVED***

func timeToMsec(t time.Time) int64 ***REMOVED***
	return t.Unix()*1000 + int64(t.Nanosecond())/1e6
***REMOVED***

func (d *dateObject) toPrimitive() Value ***REMOVED***
	return d.toPrimitiveString()
***REMOVED***

func (d *dateObject) export(*objectExportCtx) interface***REMOVED******REMOVED*** ***REMOVED***
	if d.isSet() ***REMOVED***
		return d.time()
	***REMOVED***
	return nil
***REMOVED***

func (d *dateObject) setTimeMs(ms int64) Value ***REMOVED***
	if ms >= 0 && ms <= maxTime || ms < 0 && ms >= -maxTime ***REMOVED***
		d.msec = ms
		return intToValue(ms)
	***REMOVED***

	d.unset()
	return _NaN
***REMOVED***

func (d *dateObject) isSet() bool ***REMOVED***
	return d.msec != timeUnset
***REMOVED***

func (d *dateObject) unset() ***REMOVED***
	d.msec = timeUnset
***REMOVED***

func (d *dateObject) time() time.Time ***REMOVED***
	return timeFromMsec(d.msec)
***REMOVED***

func (d *dateObject) timeUTC() time.Time ***REMOVED***
	return timeFromMsec(d.msec).In(time.UTC)
***REMOVED***
