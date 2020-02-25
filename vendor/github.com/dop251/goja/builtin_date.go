package goja

import (
	"fmt"
	"math"
	"time"
)

const (
	maxTime = 8.64e15
)

func timeFromMsec(msec int64) time.Time ***REMOVED***
	sec := msec / 1000
	nsec := (msec % 1000) * 1e6
	return time.Unix(sec, nsec)
***REMOVED***

func timeToMsec(t time.Time) int64 ***REMOVED***
	return t.Unix()*1000 + int64(t.Nanosecond())/1e6
***REMOVED***

func (r *Runtime) makeDate(args []Value, loc *time.Location) (t time.Time, valid bool) ***REMOVED***
	pick := func(index int, default_ int64) (int64, bool) ***REMOVED***
		if index >= len(args) ***REMOVED***
			return default_, true
		***REMOVED***
		value := args[index]
		if valueInt, ok := value.assertInt(); ok ***REMOVED***
			return valueInt, true
		***REMOVED***
		valueFloat := value.ToFloat()
		if math.IsNaN(valueFloat) || math.IsInf(valueFloat, 0) ***REMOVED***
			return 0, false
		***REMOVED***
		return int64(valueFloat), true
	***REMOVED***

	switch ***REMOVED***
	case len(args) >= 2:
		var year, month, day, hour, minute, second, millisecond int64
		if year, valid = pick(0, 1900); !valid ***REMOVED***
			return
		***REMOVED***
		if month, valid = pick(1, 0); !valid ***REMOVED***
			return
		***REMOVED***
		if day, valid = pick(2, 1); !valid ***REMOVED***
			return
		***REMOVED***
		if hour, valid = pick(3, 0); !valid ***REMOVED***
			return
		***REMOVED***
		if minute, valid = pick(4, 0); !valid ***REMOVED***
			return
		***REMOVED***
		if second, valid = pick(5, 0); !valid ***REMOVED***
			return
		***REMOVED***
		if millisecond, valid = pick(6, 0); !valid ***REMOVED***
			return
		***REMOVED***

		if year >= 0 && year <= 99 ***REMOVED***
			year += 1900
		***REMOVED***

		t = time.Date(int(year), time.Month(int(month)+1), int(day), int(hour), int(minute), int(second), int(millisecond)*1e6, loc)
	case len(args) == 0:
		t = r.now()
		valid = true
	default: // one argument
		pv := toPrimitiveNumber(args[0])
		if val, ok := pv.assertString(); ok ***REMOVED***
			return dateParse(val.String())
		***REMOVED***

		var n int64
		if i, ok := pv.assertInt(); ok ***REMOVED***
			n = i
		***REMOVED*** else if f, ok := pv.assertFloat(); ok ***REMOVED***
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
	msec := t.Unix()*1000 + int64(t.Nanosecond()/1e6)
	if msec < 0 ***REMOVED***
		msec = -msec
	***REMOVED***
	if msec > maxTime ***REMOVED***
		valid = false
	***REMOVED***
	return
***REMOVED***

func (r *Runtime) newDateTime(args []Value, loc *time.Location) *Object ***REMOVED***
	t, isSet := r.makeDate(args, loc)
	return r.newDateObject(t, isSet)
***REMOVED***

func (r *Runtime) builtin_newDate(args []Value) *Object ***REMOVED***
	return r.newDateTime(args, time.Local)
***REMOVED***

func (r *Runtime) builtin_date(call FunctionCall) Value ***REMOVED***
	return asciiString(dateFormat(r.now()))
***REMOVED***

func (r *Runtime) date_parse(call FunctionCall) Value ***REMOVED***
	t, set := dateParse(call.Argument(0).String())
	if set ***REMOVED***
		return intToValue(timeToMsec(t))
	***REMOVED***
	return _NaN
***REMOVED***

func (r *Runtime) date_UTC(call FunctionCall) Value ***REMOVED***
	t, valid := r.makeDate(call.Arguments, time.UTC)
	if !valid ***REMOVED***
		return _NaN
	***REMOVED***
	return intToValue(timeToMsec(t))
***REMOVED***

func (r *Runtime) date_now(call FunctionCall) Value ***REMOVED***
	return intToValue(timeToMsec(r.now()))
***REMOVED***

func (r *Runtime) dateproto_toString(call FunctionCall) Value ***REMOVED***
	obj := r.toObject(call.This)
	if d, ok := obj.self.(*dateObject); ok ***REMOVED***
		if d.isSet ***REMOVED***
			return asciiString(d.time.Format(dateTimeLayout))
		***REMOVED*** else ***REMOVED***
			return stringInvalidDate
		***REMOVED***
	***REMOVED***
	r.typeErrorResult(true, "Method Date.prototype.toString is called on incompatible receiver")
	panic("Unreachable")
