package run_auth_client

import (
	"fmt"
	"net/http"
	"time"

	"cloud.google.com/go/run/apiv2/runpb"
	"github.com/APITeamLimited/globe-test/lib"
	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/gorilla/websocket"
	dnscache "go.mercari.io/go-dnscache"
	"google.golang.org/api/idtoken"
	"google.golang.org/api/option"
)

func (config *RunAuthClient) ExecuteService(gs libOrch.BaseGlobalState, location string) (*websocket.Conn, error) {
	config.liveServicesMutex.Lock()

	var liveService libOrch.LiveService
	foundLiveService := false

	// Find the service
	for _, service := range config.liveServices {
		if service.Location == location {
			liveService = service
			foundLiveService = true
			break
		}
	}

	config.liveServicesMutex.Unlock()

	specificServiceOverride := getSpecificServiceOverride(config.serviceUrlOverride, location)

	if !foundLiveService && specificServiceOverride == "NO" {
		return nil, fmt.Errorf("service at location %s not found", location)
	}

	uri := determineUri(liveService, specificServiceOverride)

	tokenSource, err := idtoken.NewTokenSource(config.ctx, determineAudience(liveService, specificServiceOverride), option.WithCredentialsJSON(config.serviceAccount))
	if err != nil {
		fmt.Println("err", err)
		if specificServiceOverride == "NO" {
			return nil, err
		}
		tokenSource = nil
	}

	headers := http.Header{}

	if tokenSource != nil {
		// Get the token
		token, err := tokenSource.Token()
		if err != nil {
			return nil, err
		}

		headers.Add("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
	}

	// Create websocket dialer
	dialer := &websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: time.Minute,
		NetDialContext:   dnscache.DialFunc(config.resolver, nil),
	}

	conn, response, err := dialer.Dial(uri, headers)
	if err != nil {
		fmt.Println("failed to connect to websocket", err)
		return nil, err
	}

	if response.StatusCode != 101 {
		return nil, fmt.Errorf("websocket connection failed with status code %d %s", response.StatusCode, response.Status)
	}

	return conn, nil
}

func getSpecificServiceOverride(overrides []string, location string) string {
	if len(overrides) == 1 && overrides[0] == "NO" {
		return "NO"
	}

	// Loop through env variables till get location eg WORKER_0_DISPLAY_NAME, WORKER_1_DISPLAY_NAME
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
		if getSpecificServiceOverride(config.serviceUrlOverride, location) != "NO" {
			return nil
		}

		return fmt.Errorf("service client not initialised")
	}

	config.liveServicesMutex.Lock()
	defer config.liveServicesMutex.Unlock()

	var liveService *libOrch.LiveService

	// Find the service
	for _, service := range config.liveServices {
		if service.Location == location {
			liveService = &service
			break
		}
	}

	// If running locally, allow all functions
	if getSpecificServiceOverride(config.serviceUrlOverride, location) != "NO" {
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

func determineAudience(liveService libOrch.LiveService, serviceUrlOverride string) string {
	if serviceUrlOverride == "NO" {
		fmt.Println("Using service url", liveService.Uri)
		return liveService.Uri
	}

	return serviceUrlOverride
}

func determineUri(liveService libOrch.LiveService, serviceUrlOverride string) string {
	if serviceUrlOverride == "NO" {
		return replaceSchemeWithWs(liveService.Uri)
	}

	fmt.Println("Using service url override", serviceUrlOverride)
	return replaceSchemeWithWs(serviceUrlOverride)
}

func replaceSchemeWithWs(uri string) string {
	if uri[:5] == "https" {
		return "wss" + uri[5:]
	}

	return "ws" + uri[4:]
}
