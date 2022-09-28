// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package session

import (
	"sync"

	"go.mongodb.org/mongo-driver/bson"
)

// ClusterClock represents a logical clock for keeping track of cluster time.
type ClusterClock struct ***REMOVED***
	clusterTime bson.Raw
	lock        sync.Mutex
***REMOVED***

// GetClusterTime returns the cluster's current time.
func (cc *ClusterClock) GetClusterTime() bson.Raw ***REMOVED***
	var ct bson.Raw
	cc.lock.Lock()
	ct = cc.clusterTime
	cc.lock.Unlock()

	return ct
***REMOVED***

// AdvanceClusterTime updates the cluster's current time.
func (cc *ClusterClock) AdvanceClusterTime(clusterTime bson.Raw) ***REMOVED***
	cc.lock.Lock()
	cc.clusterTime = MaxClusterTime(cc.clusterTime, clusterTime)
	cc.lock.Unlock()
***REMOVED***
