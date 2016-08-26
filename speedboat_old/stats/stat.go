package stats

import (
	"time"
)

type StatType int
type StatIntent int

const (
	CounterType StatType = iota
	GaugeType
	HistogramType
)

const (
	DefaultIntent StatIntent = iota
	TimeIntent
)

type Stat struct ***REMOVED***
	Name   string
	Type   StatType
	Intent StatIntent
***REMOVED***

func ApplyIntent(v float64, intent StatIntent) interface***REMOVED******REMOVED*** ***REMOVED***
	switch intent ***REMOVED***
	case TimeIntent:
		return time.Duration(v)
	default:
		return v
	***REMOVED***
***REMOVED***
