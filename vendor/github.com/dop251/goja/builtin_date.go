package goja

import (
	"fmt"
	"math"
	"time"
)

func (r *Runtime) makeDate(args []Value, utc bool) (t time.Time, valid bool) ***REMOVED***
	switch ***REMOVED***
	case len(args) >= 2:
		t = time.Date(1970, time.January, 1, 0, 0, 0, 0, time.Local)
		t, valid = _dateSetYear(t, FunctionCall***REMOVED***Arguments: args***REMOVED***, 0, utc)
	case len(args) == 0:
		t = r.now()
		valid = true
	default: // one argument
		if o, ok := args[0].(*Object); ok ***REMOVED***
			if d, ok := o.self.(*dateObject); ok ***REMOVED***
				t = d.time()
				valid = true
			***REMOVED***
		***REMOVED***
		if !valid ***REMOVED***
			pv := toPrimitive(args[0])
			if val, ok := pv.(valueString); ok ***REMOVED***
				return dateParse(val.String())
			***REMOVED***
			pv = pv.ToNumber()
			var n int64
			if i, ok := pv.(valueInt); ok ***REMOVED***
				n = int64(i)
			***REMOVED*** else if f, ok := pv.(valueFloat); ok ***REMOVED***
				f := float64(f)
				if math.IsNaN(f) || math.IsInf(f, 0) ***REMOVED***
					return
				***REMOVED***
				if math.Abs(f) > maxTime ***REMOVED***
					return
				***REMOVED***
				n = int64(f)
			***REMOVED*** else ***REMOVED***
				n = pv.ToInteger()
			***REMOVED***
			t = timeFromMsec(n)
			valid = true
		***REMOVED***
	***REMOVED***
	if valid ***REMOVED***
		msec := t.Unix()*1000 + int64(t.Nanosecond()/1e6)
		if msec < 0 ***REMOVED***
			msec = -msec
		***REMOVED***
		if msec > maxTime ***REMOVED***
			valid = false
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

func (r *Runtime) newDateTime(args []Value, proto *Object) *Object ***REMOVED***
	t, isSet := r.makeDate(args, false)
	return r.newDateObject(t, isSet, proto)
***REMOVED***

func (r *Runtime) builtin_newDate(args []Value, proto *Object) *Object ***REMOVED***
	return r.newDateTime(args, proto)
***REMOVED***

func (r *Runtime) builtin_date(FunctionCall) Value ***REMOVED***
	return asciiString(dateFormat(r.now()))
***REMOVED***

func (r *Runtime) date_parse(call FunctionCall) Value ***REMOVED***
	t, set := dateParse(call.Argument(0).toString().String())
	if set ***REMOVED***
		return intToValue(timeToMsec(t))
	***REMOVED***
	return _NaN
***REMOVED***

func (r *Runtime) date_UTC(call FunctionCall) Value ***REMOVED***
	var args []Value
	if len(call.Arguments) < 2 ***REMOVED***
		args = []Value***REMOVED***call.Argument(0), _positiveZero***REMOVED***
	***REMOVED*** else ***REMOVED***
		args = call.Arguments
	***REMOVED***
	t, valid := r.makeDate(args, true)
	if !valid ***REMOVED***
		return _NaN
	***REMOVED***
	return intToValue(timeToMsec(t))
***REMOVED***

func (r *Runtime) date_now(FunctionCall) Value ***REMOVED***
	return intToValue(timeToMsec(r.now()))
***REMOVED***

func (r *Runtime) dateproto_toString(call FunctionCall) Value ***REMOVED***
	obj := r.toObject(call.This)
	if d, ok := obj.self.(*dateObject); ok ***REMOVED***
		if d.isSet() ***REMOVED***
			return asciiString(d.time().Format(dateTimeLayout))
		***REMOVED*** else ***REMOVED***
			return stringInvalidDate
		***REMOVED***
	***REMOVED***
	panic(r.NewTypeError("Method Date.prototype.toString is called on incompatible receiver"))
***REMOVED***

func (r *Runtime) dateproto_toUTCString(call FunctionCall) Value ***REMOVED***
	obj := r.toObject(call.This)
	if d, ok := obj.self.(*dateObject); ok ***REMOVED***
		if d.isSet() ***REMOVED***
			return asciiString(d.timeUTC().Format(utcDateTimeLayout))
		***REMOVED*** else ***REMOVED***
			return stringInvalidDate
		***REMOVED***
	***REMOVED***
	panic(r.NewTypeError("Method Date.prototype.toUTCString is called on incompatible receiver"))
***REMOVED***

func (r *Runtime) dateproto_toISOString(call FunctionCall) Value ***REMOVED***
	obj := r.toObject(call.This)
	if d, ok := obj.self.(*dateObject); ok ***REMOVED***
		if d.isSet() ***REMOVED***
			utc := d.timeUTC()
			year := utc.Year()
			if year >= -9999 && year <= 9999 ***REMOVED***
				return asciiString(utc.Format(isoDateTimeLayout))
			***REMOVED***
			// extended year
			return asciiString(fmt.Sprintf("%+06d-", year) + utc.Format(isoDateTimeLayout[5:]))
		***REMOVED*** else ***REMOVED***
			panic(r.newError(r.global.RangeError, "Invalid time value"))
		***REMOVED***
	***REMOVED***
	panic(r.NewTypeError("Method Date.prototype.toISOString is called on incompatible receiver"))
***REMOVED***

func (r *Runtime) dateproto_toJSON(call FunctionCall) Value ***REMOVED***
	obj := r.toObject(call.This)
	tv := obj.toPrimitiveNumber()
	if f, ok := tv.(valueFloat); ok ***REMOVED***
		f := float64(f)
		if math.IsNaN(f) || math.IsInf(f, 0) ***REMOVED***
			return _null
		***REMOVED***
	***REMOVED***

	if toISO, ok := obj.self.getStr("toISOString", nil).(*Object); ok ***REMOVED***
		if toISO, ok := toISO.self.assertCallable(); ok ***REMOVED***
			return toISO(FunctionCall***REMOVED***
				This: obj,
			***REMOVED***)
		***REMOVED***
	***REMOVED***

	panic(r.NewTypeError("toISOString is not a function"))
