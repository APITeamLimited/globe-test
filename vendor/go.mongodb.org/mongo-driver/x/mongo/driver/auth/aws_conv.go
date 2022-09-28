// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package auth

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver/auth/internal/awsv4"
)

type clientState int

const (
	clientStarting clientState = iota
	clientFirst
	clientFinal
	clientDone
)

type awsConversation struct ***REMOVED***
	state    clientState
	valid    bool
	nonce    []byte
	username string
	password string
	token    string
***REMOVED***

type serverMessage struct ***REMOVED***
	Nonce primitive.Binary `bson:"s"`
	Host  string           `bson:"h"`
***REMOVED***

type ecsResponse struct ***REMOVED***
	AccessKeyID     string `json:"AccessKeyId"`
	SecretAccessKey string `json:"SecretAccessKey"`
	Token           string `json:"Token"`
***REMOVED***

const (
	amzDateFormat       = "20060102T150405Z"
	awsRelativeURI      = "http://169.254.170.2/"
	awsEC2URI           = "http://169.254.169.254/"
	awsEC2RolePath      = "latest/meta-data/iam/security-credentials/"
	awsEC2TokenPath     = "latest/api/token"
	defaultRegion       = "us-east-1"
	maxHostLength       = 255
	defaultHTTPTimeout  = 10 * time.Second
	responceNonceLength = 64
)

// Step takes a string provided from a server (or just an empty string for the
// very first conversation step) and attempts to move the authentication
// conversation forward.  It returns a string to be sent to the server or an
// error if the server message is invalid.  Calling Step after a conversation
// completes is also an error.
func (ac *awsConversation) Step(challenge []byte) (response []byte, err error) ***REMOVED***
	switch ac.state ***REMOVED***
	case clientStarting:
		ac.state = clientFirst
		response = ac.firstMsg()
	case clientFirst:
		ac.state = clientFinal
		response, err = ac.finalMsg(challenge)
	case clientFinal:
		ac.state = clientDone
		ac.valid = true
	default:
		response, err = nil, errors.New("Conversation already completed")
	***REMOVED***
	return
***REMOVED***

// Done returns true if the conversation is completed or has errored.
func (ac *awsConversation) Done() bool ***REMOVED***
	return ac.state == clientDone
***REMOVED***

// Valid returns true if the conversation successfully authenticated with the
// server, including counter-validation that the server actually has the
// user's stored credentials.
func (ac *awsConversation) Valid() bool ***REMOVED***
	return ac.valid
***REMOVED***

func getRegion(host string) (string, error) ***REMOVED***
	region := defaultRegion

	if len(host) == 0 ***REMOVED***
		return "", errors.New("invalid STS host: empty")
	***REMOVED***
	if len(host) > maxHostLength ***REMOVED***
		return "", errors.New("invalid STS host: too large")
	***REMOVED***
	// The implicit region for sts.amazonaws.com is us-east-1
	if host == "sts.amazonaws.com" ***REMOVED***
		return region, nil
	***REMOVED***
	if strings.HasPrefix(host, ".") || strings.HasSuffix(host, ".") || strings.Contains(host, "..") ***REMOVED***
		return "", errors.New("invalid STS host: empty part")
	***REMOVED***

	// If the host has multiple parts, the second part is the region
	parts := strings.Split(host, ".")
	if len(parts) >= 2 ***REMOVED***
		region = parts[1]
	***REMOVED***

	return region, nil
***REMOVED***

func (ac *awsConversation) validateAndMakeCredentials() (*awsv4.StaticProvider, error) ***REMOVED***
	if ac.username != "" && ac.password == "" ***REMOVED***
		return nil, errors.New("ACCESS_KEY_ID is set, but SECRET_ACCESS_KEY is missing")
	***REMOVED***
	if ac.username == "" && ac.password != "" ***REMOVED***
		return nil, errors.New("SECRET_ACCESS_KEY is set, but ACCESS_KEY_ID is missing")
	***REMOVED***
	if ac.username == "" && ac.password == "" && ac.token != "" ***REMOVED***
		return nil, errors.New("AWS_SESSION_TOKEN is set, but ACCESS_KEY_ID and SECRET_ACCESS_KEY are missing")
	***REMOVED***
	if ac.username != "" || ac.password != "" || ac.token != "" ***REMOVED***
		return &awsv4.StaticProvider***REMOVED***Value: awsv4.Value***REMOVED***
			AccessKeyID:     ac.username,
			SecretAccessKey: ac.password,
			SessionToken:    ac.token,
		***REMOVED******REMOVED***, nil
	***REMOVED***
	return nil, nil
***REMOVED***

func executeAWSHTTPRequest(req *http.Request) ([]byte, error) ***REMOVED***
	ctx, cancel := context.WithTimeout(context.Background(), defaultHTTPTimeout)
	defer cancel()
	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer func() ***REMOVED*** _ = resp.Body.Close() ***REMOVED***()

	return ioutil.ReadAll(resp.Body)
***REMOVED***

