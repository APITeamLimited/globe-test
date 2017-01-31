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

package lib

import (
	"encoding/json"
	"github.com/loadimpact/k6/stats"
	"github.com/pkg/errors"
	"github.com/robertkrimen/otto"
)

const jsEnv = `
function p(pct) ***REMOVED***
	return __sink__.P(pct);
***REMOVED***;
`

type Threshold struct ***REMOVED***
	Source string
	Failed bool

	script *otto.Script
	vm     *otto.Otto
***REMOVED***

func NewThreshold(src string, vm *otto.Otto) (*Threshold, error) ***REMOVED***
	script, err := vm.Compile("__threshold__", src)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return &Threshold***REMOVED***
		Source: src,
		script: script,
		vm:     vm,
	***REMOVED***, nil
***REMOVED***

func (t Threshold) RunNoTaint() (bool, error) ***REMOVED***
	v, err := t.vm.Run(t.script)
	if err != nil ***REMOVED***
		return false, err
	***REMOVED***
	return v.ToBoolean()
***REMOVED***

func (t *Threshold) Run() (bool, error) ***REMOVED***
	b, err := t.RunNoTaint()
	if !b ***REMOVED***
		t.Failed = true
	***REMOVED***
	return b, err
***REMOVED***

type Thresholds struct ***REMOVED***
	VM         *otto.Otto
	Thresholds []*Threshold
***REMOVED***

func NewThresholds(sources []string) (Thresholds, error) ***REMOVED***
	vm := otto.New()

	if _, err := vm.Eval(jsEnv); err != nil ***REMOVED***
		return Thresholds***REMOVED******REMOVED***, errors.Wrap(err, "builtin")
	***REMOVED***

	ts := make([]*Threshold, len(sources))
	for i, src := range sources ***REMOVED***
		t, err := NewThreshold(src, vm)
		if err != nil ***REMOVED***
			return Thresholds***REMOVED******REMOVED***, errors.Wrapf(err, "%d", i)
		***REMOVED***
		ts[i] = t
	***REMOVED***
	return Thresholds***REMOVED***vm, ts***REMOVED***, nil
***REMOVED***

func (ts *Thresholds) UpdateVM(sink stats.Sink) error ***REMOVED***
	if err := ts.VM.Set("__sink__", sink); err != nil ***REMOVED***
		return err
	***REMOVED***
	for k, v := range sink.Format() ***REMOVED***
		if err := ts.VM.Set(k, v); err != nil ***REMOVED***
			return errors.Wrapf(err, "%s", k)
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (ts *Thresholds) RunAll() (bool, error) ***REMOVED***
	succ := true
	for i, th := range ts.Thresholds ***REMOVED***
		b, err := th.Run()
		if err != nil ***REMOVED***
			return false, errors.Wrapf(err, "%d", i)
		***REMOVED***
		if !b ***REMOVED***
			succ = false
		***REMOVED***
	***REMOVED***
	return succ, nil
***REMOVED***

func (ts *Thresholds) Run(sink stats.Sink) (bool, error) ***REMOVED***
	if err := ts.UpdateVM(sink); err != nil ***REMOVED***
		return false, err
	***REMOVED***
	return ts.RunAll()
***REMOVED***

func (ts *Thresholds) UnmarshalJSON(data []byte) error ***REMOVED***
	var sources []string
	if err := json.Unmarshal(data, &sources); err != nil ***REMOVED***
		return err
	***REMOVED***

	newts, err := NewThresholds(sources)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	*ts = newts
	return nil
***REMOVED***

func (ts Thresholds) MarshalJSON() ([]byte, error) ***REMOVED***
	sources := make([]string, len(ts.Thresholds))
	for i, t := range ts.Thresholds ***REMOVED***
		sources[i] = t.Source
	***REMOVED***
	return json.Marshal(sources)
***REMOVED***
