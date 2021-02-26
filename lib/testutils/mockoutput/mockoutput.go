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

package mockoutput

import (
	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/output"
	"github.com/loadimpact/k6/stats"
)

// New exists so that the usage from tests avoids repetition, i.e. is
// mockoutput.New() instead of &mockoutput.MockOutput***REMOVED******REMOVED***
func New() *MockOutput ***REMOVED***
	return &MockOutput***REMOVED******REMOVED***
***REMOVED***

// MockOutput can be used in tests to mock an actual output.
type MockOutput struct ***REMOVED***
	SampleContainers []stats.SampleContainer
	Samples          []stats.Sample
	RunStatus        lib.RunStatus

	DescFn  func() string
	StartFn func() error
	StopFn  func() error
***REMOVED***

var _ output.WithRunStatusUpdates = &MockOutput***REMOVED******REMOVED***

// AddMetricSamples just saves the results in memory.
func (mo *MockOutput) AddMetricSamples(scs []stats.SampleContainer) ***REMOVED***
	mo.SampleContainers = append(mo.SampleContainers, scs...)
	for _, sc := range scs ***REMOVED***
		mo.Samples = append(mo.Samples, sc.GetSamples()...)
	***REMOVED***
***REMOVED***

// SetRunStatus updates the RunStatus property.
func (mo *MockOutput) SetRunStatus(latestStatus lib.RunStatus) ***REMOVED***
	mo.RunStatus = latestStatus
***REMOVED***

// Description calls the supplied DescFn callback, if available.
func (mo *MockOutput) Description() string ***REMOVED***
	if mo.DescFn != nil ***REMOVED***
		return mo.DescFn()
	***REMOVED***
	return "mock"
***REMOVED***

// Start calls the supplied StartFn callback, if available.
func (mo *MockOutput) Start() error ***REMOVED***
	if mo.StartFn != nil ***REMOVED***
		return mo.StartFn()
	***REMOVED***
	return nil
***REMOVED***

// Stop calls the supplied StopFn callback, if available.
func (mo *MockOutput) Stop() error ***REMOVED***
	if mo.StopFn != nil ***REMOVED***
		return mo.StopFn()
	***REMOVED***
	return nil
***REMOVED***
