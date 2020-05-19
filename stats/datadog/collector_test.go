/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2019 Load Impact
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

package datadog

import (
	"strings"
	"testing"

	"github.com/loadimpact/k6/stats"
	"github.com/loadimpact/k6/stats/statsd/common"
	"github.com/loadimpact/k6/stats/statsd/common/testutil"
	"github.com/stretchr/testify/require"
)

func TestCollector(t *testing.T) ***REMOVED***
	var tagMap = stats.TagSet***REMOVED***"tag1": true, "tag2": true***REMOVED***
	var handler = tagHandler(tagMap)
	testutil.BaseTest(t, func(config common.Config) (*common.Collector, error) ***REMOVED***
		return New(NewConfig().Apply(Config***REMOVED***
			TagBlacklist: tagMap,
			Config:       config,
		***REMOVED***))
	***REMOVED***, func(t *testing.T, containers []stats.SampleContainer, expectedOutput, output string) ***REMOVED***
		var outputLines = strings.Split(output, "\n")
		var expectedOutputLines = strings.Split(expectedOutput, "\n")
		for i, container := range containers ***REMOVED***
			for j, sample := range container.GetSamples() ***REMOVED***
				var (
					expectedTagList    = handler.processTags(sample.GetTags().CloneTags())
					expectedOutputLine = expectedOutputLines[i*j+i]
					outputLine         = outputLines[i*j+i]
					outputWithoutTags  = outputLine
					outputTagList      = []string***REMOVED******REMOVED***
					tagSplit           = strings.LastIndex(outputLine, "|#")
				)

				if tagSplit != -1 ***REMOVED***
					outputWithoutTags = outputLine[:tagSplit]
					outputTagList = strings.Split(outputLine[tagSplit+len("|#"):], ",")
				***REMOVED***
				require.Equal(t, expectedOutputLine, outputWithoutTags)
				require.ElementsMatch(t, expectedTagList, outputTagList)
			***REMOVED***
		***REMOVED***
	***REMOVED***)
***REMOVED***
