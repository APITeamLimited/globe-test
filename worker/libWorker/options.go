package libWorker

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"net"
	"reflect"
	"strconv"

	"github.com/APITeamLimited/globe-test/worker/libWorker/types"
	"github.com/APITeamLimited/globe-test/worker/workerMetrics"
	"gopkg.in/guregu/null.v3"
)

// DefaultScenarioName is used as the default key/ID of the scenario config entries
// that were created due to the use of the shortcut execution control options (i.e. duration+vus,
// iterations+vus, or stages)
const DefaultScenarioName = "default"

// DefaultSummaryTrendStats are the default trend columns shown in the test summary output
//nolint:gochecknoglobals
var DefaultSummaryTrendStats = []string***REMOVED***"avg", "min", "med", "max", "p(90)", "p(95)"***REMOVED***

// Describes a TLS version. Serialised to/from JSON as a string, eg. "tls1.2".
type TLSVersion int

func (v TLSVersion) MarshalJSON() ([]byte, error) ***REMOVED***
	return []byte(`"` + SupportedTLSVersionsToString[v] + `"`), nil
***REMOVED***

func (v *TLSVersion) UnmarshalJSON(data []byte) error ***REMOVED***
	var str string
	if err := StrictJSONUnmarshal(data, &str); err != nil ***REMOVED***
		return err
	***REMOVED***
	if str == "" ***REMOVED***
		*v = 0
		return nil
	***REMOVED***
	ver, ok := SupportedTLSVersions[str]
	if !ok ***REMOVED***
		return fmt.Errorf("unknown TLS version '%s'", str)
	***REMOVED***
	*v = ver
	return nil
***REMOVED***

// Fields for TLSVersions. Unmarshalling hack.
type TLSVersionsFields struct ***REMOVED***
	Min TLSVersion `json:"min" ignored:"true"` // Minimum allowed version, 0 = any.
	Max TLSVersion `json:"max" ignored:"true"` // Maximum allowed version, 0 = any.
***REMOVED***

// Describes a set (min/max) of TLS versions.
type TLSVersions TLSVersionsFields

func (v *TLSVersions) UnmarshalJSON(data []byte) error ***REMOVED***
	var fields TLSVersionsFields
	if err := StrictJSONUnmarshal(data, &fields); err != nil ***REMOVED***
		var ver TLSVersion
		if err2 := StrictJSONUnmarshal(data, &ver); err2 != nil ***REMOVED***
			return err
		***REMOVED***
		fields.Min = ver
		fields.Max = ver
	***REMOVED***
	*v = TLSVersions(fields)
	return nil
***REMOVED***

// A list of TLS cipher suites.
// Marshals and unmarshals from a list of names, eg. "TLS_ECDHE_RSA_WITH_RC4_128_SHA".
type TLSCipherSuites []uint16

// MarshalJSON will return the JSON representation according to supported TLS cipher suites
func (s *TLSCipherSuites) MarshalJSON() ([]byte, error) ***REMOVED***
	var suiteNames []string
	for _, id := range *s ***REMOVED***
		if suiteName, ok := SupportedTLSCipherSuitesToString[id]; ok ***REMOVED***
			suiteNames = append(suiteNames, suiteName)
		***REMOVED*** else ***REMOVED***
			return nil, fmt.Errorf("unknown cipher suite id '%d'", id)
		***REMOVED***
	***REMOVED***

	return json.Marshal(suiteNames)
***REMOVED***

func (s *TLSCipherSuites) UnmarshalJSON(data []byte) error ***REMOVED***
	var suiteNames []string
	if err := StrictJSONUnmarshal(data, &suiteNames); err != nil ***REMOVED***
		return err
	***REMOVED***

	var suiteIDs []uint16
	for _, name := range suiteNames ***REMOVED***
		if suiteID, ok := SupportedTLSCipherSuites[name]; ok ***REMOVED***
			suiteIDs = append(suiteIDs, suiteID)
		***REMOVED*** else ***REMOVED***
			return fmt.Errorf("unknown cipher suite '%s'", name)
		***REMOVED***
	***REMOVED***

	*s = suiteIDs

	return nil
***REMOVED***

// Fields for TLSAuth. Unmarshalling hack.
type TLSAuthFields struct ***REMOVED***
	// Certificate and key as a PEM-encoded string, including "-----BEGIN CERTIFICATE-----".
	Cert     string      `json:"cert"`
	Key      string      `json:"key"`
	Password null.String `json:"password"`

	// Domains to present the certificate to. May contain wildcards, eg. "*.example.com".
	Domains []string `json:"domains"`
***REMOVED***

