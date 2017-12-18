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

type MultiSlotLimiter struct ***REMOVED***
	m     map[string]*SlotLimiter
	slots int
***REMOVED***

func NewMultiSlotLimiter(slots int) MultiSlotLimiter ***REMOVED***
	return MultiSlotLimiter***REMOVED***make(map[string]*SlotLimiter), slots***REMOVED***
***REMOVED***

func (l *MultiSlotLimiter) Slot(s string) *SlotLimiter ***REMOVED***
	if l.slots == 0 ***REMOVED***
		return nil
	***REMOVED***
	ll, ok := l.m[s]
	if !ok ***REMOVED***
		tmp := NewSlotLimiter(l.slots)
		ll = &tmp
		l.m[s] = ll
	***REMOVED***
	return ll
***REMOVED***
