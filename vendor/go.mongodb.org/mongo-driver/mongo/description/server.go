// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package description

import (
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/internal"
	"go.mongodb.org/mongo-driver/mongo/address"
	"go.mongodb.org/mongo-driver/tag"
)

// SelectedServer augments the Server type by also including the TopologyKind of the topology that includes the server.
// This type should be used to track the state of a server that was selected to perform an operation.
type SelectedServer struct ***REMOVED***
	Server
	Kind TopologyKind
***REMOVED***

// Server contains information about a node in a cluster. This is created from hello command responses. If the value
// of the Kind field is LoadBalancer, only the Addr and Kind fields will be set. All other fields will be set to the
// zero value of the field's type.
type Server struct ***REMOVED***
	Addr address.Address

	Arbiters              []string
	AverageRTT            time.Duration
	AverageRTTSet         bool
	Compression           []string // compression methods returned by server
	CanonicalAddr         address.Address
	ElectionID            primitive.ObjectID
	HeartbeatInterval     time.Duration
	HelloOK               bool
	Hosts                 []string
	LastError             error
	LastUpdateTime        time.Time
	LastWriteTime         time.Time
	MaxBatchCount         uint32
	MaxDocumentSize       uint32
	MaxMessageSize        uint32
	Members               []address.Address
	Passives              []string
	Passive               bool
	Primary               address.Address
	ReadOnly              bool
	ServiceID             *primitive.ObjectID // Only set for servers that are deployed behind a load balancer.
	SessionTimeoutMinutes uint32
	SetName               string
	SetVersion            uint32
	Tags                  tag.Set
	TopologyVersion       *TopologyVersion
	Kind                  ServerKind
	WireVersion           *VersionRange
***REMOVED***

