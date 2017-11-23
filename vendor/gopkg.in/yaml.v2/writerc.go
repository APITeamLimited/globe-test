package yaml

// Set the writer error and return false.
func yaml_emitter_set_writer_error(emitter *yaml_emitter_t, problem string) bool ***REMOVED***
	emitter.error = yaml_WRITER_ERROR
	emitter.problem = problem
	return false
***REMOVED***

// Flush the output buffer.
func yaml_emitter_flush(emitter *yaml_emitter_t) bool ***REMOVED***
	if emitter.write_handler == nil ***REMOVED***
		panic("write handler not set")
	***REMOVED***

	// Check if the buffer is empty.
	if emitter.buffer_pos == 0 ***REMOVED***
		return true
	***REMOVED***

	// If the output encoding is UTF-8, we don't need to recode the buffer.
	if emitter.encoding == yaml_UTF8_ENCODING ***REMOVED***
		if err := emitter.write_handler(emitter, emitter.buffer[:emitter.buffer_pos]); err != nil ***REMOVED***
			return yaml_emitter_set_writer_error(emitter, "write error: "+err.Error())
		***REMOVED***
		emitter.buffer_pos = 0
		return true
	***REMOVED***

	// Recode the buffer into the raw buffer.
	var low, high int
	if emitter.encoding == yaml_UTF16LE_ENCODING ***REMOVED***
		low, high = 0, 1
	***REMOVED*** else ***REMOVED***
		high, low = 1, 0
	***REMOVED***

	pos := 0
	for pos < emitter.buffer_pos ***REMOVED***
		// See the "reader.c" code for more details on UTF-8 encoding.  Note
		// that we assume that the buffer contains a valid UTF-8 sequence.

		// Read the next UTF-8 character.
		octet := emitter.buffer[pos]

		var w int
		var value rune
		switch ***REMOVED***
		case octet&0x80 == 0x00:
			w, value = 1, rune(octet&0x7F)
		case octet&0xE0 == 0xC0:
			w, value = 2, rune(octet&0x1F)
		case octet&0xF0 == 0xE0:
			w, value = 3, rune(octet&0x0F)
		case octet&0xF8 == 0xF0:
			w, value = 4, rune(octet&0x07)
		***REMOVED***
		for k := 1; k < w; k++ ***REMOVED***
			octet = emitter.buffer[pos+k]
			value = (value << 6) + (rune(octet) & 0x3F)
		***REMOVED***
		pos += w

		// Write the character.
		if value < 0x10000 ***REMOVED***
			var b [2]byte
			b[high] = byte(value >> 8)
			b[low] = byte(value & 0xFF)
			emitter.raw_buffer = append(emitter.raw_buffer, b[0], b[1])
		***REMOVED*** else ***REMOVED***
			// Write the character using a surrogate pair (check "reader.c").
			var b [4]byte
			value -= 0x10000
			b[high] = byte(0xD8 + (value >> 18))
			b[low] = byte((value >> 10) & 0xFF)
			b[high+2] = byte(0xDC + ((value >> 8) & 0xFF))
			b[low+2] = byte(value & 0xFF)
			emitter.raw_buffer = append(emitter.raw_buffer, b[0], b[1], b[2], b[3])
		***REMOVED***
	***REMOVED***

	// Write the raw buffer.
	if err := emitter.write_handler(emitter, emitter.raw_buffer); err != nil ***REMOVED***
		return yaml_emitter_set_writer_error(emitter, "write error: "+err.Error())
	***REMOVED***
	emitter.buffer_pos = 0
	emitter.raw_buffer = emitter.raw_buffer[:0]
	return true
***REMOVED***
