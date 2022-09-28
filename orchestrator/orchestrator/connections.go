package orchestrator

import (
	"context"
	"fmt"

	"github.com/APITeamLimited/redis/v9"
	"gitlab.com/apiteamcloud/orchestrator/lib"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func connectWorkerClients(ctx context.Context) map[string]*redis.Client ***REMOVED***
	workerClients := make(map[string]*redis.Client)

	workerClients["portsmouth"] = redis.NewClient(&redis.Options***REMOVED***
		Addr:     fmt.Sprintf("%s:%s", lib.GetEnvVariable("WORKER_PORTSMOUTH_REDIS_HOST", "localhost"), lib.GetEnvVariable("WORKER_PORTSMOUTH_REDIS_PORT", "10002")),
		Password: lib.GetEnvVariable("WORKER_PORTSMOUTH_REDIS_PASSWORD", ""),
		DB:       0,
	***REMOVED***)

	return workerClients
***REMOVED***

func getStoreMongoDB(ctx context.Context) *mongo.Database ***REMOVED***
	storeMongoUser := lib.GetEnvVariable("STORE_MONGO_USER", "")
	storeMongoPassword := lib.GetEnvVariable("STORE_MONGO_PASSWORD", "")
	storeMongoHost := lib.GetEnvVariable("STORE_MONGO_HOST", "")
	storeMongoPort := lib.GetEnvVariable("STORE_MONGO_PORT", "")
	storeMongoDatabase := lib.GetEnvVariable("STORE_MONGO_DATABASE", "")

	storeURI := fmt.Sprintf("mongodb://%s:%s@%s:%s", storeMongoUser, storeMongoPassword, storeMongoHost, storeMongoPort)

	client, err := mongo.NewClient(options.Client().ApplyURI(storeURI))
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***

	if err := client.Connect(ctx); err != nil ***REMOVED***
		panic(err)
	***REMOVED***

	mongoDB := client.Database(storeMongoDatabase)

	return mongoDB
***REMOVED***
