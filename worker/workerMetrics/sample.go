package workerMetrics

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/mailru/easyjson/jwriter"
)

// A Sample is a single measurement.
type Sample struct ***REMOVED***
	Metric *Metric
	Time   time.Time
	Tags   *SampleTags
	Value  float64
***REMOVED***

// SampleContainer is a simple abstraction that allows sample
// producers to attach extra information to samples they return
type SampleContainer interface ***REMOVED***
	GetSamples() []Sample
***REMOVED***

// Samples is just the simplest SampleContainer implementation
// that will be used when there's no need for extra information
type Samples []Sample

// GetSamples just implements the SampleContainer interface
func (s Samples) GetSamples() []Sample ***REMOVED***
	return s
***REMOVED***

// ConnectedSampleContainer is an extension of the SampleContainer
// interface that should be implemented when emitted samples
// are connected and share the same time and tags.
type ConnectedSampleContainer interface ***REMOVED***
	SampleContainer
	GetTags() *SampleTags
	GetTime() time.Time
***REMOVED***

// ConnectedSamples is the simplest ConnectedSampleContainer
// implementation that will be used when there's no need for
// extra information
type ConnectedSamples struct ***REMOVED***
	Samples []Sample
	Tags    *SampleTags
	Time    time.Time
***REMOVED***

// GetSamples implements the SampleContainer and ConnectedSampleContainer
// interfaces and returns the stored slice with samples.
func (cs ConnectedSamples) GetSamples() []Sample ***REMOVED***
	return cs.Samples
***REMOVED***

// GetTags implements ConnectedSampleContainer interface and returns stored tags.
func (cs ConnectedSamples) GetTags() *SampleTags ***REMOVED***
	return cs.Tags
***REMOVED***

// GetTime implements ConnectedSampleContainer interface and returns stored time.
func (cs ConnectedSamples) GetTime() time.Time ***REMOVED***
	return cs.Time
***REMOVED***

// GetSamples implement the ConnectedSampleContainer interface
// for a single Sample, since it's obviously connected with itself :)
func (s Sample) GetSamples() []Sample ***REMOVED***
	return []Sample***REMOVED***s***REMOVED***
***REMOVED***

// GetTags implements ConnectedSampleContainer interface
// and returns the sample's tags.
func (s Sample) GetTags() *SampleTags ***REMOVED***
	return s.Tags
***REMOVED***

// GetTime just implements ConnectedSampleContainer interface
// and returns the sample's time.
func (s Sample) GetTime() time.Time ***REMOVED***
	return s.Time
***REMOVED***

// Ensure that interfaces are implemented correctly
var (
	_ SampleContainer = Sample***REMOVED******REMOVED***
	_ SampleContainer = Samples***REMOVED******REMOVED***
)

var (
	_ ConnectedSampleContainer = Sample***REMOVED******REMOVED***
	_ ConnectedSampleContainer = ConnectedSamples***REMOVED******REMOVED***
)

// GetBufferedSamples will read all present (i.e. buffered or currently being pushed)
// values in the input channel and return them as a slice.
func GetBufferedSamples(input <-chan SampleContainer) (result []SampleContainer) ***REMOVED***
	for ***REMOVED***
		select ***REMOVED***
		case val, ok := <-input:
			if !ok ***REMOVED***
				return
			***REMOVED***
			result = append(result, val)
		default:
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

// PushIfNotDone first checks if the supplied context is done and doesn't push
// the sample container if it is.
func PushIfNotDone(ctx context.Context, output chan<- SampleContainer, sample SampleContainer) bool ***REMOVED***
	if ctx.Err() != nil ***REMOVED***
		return false
	***REMOVED***
	output <- sample
	return true
***REMOVED***

// GetResolversForTrendColumns checks if passed trend columns are valid for use in
// the summary output and then returns a map of the corresponding resolvers.
func GetResolversForTrendColumns(trendColumns []string) (map[string]func(s *TrendSink) float64, error) ***REMOVED***
	staticResolvers := map[string]func(s *TrendSink) float64***REMOVED***
		"avg":   func(s *TrendSink) float64 ***REMOVED*** return s.Avg ***REMOVED***,
		"min":   func(s *TrendSink) float64 ***REMOVED*** return s.Min ***REMOVED***,
		"med":   func(s *TrendSink) float64 ***REMOVED*** return s.Med ***REMOVED***,
		"max":   func(s *TrendSink) float64 ***REMOVED*** return s.Max ***REMOVED***,
		"count": func(s *TrendSink) float64 ***REMOVED*** return float64(s.Count) ***REMOVED***,
	***REMOVED***
	dynamicResolver := func(percentile float64) func(s *TrendSink) float64 ***REMOVED***
		return func(s *TrendSink) float64 ***REMOVED***
			return s.P(percentile / 100)
		***REMOVED***
	***REMOVED***

	result := make(map[string]func(s *TrendSink) float64, len(trendColumns))

	for _, stat := range trendColumns ***REMOVED***
		if staticStat, ok := staticResolvers[stat]; ok ***REMOVED***
			result[stat] = staticStat
			continue
		***REMOVED***

		percentile, err := parsePercentile(stat)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		result[stat] = dynamicResolver(percentile)
	***REMOVED***

	return result, nil
***REMOVED***

