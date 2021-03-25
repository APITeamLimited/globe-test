package yaml

import (
	"bytes"
	"fmt"
)

// Flush the buffer if needed.
func flush(emitter *yaml_emitter_t) bool ***REMOVED***
	if emitter.buffer_pos+5 >= len(emitter.buffer) ***REMOVED***
		return yaml_emitter_flush(emitter)
	***REMOVED***
	return true
***REMOVED***

// Put a character to the output buffer.
func put(emitter *yaml_emitter_t, value byte) bool ***REMOVED***
	if emitter.buffer_pos+5 >= len(emitter.buffer) && !yaml_emitter_flush(emitter) ***REMOVED***
		return false
	***REMOVED***
	emitter.buffer[emitter.buffer_pos] = value
	emitter.buffer_pos++
	emitter.column++
	return true
***REMOVED***

// Put a line break to the output buffer.
func put_break(emitter *yaml_emitter_t) bool ***REMOVED***
	if emitter.buffer_pos+5 >= len(emitter.buffer) && !yaml_emitter_flush(emitter) ***REMOVED***
		return false
	***REMOVED***
	switch emitter.line_break ***REMOVED***
	case yaml_CR_BREAK:
		emitter.buffer[emitter.buffer_pos] = '\r'
		emitter.buffer_pos += 1
	case yaml_LN_BREAK:
		emitter.buffer[emitter.buffer_pos] = '\n'
		emitter.buffer_pos += 1
	case yaml_CRLN_BREAK:
		emitter.buffer[emitter.buffer_pos+0] = '\r'
		emitter.buffer[emitter.buffer_pos+1] = '\n'
		emitter.buffer_pos += 2
	default:
		panic("unknown line break setting")
	***REMOVED***
	emitter.column = 0
	emitter.line++
	return true
***REMOVED***

// Copy a character from a string into buffer.
func write(emitter *yaml_emitter_t, s []byte, i *int) bool ***REMOVED***
	if emitter.buffer_pos+5 >= len(emitter.buffer) && !yaml_emitter_flush(emitter) ***REMOVED***
		return false
	***REMOVED***
	p := emitter.buffer_pos
	w := width(s[*i])
	switch w ***REMOVED***
	case 4:
		emitter.buffer[p+3] = s[*i+3]
		fallthrough
	case 3:
		emitter.buffer[p+2] = s[*i+2]
		fallthrough
	case 2:
		emitter.buffer[p+1] = s[*i+1]
		fallthrough
	case 1:
		emitter.buffer[p+0] = s[*i+0]
	default:
		panic("unknown character width")
	***REMOVED***
	emitter.column++
	emitter.buffer_pos += w
	*i += w
	return true
***REMOVED***

// Write a whole string into buffer.
func write_all(emitter *yaml_emitter_t, s []byte) bool ***REMOVED***
	for i := 0; i < len(s); ***REMOVED***
		if !write(emitter, s, &i) ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

// Copy a line break character from a string into buffer.
func write_break(emitter *yaml_emitter_t, s []byte, i *int) bool ***REMOVED***
	if s[*i] == '\n' ***REMOVED***
		if !put_break(emitter) ***REMOVED***
			return false
		***REMOVED***
		*i++
	***REMOVED*** else ***REMOVED***
		if !write(emitter, s, i) ***REMOVED***
			return false
		***REMOVED***
		emitter.column = 0
		emitter.line++
	***REMOVED***
	return true
***REMOVED***

// Set an emitter error and return false.
func yaml_emitter_set_emitter_error(emitter *yaml_emitter_t, problem string) bool ***REMOVED***
	emitter.error = yaml_EMITTER_ERROR
	emitter.problem = problem
	return false
***REMOVED***

// Emit an event.
func yaml_emitter_emit(emitter *yaml_emitter_t, event *yaml_event_t) bool ***REMOVED***
	emitter.events = append(emitter.events, *event)
	for !yaml_emitter_need_more_events(emitter) ***REMOVED***
		event := &emitter.events[emitter.events_head]
		if !yaml_emitter_analyze_event(emitter, event) ***REMOVED***
			return false
		***REMOVED***
		if !yaml_emitter_state_machine(emitter, event) ***REMOVED***
			return false
		***REMOVED***
		yaml_event_delete(event)
		emitter.events_head++
	***REMOVED***
	return true
***REMOVED***

// Check if we need to accumulate more events before emitting.
//
// We accumulate extra
//  - 1 event for DOCUMENT-START
//  - 2 events for SEQUENCE-START
//  - 3 events for MAPPING-START
//
func yaml_emitter_need_more_events(emitter *yaml_emitter_t) bool ***REMOVED***
	if emitter.events_head == len(emitter.events) ***REMOVED***
		return true
	***REMOVED***
	var accumulate int
	switch emitter.events[emitter.events_head].typ ***REMOVED***
	case yaml_DOCUMENT_START_EVENT:
		accumulate = 1
		break
	case yaml_SEQUENCE_START_EVENT:
		accumulate = 2
		break
	case yaml_MAPPING_START_EVENT:
		accumulate = 3
		break
	default:
		return false
	***REMOVED***
	if len(emitter.events)-emitter.events_head > accumulate ***REMOVED***
		return false
	***REMOVED***
	var level int
	for i := emitter.events_head; i < len(emitter.events); i++ ***REMOVED***
		switch emitter.events[i].typ ***REMOVED***
		case yaml_STREAM_START_EVENT, yaml_DOCUMENT_START_EVENT, yaml_SEQUENCE_START_EVENT, yaml_MAPPING_START_EVENT:
			level++
		case yaml_STREAM_END_EVENT, yaml_DOCUMENT_END_EVENT, yaml_SEQUENCE_END_EVENT, yaml_MAPPING_END_EVENT:
			level--
		***REMOVED***
		if level == 0 ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

// Append a directive to the directives stack.
func yaml_emitter_append_tag_directive(emitter *yaml_emitter_t, value *yaml_tag_directive_t, allow_duplicates bool) bool ***REMOVED***
	for i := 0; i < len(emitter.tag_directives); i++ ***REMOVED***
		if bytes.Equal(value.handle, emitter.tag_directives[i].handle) ***REMOVED***
			if allow_duplicates ***REMOVED***
				return true
			***REMOVED***
			return yaml_emitter_set_emitter_error(emitter, "duplicate %TAG directive")
		***REMOVED***
	***REMOVED***

	// [Go] Do we actually need to copy this given garbage collection
	// and the lack of deallocating destructors?
	tag_copy := yaml_tag_directive_t***REMOVED***
		handle: make([]byte, len(value.handle)),
		prefix: make([]byte, len(value.prefix)),
	***REMOVED***
	copy(tag_copy.handle, value.handle)
	copy(tag_copy.prefix, value.prefix)
	emitter.tag_directives = append(emitter.tag_directives, tag_copy)
	return true
***REMOVED***

// Increase the indentation level.
func yaml_emitter_increase_indent(emitter *yaml_emitter_t, flow, indentless bool) bool ***REMOVED***
	emitter.indents = append(emitter.indents, emitter.indent)
	if emitter.indent < 0 ***REMOVED***
		if flow ***REMOVED***
			emitter.indent = emitter.best_indent
		***REMOVED*** else ***REMOVED***
			emitter.indent = 0
		***REMOVED***
	***REMOVED*** else if !indentless ***REMOVED***
		emitter.indent += emitter.best_indent
	***REMOVED***
	return true
***REMOVED***