***REMOVED***

func (r *Runtime) dateproto_toUTCString(call FunctionCall) Value ***REMOVED***
	obj := r.toObject(call.This)
	if d, ok := obj.self.(*dateObject); ok ***REMOVED***
		if d.isSet ***REMOVED***
			return asciiString(d.time.In(time.UTC).Format(utcDateTimeLayout))
		***REMOVED*** else ***REMOVED***
			return stringInvalidDate
		***REMOVED***
	***REMOVED***
	r.typeErrorResult(true, "Method Date.prototype.toUTCString is called on incompatible receiver")
	panic("Unreachable")
***REMOVED***

func (r *Runtime) dateproto_toISOString(call FunctionCall) Value ***REMOVED***
	obj := r.toObject(call.This)
	if d, ok := obj.self.(*dateObject); ok ***REMOVED***
		if d.isSet ***REMOVED***
			utc := d.time.In(time.UTC)
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
	r.typeErrorResult(true, "Method Date.prototype.toISOString is called on incompatible receiver")
	panic("Unreachable")
***REMOVED***

func (r *Runtime) dateproto_toJSON(call FunctionCall) Value ***REMOVED***
	obj := r.toObject(call.This)
	tv := obj.self.toPrimitiveNumber()
	if f, ok := tv.assertFloat(); ok ***REMOVED***
		if math.IsNaN(f) || math.IsInf(f, 0) ***REMOVED***
			return _null
		***REMOVED***
	***REMOVED*** else if _, ok := tv.assertInt(); !ok ***REMOVED***
		return _null
	***REMOVED***

	if toISO, ok := obj.self.getStr("toISOString").(*Object); ok ***REMOVED***
		if toISO, ok := toISO.self.assertCallable(); ok ***REMOVED***
			return toISO(FunctionCall***REMOVED***
				This: obj,
			***REMOVED***)
		***REMOVED***
	***REMOVED***

	r.typeErrorResult(true, "toISOString is not a function")
	panic("Unreachable")
***REMOVED***

func (r *Runtime) dateproto_toDateString(call FunctionCall) Value ***REMOVED***
	obj := r.toObject(call.This)
	if d, ok := obj.self.(*dateObject); ok ***REMOVED***
		if d.isSet ***REMOVED***
			return asciiString(d.time.Format(dateLayout))
		***REMOVED*** else ***REMOVED***
			return stringInvalidDate
		***REMOVED***
	***REMOVED***
	r.typeErrorResult(true, "Method Date.prototype.toDateString is called on incompatible receiver")
	panic("Unreachable")
***REMOVED***

func (r *Runtime) dateproto_toTimeString(call FunctionCall) Value ***REMOVED***
	obj := r.toObject(call.This)
	if d, ok := obj.self.(*dateObject); ok ***REMOVED***
		if d.isSet ***REMOVED***
			return asciiString(d.time.Format(timeLayout))
		***REMOVED*** else ***REMOVED***
			return stringInvalidDate
		***REMOVED***
	***REMOVED***
	r.typeErrorResult(true, "Method Date.prototype.toTimeString is called on incompatible receiver")
	panic("Unreachable")
***REMOVED***

func (r *Runtime) dateproto_toLocaleString(call FunctionCall) Value ***REMOVED***
	obj := r.toObject(call.This)
	if d, ok := obj.self.(*dateObject); ok ***REMOVED***
		if d.isSet ***REMOVED***
			return asciiString(d.time.Format(datetimeLayout_en_GB))
		***REMOVED*** else ***REMOVED***
			return stringInvalidDate
		***REMOVED***
	***REMOVED***
	r.typeErrorResult(true, "Method Date.prototype.toLocaleString is called on incompatible receiver")
	panic("Unreachable")
***REMOVED***

func (r *Runtime) dateproto_toLocaleDateString(call FunctionCall) Value ***REMOVED***
	obj := r.toObject(call.This)
	if d, ok := obj.self.(*dateObject); ok ***REMOVED***
		if d.isSet ***REMOVED***
			return asciiString(d.time.Format(dateLayout_en_GB))
		***REMOVED*** else ***REMOVED***
			return stringInvalidDate
		***REMOVED***
	***REMOVED***
	r.typeErrorResult(true, "Method Date.prototype.toLocaleDateString is called on incompatible receiver")
	panic("Unreachable")
***REMOVED***

