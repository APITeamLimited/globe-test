package http

import (
	"errors"
	"fmt"

	"github.com/APITeamLimited/k6-worker/js/common"
	"github.com/dop251/goja"
)

//nolint:gochecknoglobals
var defaultExpectedStatuses = expectedStatuses***REMOVED***
	minmax: [][2]int***REMOVED******REMOVED***200, 399***REMOVED******REMOVED***,
***REMOVED***

// expectedStatuses is specifically totally unexported so it can't be used for anything else but
// SetResponseCallback and nothing can be done from the js side to modify it or make an instance of
// it except using ExpectedStatuses
type expectedStatuses struct ***REMOVED***
	minmax [][2]int
	exact  []int
***REMOVED***

func (e expectedStatuses) match(status int) bool ***REMOVED***
	for _, v := range e.exact ***REMOVED***
		if v == status ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***

	for _, v := range e.minmax ***REMOVED***
		if v[0] <= status && status <= v[1] ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

// expectedStatuses returns expectedStatuses object based on the provided arguments.
// The arguments must be either integers or object of `***REMOVED***min: <integer>, max: <integer>***REMOVED***`
// kind. The "integer"ness is checked by the Number.isInteger.
func (mi *ModuleInstance) expectedStatuses(args ...goja.Value) *expectedStatuses ***REMOVED***
	rt := mi.vu.Runtime()

	if len(args) == 0 ***REMOVED***
		common.Throw(rt, errors.New("no arguments"))
	***REMOVED***
	var result expectedStatuses

	jsIsInt, _ := goja.AssertFunction(rt.GlobalObject().Get("Number").ToObject(rt).Get("isInteger"))
	isInt := func(a goja.Value) bool ***REMOVED***
		v, err := jsIsInt(goja.Undefined(), a)
		return err == nil && v.ToBoolean()
	***REMOVED***

	errMsg := "argument number %d to expectedStatuses was neither an integer nor an object like ***REMOVED***min:100, max:329***REMOVED***"
	for i, arg := range args ***REMOVED***
		o := arg.ToObject(rt)
		if o == nil ***REMOVED***
			common.Throw(rt, fmt.Errorf(errMsg, i+1))
		***REMOVED***

		if isInt(arg) ***REMOVED***
			result.exact = append(result.exact, int(o.ToInteger()))
		***REMOVED*** else ***REMOVED***
			min := o.Get("min")
			max := o.Get("max")
			if min == nil || max == nil ***REMOVED***
				common.Throw(rt, fmt.Errorf(errMsg, i+1))
			***REMOVED***
			if !(isInt(min) && isInt(max)) ***REMOVED***
				common.Throw(rt, fmt.Errorf("both min and max need to be integers for argument number %d", i+1))
			***REMOVED***

			result.minmax = append(result.minmax, [2]int***REMOVED***int(min.ToInteger()), int(max.ToInteger())***REMOVED***)
		***REMOVED***
	***REMOVED***
	return &result
***REMOVED***

// SetResponseCallback sets the responseCallback to the value provided. Supported values are
// expectedStatuses object or a `null` which means that metrics shouldn't be tagged as failed and
// `http_req_failed` should not be emitted - the behaviour previous to this
func (c *Client) SetResponseCallback(val goja.Value) ***REMOVED***
	if val != nil && !goja.IsNull(val) ***REMOVED***
		// This is done this way as ExportTo exports functions to empty structs without an error
		if es, ok := val.Export().(*expectedStatuses); ok ***REMOVED***
			c.responseCallback = es.match
		***REMOVED*** else ***REMOVED***
			common.Throw(
				c.moduleInstance.vu.Runtime(),
				fmt.Errorf("unsupported argument, expected http.expectedStatuses"),
			)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		c.responseCallback = nil
	***REMOVED***
***REMOVED***
