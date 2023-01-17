package function_auth_client

import (
	"context"
	"fmt"
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
	funcUrlOverride    []string
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
		funcUrlOverride:    getFunctionUrlOverrides(),
	}

	functionAuthClient.startAutoRefreshLiveFunctions()

	return functionAuthClient
}

func getFunctionUrlOverrides() []string {
	// Loop through counter FUNCTION_URL_OVERRIDE_1, FUNCTION_URL_OVERRIDE_2, etc

	// If the env variable is not set, return nil

	// If the env variable is set, return the value

	values := []string{}
	index := 0

	for {
		value := lib.GetEnvVariableRaw(fmt.Sprintf("FUNCTION_URL_OVERRIDE_%d", index), "NONE_LEFT", true)

		if value == "NONE_LEFT" {
			if index == 0 {
				values = append(values, "NO")
			}

			break
		}

		values = append(values, value)
		index++
	}

	fmt.Printf("Function URL overrides: %v\n", values)

	return values
}
