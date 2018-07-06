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
	"bytes"
	"crypto/tls"
	"encoding/json"
	"net"

	"github.com/loadimpact/k6/lib/types"
	"github.com/loadimpact/k6/stats"
	"github.com/pkg/errors"
	"gopkg.in/guregu/null.v3"
)

// DefaultSystemTagList includes all of the system tags emitted with metrics by default.
// Other tags that are not enabled by default include: iter, vu, ocsp_status, ip
var DefaultSystemTagList = []string***REMOVED***
	"proto", "subproto", "status", "method", "url", "name", "group", "check", "error", "tls_version",
***REMOVED***

// TagSet is a string to bool map (for lookup efficiency) that is used to keep track
// which system tags should be included with with metrics.
type TagSet map[string]bool

// GetTagSet converts a the passed string tag names into the expected string to bool map.
func GetTagSet(tags ...string) TagSet ***REMOVED***
	result := TagSet***REMOVED******REMOVED***
	for _, tag := range tags ***REMOVED***
		result[tag] = true
	***REMOVED***
	return result
***REMOVED***

// MarshalJSON converts the tags map to a list (JS array).
func (t TagSet) MarshalJSON() ([]byte, error) ***REMOVED***
	var tags []string
	for tag := range t ***REMOVED***
		tags = append(tags, tag)
	***REMOVED***
	return json.Marshal(tags)
***REMOVED***

// UnmarshalJSON converts the tag list back to a the expected set (string to bool map).
func (t *TagSet) UnmarshalJSON(data []byte) error ***REMOVED***
	var tags []string
	if err := json.Unmarshal(data, &tags); err != nil ***REMOVED***
		return err
	***REMOVED***
	if len(tags) != 0 ***REMOVED***
		*t = GetTagSet(tags...)
	***REMOVED***
	return nil
***REMOVED***

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

