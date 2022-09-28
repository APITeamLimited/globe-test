package output

import (
	"fmt"
	"sync"
	"time"

	"github.com/APITeamLimited/k6-worker/metrics"
)

// SampleBuffer is a simple thread-safe buffer for metric samples. It should be
// used by most outputs, since we generally want to flush metric samples to the
// remote service asynchronously. We want to do it only every several seconds,
// and we don't want to block the Engine in the meantime.
type SampleBuffer struct ***REMOVED***
	sync.Mutex
	buffer []metrics.SampleContainer
	maxLen int
***REMOVED***

// AddMetricSamples adds the given metric samples to the internal buffer.
func (sc *SampleBuffer) AddMetricSamples(samples []metrics.SampleContainer) ***REMOVED***
	if len(samples) == 0 ***REMOVED***
		return
	***REMOVED***
	sc.Lock()
	sc.buffer = append(sc.buffer, samples...)
	sc.Unlock()
***REMOVED***

// GetBufferedSamples returns the currently buffered metric samples and makes a
// new internal buffer with some hopefully realistic size. If the internal
// buffer is empty, it will return nil.
func (sc *SampleBuffer) GetBufferedSamples() []metrics.SampleContainer ***REMOVED***
	sc.Lock()
	defer sc.Unlock()

	buffered, bufferedLen := sc.buffer, len(sc.buffer)
	if bufferedLen == 0 ***REMOVED***
		return nil
	***REMOVED***
	if bufferedLen > sc.maxLen ***REMOVED***
		sc.maxLen = bufferedLen
	***REMOVED***
	// Make the new buffer halfway between the previously allocated size and the
	// maximum buffer size we've seen so far, to hopefully reduce copying a bit.
	sc.buffer = make([]metrics.SampleContainer, 0, (bufferedLen+sc.maxLen)/2)

	return buffered
***REMOVED***

// PeriodicFlusher is a small helper for asynchronously flushing buffered metric
// samples on regular intervals. The biggest benefit is having a Stop() method
// that waits for one last flush before it returns.
type PeriodicFlusher struct ***REMOVED***
	period        time.Duration
	flushCallback func()
	stop          chan struct***REMOVED******REMOVED***
	stopped       chan struct***REMOVED******REMOVED***
	once          *sync.Once
***REMOVED***

func (pf *PeriodicFlusher) run() ***REMOVED***
	ticker := time.NewTicker(pf.period)
	defer ticker.Stop()
	for ***REMOVED***
		select ***REMOVED***
		case <-ticker.C:
			pf.flushCallback()
		case <-pf.stop:
			pf.flushCallback()
			close(pf.stopped)
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

// Stop waits for the periodic flusher flush one last time and exit. You can
// safely call Stop() multiple times from different goroutines, you just can't
// call it from inside of the flushing function.
func (pf *PeriodicFlusher) Stop() ***REMOVED***
	pf.once.Do(func() ***REMOVED***
		close(pf.stop)
	***REMOVED***)
	<-pf.stopped
***REMOVED***

// NewPeriodicFlusher creates a new PeriodicFlusher and starts its goroutine.
func NewPeriodicFlusher(period time.Duration, flushCallback func()) (*PeriodicFlusher, error) ***REMOVED***
	if period <= 0 ***REMOVED***
		return nil, fmt.Errorf("metric flush period should be positive but was %s", period)
	***REMOVED***

	pf := &PeriodicFlusher***REMOVED***
		period:        period,
		flushCallback: flushCallback,
		stop:          make(chan struct***REMOVED******REMOVED***),
		stopped:       make(chan struct***REMOVED******REMOVED***),
		once:          &sync.Once***REMOVED******REMOVED***,
	***REMOVED***

	go pf.run()

	return pf, nil
***REMOVED***