// State dispatcher.
func yaml_emitter_state_machine(emitter *yaml_emitter_t, event *yaml_event_t) bool ***REMOVED***
	switch emitter.state ***REMOVED***
	default:
	case yaml_EMIT_STREAM_START_STATE:
		return yaml_emitter_emit_stream_start(emitter, event)

	case yaml_EMIT_FIRST_DOCUMENT_START_STATE:
		return yaml_emitter_emit_document_start(emitter, event, true)

	case yaml_EMIT_DOCUMENT_START_STATE:
		return yaml_emitter_emit_document_start(emitter, event, false)

	case yaml_EMIT_DOCUMENT_CONTENT_STATE:
		return yaml_emitter_emit_document_content(emitter, event)

	case yaml_EMIT_DOCUMENT_END_STATE:
		return yaml_emitter_emit_document_end(emitter, event)

	case yaml_EMIT_FLOW_SEQUENCE_FIRST_ITEM_STATE:
		return yaml_emitter_emit_flow_sequence_item(emitter, event, true)

	case yaml_EMIT_FLOW_SEQUENCE_ITEM_STATE:
		return yaml_emitter_emit_flow_sequence_item(emitter, event, false)

	case yaml_EMIT_FLOW_MAPPING_FIRST_KEY_STATE:
		return yaml_emitter_emit_flow_mapping_key(emitter, event, true)

	case yaml_EMIT_FLOW_MAPPING_KEY_STATE:
		return yaml_emitter_emit_flow_mapping_key(emitter, event, false)

	case yaml_EMIT_FLOW_MAPPING_SIMPLE_VALUE_STATE:
		return yaml_emitter_emit_flow_mapping_value(emitter, event, true)

	case yaml_EMIT_FLOW_MAPPING_VALUE_STATE:
		return yaml_emitter_emit_flow_mapping_value(emitter, event, false)

	case yaml_EMIT_BLOCK_SEQUENCE_FIRST_ITEM_STATE:
		return yaml_emitter_emit_block_sequence_item(emitter, event, true)

	case yaml_EMIT_BLOCK_SEQUENCE_ITEM_STATE:
		return yaml_emitter_emit_block_sequence_item(emitter, event, false)

	case yaml_EMIT_BLOCK_MAPPING_FIRST_KEY_STATE:
		return yaml_emitter_emit_block_mapping_key(emitter, event, true)

	case yaml_EMIT_BLOCK_MAPPING_KEY_STATE:
		return yaml_emitter_emit_block_mapping_key(emitter, event, false)

	case yaml_EMIT_BLOCK_MAPPING_SIMPLE_VALUE_STATE:
		return yaml_emitter_emit_block_mapping_value(emitter, event, true)

	case yaml_EMIT_BLOCK_MAPPING_VALUE_STATE:
		return yaml_emitter_emit_block_mapping_value(emitter, event, false)

	case yaml_EMIT_END_STATE:
		return yaml_emitter_set_emitter_error(emitter, "expected nothing after STREAM-END")
	***REMOVED***
	panic("invalid emitter state")
***REMOVED***

// Expect STREAM-START.
func yaml_emitter_emit_stream_start(emitter *yaml_emitter_t, event *yaml_event_t) bool ***REMOVED***
	if event.typ != yaml_STREAM_START_EVENT ***REMOVED***
		return yaml_emitter_set_emitter_error(emitter, "expected STREAM-START")
	***REMOVED***
	if emitter.encoding == yaml_ANY_ENCODING ***REMOVED***
		emitter.encoding = event.encoding
		if emitter.encoding == yaml_ANY_ENCODING ***REMOVED***
			emitter.encoding = yaml_UTF8_ENCODING
		***REMOVED***
	***REMOVED***
	if emitter.best_indent < 2 || emitter.best_indent > 9 ***REMOVED***
		emitter.best_indent = 2
	***REMOVED***
	if emitter.best_width >= 0 && emitter.best_width <= emitter.best_indent*2 ***REMOVED***
		emitter.best_width = 80
	***REMOVED***
	if emitter.best_width < 0 ***REMOVED***
		emitter.best_width = 1<<31 - 1
	***REMOVED***
	if emitter.line_break == yaml_ANY_BREAK ***REMOVED***
		emitter.line_break = yaml_LN_BREAK
	***REMOVED***

	emitter.indent = -1
	emitter.line = 0
	emitter.column = 0
	emitter.whitespace = true
	emitter.indention = true

	if emitter.encoding != yaml_UTF8_ENCODING ***REMOVED***
		if !yaml_emitter_write_bom(emitter) ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	emitter.state = yaml_EMIT_FIRST_DOCUMENT_START_STATE
	return true
***REMOVED***

// Expect DOCUMENT-START or STREAM-END.
func yaml_emitter_emit_document_start(emitter *yaml_emitter_t, event *yaml_event_t, first bool) bool ***REMOVED***

	if event.typ == yaml_DOCUMENT_START_EVENT ***REMOVED***

		if event.version_directive != nil ***REMOVED***
			if !yaml_emitter_analyze_version_directive(emitter, event.version_directive) ***REMOVED***
				return false
			***REMOVED***
		***REMOVED***

		for i := 0; i < len(event.tag_directives); i++ ***REMOVED***
			tag_directive := &event.tag_directives[i]
			if !yaml_emitter_analyze_tag_directive(emitter, tag_directive) ***REMOVED***
				return false
			***REMOVED***
			if !yaml_emitter_append_tag_directive(emitter, tag_directive, false) ***REMOVED***
				return false
			***REMOVED***
		***REMOVED***

		for i := 0; i < len(default_tag_directives); i++ ***REMOVED***
			tag_directive := &default_tag_directives[i]
			if !yaml_emitter_append_tag_directive(emitter, tag_directive, true) ***REMOVED***
				return false
			***REMOVED***
		***REMOVED***

		implicit := event.implicit
		if !first || emitter.canonical ***REMOVED***
			implicit = false
		***REMOVED***

		if emitter.open_ended && (event.version_directive != nil || len(event.tag_directives) > 0) ***REMOVED***
			if !yaml_emitter_write_indicator(emitter, []byte("..."), true, false, false) ***REMOVED***
				return false
			***REMOVED***
			if !yaml_emitter_write_indent(emitter) ***REMOVED***
				return false
			***REMOVED***
		***REMOVED***

		if event.version_directive != nil ***REMOVED***
			implicit = false
			if !yaml_emitter_write_indicator(emitter, []byte("%YAML"), true, false, false) ***REMOVED***
				return false
			***REMOVED***
			if !yaml_emitter_write_indicator(emitter, []byte("1.1"), true, false, false) ***REMOVED***
				return false
			***REMOVED***
			if !yaml_emitter_write_indent(emitter) ***REMOVED***
				return false
			***REMOVED***
		***REMOVED***

		if len(event.tag_directives) > 0 ***REMOVED***
			implicit = false
			for i := 0; i < len(event.tag_directives); i++ ***REMOVED***
				tag_directive := &event.tag_directives[i]
				if !yaml_emitter_write_indicator(emitter, []byte("%TAG"), true, false, false) ***REMOVED***
					return false
				***REMOVED***
				if !yaml_emitter_write_tag_handle(emitter, tag_directive.handle) ***REMOVED***
					return false
				***REMOVED***
				if !yaml_emitter_write_tag_content(emitter, tag_directive.prefix, true) ***REMOVED***
					return false
				***REMOVED***
				if !yaml_emitter_write_indent(emitter) ***REMOVED***
					return false
				***REMOVED***
			***REMOVED***
		***REMOVED***

		if yaml_emitter_check_empty_document(emitter) ***REMOVED***
			implicit = false
		***REMOVED***
		if !implicit ***REMOVED***
			if !yaml_emitter_write_indent(emitter) ***REMOVED***
				return false
			***REMOVED***
			if !yaml_emitter_write_indicator(emitter, []byte("---"), true, false, false) ***REMOVED***
				return false
			***REMOVED***
			if emitter.canonical ***REMOVED***
				if !yaml_emitter_write_indent(emitter) ***REMOVED***
					return false
				***REMOVED***
			***REMOVED***
		***REMOVED***

		emitter.state = yaml_EMIT_DOCUMENT_CONTENT_STATE
		return true
	***REMOVED***

	if event.typ == yaml_STREAM_END_EVENT ***REMOVED***
		if emitter.open_ended ***REMOVED***
			if !yaml_emitter_write_indicator(emitter, []byte("..."), true, false, false) ***REMOVED***
				return false
			***REMOVED***
			if !yaml_emitter_write_indent(emitter) ***REMOVED***
				return false
			***REMOVED***
		***REMOVED***
		if !yaml_emitter_flush(emitter) ***REMOVED***
			return false
		***REMOVED***
		emitter.state = yaml_EMIT_END_STATE
		return true
	***REMOVED***

	return yaml_emitter_set_emitter_error(emitter, "expected DOCUMENT-START or STREAM-END")
***REMOVED***