func (r *Runtime) dateproto_toLocaleTimeString(call FunctionCall) Value ***REMOVED***
	obj := r.toObject(call.This)
	if d, ok := obj.self.(*dateObject); ok ***REMOVED***
		if d.isSet ***REMOVED***
			return asciiString(d.time.Format(timeLayout_en_GB))
		***REMOVED*** else ***REMOVED***
			return stringInvalidDate
		***REMOVED***
	***REMOVED***
	r.typeErrorResult(true, "Method Date.prototype.toLocaleTimeString is called on incompatible receiver")
	panic("Unreachable")
***REMOVED***

func (r *Runtime) dateproto_valueOf(call FunctionCall) Value ***REMOVED***
	obj := r.toObject(call.This)
	if d, ok := obj.self.(*dateObject); ok ***REMOVED***
		if d.isSet ***REMOVED***
			return intToValue(d.time.Unix()*1000 + int64(d.time.Nanosecond()/1e6))
		***REMOVED*** else ***REMOVED***
			return _NaN
		***REMOVED***
	***REMOVED***
	r.typeErrorResult(true, "Method Date.prototype.valueOf is called on incompatible receiver")
	return nil
***REMOVED***

func (r *Runtime) dateproto_getTime(call FunctionCall) Value ***REMOVED***
	obj := r.toObject(call.This)
	if d, ok := obj.self.(*dateObject); ok ***REMOVED***
		if d.isSet ***REMOVED***
			return intToValue(timeToMsec(d.time))
		***REMOVED*** else ***REMOVED***
			return _NaN
		***REMOVED***
	***REMOVED***
	r.typeErrorResult(true, "Method Date.prototype.getTime is called on incompatible receiver")
	return nil
***REMOVED***

func (r *Runtime) dateproto_getFullYear(call FunctionCall) Value ***REMOVED***
	obj := r.toObject(call.This)
	if d, ok := obj.self.(*dateObject); ok ***REMOVED***
		if d.isSet ***REMOVED***
			return intToValue(int64(d.time.Year()))
		***REMOVED*** else ***REMOVED***
			return _NaN
		***REMOVED***
	***REMOVED***
	r.typeErrorResult(true, "Method Date.prototype.getFullYear is called on incompatible receiver")
	return nil
***REMOVED***

func (r *Runtime) dateproto_getUTCFullYear(call FunctionCall) Value ***REMOVED***
	obj := r.toObject(call.This)
	if d, ok := obj.self.(*dateObject); ok ***REMOVED***
		if d.isSet ***REMOVED***
			return intToValue(int64(d.time.In(time.UTC).Year()))
		***REMOVED*** else ***REMOVED***
			return _NaN
		***REMOVED***
	***REMOVED***
	r.typeErrorResult(true, "Method Date.prototype.getUTCFullYear is called on incompatible receiver")
	return nil
***REMOVED***

func (r *Runtime) dateproto_getMonth(call FunctionCall) Value ***REMOVED***
	obj := r.toObject(call.This)
	if d, ok := obj.self.(*dateObject); ok ***REMOVED***
		if d.isSet ***REMOVED***
			return intToValue(int64(d.time.Month()) - 1)
		***REMOVED*** else ***REMOVED***
			return _NaN
		***REMOVED***
	***REMOVED***
	r.typeErrorResult(true, "Method Date.prototype.getMonth is called on incompatible receiver")
	return nil
***REMOVED***

func (r *Runtime) dateproto_getUTCMonth(call FunctionCall) Value ***REMOVED***
	obj := r.toObject(call.This)
	if d, ok := obj.self.(*dateObject); ok ***REMOVED***
		if d.isSet ***REMOVED***
			return intToValue(int64(d.time.In(time.UTC).Month()) - 1)
		***REMOVED*** else ***REMOVED***
			return _NaN
		***REMOVED***
	***REMOVED***
	r.typeErrorResult(true, "Method Date.prototype.getUTCMonth is called on incompatible receiver")
	return nil
***REMOVED***

func (r *Runtime) dateproto_getHours(call FunctionCall) Value ***REMOVED***
	obj := r.toObject(call.This)
	if d, ok := obj.self.(*dateObject); ok ***REMOVED***
		if d.isSet ***REMOVED***
			return intToValue(int64(d.time.Hour()))
		***REMOVED*** else ***REMOVED***
			return _NaN
		***REMOVED***
	***REMOVED***
	r.typeErrorResult(true, "Method Date.prototype.getHours is called on incompatible receiver")
	return nil
***REMOVED***

func (r *Runtime) dateproto_getUTCHours(call FunctionCall) Value ***REMOVED***
	obj := r.toObject(call.This)
	if d, ok := obj.self.(*dateObject); ok ***REMOVED***
		if d.isSet ***REMOVED***
			return intToValue(int64(d.time.In(time.UTC).Hour()))
		***REMOVED*** else ***REMOVED***
			return _NaN
		***REMOVED***
	***REMOVED***
	r.typeErrorResult(true, "Method Date.prototype.getUTCHours is called on incompatible receiver")
	return nil
