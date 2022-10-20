package orchMetrics

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/APITeamLimited/globe-test/worker/output"
	"github.com/APITeamLimited/globe-test/worker/output/globetest"
)

// Wrapped envelope allows for grouping
type wrappedEnvelope struct {
	globetest.SampleEnvelope `json:"sampleEnvelope"`
	workerId                 string
	location                 string
	weighting                int
}

// Cached metrics are stored before being collated and sent
type cachedMetricsStore struct {
	// Each envelope in the map is a certain metric
	envelopes map[string][]*wrappedEnvelope
	mu        sync.RWMutex
	flusher   *output.PeriodicFlusher
	gs        libOrch.BaseGlobalState
}

var (
	_ libOrch.BaseMetricsStore = &cachedMetricsStore{}
)

const globalName = "global"

func NewCachedMetricsStore(gs libOrch.BaseGlobalState) *cachedMetricsStore {
	store := &cachedMetricsStore{
		gs:        gs,
		envelopes: make(map[string][]*wrappedEnvelope),
		mu:        sync.RWMutex{},
	}

	// This will never return an error
	pf, _ := output.NewPeriodicFlusher(200*time.Millisecond, store.FlushMetrics)

	store.flusher = pf

	return store
}

func (store *cachedMetricsStore) AddMessage(message libOrch.WorkerMessage, workerLocation string) error {
	store.mu.Lock()
	defer store.mu.Unlock()

	var sampleEnvelopes []globetest.SampleEnvelope

	err := json.Unmarshal([]byte(message.Message), &sampleEnvelopes)
	if err != nil {
		return err
	}

	for _, sampleEnvelope := range sampleEnvelopes {
		metricName := sampleEnvelope.Metric.Name

		// If the metric name is not in the map, create a new slice
		if _, ok := store.envelopes[metricName]; !ok {
			store.envelopes[metricName] = make([]*wrappedEnvelope, 0)
		}

		store.envelopes[metricName] = append(store.envelopes[metricName], &wrappedEnvelope{
			SampleEnvelope: sampleEnvelope,
			workerId:       message.WorkerId,
			location:       workerLocation,
		})
	}

	return nil
}

// Empty the store and returns its contents
func (store *cachedMetricsStore) emptyStore() map[string][]*wrappedEnvelope {
	store.mu.Lock()
	defer store.mu.Unlock()

	// Copy the map
	result := make(map[string][]*wrappedEnvelope, len(store.envelopes))
	for metricName, sampleEnvelopes := range store.envelopes {
		result[metricName] = make([]*wrappedEnvelope, len(sampleEnvelopes))
		copy(result[metricName], sampleEnvelopes)
	}

	// Empty the map
	store.envelopes = make(map[string][]*wrappedEnvelope)

	return result
}

func (store *cachedMetricsStore) FlushMetrics() {
	envelopes, err := store.getMetrics()

	if err != nil {
		fmt.Println(err)
		//libOrch.HandleError(store.ctx, store.client, store.jobId, store.orchestratorId, err)
		return
	}

	if len(envelopes) == 0 {
		return
	}

	// Marshall the envelopes
	marshalledEnvelopes, err := json.Marshal(envelopes)
	if err != nil {
		fmt.Println("FlushMetrics", err)
		libOrch.HandleError(store.gs, err)
		return
	}

	libOrch.DispatchMessage(store.gs, string(marshalledEnvelopes), "METRICS")
}

