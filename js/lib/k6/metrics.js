export let CounterType = 0;
export let GaugeType = 1;
export let TrendType = 2;
export let RateType = 3;

export class Metric ***REMOVED***
	constructor(t, name, isTime=false) ***REMOVED***
		if (!__initapi__) ***REMOVED***
			throw new Error("Metrics can only be created during the init phase");
		***REMOVED***
		this._impl = __initapi__.NewMetric(t, name, !!isTime);
	***REMOVED***

	add(v, tags=***REMOVED******REMOVED***) ***REMOVED***
		if (!__jsapi__) ***REMOVED***
			throw new Error("Metric.add() needs VU context")
		***REMOVED***
		__jsapi__.MetricAdd(this._impl, v*1.0, tags);
		return v;
	***REMOVED***
***REMOVED***

export class Counter extends Metric ***REMOVED***
	constructor(name, isTime=false) ***REMOVED***
		super(CounterType, name, isTime);
	***REMOVED***
***REMOVED***

export class Gauge extends Metric ***REMOVED***
	constructor(name, isTime=false) ***REMOVED***
		super(GaugeType, name, isTime);
	***REMOVED***
***REMOVED***

export class Trend extends Metric ***REMOVED***
	constructor(name, isTime=false) ***REMOVED***
		super(TrendType, name, isTime);
	***REMOVED***
***REMOVED***

export class Rate extends Metric ***REMOVED***
	constructor(name, isTime=false) ***REMOVED***
		super(RateType, name, isTime);
	***REMOVED***
***REMOVED***

export default ***REMOVED***
	CounterType: CounterType,
	GaugeType: GaugeType,
	TrendType: TrendType,
	RateType: RateType,
	Metric: Metric,
	Counter: Counter,
	Gauge: Gauge,
	Trend: Trend,
	Rate: Rate,
***REMOVED***