***REMOVED***

func (r *Runtime) dateproto_getDate(call FunctionCall) Value ***REMOVED***
	obj := r.toObject(call.This)
	if d, ok := obj.self.(*dateObject); ok ***REMOVED***
		if d.isSet ***REMOVED***
			return intToValue(int64(d.time.Day()))
		***REMOVED*** else ***REMOVED***
			return _NaN
		***REMOVED***
	***REMOVED***
	r.typeErrorResult(true, "Method Date.prototype.getDate is called on incompatible receiver")
	return nil
***REMOVED***

func (r *Runtime) dateproto_getUTCDate(call FunctionCall) Value ***REMOVED***
	obj := r.toObject(call.This)
	if d, ok := obj.self.(*dateObject); ok ***REMOVED***
		if d.isSet ***REMOVED***
			return intToValue(int64(d.time.In(time.UTC).Day()))
		***REMOVED*** else ***REMOVED***
			return _NaN
		***REMOVED***
	***REMOVED***
	r.typeErrorResult(true, "Method Date.prototype.getUTCDate is called on incompatible receiver")
	return nil
***REMOVED***

func (r *Runtime) dateproto_getDay(call FunctionCall) Value ***REMOVED***
	obj := r.toObject(call.This)
	if d, ok := obj.self.(*dateObject); ok ***REMOVED***
		if d.isSet ***REMOVED***
			return intToValue(int64(d.time.Weekday()))
		***REMOVED*** else ***REMOVED***
			return _NaN
		***REMOVED***
	***REMOVED***
	r.typeErrorResult(true, "Method Date.prototype.getDay is called on incompatible receiver")
	return nil
***REMOVED***

func (r *Runtime) dateproto_getUTCDay(call FunctionCall) Value ***REMOVED***
	obj := r.toObject(call.This)
	if d, ok := obj.self.(*dateObject); ok ***REMOVED***
		if d.isSet ***REMOVED***
			return intToValue(int64(d.time.In(time.UTC).Weekday()))
		***REMOVED*** else ***REMOVED***
			return _NaN
		***REMOVED***
	***REMOVED***
	r.typeErrorResult(true, "Method Date.prototype.getUTCDay is called on incompatible receiver")
	return nil
***REMOVED***

func (r *Runtime) dateproto_getMinutes(call FunctionCall) Value ***REMOVED***
	obj := r.toObject(call.This)
	if d, ok := obj.self.(*dateObject); ok ***REMOVED***
		if d.isSet ***REMOVED***
			return intToValue(int64(d.time.Minute()))
		***REMOVED*** else ***REMOVED***
			return _NaN
		***REMOVED***
	***REMOVED***
	r.typeErrorResult(true, "Method Date.prototype.getMinutes is called on incompatible receiver")
	return nil
***REMOVED***

func (r *Runtime) dateproto_getUTCMinutes(call FunctionCall) Value ***REMOVED***
	obj := r.toObject(call.This)
	if d, ok := obj.self.(*dateObject); ok ***REMOVED***
		if d.isSet ***REMOVED***
			return intToValue(int64(d.time.In(time.UTC).Minute()))
		***REMOVED*** else ***REMOVED***
			return _NaN
		***REMOVED***
	***REMOVED***
	r.typeErrorResult(true, "Method Date.prototype.getUTCMinutes is called on incompatible receiver")
	return nil
***REMOVED***

func (r *Runtime) dateproto_getSeconds(call FunctionCall) Value ***REMOVED***
	obj := r.toObject(call.This)
	if d, ok := obj.self.(*dateObject); ok ***REMOVED***
		if d.isSet ***REMOVED***
			return intToValue(int64(d.time.Second()))
		***REMOVED*** else ***REMOVED***
			return _NaN
		***REMOVED***
	***REMOVED***
	r.typeErrorResult(true, "Method Date.prototype.getSeconds is called on incompatible receiver")
	return nil
***REMOVED***

func (r *Runtime) dateproto_getUTCSeconds(call FunctionCall) Value ***REMOVED***
	obj := r.toObject(call.This)
	if d, ok := obj.self.(*dateObject); ok ***REMOVED***
		if d.isSet ***REMOVED***
			return intToValue(int64(d.time.In(time.UTC).Second()))
		***REMOVED*** else ***REMOVED***
			return _NaN
		***REMOVED***
	***REMOVED***
	r.typeErrorResult(true, "Method Date.prototype.getUTCSeconds is called on incompatible receiver")
	return nil
