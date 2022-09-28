package orchestrator

import (
	"context"
	"fmt"

	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/APITeamLimited/redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func connectWorkerClients(ctx context.Context) map[string]*redis.Client ***REMOVED***
	workerClients := make(map[string]*redis.Client)

	workerClients["portsmouth"] = redis.NewClient(&redis.Options***REMOVED***
		Addr:     fmt.Sprintf("%s:%s", libOrch.GetEnvVariable("WORKER_PORTSMOUTH_REDIS_HOST", "localhost"), libOrch.GetEnvVariable("WORKER_PORTSMOUTH_REDIS_PORT", "10002")),
		Password: libOrch.GetEnvVariable("WORKER_PORTSMOUTH_REDIS_PASSWORD", ""),
		DB:       0,
	***REMOVED***)

	return workerClients
***REMOVED***

func getStoreMongoDB(ctx context.Context) *mongo.Database ***REMOVED***
	storeMongoUser := libOrch.GetEnvVariable("STORE_MONGO_USER", "")
	storeMongoPassword := libOrch.GetEnvVariable("STORE_MONGO_PASSWORD", "")
	storeMongoHost := libOrch.GetEnvVariable("STORE_MONGO_HOST", "")
	storeMongoPort := libOrch.GetEnvVariable("STORE_MONGO_PORT", "")
	storeMongoDatabase := libOrch.GetEnvVariable("STORE_MONGO_DATABASE", "")

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
