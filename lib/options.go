/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2016 Load Impact
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package lib

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net"
	"reflect"
	"regexp"
	"strings"

	"github.com/loadimpact/k6/lib/scheduler"
	"github.com/loadimpact/k6/lib/types"
	"github.com/loadimpact/k6/stats"
	"github.com/pkg/errors"
	"gopkg.in/guregu/null.v3"
)

// DefaultSchedulerName is used as the default key/ID of the scheduler config entries
// that were created due to the use of the shortcut execution control options (i.e. duration+vus,
// iterations+vus, or stages)
const DefaultSchedulerName = "default"

// DefaultSummaryTrendStats are the default trend columns shown in the test summary output
// nolint: gochecknoglobals
var DefaultSummaryTrendStats = []string***REMOVED***"avg", "min", "med", "max", "p(90)", "p(95)"***REMOVED***

// Describes a TLS version. Serialised to/from JSON as a string, eg. "tls1.2".
type TLSVersion int

func (v TLSVersion) MarshalJSON() ([]byte, error) ***REMOVED***
	return []byte(`"` + SupportedTLSVersionsToString[v] + `"`), nil
***REMOVED***

func (v *TLSVersion) UnmarshalJSON(data []byte) error ***REMOVED***
	var str string
	if err := json.Unmarshal(data, &str); err != nil ***REMOVED***
		return err
	***REMOVED***
	if str == "" ***REMOVED***
		*v = 0
		return nil
	***REMOVED***
	ver, ok := SupportedTLSVersions[str]
	if !ok ***REMOVED***
		return errors.Errorf("unknown TLS version: %s", str)
	***REMOVED***
	*v = ver
	return nil
***REMOVED***

// Fields for TLSVersions. Unmarshalling hack.
type TLSVersionsFields struct ***REMOVED***
	Min TLSVersion `json:"min"` // Minimum allowed version, 0 = any.
	Max TLSVersion `json:"max"` // Maximum allowed version, 0 = any.
***REMOVED***

// Describes a set (min/max) of TLS versions.
type TLSVersions TLSVersionsFields

func (v *TLSVersions) UnmarshalJSON(data []byte) error ***REMOVED***
	var fields TLSVersionsFields
	if err := json.Unmarshal(data, &fields); err != nil ***REMOVED***
		var ver TLSVersion
		if err2 := json.Unmarshal(data, &ver); err2 != nil ***REMOVED***
			return err
		***REMOVED***
		fields.Min = ver
		fields.Max = ver
	***REMOVED***
	*v = TLSVersions(fields)
	return nil
***REMOVED***

func (v *TLSVersions) isTLS13() bool ***REMOVED***
	return v.Min == TLSVersion13 || v.Max == TLSVersion13
***REMOVED***

// A list of TLS cipher suites.
// Marshals and unmarshals from a list of names, eg. "TLS_ECDHE_RSA_WITH_RC4_128_SHA".
// BUG: This currently doesn't marshal back to JSON properly!!
type TLSCipherSuites []uint16

func (s *TLSCipherSuites) UnmarshalJSON(data []byte) error ***REMOVED***
	var suiteNames []string
	if err := json.Unmarshal(data, &suiteNames); err != nil ***REMOVED***
		return err
	***REMOVED***

	var suiteIDs []uint16
	for _, name := range suiteNames ***REMOVED***
		if suiteID, ok := SupportedTLSCipherSuites[name]; ok ***REMOVED***
			suiteIDs = append(suiteIDs, suiteID)
		***REMOVED*** else ***REMOVED***
			return errors.New("Unknown cipher suite: " + name)
		***REMOVED***
	***REMOVED***

	*s = suiteIDs

	return nil
***REMOVED***

// Fields for TLSAuth. Unmarshalling hack.
type TLSAuthFields struct ***REMOVED***
	// Certificate and key as a PEM-encoded string, including "-----BEGIN CERTIFICATE-----".
	Cert string `json:"cert"`
	Key  string `json:"key"`

	// Domains to present the certificate to. May contain wildcards, eg. "*.example.com".
	Domains []string `json:"domains"`