***REMOVED***

func (r *Runtime) dateproto_getMilliseconds(call FunctionCall) Value ***REMOVED***
	obj := r.toObject(call.This)
	if d, ok := obj.self.(*dateObject); ok ***REMOVED***
		if d.isSet ***REMOVED***
			return intToValue(int64(d.time.Nanosecond() / 1e6))
		***REMOVED*** else ***REMOVED***
			return _NaN
		***REMOVED***
	***REMOVED***
	r.typeErrorResult(true, "Method Date.prototype.getMilliseconds is called on incompatible receiver")
	return nil
***REMOVED***

func (r *Runtime) dateproto_getUTCMilliseconds(call FunctionCall) Value ***REMOVED***
	obj := r.toObject(call.This)
	if d, ok := obj.self.(*dateObject); ok ***REMOVED***
		if d.isSet ***REMOVED***
			return intToValue(int64(d.time.In(time.UTC).Nanosecond() / 1e6))
		***REMOVED*** else ***REMOVED***
			return _NaN
		***REMOVED***
	***REMOVED***
	r.typeErrorResult(true, "Method Date.prototype.getUTCMilliseconds is called on incompatible receiver")
	return nil
***REMOVED***

func (r *Runtime) dateproto_getTimezoneOffset(call FunctionCall) Value ***REMOVED***
	obj := r.toObject(call.This)
	if d, ok := obj.self.(*dateObject); ok ***REMOVED***
		if d.isSet ***REMOVED***
			_, offset := d.time.Zone()
			return floatToValue(float64(-offset) / 60)
		***REMOVED*** else ***REMOVED***
			return _NaN
		***REMOVED***
	***REMOVED***
	r.typeErrorResult(true, "Method Date.prototype.getTimezoneOffset is called on incompatible receiver")
	return nil
***REMOVED***

func (r *Runtime) dateproto_setTime(call FunctionCall) Value ***REMOVED***
	obj := r.toObject(call.This)
	if d, ok := obj.self.(*dateObject); ok ***REMOVED***
		msec := call.Argument(0).ToInteger()
		d.time = timeFromMsec(msec)
		return intToValue(msec)
	***REMOVED***
	r.typeErrorResult(true, "Method Date.prototype.setTime is called on incompatible receiver")
	panic("Unreachable")
***REMOVED***

func (r *Runtime) dateproto_setMilliseconds(call FunctionCall) Value ***REMOVED***
	obj := r.toObject(call.This)
	if d, ok := obj.self.(*dateObject); ok ***REMOVED***
		if d.isSet ***REMOVED***
			msec := call.Argument(0).ToInteger()
			m := timeToMsec(d.time) - int64(d.time.Nanosecond())/1e6 + msec
			d.time = timeFromMsec(m)
			return intToValue(m)
		***REMOVED*** else ***REMOVED***
			return _NaN
		***REMOVED***
	***REMOVED***
	r.typeErrorResult(true, "Method Date.prototype.setMilliseconds is called on incompatible receiver")
	panic("Unreachable")
***REMOVED***

func (r *Runtime) dateproto_setUTCMilliseconds(call FunctionCall) Value ***REMOVED***
	obj := r.toObject(call.This)
	if d, ok := obj.self.(*dateObject); ok ***REMOVED***
		if d.isSet ***REMOVED***
			msec := call.Argument(0).ToInteger()
			m := timeToMsec(d.time) - int64(d.time.Nanosecond())/1e6 + msec
			d.time = timeFromMsec(m)
			return intToValue(m)
		***REMOVED*** else ***REMOVED***
			return _NaN
		***REMOVED***
	***REMOVED***
	r.typeErrorResult(true, "Method Date.prototype.setUTCMilliseconds is called on incompatible receiver")
	panic("Unreachable")
***REMOVED***

func (r *Runtime) dateproto_setSeconds(call FunctionCall) Value ***REMOVED***
	obj := r.toObject(call.This)
	if d, ok := obj.self.(*dateObject); ok ***REMOVED***
		if d.isSet ***REMOVED***
			sec := int(call.Argument(0).ToInteger())
			var nsec int
			if len(call.Arguments) > 1 ***REMOVED***
				nsec = int(call.Arguments[1].ToInteger() * 1e6)
			***REMOVED*** else ***REMOVED***
				nsec = d.time.Nanosecond()
			***REMOVED***
			d.time = time.Date(d.time.Year(), d.time.Month(), d.time.Day(), d.time.Hour(), d.time.Minute(), sec, nsec, time.Local)
			return intToValue(timeToMsec(d.time))
		***REMOVED*** else ***REMOVED***
			return _NaN
		***REMOVED***
	***REMOVED***
	r.typeErrorResult(true, "Method Date.prototype.setSeconds is called on incompatible receiver")
	panic("Unreachable")
