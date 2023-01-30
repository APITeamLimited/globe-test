package run_auth_client

import (
	"fmt"
	"net/http"

	"cloud.google.com/go/run/apiv2/runpb"
	"github.com/APITeamLimited/globe-test/lib"
	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/gorilla/websocket"
	"google.golang.org/api/idtoken"
	"google.golang.org/api/option"
)

func (config *RunAuthClient) ExecuteService(location string) (*websocket.Conn, error) {
	config.liveServicesMutex.Lock()
	defer config.liveServicesMutex.Unlock()

	var liveFunction *libOrch.LiveService

	// Find the function
	for _, function := range config.liveServices {
		if function.Location == location {
			liveFunction = &function
			break
		}
	}

	specificFuncOverride := getSpecificFuncOverride(config.serviceUrlOverride, location)

	if liveFunction == nil && specificFuncOverride == "NO" {
		return nil, fmt.Errorf("function at location %s not found", location)
	}

	uri := determineUri(liveFunction, specificFuncOverride)

	tokenSource, err := idtoken.NewTokenSource(config.ctx, uri, option.WithCredentialsJSON(config.serviceAccount))
	if err != nil {
		return nil, err
	}

	// Get the token
	token, err := tokenSource.Token()
	if err != nil {
		return nil, err
	}

	authHeader := http.Header{
		"Authorization": []string{fmt.Sprintf("Bearer %s", token.AccessToken)},
	}

	conn, response, err := websocket.DefaultDialer.Dial(uri, authHeader)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 101 {
		return nil, fmt.Errorf("websocket connection failed with status code %d", response.StatusCode)
	}

	return conn, nil
}

func getSpecificFuncOverride(overrides []string, location string) string {
	if len(overrides) == 1 && overrides[0] == "NO" {
		return "NO"
	}

	// Loop through WORKER_0_DISPLAY_NAME, WORKER_1_DISPLAY_NAME env variables till get location
	index := 0

	for {
		locationIndex := lib.GetEnvVariable(fmt.Sprintf("WORKER_%d_DISPLAY_NAME", index), "")

		if locationIndex == "" {
			break
		}

		if locationIndex == location {
			return overrides[index]
		}

		index++
	}

	panic(fmt.Sprintf("location %s not found in env variables", location))
}

func (config *RunAuthClient) CheckServiceAvailability(location string) error {
	if config.serviceClient == nil {
		if getSpecificFuncOverride(config.serviceUrlOverride, location) != "NO" {
			return nil
		}

		return fmt.Errorf("service client not initialised")
	}

	config.liveServicesMutex.Lock()
	defer config.liveServicesMutex.Unlock()

	var liveService *libOrch.LiveService

	// Find the function
	for _, service := range config.liveServices {
		if service.Location == location {
			liveService = &service
			break
		}
	}

	// If running locally, allow all functions
	if getSpecificFuncOverride(config.serviceUrlOverride, location) != "NO" {
		return nil
	}

	if liveService == nil {
		return fmt.Errorf("location %s not found", location)
	}

	state := (*liveService).State

	if state != runpb.Condition_CONDITION_SUCCEEDED && state != runpb.Condition_STATE_UNSPECIFIED {
		return fmt.Errorf("location %s is not active", location)
	}

	return nil
}

func determineUri(liveFunction *libOrch.LiveService, serviceUrlOverride string) string {
	if serviceUrlOverride == "NO" {
		return replaceSchemeWithWs(liveFunction.Uri)
	}

	fmt.Println("Using function url override", serviceUrlOverride)
	return replaceSchemeWithWs(serviceUrlOverride)
}

func replaceSchemeWithWs(uri string) string {
	if uri[:5] == "https" {
		return "wss" + uri[5:]
	}

	return "ws" + uri[4:]
}
