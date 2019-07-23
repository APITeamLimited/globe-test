package brotli

/* Copyright 2016 Google Inc. All Rights Reserved.

   Distributed under MIT license.
   See file LICENSE for detail or copy at https://opensource.org/licenses/MIT
*/

/*
Dynamically grows array capacity to at least the requested size
T: data type
A: array
C: capacity
R: requested size
*/
func brotli_ensure_capacity_uint8_t(a *[]byte, c *uint, r uint) ***REMOVED***
	if *c < r ***REMOVED***
		var new_size uint = *c
		if new_size == 0 ***REMOVED***
			new_size = r
		***REMOVED***

		for new_size < r ***REMOVED***
			new_size *= 2
		***REMOVED***
		var new_array []byte = make([]byte, new_size)
		if *c != 0 ***REMOVED***
			copy(new_array, (*a)[:*c])
		***REMOVED***

		*a = new_array
		*c = new_size
	***REMOVED***
***REMOVED***

func brotli_ensure_capacity_uint32_t(a *[]uint32, c *uint, r uint) ***REMOVED***
	var new_array []uint32
	if *c < r ***REMOVED***
		var new_size uint = *c
		if new_size == 0 ***REMOVED***
			new_size = r
		***REMOVED***

		for new_size < r ***REMOVED***
			new_size *= 2
		***REMOVED***

		new_array = make([]uint32, new_size)
		if *c != 0 ***REMOVED***
			copy(new_array, (*a)[:*c])
		***REMOVED***

		*a = new_array
		*c = new_size
	***REMOVED***
***REMOVED***
