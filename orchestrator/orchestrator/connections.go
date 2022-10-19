package orchestrator

import (
	"context"
	"fmt"

	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/APITeamLimited/redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func tryGetClient(currentIndex int) *libOrch.NamedClient ***REMOVED***
	host := libOrch.GetEnvVariable(fmt.Sprintf("WORKER_%d_HOST", currentIndex), "")
	port := libOrch.GetEnvVariable(fmt.Sprintf("WORKER_%d_PORT", currentIndex), "")
	password := libOrch.GetEnvVariable(fmt.Sprintf("WORKER_%d_PASSWORD", currentIndex), "")
	displayName := libOrch.GetEnvVariable(fmt.Sprintf("WORKER_%d_DISPLAY_NAME", currentIndex), "")

	if host == "" || port == "" || password == "" || displayName == "" ***REMOVED***
		return nil
	***REMOVED***

	return &libOrch.NamedClient***REMOVED***
		Name: displayName,
		Client: redis.NewClient(&redis.Options***REMOVED***
			Addr:     fmt.Sprintf("%s:%s", host, port),
			Password: password,
			DB:       0,
		***REMOVED***),
	***REMOVED***
***REMOVED***

func connectWorkerClients(ctx context.Context) libOrch.WorkerClients ***REMOVED***
	workerClients := libOrch.WorkerClients***REMOVED***
		Clients: make(map[string]*libOrch.NamedClient),
	***REMOVED***

	currentIndex := 0

	for ***REMOVED***
		namedClient := tryGetClient(currentIndex)

		if namedClient == nil ***REMOVED***
			if currentIndex == 0 ***REMOVED***
				panic("At least one worker client must be defined")
			***REMOVED***

			break
		***REMOVED***

		if currentIndex == 0 ***REMOVED***
			workerClients.DefaultClient = namedClient
		***REMOVED***

		workerClients.Clients[namedClient.Name] = namedClient

		currentIndex++
	***REMOVED***

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
