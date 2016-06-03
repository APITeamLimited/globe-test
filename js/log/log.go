package log

import (
	log "github.com/Sirupsen/logrus"
	"golang.org/x/net/context"
	"os"
)

type ContextKey int

const (
	loggerKey = ContextKey(iota)
)

func WithDefaultLogger(ctx context.Context) context.Context ***REMOVED***
	logger := log.New()
	logger.Out = os.Stdout
	logger.Level = log.DebugLevel
	return WithLogger(ctx, logger)
***REMOVED***

func WithLogger(ctx context.Context, c *log.Logger) context.Context ***REMOVED***
	return context.WithValue(ctx, loggerKey, c)
***REMOVED***

func GetLogger(ctx context.Context) *log.Logger ***REMOVED***
	return ctx.Value(loggerKey).(*log.Logger)
***REMOVED***

func Log(ctx context.Context, t, msg string, fields map[string]interface***REMOVED******REMOVED***) ***REMOVED***
	logger := GetLogger(ctx)
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
