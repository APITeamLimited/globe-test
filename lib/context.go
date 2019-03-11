package lib

import "context"

type ctxKey int

const (
	ctxKeyState ctxKey = iota
)

func WithState(ctx context.Context, state *State) context.Context ***REMOVED***
	return context.WithValue(ctx, ctxKeyState, state)
***REMOVED***

func GetState(ctx context.Context) *State ***REMOVED***
	v := ctx.Value(ctxKeyState)
	if v == nil ***REMOVED***
		return nil
	***REMOVED***
	return v.(*State)
***REMOVED***
