package csv

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"gopkg.in/guregu/null.v3"

	"github.com/mstoykov/envconfig"
	"github.com/sirupsen/logrus"
	"go.k6.io/k6/lib/types"
)

// Config is the config for the csv output
type Config struct ***REMOVED***
	// Samples.
	FileName     null.String        `json:"file_name" envconfig:"K6_CSV_FILENAME"`
	SaveInterval types.NullDuration `json:"save_interval" envconfig:"K6_CSV_SAVE_INTERVAL"`
	TimeFormat   null.String        `json:"time_format" envconfig:"K6_CSV_TIME_FORMAT"`
***REMOVED***

// TimeFormat custom enum type
//go:generate enumer -type=TimeFormat -transform=snake -trimprefix TimeFormat -output time_format_gen.go
type TimeFormat uint8

// valid defined values for TimeFormat
const (
	TimeFormatUnix TimeFormat = iota
	TimeFormatRFC3339
)

// NewConfig creates a new Config instance with default values for some fields.
func NewConfig() Config ***REMOVED***
	return Config***REMOVED***
		FileName:     null.NewString("file.csv", false),
		SaveInterval: types.NewNullDuration(1*time.Second, false),
		TimeFormat:   null.NewString("unix", false),
	***REMOVED***
***REMOVED***

// Apply merges two configs by overwriting properties in the old config
func (c Config) Apply(cfg Config) Config ***REMOVED***
	if cfg.FileName.Valid ***REMOVED***
		c.FileName = cfg.FileName
	***REMOVED***
	if cfg.SaveInterval.Valid ***REMOVED***
		c.SaveInterval = cfg.SaveInterval
	***REMOVED***
	if cfg.TimeFormat.Valid ***REMOVED***
		c.TimeFormat = cfg.TimeFormat
	***REMOVED***
	return c
***REMOVED***

// ParseArg takes an arg string and converts it to a config
func ParseArg(arg string, logger *logrus.Logger) (Config, error) ***REMOVED***
	c := NewConfig()

	if !strings.Contains(arg, "=") ***REMOVED***
		c.FileName = null.StringFrom(arg)
		return c, nil
	***REMOVED***

	pairs := strings.Split(arg, ",")
	for _, pair := range pairs ***REMOVED***
		r := strings.SplitN(pair, "=", 2)
		if len(r) != 2 ***REMOVED***
			return c, fmt.Errorf("couldn't parse %q as argument for csv output", arg)
		***REMOVED***
		switch r[0] ***REMOVED***
		case "save_interval":
			logger.Warnf("CSV output argument '%s' is deprecated, please use 'saveInterval' instead.", r[0])
			fallthrough
		case "saveInterval":
			err := c.SaveInterval.UnmarshalText([]byte(r[1]))
			if err != nil ***REMOVED***
				return c, err
			***REMOVED***
		case "file_name":
			logger.Warnf("CSV output argument '%s' is deprecated, please use 'fileName' instead.", r[0])
			fallthrough
		case "fileName":
			c.FileName = null.StringFrom(r[1])
		case "timeFormat":
			c.TimeFormat = null.StringFrom(r[1])

		default:
			return c, fmt.Errorf("unknown key %q as argument for csv output", r[0])
		***REMOVED***
	***REMOVED***

	return c, nil
***REMOVED***

// GetConsolidatedConfig combines ***REMOVED***default config values + JSON config +
// environment vars + arg config values***REMOVED***, and returns the final result.
func GetConsolidatedConfig(
	jsonRawConf json.RawMessage, env map[string]string, arg string, logger *logrus.Logger,
) (Config, error) ***REMOVED***
	result := NewConfig()
	if jsonRawConf != nil ***REMOVED***
		jsonConf := Config***REMOVED******REMOVED***
		if err := json.Unmarshal(jsonRawConf, &jsonConf); err != nil ***REMOVED***
			return result, err
		***REMOVED***
		result = result.Apply(jsonConf)
	***REMOVED***

	envConfig := Config***REMOVED******REMOVED***
	if err := envconfig.Process("", &envConfig, func(key string) (string, bool) ***REMOVED***
		v, ok := env[key]
		return v, ok
	***REMOVED***); err != nil ***REMOVED***
		// TODO: get rid of envconfig and actually use the env parameter...
		return result, err
	***REMOVED***
	result = result.Apply(envConfig)

	if arg != "" ***REMOVED***
		urlConf, err := ParseArg(arg, logger)
		if err != nil ***REMOVED***
			return result, err
		***REMOVED***
		result = result.Apply(urlConf)
	***REMOVED***

	return result, nil
***REMOVED***