***REMOVED***

// Defines a TLS client certificate to present to certain hosts.
type TLSAuth struct ***REMOVED***
	TLSAuthFields
	certificate *tls.Certificate
***REMOVED***

func (c *TLSAuth) UnmarshalJSON(data []byte) error ***REMOVED***
	if err := json.Unmarshal(data, &c.TLSAuthFields); err != nil ***REMOVED***
		return err
	***REMOVED***
	if _, err := c.Certificate(); err != nil ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

func (c *TLSAuth) Certificate() (*tls.Certificate, error) ***REMOVED***
	if c.certificate == nil ***REMOVED***
		cert, err := tls.X509KeyPair([]byte(c.Cert), []byte(c.Key))
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		c.certificate = &cert
	***REMOVED***
	return c.certificate, nil
***REMOVED***

// IPNet is a wrapper around net.IPNet for JSON unmarshalling
type IPNet net.IPNet

func (ipnet *IPNet) String() string ***REMOVED***
	return (*net.IPNet)(ipnet).String()
***REMOVED***

// UnmarshalText populates the IPNet from the given CIDR
func (ipnet *IPNet) UnmarshalText(b []byte) error ***REMOVED***
	newIPNet, err := ParseCIDR(string(b))
	if err != nil ***REMOVED***
		return errors.Wrap(err, "Failed to parse CIDR")
	***REMOVED***

	*ipnet = *newIPNet

	return nil
***REMOVED***

// ParseCIDR creates an IPNet out of a CIDR string
func ParseCIDR(s string) (*IPNet, error) ***REMOVED***
	_, ipnet, err := net.ParseCIDR(s)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	parsedIPNet := IPNet(*ipnet)

	return &parsedIPNet, nil
***REMOVED***

// HostnameTrie is a tree-structured list of hostname matches with support
// for wildcards exclusively at the start of the pattern. Items may only
// be inserted and searched. Internationalized hostnames are valid.
type HostnameTrie struct ***REMOVED***
	r        rune
	children []*HostnameTrie
	terminal bool // end of a valid match
***REMOVED***

// describes a valid hostname pattern to block by. Global var to avoid
// compilation penalty each call to ValidHostname.
var validHostnamePattern *regexp.Regexp = regexp.MustCompile("^\\*?(\\pL|[0-9\\.])*")

// ValidHostname returns whether the provided hostname pattern
// has an optional wildcard at the start, and is composed entirely
// of letters, numbers, or '.'s.
func ValidHostname(s string) error ***REMOVED***
	if len(validHostnamePattern.FindString(s)) != len(s) ***REMOVED***
		return fmt.Errorf("invalid hostname pattern %s", s)
	***REMOVED***
	return nil
***REMOVED***

// UnmarshalJSON forms a HostnameTrie from the provided hostname pattern
// list.
func (t *HostnameTrie) UnmarshalJSON(data []byte) error ***REMOVED***
	m := make([]string, 0)
	if err := json.Unmarshal(data, &m); err != nil ***REMOVED***
		return err
	***REMOVED***
	for _, h := range m ***REMOVED***
		if insertErr := t.Insert(h); insertErr != nil ***REMOVED***
			return insertErr
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// UnmarshalText forms a HostnameTrie from a comma-delimited list
// of hostname patterns.
func (t *HostnameTrie) UnmarshalText(b []byte) error ***REMOVED***
	for _, s := range strings.Split(string(b), ",") ***REMOVED***
		if err := t.Insert(s); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// Insert a string into the given HostnameTrie.
func (t *HostnameTrie) Insert(s string) error ***REMOVED***
	if len(s) == 0 ***REMOVED***
		return nil
	***REMOVED***

	if err := ValidHostname(s); err != nil ***REMOVED***
		return err
	***REMOVED***

	rStr := []rune(s) // need to iterate by runes for intl' names
	last := len(rStr) - 1
	for _, c := range t.children ***REMOVED***
		if c.r == rStr[last] ***REMOVED***
			return c.Insert(string(rStr[:last]))
		***REMOVED***
	***REMOVED***

	n := &HostnameTrie***REMOVED***rStr[last], nil, len(rStr) == 1***REMOVED***
	t.children = append(t.children, n)
	return n.Insert(string(rStr[:last]))
