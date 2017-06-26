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

package core

import (
	"testing"
	"time"

	"github.com/loadimpact/k6/lib"
	"github.com/stretchr/testify/assert"
	null "gopkg.in/guregu/null.v3"
)

func TestSumStages(t *testing.T) ***REMOVED***
	testdata := map[string]struct ***REMOVED***
		Time   lib.NullDuration
		Stages []lib.Stage
	***REMOVED******REMOVED***
		"Blank":    ***REMOVED***lib.NullDuration***REMOVED******REMOVED***, []lib.Stage***REMOVED******REMOVED******REMOVED***,
		"Infinite": ***REMOVED***lib.NullDuration***REMOVED******REMOVED***, []lib.Stage***REMOVED******REMOVED******REMOVED******REMOVED******REMOVED***,
		"Limit": ***REMOVED***
			lib.NullDurationFrom(10 * time.Second),
			[]lib.Stage***REMOVED***
				***REMOVED***Duration: lib.NullDurationFrom(5 * time.Second)***REMOVED***,
				***REMOVED***Duration: lib.NullDurationFrom(5 * time.Second)***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		"InfiniteTail": ***REMOVED***
			lib.NullDuration***REMOVED***Duration: lib.Duration(10 * time.Second), Valid: false***REMOVED***,
			[]lib.Stage***REMOVED***
				***REMOVED***Duration: lib.NullDurationFrom(5 * time.Second)***REMOVED***,
				***REMOVED***Duration: lib.NullDurationFrom(5 * time.Second)***REMOVED***,
				***REMOVED******REMOVED***,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***
	for name, data := range testdata ***REMOVED***
		t.Run(name, func(t *testing.T) ***REMOVED***
			assert.Equal(t, data.Time, SumStages(data.Stages))
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestProcessStages(t *testing.T) ***REMOVED***
	type checkpoint struct ***REMOVED***
		D    time.Duration
		Keep bool
		VUs  int64
	***REMOVED***
	testdata := map[string]struct ***REMOVED***
		Stages      []lib.Stage
		Checkpoints []checkpoint
	***REMOVED******REMOVED***
		"none": ***REMOVED***
			[]lib.Stage***REMOVED******REMOVED***,
			[]checkpoint***REMOVED***
				***REMOVED***0 * time.Second, false, 0***REMOVED***,
				***REMOVED***10 * time.Second, false, 0***REMOVED***,
				***REMOVED***24 * time.Hour, false, 0***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		"one": ***REMOVED***
			[]lib.Stage***REMOVED***
				***REMOVED***Duration: lib.NullDurationFrom(10 * time.Second)***REMOVED***,
			***REMOVED***,
			[]checkpoint***REMOVED***
				***REMOVED***0 * time.Second, true, 0***REMOVED***,
				***REMOVED***1 * time.Second, true, 0***REMOVED***,
				***REMOVED***10 * time.Second, true, 0***REMOVED***,
				***REMOVED***11 * time.Second, false, 0***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		"one/targeted": ***REMOVED***
			[]lib.Stage***REMOVED***
				***REMOVED***Duration: lib.NullDurationFrom(10 * time.Second), Target: null.IntFrom(100)***REMOVED***,
			***REMOVED***,
			[]checkpoint***REMOVED***
				***REMOVED***0 * time.Second, true, 0***REMOVED***,
				***REMOVED***1 * time.Second, true, 10***REMOVED***,
				***REMOVED***2 * time.Second, true, 20***REMOVED***,
				***REMOVED***3 * time.Second, true, 30***REMOVED***,
				***REMOVED***4 * time.Second, true, 40***REMOVED***,
				***REMOVED***5 * time.Second, true, 50***REMOVED***,
				***REMOVED***6 * time.Second, true, 60***REMOVED***,
				***REMOVED***7 * time.Second, true, 70***REMOVED***,
				***REMOVED***8 * time.Second, true, 80***REMOVED***,
				***REMOVED***9 * time.Second, true, 90***REMOVED***,
				***REMOVED***10 * time.Second, true, 100***REMOVED***,
				***REMOVED***11 * time.Second, false, 100***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		"two": ***REMOVED***
			[]lib.Stage***REMOVED***
				***REMOVED***Duration: lib.NullDurationFrom(5 * time.Second)***REMOVED***,
				***REMOVED***Duration: lib.NullDurationFrom(5 * time.Second)***REMOVED***,
			***REMOVED***,
			[]checkpoint***REMOVED***
				***REMOVED***0 * time.Second, true, 0***REMOVED***,
				***REMOVED***1 * time.Second, true, 0***REMOVED***,
				***REMOVED***11 * time.Second, false, 0***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		"two/targeted": ***REMOVED***
			[]lib.Stage***REMOVED***
				***REMOVED***Duration: lib.NullDurationFrom(5 * time.Second), Target: null.IntFrom(100)***REMOVED***,
				***REMOVED***Duration: lib.NullDurationFrom(5 * time.Second), Target: null.IntFrom(0)***REMOVED***,
			***REMOVED***,
			[]checkpoint***REMOVED***
				***REMOVED***0 * time.Second, true, 0***REMOVED***,
				***REMOVED***1 * time.Second, true, 20***REMOVED***,
				***REMOVED***2 * time.Second, true, 40***REMOVED***,
				***REMOVED***3 * time.Second, true, 60***REMOVED***,
				***REMOVED***4 * time.Second, true, 80***REMOVED***,
				***REMOVED***5 * time.Second, true, 100***REMOVED***,
				***REMOVED***6 * time.Second, true, 80***REMOVED***,
				***REMOVED***7 * time.Second, true, 60***REMOVED***,
				***REMOVED***8 * time.Second, true, 40***REMOVED***,
				***REMOVED***9 * time.Second, true, 20***REMOVED***,
				***REMOVED***10 * time.Second, true, 0***REMOVED***,
				***REMOVED***11 * time.Second, false, 0***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		"three": ***REMOVED***
			[]lib.Stage***REMOVED***
				***REMOVED***Duration: lib.NullDurationFrom(5 * time.Second)***REMOVED***,
				***REMOVED***Duration: lib.NullDurationFrom(10 * time.Second)***REMOVED***,
				***REMOVED***Duration: lib.NullDurationFrom(15 * time.Second)***REMOVED***,
			***REMOVED***,
			[]checkpoint***REMOVED***
				***REMOVED***0 * time.Second, true, 0***REMOVED***,
				***REMOVED***1 * time.Second, true, 0***REMOVED***,
				***REMOVED***15 * time.Second, true, 0***REMOVED***,
				***REMOVED***30 * time.Second, true, 0***REMOVED***,
				***REMOVED***31 * time.Second, false, 0***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		"three/targeted": ***REMOVED***
			[]lib.Stage***REMOVED***
				***REMOVED***Duration: lib.NullDurationFrom(5 * time.Second), Target: null.IntFrom(50)***REMOVED***,
				***REMOVED***Duration: lib.NullDurationFrom(5 * time.Second), Target: null.IntFrom(100)***REMOVED***,
				***REMOVED***Duration: lib.NullDurationFrom(5 * time.Second), Target: null.IntFrom(0)***REMOVED***,
			***REMOVED***,
			[]checkpoint***REMOVED***
				***REMOVED***0 * time.Second, true, 0***REMOVED***,
				***REMOVED***1 * time.Second, true, 10***REMOVED***,
				***REMOVED***2 * time.Second, true, 20***REMOVED***,
				***REMOVED***3 * time.Second, true, 30***REMOVED***,
				***REMOVED***4 * time.Second, true, 40***REMOVED***,
				***REMOVED***5 * time.Second, true, 50***REMOVED***,
				***REMOVED***6 * time.Second, true, 60***REMOVED***,
				***REMOVED***7 * time.Second, true, 70***REMOVED***,
				***REMOVED***8 * time.Second, true, 80***REMOVED***,
				***REMOVED***9 * time.Second, true, 90***REMOVED***,
				***REMOVED***10 * time.Second, true, 100***REMOVED***,
				***REMOVED***11 * time.Second, true, 80***REMOVED***,
				***REMOVED***12 * time.Second, true, 60***REMOVED***,
				***REMOVED***13 * time.Second, true, 40***REMOVED***,
				***REMOVED***14 * time.Second, true, 20***REMOVED***,
				***REMOVED***15 * time.Second, true, 0***REMOVED***,
				***REMOVED***16 * time.Second, false, 0***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		"mix": ***REMOVED***
			[]lib.Stage***REMOVED***
				***REMOVED***Duration: lib.NullDurationFrom(5 * time.Second), Target: null.IntFrom(20)***REMOVED***,
				***REMOVED***Duration: lib.NullDurationFrom(5 * time.Second), Target: null.IntFrom(10)***REMOVED***,
				***REMOVED***Duration: lib.NullDurationFrom(2 * time.Second)***REMOVED***,
				***REMOVED***Duration: lib.NullDurationFrom(5 * time.Second), Target: null.IntFrom(20)***REMOVED***,
				***REMOVED***Duration: lib.NullDurationFrom(2 * time.Second)***REMOVED***,
				***REMOVED***Duration: lib.NullDurationFrom(5 * time.Second), Target: null.IntFrom(10)***REMOVED***,
			***REMOVED***,
			[]checkpoint***REMOVED***
				***REMOVED***0 * time.Second, true, 0***REMOVED***,

				***REMOVED***1 * time.Second, true, 4***REMOVED***,
				***REMOVED***2 * time.Second, true, 8***REMOVED***,
				***REMOVED***3 * time.Second, true, 12***REMOVED***,
				***REMOVED***4 * time.Second, true, 16***REMOVED***,
				***REMOVED***5 * time.Second, true, 20***REMOVED***,

				***REMOVED***6 * time.Second, true, 18***REMOVED***,
				***REMOVED***7 * time.Second, true, 16***REMOVED***,
				***REMOVED***8 * time.Second, true, 14***REMOVED***,
				***REMOVED***9 * time.Second, true, 12***REMOVED***,
				***REMOVED***10 * time.Second, true, 10***REMOVED***,

				***REMOVED***11 * time.Second, true, 10***REMOVED***,
				***REMOVED***12 * time.Second, true, 10***REMOVED***,

				***REMOVED***13 * time.Second, true, 12***REMOVED***,
				***REMOVED***14 * time.Second, true, 14***REMOVED***,
				***REMOVED***15 * time.Second, true, 16***REMOVED***,
				***REMOVED***16 * time.Second, true, 18***REMOVED***,
				***REMOVED***17 * time.Second, true, 20***REMOVED***,

				***REMOVED***18 * time.Second, true, 20***REMOVED***,
				***REMOVED***19 * time.Second, true, 20***REMOVED***,

				***REMOVED***20 * time.Second, true, 18***REMOVED***,
				***REMOVED***21 * time.Second, true, 16***REMOVED***,
				***REMOVED***22 * time.Second, true, 14***REMOVED***,
				***REMOVED***23 * time.Second, true, 12***REMOVED***,
				***REMOVED***24 * time.Second, true, 10***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		"infinite": ***REMOVED***
			[]lib.Stage***REMOVED******REMOVED******REMOVED******REMOVED***,
			[]checkpoint***REMOVED***
				***REMOVED***0 * time.Second, true, 0***REMOVED***,
				***REMOVED***1 * time.Minute, true, 0***REMOVED***,
				***REMOVED***1 * time.Hour, true, 0***REMOVED***,
				***REMOVED***24 * time.Hour, true, 0***REMOVED***,
				***REMOVED***365 * 24 * time.Hour, true, 0***REMOVED***,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***
	for name, data := range testdata ***REMOVED***
		t.Run(name, func(t *testing.T) ***REMOVED***
			for _, ckp := range data.Checkpoints ***REMOVED***
				t.Run(ckp.D.String(), func(t *testing.T) ***REMOVED***
					vus, keepRunning := ProcessStages(data.Stages, ckp.D)
					assert.Equal(t, ckp.VUs, vus)
					assert.Equal(t, ckp.Keep, keepRunning)
				***REMOVED***)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***
