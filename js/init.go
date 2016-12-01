package js

import (
	"errors"
	"fmt"
	"github.com/loadimpact/k6/stats"
	"github.com/robertkrimen/otto"
	"strings"
)

type InitAPI struct ***REMOVED***
	r *Runtime
***REMOVED***

func (i InitAPI) NewMetric(it int, name string, isTime bool) *stats.Metric ***REMOVED***
	t := stats.MetricType(it)
	vt := stats.Default
	if isTime ***REMOVED***
		vt = stats.Time
	***REMOVED***

	if m, ok := i.r.Metrics[name]; ok ***REMOVED***
		if m.Type != t ***REMOVED***
			throw(i.r.VM, errors.New(fmt.Sprintf("attempted to redeclare %s with a different type (%s != %s)", name, m.Type, t)))
			return nil
		***REMOVED***
		if m.Contains != vt ***REMOVED***
			throw(i.r.VM, errors.New(fmt.Sprintf("attempted to redeclare %s with a different kind of value (%s != %s)", name, m.Contains, vt)))
		***REMOVED***
		return m
	***REMOVED***

	m := stats.New(name, t, vt)
	i.r.Metrics[name] = m
	return m
***REMOVED***

func (i InitAPI) Require(name string) otto.Value ***REMOVED***
	if !strings.HasPrefix(name, ".") ***REMOVED***
		exports, err := i.r.loadLib(name + ".js")
		if err != nil ***REMOVED***
			throw(i.r.VM, err)
		***REMOVED***
		return exports
	***REMOVED***

	exports, err := i.r.loadFile(name + ".js")
	if err != nil ***REMOVED***
		throw(i.r.VM, err)
	***REMOVED***
	return exports
***REMOVED***