***REMOVED***

func (r *Runtime) dateproto_toPrimitive(call FunctionCall) Value ***REMOVED***
	o := r.toObject(call.This)
	arg := call.Argument(0)

	if asciiString("string").StrictEquals(arg) || asciiString("default").StrictEquals(arg) ***REMOVED***
		return o.self.toPrimitiveString()
	***REMOVED***
	if asciiString("number").StrictEquals(arg) ***REMOVED***
		return o.self.toPrimitiveNumber()
	***REMOVED***
	panic(r.NewTypeError("Invalid hint: %s", arg))
***REMOVED***

func (r *Runtime) dateproto_toDateString(call FunctionCall) Value ***REMOVED***
	obj := r.toObject(call.This)
	if d, ok := obj.self.(*dateObject); ok ***REMOVED***
		if d.isSet() ***REMOVED***
			return asciiString(d.time().Format(dateLayout))
		***REMOVED*** else ***REMOVED***
			return stringInvalidDate
		***REMOVED***
	***REMOVED***
	panic(r.NewTypeError("Method Date.prototype.toDateString is called on incompatible receiver"))
***REMOVED***

func (r *Runtime) dateproto_toTimeString(call FunctionCall) Value ***REMOVED***
	obj := r.toObject(call.This)
	if d, ok := obj.self.(*dateObject); ok ***REMOVED***
		if d.isSet() ***REMOVED***
			return asciiString(d.time().Format(timeLayout))
		***REMOVED*** else ***REMOVED***
			return stringInvalidDate
		***REMOVED***
	***REMOVED***
	panic(r.NewTypeError("Method Date.prototype.toTimeString is called on incompatible receiver"))
***REMOVED***

func (r *Runtime) dateproto_toLocaleString(call FunctionCall) Value ***REMOVED***
	obj := r.toObject(call.This)
	if d, ok := obj.self.(*dateObject); ok ***REMOVED***
		if d.isSet() ***REMOVED***
			return asciiString(d.time().Format(datetimeLayout_en_GB))
		***REMOVED*** else ***REMOVED***
			return stringInvalidDate
		***REMOVED***
	***REMOVED***
	panic(r.NewTypeError("Method Date.prototype.toLocaleString is called on incompatible receiver"))
***REMOVED***

func (r *Runtime) dateproto_toLocaleDateString(call FunctionCall) Value ***REMOVED***
	obj := r.toObject(call.This)
	if d, ok := obj.self.(*dateObject); ok ***REMOVED***
		if d.isSet() ***REMOVED***
			return asciiString(d.time().Format(dateLayout_en_GB))
		***REMOVED*** else ***REMOVED***
			return stringInvalidDate
		***REMOVED***
	***REMOVED***
	panic(r.NewTypeError("Method Date.prototype.toLocaleDateString is called on incompatible receiver"))
***REMOVED***

func (r *Runtime) dateproto_toLocaleTimeString(call FunctionCall) Value ***REMOVED***
	obj := r.toObject(call.This)
	if d, ok := obj.self.(*dateObject); ok ***REMOVED***
		if d.isSet() ***REMOVED***
			return asciiString(d.time().Format(timeLayout_en_GB))
		***REMOVED*** else ***REMOVED***
			return stringInvalidDate
		***REMOVED***
	***REMOVED***
	panic(r.NewTypeError("Method Date.prototype.toLocaleTimeString is called on incompatible receiver"))
***REMOVED***

func (r *Runtime) dateproto_valueOf(call FunctionCall) Value ***REMOVED***
	obj := r.toObject(call.This)
	if d, ok := obj.self.(*dateObject); ok ***REMOVED***
		if d.isSet() ***REMOVED***
			return intToValue(d.msec)
		***REMOVED*** else ***REMOVED***
			return _NaN
		***REMOVED***
	***REMOVED***
	panic(r.NewTypeError("Method Date.prototype.valueOf is called on incompatible receiver"))
***REMOVED***

func (r *Runtime) dateproto_getTime(call FunctionCall) Value ***REMOVED***
	obj := r.toObject(call.This)
	if d, ok := obj.self.(*dateObject); ok ***REMOVED***
		if d.isSet() ***REMOVED***
			return intToValue(d.msec)
		***REMOVED*** else ***REMOVED***
			return _NaN
		***REMOVED***
	***REMOVED***
	panic(r.NewTypeError("Method Date.prototype.getTime is called on incompatible receiver"))
***REMOVED***

func (r *Runtime) dateproto_getFullYear(call FunctionCall) Value ***REMOVED***
	obj := r.toObject(call.This)
	if d, ok := obj.self.(*dateObject); ok ***REMOVED***
		if d.isSet() ***REMOVED***
			return intToValue(int64(d.time().Year()))
		***REMOVED*** else ***REMOVED***
			return _NaN
		***REMOVED***
	***REMOVED***
	panic(r.NewTypeError("Method Date.prototype.getFullYear is called on incompatible receiver"))
***REMOVED***

func (r *Runtime) dateproto_getUTCFullYear(call FunctionCall) Value ***REMOVED***
	obj := r.toObject(call.This)
	if d, ok := obj.self.(*dateObject); ok ***REMOVED***
		if d.isSet() ***REMOVED***
			return intToValue(int64(d.timeUTC().Year()))
		***REMOVED*** else ***REMOVED***
			return _NaN
		***REMOVED***
	***REMOVED***
	panic(r.NewTypeError("Method Date.prototype.getUTCFullYear is called on incompatible receiver"))
***REMOVED***

