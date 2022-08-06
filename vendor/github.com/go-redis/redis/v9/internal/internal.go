package internal

import (
	"time"

	"github.com/go-redis/redis/v9/internal/rand"
)

func RetryBackoff(retry int, minBackoff, maxBackoff time.Duration) time.Duration ***REMOVED***
	if retry < 0 ***REMOVED***
		panic("not reached")
	***REMOVED***
	if minBackoff == 0 ***REMOVED***
		return 0
	***REMOVED***

	d := minBackoff << uint(retry)
	if d < minBackoff ***REMOVED***
		return maxBackoff
	***REMOVED***

	d = minBackoff + time.Duration(rand.Int63n(int64(d)))

	if d > maxBackoff || d < minBackoff ***REMOVED***
		d = maxBackoff
	***REMOVED***

	return d
***REMOVED***