***REMOVED***

// Contains returns whether s matches a pattern in the HostnameTrie
// along with the matching pattern, if one was found.
func (t *HostnameTrie) Contains(s string) (bool, string) ***REMOVED***
	for _, c := range t.children ***REMOVED***
		if b, m := c.childContains(s, ""); b ***REMOVED***
			return b, m
		***REMOVED***
	***REMOVED***
	return false, ""
***REMOVED***

func (t *HostnameTrie) childContains(s string, match string) (bool, string) ***REMOVED***
	if len(s) == 0 ***REMOVED***
		return false, ""
	***REMOVED***

	rStr := []rune(s)
	last := len(rStr) - 1

	switch ***REMOVED***
	case t.r == '*': // wildcard encounters validate the string
		return true, string(t.r) + match
	case t.r != rStr[last]:
		return false, ""
	case len(s) == 1:
		return t.terminal, string(t.r) + match
	default:
		for _, c := range t.children ***REMOVED***
			if b, m := c.childContains(string(rStr[:last]), string(t.r)+match); b ***REMOVED***
				return b, m
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return false, ""
***REMOVED***

type Options struct ***REMOVED***
	// Should the test start in a paused state?
	Paused null.Bool `json:"paused" envconfig:"K6_PAUSED"`

	// Initial values for VUs, max VUs, duration cap, iteration cap, and stages.
	// See the Runner or Executor interfaces for more information.
	VUs null.Int `json:"vus" envconfig:"K6_VUS"`

	//TODO: deprecate this? or reuse it in the manual control "scheduler"?
	VUsMax     null.Int           `json:"vusMax" envconfig:"K6_VUS_MAX"`
	Duration   types.NullDuration `json:"duration" envconfig:"K6_DURATION"`
	Iterations null.Int           `json:"iterations" envconfig:"K6_ITERATIONS"`
	Stages     []Stage            `json:"stages" envconfig:"K6_STAGES"`

	Execution scheduler.ConfigMap `json:"execution,omitempty" envconfig:"-"`

	// Timeouts for the setup() and teardown() functions
	SetupTimeout    types.NullDuration `json:"setupTimeout" envconfig:"K6_SETUP_TIMEOUT"`
	TeardownTimeout types.NullDuration `json:"teardownTimeout" envconfig:"K6_TEARDOWN_TIMEOUT"`

	// Limit HTTP requests per second.
	RPS null.Int `json:"rps" envconfig:"K6_RPS"`

	// How many HTTP redirects do we follow?
	MaxRedirects null.Int `json:"maxRedirects" envconfig:"K6_MAX_REDIRECTS"`

	// Default User Agent string for HTTP requests.
	UserAgent null.String `json:"userAgent" envconfig:"K6_USER_AGENT"`

	// How many batch requests are allowed in parallel, in total and per host?
	Batch        null.Int `json:"batch" envconfig:"K6_BATCH"`
	BatchPerHost null.Int `json:"batchPerHost" envconfig:"K6_BATCH_PER_HOST"`

	// Should all HTTP requests and responses be logged (excluding body)?
	HTTPDebug null.String `json:"httpDebug" envconfig:"K6_HTTP_DEBUG"`

	// Accept invalid or untrusted TLS certificates.
	InsecureSkipTLSVerify null.Bool `json:"insecureSkipTLSVerify" envconfig:"K6_INSECURE_SKIP_TLS_VERIFY"`

	// Specify TLS versions and cipher suites, and present client certificates.
	TLSCipherSuites *TLSCipherSuites `json:"tlsCipherSuites" envconfig:"K6_TLS_CIPHER_SUITES"`
	TLSVersion      *TLSVersions     `json:"tlsVersion" envconfig:"K6_TLS_VERSION"`
	TLSAuth         []*TLSAuth       `json:"tlsAuth" envconfig:"K6_TLSAUTH"`

	// Throw warnings (eg. failed HTTP requests) as errors instead of simply logging them.
	Throw null.Bool `json:"throw" envconfig:"K6_THROW"`

	// Define thresholds; these take the form of 'metric=["snippet1", "snippet2"]'.
	// To create a threshold on a derived metric based on tag queries ("submetrics"), create a
	// metric on a nonexistent metric named 'real_metric***REMOVED***tagA:valueA,tagB:valueB***REMOVED***'.
	Thresholds map[string]stats.Thresholds `json:"thresholds" envconfig:"K6_THRESHOLDS"`

	// Blacklist IP ranges that tests may not contact. Mainly useful in hosted setups.
	BlacklistIPs []*IPNet `json:"blacklistIPs" envconfig:"K6_BLACKLIST_IPS"`

	// Block hostnames that tests may not contact.
	BlockedHostnames *HostnameTrie `json:"blockHostnames" envconfig:"K6_BLOCK_HOSTNAMES"`

	// Hosts overrides dns entries for given hosts
	Hosts map[string]net.IP `json:"hosts" envconfig:"K6_HOSTS"`

	// Disable keep-alive connections
	NoConnectionReuse null.Bool `json:"noConnectionReuse" envconfig:"K6_NO_CONNECTION_REUSE"`

	// Do not reuse connections between VU iterations. This gives more realistic results (depending
	// on what you're looking for), but you need to raise various kernel limits or you'll get
	// errors about running out of file handles or sockets, or being unable to bind addresses.
	NoVUConnectionReuse null.Bool `json:"noVUConnectionReuse" envconfig:"K6_NO_VU_CONNECTION_REUSE"`

	// MinIterationDuration can be used to force VUs to pause between iterations if a specific
	// iteration is shorter than the specified value.
	MinIterationDuration types.NullDuration `json:"minIterationDuration" envconfig:"K6_MIN_ITERATION_DURATION"`

	// These values are for third party collectors' benefit.
	// Can't be set through env vars.
	External map[string]json.RawMessage `json:"ext" ignored:"true"`

	// Summary trend stats for trend metrics (response times) in CLI output
	SummaryTrendStats []string `json:"summaryTrendStats" envconfig:"K6_SUMMARY_TREND_STATS"`

	// Summary time unit for summary metrics (response times) in CLI output
	SummaryTimeUnit null.String `json:"summaryTimeUnit" envconfig:"K6_SUMMARY_TIME_UNIT"`

	// Which system tags to include with metrics ("method", "vu" etc.)
	// Use pointer for identifying whether user provide any tag or not.
	SystemTags *stats.SystemTagSet `json:"systemTags" envconfig:"K6_SYSTEM_TAGS"`

	// Tags to be applied to all samples for this running
	RunTags *stats.SampleTags `json:"tags" envconfig:"K6_TAGS"`

	// Buffer size of the channel for metric samples; 0 means unbuffered
	MetricSamplesBufferSize null.Int `json:"metricSamplesBufferSize" envconfig:"K6_METRIC_SAMPLES_BUFFER_SIZE"`

	// Do not reset cookies after a VU iteration
	NoCookiesReset null.Bool `json:"noCookiesReset" envconfig:"K6_NO_COOKIES_RESET"`

	// Discard Http Responses Body
	DiscardResponseBodies null.Bool `json:"discardResponseBodies" envconfig:"K6_DISCARD_RESPONSE_BODIES"`

	// Redirect console logging to a file
	ConsoleOutput null.String `json:"-" envconfig:"K6_CONSOLE_OUTPUT"`
