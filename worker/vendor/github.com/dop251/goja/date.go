package goja

import (
	"math"
	"reflect"
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

type dateLayoutDesc struct ***REMOVED***
	layout   string
	dateOnly bool
***REMOVED***

var (
	dateLayoutsNumeric = []dateLayoutDesc***REMOVED***
		***REMOVED***layout: "2006-01-02T15:04:05Z0700"***REMOVED***,
		***REMOVED***layout: "2006-01-02T15:04:05"***REMOVED***,
		***REMOVED***layout: "2006-01-02", dateOnly: true***REMOVED***,
		***REMOVED***layout: "2006-01-02 15:04:05"***REMOVED***,

		***REMOVED***layout: "2006", dateOnly: true***REMOVED***,
		***REMOVED***layout: "2006-01", dateOnly: true***REMOVED***,

		***REMOVED***layout: "2006T15:04"***REMOVED***,
		***REMOVED***layout: "2006-01T15:04"***REMOVED***,
		***REMOVED***layout: "2006-01-02T15:04"***REMOVED***,

		***REMOVED***layout: "2006T15:04:05"***REMOVED***,
		***REMOVED***layout: "2006-01T15:04:05"***REMOVED***,

		***REMOVED***layout: "2006T15:04Z0700"***REMOVED***,
		***REMOVED***layout: "2006-01T15:04Z0700"***REMOVED***,
		***REMOVED***layout: "2006-01-02T15:04Z0700"***REMOVED***,

		***REMOVED***layout: "2006T15:04:05Z0700"***REMOVED***,
		***REMOVED***layout: "2006-01T15:04:05Z0700"***REMOVED***,
	***REMOVED***

	dateLayoutsAlpha = []dateLayoutDesc***REMOVED***
		***REMOVED***layout: time.RFC1123***REMOVED***,
		***REMOVED***layout: time.RFC1123Z***REMOVED***,
		***REMOVED***layout: dateTimeLayout***REMOVED***,
		***REMOVED***layout: time.UnixDate***REMOVED***,
		***REMOVED***layout: time.ANSIC***REMOVED***,
		***REMOVED***layout: time.RubyDate***REMOVED***,
		***REMOVED***layout: "Mon, _2 Jan 2006 15:04:05 GMT-0700 (MST)"***REMOVED***,
		***REMOVED***layout: "Mon, _2 Jan 2006 15:04:05 -0700 (MST)"***REMOVED***,
		***REMOVED***layout: "Jan _2, 2006", dateOnly: true***REMOVED***,
	***REMOVED***
)

func dateParse(date string) (time.Time, bool) ***REMOVED***
	var t time.Time
	var err error
	var layouts []dateLayoutDesc
	if len(date) > 0 ***REMOVED***
		first := date[0]
		if first <= '9' && (first >= '0' || first == '-' || first == '+') ***REMOVED***
			layouts = dateLayoutsNumeric
		***REMOVED*** else ***REMOVED***
			layouts = dateLayoutsAlpha
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		return time.Time***REMOVED******REMOVED***, false
	***REMOVED***
	for _, desc := range layouts ***REMOVED***
		var defLoc *time.Location
		if desc.dateOnly ***REMOVED***
			defLoc = time.UTC
		***REMOVED*** else ***REMOVED***
			defLoc = time.Local
		***REMOVED***
		t, err = parseDate(desc.layout, date, defLoc)
		if err == nil ***REMOVED***
			break
		***REMOVED***
	***REMOVED***
	if err != nil ***REMOVED***
		return time.Time***REMOVED******REMOVED***, false
	***REMOVED***
	unix := timeToMsec(t)
	return t, unix >= -maxTime && unix <= maxTime
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

func (d *dateObject) exportType() reflect.Type ***REMOVED***
	return typeTime
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
