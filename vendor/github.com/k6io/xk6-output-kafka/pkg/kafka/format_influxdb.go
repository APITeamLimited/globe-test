/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2021 Load Impact
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
	"fmt"
	"strconv"
	"strings"

	client "github.com/influxdata/influxdb1-client/v2"
	"go.k6.io/k6/lib/types"
	"go.k6.io/k6/stats"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
)

type extractTagsToValuesFunc func(map[string]string, map[string]interface***REMOVED******REMOVED***) map[string]interface***REMOVED******REMOVED***

// format returns a string array of metrics in influx line-protocol
func formatAsInfluxdbV1(
	logger logrus.FieldLogger, samples []stats.Sample, extractTagsToValues extractTagsToValuesFunc,
) ([]string, error) ***REMOVED***
	var metrics []string
	type cacheItem struct ***REMOVED***
		tags   map[string]string
		values map[string]interface***REMOVED******REMOVED***
	***REMOVED***
	cache := map[*stats.SampleTags]cacheItem***REMOVED******REMOVED***
	for _, sample := range samples ***REMOVED***
		var tags map[string]string
		values := make(map[string]interface***REMOVED******REMOVED***)
		if cached, ok := cache[sample.Tags]; ok ***REMOVED***
			tags = cached.tags
			for k, v := range cached.values ***REMOVED***
				values[k] = v
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			tags = sample.Tags.CloneTags()
			extractTagsToValues(tags, values)
			cache[sample.Tags] = cacheItem***REMOVED***tags, values***REMOVED***
		***REMOVED***
		values["value"] = sample.Value
		p, err := client.NewPoint(
			sample.Metric.Name,
			tags,
			values,
			sample.Time,
		)
		if err != nil ***REMOVED***
			logger.WithError(err).Error("InfluxDB: Couldn't make point from sample!")
			return nil, err
		***REMOVED***
		metrics = append(metrics, p.String())
	***REMOVED***

	return metrics, nil
***REMOVED***

// FieldKind defines Enum for tag-to-field type conversion
type FieldKind int

const (
	// String field (default)
	String FieldKind = iota
	// Int field
	Int
	// Float field
	Float
	// Bool field
	Bool
)

func newExtractTagsFields(fieldKinds map[string]FieldKind) extractTagsToValuesFunc ***REMOVED***
	return func(tags map[string]string, values map[string]interface***REMOVED******REMOVED***) map[string]interface***REMOVED******REMOVED*** ***REMOVED***
		for tag, kind := range fieldKinds ***REMOVED***
			if val, ok := tags[tag]; ok ***REMOVED***
				var v interface***REMOVED******REMOVED***
				var err error

				switch kind ***REMOVED***
				case String:
					v = val
				case Bool:
					v, err = strconv.ParseBool(val)
				case Float:
					v, err = strconv.ParseFloat(val, 64)
				case Int:
					v, err = strconv.ParseInt(val, 10, 64)
				***REMOVED***
				if err == nil ***REMOVED***
					values[tag] = v
				***REMOVED*** else ***REMOVED***
					values[tag] = val
				***REMOVED***

				delete(tags, tag)
			***REMOVED***
		***REMOVED***
		return values
	***REMOVED***
***REMOVED***

// makeFieldKinds reads the Config and returns a lookup map of tag names to
// the field type their values should be converted to.
func makeInfluxdbFieldKinds(tagsAsFields []string) (map[string]FieldKind, error) ***REMOVED***
	fieldKinds := make(map[string]FieldKind)
	for _, tag := range tagsAsFields ***REMOVED***
		var fieldName, fieldType string
		s := strings.SplitN(tag, ":", 2)
		if len(s) == 1 ***REMOVED***
			fieldName, fieldType = s[0], "string"
		***REMOVED*** else ***REMOVED***
			fieldName, fieldType = s[0], s[1]
		***REMOVED***

		err := checkDuplicatedTypeDefinitions(fieldKinds, fieldName)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		switch fieldType ***REMOVED***
		case "string":
			fieldKinds[fieldName] = String
		case "bool":
			fieldKinds[fieldName] = Bool
		case "float":
			fieldKinds[fieldName] = Float
		case "int":
			fieldKinds[fieldName] = Int
		default:
			return nil, fmt.Errorf("An invalid type (%s) is specified for an InfluxDB field (%s).",
				fieldType, fieldName)
		***REMOVED***
	***REMOVED***

	return fieldKinds, nil
***REMOVED***

func checkDuplicatedTypeDefinitions(fieldKinds map[string]FieldKind, tag string) error ***REMOVED***
	if _, found := fieldKinds[tag]; found ***REMOVED***
		return fmt.Errorf("A tag name (%s) shows up more than once in InfluxDB field type configurations.", tag)
	***REMOVED***
	return nil
***REMOVED***

func (c influxdbConfig) Apply(cfg influxdbConfig) influxdbConfig ***REMOVED***
	if len(cfg.TagsAsFields) > 0 ***REMOVED***
		c.TagsAsFields = cfg.TagsAsFields
	***REMOVED***
	return c
***REMOVED***

// ParseMap parses a map[string]interface***REMOVED******REMOVED*** into a Config
func influxdbParseMap(m map[string]interface***REMOVED******REMOVED***) (influxdbConfig, error) ***REMOVED***
	c := influxdbConfig***REMOVED******REMOVED***
	if v, ok := m["tagsAsFields"].(string); ok ***REMOVED***
		m["tagsAsFields"] = []string***REMOVED***v***REMOVED***
	***REMOVED***
	dec, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig***REMOVED***
		DecodeHook: types.NullDecoder,
		Result:     &c,
	***REMOVED***)
	if err != nil ***REMOVED***
		return c, err
	***REMOVED***

	err = dec.Decode(m)
	return c, err
***REMOVED***

type influxdbConfig struct ***REMOVED***
	TagsAsFields []string `json:"tagsAsFields,omitempty" envconfig:"K6_INFLUXDB_TAGS_AS_FIELDS"`
***REMOVED***

func newInfluxdbConfig() influxdbConfig ***REMOVED***
	c := influxdbConfig***REMOVED***
		TagsAsFields: []string***REMOVED***"vu", "iter", "url"***REMOVED***,
	***REMOVED***
	return c
***REMOVED***