// NewServer creates a new server description from the given hello command response.
func NewServer(addr address.Address, response bson.Raw) Server ***REMOVED***
	desc := Server***REMOVED***Addr: addr, CanonicalAddr: addr, LastUpdateTime: time.Now().UTC()***REMOVED***
	elements, err := response.Elements()
	if err != nil ***REMOVED***
		desc.LastError = err
		return desc
	***REMOVED***
	var ok bool
	var isReplicaSet, isWritablePrimary, hidden, secondary, arbiterOnly bool
	var msg string
	var version VersionRange
	for _, element := range elements ***REMOVED***
		switch element.Key() ***REMOVED***
		case "arbiters":
			var err error
			desc.Arbiters, err = internal.StringSliceFromRawElement(element)
			if err != nil ***REMOVED***
				desc.LastError = err
				return desc
			***REMOVED***
		case "arbiterOnly":
			arbiterOnly, ok = element.Value().BooleanOK()
			if !ok ***REMOVED***
				desc.LastError = fmt.Errorf("expected 'arbiterOnly' to be a boolean but it's a BSON %s", element.Value().Type)
				return desc
			***REMOVED***
		case "compression":
			var err error
			desc.Compression, err = internal.StringSliceFromRawElement(element)
			if err != nil ***REMOVED***
				desc.LastError = err
				return desc
			***REMOVED***
		case "electionId":
			desc.ElectionID, ok = element.Value().ObjectIDOK()
			if !ok ***REMOVED***
				desc.LastError = fmt.Errorf("expected 'electionId' to be a objectID but it's a BSON %s", element.Value().Type)
				return desc
			***REMOVED***
		case "helloOk":
			desc.HelloOK, ok = element.Value().BooleanOK()
			if !ok ***REMOVED***
				desc.LastError = fmt.Errorf("expected 'helloOk' to be a boolean but it's a BSON %s", element.Value().Type)
				return desc
			***REMOVED***
		case "hidden":
			hidden, ok = element.Value().BooleanOK()
			if !ok ***REMOVED***
				desc.LastError = fmt.Errorf("expected 'hidden' to be a boolean but it's a BSON %s", element.Value().Type)
				return desc
			***REMOVED***
		case "hosts":
			var err error
			desc.Hosts, err = internal.StringSliceFromRawElement(element)
			if err != nil ***REMOVED***
				desc.LastError = err
				return desc
			***REMOVED***
		case "isWritablePrimary":
			isWritablePrimary, ok = element.Value().BooleanOK()
			if !ok ***REMOVED***
				desc.LastError = fmt.Errorf("expected 'isWritablePrimary' to be a boolean but it's a BSON %s", element.Value().Type)
				return desc
			***REMOVED***
		case internal.LegacyHelloLowercase:
			isWritablePrimary, ok = element.Value().BooleanOK()
			if !ok ***REMOVED***
				desc.LastError = fmt.Errorf("expected legacy hello to be a boolean but it's a BSON %s", element.Value().Type)
				return desc
			***REMOVED***
		case "isreplicaset":
			isReplicaSet, ok = element.Value().BooleanOK()
			if !ok ***REMOVED***
				desc.LastError = fmt.Errorf("expected 'isreplicaset' to be a boolean but it's a BSON %s", element.Value().Type)
				return desc
			***REMOVED***
		case "lastWrite":
			lastWrite, ok := element.Value().DocumentOK()
			if !ok ***REMOVED***
				desc.LastError = fmt.Errorf("expected 'lastWrite' to be a document but it's a BSON %s", element.Value().Type)
				return desc
			***REMOVED***
			dateTime, err := lastWrite.LookupErr("lastWriteDate")
			if err == nil ***REMOVED***
				dt, ok := dateTime.DateTimeOK()
				if !ok ***REMOVED***
					desc.LastError = fmt.Errorf("expected 'lastWriteDate' to be a datetime but it's a BSON %s", dateTime.Type)
					return desc
				***REMOVED***
				desc.LastWriteTime = time.Unix(dt/1000, dt%1000*1000000).UTC()
			***REMOVED***
		case "logicalSessionTimeoutMinutes":
			i64, ok := element.Value().AsInt64OK()
			if !ok ***REMOVED***
				desc.LastError = fmt.Errorf("expected 'logicalSessionTimeoutMinutes' to be an integer but it's a BSON %s", element.Value().Type)
				return desc
			***REMOVED***
			desc.SessionTimeoutMinutes = uint32(i64)
		case "maxBsonObjectSize":
			i64, ok := element.Value().AsInt64OK()
			if !ok ***REMOVED***
				desc.LastError = fmt.Errorf("expected 'maxBsonObjectSize' to be an integer but it's a BSON %s", element.Value().Type)
				return desc
			***REMOVED***
			desc.MaxDocumentSize = uint32(i64)
		case "maxMessageSizeBytes":
			i64, ok := element.Value().AsInt64OK()
			if !ok ***REMOVED***
				desc.LastError = fmt.Errorf("expected 'maxMessageSizeBytes' to be an integer but it's a BSON %s", element.Value().Type)
				return desc
			***REMOVED***
			desc.MaxMessageSize = uint32(i64)
		case "maxWriteBatchSize":
			i64, ok := element.Value().AsInt64OK()
			if !ok ***REMOVED***
				desc.LastError = fmt.Errorf("expected 'maxWriteBatchSize' to be an integer but it's a BSON %s", element.Value().Type)
				return desc
			***REMOVED***
			desc.MaxBatchCount = uint32(i64)
		case "me":
			me, ok := element.Value().StringValueOK()
			if !ok ***REMOVED***
				desc.LastError = fmt.Errorf("expected 'me' to be a string but it's a BSON %s", element.Value().Type)
				return desc
			***REMOVED***
			desc.CanonicalAddr = address.Address(me).Canonicalize()
		case "maxWireVersion":
			version.Max, ok = element.Value().AsInt32OK()
			if !ok ***REMOVED***
				desc.LastError = fmt.Errorf("expected 'maxWireVersion' to be an integer but it's a BSON %s", element.Value().Type)
				return desc
			***REMOVED***
		case "minWireVersion":
			version.Min, ok = element.Value().AsInt32OK()
			if !ok ***REMOVED***
				desc.LastError = fmt.Errorf("expected 'minWireVersion' to be an integer but it's a BSON %s", element.Value().Type)
				return desc
			***REMOVED***
		case "msg":
			msg, ok = element.Value().StringValueOK()
			if !ok ***REMOVED***
				desc.LastError = fmt.Errorf("expected 'msg' to be a string but it's a BSON %s", element.Value().Type)
				return desc
			***REMOVED***
		case "ok":
			okay, ok := element.Value().AsInt32OK()
			if !ok ***REMOVED***
				desc.LastError = fmt.Errorf("expected 'ok' to be a boolean but it's a BSON %s", element.Value().Type)
				return desc
			***REMOVED***
			if okay != 1 ***REMOVED***
				desc.LastError = errors.New("not ok")
				return desc
			***REMOVED***
		case "passives":
			var err error
			desc.Passives, err = internal.StringSliceFromRawElement(element)
			if err != nil ***REMOVED***
				desc.LastError = err
				return desc
			***REMOVED***
		case "passive":
			desc.Passive, ok = element.Value().BooleanOK()
			if !ok ***REMOVED***
				desc.LastError = fmt.Errorf("expected 'passive' to be a boolean but it's a BSON %s", element.Value().Type)
				return desc
			***REMOVED***
		case "primary":
			primary, ok := element.Value().StringValueOK()
			if !ok ***REMOVED***
				desc.LastError = fmt.Errorf("expected 'primary' to be a string but it's a BSON %s", element.Value().Type)
				return desc
			***REMOVED***
			desc.Primary = address.Address(primary)
		case "readOnly":
			desc.ReadOnly, ok = element.Value().BooleanOK()
			if !ok ***REMOVED***
				desc.LastError = fmt.Errorf("expected 'readOnly' to be a boolean but it's a BSON %s", element.Value().Type)
				return desc
			***REMOVED***
		case "secondary":
			secondary, ok = element.Value().BooleanOK()
			if !ok ***REMOVED***
				desc.LastError = fmt.Errorf("expected 'secondary' to be a boolean but it's a BSON %s", element.Value().Type)
				return desc
			***REMOVED***
		case "serviceId":
			oid, ok := element.Value().ObjectIDOK()
			if !ok ***REMOVED***
				desc.LastError = fmt.Errorf("expected 'serviceId' to be an ObjectId but it's a BSON %s", element.Value().Type)
			***REMOVED***
			desc.ServiceID = &oid
		case "setName":
			desc.SetName, ok = element.Value().StringValueOK()
			if !ok ***REMOVED***
				desc.LastError = fmt.Errorf("expected 'setName' to be a string but it's a BSON %s", element.Value().Type)
				return desc
			***REMOVED***
		case "setVersion":
			i64, ok := element.Value().AsInt64OK()
			if !ok ***REMOVED***
				desc.LastError = fmt.Errorf("expected 'setVersion' to be an integer but it's a BSON %s", element.Value().Type)
				return desc
			***REMOVED***
			desc.SetVersion = uint32(i64)
		case "tags":
			m, err := decodeStringMap(element, "tags")
			if err != nil ***REMOVED***
				desc.LastError = err
				return desc
			***REMOVED***
			desc.Tags = tag.NewTagSetFromMap(m)
		case "topologyVersion":
			doc, ok := element.Value().DocumentOK()
			if !ok ***REMOVED***
				desc.LastError = fmt.Errorf("expected 'topologyVersion' to be a document but it's a BSON %s", element.Value().Type)
				return desc
			***REMOVED***

			desc.TopologyVersion, err = NewTopologyVersion(doc)
			if err != nil ***REMOVED***
				desc.LastError = err
				return desc
			***REMOVED***
		***REMOVED***
	***REMOVED***

	for _, host := range desc.Hosts ***REMOVED***
		desc.Members = append(desc.Members, address.Address(host).Canonicalize())
	***REMOVED***

	for _, passive := range desc.Passives ***REMOVED***
		desc.Members = append(desc.Members, address.Address(passive).Canonicalize())
	***REMOVED***

	for _, arbiter := range desc.Arbiters ***REMOVED***
		desc.Members = append(desc.Members, address.Address(arbiter).Canonicalize())
	***REMOVED***

	desc.Kind = Standalone

	if isReplicaSet ***REMOVED***
		desc.Kind = RSGhost
	***REMOVED*** else if desc.SetName != "" ***REMOVED***
		if isWritablePrimary ***REMOVED***
			desc.Kind = RSPrimary
		***REMOVED*** else if hidden ***REMOVED***
			desc.Kind = RSMember
		***REMOVED*** else if secondary ***REMOVED***
			desc.Kind = RSSecondary
		***REMOVED*** else if arbiterOnly ***REMOVED***
			desc.Kind = RSArbiter
		***REMOVED*** else ***REMOVED***
			desc.Kind = RSMember
		***REMOVED***
	***REMOVED*** else if msg == "isdbgrid" ***REMOVED***
		desc.Kind = Mongos
	***REMOVED***

	desc.WireVersion = &version

	return desc
