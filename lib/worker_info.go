package lib

import "github.com/go-redis/redis/v9"

type WorkerInfo struct {
	Client         *redis.Client
	JobId          string
	ScopeId        string
	OrchestratorId string
	WorkerId       string
}

func GetTestWorkerInfo() *WorkerInfo {
	return &WorkerInfo{
		Client: redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		}),
		JobId:          "4d2b8a88-07e6-4e70-9a53-de45c273b3d6",
		ScopeId:        "7faae966-211d-4b41-a9da-d9ae634ad085",
		OrchestratorId: "33f39131-3cec-4e9c-aff9-66d7c7b0e4b8",
		WorkerId:       "46221780-2f61-4733-a181-9d34684734b9",
	}
}
