package lib

import "github.com/go-redis/redis/v9"

type WorkerInfo struct {
	Client *redis.Client
}

func GetTestWorkerInfo() *WorkerInfo {
	return &WorkerInfo{
		Client: redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		}),
	}
}
