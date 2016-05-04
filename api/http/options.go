package http

func (ctx *context) SetMaxConnsPerHost(conns int) ***REMOVED***
	ctx.client.MaxConnsPerHost = conns
***REMOVED***

func (ctx *context) SetDefaults(args RequestArgs) ***REMOVED***
	ctx.defaults = args
***REMOVED***
