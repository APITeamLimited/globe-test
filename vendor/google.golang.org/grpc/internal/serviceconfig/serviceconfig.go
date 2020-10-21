/*
 *
 * Copyright 2020 gRPC authors.
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

// Package serviceconfig contains utility functions to parse service config.
package serviceconfig

import (
	"encoding/json"
	"fmt"

	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/grpclog"
	externalserviceconfig "google.golang.org/grpc/serviceconfig"
)

// BalancerConfig is the balancer config part that service config's
// loadBalancingConfig fields can be unmarshalled to. It's a json unmarshaller.
//
// https://github.com/grpc/grpc-proto/blob/54713b1e8bc6ed2d4f25fb4dff527842150b91b2/grpc/service_config/service_config.proto#L247
type BalancerConfig struct ***REMOVED***
	Name   string
	Config externalserviceconfig.LoadBalancingConfig
***REMOVED***

type intermediateBalancerConfig []map[string]json.RawMessage

// UnmarshalJSON implements json unmarshaller.
func (bc *BalancerConfig) UnmarshalJSON(b []byte) error ***REMOVED***
	var ir intermediateBalancerConfig
	err := json.Unmarshal(b, &ir)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	for i, lbcfg := range ir ***REMOVED***
		if len(lbcfg) != 1 ***REMOVED***
			return fmt.Errorf("invalid loadBalancingConfig: entry %v does not contain exactly 1 policy/config pair: %q", i, lbcfg)
		***REMOVED***
		var (
			name    string
			jsonCfg json.RawMessage
		)
		// Get the key:value pair from the map.
		for name, jsonCfg = range lbcfg ***REMOVED***
		***REMOVED***
		builder := balancer.Get(name)
		if builder == nil ***REMOVED***
			// If the balancer is not registered, move on to the next config.
			// This is not an error.
			continue
		***REMOVED***
		bc.Name = name

		parser, ok := builder.(balancer.ConfigParser)
		if !ok ***REMOVED***
			if string(jsonCfg) != "***REMOVED******REMOVED***" ***REMOVED***
				grpclog.Warningf("non-empty balancer configuration %q, but balancer does not implement ParseConfig", string(jsonCfg))
			***REMOVED***
			// Stop at this, though the builder doesn't support parsing config.
			return nil
		***REMOVED***

		cfg, err := parser.ParseConfig(jsonCfg)
		if err != nil ***REMOVED***
			return fmt.Errorf("error parsing loadBalancingConfig for policy %q: %v", name, err)
		***REMOVED***
		bc.Config = cfg
		return nil
	***REMOVED***
	// This is reached when the for loop iterates over all entries, but didn't
	// return. This means we had a loadBalancingConfig slice but did not
	// encounter a registered policy. The config is considered invalid in this
	// case.
	return fmt.Errorf("invalid loadBalancingConfig: no supported policies found")
***REMOVED***