***REMOVED***

func (r *Runtime) dateproto_setUTCSeconds(call FunctionCall) Value ***REMOVED***
	obj := r.toObject(call.This)
	if d, ok := obj.self.(*dateObject); ok ***REMOVED***
		if d.isSet ***REMOVED***
			sec := int(call.Argument(0).ToInteger())
			var nsec int
			t := d.time.In(time.UTC)
			if len(call.Arguments) > 1 ***REMOVED***
				nsec = int(call.Arguments[1].ToInteger() * 1e6)
			***REMOVED*** else ***REMOVED***
				nsec = t.Nanosecond()
			***REMOVED***
			d.time = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), sec, nsec, time.UTC).In(time.Local)
			return intToValue(timeToMsec(d.time))
		***REMOVED*** else ***REMOVED***
			return _NaN
		***REMOVED***
	***REMOVED***
	r.typeErrorResult(true, "Method Date.prototype.setUTCSeconds is called on incompatible receiver")
	panic("Unreachable")
***REMOVED***

func (r *Runtime) dateproto_setMinutes(call FunctionCall) Value ***REMOVED***
	obj := r.toObject(call.This)
	if d, ok := obj.self.(*dateObject); ok ***REMOVED***
		if d.isSet ***REMOVED***
			min := int(call.Argument(0).ToInteger())
			var sec, nsec int
			if len(call.Arguments) > 1 ***REMOVED***
				sec = int(call.Arguments[1].ToInteger())
			***REMOVED*** else ***REMOVED***
				sec = d.time.Second()
			***REMOVED***
			if len(call.Arguments) > 2 ***REMOVED***
				nsec = int(call.Arguments[2].ToInteger() * 1e6)
			***REMOVED*** else ***REMOVED***
				nsec = d.time.Nanosecond()
			***REMOVED***
			d.time = time.Date(d.time.Year(), d.time.Month(), d.time.Day(), d.time.Hour(), min, sec, nsec, time.Local)
			return intToValue(timeToMsec(d.time))
		***REMOVED*** else ***REMOVED***
			return _NaN
		***REMOVED***
	***REMOVED***
	r.typeErrorResult(true, "Method Date.prototype.setMinutes is called on incompatible receiver")
	panic("Unreachable")
***REMOVED***

func (r *Runtime) dateproto_setUTCMinutes(call FunctionCall) Value ***REMOVED***
	obj := r.toObject(call.This)
	if d, ok := obj.self.(*dateObject); ok ***REMOVED***
		if d.isSet ***REMOVED***
			min := int(call.Argument(0).ToInteger())
			var sec, nsec int
			t := d.time.In(time.UTC)
			if len(call.Arguments) > 1 ***REMOVED***
				sec = int(call.Arguments[1].ToInteger())
			***REMOVED*** else ***REMOVED***
				sec = t.Second()
			***REMOVED***
			if len(call.Arguments) > 2 ***REMOVED***
				nsec = int(call.Arguments[2].ToInteger() * 1e6)
			***REMOVED*** else ***REMOVED***
				nsec = t.Nanosecond()
			***REMOVED***
			d.time = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), min, sec, nsec, time.UTC).In(time.Local)
			return intToValue(timeToMsec(d.time))
		***REMOVED*** else ***REMOVED***
			return _NaN
		***REMOVED***
	***REMOVED***
	r.typeErrorResult(true, "Method Date.prototype.setUTCMinutes is called on incompatible receiver")
	panic("Unreachable")
***REMOVED***

func (r *Runtime) dateproto_setHours(call FunctionCall) Value ***REMOVED***
	obj := r.toObject(call.This)
	if d, ok := obj.self.(*dateObject); ok ***REMOVED***
		if d.isSet ***REMOVED***
			hour := int(call.Argument(0).ToInteger())
			var min, sec, nsec int
			if len(call.Arguments) > 1 ***REMOVED***
				min = int(call.Arguments[1].ToInteger())
			***REMOVED*** else ***REMOVED***
				min = d.time.Minute()
			***REMOVED***
			if len(call.Arguments) > 2 ***REMOVED***
				sec = int(call.Arguments[2].ToInteger())
			***REMOVED*** else ***REMOVED***
				sec = d.time.Second()
			***REMOVED***
			if len(call.Arguments) > 3 ***REMOVED***
				nsec = int(call.Arguments[3].ToInteger() * 1e6)
			***REMOVED*** else ***REMOVED***
				nsec = d.time.Nanosecond()
			***REMOVED***
			d.time = time.Date(d.time.Year(), d.time.Month(), d.time.Day(), hour, min, sec, nsec, time.Local)
			return intToValue(timeToMsec(d.time))
		***REMOVED*** else ***REMOVED***
			return _NaN
		***REMOVED***
	***REMOVED***
	r.typeErrorResult(true, "Method Date.prototype.setHours is called on incompatible receiver")
	panic("Unreachable")
