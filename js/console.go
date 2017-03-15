package js

import (
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/robertkrimen/otto"
)

type Console struct ***REMOVED***
	Logger *log.Logger
***REMOVED***

func (c Console) Log(level int, msg string, args []otto.Value) ***REMOVED***
	fields := make(log.Fields, len(args))
	for i, arg := range args ***REMOVED***
		if arg.IsObject() ***REMOVED***
			obj := arg.Object()
			for _, key := range obj.Keys() ***REMOVED***
				v, _ := obj.Get(key)
				fields[key] = v.String()
			***REMOVED***
			continue
		***REMOVED***
		fields["arg"+strconv.Itoa(i)] = arg.String()
	***REMOVED***

	entry := c.Logger.WithFields(fields)
	switch level ***REMOVED***
	case 0:
		entry.Debug(msg)
	case 1:
		entry.Info(msg)
	case 2:
		entry.Warn(msg)
	case 3:
		entry.Error(msg)
	***REMOVED***
***REMOVED***
