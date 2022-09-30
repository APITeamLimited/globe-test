package orchMetrics

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/APITeamLimited/globe-test/worker/output"
	"github.com/APITeamLimited/globe-test/worker/output/globetest"
	"github.com/APITeamLimited/redis/v9"
)

// Wrapped envelope allows for grouping
type wrappedEnvelope struct ***REMOVED***
	globetest.SampleEnvelope `json:"sampleEnvelope"`
	workerId                 string
	location                 string
	weighting                int
***REMOVED***

// Cached metrics are stored before being collated and sent
type cachedMetricsStore struct ***REMOVED***
	// Each envelope in the map is a certain metric
	envelopes      map[string][]*wrappedEnvelope
	mu             sync.RWMutex
	flusher        *output.PeriodicFlusher
	ctx            context.Context
	client         *redis.Client
	orchestratorId string
	jobId          string
***REMOVED***

var (
	_ libOrch.BaseMetricsStore = &cachedMetricsStore***REMOVED******REMOVED***
)

const globalName = "global"

func NewCachedMetricsStore(ctx context.Context, client *redis.Client, orchestratorId string, jobId string) *cachedMetricsStore ***REMOVED***
	store := &cachedMetricsStore***REMOVED***
		ctx:            ctx,
		client:         client,
		orchestratorId: orchestratorId,
		jobId:          jobId,
		envelopes:      make(map[string][]*wrappedEnvelope),
		// Mutex doesn't need to be initialised
	***REMOVED***

	// This will never return an error
	pf, _ := output.NewPeriodicFlusher(200*time.Millisecond, store.FlushMetrics)

	store.flusher = pf

	return store
***REMOVED***

func (store *cachedMetricsStore) AddMessage(message libOrch.WorkerMessage, workerLocation string) error ***REMOVED***
	store.mu.Lock()
	defer store.mu.Unlock()

	var sampleEnvelopes []globetest.SampleEnvelope

	err := json.Unmarshal([]byte(message.Message), &sampleEnvelopes)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	for _, sampleEnvelope := range sampleEnvelopes ***REMOVED***
		metricName := sampleEnvelope.Metric.Name

		// If the metric name is not in the map, create a new slice
		if _, ok := store.envelopes[metricName]; !ok ***REMOVED***
			store.envelopes[metricName] = make([]*wrappedEnvelope, 0)
		***REMOVED***

		store.envelopes[metricName] = append(store.envelopes[metricName], &wrappedEnvelope***REMOVED***
			SampleEnvelope: sampleEnvelope,
			workerId:       message.WorkerId,
			location:       workerLocation,
		***REMOVED***)
	***REMOVED***

	return nil
***REMOVED***

// Empty the store and returns its contents
func (store *cachedMetricsStore) emptyStore() map[string][]*wrappedEnvelope ***REMOVED***
	store.mu.Lock()
	defer store.mu.Unlock()

	// Copy the map
	result := make(map[string][]*wrappedEnvelope, len(store.envelopes))
	for metricName, sampleEnvelopes := range store.envelopes ***REMOVED***
		result[metricName] = make([]*wrappedEnvelope, len(sampleEnvelopes))
		copy(result[metricName], sampleEnvelopes)
	***REMOVED***

	// Empty the map
	store.envelopes = make(map[string][]*wrappedEnvelope)

	return result
***REMOVED***

func (store *cachedMetricsStore) FlushMetrics() ***REMOVED***
	envelopes, err := store.getMetrics()

	if err != nil ***REMOVED***
		fmt.Println(err)
		//libOrch.HandleError(store.ctx, store.client, store.jobId, store.orchestratorId, err)
		return
	***REMOVED***

	if len(envelopes) == 0 ***REMOVED***
		return
	***REMOVED***

	// Marshall the envelopes
	marshalledEnvelopes, err := json.Marshal(envelopes)
	if err != nil ***REMOVED***
		libOrch.HandleError(store.ctx, store.client, store.jobId, store.orchestratorId, err)
		return
	***REMOVED***

	libOrch.DispatchMessage(store.ctx, store.client, store.jobId, store.orchestratorId, string(marshalledEnvelopes), "METRICS")
***REMOVED***