type Options struct ***REMOVED***
	// Should the test start in a paused state?
	Paused null.Bool `json:"paused" envconfig:"paused"`

	// Initial values for VUs, max VUs, duration cap, iteration cap, and stages.
	// See the Runner or Executor interfaces for more information.
	VUs        null.Int           `json:"vus" envconfig:"vus"`
	VUsMax     null.Int           `json:"vusMax" envconfig:"vus_max"`
	Duration   types.NullDuration `json:"duration" envconfig:"duration"`
	Iterations null.Int           `json:"iterations" envconfig:"iterations"`
	Stages     []Stage            `json:"stages" envconfig:"stages"`

	// Timeouts for the setup() and teardown() functions
	SetupTimeout    types.NullDuration `json:"setupTimeout" envconfig:"setup_timeout"`
	TeardownTimeout types.NullDuration `json:"teardownTimeout" envconfig:"teardown_timeout"`

	// Limit HTTP requests per second.
	RPS null.Int `json:"rps" envconfig:"rps"`

	// How many HTTP redirects do we follow?
	MaxRedirects null.Int `json:"maxRedirects" envconfig:"max_redirects"`

	// Default User Agent string for HTTP requests.
	UserAgent null.String `json:"userAgent" envconfig:"user_agent"`

	// How many batch requests are allowed in parallel, in total and per host?
	Batch        null.Int `json:"batch" envconfig:"batch"`
	BatchPerHost null.Int `json:"batchPerHost" envconfig:"batch_per_host"`

	// Should all HTTP requests and responses be logged (excluding body)?
	HttpDebug null.String `json:"httpDebug" envconfig:"http_debug"`

	// Accept invalid or untrusted TLS certificates.
	InsecureSkipTLSVerify null.Bool `json:"insecureSkipTLSVerify" envconfig:"insecure_skip_tls_verify"`

	// Specify TLS versions and cipher suites, and present client certificates.
	TLSCipherSuites *TLSCipherSuites `json:"tlsCipherSuites" envconfig:"tls_cipher_suites"`
	TLSVersion      *TLSVersions     `json:"tlsVersion" envconfig:"tls_version"`
	TLSAuth         []*TLSAuth       `json:"tlsAuth" envconfig:"tlsauth"`

	// Throw warnings (eg. failed HTTP requests) as errors instead of simply logging them.
	Throw null.Bool `json:"throw" envconfig:"throw"`

	// Define thresholds; these take the form of 'metric=["snippet1", "snippet2"]'.
	// To create a threshold on a derived metric based on tag queries ("submetrics"), create a
	// metric on a nonexistent metric named 'real_metric***REMOVED***tagA:valueA,tagB:valueB***REMOVED***'.
	Thresholds map[string]stats.Thresholds `json:"thresholds" envconfig:"thresholds"`

	// Blacklist IP ranges that tests may not contact. Mainly useful in hosted setups.
	BlacklistIPs []*net.IPNet `json:"blacklistIPs" envconfig:"blacklist_ips"`

	// Hosts overrides dns entries for given hosts
	Hosts map[string]net.IP `json:"hosts" envconfig:"hosts"`

	// Disable keep-alive connections
	NoConnectionReuse null.Bool `json:"noConnectionReuse" envconfig:"no_connection_reuse"`

	// Do not reuse connections between VU iterations. This gives more realistic results (depending
	// on what you're looking for), but you need to raise various kernel limits or you'll get
	// errors about running out of file handles or sockets, or being unable to bind addresses.
	NoVUConnectionReuse null.Bool `json:"noVUConnectionReuse" envconfig:"no_vu_connection_reuse"`

	// These values are for third party collectors' benefit.
	// Can't be set through env vars.
	External map[string]json.RawMessage `json:"ext" ignored:"true"`

	// Summary trend stats for trend metrics (response times) in CLI output
	SummaryTrendStats []string `json:"summaryTrendStats" envconfig:"summary_trend_stats"`

	// Summary time unit for summary metrics (response times) in CLI output
	SummaryTimeUnit null.String `json:"summaryTimeUnit" envconfig:"summary_time_unit"`

	// Which system tags to include with metrics ("method", "vu" etc.)
	SystemTags TagSet `json:"systemTags" envconfig:"system_tags"`

	// Tags to be applied to all samples for this running
	RunTags *stats.SampleTags `json:"tags" envconfig:"tags"`

	// Buffer size of the channel for metric samples; 0 means unbuffered
	MetricSamplesBufferSize null.Int `json:"metricSamplesBufferSize" envconfig:"metric_samples_buffer_size"`
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
	if opts.Duration.Valid ***REMOVED***
		o.Duration = opts.Duration
	***REMOVED***
	if opts.Iterations.Valid ***REMOVED***
		o.Iterations = opts.Iterations
	***REMOVED***
	if len(opts.Stages) > 0 ***REMOVED***
		for _, s := range opts.Stages ***REMOVED***
			if s.Duration.Valid ***REMOVED***
				o.Stages = append(o.Stages, s)
			***REMOVED***
		***REMOVED***
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
	if opts.HttpDebug.Valid ***REMOVED***
		o.HttpDebug = opts.HttpDebug
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
	if opts.Hosts != nil ***REMOVED***
		o.Hosts = opts.Hosts
	***REMOVED***
	if opts.NoConnectionReuse.Valid ***REMOVED***
		o.NoConnectionReuse = opts.NoConnectionReuse
	***REMOVED***
	if opts.NoVUConnectionReuse.Valid ***REMOVED***
		o.NoVUConnectionReuse = opts.NoVUConnectionReuse
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
	return o
***REMOVED***

// GetCleanJSON is a massive hack that works arround the fact that some
// of the null-able types used in Options are marshalled to `null` when
// their `valid` flag is false.
func (o Options) GetCleanJSON() ([]byte, error) ***REMOVED***
	nullyResult, err := json.Marshal(o)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var tmpMap map[string]json.RawMessage
	if err := json.Unmarshal(nullyResult, &tmpMap); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	null := []byte("null")
	for k, v := range tmpMap ***REMOVED***
		if bytes.Equal(v, null) ***REMOVED***
			delete(tmpMap, k)
		***REMOVED***
	***REMOVED***
	return json.Marshal(tmpMap)
***REMOVED***
