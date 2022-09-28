// Copyright 2018 by David A. Golden. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package scram

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type c1Msg struct ***REMOVED***
	gs2Header string
	authzID   string
	username  string
	nonce     string
	c1b       string
***REMOVED***

type c2Msg struct ***REMOVED***
	cbind []byte
	nonce string
	proof []byte
	c2wop string
***REMOVED***

type s1Msg struct ***REMOVED***
	nonce string
	salt  []byte
	iters int
***REMOVED***

type s2Msg struct ***REMOVED***
	verifier []byte
	err      string
***REMOVED***

func parseField(s, k string) (string, error) ***REMOVED***
	t := strings.TrimPrefix(s, k+"=")
	if t == s ***REMOVED***
		return "", fmt.Errorf("error parsing '%s' for field '%s'", s, k)
	***REMOVED***
	return t, nil
***REMOVED***

func parseGS2Flag(s string) (string, error) ***REMOVED***
	if s[0] == 'p' ***REMOVED***
		return "", fmt.Errorf("channel binding requested but not supported")
	***REMOVED***

	if s == "n" || s == "y" ***REMOVED***
		return s, nil
	***REMOVED***

	return "", fmt.Errorf("error parsing '%s' for gs2 flag", s)
***REMOVED***

func parseFieldBase64(s, k string) ([]byte, error) ***REMOVED***
	raw, err := parseField(s, k)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	dec, err := base64.StdEncoding.DecodeString(raw)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return dec, nil
***REMOVED***

func parseFieldInt(s, k string) (int, error) ***REMOVED***
	raw, err := parseField(s, k)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***

	num, err := strconv.Atoi(raw)
	if err != nil ***REMOVED***
		return 0, fmt.Errorf("error parsing field '%s': %v", k, err)
	***REMOVED***

	return num, nil
***REMOVED***

func parseClientFirst(c1 string) (msg c1Msg, err error) ***REMOVED***

	fields := strings.Split(c1, ",")
	if len(fields) < 4 ***REMOVED***
		err = errors.New("not enough fields in first server message")
		return
	***REMOVED***

	gs2flag, err := parseGS2Flag(fields[0])
	if err != nil ***REMOVED***
		return
	***REMOVED***

	// 'a' field is optional
	if len(fields[1]) > 0 ***REMOVED***
		msg.authzID, err = parseField(fields[1], "a")
		if err != nil ***REMOVED***
			return
		***REMOVED***
	***REMOVED***

	// Recombine and save the gs2 header
	msg.gs2Header = gs2flag + "," + msg.authzID + ","

	// Check for unsupported extensions field "m".
	if strings.HasPrefix(fields[2], "m=") ***REMOVED***
		err = errors.New("SCRAM message extensions are not supported")
		return
	***REMOVED***

	msg.username, err = parseField(fields[2], "n")
	if err != nil ***REMOVED***
		return
	***REMOVED***

	msg.nonce, err = parseField(fields[3], "r")
	if err != nil ***REMOVED***
		return
	***REMOVED***

	msg.c1b = strings.Join(fields[2:], ",")

	return
***REMOVED***

func parseClientFinal(c2 string) (msg c2Msg, err error) ***REMOVED***
	fields := strings.Split(c2, ",")
	if len(fields) < 3 ***REMOVED***
		err = errors.New("not enough fields in first server message")
		return
	***REMOVED***

	msg.cbind, err = parseFieldBase64(fields[0], "c")
	if err != nil ***REMOVED***
		return
	***REMOVED***

	msg.nonce, err = parseField(fields[1], "r")
	if err != nil ***REMOVED***
		return
	***REMOVED***

	// Extension fields may come between nonce and proof, so we
	// grab the *last* fields as proof.
	msg.proof, err = parseFieldBase64(fields[len(fields)-1], "p")
	if err != nil ***REMOVED***
		return
	***REMOVED***

	msg.c2wop = c2[:strings.LastIndex(c2, ",")]

	return
***REMOVED***

func parseServerFirst(s1 string) (msg s1Msg, err error) ***REMOVED***

	// Check for unsupported extensions field "m".
	if strings.HasPrefix(s1, "m=") ***REMOVED***
		err = errors.New("SCRAM message extensions are not supported")
		return
	***REMOVED***

	fields := strings.Split(s1, ",")
	if len(fields) < 3 ***REMOVED***
		err = errors.New("not enough fields in first server message")
		return
	***REMOVED***

	msg.nonce, err = parseField(fields[0], "r")
	if err != nil ***REMOVED***
		return
	***REMOVED***

	msg.salt, err = parseFieldBase64(fields[1], "s")
	if err != nil ***REMOVED***
		return
	***REMOVED***

	msg.iters, err = parseFieldInt(fields[2], "i")

	return
***REMOVED***

func parseServerFinal(s2 string) (msg s2Msg, err error) ***REMOVED***
	fields := strings.Split(s2, ",")

	msg.verifier, err = parseFieldBase64(fields[0], "v")
	if err == nil ***REMOVED***
		return
	***REMOVED***

	msg.err, err = parseField(fields[0], "e")

	return
***REMOVED***