***REMOVED***

func (r *Runtime) dateproto_setUTCHours(call FunctionCall) Value ***REMOVED***
	obj := r.toObject(call.This)
	if d, ok := obj.self.(*dateObject); ok ***REMOVED***
		if d.isSet ***REMOVED***
			hour := int(call.Argument(0).ToInteger())
			var min, sec, nsec int
			t := d.time.In(time.UTC)
			if len(call.Arguments) > 1 ***REMOVED***
				min = int(call.Arguments[1].ToInteger())
			***REMOVED*** else ***REMOVED***
				min = t.Minute()
			***REMOVED***
			if len(call.Arguments) > 2 ***REMOVED***
				sec = int(call.Arguments[2].ToInteger())
			***REMOVED*** else ***REMOVED***
				sec = t.Second()
			***REMOVED***
			if len(call.Arguments) > 3 ***REMOVED***
				nsec = int(call.Arguments[3].ToInteger() * 1e6)
			***REMOVED*** else ***REMOVED***
				nsec = t.Nanosecond()
			***REMOVED***
			d.time = time.Date(d.time.Year(), d.time.Month(), d.time.Day(), hour, min, sec, nsec, time.UTC).In(time.Local)
			return intToValue(timeToMsec(d.time))
		***REMOVED*** else ***REMOVED***
			return _NaN
		***REMOVED***
	***REMOVED***
	r.typeErrorResult(true, "Method Date.prototype.setUTCHours is called on incompatible receiver")
	panic("Unreachable")
***REMOVED***

func (r *Runtime) dateproto_setDate(call FunctionCall) Value ***REMOVED***
	obj := r.toObject(call.This)
	if d, ok := obj.self.(*dateObject); ok ***REMOVED***
		if d.isSet ***REMOVED***
			d.time = time.Date(d.time.Year(), d.time.Month(), int(call.Argument(0).ToInteger()), d.time.Hour(), d.time.Minute(), d.time.Second(), d.time.Nanosecond(), time.Local)
			return intToValue(timeToMsec(d.time))
		***REMOVED*** else ***REMOVED***
			return _NaN
		***REMOVED***
	***REMOVED***
	r.typeErrorResult(true, "Method Date.prototype.setDate is called on incompatible receiver")
	panic("Unreachable")
***REMOVED***

func (r *Runtime) dateproto_setUTCDate(call FunctionCall) Value ***REMOVED***
	obj := r.toObject(call.This)
	if d, ok := obj.self.(*dateObject); ok ***REMOVED***
		if d.isSet ***REMOVED***
			t := d.time.In(time.UTC)
			d.time = time.Date(t.Year(), t.Month(), int(call.Argument(0).ToInteger()), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), time.UTC).In(time.Local)
			return intToValue(timeToMsec(d.time))
		***REMOVED*** else ***REMOVED***
			return _NaN
		***REMOVED***
	***REMOVED***
	r.typeErrorResult(true, "Method Date.prototype.setUTCDate is called on incompatible receiver")
	panic("Unreachable")
***REMOVED***

func (r *Runtime) dateproto_setMonth(call FunctionCall) Value ***REMOVED***
	obj := r.toObject(call.This)
	if d, ok := obj.self.(*dateObject); ok ***REMOVED***
		if d.isSet ***REMOVED***
			month := time.Month(int(call.Argument(0).ToInteger()) + 1)
			var day int
			if len(call.Arguments) > 1 ***REMOVED***
				day = int(call.Arguments[1].ToInteger())
			***REMOVED*** else ***REMOVED***
				day = d.time.Day()
			***REMOVED***
			d.time = time.Date(d.time.Year(), month, day, d.time.Hour(), d.time.Minute(), d.time.Second(), d.time.Nanosecond(), time.Local)
			return intToValue(timeToMsec(d.time))
		***REMOVED*** else ***REMOVED***
			return _NaN
		***REMOVED***
	***REMOVED***
	r.typeErrorResult(true, "Method Date.prototype.setMonth is called on incompatible receiver")
	panic("Unreachable")
***REMOVED***

