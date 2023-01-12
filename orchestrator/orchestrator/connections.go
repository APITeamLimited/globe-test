package orchestrator

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"

	"github.com/APITeamLimited/globe-test/lib"
	"github.com/APITeamLimited/globe-test/lib/agent"
	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/APITeamLimited/redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type clientOptions struct {
	redisOptions *redis.Options
	displayName  string
}

func connectWorkerClients(ctx context.Context, standalone, independentWorkerRedisHosts bool) libOrch.WorkerClients {
	// Agent mode
	if !standalone {
		return getAgentClient()
	}

	if independentWorkerRedisHosts {
		return getIndependentClients()
	}

	return getUnifiedClient()
}

// Several independent worker redis clients
func getIndependentClients() libOrch.WorkerClients {
	workerClients := libOrch.WorkerClients{
		Clients: make(map[string]*libOrch.NamedClient),
	}

	currentIndex := 0

	for {
		clientOptions := tryGetClientOptions(fmt.Sprint(currentIndex), currentIndex)

		if clientOptions == nil {
			if currentIndex == 0 {
				panic("At least one worker client must be defined")
			}

			break
		}

		client := redis.NewClient(clientOptions.redisOptions)

		// Ensure that the client is connected
		if err := client.Ping(context.Background()).Err(); err != nil {
			panic(err)
		}

		namedClient := &libOrch.NamedClient{
			Name:   clientOptions.displayName,
			Client: client,
		}

		if currentIndex == 0 {
			workerClients.DefaultClient = namedClient
		}

		workerClients.Clients[namedClient.Name] = namedClient

		currentIndex++
	}

	return workerClients
}

// A unified worker redis client, supporting different display names
func getUnifiedClient() libOrch.WorkerClients {
	workerClients := libOrch.WorkerClients{
		Clients: make(map[string]*libOrch.NamedClient),
	}

	currentIndex := 0
	var unifiedClient *redis.Client

	for {
		clientOptions := tryGetClientOptions("UNIFIED", currentIndex)

		if clientOptions == nil {
			if currentIndex == 0 {
				panic("At least one worker client must be defined")
			}

			break
		}

		if currentIndex == 0 {
			unifiedClient = redis.NewClient(clientOptions.redisOptions)

			// Ensure that the client is connected
			if err := unifiedClient.Ping(context.Background()).Err(); err != nil {
				panic(err)
			}
		}

		namedClient := &libOrch.NamedClient{
			Name:   clientOptions.displayName,
			Client: unifiedClient,
		}

		if currentIndex == 0 {
			workerClients.DefaultClient = namedClient
		}

		workerClients.Clients[namedClient.Name] = namedClient

		currentIndex++
	}

	return workerClients
}

func tryGetClientOptions(identifier string, currentIndex int) *clientOptions {
	host := lib.GetEnvVariableRaw(fmt.Sprintf("WORKER_%s_HOST", identifier), "NONE", true)
	port := lib.GetEnvVariableRaw(fmt.Sprintf("WORKER_%s_PORT", identifier), "NONE", true)
	password := lib.GetEnvVariableRaw(fmt.Sprintf("WORKER_%s_PASSWORD", identifier), "NONE", true)
	displayName := lib.GetEnvVariableRaw(fmt.Sprintf("WORKER_%d_DISPLAY_NAME", currentIndex), "NONE", true)

	if host == "NONE" || port == "NONE" || password == "NONE" || displayName == "NONE" {
		return nil
	}

	options := &redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password,
		DB:       0,
	}

	isSecure := lib.GetEnvVariableBool(fmt.Sprintf("WORKER_%s_IS_SECURE", identifier), false)

	if isSecure {
		clientCert := lib.GetHexEnvVariable(fmt.Sprintf("WORKER_%s_CERT_HEX", identifier), "")
		clientKey := lib.GetHexEnvVariable(fmt.Sprintf("WORKER_%s_KEY_HEX", identifier), "")

		cert, err := tls.X509KeyPair(clientCert, clientKey)
		if err != nil {
			panic(fmt.Errorf("error loading orchestrator cert: %s", err))
		}

		// Load CA cert
		caCertPool := x509.NewCertPool()
		caCert := lib.GetHexEnvVariable(fmt.Sprintf("WORKER_%s_CA_CERT_HEX", identifier), "")
		ok := caCertPool.AppendCertsFromPEM(caCert)
		if !ok {
			panic("failed to parse root certificate")
		}

		options.TLSConfig = &tls.Config{
			MinVersion:         tls.VersionTLS12,
			InsecureSkipVerify: lib.GetEnvVariableBool(fmt.Sprintf("WORKER_%s_INSECURE_SKIP_VERIFY", identifier), false),
			Certificates:       []tls.Certificate{cert},
			RootCAs:            caCertPool,
		}
	}

	return &clientOptions{
		redisOptions: options,
		displayName:  displayName,
	}
}

// Gets redis client for when running as an agent
func getAgentClient() libOrch.WorkerClients {
	workerClients := libOrch.WorkerClients{
		Clients: make(map[string]*libOrch.NamedClient),
	}

	workerClients.Clients[agent.AgentWorkerName] = &libOrch.NamedClient{
		Name: agent.AgentWorkerName,
		Client: redis.NewClient(&redis.Options{
			Addr: fmt.Sprintf("%s:%s", agent.WorkerRedisHost, agent.WorkerRedisPort),
		}),
	}
	workerClients.DefaultClient = workerClients.Clients[agent.AgentWorkerName]

	return workerClients
}

func getOrchestratorClient(standalone bool) *redis.Client {
	if !standalone {
		return redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", agent.OrchestratorRedisHost, agent.OrchestratorRedisPort),
			Username: "default",
			Password: "",
		},
		)
	}

	orchestratorHost := lib.GetEnvVariable("ORCHESTRATOR_REDIS_HOST", "localhost")

	options := &redis.Options{
		Addr:     fmt.Sprintf("%s:%s", orchestratorHost, lib.GetEnvVariable("ORCHESTRATOR_REDIS_PORT", "10000")),
		Username: "default",
		Password: lib.GetEnvVariable("ORCHESTRATOR_REDIS_PASSWORD", ""),
	}

	isSecure := lib.GetEnvVariableBool("ORCHESTRATOR_REDIS_IS_SECURE", false)

	if isSecure {
		clientCert := lib.GetHexEnvVariable("ORCHESTRATOR_REDIS_CERT_HEX", "")
		clientKey := lib.GetHexEnvVariable("ORCHESTRATOR_REDIS_KEY_HEX", "")

		cert, err := tls.X509KeyPair([]byte(clientCert), []byte(clientKey))
		if err != nil {
			panic(fmt.Errorf("error loading orchestrator cert: %s", err))
		}

		// Load CA cert
		caCertPool := x509.NewCertPool()
		caCert := lib.GetHexEnvVariable("ORCHESTRATOR_REDIS_CA_CERT_HEX", "")
		ok := caCertPool.AppendCertsFromPEM(caCert)
		if !ok {
			panic("failed to parse root certificate")
		}

		options.TLSConfig = &tls.Config{
			MinVersion:         tls.VersionTLS12,
			InsecureSkipVerify: lib.GetEnvVariableBool("ORCHESTRATOR_REDIS_INSECURE_SKIP_VERIFY", false),
			Certificates:       []tls.Certificate{cert},
			RootCAs:            caCertPool,
		}
	}

	client := redis.NewClient(options)

	// Ensure that the client is connected
	if err := client.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}

	return client
}

func getStoreMongoDB(ctx context.Context, standalone bool) *mongo.Database {
	if !standalone {
		return nil
	}

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