func (r *Runtime) dateproto_getMonth(call FunctionCall) Value ***REMOVED***
	obj := r.toObject(call.This)
	if d, ok := obj.self.(*dateObject); ok ***REMOVED***
		if d.isSet() ***REMOVED***
			return intToValue(int64(d.time().Month()) - 1)
		***REMOVED*** else ***REMOVED***
			return _NaN
		***REMOVED***
	***REMOVED***
	panic(r.NewTypeError("Method Date.prototype.getMonth is called on incompatible receiver"))
***REMOVED***

func (r *Runtime) dateproto_getUTCMonth(call FunctionCall) Value ***REMOVED***
	obj := r.toObject(call.This)
	if d, ok := obj.self.(*dateObject); ok ***REMOVED***
		if d.isSet() ***REMOVED***
			return intToValue(int64(d.timeUTC().Month()) - 1)
		***REMOVED*** else ***REMOVED***
			return _NaN
		***REMOVED***
	***REMOVED***
	panic(r.NewTypeError("Method Date.prototype.getUTCMonth is called on incompatible receiver"))
***REMOVED***

func (r *Runtime) dateproto_getHours(call FunctionCall) Value ***REMOVED***
	obj := r.toObject(call.This)
	if d, ok := obj.self.(*dateObject); ok ***REMOVED***
		if d.isSet() ***REMOVED***
			return intToValue(int64(d.time().Hour()))
		***REMOVED*** else ***REMOVED***
			return _NaN
		***REMOVED***
	***REMOVED***
	panic(r.NewTypeError("Method Date.prototype.getHours is called on incompatible receiver"))
***REMOVED***

func (r *Runtime) dateproto_getUTCHours(call FunctionCall) Value ***REMOVED***
	obj := r.toObject(call.This)
	if d, ok := obj.self.(*dateObject); ok ***REMOVED***
		if d.isSet() ***REMOVED***
			return intToValue(int64(d.timeUTC().Hour()))
		***REMOVED*** else ***REMOVED***
			return _NaN
		***REMOVED***
	***REMOVED***
	panic(r.NewTypeError("Method Date.prototype.getUTCHours is called on incompatible receiver"))
***REMOVED***

func (r *Runtime) dateproto_getDate(call FunctionCall) Value ***REMOVED***
	obj := r.toObject(call.This)
	if d, ok := obj.self.(*dateObject); ok ***REMOVED***
		if d.isSet() ***REMOVED***
			return intToValue(int64(d.time().Day()))
		***REMOVED*** else ***REMOVED***
			return _NaN
		***REMOVED***
	***REMOVED***
	panic(r.NewTypeError("Method Date.prototype.getDate is called on incompatible receiver"))
***REMOVED***

func (r *Runtime) dateproto_getUTCDate(call FunctionCall) Value ***REMOVED***
	obj := r.toObject(call.This)
	if d, ok := obj.self.(*dateObject); ok ***REMOVED***
		if d.isSet() ***REMOVED***
			return intToValue(int64(d.timeUTC().Day()))
		***REMOVED*** else ***REMOVED***
			return _NaN
		***REMOVED***
	***REMOVED***
	panic(r.NewTypeError("Method Date.prototype.getUTCDate is called on incompatible receiver"))
***REMOVED***

func (r *Runtime) dateproto_getDay(call FunctionCall) Value ***REMOVED***
	obj := r.toObject(call.This)
	if d, ok := obj.self.(*dateObject); ok ***REMOVED***
		if d.isSet() ***REMOVED***
			return intToValue(int64(d.time().Weekday()))
		***REMOVED*** else ***REMOVED***
			return _NaN
		***REMOVED***
	***REMOVED***
	panic(r.NewTypeError("Method Date.prototype.getDay is called on incompatible receiver"))
***REMOVED***

func (r *Runtime) dateproto_getUTCDay(call FunctionCall) Value ***REMOVED***
	obj := r.toObject(call.This)
	if d, ok := obj.self.(*dateObject); ok ***REMOVED***
		if d.isSet() ***REMOVED***
			return intToValue(int64(d.timeUTC().Weekday()))
		***REMOVED*** else ***REMOVED***
			return _NaN
		***REMOVED***
	***REMOVED***
	panic(r.NewTypeError("Method Date.prototype.getUTCDay is called on incompatible receiver"))
***REMOVED***

func (r *Runtime) dateproto_getMinutes(call FunctionCall) Value ***REMOVED***
	obj := r.toObject(call.This)
	if d, ok := obj.self.(*dateObject); ok ***REMOVED***
		if d.isSet() ***REMOVED***
			return intToValue(int64(d.time().Minute()))
		***REMOVED*** else ***REMOVED***
			return _NaN
		***REMOVED***
	***REMOVED***
	panic(r.NewTypeError("Method Date.prototype.getMinutes is called on incompatible receiver"))
***REMOVED***

func (r *Runtime) dateproto_getUTCMinutes(call FunctionCall) Value ***REMOVED***
	obj := r.toObject(call.This)
	if d, ok := obj.self.(*dateObject); ok ***REMOVED***
		if d.isSet() ***REMOVED***
			return intToValue(int64(d.timeUTC().Minute()))
		***REMOVED*** else ***REMOVED***
			return _NaN
		***REMOVED***
	***REMOVED***
	panic(r.NewTypeError("Method Date.prototype.getUTCMinutes is called on incompatible receiver"))
***REMOVED***

func (r *Runtime) dateproto_getSeconds(call FunctionCall) Value ***REMOVED***
	obj := r.toObject(call.This)
	if d, ok := obj.self.(*dateObject); ok ***REMOVED***
		if d.isSet() ***REMOVED***
			return intToValue(int64(d.time().Second()))
		***REMOVED*** else ***REMOVED***
			return _NaN
		***REMOVED***
	***REMOVED***
	panic(r.NewTypeError("Method Date.prototype.getSeconds is called on incompatible receiver"))
