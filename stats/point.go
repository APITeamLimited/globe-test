package stats

import (
	"time"
)

type Tags map[string]interface***REMOVED******REMOVED***
type Values map[string]float64

type Point struct ***REMOVED***
	Stat   *Stat
	Time   time.Time
	Tags   Tags
	Values Values
***REMOVED***

func Value(val float64) Values ***REMOVED***
	return Values***REMOVED***"value": val***REMOVED***
***REMOVED***
