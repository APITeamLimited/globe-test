package http

type SlotLimiter struct ***REMOVED***
	ch chan struct***REMOVED******REMOVED***
***REMOVED***

func NewSlotLimiter(slots int) SlotLimiter ***REMOVED***
	if slots <= 0 ***REMOVED***
		return SlotLimiter***REMOVED***nil***REMOVED***
	***REMOVED***

	ch := make(chan struct***REMOVED******REMOVED***, slots)
	for i := 0; i < slots; i++ ***REMOVED***
		ch <- struct***REMOVED******REMOVED******REMOVED******REMOVED***
	***REMOVED***
	return SlotLimiter***REMOVED***ch***REMOVED***
***REMOVED***

func (l *SlotLimiter) Begin() ***REMOVED***
	if l.ch != nil ***REMOVED***
		<-l.ch
	***REMOVED***
***REMOVED***

func (l *SlotLimiter) End() ***REMOVED***
	if l.ch != nil ***REMOVED***
		l.ch <- struct***REMOVED******REMOVED******REMOVED******REMOVED***
	***REMOVED***
***REMOVED***
