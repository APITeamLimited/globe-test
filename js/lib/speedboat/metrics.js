export let CounterType = 1;
export let GaugeType = 2;
export let TrendType = 3;

export class Metric ***REMOVED***
	constructor(t, name) ***REMOVED***
		if (!__initapi__) ***REMOVED***
			throw new Error("Metrics can only be created during the init phase");
		***REMOVED***
		this._impl = __initapi__.NewMetric(t, name);
	***REMOVED***
***REMOVED***

export class Counter extends Metric ***REMOVED***
	constructor(name) ***REMOVED***
		super(CounterType, name);
	***REMOVED***
***REMOVED***

export class Gauge extends Metric ***REMOVED***
	constructor(name) ***REMOVED***
		super(GaugeType, name);
	***REMOVED***
***REMOVED***

export class Trend extends Metric ***REMOVED***
	constructor(name) ***REMOVED***
		super(TrendType, name);
	***REMOVED***
***REMOVED***

export default ***REMOVED***
	CounterType: CounterType,
	GaugeType: GaugeType,
	TrendType: TrendType,
	Metric: Metric,
	Counter: Counter,
	Gauge: Gauge,
	Trend: Trend,
***REMOVED***