// parsePercentile is a helper function to parse and validate percentile notations
func parsePercentile(stat string) (float64, error) ***REMOVED***
	if !strings.HasPrefix(stat, "p(") || !strings.HasSuffix(stat, ")") ***REMOVED***
		return 0, fmt.Errorf("invalid trend stat '%s', unknown format", stat)
	***REMOVED***

	percentile, err := strconv.ParseFloat(stat[2:len(stat)-1], 64)

	if err != nil || (percentile < 0) || (percentile > 100) ***REMOVED***
		return 0, fmt.Errorf("invalid percentile trend stat value '%s', provide a number between 0 and 100", stat)
	***REMOVED***

	return percentile, nil
***REMOVED***

// SampleTags is an immutable string[string] map for tags. Once a tag
// set is created, direct modification is prohibited. It has
// copy-on-write semantics and uses pointers for faster comparison
// between maps, since the same tag set is often used for multiple samples.
// All methods should not panic, even if they are called on a nil pointer.
//easyjson:skip
type SampleTags struct ***REMOVED***
	tags map[string]string
	json []byte
***REMOVED***

// Get returns an empty string and false if the the requested key is not
// present or its value and true if it is.
func (st *SampleTags) Get(key string) (string, bool) ***REMOVED***
	if st == nil ***REMOVED***
		return "", false
	***REMOVED***
	val, ok := st.tags[key]
	return val, ok
***REMOVED***

// IsEmpty checks for a nil pointer or zero tags.
// It's necessary because of this envconfig issue: https://github.com/kelseyhightower/envconfig/issues/113
func (st *SampleTags) IsEmpty() bool ***REMOVED***
	return st == nil || len(st.tags) == 0
***REMOVED***

// IsEqual tries to compare two tag sets with maximum efficiency.
func (st *SampleTags) IsEqual(other *SampleTags) bool ***REMOVED***
	if st == other ***REMOVED***
		return true
	***REMOVED***
	if st == nil || other == nil || len(st.tags) != len(other.tags) ***REMOVED***
		return false
	***REMOVED***
	for k, v := range st.tags ***REMOVED***
		if otherv, ok := other.tags[k]; !ok || v != otherv ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

func (st *SampleTags) Contains(other *SampleTags) bool ***REMOVED***
	if st == other || other == nil ***REMOVED***
		return true
	***REMOVED***
	if st == nil || len(st.tags) < len(other.tags) ***REMOVED***
		return false
	***REMOVED***

	for k, v := range other.tags ***REMOVED***
		if myv, ok := st.tags[k]; !ok || myv != v ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***

	return true
***REMOVED***

// MarshalJSON serializes SampleTags to a JSON string and caches
// the result. It is not thread safe in the sense that the Go race
// detector will complain if it's used concurrently, but no data
// should be corrupted.
func (st *SampleTags) MarshalJSON() ([]byte, error) ***REMOVED***
	if st.IsEmpty() ***REMOVED***
		return []byte("null"), nil
	***REMOVED***
	if st.json != nil ***REMOVED***
		return st.json, nil
	***REMOVED***
	res, err := json.Marshal(st.tags)
	if err != nil ***REMOVED***
		return res, err
	***REMOVED***
	st.json = res
	return res, nil
***REMOVED***

// MarshalEasyJSON supports easyjson.Marshaler interface
func (st *SampleTags) MarshalEasyJSON(w *jwriter.Writer) ***REMOVED***
	w.RawByte('***REMOVED***')
	first := true
	for k, v := range st.tags ***REMOVED***
		if first ***REMOVED***
			first = false
		***REMOVED*** else ***REMOVED***
			w.RawByte(',')
		***REMOVED***
		w.String(k)
		w.RawByte(':')
		w.String(v)
	***REMOVED***
	w.RawByte('***REMOVED***')
***REMOVED***

// UnmarshalJSON deserializes SampleTags from a JSON string.
func (st *SampleTags) UnmarshalJSON(data []byte) error ***REMOVED***
	if st == nil ***REMOVED***
		*st = SampleTags***REMOVED******REMOVED***
	***REMOVED***
	return json.Unmarshal(data, &st.tags)
***REMOVED***

// CloneTags copies the underlying set of a sample tags and
// returns it. If the receiver is nil, it returns an empty non-nil map.
func (st *SampleTags) CloneTags() map[string]string ***REMOVED***
	if st == nil ***REMOVED***
		return map[string]string***REMOVED******REMOVED***
	***REMOVED***
	res := make(map[string]string, len(st.tags))
	for k, v := range st.tags ***REMOVED***
		res[k] = v
	***REMOVED***
	return res
***REMOVED***

// NewSampleTags *copies* the supplied tag set and returns a new SampleTags
// instance with the key-value pairs from it.
func NewSampleTags(data map[string]string) *SampleTags ***REMOVED***
	if len(data) == 0 ***REMOVED***
		return nil
	***REMOVED***

	tags := map[string]string***REMOVED******REMOVED***
	for k, v := range data ***REMOVED***
		tags[k] = v
	***REMOVED***
	return &SampleTags***REMOVED***tags: tags***REMOVED***
***REMOVED***

// IntoSampleTags "consumes" the passed map and creates a new SampleTags
// struct with the data. The map is set to nil as a hint that it shouldn't
// be changed after it has been transformed into an "immutable" tag set.
// Oh, how I miss Rust and move semantics... :)
func IntoSampleTags(data *map[string]string) *SampleTags ***REMOVED***
	if len(*data) == 0 ***REMOVED***
		return nil
	***REMOVED***

	res := SampleTags***REMOVED***tags: *data***REMOVED***
	*data = nil
	return &res
***REMOVED***
