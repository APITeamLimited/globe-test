// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package dns

import (
	"errors"
	"fmt"
	"net"
	"runtime"
	"strings"
)

// Resolver resolves DNS records.
type Resolver struct ***REMOVED***
	// Holds the functions to use for DNS lookups
	LookupSRV func(string, string, string) (string, []*net.SRV, error)
	LookupTXT func(string) ([]string, error)
***REMOVED***

// DefaultResolver is a Resolver that uses the default Resolver from the net package.
var DefaultResolver = &Resolver***REMOVED***net.LookupSRV, net.LookupTXT***REMOVED***

// ParseHosts uses the srv string and service name to get the hosts.
func (r *Resolver) ParseHosts(host string, srvName string, stopOnErr bool) ([]string, error) ***REMOVED***
	parsedHosts := strings.Split(host, ",")

	if len(parsedHosts) != 1 ***REMOVED***
		return nil, fmt.Errorf("URI with SRV must include one and only one hostname")
	***REMOVED***
	return r.fetchSeedlistFromSRV(parsedHosts[0], srvName, stopOnErr)
***REMOVED***

// GetConnectionArgsFromTXT gets the TXT record associated with the host and returns the connection arguments.
func (r *Resolver) GetConnectionArgsFromTXT(host string) ([]string, error) ***REMOVED***
	var connectionArgsFromTXT []string

	// error ignored because not finding a TXT record should not be
	// considered an error.
	recordsFromTXT, _ := r.LookupTXT(host)

	// This is a temporary fix to get around bug https://github.com/golang/go/issues/21472.
	// It will currently incorrectly concatenate multiple TXT records to one
	// on windows.
	if runtime.GOOS == "windows" ***REMOVED***
		recordsFromTXT = []string***REMOVED***strings.Join(recordsFromTXT, "")***REMOVED***
	***REMOVED***

	if len(recordsFromTXT) > 1 ***REMOVED***
		return nil, errors.New("multiple records from TXT not supported")
	***REMOVED***
	if len(recordsFromTXT) > 0 ***REMOVED***
		connectionArgsFromTXT = strings.FieldsFunc(recordsFromTXT[0], func(r rune) bool ***REMOVED*** return r == ';' || r == '&' ***REMOVED***)

		err := validateTXTResult(connectionArgsFromTXT)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	return connectionArgsFromTXT, nil
***REMOVED***

func (r *Resolver) fetchSeedlistFromSRV(host string, srvName string, stopOnErr bool) ([]string, error) ***REMOVED***
	var err error

	_, _, err = net.SplitHostPort(host)

	if err == nil ***REMOVED***
		// we were able to successfully extract a port from the host,
		// but should not be able to when using SRV
		return nil, fmt.Errorf("URI with srv must not include a port number")
	***REMOVED***

	// default to "mongodb" as service name if not supplied
	if srvName == "" ***REMOVED***
		srvName = "mongodb"
	***REMOVED***
	_, addresses, err := r.LookupSRV(srvName, "tcp", host)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	trimmedHost := strings.TrimSuffix(host, ".")

	parsedHosts := make([]string, 0, len(addresses))
	for _, address := range addresses ***REMOVED***
		trimmedAddressTarget := strings.TrimSuffix(address.Target, ".")
		err := validateSRVResult(trimmedAddressTarget, trimmedHost)
		if err != nil ***REMOVED***
			if stopOnErr ***REMOVED***
				return nil, err
			***REMOVED***
			continue
		***REMOVED***
		parsedHosts = append(parsedHosts, fmt.Sprintf("%s:%d", trimmedAddressTarget, address.Port))
	***REMOVED***
	return parsedHosts, nil
***REMOVED***

func validateSRVResult(recordFromSRV, inputHostName string) error ***REMOVED***
	separatedInputDomain := strings.Split(inputHostName, ".")
	separatedRecord := strings.Split(recordFromSRV, ".")
	if len(separatedRecord) < 2 ***REMOVED***
		return errors.New("DNS name must contain at least 2 labels")
	***REMOVED***
	if len(separatedRecord) < len(separatedInputDomain) ***REMOVED***
		return errors.New("Domain suffix from SRV record not matched input domain")
	***REMOVED***

	inputDomainSuffix := separatedInputDomain[1:]
	domainSuffixOffset := len(separatedRecord) - (len(separatedInputDomain) - 1)

	recordDomainSuffix := separatedRecord[domainSuffixOffset:]
	for ix, label := range inputDomainSuffix ***REMOVED***
		if label != recordDomainSuffix[ix] ***REMOVED***
			return errors.New("Domain suffix from SRV record not matched input domain")
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

var allowedTXTOptions = map[string]struct***REMOVED******REMOVED******REMOVED***
	"authsource":   ***REMOVED******REMOVED***,
	"replicaset":   ***REMOVED******REMOVED***,
	"loadbalanced": ***REMOVED******REMOVED***,
***REMOVED***

func validateTXTResult(paramsFromTXT []string) error ***REMOVED***
	for _, param := range paramsFromTXT ***REMOVED***
		kv := strings.SplitN(param, "=", 2)
		if len(kv) != 2 ***REMOVED***
			return errors.New("Invalid TXT record")
		***REMOVED***
		key := strings.ToLower(kv[0])
		if _, ok := allowedTXTOptions[key]; !ok ***REMOVED***
			return fmt.Errorf("Cannot specify option '%s' in TXT record", kv[0])
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***
