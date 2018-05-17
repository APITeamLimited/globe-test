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

package kafka

import (
	"testing"

	"github.com/loadimpact/k6/stats/influxdb"
	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v3"
)

func TestConfigParseArg(t *testing.T) ***REMOVED***
	c, err := ParseArg("brokers=broker1,topic=someTopic,format=influxdb")
	expInfluxConfig := influxdb.Config***REMOVED******REMOVED***
	assert.Nil(t, err)
	assert.Equal(t, []string***REMOVED***"broker1"***REMOVED***, c.Brokers)
	assert.Equal(t, null.StringFrom("someTopic"), c.Topic)
	assert.Equal(t, null.StringFrom("influxdb"), c.Format)
	assert.Equal(t, expInfluxConfig, c.InfluxDBConfig)

	c, err = ParseArg("brokers=***REMOVED***broker2,broker3:9092***REMOVED***,topic=someTopic2,format=json")
	assert.Nil(t, err)
	assert.Equal(t, []string***REMOVED***"broker2", "broker3:9092"***REMOVED***, c.Brokers)
	assert.Equal(t, null.StringFrom("someTopic2"), c.Topic)
	assert.Equal(t, null.StringFrom("json"), c.Format)

	c, err = ParseArg("brokers=***REMOVED***broker2,broker3:9092***REMOVED***,topic=someTopic,format=influxdb,influxdb.tagsAsFields=fake")
	expInfluxConfig = influxdb.Config***REMOVED***
		TagsAsFields: []string***REMOVED***"fake"***REMOVED***,
	***REMOVED***
	assert.Nil(t, err)
	assert.Equal(t, []string***REMOVED***"broker2", "broker3:9092"***REMOVED***, c.Brokers)
	assert.Equal(t, null.StringFrom("someTopic"), c.Topic)
	assert.Equal(t, null.StringFrom("influxdb"), c.Format)
	assert.Equal(t, expInfluxConfig, c.InfluxDBConfig)

	c, err = ParseArg("brokers=***REMOVED***broker2,broker3:9092***REMOVED***,topic=someTopic,format=influxdb,influxdb.tagsAsFields=***REMOVED***fake,anotherFake***REMOVED***")
	expInfluxConfig = influxdb.Config***REMOVED***
		TagsAsFields: []string***REMOVED***"fake", "anotherFake"***REMOVED***,
	***REMOVED***
	assert.Nil(t, err)
	assert.Equal(t, []string***REMOVED***"broker2", "broker3:9092"***REMOVED***, c.Brokers)
	assert.Equal(t, null.StringFrom("someTopic"), c.Topic)
	assert.Equal(t, null.StringFrom("influxdb"), c.Format)
	assert.Equal(t, expInfluxConfig, c.InfluxDBConfig)
***REMOVED***
