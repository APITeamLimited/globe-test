package internal

import (
	"context"
	"time"

	"github.com/go-redis/redis/v9/internal/util"
)

func Sleep(ctx context.Context, dur time.Duration) error ***REMOVED***
	t := time.NewTimer(dur)
	defer t.Stop()

	select ***REMOVED***
	case <-t.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	***REMOVED***
***REMOVED***

func ToLower(s string) string ***REMOVED***
	if isLower(s) ***REMOVED***
		return s
	***REMOVED***

	b := make([]byte, len(s))
	for i := range b ***REMOVED***
		c := s[i]
		if c >= 'A' && c <= 'Z' ***REMOVED***
			c += 'a' - 'A'
		***REMOVED***
		b[i] = c
	***REMOVED***
	return util.BytesToString(b)
***REMOVED***

func isLower(s string) bool ***REMOVED***
	for i := 0; i < len(s); i++ ***REMOVED***
		c := s[i]
		if c >= 'A' && c <= 'Z' ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***
