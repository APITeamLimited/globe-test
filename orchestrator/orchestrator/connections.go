package orchestrator

import (
	"context"
	"fmt"

	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/APITeamLimited/redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func tryGetClient(currentIndex int) *libOrch.NamedClient {
	host := libOrch.GetEnvVariable(fmt.Sprintf("WORKER_%d_HOST", currentIndex), "")
	port := libOrch.GetEnvVariable(fmt.Sprintf("WORKER_%d_PORT", currentIndex), "")
	password := libOrch.GetEnvVariable(fmt.Sprintf("WORKER_%d_PASSWORD", currentIndex), "")
	displayName := libOrch.GetEnvVariable(fmt.Sprintf("WORKER_%d_DISPLAY_NAME", currentIndex), "")

	if host == "" || port == "" || password == "" || displayName == "" {
		return nil
	}

	return &libOrch.NamedClient{
		Name: displayName,
		Client: redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", host, port),
			Password: password,
			DB:       0,
		}),
	}
}

func connectWorkerClients(ctx context.Context) libOrch.WorkerClients {
	workerClients := libOrch.WorkerClients{
		Clients: make(map[string]*libOrch.NamedClient),
	}

	currentIndex := 0

	for {
		namedClient := tryGetClient(currentIndex)

		if namedClient == nil {
			if currentIndex == 0 {
				panic("At least one worker client must be defined")
			}

			break
		}

		if currentIndex == 0 {
			workerClients.DefaultClient = namedClient
		}

		workerClients.Clients[namedClient.Name] = namedClient

		currentIndex++
	}

	return workerClients
}

func getStoreMongoDB(ctx context.Context) *mongo.Database {
	storeMongoUser := libOrch.GetEnvVariable("STORE_MONGO_USER", "")
	storeMongoPassword := libOrch.GetEnvVariable("STORE_MONGO_PASSWORD", "")
	storeMongoHost := libOrch.GetEnvVariable("STORE_MONGO_HOST", "")
	storeMongoPort := libOrch.GetEnvVariable("STORE_MONGO_PORT", "")
	storeMongoDatabase := libOrch.GetEnvVariable("STORE_MONGO_DATABASE", "")

	storeURI := fmt.Sprintf("mongodb://%s:%s@%s:%s", storeMongoUser, storeMongoPassword, storeMongoHost, storeMongoPort)

	client, err := mongo.NewClient(options.Client().ApplyURI(storeURI))
	if err != nil {
		panic(err)
	}

	if err := client.Connect(ctx); err != nil {
		panic(err)
	}

	mongoDB := client.Database(storeMongoDatabase)

	return mongoDB
}
