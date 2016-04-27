package runner

import (
	log "github.com/Sirupsen/logrus"
	"github.com/loadimpact/speedboat/loadtest"
	"golang.org/x/net/context"
	"sync"
	"time"
)

type Runner interface ***REMOVED***
	Run(ctx context.Context, id int64) <-chan Result
***REMOVED***

type Result struct ***REMOVED***
	Text  string
	Time  time.Duration
	Error error
***REMOVED***

type VU struct ***REMOVED***
	Cancel context.CancelFunc
***REMOVED***

func Run(ctx context.Context, r Runner, scale <-chan int) <-chan Result ***REMOVED***
	ch := make(chan Result)

	go func() ***REMOVED***
		wg := sync.WaitGroup***REMOVED******REMOVED***
		defer func() ***REMOVED***
			wg.Wait()
			close(ch)
		***REMOVED***()

		currentVUs := make([]VU, 0, 100)
		currentID := int64(0)
		for ***REMOVED***
			select ***REMOVED***
			case vus := <-scale:
				for vus > len(currentVUs) ***REMOVED***
					currentID += 1
					currentID := currentID
					wg.Add(1)
					c, cancel := context.WithCancel(ctx)
					currentVUs = append(currentVUs, VU***REMOVED***Cancel: cancel***REMOVED***)
					go func() ***REMOVED***
						defer wg.Done()
						for res := range r.Run(c, currentID) ***REMOVED***
							ch <- res
						***REMOVED***
					***REMOVED***()
				***REMOVED***
				for vus < len(currentVUs) ***REMOVED***
					currentVUs[len(currentVUs)-1].Cancel()
					currentVUs = currentVUs[:len(currentVUs)-1]
				***REMOVED***
			case <-ctx.Done():
				return
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	return ch
***REMOVED***

func Ramp(t *loadtest.LoadTest, scale chan int, in <-chan Result) <-chan Result ***REMOVED***
	ch := make(chan Result)

	go func() ***REMOVED***
		defer close(ch)

		ticker := time.NewTicker(time.Duration(1) * time.Second)
		startTime := time.Now()
		for ***REMOVED***
			select ***REMOVED***
			case <-ticker.C:
				vus, _ := t.VUsAt(time.Since(startTime))
				scale <- vus
			case res, ok := <-in:
				if !ok ***REMOVED***
					return
				***REMOVED***
				ch <- res
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	return ch
***REMOVED***