// Expect the root node.
func yaml_emitter_emit_document_content(emitter *yaml_emitter_t, event *yaml_event_t) bool ***REMOVED***
	emitter.states = append(emitter.states, yaml_EMIT_DOCUMENT_END_STATE)
	return yaml_emitter_emit_node(emitter, event, true, false, false, false)
***REMOVED***

// Expect DOCUMENT-END.
func yaml_emitter_emit_document_end(emitter *yaml_emitter_t, event *yaml_event_t) bool ***REMOVED***
	if event.typ != yaml_DOCUMENT_END_EVENT ***REMOVED***
		return yaml_emitter_set_emitter_error(emitter, "expected DOCUMENT-END")
	***REMOVED***
	if !yaml_emitter_write_indent(emitter) ***REMOVED***
		return false
	***REMOVED***
	if !event.implicit ***REMOVED***
		// [Go] Allocate the slice elsewhere.
		if !yaml_emitter_write_indicator(emitter, []byte("..."), true, false, false) ***REMOVED***
			return false
		***REMOVED***
		if !yaml_emitter_write_indent(emitter) ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	if !yaml_emitter_flush(emitter) ***REMOVED***
		return false
	***REMOVED***
	emitter.state = yaml_EMIT_DOCUMENT_START_STATE
	emitter.tag_directives = emitter.tag_directives[:0]
	return true
***REMOVED***

// Expect a flow item node.
func yaml_emitter_emit_flow_sequence_item(emitter *yaml_emitter_t, event *yaml_event_t, first bool) bool ***REMOVED***
	if first ***REMOVED***
		if !yaml_emitter_write_indicator(emitter, []byte***REMOVED***'['***REMOVED***, true, true, false) ***REMOVED***
			return false
		***REMOVED***
		if !yaml_emitter_increase_indent(emitter, true, false) ***REMOVED***
			return false
		***REMOVED***
		emitter.flow_level++
	***REMOVED***

	if event.typ == yaml_SEQUENCE_END_EVENT ***REMOVED***
		emitter.flow_level--
		emitter.indent = emitter.indents[len(emitter.indents)-1]
		emitter.indents = emitter.indents[:len(emitter.indents)-1]
		if emitter.canonical && !first ***REMOVED***
			if !yaml_emitter_write_indicator(emitter, []byte***REMOVED***','***REMOVED***, false, false, false) ***REMOVED***
				return false
			***REMOVED***
			if !yaml_emitter_write_indent(emitter) ***REMOVED***
				return false
			***REMOVED***
		***REMOVED***
		if !yaml_emitter_write_indicator(emitter, []byte***REMOVED***']'***REMOVED***, false, false, false) ***REMOVED***
			return false
		***REMOVED***
		emitter.state = emitter.states[len(emitter.states)-1]
		emitter.states = emitter.states[:len(emitter.states)-1]

		return true
	***REMOVED***

	if !first ***REMOVED***
		if !yaml_emitter_write_indicator(emitter, []byte***REMOVED***','***REMOVED***, false, false, false) ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***

	if emitter.canonical || emitter.column > emitter.best_width ***REMOVED***
		if !yaml_emitter_write_indent(emitter) ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	emitter.states = append(emitter.states, yaml_EMIT_FLOW_SEQUENCE_ITEM_STATE)
	return yaml_emitter_emit_node(emitter, event, false, true, false, false)
***REMOVED***

// Expect a flow key node.
func yaml_emitter_emit_flow_mapping_key(emitter *yaml_emitter_t, event *yaml_event_t, first bool) bool ***REMOVED***
	if first ***REMOVED***
		if !yaml_emitter_write_indicator(emitter, []byte***REMOVED***'***REMOVED***'***REMOVED***, true, true, false) ***REMOVED***
			return false
		***REMOVED***
		if !yaml_emitter_increase_indent(emitter, true, false) ***REMOVED***
			return false
		***REMOVED***
		emitter.flow_level++
	***REMOVED***

	if event.typ == yaml_MAPPING_END_EVENT ***REMOVED***
		emitter.flow_level--
		emitter.indent = emitter.indents[len(emitter.indents)-1]
		emitter.indents = emitter.indents[:len(emitter.indents)-1]
		if emitter.canonical && !first ***REMOVED***
			if !yaml_emitter_write_indicator(emitter, []byte***REMOVED***','***REMOVED***, false, false, false) ***REMOVED***
				return false
			***REMOVED***
			if !yaml_emitter_write_indent(emitter) ***REMOVED***
				return false
			***REMOVED***
		***REMOVED***
		if !yaml_emitter_write_indicator(emitter, []byte***REMOVED***'***REMOVED***'***REMOVED***, false, false, false) ***REMOVED***
			return false
		***REMOVED***
		emitter.state = emitter.states[len(emitter.states)-1]
		emitter.states = emitter.states[:len(emitter.states)-1]
		return true
	***REMOVED***

	if !first ***REMOVED***
		if !yaml_emitter_write_indicator(emitter, []byte***REMOVED***','***REMOVED***, false, false, false) ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	if emitter.canonical || emitter.column > emitter.best_width ***REMOVED***
		if !yaml_emitter_write_indent(emitter) ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***

	if !emitter.canonical && yaml_emitter_check_simple_key(emitter) ***REMOVED***
		emitter.states = append(emitter.states, yaml_EMIT_FLOW_MAPPING_SIMPLE_VALUE_STATE)
		return yaml_emitter_emit_node(emitter, event, false, false, true, true)
	***REMOVED***
	if !yaml_emitter_write_indicator(emitter, []byte***REMOVED***'?'***REMOVED***, true, false, false) ***REMOVED***
		return false
	***REMOVED***
	emitter.states = append(emitter.states, yaml_EMIT_FLOW_MAPPING_VALUE_STATE)
	return yaml_emitter_emit_node(emitter, event, false, false, true, false)
***REMOVED***

// Expect a flow value node.
func yaml_emitter_emit_flow_mapping_value(emitter *yaml_emitter_t, event *yaml_event_t, simple bool) bool ***REMOVED***
	if simple ***REMOVED***
		if !yaml_emitter_write_indicator(emitter, []byte***REMOVED***':'***REMOVED***, false, false, false) ***REMOVED***
			return false
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if emitter.canonical || emitter.column > emitter.best_width ***REMOVED***
			if !yaml_emitter_write_indent(emitter) ***REMOVED***
				return false
			***REMOVED***
		***REMOVED***
		if !yaml_emitter_write_indicator(emitter, []byte***REMOVED***':'***REMOVED***, true, false, false) ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	emitter.states = append(emitter.states, yaml_EMIT_FLOW_MAPPING_KEY_STATE)
	return yaml_emitter_emit_node(emitter, event, false, false, true, false)
***REMOVED***

// Expect a block item node.
func yaml_emitter_emit_block_sequence_item(emitter *yaml_emitter_t, event *yaml_event_t, first bool) bool ***REMOVED***
	if first ***REMOVED***
		if !yaml_emitter_increase_indent(emitter, false, emitter.mapping_context && !emitter.indention) ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	if event.typ == yaml_SEQUENCE_END_EVENT ***REMOVED***
		emitter.indent = emitter.indents[len(emitter.indents)-1]
		emitter.indents = emitter.indents[:len(emitter.indents)-1]
		emitter.state = emitter.states[len(emitter.states)-1]
		emitter.states = emitter.states[:len(emitter.states)-1]
		return true
	***REMOVED***
	if !yaml_emitter_write_indent(emitter) ***REMOVED***
		return false
	***REMOVED***
	if !yaml_emitter_write_indicator(emitter, []byte***REMOVED***'-'***REMOVED***, true, false, true) ***REMOVED***
		return false
	***REMOVED***
	emitter.states = append(emitter.states, yaml_EMIT_BLOCK_SEQUENCE_ITEM_STATE)
	return yaml_emitter_emit_node(emitter, event, false, true, false, false)
***REMOVED***

