package executor

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"gopkg.in/guregu/null.v3"

	"github.com/APITeamLimited/k6-worker/lib/consts"
	"github.com/APITeamLimited/k6-worker/lib/types"
)

// DefaultGracefulStopValue is the graceful top value for all executors, unless
// it's manually changed by the gracefulStop in each one.
// TODO?: Discard? Or make this actually user-configurable somehow? hello #883...
var DefaultGracefulStopValue = 30 * time.Second //nolint:gochecknoglobals

var executorNameWhitelist = regexp.MustCompile(`^[0-9a-zA-Z_-]+$`) //nolint:gochecknoglobals
const executorNameErr = "the executor name should contain only numbers, latin letters, underscores, and dashes"

// BaseConfig contains the common config fields for all executors
type BaseConfig struct ***REMOVED***
	Name         string             `json:"-"` // set via the JS object key
	Type         string             `json:"executor"`
	StartTime    types.NullDuration `json:"startTime"`
	GracefulStop types.NullDuration `json:"gracefulStop"`
	Env          map[string]string  `json:"env"`
	Exec         null.String        `json:"exec"` // function name, externally validated
	Tags         map[string]string  `json:"tags"`

	// TODO: future extensions like distribution, others?
***REMOVED***

// NewBaseConfig returns a default base config with the default values
func NewBaseConfig(name, configType string) BaseConfig ***REMOVED***
	return BaseConfig***REMOVED***
		Name:         name,
		Type:         configType,
		GracefulStop: types.NewNullDuration(DefaultGracefulStopValue, false),
	***REMOVED***
***REMOVED***

// Validate checks some basic things like present name, type, and a positive start time
func (bc BaseConfig) Validate() (errors []error) ***REMOVED***
	// Some just-in-case checks, since those things are likely checked in other places or
	// even assigned by us:
	if bc.Name == "" ***REMOVED***
		errors = append(errors, fmt.Errorf("executor name can't be empty"))
	***REMOVED***
	if !executorNameWhitelist.MatchString(bc.Name) ***REMOVED***
		errors = append(errors, fmt.Errorf(executorNameErr))
	***REMOVED***
	if bc.Exec.Valid && bc.Exec.String == "" ***REMOVED***
		errors = append(errors, fmt.Errorf("exec value cannot be empty"))
	***REMOVED***
	if bc.Type == "" ***REMOVED***
		errors = append(errors, fmt.Errorf("missing or empty type field"))
	***REMOVED***
	// The actually reasonable checks:
	if bc.StartTime.Duration < 0 ***REMOVED***
		errors = append(errors, fmt.Errorf("the startTime can't be negative"))
	***REMOVED***
	if bc.GracefulStop.Duration < 0 ***REMOVED***
		errors = append(errors, fmt.Errorf("the gracefulStop timeout can't be negative"))
	***REMOVED***
	return errors
***REMOVED***

// GetName returns the name of the executor.
func (bc BaseConfig) GetName() string ***REMOVED***
	return bc.Name
***REMOVED***

// GetType returns the executor's type as a string ID.
func (bc BaseConfig) GetType() string ***REMOVED***
	return bc.Type
***REMOVED***

// GetStartTime returns the starting time, relative to the beginning of the
// actual test, that this executor is supposed to execute.
func (bc BaseConfig) GetStartTime() time.Duration ***REMOVED***
	return bc.StartTime.TimeDuration()
***REMOVED***

// GetGracefulStop returns how long k6 is supposed to wait for any still
// running iterations to finish executing at the end of the normal executor
// duration, before it actually kills them.
//
// Of course, that doesn't count when the user manually interrupts the test,
// then iterations are immediately stopped.
func (bc BaseConfig) GetGracefulStop() time.Duration ***REMOVED***
	return bc.GracefulStop.TimeDuration()
***REMOVED***

// GetEnv returns any specific environment key=value pairs that
// are configured for the executor.
func (bc BaseConfig) GetEnv() map[string]string ***REMOVED***
	return bc.Env
***REMOVED***

// GetExec returns the configured custom exec value, if any.
func (bc BaseConfig) GetExec() string ***REMOVED***
	exec := bc.Exec.ValueOrZero()
	if exec == "" ***REMOVED***
		exec = consts.DefaultFn
	***REMOVED***
	return exec
***REMOVED***

// GetTags returns any custom tags configured for the executor.
func (bc BaseConfig) GetTags() map[string]string ***REMOVED***
	return bc.Tags
***REMOVED***

// IsDistributable returns true since by default all executors could be run in
// a distributed manner.
func (bc BaseConfig) IsDistributable() bool ***REMOVED***
	return true
***REMOVED***

// getBaseInfo is a helper method for the "parent" String methods.
func (bc BaseConfig) getBaseInfo(facts ...string) string ***REMOVED***
	if bc.Exec.Valid ***REMOVED***
		facts = append(facts, fmt.Sprintf("exec: %s", bc.Exec.String))
	***REMOVED***
	if bc.StartTime.Duration > 0 ***REMOVED***
		facts = append(facts, fmt.Sprintf("startTime: %s", bc.StartTime.Duration))
	***REMOVED***
	if bc.GracefulStop.Duration > 0 ***REMOVED***
		facts = append(facts, fmt.Sprintf("gracefulStop: %s", bc.GracefulStop.Duration))
	***REMOVED***
	if len(facts) == 0 ***REMOVED***
		return ""
	***REMOVED***
	return " (" + strings.Join(facts, ", ") + ")"
***REMOVED***
