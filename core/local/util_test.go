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

package local

//TODO: translate this test to the new paradigm
/*
func TestProcessStages(t *testing.T) ***REMOVED***
	type checkpoint struct ***REMOVED***
		D    time.Duration
		Keep bool
		VUs  null.Int
	***REMOVED***
	testdata := map[string]struct ***REMOVED***
		Start       int64
		Stages      []lib.Stage
		Checkpoints []checkpoint
	***REMOVED******REMOVED***
		"none": ***REMOVED***
			0,
			[]lib.Stage***REMOVED******REMOVED***,
			[]checkpoint***REMOVED***
				***REMOVED***0 * time.Second, false, null.NewInt(0, false)***REMOVED***,
				***REMOVED***10 * time.Second, false, null.NewInt(0, false)***REMOVED***,
				***REMOVED***24 * time.Hour, false, null.NewInt(0, false)***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		"one": ***REMOVED***
			0,
			[]lib.Stage***REMOVED***
				***REMOVED***Duration: types.NullDurationFrom(10 * time.Second)***REMOVED***,
			***REMOVED***,
			[]checkpoint***REMOVED***
				***REMOVED***0 * time.Second, true, null.NewInt(0, false)***REMOVED***,
				***REMOVED***1 * time.Second, true, null.NewInt(0, false)***REMOVED***,
				***REMOVED***10 * time.Second, true, null.NewInt(0, false)***REMOVED***,
				***REMOVED***11 * time.Second, false, null.NewInt(0, false)***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		"one/start": ***REMOVED***
			5,
			[]lib.Stage***REMOVED***
				***REMOVED***Duration: types.NullDurationFrom(10 * time.Second)***REMOVED***,
			***REMOVED***,
			[]checkpoint***REMOVED***
				***REMOVED***0 * time.Second, true, null.NewInt(5, false)***REMOVED***,
				***REMOVED***1 * time.Second, true, null.NewInt(5, false)***REMOVED***,
				***REMOVED***10 * time.Second, true, null.NewInt(5, false)***REMOVED***,
				***REMOVED***11 * time.Second, false, null.NewInt(5, false)***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		"one/targeted": ***REMOVED***
			0,
			[]lib.Stage***REMOVED***
				***REMOVED***Duration: types.NullDurationFrom(10 * time.Second), Target: null.IntFrom(100)***REMOVED***,
			***REMOVED***,
			[]checkpoint***REMOVED***
				***REMOVED***0 * time.Second, true, null.IntFrom(0)***REMOVED***,
				***REMOVED***1 * time.Second, true, null.IntFrom(10)***REMOVED***,
				***REMOVED***2 * time.Second, true, null.IntFrom(20)***REMOVED***,
				***REMOVED***3 * time.Second, true, null.IntFrom(30)***REMOVED***,
				***REMOVED***4 * time.Second, true, null.IntFrom(40)***REMOVED***,
				***REMOVED***5 * time.Second, true, null.IntFrom(50)***REMOVED***,
				***REMOVED***6 * time.Second, true, null.IntFrom(60)***REMOVED***,
				***REMOVED***7 * time.Second, true, null.IntFrom(70)***REMOVED***,
				***REMOVED***8 * time.Second, true, null.IntFrom(80)***REMOVED***,
				***REMOVED***9 * time.Second, true, null.IntFrom(90)***REMOVED***,
				***REMOVED***10 * time.Second, true, null.IntFrom(100)***REMOVED***,
				***REMOVED***11 * time.Second, false, null.IntFrom(100)***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		"one/targeted/start": ***REMOVED***
			50,
			[]lib.Stage***REMOVED***
				***REMOVED***Duration: types.NullDurationFrom(10 * time.Second), Target: null.IntFrom(100)***REMOVED***,
			***REMOVED***,
			[]checkpoint***REMOVED***
				***REMOVED***0 * time.Second, true, null.IntFrom(50)***REMOVED***,
				***REMOVED***1 * time.Second, true, null.IntFrom(55)***REMOVED***,
				***REMOVED***2 * time.Second, true, null.IntFrom(60)***REMOVED***,
				***REMOVED***3 * time.Second, true, null.IntFrom(65)***REMOVED***,
				***REMOVED***4 * time.Second, true, null.IntFrom(70)***REMOVED***,
				***REMOVED***5 * time.Second, true, null.IntFrom(75)***REMOVED***,
				***REMOVED***6 * time.Second, true, null.IntFrom(80)***REMOVED***,
				***REMOVED***7 * time.Second, true, null.IntFrom(85)***REMOVED***,
				***REMOVED***8 * time.Second, true, null.IntFrom(90)***REMOVED***,
				***REMOVED***9 * time.Second, true, null.IntFrom(95)***REMOVED***,
				***REMOVED***10 * time.Second, true, null.IntFrom(100)***REMOVED***,
				***REMOVED***11 * time.Second, false, null.IntFrom(100)***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		"two": ***REMOVED***
			0,
			[]lib.Stage***REMOVED***
				***REMOVED***Duration: types.NullDurationFrom(5 * time.Second)***REMOVED***,
				***REMOVED***Duration: types.NullDurationFrom(5 * time.Second)***REMOVED***,
			***REMOVED***,
			[]checkpoint***REMOVED***
				***REMOVED***0 * time.Second, true, null.NewInt(0, false)***REMOVED***,
				***REMOVED***1 * time.Second, true, null.NewInt(0, false)***REMOVED***,
				***REMOVED***11 * time.Second, false, null.NewInt(0, false)***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		"two/start": ***REMOVED***
			5,
			[]lib.Stage***REMOVED***
				***REMOVED***Duration: types.NullDurationFrom(5 * time.Second)***REMOVED***,
				***REMOVED***Duration: types.NullDurationFrom(5 * time.Second)***REMOVED***,
			***REMOVED***,
			[]checkpoint***REMOVED***
				***REMOVED***0 * time.Second, true, null.NewInt(5, false)***REMOVED***,
				***REMOVED***1 * time.Second, true, null.NewInt(5, false)***REMOVED***,
				***REMOVED***11 * time.Second, false, null.NewInt(5, false)***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		"two/targeted": ***REMOVED***
			0,
			[]lib.Stage***REMOVED***
				***REMOVED***Duration: types.NullDurationFrom(5 * time.Second), Target: null.IntFrom(100)***REMOVED***,
				***REMOVED***Duration: types.NullDurationFrom(5 * time.Second), Target: null.IntFrom(0)***REMOVED***,
			***REMOVED***,
			[]checkpoint***REMOVED***
				***REMOVED***0 * time.Second, true, null.IntFrom(0)***REMOVED***,
				***REMOVED***1 * time.Second, true, null.IntFrom(20)***REMOVED***,
				***REMOVED***2 * time.Second, true, null.IntFrom(40)***REMOVED***,
				***REMOVED***3 * time.Second, true, null.IntFrom(60)***REMOVED***,
				***REMOVED***4 * time.Second, true, null.IntFrom(80)***REMOVED***,
				***REMOVED***5 * time.Second, true, null.IntFrom(100)***REMOVED***,
				***REMOVED***6 * time.Second, true, null.IntFrom(80)***REMOVED***,
				***REMOVED***7 * time.Second, true, null.IntFrom(60)***REMOVED***,
				***REMOVED***8 * time.Second, true, null.IntFrom(40)***REMOVED***,
				***REMOVED***9 * time.Second, true, null.IntFrom(20)***REMOVED***,
				***REMOVED***10 * time.Second, true, null.IntFrom(0)***REMOVED***,
				***REMOVED***11 * time.Second, false, null.IntFrom(0)***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		"three": ***REMOVED***
			0,
			[]lib.Stage***REMOVED***
				***REMOVED***Duration: types.NullDurationFrom(5 * time.Second)***REMOVED***,
				***REMOVED***Duration: types.NullDurationFrom(10 * time.Second)***REMOVED***,
				***REMOVED***Duration: types.NullDurationFrom(15 * time.Second)***REMOVED***,
			***REMOVED***,
			[]checkpoint***REMOVED***
				***REMOVED***0 * time.Second, true, null.NewInt(0, false)***REMOVED***,
				***REMOVED***1 * time.Second, true, null.NewInt(0, false)***REMOVED***,
				***REMOVED***15 * time.Second, true, null.NewInt(0, false)***REMOVED***,
				***REMOVED***30 * time.Second, true, null.NewInt(0, false)***REMOVED***,
				***REMOVED***31 * time.Second, false, null.NewInt(0, false)***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		"three/targeted": ***REMOVED***
			0,
			[]lib.Stage***REMOVED***
				***REMOVED***Duration: types.NullDurationFrom(5 * time.Second), Target: null.IntFrom(50)***REMOVED***,
				***REMOVED***Duration: types.NullDurationFrom(5 * time.Second), Target: null.IntFrom(100)***REMOVED***,
				***REMOVED***Duration: types.NullDurationFrom(5 * time.Second), Target: null.IntFrom(0)***REMOVED***,
			***REMOVED***,
			[]checkpoint***REMOVED***
				***REMOVED***0 * time.Second, true, null.IntFrom(0)***REMOVED***,
				***REMOVED***1 * time.Second, true, null.IntFrom(10)***REMOVED***,
				***REMOVED***2 * time.Second, true, null.IntFrom(20)***REMOVED***,
				***REMOVED***3 * time.Second, true, null.IntFrom(30)***REMOVED***,
				***REMOVED***4 * time.Second, true, null.IntFrom(40)***REMOVED***,
				***REMOVED***5 * time.Second, true, null.IntFrom(50)***REMOVED***,
				***REMOVED***6 * time.Second, true, null.IntFrom(60)***REMOVED***,
				***REMOVED***7 * time.Second, true, null.IntFrom(70)***REMOVED***,
				***REMOVED***8 * time.Second, true, null.IntFrom(80)***REMOVED***,
				***REMOVED***9 * time.Second, true, null.IntFrom(90)***REMOVED***,
				***REMOVED***10 * time.Second, true, null.IntFrom(100)***REMOVED***,
				***REMOVED***11 * time.Second, true, null.IntFrom(80)***REMOVED***,
				***REMOVED***12 * time.Second, true, null.IntFrom(60)***REMOVED***,
				***REMOVED***13 * time.Second, true, null.IntFrom(40)***REMOVED***,
				***REMOVED***14 * time.Second, true, null.IntFrom(20)***REMOVED***,
				***REMOVED***15 * time.Second, true, null.IntFrom(0)***REMOVED***,
				***REMOVED***16 * time.Second, false, null.IntFrom(0)***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		"mix": ***REMOVED***
			0,
			[]lib.Stage***REMOVED***
				***REMOVED***Duration: types.NullDurationFrom(5 * time.Second), Target: null.IntFrom(20)***REMOVED***,
				***REMOVED***Duration: types.NullDurationFrom(5 * time.Second), Target: null.IntFrom(10)***REMOVED***,
				***REMOVED***Duration: types.NullDurationFrom(2 * time.Second)***REMOVED***,
				***REMOVED***Duration: types.NullDurationFrom(5 * time.Second), Target: null.IntFrom(20)***REMOVED***,
				***REMOVED***Duration: types.NullDurationFrom(2 * time.Second)***REMOVED***,
				***REMOVED***Duration: types.NullDurationFrom(5 * time.Second), Target: null.IntFrom(10)***REMOVED***,
			***REMOVED***,
			[]checkpoint***REMOVED***
				***REMOVED***0 * time.Second, true, null.IntFrom(0)***REMOVED***,

				***REMOVED***1 * time.Second, true, null.IntFrom(4)***REMOVED***,
				***REMOVED***2 * time.Second, true, null.IntFrom(8)***REMOVED***,
				***REMOVED***3 * time.Second, true, null.IntFrom(12)***REMOVED***,
				***REMOVED***4 * time.Second, true, null.IntFrom(16)***REMOVED***,
				***REMOVED***5 * time.Second, true, null.IntFrom(20)***REMOVED***,

				***REMOVED***6 * time.Second, true, null.IntFrom(18)***REMOVED***,
				***REMOVED***7 * time.Second, true, null.IntFrom(16)***REMOVED***,
				***REMOVED***8 * time.Second, true, null.IntFrom(14)***REMOVED***,
				***REMOVED***9 * time.Second, true, null.IntFrom(12)***REMOVED***,
				***REMOVED***10 * time.Second, true, null.IntFrom(10)***REMOVED***,

				***REMOVED***11 * time.Second, true, null.IntFrom(10)***REMOVED***,
				***REMOVED***12 * time.Second, true, null.IntFrom(10)***REMOVED***,

				***REMOVED***13 * time.Second, true, null.IntFrom(12)***REMOVED***,
				***REMOVED***14 * time.Second, true, null.IntFrom(14)***REMOVED***,
				***REMOVED***15 * time.Second, true, null.IntFrom(16)***REMOVED***,
				***REMOVED***16 * time.Second, true, null.IntFrom(18)***REMOVED***,
				***REMOVED***17 * time.Second, true, null.IntFrom(20)***REMOVED***,

				***REMOVED***18 * time.Second, true, null.IntFrom(20)***REMOVED***,
				***REMOVED***19 * time.Second, true, null.IntFrom(20)***REMOVED***,

				***REMOVED***20 * time.Second, true, null.IntFrom(18)***REMOVED***,
				***REMOVED***21 * time.Second, true, null.IntFrom(16)***REMOVED***,
				***REMOVED***22 * time.Second, true, null.IntFrom(14)***REMOVED***,
				***REMOVED***23 * time.Second, true, null.IntFrom(12)***REMOVED***,
				***REMOVED***24 * time.Second, true, null.IntFrom(10)***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		"mix/start": ***REMOVED***
			5,
			[]lib.Stage***REMOVED***
				***REMOVED***Duration: types.NullDurationFrom(5 * time.Second)***REMOVED***,
				***REMOVED***Duration: types.NullDurationFrom(5 * time.Second), Target: null.IntFrom(10)***REMOVED***,
			***REMOVED***,
			[]checkpoint***REMOVED***
				***REMOVED***0 * time.Second, true, null.NewInt(5, false)***REMOVED***,

				***REMOVED***1 * time.Second, true, null.NewInt(5, false)***REMOVED***,
				***REMOVED***2 * time.Second, true, null.NewInt(5, false)***REMOVED***,
				***REMOVED***3 * time.Second, true, null.NewInt(5, false)***REMOVED***,
				***REMOVED***4 * time.Second, true, null.NewInt(5, false)***REMOVED***,
				***REMOVED***5 * time.Second, true, null.NewInt(5, false)***REMOVED***,

				***REMOVED***6 * time.Second, true, null.NewInt(6, true)***REMOVED***,
				***REMOVED***7 * time.Second, true, null.NewInt(7, true)***REMOVED***,
				***REMOVED***8 * time.Second, true, null.NewInt(8, true)***REMOVED***,
				***REMOVED***9 * time.Second, true, null.NewInt(9, true)***REMOVED***,
				***REMOVED***10 * time.Second, true, null.NewInt(10, true)***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		"infinite": ***REMOVED***
			0,
			[]lib.Stage***REMOVED******REMOVED******REMOVED******REMOVED***,
			[]checkpoint***REMOVED***
				***REMOVED***0 * time.Second, true, null.NewInt(0, false)***REMOVED***,
				***REMOVED***1 * time.Minute, true, null.NewInt(0, false)***REMOVED***,
				***REMOVED***1 * time.Hour, true, null.NewInt(0, false)***REMOVED***,
				***REMOVED***24 * time.Hour, true, null.NewInt(0, false)***REMOVED***,
				***REMOVED***365 * 24 * time.Hour, true, null.NewInt(0, false)***REMOVED***,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***
	for name, data := range testdata ***REMOVED***
		t.Run(name, func(t *testing.T) ***REMOVED***
			for _, ckp := range data.Checkpoints ***REMOVED***
				t.Run(ckp.D.String(), func(t *testing.T) ***REMOVED***
					vus, keepRunning := ProcessStages(data.Start, data.Stages, ckp.D)
					assert.Equal(t, ckp.VUs, vus)
					assert.Equal(t, ckp.Keep, keepRunning)
				***REMOVED***)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***
*/