func (store *cachedMetricsStore) getMetrics() (map[string][]*wrappedEnvelope, error) {
	envelopeMap := store.emptyStore()
	currentTime := time.Now()

	// Combined metrics is the collated metrics
	combinedMetrics := make(map[string][]*wrappedEnvelope)

	// Create global name entry
	combinedMetrics[globalName] = make([]*wrappedEnvelope, 0)

	for _, envelopes := range envelopeMap {
		for _, envelope := range envelopes {
			locationName := envelope.location

			// If locationName is not in the map, create a new entry
			if _, ok := combinedMetrics[locationName]; !ok {
				combinedMetrics[locationName] = make([]*wrappedEnvelope, 0)
			}
		}
	}

	// sortedEnvelopeMap is the raw envelopes, indexed by location
	sortedEnvelopeMap := make(map[string][]*wrappedEnvelope)

	sortedEnvelopeMap[globalName] = make([]*wrappedEnvelope, 0)

	for _, envelopes := range envelopeMap {
		// Sort the envelopes by location
		for _, envelope := range envelopes {
			locationName := envelope.location

			// If locationName is not in the map, create a new entry
			if _, ok := sortedEnvelopeMap[locationName]; !ok {
				sortedEnvelopeMap[locationName] = make([]*wrappedEnvelope, 0)
			}

			sortedEnvelopeMap[locationName] = append(sortedEnvelopeMap[locationName], envelope)

			// Add to global
			sortedEnvelopeMap[globalName] = append(sortedEnvelopeMap[globalName], envelope)
		}
	}

	// Loop through the messages again and call the accumulator function for the
	// correct type

	// Store found vu count for that worker at that time
	workerVuCount := make(map[string]int)

	for location, envelopes := range sortedEnvelopeMap {
		for _, envelope := range envelopes {
			if envelope.Type != "Point" {
				continue
			}

			// Find worker VU count
			vuCount, err := findVuCount(&workerVuCount, envelopeMap, envelope.workerId)
			if err != nil {
				vuCount = -1 // Set to -1 to indicate that the vu count was not found
				//return nil, err
			}

			// Determine based off metric type what accumulator function to use
			metricName := envelope.Metric.Name

			if metricName == "http_reqs" {
				combinedMetrics[location] = calculateTotal(combinedMetrics[location], envelope, currentTime)
				continue
			}

			if metricName == "vus" {
				combinedMetrics[location] = calculateTotal(combinedMetrics[location], envelope, currentTime)
				continue
			}

			if metricName == "vus_max" {
				combinedMetrics[location] = calculateTotal(combinedMetrics[location], envelope, currentTime)
				continue
			}

			// Think trends are just averages here

			if metricName == "http_req_duration" && vuCount != -1 {
				combinedMetrics[location] = calculateMean(combinedMetrics[location], envelope, currentTime, vuCount)
				continue
			}

			if metricName == "http_req_blocked" && vuCount != -1 {
				combinedMetrics[location] = calculateMean(combinedMetrics[location], envelope, currentTime, vuCount)
				continue
			}

			if metricName == "http_req_connecting" && vuCount != -1 {
				combinedMetrics[location] = calculateMean(combinedMetrics[location], envelope, currentTime, vuCount)
				continue
			}

			if metricName == "http_req_sending" && vuCount != -1 {
				combinedMetrics[location] = calculateMean(combinedMetrics[location], envelope, currentTime, vuCount)
				continue
			}

			if metricName == "http_req_waiting" && vuCount != -1 {
				combinedMetrics[location] = calculateMean(combinedMetrics[location], envelope, currentTime, vuCount)
				continue
			}

			if metricName == "http_req_receiving" && vuCount != -1 {
				combinedMetrics[location] = calculateMean(combinedMetrics[location], envelope, currentTime, vuCount)
				continue
			}

			if metricName == "http_req_tls_handshaking" && vuCount != -1 {
				combinedMetrics[location] = calculateMean(combinedMetrics[location], envelope, currentTime, vuCount)
				continue
			}

			if metricName == "http_req_failed" && vuCount != -1 {
				combinedMetrics[location] = calculateMean(combinedMetrics[location], envelope, currentTime, vuCount)
				continue
			}
		}
	}

	// If any keys are empty, remove them
	for key, value := range combinedMetrics {
		if len(value) == 0 {
			delete(combinedMetrics, key)
		}
	}

	return combinedMetrics, nil
}

func findVuCount(workerVuCount *map[string]int, envelopeMap map[string][]*wrappedEnvelope, workerId string) (int, error) {
	// Find worker Vu count in workerVuCount map first
	if vuCount, ok := (*workerVuCount)[workerId]; ok {
		return vuCount, nil
	}

	// Check if vus is keyed in the envelope map
	if _, ok := (envelopeMap)["vus"]; !ok {
		return -1, errors.New("vu count not found")
	}

	// Find vu count for this workerId
	for _, envelope := range envelopeMap["vus"] {
		if envelope.workerId == workerId {
			newCount := int((*envelope).SampleEnvelope.Data.Value)
			(*workerVuCount)[workerId] = newCount
			return newCount, nil
		}
	}

	return -1, errors.New("vu count not found")
}

func (store *cachedMetricsStore) Stop() {
	store.flusher.Stop()
}
