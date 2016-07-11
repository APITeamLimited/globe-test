package stats

type Filter struct ***REMOVED***
	exclude map[string]bool
	only    map[string]bool
***REMOVED***

func MakeFilter(exclude, only []string) Filter ***REMOVED***
	f := Filter***REMOVED******REMOVED***

	if len(exclude) > 0 ***REMOVED***
		f.exclude = make(map[string]bool)
		for _, stat := range exclude ***REMOVED***
			f.exclude[stat] = true
		***REMOVED***
	***REMOVED***
	if len(only) > 0 ***REMOVED***
		f.only = make(map[string]bool)
		for _, stat := range only ***REMOVED***
			f.only[stat] = true
		***REMOVED***
	***REMOVED***

	return f
***REMOVED***

func (f Filter) Check(s Sample) bool ***REMOVED***
	if f.only != nil ***REMOVED***
		return f.only[s.Stat.Name]
	***REMOVED***
	return !f.exclude[s.Stat.Name]
***REMOVED***
