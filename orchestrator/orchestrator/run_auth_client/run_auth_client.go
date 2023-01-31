package run_auth_client

import (
	"context"
	"fmt"
	"sync"

	run "cloud.google.com/go/run/apiv2"
	"github.com/APITeamLimited/globe-test/lib"
	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"google.golang.org/api/option"
)

type RunAuthClient struct {
	serviceClient      *run.ServicesClient
	liveServices       []libOrch.LiveService
	liveServicesMutex  sync.Mutex
	ctx                context.Context
	serviceAccount     []byte
	serviceUrlOverride []string
	loadZones          []string
}

var _ = libOrch.RunAuthClient(&RunAuthClient{})

func CreateServicesClient(ctx context.Context, standalone bool, loadZones []string) *RunAuthClient {
	var serviceClient *run.ServicesClient

	// Convert hex to bytes
	serviceAccount := lib.GetHexEnvVariable("SERVICE_ACCOUNT_KEY_HEX", "NONE")
	if string(serviceAccount) != "NONE" {
		loadedServiceClient, err := run.NewServicesClient(ctx, option.WithCredentialsJSON(serviceAccount))
		if err != nil {
			panic(err)
		}

		serviceClient = loadedServiceClient
	} else if !standalone {
		fmt.Println("Service account key not found, assuming using overrides")
		serviceClient = nil
	} else {
		panic("Service account key not found")
	}

	runAuthClient := &RunAuthClient{
		serviceClient:      serviceClient,
		ctx:                ctx,
		liveServices:       []libOrch.LiveService{},
		liveServicesMutex:  sync.Mutex{},
		serviceAccount:     serviceAccount,
		serviceUrlOverride: getServiceUrlOverrides(),
		loadZones:          loadZones,
	}

	runAuthClient.startAutoRefreshLiveServices()

	return runAuthClient
}

func getServiceUrlOverrides() []string {
	// Loop through counter SERVICE_URL_OVERRIDE_1, SERVICE_URL_OVERRIDE_2, etc

	// If the env variable is not set, return nil

	// If the env variable is set, return the value

	values := []string{}
	index := 0

	for {
		value := lib.GetEnvVariableRaw(fmt.Sprintf("SERVICE_URL_OVERRIDE_%d", index), "NONE_LEFT", true)

		if value == "NONE_LEFT" {
			if index == 0 {
				values = append(values, "NO")
			}

			break
		}

		values = append(values, value)
		index++
	}

	return values
}