***REMOVED***

func (r *Runtime) dateproto_getUTCSeconds(call FunctionCall) Value ***REMOVED***
	obj := r.toObject(call.This)
	if d, ok := obj.self.(*dateObject); ok ***REMOVED***
		if d.isSet() ***REMOVED***
			return intToValue(int64(d.timeUTC().Second()))
		***REMOVED*** else ***REMOVED***
			return _NaN
		***REMOVED***
	***REMOVED***
	panic(r.NewTypeError("Method Date.prototype.getUTCSeconds is called on incompatible receiver"))
***REMOVED***

func (r *Runtime) dateproto_getMilliseconds(call FunctionCall) Value ***REMOVED***
	obj := r.toObject(call.This)
	if d, ok := obj.self.(*dateObject); ok ***REMOVED***
		if d.isSet() ***REMOVED***
			return intToValue(int64(d.time().Nanosecond() / 1e6))
		***REMOVED*** else ***REMOVED***
			return _NaN
		***REMOVED***
	***REMOVED***
	panic(r.NewTypeError("Method Date.prototype.getMilliseconds is called on incompatible receiver"))
***REMOVED***

func (r *Runtime) dateproto_getUTCMilliseconds(call FunctionCall) Value ***REMOVED***
	obj := r.toObject(call.This)
	if d, ok := obj.self.(*dateObject); ok ***REMOVED***
		if d.isSet() ***REMOVED***
			return intToValue(int64(d.timeUTC().Nanosecond() / 1e6))
		***REMOVED*** else ***REMOVED***
			return _NaN
		***REMOVED***
	***REMOVED***
	panic(r.NewTypeError("Method Date.prototype.getUTCMilliseconds is called on incompatible receiver"))
***REMOVED***

func (r *Runtime) dateproto_getTimezoneOffset(call FunctionCall) Value ***REMOVED***
	obj := r.toObject(call.This)
	if d, ok := obj.self.(*dateObject); ok ***REMOVED***
		if d.isSet() ***REMOVED***
			_, offset := d.time().Zone()
			return floatToValue(float64(-offset) / 60)
		***REMOVED*** else ***REMOVED***
			return _NaN
		***REMOVED***
	***REMOVED***
	panic(r.NewTypeError("Method Date.prototype.getTimezoneOffset is called on incompatible receiver"))
***REMOVED***

func (r *Runtime) dateproto_setTime(call FunctionCall) Value ***REMOVED***
	obj := r.toObject(call.This)
	if d, ok := obj.self.(*dateObject); ok ***REMOVED***
		n := call.Argument(0).ToNumber()
		if IsNaN(n) ***REMOVED***
			d.unset()
			return _NaN
		***REMOVED***
		return d.setTimeMs(n.ToInteger())
	***REMOVED***
	panic(r.NewTypeError("Method Date.prototype.setTime is called on incompatible receiver"))
***REMOVED***

// _norm returns nhi, nlo such that
//	hi * base + lo == nhi * base + nlo
//	0 <= nlo < base
func _norm(hi, lo, base int64) (nhi, nlo int64, ok bool) ***REMOVED***
	if lo < 0 ***REMOVED***
		if hi == math.MinInt64 && lo <= -base ***REMOVED***
			// underflow
			ok = false
			return
		***REMOVED***
		n := (-lo-1)/base + 1
		hi -= n
		lo += n * base
	***REMOVED***
	if lo >= base ***REMOVED***
		if hi == math.MaxInt64 ***REMOVED***
			// overflow
			ok = false
			return
		***REMOVED***
		n := lo / base
		hi += n
		lo -= n * base
	***REMOVED***
	return hi, lo, true
***REMOVED***

func mkTime(year, m, day, hour, min, sec, nsec int64, loc *time.Location) (t time.Time, ok bool) ***REMOVED***
	year, m, ok = _norm(year, m, 12)
	if !ok ***REMOVED***
		return
	***REMOVED***

	// Normalise nsec, sec, min, hour, overflowing into day.
	sec, nsec, ok = _norm(sec, nsec, 1e9)
	if !ok ***REMOVED***
		return
	***REMOVED***
	min, sec, ok = _norm(min, sec, 60)
	if !ok ***REMOVED***
		return
	***REMOVED***
	hour, min, ok = _norm(hour, min, 60)
	if !ok ***REMOVED***
		return
	***REMOVED***
	day, hour, ok = _norm(day, hour, 24)
	if !ok ***REMOVED***
		return
	***REMOVED***
	if year > math.MaxInt32 || year < math.MinInt32 ||
		day > math.MaxInt32 || day < math.MinInt32 ||
		m >= math.MaxInt32 || m < math.MinInt32-1 ***REMOVED***
		return time.Time***REMOVED******REMOVED***, false
	***REMOVED***
	month := time.Month(m) + 1
	return time.Date(int(year), month, int(day), int(hour), int(min), int(sec), int(nsec), loc), true
***REMOVED***

func _intArg(call FunctionCall, argNum int) (int64, bool) ***REMOVED***
	n := call.Argument(argNum).ToNumber()
	if IsNaN(n) ***REMOVED***
		return 0, false
	***REMOVED***
	return n.ToInteger(), true
***REMOVED***

func _dateSetYear(t time.Time, call FunctionCall, argNum int, utc bool) (time.Time, bool) ***REMOVED***
	var year int64
	if argNum == 0 || argNum > 0 && argNum < len(call.Arguments) ***REMOVED***
		var ok bool
		year, ok = _intArg(call, argNum)
		if !ok ***REMOVED***
			return time.Time***REMOVED******REMOVED***, false
		***REMOVED***
		if year >= 0 && year <= 99 ***REMOVED***
			year += 1900
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		year = int64(t.Year())
	***REMOVED***

	return _dateSetMonth(year, t, call, argNum+1, utc)
***REMOVED***