// Defines a TLS client certificate to present to certain hosts.
type TLSAuth struct ***REMOVED***
	TLSAuthFields
	certificate *tls.Certificate
***REMOVED***

func (c *TLSAuth) UnmarshalJSON(data []byte) error ***REMOVED***
	if err := StrictJSONUnmarshal(data, &c.TLSAuthFields); err != nil ***REMOVED***
		return err
	***REMOVED***
	if _, err := c.Certificate(); err != nil ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

func (c *TLSAuth) Certificate() (*tls.Certificate, error) ***REMOVED***
	key := []byte(c.Key)
	var err error
	if c.Password.Valid ***REMOVED***
		key, err = decryptPrivateKey(c.Key, c.Password.String)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	if c.certificate == nil ***REMOVED***
		cert, err := tls.X509KeyPair([]byte(c.Cert), key)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		c.certificate = &cert
	***REMOVED***
	return c.certificate, nil
***REMOVED***

func decryptPrivateKey(privKey, password string) ([]byte, error) ***REMOVED***
	key := []byte(privKey)

	block, _ := pem.Decode(key)
	if block == nil ***REMOVED***
		return nil, fmt.Errorf("failed to decode PEM key")
	***REMOVED***

	blockType := block.Type
	if blockType == "ENCRYPTED PRIVATE KEY" ***REMOVED***
		return nil, fmt.Errorf("encrypted pkcs8 formatted key is not supported")
	***REMOVED***
	/*
	   Even though `DecryptPEMBlock` has been deprecated since 1.16.x it is still
	   being used here because it is deprecated due to it not supporting *good* crypography
	   ultimately though we want to support something so we will be using it for now.
	*/
	decryptedKey, err := x509.DecryptPEMBlock(block, []byte(password)) //nolint:staticcheck
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	key = pem.EncodeToMemory(&pem.Block***REMOVED***
		Type:  blockType,
		Bytes: decryptedKey,
	***REMOVED***)
	return key, nil
***REMOVED***

// IPNet is a wrapper around net.IPNet for JSON unmarshalling
type IPNet struct ***REMOVED***
	net.IPNet
***REMOVED***

// UnmarshalText populates the IPNet from the given CIDR
func (ipnet *IPNet) UnmarshalText(b []byte) error ***REMOVED***
	newIPNet, err := ParseCIDR(string(b))
	if err != nil ***REMOVED***
		return fmt.Errorf("failed to parse CIDR '%s': %w", string(b), err)
	***REMOVED***

	*ipnet = *newIPNet
	return nil
***REMOVED***

// MarshalText encodes the IPNet representation using CIDR notation.
func (ipnet *IPNet) MarshalText() ([]byte, error) ***REMOVED***
	return []byte(ipnet.String()), nil
***REMOVED***

// HostAddress stores information about IP and port
// for a host.
type HostAddress net.TCPAddr

// NewHostAddress creates a pointer to a new address with an IP object.
func NewHostAddress(ip net.IP, portString string) (*HostAddress, error) ***REMOVED***
	var port int
	if portString != "" ***REMOVED***
		var err error
		if port, err = strconv.Atoi(portString); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	return &HostAddress***REMOVED***
		IP:   ip,
		Port: port,
	***REMOVED***, nil
***REMOVED***

// String converts a HostAddress into a string.
func (h *HostAddress) String() string ***REMOVED***
	return (*net.TCPAddr)(h).String()
***REMOVED***

// MarshalText implements the encoding.TextMarshaler interface.
// The encoding is the same as returned by String, with one exception:
// When len(ip) is zero, it returns an empty slice.
func (h *HostAddress) MarshalText() ([]byte, error) ***REMOVED***
	if h == nil || len(h.IP) == 0 ***REMOVED***
		return []byte(""), nil
	***REMOVED***

	if len(h.IP) != net.IPv4len && len(h.IP) != net.IPv6len ***REMOVED***
		return nil, &net.AddrError***REMOVED***Err: "invalid IP address", Addr: h.IP.String()***REMOVED***
	***REMOVED***

	return []byte(h.String()), nil
***REMOVED***

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// The IP address is expected in a form accepted by ParseIP.
func (h *HostAddress) UnmarshalText(text []byte) error ***REMOVED***
	if len(text) == 0 ***REMOVED***
		return &net.ParseError***REMOVED***Type: "IP address", Text: "<nil>"***REMOVED***
	***REMOVED***

	ip, port, err := splitHostPort(text)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	nh, err := NewHostAddress(ip, port)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	*h = *nh

	return nil
***REMOVED***

