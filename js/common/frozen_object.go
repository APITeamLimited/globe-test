package common

import (
	"fmt"

	"github.com/dop251/goja"
)

// FreezeObject replicates the JavaScript Object.freeze function.
func FreezeObject(rt *goja.Runtime, obj goja.Value) error ***REMOVED***
	global := rt.GlobalObject().Get("Object").ToObject(rt)
	freeze, ok := goja.AssertFunction(global.Get("freeze"))
	if !ok ***REMOVED***
		panic("failed to get the Object.freeze function from the runtime")
	***REMOVED***
	isFrozen, ok := goja.AssertFunction(global.Get("isFrozen"))
	if !ok ***REMOVED***
		panic("failed to get the Object.isFrozen function from the runtime")
	***REMOVED***
	fobj := &freezing***REMOVED***
		global:   global,
		rt:       rt,
		freeze:   freeze,
		isFrozen: isFrozen,
	***REMOVED***
	return fobj.deepFreeze(obj)
***REMOVED***

type freezing struct ***REMOVED***
	rt       *goja.Runtime
	global   goja.Value
	freeze   goja.Callable
	isFrozen goja.Callable
***REMOVED***

func (f *freezing) deepFreeze(val goja.Value) error ***REMOVED***
	if val != nil && goja.IsNull(val) ***REMOVED***
		return nil
	***REMOVED***

	_, err := f.freeze(goja.Undefined(), val)
	if err != nil ***REMOVED***
		return fmt.Errorf("object freeze failed: %w", err)
	***REMOVED***

	o := val.ToObject(f.rt)
	if o == nil ***REMOVED***
		return nil
	***REMOVED***

	for _, key := range o.Keys() ***REMOVED***
		prop := o.Get(key)
		if prop == nil ***REMOVED***
			continue
		***REMOVED***
		frozen, err := f.isFrozen(goja.Undefined(), prop)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if frozen.ToBoolean() ***REMOVED*** // prevent cycles
			continue
		***REMOVED***
		if err = f.deepFreeze(prop); err != nil ***REMOVED***
			return fmt.Errorf("deep freezing the property %s failed: %w", key, err)
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***