***REMOVED***

// Returns the result of overwriting any fields with any that are set on the argument.
//
// Example:
//   a := Options***REMOVED***VUs: null.IntFrom(10), VUsMax: null.IntFrom(10)***REMOVED***
//   b := Options***REMOVED***VUs: null.IntFrom(5)***REMOVED***
//   a.Apply(b) // Options***REMOVED***VUs: null.IntFrom(5), VUsMax: null.IntFrom(10)***REMOVED***
func (o Options) Apply(opts Options) Options ***REMOVED***
	if opts.Paused.Valid ***REMOVED***
		o.Paused = opts.Paused
	***REMOVED***
	if opts.VUs.Valid ***REMOVED***
		o.VUs = opts.VUs
	***REMOVED***
	if opts.VUsMax.Valid ***REMOVED***
		o.VUsMax = opts.VUsMax
	***REMOVED***

	// Specifying duration, iterations, stages, or execution in a "higher" config tier
	// will overwrite all of the the previous execution settings (if any) from any
	// "lower" config tiers
	// Still, if more than one of those options is simultaneously specified in the same
	// config tier, they will be preserved, so the validation after we've consolidated
	// all of the options can return an error.
	if opts.Duration.Valid || opts.Iterations.Valid || opts.Stages != nil || opts.Execution != nil ***REMOVED***
		//TODO: uncomment this after we start using the new schedulers
		/*
			o.Duration = types.NewNullDuration(0, false)
			o.Iterations = null.NewInt(0, false)
			o.Stages = nil
		*/
		o.Execution = nil
	***REMOVED***

	if opts.Duration.Valid ***REMOVED***
		o.Duration = opts.Duration
	***REMOVED***
	if opts.Iterations.Valid ***REMOVED***
		o.Iterations = opts.Iterations
	***REMOVED***
	if opts.Stages != nil ***REMOVED***
		o.Stages = []Stage***REMOVED******REMOVED***
		for _, s := range opts.Stages ***REMOVED***
			if s.Duration.Valid ***REMOVED***
				o.Stages = append(o.Stages, s)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	// o.Execution can also be populated by the duration/iterations/stages config shortcuts, but
	// that happens after the configuration from the different sources is consolidated. It can't
	// happen here, because something like `K6_ITERATIONS=10 k6 run --vus 5 script.js` wont't
	// work correctly at this level.
	if opts.Execution != nil ***REMOVED***
		o.Execution = opts.Execution
	***REMOVED***
	if opts.SetupTimeout.Valid ***REMOVED***
		o.SetupTimeout = opts.SetupTimeout
	***REMOVED***
	if opts.TeardownTimeout.Valid ***REMOVED***
		o.TeardownTimeout = opts.TeardownTimeout
	***REMOVED***
	if opts.RPS.Valid ***REMOVED***
		o.RPS = opts.RPS
	***REMOVED***
	if opts.MaxRedirects.Valid ***REMOVED***
		o.MaxRedirects = opts.MaxRedirects
	***REMOVED***
	if opts.UserAgent.Valid ***REMOVED***
		o.UserAgent = opts.UserAgent
	***REMOVED***
	if opts.Batch.Valid ***REMOVED***
		o.Batch = opts.Batch
	***REMOVED***
	if opts.BatchPerHost.Valid ***REMOVED***
		o.BatchPerHost = opts.BatchPerHost
	***REMOVED***
	if opts.HTTPDebug.Valid ***REMOVED***
		o.HTTPDebug = opts.HTTPDebug
	***REMOVED***
	if opts.InsecureSkipTLSVerify.Valid ***REMOVED***
		o.InsecureSkipTLSVerify = opts.InsecureSkipTLSVerify
	***REMOVED***
	if opts.TLSCipherSuites != nil ***REMOVED***
		o.TLSCipherSuites = opts.TLSCipherSuites
	***REMOVED***
	if opts.TLSVersion != nil ***REMOVED***
		o.TLSVersion = opts.TLSVersion
		if o.TLSVersion.isTLS13() ***REMOVED***
			enableTLS13()
		***REMOVED***
	***REMOVED***
	if opts.TLSAuth != nil ***REMOVED***
		o.TLSAuth = opts.TLSAuth
	***REMOVED***
	if opts.Throw.Valid ***REMOVED***
		o.Throw = opts.Throw
	***REMOVED***
	if opts.Thresholds != nil ***REMOVED***
		o.Thresholds = opts.Thresholds
	***REMOVED***
	if opts.BlacklistIPs != nil ***REMOVED***
		o.BlacklistIPs = opts.BlacklistIPs
	***REMOVED***
	if opts.BlockedHostnames != nil ***REMOVED***
		o.BlockedHostnames = opts.BlockedHostnames
	***REMOVED***
	if opts.Hosts != nil ***REMOVED***
		o.Hosts = opts.Hosts
	***REMOVED***
	if opts.NoConnectionReuse.Valid ***REMOVED***
		o.NoConnectionReuse = opts.NoConnectionReuse
	***REMOVED***
	if opts.NoVUConnectionReuse.Valid ***REMOVED***
		o.NoVUConnectionReuse = opts.NoVUConnectionReuse
	***REMOVED***
	if opts.MinIterationDuration.Valid ***REMOVED***
		o.MinIterationDuration = opts.MinIterationDuration
	***REMOVED***
	if opts.NoCookiesReset.Valid ***REMOVED***
		o.NoCookiesReset = opts.NoCookiesReset
	***REMOVED***
	if opts.External != nil ***REMOVED***
		o.External = opts.External
	***REMOVED***
	if opts.SummaryTrendStats != nil ***REMOVED***
		o.SummaryTrendStats = opts.SummaryTrendStats
	***REMOVED***
	if opts.SummaryTimeUnit.Valid ***REMOVED***
		o.SummaryTimeUnit = opts.SummaryTimeUnit
	***REMOVED***
	if opts.SystemTags != nil ***REMOVED***
		o.SystemTags = opts.SystemTags
	***REMOVED***
	if !opts.RunTags.IsEmpty() ***REMOVED***
		o.RunTags = opts.RunTags
	***REMOVED***
	if opts.MetricSamplesBufferSize.Valid ***REMOVED***
		o.MetricSamplesBufferSize = opts.MetricSamplesBufferSize
	***REMOVED***
	if opts.DiscardResponseBodies.Valid ***REMOVED***
		o.DiscardResponseBodies = opts.DiscardResponseBodies
	***REMOVED***
	if opts.ConsoleOutput.Valid ***REMOVED***
		o.ConsoleOutput = opts.ConsoleOutput
	***REMOVED***

	return o
