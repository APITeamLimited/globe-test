package scheduler

// Config is an interface that should be implemented by all scheduler config types
type Config interface ***REMOVED***
	GetBaseConfig() BaseConfig
	Validate() []error
	// TODO: Split(percentages []float64) ([]Config, error)
	// TODO: GetMaxVUs() int64
	// TODO: GetMaxDuration() time.Duration // inclusind max timeouts, if we want to share VUs between schedulers in the future?
***REMOVED***

// Scheduler is an interface that should be implemented by all scheduler implementations
type Scheduler interface ***REMOVED***
	Initialize(Config) error
	Start()
***REMOVED***