func splitHostPort(text []byte) (net.IP, string, error) ***REMOVED***
	host, port, err := net.SplitHostPort(string(text))
	if err != nil ***REMOVED***
		// This error means that there is no port.
		// Make host the full text.
		host = string(text)
	***REMOVED***

	ip := net.ParseIP(host)
	if ip == nil ***REMOVED***
		return nil, "", &net.ParseError***REMOVED***Type: "IP address", Text: host***REMOVED***
	***REMOVED***

	return ip, port, nil
***REMOVED***

// ParseCIDR creates an IPNet out of a CIDR string
func ParseCIDR(s string) (*IPNet, error) ***REMOVED***
	_, ipnet, err := net.ParseCIDR(s)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	parsedIPNet := IPNet***REMOVED***IPNet: *ipnet***REMOVED***

	return &parsedIPNet, nil
***REMOVED***

type Options struct ***REMOVED***
	ExecutionMode    types.NullExecutionMode    `json:"executionMode"`
	LoadDistribution types.NullLoadDistribution `json:"loadDistribution"`

	MaxPossibleVUs null.Int `json:"maxPossibleVUs"`

	// Should the test start in a paused state?
	Paused null.Bool `json:"paused" envconfig:"K6_PAUSED"`

	// Initial values for VUs, max VUs, duration cap, iteration cap, and stages.
	// See the Runner or Executor interfaces for more information.
	VUs        null.Int           `json:"vus" envconfig:"K6_VUS"`
	Duration   types.NullDuration `json:"duration" envconfig:"K6_DURATION"`
	Iterations null.Int           `json:"iterations" envconfig:"K6_ITERATIONS"`
	Stages     []Stage            `json:"stages" envconfig:"K6_STAGES"`

	// TODO: remove the `ignored:"true"` from the field tags, it's there so that
	// the envconfig library will ignore those fields.
	//
	// We should support specifying execution segments via environment
	// variables, but we currently can't, because envconfig has this nasty bug
	// (among others): https://github.com/kelseyhightower/envconfig/issues/113
	Scenarios                ScenarioConfigs           `json:"scenarios" ignored:"true"`
	ExecutionSegment         *ExecutionSegment         `json:"executionSegment" ignored:"true"`
	ExecutionSegmentSequence *ExecutionSegmentSequence `json:"executionSegmentSequence" ignored:"true"`

	// Timeouts for the setup() and teardown() functions
	NoSetup         null.Bool          `json:"noSetup" envconfig:"K6_NO_SETUP"`
	SetupTimeout    types.NullDuration `json:"setupTimeout" envconfig:"K6_SETUP_TIMEOUT"`
	NoTeardown      null.Bool          `json:"noTeardown" envconfig:"K6_NO_TEARDOWN"`
	TeardownTimeout types.NullDuration `json:"teardownTimeout" envconfig:"K6_TEARDOWN_TIMEOUT"`

	// Limit HTTP requests per second.
	RPS null.Int `json:"rps" envconfig:"K6_RPS"`

	// DNS handling configuration.
	DNS types.DNSConfig `json:"dns" envconfig:"K6_DNS"`

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
	TLSVersion      *TLSVersions     `json:"tlsVersion" ignored:"true"`
	TLSAuth         []*TLSAuth       `json:"tlsAuth" envconfig:"K6_TLSAUTH"`

	// Throw warnings (eg. failed HTTP requests) as errors instead of simply logging them.
	Throw null.Bool `json:"throw" envconfig:"K6_THROW"`

	// Define thresholds; these take the form of 'metric=["snippet1", "snippet2"]'.
	// To create a threshold on a derived metric based on tag queries ("submetrics"), create a
	// metric on a nonexistent metric named 'real_metric***REMOVED***tagA:valueA,tagB:valueB***REMOVED***'.
	Thresholds map[string]workerMetrics.Thresholds `json:"thresholds" envconfig:"K6_THRESHOLDS"`

	// Blacklist IP ranges that tests may not contact. Mainly useful in hosted setups.
	BlacklistIPs []*IPNet `json:"blacklistIPs" envconfig:"K6_BLACKLIST_IPS"`

	// Block hostname patterns that tests may not contact.
	BlockedHostnames types.NullHostnameTrie `json:"blockHostnames" envconfig:"K6_BLOCK_HOSTNAMES"`

	// Hosts overrides dns entries for given hosts
	Hosts map[string]*HostAddress `json:"hosts" envconfig:"K6_HOSTS"`

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
	SystemTags *workerMetrics.SystemTagSet `json:"systemTags" envconfig:"K6_SYSTEM_TAGS"`

	// Tags are key-value pairs to be applied to all samples for the run.
	RunTags map[string]string `json:"tags" envconfig:"K6_TAGS"`

	// Buffer size of the channel for metric samples; 0 means unbuffered
	MetricSamplesBufferSize null.Int `json:"metricSamplesBufferSize" envconfig:"K6_METRIC_SAMPLES_BUFFER_SIZE"`

	// Do not reset cookies after a VU iteration
	NoCookiesReset null.Bool `json:"noCookiesReset" envconfig:"K6_NO_COOKIES_RESET"`

	// Discard Http Responses Body
	DiscardResponseBodies null.Bool `json:"discardResponseBodies" envconfig:"K6_DISCARD_RESPONSE_BODIES"`
