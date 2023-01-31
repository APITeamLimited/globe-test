package run_auth_client

import (
	"fmt"
	"strings"
	"time"

	"cloud.google.com/go/run/apiv2/runpb"
	"github.com/APITeamLimited/globe-test/lib"
	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
)

func (config *RunAuthClient) Services() []libOrch.LiveService {
	config.liveServicesMutex.Lock()
	defer config.liveServicesMutex.Unlock()
	return config.liveServices
}

func (config *RunAuthClient) startAutoRefreshLiveServices() {
	if config.serviceClient == nil {
		return
	}

	// Get a new token straight away
	config.updateLiveServices()

	go func() {
		for {
			time.Sleep(time.Second * 10)
			config.updateLiveServices()
		}
	}()
}

// Queries the google cloud functions API to get the list of function URLs
func (config *RunAuthClient) updateLiveServices() {
	findServices := func(location string) {
		servicesIterator := config.serviceClient.ListServices(config.ctx, &runpb.ListServicesRequest{
			Parent: fmt.Sprintf("projects/apiteam-%s/locations/%s", lib.GetEnvVariable("ENVIRONMENT", ""), location),
		})

		for {
			service, err := servicesIterator.Next()
			if err != nil {
				if !strings.Contains(err.Error(), "no more items") {
					fmt.Println(err)
				}

				break
			}

			if service == nil {
				break
			}

			location, err := parseLocation(service.Description)
			if err != nil {
				continue
			}

			config.liveServicesMutex.Lock()
			// Check if the service is already in the list and update it if it is remove it and add it again

			isNewService := true

			for i, liveService := range config.liveServices {
				if liveService.Location == location {
					liveService.Uri = service.Uri
					liveService.State = service.TerminalCondition.State
					config.liveServices[i] = liveService
					isNewService = false
					break
				}
			}

			if isNewService {
				config.liveServices = append(config.liveServices, libOrch.LiveService{
					Location: location,
					Uri:      service.Uri,
					State:    service.TerminalCondition.State,
				})
			}

			config.liveServicesMutex.Unlock()
		}
	}

	for _, location := range config.loadZones {
		go findServices(location)
	}
}

func parseLocation(description string) (string, error) {
	// Format location:{location};

	if description == "" {
		return "", fmt.Errorf("invalid function description: %s", description)
	}

	// Split at semicolon
	parts := strings.Split(description, ";")

	if len(parts) == 0 {
		return "", fmt.Errorf("invalid function description: %s", description)
	}

	locationParts := strings.Split(parts[0], ":")

	return locationParts[1], nil
}