// Expect a block key node.
func yaml_emitter_emit_block_mapping_key(emitter *yaml_emitter_t, event *yaml_event_t, first bool) bool ***REMOVED***
	if first ***REMOVED***
		if !yaml_emitter_increase_indent(emitter, false, false) ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	if event.typ == yaml_MAPPING_END_EVENT ***REMOVED***
		emitter.indent = emitter.indents[len(emitter.indents)-1]
		emitter.indents = emitter.indents[:len(emitter.indents)-1]
		emitter.state = emitter.states[len(emitter.states)-1]
		emitter.states = emitter.states[:len(emitter.states)-1]
		return true
	***REMOVED***
	if !yaml_emitter_write_indent(emitter) ***REMOVED***
		return false
	***REMOVED***
	if yaml_emitter_check_simple_key(emitter) ***REMOVED***
		emitter.states = append(emitter.states, yaml_EMIT_BLOCK_MAPPING_SIMPLE_VALUE_STATE)
		return yaml_emitter_emit_node(emitter, event, false, false, true, true)
	***REMOVED***
	if !yaml_emitter_write_indicator(emitter, []byte***REMOVED***'?'***REMOVED***, true, false, true) ***REMOVED***
		return false
	***REMOVED***
	emitter.states = append(emitter.states, yaml_EMIT_BLOCK_MAPPING_VALUE_STATE)
	return yaml_emitter_emit_node(emitter, event, false, false, true, false)
***REMOVED***

// Expect a block value node.
func yaml_emitter_emit_block_mapping_value(emitter *yaml_emitter_t, event *yaml_event_t, simple bool) bool ***REMOVED***
	if simple ***REMOVED***
		if !yaml_emitter_write_indicator(emitter, []byte***REMOVED***':'***REMOVED***, false, false, false) ***REMOVED***
			return false
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if !yaml_emitter_write_indent(emitter) ***REMOVED***
			return false
		***REMOVED***
		if !yaml_emitter_write_indicator(emitter, []byte***REMOVED***':'***REMOVED***, true, false, true) ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	emitter.states = append(emitter.states, yaml_EMIT_BLOCK_MAPPING_KEY_STATE)
	return yaml_emitter_emit_node(emitter, event, false, false, true, false)
***REMOVED***

// Expect a node.
func yaml_emitter_emit_node(emitter *yaml_emitter_t, event *yaml_event_t,
	root bool, sequence bool, mapping bool, simple_key bool) bool ***REMOVED***

	emitter.root_context = root
	emitter.sequence_context = sequence
	emitter.mapping_context = mapping
	emitter.simple_key_context = simple_key

	switch event.typ ***REMOVED***
	case yaml_ALIAS_EVENT:
		return yaml_emitter_emit_alias(emitter, event)
	case yaml_SCALAR_EVENT:
		return yaml_emitter_emit_scalar(emitter, event)
	case yaml_SEQUENCE_START_EVENT:
		return yaml_emitter_emit_sequence_start(emitter, event)
	case yaml_MAPPING_START_EVENT:
		return yaml_emitter_emit_mapping_start(emitter, event)
	default:
		return yaml_emitter_set_emitter_error(emitter,
			fmt.Sprintf("expected SCALAR, SEQUENCE-START, MAPPING-START, or ALIAS, but got %v", event.typ))
	***REMOVED***
***REMOVED***

// Expect ALIAS.
func yaml_emitter_emit_alias(emitter *yaml_emitter_t, event *yaml_event_t) bool ***REMOVED***
	if !yaml_emitter_process_anchor(emitter) ***REMOVED***
		return false
	***REMOVED***
	emitter.state = emitter.states[len(emitter.states)-1]
	emitter.states = emitter.states[:len(emitter.states)-1]
	return true
***REMOVED***

// Expect SCALAR.
func yaml_emitter_emit_scalar(emitter *yaml_emitter_t, event *yaml_event_t) bool ***REMOVED***
	if !yaml_emitter_select_scalar_style(emitter, event) ***REMOVED***
		return false
	***REMOVED***
	if !yaml_emitter_process_anchor(emitter) ***REMOVED***
		return false
	***REMOVED***
	if !yaml_emitter_process_tag(emitter) ***REMOVED***
		return false
	***REMOVED***
	if !yaml_emitter_increase_indent(emitter, true, false) ***REMOVED***
		return false
	***REMOVED***
	if !yaml_emitter_process_scalar(emitter) ***REMOVED***
		return false
	***REMOVED***
	emitter.indent = emitter.indents[len(emitter.indents)-1]
	emitter.indents = emitter.indents[:len(emitter.indents)-1]
	emitter.state = emitter.states[len(emitter.states)-1]
	emitter.states = emitter.states[:len(emitter.states)-1]
	return true
***REMOVED***

// Expect SEQUENCE-START.
func yaml_emitter_emit_sequence_start(emitter *yaml_emitter_t, event *yaml_event_t) bool ***REMOVED***
	if !yaml_emitter_process_anchor(emitter) ***REMOVED***
		return false
	***REMOVED***
	if !yaml_emitter_process_tag(emitter) ***REMOVED***
		return false
	***REMOVED***
	if emitter.flow_level > 0 || emitter.canonical || event.sequence_style() == yaml_FLOW_SEQUENCE_STYLE ||
		yaml_emitter_check_empty_sequence(emitter) ***REMOVED***
		emitter.state = yaml_EMIT_FLOW_SEQUENCE_FIRST_ITEM_STATE
	***REMOVED*** else ***REMOVED***
		emitter.state = yaml_EMIT_BLOCK_SEQUENCE_FIRST_ITEM_STATE
	***REMOVED***
	return true
***REMOVED***

// Expect MAPPING-START.
func yaml_emitter_emit_mapping_start(emitter *yaml_emitter_t, event *yaml_event_t) bool ***REMOVED***
	if !yaml_emitter_process_anchor(emitter) ***REMOVED***
		return false
	***REMOVED***
	if !yaml_emitter_process_tag(emitter) ***REMOVED***
		return false
	***REMOVED***
	if emitter.flow_level > 0 || emitter.canonical || event.mapping_style() == yaml_FLOW_MAPPING_STYLE ||
		yaml_emitter_check_empty_mapping(emitter) ***REMOVED***
		emitter.state = yaml_EMIT_FLOW_MAPPING_FIRST_KEY_STATE
	***REMOVED*** else ***REMOVED***
		emitter.state = yaml_EMIT_BLOCK_MAPPING_FIRST_KEY_STATE
	***REMOVED***
	return true
***REMOVED***

// Check if the document content is an empty scalar.
func yaml_emitter_check_empty_document(emitter *yaml_emitter_t) bool ***REMOVED***
	return false // [Go] Huh?
***REMOVED***

// Check if the next events represent an empty sequence.
func yaml_emitter_check_empty_sequence(emitter *yaml_emitter_t) bool ***REMOVED***
	if len(emitter.events)-emitter.events_head < 2 ***REMOVED***
		return false
	***REMOVED***
	return emitter.events[emitter.events_head].typ == yaml_SEQUENCE_START_EVENT &&
		emitter.events[emitter.events_head+1].typ == yaml_SEQUENCE_END_EVENT
***REMOVED***

// Check if the next events represent an empty mapping.
func yaml_emitter_check_empty_mapping(emitter *yaml_emitter_t) bool ***REMOVED***
	if len(emitter.events)-emitter.events_head < 2 ***REMOVED***
		return false
	***REMOVED***
	return emitter.events[emitter.events_head].typ == yaml_MAPPING_START_EVENT &&
		emitter.events[emitter.events_head+1].typ == yaml_MAPPING_END_EVENT
***REMOVED***

// Check if the next node can be expressed as a simple key.
func yaml_emitter_check_simple_key(emitter *yaml_emitter_t) bool ***REMOVED***
	length := 0
	switch emitter.events[emitter.events_head].typ ***REMOVED***
	case yaml_ALIAS_EVENT:
		length += len(emitter.anchor_data.anchor)
	case yaml_SCALAR_EVENT:
		if emitter.scalar_data.multiline ***REMOVED***
			return false
		***REMOVED***
		length += len(emitter.anchor_data.anchor) +
			len(emitter.tag_data.handle) +
			len(emitter.tag_data.suffix) +
			len(emitter.scalar_data.value)
	case yaml_SEQUENCE_START_EVENT:
		if !yaml_emitter_check_empty_sequence(emitter) ***REMOVED***
			return false
		***REMOVED***
		length += len(emitter.anchor_data.anchor) +
			len(emitter.tag_data.handle) +
			len(emitter.tag_data.suffix)
	case yaml_MAPPING_START_EVENT:
		if !yaml_emitter_check_empty_mapping(emitter) ***REMOVED***
			return false
		***REMOVED***
		length += len(emitter.anchor_data.anchor) +
			len(emitter.tag_data.handle) +
			len(emitter.tag_data.suffix)
	default:
		return false
	***REMOVED***
	return length <= 128
