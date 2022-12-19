package function_auth_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"cloud.google.com/go/functions/apiv2/functionspb"
	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"google.golang.org/api/idtoken"
)

func (config *FunctionAuthClient) ExecuteFunction(location string, childJob libOrch.ChildJob) (*(chan libOrch.FunctionResult), error) {
	config.liveFunctionsMutex.Lock()
	defer config.liveFunctionsMutex.Unlock()

	var liveFunction *libOrch.LiveFunction

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

	// Encode the job in the body of the request

	// Authenticate the function
	client, err := idtoken.NewClient(config.ctx, liveFunction.Uri)
	if err != nil {
		return nil, err
	}

	// Send job as body of request
	var buff bytes.Buffer
	err = json.NewEncoder(&buff).Encode(childJob)
	if err != nil {
		return nil, err
	}

	request := &http.Request{
		Method:     "POST",
		RequestURI: liveFunction.Uri,
		Body:       io.NopCloser(&buff),
	}

	var responseCh chan libOrch.FunctionResult

	go func() {
		// Make the request
		response, err := client.Do(request)

		fmt.Println("response", response)
		fmt.Println("err", err)

		responseCh <- libOrch.FunctionResult{
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

<<<<<<< HEAD
	var liveFunction *libOrch.LiveFunction
=======
	var liveFunction *lib.LiveFunction
>>>>>>> 91df2ad531f51a3494f2b358576d5bfedca06952

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
