package metrics

// Healthchecks hold an error value describing an arbitrary up/down status.
type Healthcheck interface ***REMOVED***
	Check()
	Error() error
	Healthy()
	Unhealthy(error)
***REMOVED***

// NewHealthcheck constructs a new Healthcheck which will use the given
// function to update its status.
func NewHealthcheck(f func(Healthcheck)) Healthcheck ***REMOVED***
	if UseNilMetrics ***REMOVED***
		return NilHealthcheck***REMOVED******REMOVED***
	***REMOVED***
	return &StandardHealthcheck***REMOVED***nil, f***REMOVED***
***REMOVED***

// NilHealthcheck is a no-op.
type NilHealthcheck struct***REMOVED******REMOVED***

// Check is a no-op.
func (NilHealthcheck) Check() ***REMOVED******REMOVED***

// Error is a no-op.
func (NilHealthcheck) Error() error ***REMOVED*** return nil ***REMOVED***

// Healthy is a no-op.
func (NilHealthcheck) Healthy() ***REMOVED******REMOVED***

// Unhealthy is a no-op.
func (NilHealthcheck) Unhealthy(error) ***REMOVED******REMOVED***

// StandardHealthcheck is the standard implementation of a Healthcheck and
// stores the status and a function to call to update the status.
type StandardHealthcheck struct ***REMOVED***
	err error
	f   func(Healthcheck)
***REMOVED***

// Check runs the healthcheck function to update the healthcheck's status.
func (h *StandardHealthcheck) Check() ***REMOVED***
	h.f(h)
***REMOVED***

// Error returns the healthcheck's status, which will be nil if it is healthy.
func (h *StandardHealthcheck) Error() error ***REMOVED***
	return h.err
***REMOVED***

// Healthy marks the healthcheck as healthy.
func (h *StandardHealthcheck) Healthy() ***REMOVED***
	h.err = nil
***REMOVED***

// Unhealthy marks the healthcheck as unhealthy.  The error is stored and
// may be retrieved by the Error method.
func (h *StandardHealthcheck) Unhealthy(err error) ***REMOVED***
	h.err = err
***REMOVED***