func _dateSetFullYear(t time.Time, call FunctionCall, argNum int, utc bool) (time.Time, bool) ***REMOVED***
	var year int64
	if argNum == 0 || argNum > 0 && argNum < len(call.Arguments) ***REMOVED***
		var ok bool
		year, ok = _intArg(call, argNum)
		if !ok ***REMOVED***
			return time.Time***REMOVED******REMOVED***, false
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		year = int64(t.Year())
	***REMOVED***
	return _dateSetMonth(year, t, call, argNum+1, utc)
***REMOVED***

func _dateSetMonth(year int64, t time.Time, call FunctionCall, argNum int, utc bool) (time.Time, bool) ***REMOVED***
	var mon int64
	if argNum == 0 || argNum > 0 && argNum < len(call.Arguments) ***REMOVED***
		var ok bool
		mon, ok = _intArg(call, argNum)
		if !ok ***REMOVED***
			return time.Time***REMOVED******REMOVED***, false
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		mon = int64(t.Month()) - 1
	***REMOVED***

	return _dateSetDay(year, mon, t, call, argNum+1, utc)
***REMOVED***

func _dateSetDay(year, mon int64, t time.Time, call FunctionCall, argNum int, utc bool) (time.Time, bool) ***REMOVED***
	var day int64
	if argNum == 0 || argNum > 0 && argNum < len(call.Arguments) ***REMOVED***
		var ok bool
		day, ok = _intArg(call, argNum)
		if !ok ***REMOVED***
			return time.Time***REMOVED******REMOVED***, false
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		day = int64(t.Day())
	***REMOVED***

	return _dateSetHours(year, mon, day, t, call, argNum+1, utc)
***REMOVED***

func _dateSetHours(year, mon, day int64, t time.Time, call FunctionCall, argNum int, utc bool) (time.Time, bool) ***REMOVED***
	var hours int64
	if argNum == 0 || argNum > 0 && argNum < len(call.Arguments) ***REMOVED***
		var ok bool
		hours, ok = _intArg(call, argNum)
		if !ok ***REMOVED***
			return time.Time***REMOVED******REMOVED***, false
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		hours = int64(t.Hour())
	***REMOVED***
	return _dateSetMinutes(year, mon, day, hours, t, call, argNum+1, utc)
***REMOVED***

func _dateSetMinutes(year, mon, day, hours int64, t time.Time, call FunctionCall, argNum int, utc bool) (time.Time, bool) ***REMOVED***
	var min int64
	if argNum == 0 || argNum > 0 && argNum < len(call.Arguments) ***REMOVED***
		var ok bool
		min, ok = _intArg(call, argNum)
		if !ok ***REMOVED***
			return time.Time***REMOVED******REMOVED***, false
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		min = int64(t.Minute())
	***REMOVED***
	return _dateSetSeconds(year, mon, day, hours, min, t, call, argNum+1, utc)
***REMOVED***

func _dateSetSeconds(year, mon, day, hours, min int64, t time.Time, call FunctionCall, argNum int, utc bool) (time.Time, bool) ***REMOVED***
	var sec int64
	if argNum == 0 || argNum > 0 && argNum < len(call.Arguments) ***REMOVED***
		var ok bool
		sec, ok = _intArg(call, argNum)
		if !ok ***REMOVED***
			return time.Time***REMOVED******REMOVED***, false
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		sec = int64(t.Second())
	***REMOVED***
	return _dateSetMilliseconds(year, mon, day, hours, min, sec, t, call, argNum+1, utc)
***REMOVED***

func _dateSetMilliseconds(year, mon, day, hours, min, sec int64, t time.Time, call FunctionCall, argNum int, utc bool) (time.Time, bool) ***REMOVED***
	var msec int64
	if argNum == 0 || argNum > 0 && argNum < len(call.Arguments) ***REMOVED***
		var ok bool
		msec, ok = _intArg(call, argNum)
		if !ok ***REMOVED***
			return time.Time***REMOVED******REMOVED***, false
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		msec = int64(t.Nanosecond() / 1e6)
	***REMOVED***
	var ok bool
	sec, msec, ok = _norm(sec, msec, 1e3)
	if !ok ***REMOVED***
		return time.Time***REMOVED******REMOVED***, false
	***REMOVED***

	var loc *time.Location
	if utc ***REMOVED***
		loc = time.UTC
	***REMOVED*** else ***REMOVED***
		loc = time.Local
	***REMOVED***
	r, ok := mkTime(year, mon, day, hours, min, sec, msec*1e6, loc)
	if !ok ***REMOVED***
		return time.Time***REMOVED******REMOVED***, false
	***REMOVED***
	if utc ***REMOVED***
		return r.In(time.Local), true
	***REMOVED***
	return r, true
***REMOVED***

func (r *Runtime) dateproto_setMilliseconds(call FunctionCall) Value ***REMOVED***
	obj := r.toObject(call.This)
	if d, ok := obj.self.(*dateObject); ok ***REMOVED***
		if d.isSet() ***REMOVED***
			n := call.Argument(0).ToNumber()
			if IsNaN(n) ***REMOVED***
				d.unset()
				return _NaN
			***REMOVED***
			msec := n.ToInteger()
			sec := d.msec / 1e3
			var ok bool
			sec, msec, ok = _norm(sec, msec, 1e3)
			if !ok ***REMOVED***
				d.unset()
				return _NaN
			***REMOVED***
			return d.setTimeMs(sec*1e3 + msec)
		***REMOVED*** else ***REMOVED***
			return _NaN
		***REMOVED***
	***REMOVED***
	panic(r.NewTypeError("Method Date.prototype.setMilliseconds is called on incompatible receiver"))
***REMOVED***

