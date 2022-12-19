package function_auth_client

import (
	"context"
	"sync"

	functions "cloud.google.com/go/functions/apiv2"
	"github.com/APITeamLimited/globe-test/lib"
	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"google.golang.org/api/option"
)

type FunctionAuthClient struct {
	functionClient     *functions.FunctionClient
	liveFunctions      []libOrch.LiveFunction
	liveFunctionsMutex sync.Mutex
	ctx                context.Context
	serviceAccount     []byte
	funcUrlOverride    string
}

var _ = libOrch.FunctionAuthClient(&FunctionAuthClient{})

func CreateFunctionAuthClient(ctx context.Context, funcMode bool) *FunctionAuthClient {
	if !funcMode {
		return nil
	}

	// Convert hex to bytes
	serviceAccount := lib.GetHexEnvVariable("SERVICE_ACCOUNT_KEY_HEX", "")

	functionClient, err := functions.NewFunctionClient(ctx, option.WithCredentialsJSON(serviceAccount))
	if err != nil {
		panic(err)
	}

	functionAuthClient := &FunctionAuthClient{
		functionClient:     functionClient,
		ctx:                ctx,
		liveFunctions:      []libOrch.LiveFunction{},
		liveFunctionsMutex: sync.Mutex{},
		serviceAccount:     serviceAccount,
		funcUrlOverride:    lib.GetEnvVariable("FUNCTION_URL_OVERRIDE", "NO"),
	}

	functionAuthClient.startAutoRefreshLiveFunctions()

	return functionAuthClient
}
