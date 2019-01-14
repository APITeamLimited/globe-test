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

package stats

import (
	"encoding/json"
	"time"

	"github.com/dop251/goja"
	"github.com/loadimpact/k6/lib/types"
	"github.com/pkg/errors"
)

const jsEnvSrc = `
function p(pct) ***REMOVED***
	return __sink__.P(pct/100.0);
***REMOVED***;
`

var jsEnv *goja.Program

func init() ***REMOVED***
	pgm, err := goja.Compile("__env__", jsEnvSrc, true)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	jsEnv = pgm
***REMOVED***

type Threshold struct ***REMOVED***
	Source           string
	Failed           bool
	AbortOnFail      bool
	AbortGracePeriod types.NullDuration

	pgm *goja.Program
	rt  *goja.Runtime
***REMOVED***

func NewThreshold(src string, rt *goja.Runtime, abortOnFail bool, gracePeriod types.NullDuration) (*Threshold, error) ***REMOVED***
	pgm, err := goja.Compile("__threshold__", src, true)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return &Threshold***REMOVED***
		Source:           src,
		AbortOnFail:      abortOnFail,
		AbortGracePeriod: gracePeriod,
		pgm:              pgm,
		rt:               rt,
	***REMOVED***, nil
***REMOVED***

func (t Threshold) RunNoTaint() (bool, error) ***REMOVED***
	v, err := t.rt.RunProgram(t.pgm)
	if err != nil ***REMOVED***
		return false, err
	***REMOVED***
	return v.ToBoolean(), nil
***REMOVED***

func (t *Threshold) Run() (bool, error) ***REMOVED***
	b, err := t.RunNoTaint()
	t.Failed = !b
	return b, err
***REMOVED***

type ThresholdConfig struct ***REMOVED***
	Threshold        string             `json:"threshold"`
	AbortOnFail      bool               `json:"abortOnFail"`
	AbortGracePeriod types.NullDuration `json:"delayAbortEval"`
***REMOVED***

//used internally for JSON marshalling
type rawThresholdConfig ThresholdConfig

func (tc *ThresholdConfig) UnmarshalJSON(data []byte) error ***REMOVED***
	//shortcircuit unmarshalling for simple string format
	if err := json.Unmarshal(data, &tc.Threshold); err == nil ***REMOVED***
		return nil
	***REMOVED***

	rawConfig := (*rawThresholdConfig)(tc)
	return json.Unmarshal(data, rawConfig)
***REMOVED***

func (tc ThresholdConfig) MarshalJSON() ([]byte, error) ***REMOVED***
	if tc.AbortOnFail ***REMOVED***
		return json.Marshal(rawThresholdConfig(tc))
	***REMOVED***
	return json.Marshal(tc.Threshold)
***REMOVED***

type Thresholds struct ***REMOVED***
	Runtime    *goja.Runtime
	Thresholds []*Threshold
	Abort      bool
***REMOVED***

func NewThresholds(sources []string) (Thresholds, error) ***REMOVED***
	tcs := make([]ThresholdConfig, len(sources))
	for i, source := range sources ***REMOVED***
		tcs[i].Threshold = source
	***REMOVED***

	return NewThresholdsWithConfig(tcs)
***REMOVED***

func NewThresholdsWithConfig(configs []ThresholdConfig) (Thresholds, error) ***REMOVED***
	rt := goja.New()
	if _, err := rt.RunProgram(jsEnv); err != nil ***REMOVED***
		return Thresholds***REMOVED******REMOVED***, errors.Wrap(err, "builtin")
	***REMOVED***

	ts := make([]*Threshold, len(configs))
	for i, config := range configs ***REMOVED***
		t, err := NewThreshold(config.Threshold, rt, config.AbortOnFail, config.AbortGracePeriod)
		if err != nil ***REMOVED***
			return Thresholds***REMOVED******REMOVED***, errors.Wrapf(err, "%d", i)
		***REMOVED***
		ts[i] = t
	***REMOVED***

	return Thresholds***REMOVED***rt, ts, false***REMOVED***, nil
***REMOVED***

func (ts *Thresholds) UpdateVM(sink Sink, t time.Duration) error ***REMOVED***
	ts.Runtime.Set("__sink__", sink)
	f := sink.Format(t)
	for k, v := range f ***REMOVED***
		ts.Runtime.Set(k, v)
	***REMOVED***
	return nil
***REMOVED***

func (ts *Thresholds) RunAll(t time.Duration) (bool, error) ***REMOVED***
	succ := true
	for i, th := range ts.Thresholds ***REMOVED***
		b, err := th.Run()
		if err != nil ***REMOVED***
			return false, errors.Wrapf(err, "%d", i)
		***REMOVED***
		if !b ***REMOVED***
			succ = false

			if ts.Abort || !th.AbortOnFail ***REMOVED***
				continue
			***REMOVED***

			ts.Abort = !th.AbortGracePeriod.Valid ||
				th.AbortGracePeriod.Duration < types.Duration(t)
		***REMOVED***
	***REMOVED***
	return succ, nil
***REMOVED***

func (ts *Thresholds) Run(sink Sink, t time.Duration) (bool, error) ***REMOVED***
	if err := ts.UpdateVM(sink, t); err != nil ***REMOVED***
		return false, err
	***REMOVED***
	return ts.RunAll(t)
***REMOVED***

func (ts *Thresholds) UnmarshalJSON(data []byte) error ***REMOVED***
	var configs []ThresholdConfig
	if err := json.Unmarshal(data, &configs); err != nil ***REMOVED***
		return err
	***REMOVED***
	newts, err := NewThresholdsWithConfig(configs)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	*ts = newts
	return nil
***REMOVED***

func (ts Thresholds) MarshalJSON() ([]byte, error) ***REMOVED***
	configs := make([]ThresholdConfig, len(ts.Thresholds))
	for i, t := range ts.Thresholds ***REMOVED***
		configs[i].Threshold = t.Source
		configs[i].AbortOnFail = t.AbortOnFail
		configs[i].AbortGracePeriod = t.AbortGracePeriod
	***REMOVED***
	return json.Marshal(configs)
***REMOVED***
