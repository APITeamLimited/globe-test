package common

import (
	"context"
	"github.com/loadimpact/k6/lib"
)

type ContextKey int

const ctxKeyEngine = ContextKey(1)

func WithEngine(ctx context.Context, engine *lib.Engine) context.Context ***REMOVED***
	return context.WithValue(ctx, ctxKeyEngine, engine)
***REMOVED***

func GetEngine(ctx context.Context) *lib.Engine ***REMOVED***
	return ctx.Value(ctxKeyEngine).(*lib.Engine)
***REMOVED***