***REMOVED***

// NewDefaultServer creates a new unknown server description with the given address.
func NewDefaultServer(addr address.Address) Server ***REMOVED***
	return NewServerFromError(addr, nil, nil)
***REMOVED***

// NewServerFromError creates a new unknown server description with the given parameters.
func NewServerFromError(addr address.Address, err error, tv *TopologyVersion) Server ***REMOVED***
	return Server***REMOVED***
		Addr:            addr,
		LastError:       err,
		Kind:            Unknown,
		TopologyVersion: tv,
	***REMOVED***
***REMOVED***

// SetAverageRTT sets the average round trip time for this server description.
func (s Server) SetAverageRTT(rtt time.Duration) Server ***REMOVED***
	s.AverageRTT = rtt
	s.AverageRTTSet = true
	return s
***REMOVED***

// DataBearing returns true if the server is a data bearing server.
func (s Server) DataBearing() bool ***REMOVED***
	return s.Kind == RSPrimary ||
		s.Kind == RSSecondary ||
		s.Kind == Mongos ||
		s.Kind == Standalone
***REMOVED***

// LoadBalanced returns true if the server is a load balancer or is behind a load balancer.
func (s Server) LoadBalanced() bool ***REMOVED***
	return s.Kind == LoadBalancer || s.ServiceID != nil
