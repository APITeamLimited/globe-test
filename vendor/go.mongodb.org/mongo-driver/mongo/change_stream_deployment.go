// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package mongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo/description"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
)

type changeStreamDeployment struct ***REMOVED***
	topologyKind description.TopologyKind
	server       driver.Server
	conn         driver.Connection
***REMOVED***

var _ driver.Deployment = (*changeStreamDeployment)(nil)
var _ driver.Server = (*changeStreamDeployment)(nil)
var _ driver.ErrorProcessor = (*changeStreamDeployment)(nil)

func (c *changeStreamDeployment) SelectServer(context.Context, description.ServerSelector) (driver.Server, error) ***REMOVED***
	return c, nil
***REMOVED***

func (c *changeStreamDeployment) Kind() description.TopologyKind ***REMOVED***
	return c.topologyKind
***REMOVED***

func (c *changeStreamDeployment) Connection(context.Context) (driver.Connection, error) ***REMOVED***
	return c.conn, nil
***REMOVED***

func (c *changeStreamDeployment) MinRTT() time.Duration ***REMOVED***
	return c.server.MinRTT()
***REMOVED***

func (c *changeStreamDeployment) RTT90() time.Duration ***REMOVED***
	return c.server.RTT90()
***REMOVED***

func (c *changeStreamDeployment) ProcessError(err error, conn driver.Connection) driver.ProcessErrorResult ***REMOVED***
	ep, ok := c.server.(driver.ErrorProcessor)
	if !ok ***REMOVED***
		return driver.NoChange
	***REMOVED***

	return ep.ProcessError(err, conn)
***REMOVED***
