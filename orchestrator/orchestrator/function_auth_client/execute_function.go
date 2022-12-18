package function_auth_client

import (
	"fmt"

	"cloud.google.com/go/functions/apiv2/functionspb"
	"github.com/APITeamLimited/globe-test/lib"
	"google.golang.org/api/idtoken"
)

func (config *FunctionAuthClient) ExecuteFunction(location string) (*(chan lib.FunctionResult), error) {
	config.liveFunctionsMutex.Lock()
	defer config.liveFunctionsMutex.Unlock()

	var liveFunction *lib.LiveFunction

	// Find the function
	for _, function := range config.liveFunctions {
		if function.Location == location {
			liveFunction = &function
			break
		}
	}

	if liveFunction == nil {
		return nil, fmt.Errorf("function at location %s not found", location)
	}

	// Authenticate the function
	client, err := idtoken.NewClient(config.ctx, liveFunction.Uri)
	if err != nil {
		return nil, err
	}

	var responseCh chan lib.FunctionResult

	go func() {
		// Make the request
		response, err := client.Get(liveFunction.Uri)

		responseCh <- lib.FunctionResult{
			Response: response,
			Error:    err,
		}

		close(responseCh)
	}()

	return &responseCh, nil
}

func (config *FunctionAuthClient) CheckFunctionAvailability(location string) error {
	config.liveFunctionsMutex.Lock()
	defer config.liveFunctionsMutex.Unlock()

	var liveFunction *lib.LiveFunction

	// Find the function
	for _, function := range config.liveFunctions {
		if function.Location == location {
			liveFunction = &function
			break
		}
	}

	if liveFunction == nil {
		return fmt.Errorf("location %s not found", location)
	}

	if (*liveFunction).State != functionspb.Function_ACTIVE {
		return fmt.Errorf("location %s is not active", location)
	}

	return nil
}