***REMOVED***

// Determine an acceptable scalar style.
func yaml_emitter_select_scalar_style(emitter *yaml_emitter_t, event *yaml_event_t) bool ***REMOVED***

	no_tag := len(emitter.tag_data.handle) == 0 && len(emitter.tag_data.suffix) == 0
	if no_tag && !event.implicit && !event.quoted_implicit ***REMOVED***
		return yaml_emitter_set_emitter_error(emitter, "neither tag nor implicit flags are specified")
	***REMOVED***

	style := event.scalar_style()
	if style == yaml_ANY_SCALAR_STYLE ***REMOVED***
		style = yaml_PLAIN_SCALAR_STYLE
	***REMOVED***
	if emitter.canonical ***REMOVED***
		style = yaml_DOUBLE_QUOTED_SCALAR_STYLE
	***REMOVED***
	if emitter.simple_key_context && emitter.scalar_data.multiline ***REMOVED***
		style = yaml_DOUBLE_QUOTED_SCALAR_STYLE
	***REMOVED***

	if style == yaml_PLAIN_SCALAR_STYLE ***REMOVED***
		if emitter.flow_level > 0 && !emitter.scalar_data.flow_plain_allowed ||
			emitter.flow_level == 0 && !emitter.scalar_data.block_plain_allowed ***REMOVED***
			style = yaml_SINGLE_QUOTED_SCALAR_STYLE
		***REMOVED***
		if len(emitter.scalar_data.value) == 0 && (emitter.flow_level > 0 || emitter.simple_key_context) ***REMOVED***
			style = yaml_SINGLE_QUOTED_SCALAR_STYLE
		***REMOVED***
		if no_tag && !event.implicit ***REMOVED***
			style = yaml_SINGLE_QUOTED_SCALAR_STYLE
		***REMOVED***
	***REMOVED***
	if style == yaml_SINGLE_QUOTED_SCALAR_STYLE ***REMOVED***
		if !emitter.scalar_data.single_quoted_allowed ***REMOVED***
			style = yaml_DOUBLE_QUOTED_SCALAR_STYLE
		***REMOVED***
	***REMOVED***
	if style == yaml_LITERAL_SCALAR_STYLE || style == yaml_FOLDED_SCALAR_STYLE ***REMOVED***
		if !emitter.scalar_data.block_allowed || emitter.flow_level > 0 || emitter.simple_key_context ***REMOVED***
			style = yaml_DOUBLE_QUOTED_SCALAR_STYLE
		***REMOVED***
	***REMOVED***

	if no_tag && !event.quoted_implicit && style != yaml_PLAIN_SCALAR_STYLE ***REMOVED***
		emitter.tag_data.handle = []byte***REMOVED***'!'***REMOVED***
	***REMOVED***
	emitter.scalar_data.style = style
	return true
***REMOVED***

// Write an anchor.
func yaml_emitter_process_anchor(emitter *yaml_emitter_t) bool ***REMOVED***
	if emitter.anchor_data.anchor == nil ***REMOVED***
		return true
	***REMOVED***
	c := []byte***REMOVED***'&'***REMOVED***
	if emitter.anchor_data.alias ***REMOVED***
		c[0] = '*'
	***REMOVED***
	if !yaml_emitter_write_indicator(emitter, c, true, false, false) ***REMOVED***
		return false
	***REMOVED***
	return yaml_emitter_write_anchor(emitter, emitter.anchor_data.anchor)
***REMOVED***

// Write a tag.
func yaml_emitter_process_tag(emitter *yaml_emitter_t) bool ***REMOVED***
	if len(emitter.tag_data.handle) == 0 && len(emitter.tag_data.suffix) == 0 ***REMOVED***
		return true
	***REMOVED***
	if len(emitter.tag_data.handle) > 0 ***REMOVED***
		if !yaml_emitter_write_tag_handle(emitter, emitter.tag_data.handle) ***REMOVED***
			return false
		***REMOVED***
		if len(emitter.tag_data.suffix) > 0 ***REMOVED***
			if !yaml_emitter_write_tag_content(emitter, emitter.tag_data.suffix, false) ***REMOVED***
				return false
			***REMOVED***
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		// [Go] Allocate these slices elsewhere.
		if !yaml_emitter_write_indicator(emitter, []byte("!<"), true, false, false) ***REMOVED***
			return false
		***REMOVED***
		if !yaml_emitter_write_tag_content(emitter, emitter.tag_data.suffix, false) ***REMOVED***
			return false
		***REMOVED***
		if !yaml_emitter_write_indicator(emitter, []byte***REMOVED***'>'***REMOVED***, false, false, false) ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

// Write a scalar.
func yaml_emitter_process_scalar(emitter *yaml_emitter_t) bool ***REMOVED***
	switch emitter.scalar_data.style ***REMOVED***
	case yaml_PLAIN_SCALAR_STYLE:
		return yaml_emitter_write_plain_scalar(emitter, emitter.scalar_data.value, !emitter.simple_key_context)

	case yaml_SINGLE_QUOTED_SCALAR_STYLE:
		return yaml_emitter_write_single_quoted_scalar(emitter, emitter.scalar_data.value, !emitter.simple_key_context)

	case yaml_DOUBLE_QUOTED_SCALAR_STYLE:
		return yaml_emitter_write_double_quoted_scalar(emitter, emitter.scalar_data.value, !emitter.simple_key_context)

	case yaml_LITERAL_SCALAR_STYLE:
		return yaml_emitter_write_literal_scalar(emitter, emitter.scalar_data.value)

	case yaml_FOLDED_SCALAR_STYLE:
		return yaml_emitter_write_folded_scalar(emitter, emitter.scalar_data.value)
	***REMOVED***
	panic("unknown scalar style")
***REMOVED***

// Check if a %YAML directive is valid.
func yaml_emitter_analyze_version_directive(emitter *yaml_emitter_t, version_directive *yaml_version_directive_t) bool ***REMOVED***
	if version_directive.major != 1 || version_directive.minor != 1 ***REMOVED***
		return yaml_emitter_set_emitter_error(emitter, "incompatible %YAML directive")
	***REMOVED***
	return true
***REMOVED***

// Check if a %TAG directive is valid.
func yaml_emitter_analyze_tag_directive(emitter *yaml_emitter_t, tag_directive *yaml_tag_directive_t) bool ***REMOVED***
	handle := tag_directive.handle
	prefix := tag_directive.prefix
	if len(handle) == 0 ***REMOVED***
		return yaml_emitter_set_emitter_error(emitter, "tag handle must not be empty")
	***REMOVED***
	if handle[0] != '!' ***REMOVED***
		return yaml_emitter_set_emitter_error(emitter, "tag handle must start with '!'")
	***REMOVED***
	if handle[len(handle)-1] != '!' ***REMOVED***
		return yaml_emitter_set_emitter_error(emitter, "tag handle must end with '!'")
	***REMOVED***
	for i := 1; i < len(handle)-1; i += width(handle[i]) ***REMOVED***
		if !is_alpha(handle, i) ***REMOVED***
			return yaml_emitter_set_emitter_error(emitter, "tag handle must contain alphanumerical characters only")
		***REMOVED***
	***REMOVED***
	if len(prefix) == 0 ***REMOVED***
		return yaml_emitter_set_emitter_error(emitter, "tag prefix must not be empty")
	***REMOVED***
	return true
***REMOVED***

// Check if an anchor is valid.
func yaml_emitter_analyze_anchor(emitter *yaml_emitter_t, anchor []byte, alias bool) bool ***REMOVED***
	if len(anchor) == 0 ***REMOVED***
		problem := "anchor value must not be empty"
		if alias ***REMOVED***
			problem = "alias value must not be empty"
		***REMOVED***
		return yaml_emitter_set_emitter_error(emitter, problem)
	***REMOVED***
	for i := 0; i < len(anchor); i += width(anchor[i]) ***REMOVED***
		if !is_alpha(anchor, i) ***REMOVED***
			problem := "anchor value must contain alphanumerical characters only"
			if alias ***REMOVED***
				problem = "alias value must contain alphanumerical characters only"
			***REMOVED***
			return yaml_emitter_set_emitter_error(emitter, problem)
		***REMOVED***
	***REMOVED***
	emitter.anchor_data.anchor = anchor
	emitter.anchor_data.alias = alias
	return true
