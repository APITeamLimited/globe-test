package log

import (
	log "github.com/Sirupsen/logrus"
)

func Log(logger *log.Logger, t, msg string, fields map[string]interface***REMOVED******REMOVED***) ***REMOVED***
	e := logger.WithFields(log.Fields(fields))
	switch t ***REMOVED***
	case "error":
		e.Error(msg)
	case "warn":
		e.Warn(msg)
	case "info":
		e.Info(msg)
	case "debug":
		e.Debug(msg)
	***REMOVED***
***REMOVED***
