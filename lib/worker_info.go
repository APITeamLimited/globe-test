package lib

import "github.com/go-redis/redis/v9"

type WorkerInfo struct ***REMOVED***
	Client *redis.Client
***REMOVED***

func GetTestWorkerInfo() *WorkerInfo ***REMOVED***
	return &WorkerInfo***REMOVED***
		Client: redis.NewClient(&redis.Options***REMOVED***
			Addr:     "localhost:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		***REMOVED***),
	***REMOVED***
***REMOVED***
