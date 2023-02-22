package orchestrator

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"

	"github.com/APITeamLimited/globe-test/lib"
	"github.com/APITeamLimited/globe-test/lib/agent"
	"github.com/APITeamLimited/redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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
