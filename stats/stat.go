package stats

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