func (ac *awsConversation) getEC2Credentials() (*awsv4.StaticProvider, error) ***REMOVED***
	// get token
	req, err := http.NewRequest("PUT", awsEC2URI+awsEC2TokenPath, nil)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	req.Header.Set("X-aws-ec2-metadata-token-ttl-seconds", "30")

	token, err := executeAWSHTTPRequest(req)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if len(token) == 0 ***REMOVED***
		return nil, errors.New("unable to retrieve token from EC2 metadata")
	***REMOVED***
	tokenStr := string(token)

	// get role name
	req, err = http.NewRequest("GET", awsEC2URI+awsEC2RolePath, nil)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	req.Header.Set("X-aws-ec2-metadata-token", tokenStr)

	role, err := executeAWSHTTPRequest(req)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if len(role) == 0 ***REMOVED***
		return nil, errors.New("unable to retrieve role_name from EC2 metadata")
	***REMOVED***

	// get credentials
	pathWithRole := awsEC2URI + awsEC2RolePath + string(role)
	req, err = http.NewRequest("GET", pathWithRole, nil)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	req.Header.Set("X-aws-ec2-metadata-token", tokenStr)
	creds, err := executeAWSHTTPRequest(req)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var es2Resp ecsResponse
	err = json.Unmarshal(creds, &es2Resp)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	ac.username = es2Resp.AccessKeyID
	ac.password = es2Resp.SecretAccessKey
	ac.token = es2Resp.Token

	return ac.validateAndMakeCredentials()
***REMOVED***

func (ac *awsConversation) getCredentials() (*awsv4.StaticProvider, error) ***REMOVED***
	// Credentials passed through URI
	creds, err := ac.validateAndMakeCredentials()
	if creds != nil || err != nil ***REMOVED***
		return creds, err
	***REMOVED***

	// Credentials from environment variables
	ac.username = os.Getenv("AWS_ACCESS_KEY_ID")
	ac.password = os.Getenv("AWS_SECRET_ACCESS_KEY")
	ac.token = os.Getenv("AWS_SESSION_TOKEN")

	creds, err = ac.validateAndMakeCredentials()
	if creds != nil || err != nil ***REMOVED***
		return creds, err
	***REMOVED***

	// Credentials from ECS metadata
	relativeEcsURI := os.Getenv("AWS_CONTAINER_CREDENTIALS_RELATIVE_URI")
	if len(relativeEcsURI) > 0 ***REMOVED***
		fullURI := awsRelativeURI + relativeEcsURI

		req, err := http.NewRequest("GET", fullURI, nil)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		body, err := executeAWSHTTPRequest(req)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		var espResp ecsResponse
		err = json.Unmarshal(body, &espResp)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		ac.username = espResp.AccessKeyID
		ac.password = espResp.SecretAccessKey
		ac.token = espResp.Token

		creds, err = ac.validateAndMakeCredentials()
		if creds != nil || err != nil ***REMOVED***
			return creds, err
		***REMOVED***
	***REMOVED***

	// Credentials from EC2 metadata
	creds, err = ac.getEC2Credentials()
	if creds == nil && err == nil ***REMOVED***
		return nil, errors.New("unable to get credentials")
	***REMOVED***
	return creds, err
***REMOVED***

func (ac *awsConversation) firstMsg() []byte ***REMOVED***
	// Values are cached for use in final message parameters
	ac.nonce = make([]byte, 32)
	_, _ = rand.Read(ac.nonce)

	idx, msg := bsoncore.AppendDocumentStart(nil)
	msg = bsoncore.AppendInt32Element(msg, "p", 110)
	msg = bsoncore.AppendBinaryElement(msg, "r", 0x00, ac.nonce)
	msg, _ = bsoncore.AppendDocumentEnd(msg, idx)
	return msg
***REMOVED***

func (ac *awsConversation) finalMsg(s1 []byte) ([]byte, error) ***REMOVED***
	var sm serverMessage
	err := bson.Unmarshal(s1, &sm)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// Check nonce prefix
	if sm.Nonce.Subtype != 0x00 ***REMOVED***
		return nil, errors.New("server reply contained unexpected binary subtype")
	***REMOVED***
	if len(sm.Nonce.Data) != responceNonceLength ***REMOVED***
		return nil, fmt.Errorf("server reply nonce was not %v bytes", responceNonceLength)
	***REMOVED***
	if !bytes.HasPrefix(sm.Nonce.Data, ac.nonce) ***REMOVED***
		return nil, errors.New("server nonce did not extend client nonce")
	***REMOVED***

	region, err := getRegion(sm.Host)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	creds, err := ac.getCredentials()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	currentTime := time.Now().UTC()
	body := "Action=GetCallerIdentity&Version=2011-06-15"

	// Create http.Request
	req, _ := http.NewRequest("POST", "/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", "43")
	req.Host = sm.Host
	req.Header.Set("X-Amz-Date", currentTime.Format(amzDateFormat))
	if len(ac.token) > 0 ***REMOVED***
		req.Header.Set("X-Amz-Security-Token", ac.token)
	***REMOVED***
	req.Header.Set("X-MongoDB-Server-Nonce", base64.StdEncoding.EncodeToString(sm.Nonce.Data))
	req.Header.Set("X-MongoDB-GS2-CB-Flag", "n")

	// Create signer with credentials
	signer := awsv4.NewSigner(creds)

	// Get signed header
	_, err = signer.Sign(req, strings.NewReader(body), "sts", region, currentTime)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// create message
	idx, msg := bsoncore.AppendDocumentStart(nil)
	msg = bsoncore.AppendStringElement(msg, "a", req.Header.Get("Authorization"))
	msg = bsoncore.AppendStringElement(msg, "d", req.Header.Get("X-Amz-Date"))
	if len(ac.token) > 0 ***REMOVED***
		msg = bsoncore.AppendStringElement(msg, "t", ac.token)
	***REMOVED***
	msg, _ = bsoncore.AppendDocumentEnd(msg, idx)

	return msg, nil
***REMOVED***