func (r *Runtime) dateproto_setUTCMilliseconds(call FunctionCall) Value ***REMOVED***
	obj := r.toObject(call.This)
	if d, ok := obj.self.(*dateObject); ok ***REMOVED***
		if d.isSet() ***REMOVED***
			n := call.Argument(0).ToNumber()
			if IsNaN(n) ***REMOVED***
				d.unset()
				return _NaN
			***REMOVED***
			msec := n.ToInteger()
			sec := d.msec / 1e3
			var ok bool
			sec, msec, ok = _norm(sec, msec, 1e3)
			if !ok ***REMOVED***
				d.unset()
				return _NaN
			***REMOVED***
			return d.setTimeMs(sec*1e3 + msec)
		***REMOVED*** else ***REMOVED***
			return _NaN
		***REMOVED***
	***REMOVED***
	panic(r.NewTypeError("Method Date.prototype.setUTCMilliseconds is called on incompatible receiver"))
***REMOVED***

func (r *Runtime) dateproto_setSeconds(call FunctionCall) Value ***REMOVED***
	obj := r.toObject(call.This)
	if d, ok := obj.self.(*dateObject); ok ***REMOVED***
		if d.isSet() ***REMOVED***
			t, ok := _dateSetFullYear(d.time(), call, -5, false)
			if !ok ***REMOVED***
				d.unset()
				return _NaN
			***REMOVED***
			return d.setTimeMs(timeToMsec(t))
		***REMOVED*** else ***REMOVED***
			return _NaN
		***REMOVED***
	***REMOVED***
	panic(r.NewTypeError("Method Date.prototype.setSeconds is called on incompatible receiver"))
***REMOVED***

func (r *Runtime) dateproto_setUTCSeconds(call FunctionCall) Value ***REMOVED***
	obj := r.toObject(call.This)
	if d, ok := obj.self.(*dateObject); ok ***REMOVED***
		if d.isSet() ***REMOVED***
			t, ok := _dateSetFullYear(d.timeUTC(), call, -5, true)
			if !ok ***REMOVED***
				d.unset()
				return _NaN
			***REMOVED***
			return d.setTimeMs(timeToMsec(t))
		***REMOVED*** else ***REMOVED***
			return _NaN
		***REMOVED***
	***REMOVED***
	panic(r.NewTypeError("Method Date.prototype.setUTCSeconds is called on incompatible receiver"))
***REMOVED***

func (r *Runtime) dateproto_setMinutes(call FunctionCall) Value ***REMOVED***
	obj := r.toObject(call.This)
	if d, ok := obj.self.(*dateObject); ok ***REMOVED***
		if d.isSet() ***REMOVED***
			t, ok := _dateSetFullYear(d.time(), call, -4, false)
			if !ok ***REMOVED***
				d.unset()
				return _NaN
			***REMOVED***
			return d.setTimeMs(timeToMsec(t))
		***REMOVED*** else ***REMOVED***
			return _NaN
		***REMOVED***
	***REMOVED***
	panic(r.NewTypeError("Method Date.prototype.setMinutes is called on incompatible receiver"))
***REMOVED***

func (r *Runtime) dateproto_setUTCMinutes(call FunctionCall) Value ***REMOVED***
	obj := r.toObject(call.This)
	if d, ok := obj.self.(*dateObject); ok ***REMOVED***
		if d.isSet() ***REMOVED***
			t, ok := _dateSetFullYear(d.timeUTC(), call, -4, true)
			if !ok ***REMOVED***
				d.unset()
				return _NaN
			***REMOVED***
			return d.setTimeMs(timeToMsec(t))
		***REMOVED*** else ***REMOVED***
			return _NaN
		***REMOVED***
	***REMOVED***
	panic(r.NewTypeError("Method Date.prototype.setUTCMinutes is called on incompatible receiver"))
***REMOVED***

func (r *Runtime) dateproto_setHours(call FunctionCall) Value ***REMOVED***
	obj := r.toObject(call.This)
	if d, ok := obj.self.(*dateObject); ok ***REMOVED***
		if d.isSet() ***REMOVED***
			t, ok := _dateSetFullYear(d.time(), call, -3, false)
			if !ok ***REMOVED***
				d.unset()
				return _NaN
			***REMOVED***
			return d.setTimeMs(timeToMsec(t))
		***REMOVED*** else ***REMOVED***
			return _NaN
		***REMOVED***
	***REMOVED***
	panic(r.NewTypeError("Method Date.prototype.setHours is called on incompatible receiver"))
***REMOVED***

func (r *Runtime) dateproto_setUTCHours(call FunctionCall) Value ***REMOVED***
	obj := r.toObject(call.This)
	if d, ok := obj.self.(*dateObject); ok ***REMOVED***
		if d.isSet() ***REMOVED***
			t, ok := _dateSetFullYear(d.timeUTC(), call, -3, true)
			if !ok ***REMOVED***
				d.unset()
				return _NaN
			***REMOVED***
			return d.setTimeMs(timeToMsec(t))
		***REMOVED*** else ***REMOVED***
			return _NaN
		***REMOVED***
	***REMOVED***
	panic(r.NewTypeError("Method Date.prototype.setUTCHours is called on incompatible receiver"))
***REMOVED***

func (r *Runtime) dateproto_setDate(call FunctionCall) Value ***REMOVED***
	obj := r.toObject(call.This)
	if d, ok := obj.self.(*dateObject); ok ***REMOVED***
		if d.isSet() ***REMOVED***
			t, ok := _dateSetFullYear(d.time(), limitCallArgs(call, 1), -2, false)
			if !ok ***REMOVED***
				d.unset()
				return _NaN
			***REMOVED***
			return d.setTimeMs(timeToMsec(t))
		***REMOVED*** else ***REMOVED***
			return _NaN
		***REMOVED***
	***REMOVED***
	panic(r.NewTypeError("Method Date.prototype.setDate is called on incompatible receiver"))
***REMOVED***

