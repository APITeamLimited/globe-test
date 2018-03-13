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

	if err := emitter.write_handler(emitter, emitter.buffer[:emitter.buffer_pos]); err != nil ***REMOVED***
		return yaml_emitter_set_writer_error(emitter, "write error: "+err.Error())
	***REMOVED***
	emitter.buffer_pos = 0
	return true
***REMOVED***
