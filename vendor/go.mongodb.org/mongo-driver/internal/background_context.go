// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package internal

import "context"

// backgroundContext is an implementation of the context.Context interface that wraps a child Context. Value requests
// are forwarded to the child Context but the Done and Err functions are overridden to ensure the new context does not
// time out or get cancelled.
type backgroundContext struct ***REMOVED***
	context.Context
	childValuesCtx context.Context
***REMOVED***

// NewBackgroundContext creates a new Context whose behavior matches that of context.Background(), but Value calls are
// forwarded to the provided ctx parameter. If ctx is nil, context.Background() is returned.
func NewBackgroundContext(ctx context.Context) context.Context ***REMOVED***
	if ctx == nil ***REMOVED***
		return context.Background()
	***REMOVED***

	return &backgroundContext***REMOVED***
		Context:        context.Background(),
		childValuesCtx: ctx,
	***REMOVED***
***REMOVED***

func (b *backgroundContext) Value(key interface***REMOVED******REMOVED***) interface***REMOVED******REMOVED*** ***REMOVED***
	return b.childValuesCtx.Value(key)
***REMOVED***