func (r *Runtime) dateproto_setUTCDate(call FunctionCall) Value ***REMOVED***
	obj := r.toObject(call.This)
	if d, ok := obj.self.(*dateObject); ok ***REMOVED***
		if d.isSet() ***REMOVED***
			t, ok := _dateSetFullYear(d.timeUTC(), limitCallArgs(call, 1), -2, true)
			if !ok ***REMOVED***
				d.unset()
				return _NaN
			***REMOVED***
			return d.setTimeMs(timeToMsec(t))
		***REMOVED*** else ***REMOVED***
			return _NaN
		***REMOVED***
	***REMOVED***
	panic(r.NewTypeError("Method Date.prototype.setUTCDate is called on incompatible receiver"))
***REMOVED***

func (r *Runtime) dateproto_setMonth(call FunctionCall) Value ***REMOVED***
	obj := r.toObject(call.This)
	if d, ok := obj.self.(*dateObject); ok ***REMOVED***
		if d.isSet() ***REMOVED***
			t, ok := _dateSetFullYear(d.time(), limitCallArgs(call, 2), -1, false)
			if !ok ***REMOVED***
				d.unset()
				return _NaN
			***REMOVED***
			return d.setTimeMs(timeToMsec(t))
		***REMOVED*** else ***REMOVED***
			return _NaN
		***REMOVED***
	***REMOVED***
	panic(r.NewTypeError("Method Date.prototype.setMonth is called on incompatible receiver"))
***REMOVED***

func (r *Runtime) dateproto_setUTCMonth(call FunctionCall) Value ***REMOVED***
	obj := r.toObject(call.This)
	if d, ok := obj.self.(*dateObject); ok ***REMOVED***
		if d.isSet() ***REMOVED***
			t, ok := _dateSetFullYear(d.timeUTC(), limitCallArgs(call, 2), -1, true)
			if !ok ***REMOVED***
				d.unset()
				return _NaN
			***REMOVED***
			return d.setTimeMs(timeToMsec(t))
		***REMOVED*** else ***REMOVED***
			return _NaN
		***REMOVED***
	***REMOVED***
	panic(r.NewTypeError("Method Date.prototype.setUTCMonth is called on incompatible receiver"))
***REMOVED***

func (r *Runtime) dateproto_setFullYear(call FunctionCall) Value ***REMOVED***
	obj := r.toObject(call.This)
	if d, ok := obj.self.(*dateObject); ok ***REMOVED***
		var t time.Time
		if d.isSet() ***REMOVED***
			t = d.time()
		***REMOVED*** else ***REMOVED***
			t = time.Date(1970, time.January, 1, 0, 0, 0, 0, time.Local)
		***REMOVED***
		t, ok := _dateSetFullYear(t, limitCallArgs(call, 3), 0, false)
		if !ok ***REMOVED***
			d.unset()
			return _NaN
		***REMOVED***
		return d.setTimeMs(timeToMsec(t))
	***REMOVED***
	panic(r.NewTypeError("Method Date.prototype.setFullYear is called on incompatible receiver"))
***REMOVED***

func (r *Runtime) dateproto_setUTCFullYear(call FunctionCall) Value ***REMOVED***
	obj := r.toObject(call.This)
	if d, ok := obj.self.(*dateObject); ok ***REMOVED***
		var t time.Time
		if d.isSet() ***REMOVED***
			t = d.timeUTC()
		***REMOVED*** else ***REMOVED***
			t = time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)
		***REMOVED***
		t, ok := _dateSetFullYear(t, limitCallArgs(call, 3), 0, true)
		if !ok ***REMOVED***
			d.unset()
			return _NaN
		***REMOVED***
		return d.setTimeMs(timeToMsec(t))
	***REMOVED***
	panic(r.NewTypeError("Method Date.prototype.setUTCFullYear is called on incompatible receiver"))
***REMOVED***