func (store *cachedMetricsStore) getMetrics() (map[string][]*wrappedEnvelope, error) ***REMOVED***
	envelopeMap := store.emptyStore()
	currentTime := time.Now()

	// Combined metrics is indexed by location NOT metric name
	combinedMetrics := make(map[string][]*wrappedEnvelope)

	// Create names for the groups
	combinedMetrics[globalName] = make([]*wrappedEnvelope, 0)

	for _, envelopes := range envelopeMap ***REMOVED***
		for _, envelope := range envelopes ***REMOVED***
			locationName := envelope.location

			// If locationName is not in the map, create a new entry
			if _, ok := combinedMetrics[locationName]; !ok ***REMOVED***
				combinedMetrics[locationName] = make([]*wrappedEnvelope, 0)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// Sorted envelope map, idnex by location
	sortedEnvelopeMap := make(map[string][]*wrappedEnvelope)

	sortedEnvelopeMap[globalName] = make([]*wrappedEnvelope, 0)

	for _, envelopes := range envelopeMap ***REMOVED***
		// Sort the envelopes by location
		for _, envelope := range envelopes ***REMOVED***
			locationName := envelope.location

			// If locationName is not in the map, create a new entry
			if _, ok := sortedEnvelopeMap[locationName]; !ok ***REMOVED***
				sortedEnvelopeMap[locationName] = make([]*wrappedEnvelope, 0)
			***REMOVED***

			sortedEnvelopeMap[locationName] = append(sortedEnvelopeMap[locationName], envelope)

			// Add to global
			sortedEnvelopeMap[globalName] = append(sortedEnvelopeMap[globalName], envelope)
		***REMOVED***
	***REMOVED***

	// Loop through the messages again and call the accumulator function for the
	// correct type

	// Store found vu count for that worker at that time
	workerVuCount := make(map[string]int)

	for location, envelopes := range sortedEnvelopeMap ***REMOVED***
		for _, envelope := range envelopes ***REMOVED***
			if envelope.Type != "Point" ***REMOVED***
				continue
			***REMOVED***

			// Find worker Vu count
			vuCount, err := findVuCount(&workerVuCount, envelopeMap, envelope.workerId)
			if err != nil ***REMOVED***
				vuCount = -1 // Set to -1 to indicate that the vu count was not found
				//return nil, err
			***REMOVED***

			// Determine based off metric type what accumulator function to use
			metricName := envelope.Metric.Name

			if metricName == "http_reqs" ***REMOVED***
				combinedMetrics[location] = calculateTotal(combinedMetrics[location], envelope, currentTime)
				continue
			***REMOVED***

			if metricName == "vus" ***REMOVED***
				combinedMetrics[location] = calculateTotal(combinedMetrics[location], envelope, currentTime)
				continue
			***REMOVED***

			if metricName == "vus_max" ***REMOVED***
				combinedMetrics[location] = calculateTotal(combinedMetrics[location], envelope, currentTime)
				continue
			***REMOVED***

			// Think trends are just averages here

			if metricName == "http_req_duration" && vuCount != -1 ***REMOVED***
				combinedMetrics[location] = calculateMean(combinedMetrics[location], envelope, currentTime, vuCount)
				continue
			***REMOVED***

			if metricName == "http_req_blocked" && vuCount != -1 ***REMOVED***
				combinedMetrics[location] = calculateMean(combinedMetrics[location], envelope, currentTime, vuCount)
				continue
			***REMOVED***

			if metricName == "http_req_connecting" && vuCount != -1 ***REMOVED***
				combinedMetrics[location] = calculateMean(combinedMetrics[location], envelope, currentTime, vuCount)
				continue
			***REMOVED***

			if metricName == "http_req_sending" && vuCount != -1 ***REMOVED***
				combinedMetrics[location] = calculateMean(combinedMetrics[location], envelope, currentTime, vuCount)
				continue
			***REMOVED***

			if metricName == "http_req_waiting" && vuCount != -1 ***REMOVED***
				combinedMetrics[location] = calculateMean(combinedMetrics[location], envelope, currentTime, vuCount)
				continue
			***REMOVED***

			if metricName == "http_req_receiving" && vuCount != -1 ***REMOVED***
				combinedMetrics[location] = calculateMean(combinedMetrics[location], envelope, currentTime, vuCount)
				continue
			***REMOVED***

			if metricName == "http_req_tls_handshaking" && vuCount != -1 ***REMOVED***
				combinedMetrics[location] = calculateMean(combinedMetrics[location], envelope, currentTime, vuCount)
				continue
			***REMOVED***

			if metricName == "http_req_failed" && vuCount != -1 ***REMOVED***
				combinedMetrics[location] = calculateMean(combinedMetrics[location], envelope, currentTime, vuCount)
				continue
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// If any keys are empty, remove them
	for key, value := range combinedMetrics ***REMOVED***
		if len(value) == 0 ***REMOVED***
			delete(combinedMetrics, key)
		***REMOVED***
	***REMOVED***

	return combinedMetrics, nil
***REMOVED***

func findVuCount(workerVuCount *map[string]int, envelopeMap map[string][]*wrappedEnvelope, workerId string) (int, error) ***REMOVED***
	// Find worker Vu count in workerVuCount map first
	if vuCount, ok := (*workerVuCount)[workerId]; ok ***REMOVED***
		return vuCount, nil
	***REMOVED***

	// Check if vus is keyed in the envelope map
	if _, ok := (envelopeMap)["vus"]; !ok ***REMOVED***
		return -1, errors.New("vu count not found")
	***REMOVED***

	// Find vu count for this workerId
	for _, envelope := range envelopeMap["vus"] ***REMOVED***
		if envelope.workerId == workerId ***REMOVED***
			newCount := int((*envelope).SampleEnvelope.Data.Value)
			(*workerVuCount)[workerId] = newCount
			return newCount, nil
		***REMOVED***
	***REMOVED***

	return -1, errors.New("vu count not found")
***REMOVED***

func (store *cachedMetricsStore) Stop() ***REMOVED***
	store.flusher.Stop()
***REMOVED***
