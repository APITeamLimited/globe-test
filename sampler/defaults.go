package sampler

var DefaultSampler = New()

func Get(name string) *Metric ***REMOVED***
	return DefaultSampler.Get(name)
***REMOVED***

func GetAs(name string, t int) *Metric ***REMOVED***
	return DefaultSampler.GetAs(name, t)
***REMOVED***

func Gauge(name string) *Metric ***REMOVED***
	return DefaultSampler.Gauge(name)
***REMOVED***

func Stats(name string) *Metric ***REMOVED***
	return DefaultSampler.Stats(name)
***REMOVED***