***REMOVED***

// Returns the result of overwriting any fields with any that are set on the argument.
//
// Example:
//   a := Options***REMOVED***VUs: null.IntFrom(10)***REMOVED***
//   b := Options***REMOVED***VUs: null.IntFrom(5)***REMOVED***
//   a.Apply(b) // Options***REMOVED***VUs: null.IntFrom(5)***REMOVED***
func (o Options) Apply(opts Options) Options ***REMOVED***
	if opts.Paused.Valid ***REMOVED***
		o.Paused = opts.Paused
	***REMOVED***
	if opts.VUs.Valid ***REMOVED***
		o.VUs = opts.VUs
	***REMOVED***

	// Specifying duration, iterations, stages, or execution in a "higher" config tier
	// will overwrite all of the the previous execution settings (if any) from any
	// "lower" config tiers
	// Still, if more than one of those options is simultaneously specified in the same
	// config tier, they will be preserved, so the validation after we've consolidated
	// all of the options can return an error.
	if opts.Duration.Valid || opts.Iterations.Valid || opts.Stages != nil || opts.Scenarios != nil ***REMOVED***
		// TODO: emit a warning or a notice log message if overwrite lower tier config options?
		o.Duration = types.NewNullDuration(0, false)
		o.Iterations = null.NewInt(0, false)
		o.Stages = nil
		o.Scenarios = nil
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
	if opts.Scenarios != nil ***REMOVED***
		o.Scenarios = opts.Scenarios
	***REMOVED***
	if opts.ExecutionSegment != nil ***REMOVED***
		o.ExecutionSegment = opts.ExecutionSegment
	***REMOVED***

	if opts.ExecutionSegmentSequence != nil ***REMOVED***
		o.ExecutionSegmentSequence = opts.ExecutionSegmentSequence
	***REMOVED***
	if opts.NoSetup.Valid ***REMOVED***
		o.NoSetup = opts.NoSetup
	***REMOVED***
	if opts.SetupTimeout.Valid ***REMOVED***
		o.SetupTimeout = opts.SetupTimeout
	***REMOVED***
	if opts.NoTeardown.Valid ***REMOVED***
		o.NoTeardown = opts.NoTeardown
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
	if opts.BlockedHostnames.Valid ***REMOVED***
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
	if len(opts.RunTags) > 0 ***REMOVED***
		o.RunTags = opts.RunTags
	***REMOVED***
	if opts.MetricSamplesBufferSize.Valid ***REMOVED***
		o.MetricSamplesBufferSize = opts.MetricSamplesBufferSize
	***REMOVED***
	if opts.DiscardResponseBodies.Valid ***REMOVED***
		o.DiscardResponseBodies = opts.DiscardResponseBodies
	***REMOVED***
	if opts.DNS.TTL.Valid ***REMOVED***
		o.DNS.TTL = opts.DNS.TTL
	***REMOVED***
	if opts.DNS.Select.Valid ***REMOVED***
		o.DNS.Select = opts.DNS.Select
	***REMOVED***
	if opts.DNS.Policy.Valid ***REMOVED***
		o.DNS.Policy = opts.DNS.Policy
	***REMOVED***

	return o
***REMOVED***

// Validate checks if all of the specified options make sense
func (o Options) Validate() []error ***REMOVED***
	// TODO: validate all of the other options... that we should have already been validating...
	// TODO: maybe integrate an external validation lib: https://github.com/avelino/awesome-go#validation
	var errors []error
	if o.ExecutionSegmentSequence != nil ***REMOVED***
		var segmentFound bool
		for _, segment := range *o.ExecutionSegmentSequence ***REMOVED***
			if o.ExecutionSegment.Equal(segment) ***REMOVED***
				segmentFound = true
				break
			***REMOVED***
		***REMOVED***
		if !segmentFound ***REMOVED***
			errors = append(errors,
				fmt.Errorf("provided segment %s can't be found in sequence %s",
					o.ExecutionSegment, o.ExecutionSegmentSequence))
		***REMOVED***
	***REMOVED***
	return append(errors, o.Scenarios.Validate()...)
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

func (o Options) Clone() Options ***REMOVED***
	return o
***REMOVED***
