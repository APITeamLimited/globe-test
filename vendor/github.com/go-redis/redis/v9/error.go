package redis

import (
	"context"
	"io"
	"net"
	"strings"

	"github.com/go-redis/redis/v9/internal/pool"
	"github.com/go-redis/redis/v9/internal/proto"
)

// ErrClosed performs any operation on the closed client will return this error.
var ErrClosed = pool.ErrClosed

type Error interface ***REMOVED***
	error

	// RedisError is a no-op function but
	// serves to distinguish types that are Redis
	// errors from ordinary errors: a type is a
	// Redis error if it has a RedisError method.
	RedisError()
***REMOVED***

var _ Error = proto.RedisError("")

func shouldRetry(err error, retryTimeout bool) bool ***REMOVED***
	switch err ***REMOVED***
	case io.EOF, io.ErrUnexpectedEOF:
		return true
	case nil, context.Canceled, context.DeadlineExceeded:
		return false
	***REMOVED***

	if v, ok := err.(timeoutError); ok ***REMOVED***
		if v.Timeout() ***REMOVED***
			return retryTimeout
		***REMOVED***
		return true
	***REMOVED***

	s := err.Error()
	if s == "ERR max number of clients reached" ***REMOVED***
		return true
	***REMOVED***
	if strings.HasPrefix(s, "LOADING ") ***REMOVED***
		return true
	***REMOVED***
	if strings.HasPrefix(s, "READONLY ") ***REMOVED***
		return true
	***REMOVED***
	if strings.HasPrefix(s, "CLUSTERDOWN ") ***REMOVED***
		return true
	***REMOVED***
	if strings.HasPrefix(s, "TRYAGAIN ") ***REMOVED***
		return true
	***REMOVED***

	return false
***REMOVED***

func isRedisError(err error) bool ***REMOVED***
	_, ok := err.(proto.RedisError)
	return ok
***REMOVED***

func isBadConn(err error, allowTimeout bool, addr string) bool ***REMOVED***
	switch err ***REMOVED***
	case nil:
		return false
	case context.Canceled, context.DeadlineExceeded:
		return true
	***REMOVED***

	if isRedisError(err) ***REMOVED***
		switch ***REMOVED***
		case isReadOnlyError(err):
			// Close connections in read only state in case domain addr is used
			// and domain resolves to a different Redis Server. See #790.
			return true
		case isMovedSameConnAddr(err, addr):
			// Close connections when we are asked to move to the same addr
			// of the connection. Force a DNS resolution when all connections
			// of the pool are recycled
			return true
		default:
			return false
		***REMOVED***
	***REMOVED***

	if allowTimeout ***REMOVED***
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***

	return true
***REMOVED***

func isMovedError(err error) (moved bool, ask bool, addr string) ***REMOVED***
	if !isRedisError(err) ***REMOVED***
		return
	***REMOVED***

	s := err.Error()
	switch ***REMOVED***
	case strings.HasPrefix(s, "MOVED "):
		moved = true
	case strings.HasPrefix(s, "ASK "):
		ask = true
	default:
		return
	***REMOVED***

	ind := strings.LastIndex(s, " ")
	if ind == -1 ***REMOVED***
		return false, false, ""
	***REMOVED***
	addr = s[ind+1:]
	return
***REMOVED***

func isLoadingError(err error) bool ***REMOVED***
	return strings.HasPrefix(err.Error(), "LOADING ")
***REMOVED***

func isReadOnlyError(err error) bool ***REMOVED***
	return strings.HasPrefix(err.Error(), "READONLY ")
***REMOVED***

func isMovedSameConnAddr(err error, addr string) bool ***REMOVED***
	redisError := err.Error()
	if !strings.HasPrefix(redisError, "MOVED ") ***REMOVED***
		return false
	***REMOVED***
	return strings.HasSuffix(redisError, " "+addr)
***REMOVED***

//------------------------------------------------------------------------------

type timeoutError interface ***REMOVED***
	Timeout() bool
***REMOVED***