***REMOVED***

// String implements the Stringer interface
func (s Server) String() string ***REMOVED***
	str := fmt.Sprintf("Addr: %s, Type: %s",
		s.Addr, s.Kind)
	if len(s.Tags) != 0 ***REMOVED***
		str += fmt.Sprintf(", Tag sets: %s", s.Tags)
	***REMOVED***

	if s.AverageRTTSet ***REMOVED***
		str += fmt.Sprintf(", Average RTT: %d", s.AverageRTT)
	***REMOVED***

	if s.LastError != nil ***REMOVED***
		str += fmt.Sprintf(", Last error: %s", s.LastError)
	***REMOVED***
	return str
***REMOVED***

func decodeStringMap(element bson.RawElement, name string) (map[string]string, error) ***REMOVED***
	doc, ok := element.Value().DocumentOK()
	if !ok ***REMOVED***
		return nil, fmt.Errorf("expected '%s' to be a document but it's a BSON %s", name, element.Value().Type)
	***REMOVED***
	elements, err := doc.Elements()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	m := make(map[string]string)
	for _, element := range elements ***REMOVED***
		key := element.Key()
		value, ok := element.Value().StringValueOK()
		if !ok ***REMOVED***
			return nil, fmt.Errorf("expected '%s' to be a document of strings, but found a BSON %s", name, element.Value().Type)
		***REMOVED***
		m[key] = value
	***REMOVED***
	return m, nil
***REMOVED***

// Equal compares two server descriptions and returns true if they are equal
func (s Server) Equal(other Server) bool ***REMOVED***
	if s.CanonicalAddr.String() != other.CanonicalAddr.String() ***REMOVED***
		return false
	***REMOVED***

	if !sliceStringEqual(s.Arbiters, other.Arbiters) ***REMOVED***
		return false
	***REMOVED***

	if !sliceStringEqual(s.Hosts, other.Hosts) ***REMOVED***
		return false
	***REMOVED***

	if !sliceStringEqual(s.Passives, other.Passives) ***REMOVED***
		return false
	***REMOVED***

	if s.Primary != other.Primary ***REMOVED***
		return false
	***REMOVED***

	if s.SetName != other.SetName ***REMOVED***
		return false
	***REMOVED***

	if s.Kind != other.Kind ***REMOVED***
		return false
	***REMOVED***

	if s.LastError != nil || other.LastError != nil ***REMOVED***
		if s.LastError == nil || other.LastError == nil ***REMOVED***
			return false
		***REMOVED***
		if s.LastError.Error() != other.LastError.Error() ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***

	if !s.WireVersion.Equals(other.WireVersion) ***REMOVED***
		return false
	***REMOVED***

	if len(s.Tags) != len(other.Tags) || !s.Tags.ContainsAll(other.Tags) ***REMOVED***
		return false
	***REMOVED***

	if s.SetVersion != other.SetVersion ***REMOVED***
		return false
	***REMOVED***

	if s.ElectionID != other.ElectionID ***REMOVED***
		return false
	***REMOVED***

	if s.SessionTimeoutMinutes != other.SessionTimeoutMinutes ***REMOVED***
		return false
	***REMOVED***

	// If TopologyVersion is nil for both servers, CompareToIncoming will return -1 because it assumes that the
	// incoming response is newer. We want the descriptions to be considered equal in this case, though, so an
	// explicit check is required.
	if s.TopologyVersion == nil && other.TopologyVersion == nil ***REMOVED***
		return true
	***REMOVED***
	return s.TopologyVersion.CompareToIncoming(other.TopologyVersion) == 0
***REMOVED***

func sliceStringEqual(a []string, b []string) bool ***REMOVED***
	if len(a) != len(b) ***REMOVED***
		return false
	***REMOVED***
	for i, v := range a ***REMOVED***
		if v != b[i] ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***
