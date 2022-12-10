package orchestrator

import (
	"context"
	"crypto/tls"
	"fmt"

	"github.com/APITeamLimited/globe-test/lib"
	"github.com/APITeamLimited/globe-test/lib/agent"
	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/APITeamLimited/redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func tryGetClient(currentIndex int) *libOrch.NamedClient ***REMOVED***
	host := lib.GetEnvVariableHideError(fmt.Sprintf("WORKER_%d_HOST", currentIndex), "NONE", false)
	port := lib.GetEnvVariableHideError(fmt.Sprintf("WORKER_%d_PORT", currentIndex), "NONE", false)
	password := lib.GetEnvVariableHideError(fmt.Sprintf("WORKER_%d_PASSWORD", currentIndex), "NONE", false)
	displayName := lib.GetEnvVariableHideError(fmt.Sprintf("WORKER_%d_DISPLAY_NAME", currentIndex), "NONE", false)

	if host == "NONE" || port == "NONE" || password == "NONE" || displayName == "NONE" ***REMOVED***
		return nil
	***REMOVED***

	options := &redis.Options***REMOVED***
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password,
		DB:       0,
	***REMOVED***

	isSecure := lib.GetEnvVariable(fmt.Sprintf("WORKER_%d_IS_SECURE", currentIndex), "false") == "true"

	if isSecure ***REMOVED***
		clientCert := lib.GetEnvVariable(fmt.Sprintf("WORKER_%d_CERT", currentIndex), "")
		clientKey := lib.GetEnvVariable(fmt.Sprintf("WORKER_%d_KEY", currentIndex), "")

		cert, err := tls.X509KeyPair([]byte(clientCert), []byte(clientKey))
		if err != nil ***REMOVED***
			panic(fmt.Errorf("error loading orchestrator cert: %s", err))
		***REMOVED***

		options.TLSConfig = &tls.Config***REMOVED***
			MinVersion:         tls.VersionTLS12,
			InsecureSkipVerify: lib.GetEnvVariable(fmt.Sprintf("WORKER_%d_INSECURE_SKIP_VERIFY", currentIndex), "false") == "true",
			Certificates:       []tls.Certificate***REMOVED***cert***REMOVED***,
		***REMOVED***
	***REMOVED***

	return &libOrch.NamedClient***REMOVED***
		Name:   displayName,
		Client: redis.NewClient(options),
	***REMOVED***
***REMOVED***

func connectWorkerClients(ctx context.Context, standalone bool) libOrch.WorkerClients ***REMOVED***
	workerClients := libOrch.WorkerClients***REMOVED***
		Clients: make(map[string]*libOrch.NamedClient),
	***REMOVED***

	if !standalone ***REMOVED***
		// Just get a single worker client for the agent
		workerClients.Clients[agent.AgentWorkerName] = &libOrch.NamedClient***REMOVED***
			Name: agent.AgentWorkerName,
			Client: redis.NewClient(&redis.Options***REMOVED***
				Addr: fmt.Sprintf("%s:%s", agent.WorkerRedisHost, agent.WorkerRedisPort),
			***REMOVED***),
		***REMOVED***
		workerClients.DefaultClient = workerClients.Clients[agent.AgentWorkerName]

		return workerClients
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

func getStoreMongoDB(ctx context.Context, standalone bool) *mongo.Database ***REMOVED***
	if !standalone ***REMOVED***
		return nil
	***REMOVED***

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