func (r *Runtime) dateproto_setUTCMonth(call FunctionCall) Value ***REMOVED***
	obj := r.toObject(call.This)
	if d, ok := obj.self.(*dateObject); ok ***REMOVED***
		if d.isSet ***REMOVED***
			month := time.Month(int(call.Argument(0).ToInteger()) + 1)
			var day int
			t := d.time.In(time.UTC)
			if len(call.Arguments) > 1 ***REMOVED***
				day = int(call.Arguments[1].ToInteger())
			***REMOVED*** else ***REMOVED***
				day = t.Day()
			***REMOVED***
			d.time = time.Date(t.Year(), month, day, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), time.UTC).In(time.Local)
			return intToValue(timeToMsec(d.time))
		***REMOVED*** else ***REMOVED***
			return _NaN
		***REMOVED***
	***REMOVED***
	r.typeErrorResult(true, "Method Date.prototype.setUTCMonth is called on incompatible receiver")
	panic("Unreachable")
***REMOVED***

func (r *Runtime) dateproto_setFullYear(call FunctionCall) Value ***REMOVED***
	obj := r.toObject(call.This)
	if d, ok := obj.self.(*dateObject); ok ***REMOVED***
		if !d.isSet ***REMOVED***
			d.time = time.Unix(0, 0)
		***REMOVED***
		year := int(call.Argument(0).ToInteger())
		var month time.Month
		var day int
		if len(call.Arguments) > 1 ***REMOVED***
			month = time.Month(call.Arguments[1].ToInteger() + 1)
		***REMOVED*** else ***REMOVED***
			month = d.time.Month()
		***REMOVED***
		if len(call.Arguments) > 2 ***REMOVED***
			day = int(call.Arguments[2].ToInteger())
		***REMOVED*** else ***REMOVED***
			day = d.time.Day()
		***REMOVED***
		d.time = time.Date(year, month, day, d.time.Hour(), d.time.Minute(), d.time.Second(), d.time.Nanosecond(), time.Local)
		return intToValue(timeToMsec(d.time))
	***REMOVED***
	r.typeErrorResult(true, "Method Date.prototype.setFullYear is called on incompatible receiver")
	panic("Unreachable")
***REMOVED***

func (r *Runtime) dateproto_setUTCFullYear(call FunctionCall) Value ***REMOVED***
	obj := r.toObject(call.This)
	if d, ok := obj.self.(*dateObject); ok ***REMOVED***
		if !d.isSet ***REMOVED***
			d.time = time.Unix(0, 0)
		***REMOVED***
		year := int(call.Argument(0).ToInteger())
		var month time.Month
		var day int
		t := d.time.In(time.UTC)
		if len(call.Arguments) > 1 ***REMOVED***
			month = time.Month(call.Arguments[1].ToInteger() + 1)
		***REMOVED*** else ***REMOVED***
			month = t.Month()
		***REMOVED***
		if len(call.Arguments) > 2 ***REMOVED***
			day = int(call.Arguments[2].ToInteger())
		***REMOVED*** else ***REMOVED***
			day = t.Day()
		***REMOVED***
		d.time = time.Date(year, month, day, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), time.UTC).In(time.Local)
		return intToValue(timeToMsec(d.time))
	***REMOVED***
	r.typeErrorResult(true, "Method Date.prototype.setUTCFullYear is called on incompatible receiver")
	panic("Unreachable")
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

	return o
***REMOVED***

func (r *Runtime) createDate(val *Object) objectImpl ***REMOVED***
	o := r.newNativeFuncObj(val, r.builtin_date, r.builtin_newDate, "Date", r.global.DatePrototype, 7)

	o._putProp("parse", r.newNativeFunc(r.date_parse, nil, "parse", nil, 1), true, false, true)
	o._putProp("UTC", r.newNativeFunc(r.date_UTC, nil, "UTC", nil, 7), true, false, true)
	o._putProp("now", r.newNativeFunc(r.date_now, nil, "now", nil, 0), true, false, true)

	return o
***REMOVED***

func (r *Runtime) newLazyObject(create func(*Object) objectImpl) *Object ***REMOVED***
	val := &Object***REMOVED***runtime: r***REMOVED***
	o := &lazyObject***REMOVED***
		val:    val,
		create: create,
	***REMOVED***
	val.self = o
	return val
***REMOVED***

func (r *Runtime) initDate() ***REMOVED***
	//r.global.DatePrototype = r.newObject()
	//o := r.global.DatePrototype.self
	r.global.DatePrototype = r.newLazyObject(r.createDateProto)

	//r.global.Date = r.newNativeFunc(r.builtin_date, r.builtin_newDate, "Date", r.global.DatePrototype, 7)
	//o := r.global.Date.self
	r.global.Date = r.newLazyObject(r.createDate)

	r.addToGlobal("Date", r.global.Date)
***REMOVED***
