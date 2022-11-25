package orchestrator

import (
	"context"
	"crypto/tls"
	"fmt"

	"github.com/APITeamLimited/globe-test/lib"
	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/APITeamLimited/redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func tryGetClient(currentIndex int) *libOrch.NamedClient {
	host := lib.GetEnvVariableHideError(fmt.Sprintf("WORKER_%d_HOST", currentIndex), "NONE", false)
	port := lib.GetEnvVariableHideError(fmt.Sprintf("WORKER_%d_PORT", currentIndex), "NONE", false)
	password := lib.GetEnvVariableHideError(fmt.Sprintf("WORKER_%d_PASSWORD", currentIndex), "NONE", false)
	displayName := lib.GetEnvVariableHideError(fmt.Sprintf("WORKER_%d_DISPLAY_NAME", currentIndex), "NONE", false)

	if host == "NONE" || port == "NONE" || password == "NONE" || displayName == "NONE" {
		return nil
	}

	options := &redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password,
		DB:       0,
	}

	isSecure := lib.GetEnvVariable(fmt.Sprintf("WORKER_%d_IS_SECURE", currentIndex), "false") == "true"

	if isSecure {
		clientCert := lib.GetEnvVariable(fmt.Sprintf("WORKER_%d_CERT", currentIndex), "")
		clientKey := lib.GetEnvVariable(fmt.Sprintf("WORKER_%d_KEY", currentIndex), "")

		cert, err := tls.X509KeyPair([]byte(clientCert), []byte(clientKey))
		if err != nil {
			panic(fmt.Errorf("error loading orchestrator cert: %s", err))
		}

		options.TLSConfig = &tls.Config{
			MinVersion:         tls.VersionTLS12,
			InsecureSkipVerify: lib.GetEnvVariable(fmt.Sprintf("WORKER_%d_INSECURE_SKIP_VERIFY", currentIndex), "false") == "true",
			Certificates:       []tls.Certificate{cert},
		}
	}

	return &libOrch.NamedClient{
		Name:   displayName,
		Client: redis.NewClient(options),
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
	storeMongoUser := lib.GetEnvVariable("STORE_MONGO_USER", "")
	storeMongoPassword := lib.GetEnvVariable("STORE_MONGO_PASSWORD", "")
	storeMongoHost := lib.GetEnvVariable("STORE_MONGO_HOST", "")
	storeMongoPort := lib.GetEnvVariable("STORE_MONGO_PORT", "")
	storeMongoDatabase := lib.GetEnvVariable("STORE_MONGO_DATABASE", "")

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
