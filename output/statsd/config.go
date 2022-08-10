package statsd

import (
	"encoding/json"
	"time"

	"github.com/mstoykov/envconfig"
	"gopkg.in/guregu/null.v3"

	"go.k6.io/k6/lib/types"
	"go.k6.io/k6/metrics"
)

// config defines the StatsD configuration.
type config struct ***REMOVED***
	Addr         null.String         `json:"addr,omitempty" envconfig:"K6_STATSD_ADDR"`
	BufferSize   null.Int            `json:"bufferSize,omitempty" envconfig:"K6_STATSD_BUFFER_SIZE"`
	Namespace    null.String         `json:"namespace,omitempty" envconfig:"K6_STATSD_NAMESPACE"`
	PushInterval types.NullDuration  `json:"pushInterval,omitempty" envconfig:"K6_STATSD_PUSH_INTERVAL"`
	TagBlocklist metrics.EnabledTags `json:"tagBlocklist,omitempty" envconfig:"K6_STATSD_TAG_BLOCKLIST"`
	EnableTags   null.Bool           `json:"enableTags,omitempty" envconfig:"K6_STATSD_ENABLE_TAGS"`
***REMOVED***

func processTags(t metrics.EnabledTags, tags map[string]string) []string ***REMOVED***
	var res []string
	for key, value := range tags ***REMOVED***
		if value != "" && !t[key] ***REMOVED***
			res = append(res, key+":"+value)
		***REMOVED***
	***REMOVED***
	return res
***REMOVED***

// Apply saves config non-zero config values from the passed config in the receiver.
func (c config) Apply(cfg config) config ***REMOVED***
	if cfg.Addr.Valid ***REMOVED***
		c.Addr = cfg.Addr
	***REMOVED***
	if cfg.BufferSize.Valid ***REMOVED***
		c.BufferSize = cfg.BufferSize
	***REMOVED***
	if cfg.Namespace.Valid ***REMOVED***
		c.Namespace = cfg.Namespace
	***REMOVED***
	if cfg.PushInterval.Valid ***REMOVED***
		c.PushInterval = cfg.PushInterval
	***REMOVED***
	if cfg.TagBlocklist != nil ***REMOVED***
		c.TagBlocklist = cfg.TagBlocklist
	***REMOVED***
	if cfg.EnableTags.Valid ***REMOVED***
		c.EnableTags = cfg.EnableTags
	***REMOVED***

	return c
***REMOVED***

// newConfig creates a new Config instance with default values for some fields.
func newConfig() config ***REMOVED***
	return config***REMOVED***
		Addr:         null.NewString("localhost:8125", false),
		BufferSize:   null.NewInt(20, false),
		Namespace:    null.NewString("k6.", false),
		PushInterval: types.NewNullDuration(1*time.Second, false),
		TagBlocklist: (metrics.TagVU | metrics.TagIter | metrics.TagURL).Map(),
		EnableTags:   null.NewBool(false, false),
	***REMOVED***
***REMOVED***

// getConsolidatedConfig combines ***REMOVED***default config values + JSON config +
// environment vars***REMOVED***, and returns the final result.
func getConsolidatedConfig(jsonRawConf json.RawMessage, env map[string]string, _ string) (config, error) ***REMOVED***
	result := newConfig()
	if jsonRawConf != nil ***REMOVED***
		jsonConf := config***REMOVED******REMOVED***
		if err := json.Unmarshal(jsonRawConf, &jsonConf); err != nil ***REMOVED***
			return result, err
		***REMOVED***
		result = result.Apply(jsonConf)
	***REMOVED***

	envConfig := config***REMOVED******REMOVED***
	_ = env // TODO: get rid of envconfig and actually use the env parameter...
	if err := envconfig.Process("", &envConfig, func(key string) (string, bool) ***REMOVED***
		v, ok := env[key]
		return v, ok
	***REMOVED***); err != nil ***REMOVED***
		return result, err
	***REMOVED***
	result = result.Apply(envConfig)

	return result, nil
***REMOVED***
