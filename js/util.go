package js

import (
	"github.com/robertkrimen/otto"
)

func paramsFromObject(o *otto.Object) (params HTTPParams, err error) ***REMOVED***
	if o == nil ***REMOVED***
		return params, nil
	***REMOVED***

	for _, key := range o.Keys() ***REMOVED***
		switch key ***REMOVED***
		case "quiet":
			v, err := o.Get(key)
			if err != nil ***REMOVED***
				return params, err
			***REMOVED***
			quiet, err := v.ToBoolean()
			if err != nil ***REMOVED***
				return params, err
			***REMOVED***
			params.Quiet = quiet
		case "headers":
			v, err := o.Get(key)
			if err != nil ***REMOVED***
				return params, err
			***REMOVED***
			obj := v.Object()
			if obj == nil ***REMOVED***
				continue
			***REMOVED***

			params.Headers = make(map[string]string)
			for _, name := range obj.Keys() ***REMOVED***
				hv, err := obj.Get(name)
				if err != nil ***REMOVED***
					return params, err
				***REMOVED***
				value, err := hv.ToString()
				if err != nil ***REMOVED***
					return params, err
				***REMOVED***
				params.Headers[name] = value
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return params, nil
***REMOVED***