***REMOVED***

// Validate checks if all of the specified options make sense
func (o Options) Validate() []error ***REMOVED***
	//TODO: validate all of the other options... that we should have already been validating...
	//TODO: maybe integrate an external validation lib: https://github.com/avelino/awesome-go#validation
	return o.Execution.Validate()
***REMOVED***

// ForEachSpecified enumerates all struct fields and calls the supplied function with each
// element that is valid. It panics for any unfamiliar or unexpected fields, so make sure
// new fields in Options are accounted for.
func (o Options) ForEachSpecified(structTag string, callback func(key string, value interface***REMOVED******REMOVED***)) ***REMOVED***
	structType := reflect.TypeOf(o)
	structVal := reflect.ValueOf(o)
	for i := 0; i < structType.NumField(); i++ ***REMOVED***
		fieldType := structType.Field(i)
		fieldVal := structVal.Field(i)
		value := fieldVal.Interface()

		shouldCall := false
		switch fieldType.Type.Kind() ***REMOVED***
		case reflect.Struct:
			// Unpack any guregu/null values
			shouldCall = fieldVal.FieldByName("Valid").Bool()
			valOrZero := fieldVal.MethodByName("ValueOrZero")
			if shouldCall && valOrZero.IsValid() ***REMOVED***
				value = valOrZero.Call([]reflect.Value***REMOVED******REMOVED***)[0].Interface()
				if v, ok := value.(types.Duration); ok ***REMOVED***
					value = v.String()
				***REMOVED***
			***REMOVED***
		case reflect.Slice:
			shouldCall = fieldVal.Len() > 0
		case reflect.Map:
			shouldCall = fieldVal.Len() > 0
		case reflect.Ptr:
			shouldCall = !fieldVal.IsNil()
		default:
			panic(fmt.Sprintf("Unknown Options field %#v", fieldType))
		***REMOVED***

		if shouldCall ***REMOVED***
			key, ok := fieldType.Tag.Lookup(structTag)
			if !ok ***REMOVED***
				key = fieldType.Name
			***REMOVED***

			callback(key, value)
		***REMOVED***
	***REMOVED***
***REMOVED***