func (r *Runtime) createDateProto(val *Object) objectImpl ***REMOVED***
	o := &baseObject***REMOVED***
		class:      classObject,
		val:        val,
		extensible: true,
		prototype:  r.global.ObjectPrototype,
	***REMOVED***
	o.init()

	o._putProp("constructor", r.global.Date, true, false, true)
	o._putProp("toString", r.newNativeFunc(r.dateproto_toString, nil, "toString", nil, 0), true, false, true)
	o._putProp("toDateString", r.newNativeFunc(r.dateproto_toDateString, nil, "toDateString", nil, 0), true, false, true)
	o._putProp("toTimeString", r.newNativeFunc(r.dateproto_toTimeString, nil, "toTimeString", nil, 0), true, false, true)
	o._putProp("toLocaleString", r.newNativeFunc(r.dateproto_toLocaleString, nil, "toLocaleString", nil, 0), true, false, true)
	o._putProp("toLocaleDateString", r.newNativeFunc(r.dateproto_toLocaleDateString, nil, "toLocaleDateString", nil, 0), true, false, true)
	o._putProp("toLocaleTimeString", r.newNativeFunc(r.dateproto_toLocaleTimeString, nil, "toLocaleTimeString", nil, 0), true, false, true)
	o._putProp("valueOf", r.newNativeFunc(r.dateproto_valueOf, nil, "valueOf", nil, 0), true, false, true)
	o._putProp("getTime", r.newNativeFunc(r.dateproto_getTime, nil, "getTime", nil, 0), true, false, true)
	o._putProp("getFullYear", r.newNativeFunc(r.dateproto_getFullYear, nil, "getFullYear", nil, 0), true, false, true)
	o._putProp("getUTCFullYear", r.newNativeFunc(r.dateproto_getUTCFullYear, nil, "getUTCFullYear", nil, 0), true, false, true)
	o._putProp("getMonth", r.newNativeFunc(r.dateproto_getMonth, nil, "getMonth", nil, 0), true, false, true)
	o._putProp("getUTCMonth", r.newNativeFunc(r.dateproto_getUTCMonth, nil, "getUTCMonth", nil, 0), true, false, true)
	o._putProp("getDate", r.newNativeFunc(r.dateproto_getDate, nil, "getDate", nil, 0), true, false, true)
	o._putProp("getUTCDate", r.newNativeFunc(r.dateproto_getUTCDate, nil, "getUTCDate", nil, 0), true, false, true)
	o._putProp("getDay", r.newNativeFunc(r.dateproto_getDay, nil, "getDay", nil, 0), true, false, true)
	o._putProp("getUTCDay", r.newNativeFunc(r.dateproto_getUTCDay, nil, "getUTCDay", nil, 0), true, false, true)
	o._putProp("getHours", r.newNativeFunc(r.dateproto_getHours, nil, "getHours", nil, 0), true, false, true)
	o._putProp("getUTCHours", r.newNativeFunc(r.dateproto_getUTCHours, nil, "getUTCHours", nil, 0), true, false, true)
	o._putProp("getMinutes", r.newNativeFunc(r.dateproto_getMinutes, nil, "getMinutes", nil, 0), true, false, true)
	o._putProp("getUTCMinutes", r.newNativeFunc(r.dateproto_getUTCMinutes, nil, "getUTCMinutes", nil, 0), true, false, true)
	o._putProp("getSeconds", r.newNativeFunc(r.dateproto_getSeconds, nil, "getSeconds", nil, 0), true, false, true)
	o._putProp("getUTCSeconds", r.newNativeFunc(r.dateproto_getUTCSeconds, nil, "getUTCSeconds", nil, 0), true, false, true)
	o._putProp("getMilliseconds", r.newNativeFunc(r.dateproto_getMilliseconds, nil, "getMilliseconds", nil, 0), true, false, true)
	o._putProp("getUTCMilliseconds", r.newNativeFunc(r.dateproto_getUTCMilliseconds, nil, "getUTCMilliseconds", nil, 0), true, false, true)
	o._putProp("getTimezoneOffset", r.newNativeFunc(r.dateproto_getTimezoneOffset, nil, "getTimezoneOffset", nil, 0), true, false, true)
	o._putProp("setTime", r.newNativeFunc(r.dateproto_setTime, nil, "setTime", nil, 1), true, false, true)
	o._putProp("setMilliseconds", r.newNativeFunc(r.dateproto_setMilliseconds, nil, "setMilliseconds", nil, 1), true, false, true)
	o._putProp("setUTCMilliseconds", r.newNativeFunc(r.dateproto_setUTCMilliseconds, nil, "setUTCMilliseconds", nil, 1), true, false, true)
	o._putProp("setSeconds", r.newNativeFunc(r.dateproto_setSeconds, nil, "setSeconds", nil, 2), true, false, true)
	o._putProp("setUTCSeconds", r.newNativeFunc(r.dateproto_setUTCSeconds, nil, "setUTCSeconds", nil, 2), true, false, true)
	o._putProp("setMinutes", r.newNativeFunc(r.dateproto_setMinutes, nil, "setMinutes", nil, 3), true, false, true)
	o._putProp("setUTCMinutes", r.newNativeFunc(r.dateproto_setUTCMinutes, nil, "setUTCMinutes", nil, 3), true, false, true)
	o._putProp("setHours", r.newNativeFunc(r.dateproto_setHours, nil, "setHours", nil, 4), true, false, true)
	o._putProp("setUTCHours", r.newNativeFunc(r.dateproto_setUTCHours, nil, "setUTCHours", nil, 4), true, false, true)
	o._putProp("setDate", r.newNativeFunc(r.dateproto_setDate, nil, "setDate", nil, 1), true, false, true)
	o._putProp("setUTCDate", r.newNativeFunc(r.dateproto_setUTCDate, nil, "setUTCDate", nil, 1), true, false, true)
	o._putProp("setMonth", r.newNativeFunc(r.dateproto_setMonth, nil, "setMonth", nil, 2), true, false, true)
	o._putProp("setUTCMonth", r.newNativeFunc(r.dateproto_setUTCMonth, nil, "setUTCMonth", nil, 2), true, false, true)
	o._putProp("setFullYear", r.newNativeFunc(r.dateproto_setFullYear, nil, "setFullYear", nil, 3), true, false, true)
	o._putProp("setUTCFullYear", r.newNativeFunc(r.dateproto_setUTCFullYear, nil, "setUTCFullYear", nil, 3), true, false, true)
	o._putProp("toUTCString", r.newNativeFunc(r.dateproto_toUTCString, nil, "toUTCString", nil, 0), true, false, true)
	o._putProp("toISOString", r.newNativeFunc(r.dateproto_toISOString, nil, "toISOString", nil, 0), true, false, true)
	o._putProp("toJSON", r.newNativeFunc(r.dateproto_toJSON, nil, "toJSON", nil, 1), true, false, true)

	o._putSym(SymToPrimitive, valueProp(r.newNativeFunc(r.dateproto_toPrimitive, nil, "[Symbol.toPrimitive]", nil, 1), false, false, true))

	return o
***REMOVED***

func (r *Runtime) createDate(val *Object) objectImpl ***REMOVED***
	o := r.newNativeFuncObj(val, r.builtin_date, r.builtin_newDate, "Date", r.global.DatePrototype, 7)

	o._putProp("parse", r.newNativeFunc(r.date_parse, nil, "parse", nil, 1), true, false, true)
	o._putProp("UTC", r.newNativeFunc(r.date_UTC, nil, "UTC", nil, 7), true, false, true)
	o._putProp("now", r.newNativeFunc(r.date_now, nil, "now", nil, 0), true, false, true)

	return o
***REMOVED***

func (r *Runtime) initDate() ***REMOVED***
	r.global.DatePrototype = r.newLazyObject(r.createDateProto)

	r.global.Date = r.newLazyObject(r.createDate)
	r.addToGlobal("Date", r.global.Date)
***REMOVED***