***REMOVED***

// Check if a tag is valid.
func yaml_emitter_analyze_tag(emitter *yaml_emitter_t, tag []byte) bool ***REMOVED***
	if len(tag) == 0 ***REMOVED***
		return yaml_emitter_set_emitter_error(emitter, "tag value must not be empty")
	***REMOVED***
	for i := 0; i < len(emitter.tag_directives); i++ ***REMOVED***
		tag_directive := &emitter.tag_directives[i]
		if bytes.HasPrefix(tag, tag_directive.prefix) ***REMOVED***
			emitter.tag_data.handle = tag_directive.handle
			emitter.tag_data.suffix = tag[len(tag_directive.prefix):]
			return true
		***REMOVED***
	***REMOVED***
	emitter.tag_data.suffix = tag
	return true
***REMOVED***

// Check if a scalar is valid.
func yaml_emitter_analyze_scalar(emitter *yaml_emitter_t, value []byte) bool ***REMOVED***
	var (
		block_indicators   = false
		flow_indicators    = false
		line_breaks        = false
		special_characters = false

		leading_space  = false
		leading_break  = false
		trailing_space = false
		trailing_break = false
		break_space    = false
		space_break    = false

		preceded_by_whitespace = false
		followed_by_whitespace = false
		previous_space         = false
		previous_break         = false
	)

	emitter.scalar_data.value = value

	if len(value) == 0 ***REMOVED***
		emitter.scalar_data.multiline = false
		emitter.scalar_data.flow_plain_allowed = false
		emitter.scalar_data.block_plain_allowed = true
		emitter.scalar_data.single_quoted_allowed = true
		emitter.scalar_data.block_allowed = false
		return true
	***REMOVED***

	if len(value) >= 3 && ((value[0] == '-' && value[1] == '-' && value[2] == '-') || (value[0] == '.' && value[1] == '.' && value[2] == '.')) ***REMOVED***
		block_indicators = true
		flow_indicators = true
	***REMOVED***

	preceded_by_whitespace = true
	for i, w := 0, 0; i < len(value); i += w ***REMOVED***
		w = width(value[i])
		followed_by_whitespace = i+w >= len(value) || is_blank(value, i+w)

		if i == 0 ***REMOVED***
			switch value[i] ***REMOVED***
			case '#', ',', '[', ']', '***REMOVED***', '***REMOVED***', '&', '*', '!', '|', '>', '\'', '"', '%', '@', '`':
				flow_indicators = true
				block_indicators = true
			case '?', ':':
				flow_indicators = true
				if followed_by_whitespace ***REMOVED***
					block_indicators = true
				***REMOVED***
			case '-':
				if followed_by_whitespace ***REMOVED***
					flow_indicators = true
					block_indicators = true
				***REMOVED***
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			switch value[i] ***REMOVED***
			case ',', '?', '[', ']', '***REMOVED***', '***REMOVED***':
				flow_indicators = true
			case ':':
				flow_indicators = true
				if followed_by_whitespace ***REMOVED***
					block_indicators = true
				***REMOVED***
			case '#':
				if preceded_by_whitespace ***REMOVED***
					flow_indicators = true
					block_indicators = true
				***REMOVED***
			***REMOVED***
		***REMOVED***

		if !is_printable(value, i) || !is_ascii(value, i) && !emitter.unicode ***REMOVED***
			special_characters = true
		***REMOVED***
		if is_space(value, i) ***REMOVED***
			if i == 0 ***REMOVED***
				leading_space = true
			***REMOVED***
			if i+width(value[i]) == len(value) ***REMOVED***
				trailing_space = true
			***REMOVED***
			if previous_break ***REMOVED***
				break_space = true
			***REMOVED***
			previous_space = true
			previous_break = false
		***REMOVED*** else if is_break(value, i) ***REMOVED***
			line_breaks = true
			if i == 0 ***REMOVED***
				leading_break = true
			***REMOVED***
			if i+width(value[i]) == len(value) ***REMOVED***
				trailing_break = true
			***REMOVED***
			if previous_space ***REMOVED***
				space_break = true
			***REMOVED***
			previous_space = false
			previous_break = true
		***REMOVED*** else ***REMOVED***
			previous_space = false
			previous_break = false
		***REMOVED***

		// [Go]: Why 'z'? Couldn't be the end of the string as that's the loop condition.
		preceded_by_whitespace = is_blankz(value, i)
	***REMOVED***

	emitter.scalar_data.multiline = line_breaks
	emitter.scalar_data.flow_plain_allowed = true
	emitter.scalar_data.block_plain_allowed = true
	emitter.scalar_data.single_quoted_allowed = true
	emitter.scalar_data.block_allowed = true

	if leading_space || leading_break || trailing_space || trailing_break ***REMOVED***
		emitter.scalar_data.flow_plain_allowed = false
		emitter.scalar_data.block_plain_allowed = false
	***REMOVED***
	if trailing_space ***REMOVED***
		emitter.scalar_data.block_allowed = false
	***REMOVED***
	if break_space ***REMOVED***
		emitter.scalar_data.flow_plain_allowed = false
		emitter.scalar_data.block_plain_allowed = false
		emitter.scalar_data.single_quoted_allowed = false
	***REMOVED***
	if space_break || special_characters ***REMOVED***
		emitter.scalar_data.flow_plain_allowed = false
		emitter.scalar_data.block_plain_allowed = false
		emitter.scalar_data.single_quoted_allowed = false
		emitter.scalar_data.block_allowed = false
	***REMOVED***
	if line_breaks ***REMOVED***
		emitter.scalar_data.flow_plain_allowed = false
		emitter.scalar_data.block_plain_allowed = false
	***REMOVED***
	if flow_indicators ***REMOVED***
		emitter.scalar_data.flow_plain_allowed = false
	***REMOVED***
	if block_indicators ***REMOVED***
		emitter.scalar_data.block_plain_allowed = false
	***REMOVED***
	return true
***REMOVED***

// Check if the event data is valid.
func yaml_emitter_analyze_event(emitter *yaml_emitter_t, event *yaml_event_t) bool ***REMOVED***

	emitter.anchor_data.anchor = nil
	emitter.tag_data.handle = nil
	emitter.tag_data.suffix = nil
	emitter.scalar_data.value = nil

	switch event.typ ***REMOVED***
	case yaml_ALIAS_EVENT:
		if !yaml_emitter_analyze_anchor(emitter, event.anchor, true) ***REMOVED***
			return false
		***REMOVED***

	case yaml_SCALAR_EVENT:
		if len(event.anchor) > 0 ***REMOVED***
			if !yaml_emitter_analyze_anchor(emitter, event.anchor, false) ***REMOVED***
				return false
			***REMOVED***
		***REMOVED***
		if len(event.tag) > 0 && (emitter.canonical || (!event.implicit && !event.quoted_implicit)) ***REMOVED***
			if !yaml_emitter_analyze_tag(emitter, event.tag) ***REMOVED***
				return false
			***REMOVED***
		***REMOVED***
		if !yaml_emitter_analyze_scalar(emitter, event.value) ***REMOVED***
			return false
		***REMOVED***

	case yaml_SEQUENCE_START_EVENT:
		if len(event.anchor) > 0 ***REMOVED***
			if !yaml_emitter_analyze_anchor(emitter, event.anchor, false) ***REMOVED***
				return false
			***REMOVED***
		***REMOVED***
		if len(event.tag) > 0 && (emitter.canonical || !event.implicit) ***REMOVED***
			if !yaml_emitter_analyze_tag(emitter, event.tag) ***REMOVED***
				return false
			***REMOVED***
		***REMOVED***

	case yaml_MAPPING_START_EVENT:
		if len(event.anchor) > 0 ***REMOVED***
			if !yaml_emitter_analyze_anchor(emitter, event.anchor, false) ***REMOVED***
				return false
			***REMOVED***
		***REMOVED***
		if len(event.tag) > 0 && (emitter.canonical || !event.implicit) ***REMOVED***
			if !yaml_emitter_analyze_tag(emitter, event.tag) ***REMOVED***
				return false
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

