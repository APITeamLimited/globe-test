/*
 *
 * Copyright 2017 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package grpc

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/internal"
	internalserviceconfig "google.golang.org/grpc/internal/serviceconfig"
	"google.golang.org/grpc/serviceconfig"
)

const maxInt = int(^uint(0) >> 1)

// MethodConfig defines the configuration recommended by the service providers for a
// particular method.
//
// Deprecated: Users should not use this struct. Service config should be received
// through name resolver, as specified here
// https://github.com/grpc/grpc/blob/master/doc/service_config.md
type MethodConfig = internalserviceconfig.MethodConfig

type lbConfig struct ***REMOVED***
	name string
	cfg  serviceconfig.LoadBalancingConfig
***REMOVED***

// ServiceConfig is provided by the service provider and contains parameters for how
// clients that connect to the service should behave.
//
// Deprecated: Users should not use this struct. Service config should be received
// through name resolver, as specified here
// https://github.com/grpc/grpc/blob/master/doc/service_config.md
type ServiceConfig struct ***REMOVED***
	serviceconfig.Config

	// LB is the load balancer the service providers recommends. The balancer
	// specified via grpc.WithBalancerName will override this.  This is deprecated;
	// lbConfigs is preferred.  If lbConfig and LB are both present, lbConfig
	// will be used.
	LB *string

	// lbConfig is the service config's load balancing configuration.  If
	// lbConfig and LB are both present, lbConfig will be used.
	lbConfig *lbConfig

	// Methods contains a map for the methods in this service.  If there is an
	// exact match for a method (i.e. /service/method) in the map, use the
	// corresponding MethodConfig.  If there's no exact match, look for the
	// default config for the service (/service/) and use the corresponding
	// MethodConfig if it exists.  Otherwise, the method has no MethodConfig to
	// use.
	Methods map[string]MethodConfig

	// If a retryThrottlingPolicy is provided, gRPC will automatically throttle
	// retry attempts and hedged RPCs when the client’s ratio of failures to
	// successes exceeds a threshold.
	//
	// For each server name, the gRPC client will maintain a token_count which is
	// initially set to maxTokens, and can take values between 0 and maxTokens.
	//
	// Every outgoing RPC (regardless of service or method invoked) will change
	// token_count as follows:
	//
	//   - Every failed RPC will decrement the token_count by 1.
	//   - Every successful RPC will increment the token_count by tokenRatio.
	//
	// If token_count is less than or equal to maxTokens / 2, then RPCs will not
	// be retried and hedged RPCs will not be sent.
	retryThrottling *retryThrottlingPolicy
	// healthCheckConfig must be set as one of the requirement to enable LB channel
	// health check.
	healthCheckConfig *healthCheckConfig
	// rawJSONString stores service config json string that get parsed into
	// this service config struct.
	rawJSONString string
***REMOVED***

// healthCheckConfig defines the go-native version of the LB channel health check config.
type healthCheckConfig struct ***REMOVED***
	// serviceName is the service name to use in the health-checking request.
	ServiceName string
***REMOVED***

type jsonRetryPolicy struct ***REMOVED***
	MaxAttempts          int
	InitialBackoff       string
	MaxBackoff           string
	BackoffMultiplier    float64
	RetryableStatusCodes []codes.Code
***REMOVED***

// retryThrottlingPolicy defines the go-native version of the retry throttling
// policy defined by the service config here:
// https://github.com/grpc/proposal/blob/master/A6-client-retries.md#integration-with-service-config
type retryThrottlingPolicy struct ***REMOVED***
	// The number of tokens starts at maxTokens. The token_count will always be
	// between 0 and maxTokens.
	//
	// This field is required and must be greater than zero.
	MaxTokens float64
	// The amount of tokens to add on each successful RPC. Typically this will
	// be some number between 0 and 1, e.g., 0.1.
	//
	// This field is required and must be greater than zero. Up to 3 decimal
	// places are supported.
	TokenRatio float64
***REMOVED***

func parseDuration(s *string) (*time.Duration, error) ***REMOVED***
	if s == nil ***REMOVED***
		return nil, nil
	***REMOVED***
	if !strings.HasSuffix(*s, "s") ***REMOVED***
		return nil, fmt.Errorf("malformed duration %q", *s)
	***REMOVED***
	ss := strings.SplitN((*s)[:len(*s)-1], ".", 3)
	if len(ss) > 2 ***REMOVED***
		return nil, fmt.Errorf("malformed duration %q", *s)
	***REMOVED***
	// hasDigits is set if either the whole or fractional part of the number is
	// present, since both are optional but one is required.
	hasDigits := false
	var d time.Duration
	if len(ss[0]) > 0 ***REMOVED***
		i, err := strconv.ParseInt(ss[0], 10, 32)
		if err != nil ***REMOVED***
			return nil, fmt.Errorf("malformed duration %q: %v", *s, err)
		***REMOVED***
		d = time.Duration(i) * time.Second
		hasDigits = true
	***REMOVED***
	if len(ss) == 2 && len(ss[1]) > 0 ***REMOVED***
		if len(ss[1]) > 9 ***REMOVED***
			return nil, fmt.Errorf("malformed duration %q", *s)
		***REMOVED***
		f, err := strconv.ParseInt(ss[1], 10, 64)
		if err != nil ***REMOVED***
			return nil, fmt.Errorf("malformed duration %q: %v", *s, err)
		***REMOVED***
		for i := 9; i > len(ss[1]); i-- ***REMOVED***
			f *= 10
		***REMOVED***
		d += time.Duration(f)
		hasDigits = true
	***REMOVED***
	if !hasDigits ***REMOVED***
		return nil, fmt.Errorf("malformed duration %q", *s)
	***REMOVED***

	return &d, nil
***REMOVED***

type jsonName struct ***REMOVED***
	Service string
	Method  string
***REMOVED***

var (
	errDuplicatedName             = errors.New("duplicated name")
	errEmptyServiceNonEmptyMethod = errors.New("cannot combine empty 'service' and non-empty 'method'")
)

func (j jsonName) generatePath() (string, error) ***REMOVED***
	if j.Service == "" ***REMOVED***
		if j.Method != "" ***REMOVED***
			return "", errEmptyServiceNonEmptyMethod
		***REMOVED***
		return "", nil
	***REMOVED***
	res := "/" + j.Service + "/"
	if j.Method != "" ***REMOVED***
		res += j.Method
	***REMOVED***
	return res, nil
***REMOVED***

// TODO(lyuxuan): delete this struct after cleaning up old service config implementation.
type jsonMC struct ***REMOVED***
	Name                    *[]jsonName
	WaitForReady            *bool
	Timeout                 *string
	MaxRequestMessageBytes  *int64
	MaxResponseMessageBytes *int64
	RetryPolicy             *jsonRetryPolicy
***REMOVED***

// TODO(lyuxuan): delete this struct after cleaning up old service config implementation.
type jsonSC struct ***REMOVED***
	LoadBalancingPolicy *string
	LoadBalancingConfig *internalserviceconfig.BalancerConfig
	MethodConfig        *[]jsonMC
	RetryThrottling     *retryThrottlingPolicy
	HealthCheckConfig   *healthCheckConfig
***REMOVED***

func init() ***REMOVED***
	internal.ParseServiceConfig = parseServiceConfig
***REMOVED***
func parseServiceConfig(js string) *serviceconfig.ParseResult ***REMOVED***
	if len(js) == 0 ***REMOVED***
		return &serviceconfig.ParseResult***REMOVED***Err: fmt.Errorf("no JSON service config provided")***REMOVED***
	***REMOVED***
	var rsc jsonSC
	err := json.Unmarshal([]byte(js), &rsc)
	if err != nil ***REMOVED***
		logger.Warningf("grpc: parseServiceConfig error unmarshaling %s due to %v", js, err)
		return &serviceconfig.ParseResult***REMOVED***Err: err***REMOVED***
	***REMOVED***
	sc := ServiceConfig***REMOVED***
		LB:                rsc.LoadBalancingPolicy,
		Methods:           make(map[string]MethodConfig),
		retryThrottling:   rsc.RetryThrottling,
		healthCheckConfig: rsc.HealthCheckConfig,
		rawJSONString:     js,
	***REMOVED***
	if c := rsc.LoadBalancingConfig; c != nil ***REMOVED***
		sc.lbConfig = &lbConfig***REMOVED***
			name: c.Name,
			cfg:  c.Config,
		***REMOVED***
	***REMOVED***

	if rsc.MethodConfig == nil ***REMOVED***
		return &serviceconfig.ParseResult***REMOVED***Config: &sc***REMOVED***
	***REMOVED***

	paths := map[string]struct***REMOVED******REMOVED******REMOVED******REMOVED***
	for _, m := range *rsc.MethodConfig ***REMOVED***
		if m.Name == nil ***REMOVED***
			continue
		***REMOVED***
		d, err := parseDuration(m.Timeout)
		if err != nil ***REMOVED***
			logger.Warningf("grpc: parseServiceConfig error unmarshaling %s due to %v", js, err)
			return &serviceconfig.ParseResult***REMOVED***Err: err***REMOVED***
		***REMOVED***

		mc := MethodConfig***REMOVED***
			WaitForReady: m.WaitForReady,
			Timeout:      d,
		***REMOVED***
		if mc.RetryPolicy, err = convertRetryPolicy(m.RetryPolicy); err != nil ***REMOVED***
			logger.Warningf("grpc: parseServiceConfig error unmarshaling %s due to %v", js, err)
			return &serviceconfig.ParseResult***REMOVED***Err: err***REMOVED***
		***REMOVED***
		if m.MaxRequestMessageBytes != nil ***REMOVED***
			if *m.MaxRequestMessageBytes > int64(maxInt) ***REMOVED***
				mc.MaxReqSize = newInt(maxInt)
			***REMOVED*** else ***REMOVED***
				mc.MaxReqSize = newInt(int(*m.MaxRequestMessageBytes))
			***REMOVED***
		***REMOVED***
		if m.MaxResponseMessageBytes != nil ***REMOVED***
			if *m.MaxResponseMessageBytes > int64(maxInt) ***REMOVED***
				mc.MaxRespSize = newInt(maxInt)
			***REMOVED*** else ***REMOVED***
				mc.MaxRespSize = newInt(int(*m.MaxResponseMessageBytes))
			***REMOVED***
		***REMOVED***
		for i, n := range *m.Name ***REMOVED***
			path, err := n.generatePath()
			if err != nil ***REMOVED***
				logger.Warningf("grpc: parseServiceConfig error unmarshaling %s due to methodConfig[%d]: %v", js, i, err)
				return &serviceconfig.ParseResult***REMOVED***Err: err***REMOVED***
			***REMOVED***

			if _, ok := paths[path]; ok ***REMOVED***
				err = errDuplicatedName
				logger.Warningf("grpc: parseServiceConfig error unmarshaling %s due to methodConfig[%d]: %v", js, i, err)
				return &serviceconfig.ParseResult***REMOVED***Err: err***REMOVED***
			***REMOVED***
			paths[path] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
			sc.Methods[path] = mc
		***REMOVED***
	***REMOVED***

	if sc.retryThrottling != nil ***REMOVED***
		if mt := sc.retryThrottling.MaxTokens; mt <= 0 || mt > 1000 ***REMOVED***
			return &serviceconfig.ParseResult***REMOVED***Err: fmt.Errorf("invalid retry throttling config: maxTokens (%v) out of range (0, 1000]", mt)***REMOVED***
		***REMOVED***
		if tr := sc.retryThrottling.TokenRatio; tr <= 0 ***REMOVED***
			return &serviceconfig.ParseResult***REMOVED***Err: fmt.Errorf("invalid retry throttling config: tokenRatio (%v) may not be negative", tr)***REMOVED***
		***REMOVED***
	***REMOVED***
	return &serviceconfig.ParseResult***REMOVED***Config: &sc***REMOVED***
***REMOVED***

func convertRetryPolicy(jrp *jsonRetryPolicy) (p *internalserviceconfig.RetryPolicy, err error) ***REMOVED***
	if jrp == nil ***REMOVED***
		return nil, nil
	***REMOVED***
	ib, err := parseDuration(&jrp.InitialBackoff)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	mb, err := parseDuration(&jrp.MaxBackoff)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if jrp.MaxAttempts <= 1 ||
		*ib <= 0 ||
		*mb <= 0 ||
		jrp.BackoffMultiplier <= 0 ||
		len(jrp.RetryableStatusCodes) == 0 ***REMOVED***
		logger.Warningf("grpc: ignoring retry policy %v due to illegal configuration", jrp)
		return nil, nil
	***REMOVED***

	rp := &internalserviceconfig.RetryPolicy***REMOVED***
		MaxAttempts:          jrp.MaxAttempts,
		InitialBackoff:       *ib,
		MaxBackoff:           *mb,
		BackoffMultiplier:    jrp.BackoffMultiplier,
		RetryableStatusCodes: make(map[codes.Code]bool),
	***REMOVED***
	if rp.MaxAttempts > 5 ***REMOVED***
		// TODO(retry): Make the max maxAttempts configurable.
		rp.MaxAttempts = 5
	***REMOVED***
	for _, code := range jrp.RetryableStatusCodes ***REMOVED***
		rp.RetryableStatusCodes[code] = true
	***REMOVED***
	return rp, nil
***REMOVED***

func min(a, b *int) *int ***REMOVED***
	if *a < *b ***REMOVED***
		return a
	***REMOVED***
	return b
***REMOVED***

func getMaxSize(mcMax, doptMax *int, defaultVal int) *int ***REMOVED***
	if mcMax == nil && doptMax == nil ***REMOVED***
		return &defaultVal
	***REMOVED***
	if mcMax != nil && doptMax != nil ***REMOVED***
		return min(mcMax, doptMax)
	***REMOVED***
	if mcMax != nil ***REMOVED***
		return mcMax
	***REMOVED***
	return doptMax
***REMOVED***

func newInt(b int) *int ***REMOVED***
	return &b
***REMOVED***

func init() ***REMOVED***
	internal.EqualServiceConfigForTesting = equalServiceConfig
***REMOVED***

// equalServiceConfig compares two configs. The rawJSONString field is ignored,
// because they may diff in white spaces.
//
// If any of them is NOT *ServiceConfig, return false.
func equalServiceConfig(a, b serviceconfig.Config) bool ***REMOVED***
	if a == nil && b == nil ***REMOVED***
		return true
	***REMOVED***
	aa, ok := a.(*ServiceConfig)
	if !ok ***REMOVED***
		return false
	***REMOVED***
	bb, ok := b.(*ServiceConfig)
	if !ok ***REMOVED***
		return false
	***REMOVED***
	aaRaw := aa.rawJSONString
	aa.rawJSONString = ""
	bbRaw := bb.rawJSONString
	bb.rawJSONString = ""
	defer func() ***REMOVED***
		aa.rawJSONString = aaRaw
		bb.rawJSONString = bbRaw
	***REMOVED***()
	// Using reflect.DeepEqual instead of cmp.Equal because many balancer
	// configs are unexported, and cmp.Equal cannot compare unexported fields
	// from unexported structs.
	return reflect.DeepEqual(aa, bb)
***REMOVED***
