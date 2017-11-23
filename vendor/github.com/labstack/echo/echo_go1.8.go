// +build go1.8

package echo

import (
	stdContext "context"
)

// Close immediately stops the server.
// It internally calls `http.Server#Close()`.
func (e *Echo) Close() error ***REMOVED***
	if err := e.TLSServer.Close(); err != nil ***REMOVED***
		return err
	***REMOVED***
	return e.Server.Close()
***REMOVED***

// Shutdown stops server the gracefully.
// It internally calls `http.Server#Shutdown()`.
func (e *Echo) Shutdown(ctx stdContext.Context) error ***REMOVED***
	if err := e.TLSServer.Shutdown(ctx); err != nil ***REMOVED***
		return err
	***REMOVED***
	return e.Server.Shutdown(ctx)
***REMOVED***