// Write the BOM character.
func yaml_emitter_write_bom(emitter *yaml_emitter_t) bool ***REMOVED***
	if !flush(emitter) ***REMOVED***
		return false
	***REMOVED***
	pos := emitter.buffer_pos
	emitter.buffer[pos+0] = '\xEF'
	emitter.buffer[pos+1] = '\xBB'
	emitter.buffer[pos+2] = '\xBF'
	emitter.buffer_pos += 3
	return true
***REMOVED***

func yaml_emitter_write_indent(emitter *yaml_emitter_t) bool ***REMOVED***
	indent := emitter.indent
	if indent < 0 ***REMOVED***
		indent = 0
	***REMOVED***
	if !emitter.indention || emitter.column > indent || (emitter.column == indent && !emitter.whitespace) ***REMOVED***
		if !put_break(emitter) ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	for emitter.column < indent ***REMOVED***
		if !put(emitter, ' ') ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	emitter.whitespace = true
	emitter.indention = true
	return true
***REMOVED***

func yaml_emitter_write_indicator(emitter *yaml_emitter_t, indicator []byte, need_whitespace, is_whitespace, is_indention bool) bool ***REMOVED***
	if need_whitespace && !emitter.whitespace ***REMOVED***
		if !put(emitter, ' ') ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	if !write_all(emitter, indicator) ***REMOVED***
		return false
	***REMOVED***
	emitter.whitespace = is_whitespace
	emitter.indention = (emitter.indention && is_indention)
	emitter.open_ended = false
	return true
***REMOVED***

func yaml_emitter_write_anchor(emitter *yaml_emitter_t, value []byte) bool ***REMOVED***
	if !write_all(emitter, value) ***REMOVED***
		return false
	***REMOVED***
	emitter.whitespace = false
	emitter.indention = false
	return true
***REMOVED***

func yaml_emitter_write_tag_handle(emitter *yaml_emitter_t, value []byte) bool ***REMOVED***
	if !emitter.whitespace ***REMOVED***
		if !put(emitter, ' ') ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	if !write_all(emitter, value) ***REMOVED***
		return false
	***REMOVED***
	emitter.whitespace = false
	emitter.indention = false
	return true
***REMOVED***

func yaml_emitter_write_tag_content(emitter *yaml_emitter_t, value []byte, need_whitespace bool) bool ***REMOVED***
	if need_whitespace && !emitter.whitespace ***REMOVED***
		if !put(emitter, ' ') ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	for i := 0; i < len(value); ***REMOVED***
		var must_write bool
		switch value[i] ***REMOVED***
		case ';', '/', '?', ':', '@', '&', '=', '+', '$', ',', '_', '.', '~', '*', '\'', '(', ')', '[', ']':
			must_write = true
		default:
			must_write = is_alpha(value, i)
		***REMOVED***
		if must_write ***REMOVED***
			if !write(emitter, value, &i) ***REMOVED***
				return false
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			w := width(value[i])
			for k := 0; k < w; k++ ***REMOVED***
				octet := value[i]
				i++
				if !put(emitter, '%') ***REMOVED***
					return false
				***REMOVED***

				c := octet >> 4
				if c < 10 ***REMOVED***
					c += '0'
				***REMOVED*** else ***REMOVED***
					c += 'A' - 10
				***REMOVED***
				if !put(emitter, c) ***REMOVED***
					return false
				***REMOVED***

				c = octet & 0x0f
				if c < 10 ***REMOVED***
					c += '0'
				***REMOVED*** else ***REMOVED***
					c += 'A' - 10
				***REMOVED***
				if !put(emitter, c) ***REMOVED***
					return false
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
	emitter.whitespace = false
	emitter.indention = false
	return true
***REMOVED***

func yaml_emitter_write_plain_scalar(emitter *yaml_emitter_t, value []byte, allow_breaks bool) bool ***REMOVED***
	if !emitter.whitespace ***REMOVED***
		if !put(emitter, ' ') ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***

	spaces := false
	breaks := false
	for i := 0; i < len(value); ***REMOVED***
		if is_space(value, i) ***REMOVED***
			if allow_breaks && !spaces && emitter.column > emitter.best_width && !is_space(value, i+1) ***REMOVED***
				if !yaml_emitter_write_indent(emitter) ***REMOVED***
					return false
				***REMOVED***
				i += width(value[i])
			***REMOVED*** else ***REMOVED***
				if !write(emitter, value, &i) ***REMOVED***
					return false
				***REMOVED***
			***REMOVED***
			spaces = true
		***REMOVED*** else if is_break(value, i) ***REMOVED***
			if !breaks && value[i] == '\n' ***REMOVED***
				if !put_break(emitter) ***REMOVED***
					return false
				***REMOVED***
			***REMOVED***
			if !write_break(emitter, value, &i) ***REMOVED***
				return false
			***REMOVED***
			emitter.indention = true
			breaks = true
		***REMOVED*** else ***REMOVED***
			if breaks ***REMOVED***
				if !yaml_emitter_write_indent(emitter) ***REMOVED***
					return false
				***REMOVED***
			***REMOVED***
			if !write(emitter, value, &i) ***REMOVED***
				return false
			***REMOVED***
			emitter.indention = false
			spaces = false
			breaks = false
		***REMOVED***
	***REMOVED***

	emitter.whitespace = false
	emitter.indention = false
	if emitter.root_context ***REMOVED***
		emitter.open_ended = true
	***REMOVED***

	return true
***REMOVED***

func yaml_emitter_write_single_quoted_scalar(emitter *yaml_emitter_t, value []byte, allow_breaks bool) bool ***REMOVED***

	if !yaml_emitter_write_indicator(emitter, []byte***REMOVED***'\''***REMOVED***, true, false, false) ***REMOVED***
		return false
	***REMOVED***

	spaces := false
	breaks := false
	for i := 0; i < len(value); ***REMOVED***
		if is_space(value, i) ***REMOVED***
			if allow_breaks && !spaces && emitter.column > emitter.best_width && i > 0 && i < len(value)-1 && !is_space(value, i+1) ***REMOVED***
				if !yaml_emitter_write_indent(emitter) ***REMOVED***
					return false
				***REMOVED***
				i += width(value[i])
			***REMOVED*** else ***REMOVED***
				if !write(emitter, value, &i) ***REMOVED***
					return false
				***REMOVED***
			***REMOVED***
			spaces = true
		***REMOVED*** else if is_break(value, i) ***REMOVED***
			if !breaks && value[i] == '\n' ***REMOVED***
				if !put_break(emitter) ***REMOVED***
					return false
				***REMOVED***
			***REMOVED***
			if !write_break(emitter, value, &i) ***REMOVED***
				return false
			***REMOVED***
			emitter.indention = true
			breaks = true
		***REMOVED*** else ***REMOVED***
			if breaks ***REMOVED***
				if !yaml_emitter_write_indent(emitter) ***REMOVED***
					return false
				***REMOVED***
			***REMOVED***
			if value[i] == '\'' ***REMOVED***
				if !put(emitter, '\'') ***REMOVED***
					return false
				***REMOVED***
			***REMOVED***
			if !write(emitter, value, &i) ***REMOVED***
				return false
			***REMOVED***
			emitter.indention = false
			spaces = false
			breaks = false
		***REMOVED***
	***REMOVED***
	if !yaml_emitter_write_indicator(emitter, []byte***REMOVED***'\''***REMOVED***, false, false, false) ***REMOVED***
		return false
	***REMOVED***
	emitter.whitespace = false
	emitter.indention = false
	return true
***REMOVED***

