package worker

import (
	"context"

	"github.com/go-redis/redis/v9"
	"github.com/sirupsen/logrus"
)

// RedisHook to sends logs to Redis server
type RedisHook struct ***REMOVED***
	ctx      context.Context
	client   *redis.Client
	jobId    string
	workerId string
***REMOVED***

// NewRedisHook creates a hook to be added to an instance of logger
func NewRedisHook(client *redis.Client, ctx context.Context, jobId string, workerId string) (*RedisHook, error) ***REMOVED***
	return &RedisHook***REMOVED***
		client:   client,
		ctx:      ctx,
		jobId:    jobId,
		workerId: workerId,
	***REMOVED***, nil

***REMOVED***

// Fire is called when a log event is fired.
func (hook *RedisHook) Fire(entry *logrus.Entry) error ***REMOVED***
	// This doesn't work for some reason but redundant
	//go dispatchMessage(hook.ctx, hook.client, hook.jobId, hook.workerId, entry.Message, "MESSAGE")
	return nil
***REMOVED***

// Levels returns the available logging levels.
func (hook *RedisHook) Levels() []logrus.Level ***REMOVED***
	return []logrus.Level***REMOVED***
		logrus.TraceLevel,
		logrus.DebugLevel,
		logrus.InfoLevel,
		logrus.WarnLevel,
		logrus.ErrorLevel,
		logrus.FatalLevel,
		logrus.PanicLevel,
	***REMOVED***
***REMOVED***