func yaml_emitter_write_double_quoted_scalar(emitter *yaml_emitter_t, value []byte, allow_breaks bool) bool ***REMOVED***
	spaces := false
	if !yaml_emitter_write_indicator(emitter, []byte***REMOVED***'"'***REMOVED***, true, false, false) ***REMOVED***
		return false
	***REMOVED***

	for i := 0; i < len(value); ***REMOVED***
		if !is_printable(value, i) || (!emitter.unicode && !is_ascii(value, i)) ||
			is_bom(value, i) || is_break(value, i) ||
			value[i] == '"' || value[i] == '\\' ***REMOVED***

			octet := value[i]

			var w int
			var v rune
			switch ***REMOVED***
			case octet&0x80 == 0x00:
				w, v = 1, rune(octet&0x7F)
			case octet&0xE0 == 0xC0:
				w, v = 2, rune(octet&0x1F)
			case octet&0xF0 == 0xE0:
				w, v = 3, rune(octet&0x0F)
			case octet&0xF8 == 0xF0:
				w, v = 4, rune(octet&0x07)
			***REMOVED***
			for k := 1; k < w; k++ ***REMOVED***
				octet = value[i+k]
				v = (v << 6) + (rune(octet) & 0x3F)
			***REMOVED***
			i += w

			if !put(emitter, '\\') ***REMOVED***
				return false
			***REMOVED***

			var ok bool
			switch v ***REMOVED***
			case 0x00:
				ok = put(emitter, '0')
			case 0x07:
				ok = put(emitter, 'a')
			case 0x08:
				ok = put(emitter, 'b')
			case 0x09:
				ok = put(emitter, 't')
			case 0x0A:
				ok = put(emitter, 'n')
			case 0x0b:
				ok = put(emitter, 'v')
			case 0x0c:
				ok = put(emitter, 'f')
			case 0x0d:
				ok = put(emitter, 'r')
			case 0x1b:
				ok = put(emitter, 'e')
			case 0x22:
				ok = put(emitter, '"')
			case 0x5c:
				ok = put(emitter, '\\')
			case 0x85:
				ok = put(emitter, 'N')
			case 0xA0:
				ok = put(emitter, '_')
			case 0x2028:
				ok = put(emitter, 'L')
			case 0x2029:
				ok = put(emitter, 'P')
			default:
				if v <= 0xFF ***REMOVED***
					ok = put(emitter, 'x')
					w = 2
				***REMOVED*** else if v <= 0xFFFF ***REMOVED***
					ok = put(emitter, 'u')
					w = 4
				***REMOVED*** else ***REMOVED***
					ok = put(emitter, 'U')
					w = 8
				***REMOVED***
				for k := (w - 1) * 4; ok && k >= 0; k -= 4 ***REMOVED***
					digit := byte((v >> uint(k)) & 0x0F)
					if digit < 10 ***REMOVED***
						ok = put(emitter, digit+'0')
					***REMOVED*** else ***REMOVED***
						ok = put(emitter, digit+'A'-10)
					***REMOVED***
				***REMOVED***
			***REMOVED***
			if !ok ***REMOVED***
				return false
			***REMOVED***
			spaces = false
		***REMOVED*** else if is_space(value, i) ***REMOVED***
			if allow_breaks && !spaces && emitter.column > emitter.best_width && i > 0 && i < len(value)-1 ***REMOVED***
				if !yaml_emitter_write_indent(emitter) ***REMOVED***
					return false
				***REMOVED***
				if is_space(value, i+1) ***REMOVED***
					if !put(emitter, '\\') ***REMOVED***
						return false
					***REMOVED***
				***REMOVED***
				i += width(value[i])
			***REMOVED*** else if !write(emitter, value, &i) ***REMOVED***
				return false
			***REMOVED***
			spaces = true
		***REMOVED*** else ***REMOVED***
			if !write(emitter, value, &i) ***REMOVED***
				return false
			***REMOVED***
			spaces = false
		***REMOVED***
	***REMOVED***
	if !yaml_emitter_write_indicator(emitter, []byte***REMOVED***'"'***REMOVED***, false, false, false) ***REMOVED***
		return false
	***REMOVED***
	emitter.whitespace = false
	emitter.indention = false
	return true
***REMOVED***

func yaml_emitter_write_block_scalar_hints(emitter *yaml_emitter_t, value []byte) bool ***REMOVED***
	if is_space(value, 0) || is_break(value, 0) ***REMOVED***
		indent_hint := []byte***REMOVED***'0' + byte(emitter.best_indent)***REMOVED***
		if !yaml_emitter_write_indicator(emitter, indent_hint, false, false, false) ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***

	emitter.open_ended = false

	var chomp_hint [1]byte
	if len(value) == 0 ***REMOVED***
		chomp_hint[0] = '-'
	***REMOVED*** else ***REMOVED***
		i := len(value) - 1
		for value[i]&0xC0 == 0x80 ***REMOVED***
			i--
		***REMOVED***
		if !is_break(value, i) ***REMOVED***
			chomp_hint[0] = '-'
		***REMOVED*** else if i == 0 ***REMOVED***
			chomp_hint[0] = '+'
			emitter.open_ended = true
		***REMOVED*** else ***REMOVED***
			i--
			for value[i]&0xC0 == 0x80 ***REMOVED***
				i--
			***REMOVED***
			if is_break(value, i) ***REMOVED***
				chomp_hint[0] = '+'
				emitter.open_ended = true
			***REMOVED***
		***REMOVED***
	***REMOVED***
	if chomp_hint[0] != 0 ***REMOVED***
		if !yaml_emitter_write_indicator(emitter, chomp_hint[:], false, false, false) ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

func yaml_emitter_write_literal_scalar(emitter *yaml_emitter_t, value []byte) bool ***REMOVED***
	if !yaml_emitter_write_indicator(emitter, []byte***REMOVED***'|'***REMOVED***, true, false, false) ***REMOVED***
		return false
	***REMOVED***
	if !yaml_emitter_write_block_scalar_hints(emitter, value) ***REMOVED***
		return false
	***REMOVED***
	if !put_break(emitter) ***REMOVED***
		return false
	***REMOVED***
	emitter.indention = true
	emitter.whitespace = true
	breaks := true
	for i := 0; i < len(value); ***REMOVED***
		if is_break(value, i) ***REMOVED***
			if !write_break(emitter, value, &i) ***REMOVED***
				return false
			***REMOVED***
			emitter.indention = true
			breaks = true
		***REMOVED*** else ***REMOVED***
			if breaks ***REMOVED***
				if !yaml_emitter_write_indent(emitter) ***REMOVED***
					return false
				***REMOVED***
			***REMOVED***
			if !write(emitter, value, &i) ***REMOVED***
				return false
			***REMOVED***
			emitter.indention = false
			breaks = false
		***REMOVED***
	***REMOVED***

	return true
***REMOVED***

func yaml_emitter_write_folded_scalar(emitter *yaml_emitter_t, value []byte) bool ***REMOVED***
	if !yaml_emitter_write_indicator(emitter, []byte***REMOVED***'>'***REMOVED***, true, false, false) ***REMOVED***
		return false
	***REMOVED***
	if !yaml_emitter_write_block_scalar_hints(emitter, value) ***REMOVED***
		return false
	***REMOVED***

	if !put_break(emitter) ***REMOVED***
		return false
	***REMOVED***
	emitter.indention = true
	emitter.whitespace = true

	breaks := true
	leading_spaces := true
	for i := 0; i < len(value); ***REMOVED***
		if is_break(value, i) ***REMOVED***
			if !breaks && !leading_spaces && value[i] == '\n' ***REMOVED***
				k := 0
				for is_break(value, k) ***REMOVED***
					k += width(value[k])
				***REMOVED***
				if !is_blankz(value, k) ***REMOVED***
					if !put_break(emitter) ***REMOVED***
						return false
					***REMOVED***
				***REMOVED***
			***REMOVED***
			if !write_break(emitter, value, &i) ***REMOVED***
				return false
			***REMOVED***
			emitter.indention = true
			breaks = true
		***REMOVED*** else ***REMOVED***
			if breaks ***REMOVED***
				if !yaml_emitter_write_indent(emitter) ***REMOVED***
					return false
				***REMOVED***
				leading_spaces = is_blank(value, i)
			***REMOVED***
			if !breaks && is_space(value, i) && !is_space(value, i+1) && emitter.column > emitter.best_width ***REMOVED***
				if !yaml_emitter_write_indent(emitter) ***REMOVED***
					return false
				***REMOVED***
				i += width(value[i])
			***REMOVED*** else ***REMOVED***
				if !write(emitter, value, &i) ***REMOVED***
					return false
				***REMOVED***
			***REMOVED***
			emitter.indention = false
			breaks = false
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***
