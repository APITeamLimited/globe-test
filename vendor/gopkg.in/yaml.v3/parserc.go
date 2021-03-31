//
// Copyright (c) 2011-2019 Canonical Ltd
// Copyright (c) 2006-2010 Kirill Simonov
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies
// of the Software, and to permit persons to whom the Software is furnished to do
// so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package yaml

import (
	"bytes"
)

// The parser implements the following grammar:
//
// stream               ::= STREAM-START implicit_document? explicit_document* STREAM-END
// implicit_document    ::= block_node DOCUMENT-END*
// explicit_document    ::= DIRECTIVE* DOCUMENT-START block_node? DOCUMENT-END*
// block_node_or_indentless_sequence    ::=
//                          ALIAS
//                          | properties (block_content | indentless_block_sequence)?
//                          | block_content
//                          | indentless_block_sequence
// block_node           ::= ALIAS
//                          | properties block_content?
//                          | block_content
// flow_node            ::= ALIAS
//                          | properties flow_content?
//                          | flow_content
// properties           ::= TAG ANCHOR? | ANCHOR TAG?
// block_content        ::= block_collection | flow_collection | SCALAR
// flow_content         ::= flow_collection | SCALAR
// block_collection     ::= block_sequence | block_mapping
// flow_collection      ::= flow_sequence | flow_mapping
// block_sequence       ::= BLOCK-SEQUENCE-START (BLOCK-ENTRY block_node?)* BLOCK-END
// indentless_sequence  ::= (BLOCK-ENTRY block_node?)+
// block_mapping        ::= BLOCK-MAPPING_START
//                          ((KEY block_node_or_indentless_sequence?)?
//                          (VALUE block_node_or_indentless_sequence?)?)*
//                          BLOCK-END
// flow_sequence        ::= FLOW-SEQUENCE-START
//                          (flow_sequence_entry FLOW-ENTRY)*
//                          flow_sequence_entry?
//                          FLOW-SEQUENCE-END
// flow_sequence_entry  ::= flow_node | KEY flow_node? (VALUE flow_node?)?
// flow_mapping         ::= FLOW-MAPPING-START
//                          (flow_mapping_entry FLOW-ENTRY)*
//                          flow_mapping_entry?
//                          FLOW-MAPPING-END
// flow_mapping_entry   ::= flow_node | KEY flow_node? (VALUE flow_node?)?

// Peek the next token in the token queue.
func peek_token(parser *yaml_parser_t) *yaml_token_t ***REMOVED***
	if parser.token_available || yaml_parser_fetch_more_tokens(parser) ***REMOVED***
		token := &parser.tokens[parser.tokens_head]
		yaml_parser_unfold_comments(parser, token)
		return token
	***REMOVED***
	return nil
***REMOVED***

// yaml_parser_unfold_comments walks through the comments queue and joins all
// comments behind the position of the provided token into the respective
// top-level comment slices in the parser.
func yaml_parser_unfold_comments(parser *yaml_parser_t, token *yaml_token_t) ***REMOVED***
	for parser.comments_head < len(parser.comments) && token.start_mark.index >= parser.comments[parser.comments_head].token_mark.index ***REMOVED***
		comment := &parser.comments[parser.comments_head]
		if len(comment.head) > 0 ***REMOVED***
			if token.typ == yaml_BLOCK_END_TOKEN ***REMOVED***
				// No heads on ends, so keep comment.head for a follow up token.
				break
			***REMOVED***
			if len(parser.head_comment) > 0 ***REMOVED***
				parser.head_comment = append(parser.head_comment, '\n')
			***REMOVED***
			parser.head_comment = append(parser.head_comment, comment.head...)
		***REMOVED***
		if len(comment.foot) > 0 ***REMOVED***
			if len(parser.foot_comment) > 0 ***REMOVED***
				parser.foot_comment = append(parser.foot_comment, '\n')
			***REMOVED***
			parser.foot_comment = append(parser.foot_comment, comment.foot...)
		***REMOVED***
		if len(comment.line) > 0 ***REMOVED***
			if len(parser.line_comment) > 0 ***REMOVED***
				parser.line_comment = append(parser.line_comment, '\n')
			***REMOVED***
			parser.line_comment = append(parser.line_comment, comment.line...)
		***REMOVED***
		*comment = yaml_comment_t***REMOVED******REMOVED***
		parser.comments_head++
	***REMOVED***
***REMOVED***

// Remove the next token from the queue (must be called after peek_token).
func skip_token(parser *yaml_parser_t) ***REMOVED***
	parser.token_available = false
	parser.tokens_parsed++
	parser.stream_end_produced = parser.tokens[parser.tokens_head].typ == yaml_STREAM_END_TOKEN
	parser.tokens_head++
***REMOVED***

// Get the next event.
func yaml_parser_parse(parser *yaml_parser_t, event *yaml_event_t) bool ***REMOVED***
	// Erase the event object.
	*event = yaml_event_t***REMOVED******REMOVED***

	// No events after the end of the stream or error.
	if parser.stream_end_produced || parser.error != yaml_NO_ERROR || parser.state == yaml_PARSE_END_STATE ***REMOVED***
		return true
	***REMOVED***

	// Generate the next event.
	return yaml_parser_state_machine(parser, event)
***REMOVED***

// Set parser error.
func yaml_parser_set_parser_error(parser *yaml_parser_t, problem string, problem_mark yaml_mark_t) bool ***REMOVED***
	parser.error = yaml_PARSER_ERROR
	parser.problem = problem
	parser.problem_mark = problem_mark
	return false
***REMOVED***

func yaml_parser_set_parser_error_context(parser *yaml_parser_t, context string, context_mark yaml_mark_t, problem string, problem_mark yaml_mark_t) bool ***REMOVED***
	parser.error = yaml_PARSER_ERROR
	parser.context = context
	parser.context_mark = context_mark
	parser.problem = problem
	parser.problem_mark = problem_mark
	return false
***REMOVED***

// State dispatcher.
func yaml_parser_state_machine(parser *yaml_parser_t, event *yaml_event_t) bool ***REMOVED***
	//trace("yaml_parser_state_machine", "state:", parser.state.String())

	switch parser.state ***REMOVED***
	case yaml_PARSE_STREAM_START_STATE:
		return yaml_parser_parse_stream_start(parser, event)

	case yaml_PARSE_IMPLICIT_DOCUMENT_START_STATE:
		return yaml_parser_parse_document_start(parser, event, true)

	case yaml_PARSE_DOCUMENT_START_STATE:
		return yaml_parser_parse_document_start(parser, event, false)

	case yaml_PARSE_DOCUMENT_CONTENT_STATE:
		return yaml_parser_parse_document_content(parser, event)

	case yaml_PARSE_DOCUMENT_END_STATE:
		return yaml_parser_parse_document_end(parser, event)

	case yaml_PARSE_BLOCK_NODE_STATE:
		return yaml_parser_parse_node(parser, event, true, false)

	case yaml_PARSE_BLOCK_NODE_OR_INDENTLESS_SEQUENCE_STATE:
		return yaml_parser_parse_node(parser, event, true, true)

	case yaml_PARSE_FLOW_NODE_STATE:
		return yaml_parser_parse_node(parser, event, false, false)

	case yaml_PARSE_BLOCK_SEQUENCE_FIRST_ENTRY_STATE:
		return yaml_parser_parse_block_sequence_entry(parser, event, true)

	case yaml_PARSE_BLOCK_SEQUENCE_ENTRY_STATE:
		return yaml_parser_parse_block_sequence_entry(parser, event, false)

	case yaml_PARSE_INDENTLESS_SEQUENCE_ENTRY_STATE:
		return yaml_parser_parse_indentless_sequence_entry(parser, event)

	case yaml_PARSE_BLOCK_MAPPING_FIRST_KEY_STATE:
		return yaml_parser_parse_block_mapping_key(parser, event, true)

	case yaml_PARSE_BLOCK_MAPPING_KEY_STATE:
		return yaml_parser_parse_block_mapping_key(parser, event, false)

	case yaml_PARSE_BLOCK_MAPPING_VALUE_STATE:
		return yaml_parser_parse_block_mapping_value(parser, event)

	case yaml_PARSE_FLOW_SEQUENCE_FIRST_ENTRY_STATE:
		return yaml_parser_parse_flow_sequence_entry(parser, event, true)

	case yaml_PARSE_FLOW_SEQUENCE_ENTRY_STATE:
		return yaml_parser_parse_flow_sequence_entry(parser, event, false)

	case yaml_PARSE_FLOW_SEQUENCE_ENTRY_MAPPING_KEY_STATE:
		return yaml_parser_parse_flow_sequence_entry_mapping_key(parser, event)

	case yaml_PARSE_FLOW_SEQUENCE_ENTRY_MAPPING_VALUE_STATE:
		return yaml_parser_parse_flow_sequence_entry_mapping_value(parser, event)

	case yaml_PARSE_FLOW_SEQUENCE_ENTRY_MAPPING_END_STATE:
		return yaml_parser_parse_flow_sequence_entry_mapping_end(parser, event)

	case yaml_PARSE_FLOW_MAPPING_FIRST_KEY_STATE:
		return yaml_parser_parse_flow_mapping_key(parser, event, true)

	case yaml_PARSE_FLOW_MAPPING_KEY_STATE:
		return yaml_parser_parse_flow_mapping_key(parser, event, false)

	case yaml_PARSE_FLOW_MAPPING_VALUE_STATE:
		return yaml_parser_parse_flow_mapping_value(parser, event, false)

	case yaml_PARSE_FLOW_MAPPING_EMPTY_VALUE_STATE:
		return yaml_parser_parse_flow_mapping_value(parser, event, true)

	default:
		panic("invalid parser state")
	***REMOVED***
***REMOVED***

// Parse the production:
// stream   ::= STREAM-START implicit_document? explicit_document* STREAM-END
//              ************
func yaml_parser_parse_stream_start(parser *yaml_parser_t, event *yaml_event_t) bool ***REMOVED***
	token := peek_token(parser)
	if token == nil ***REMOVED***
		return false
	***REMOVED***
	if token.typ != yaml_STREAM_START_TOKEN ***REMOVED***
		return yaml_parser_set_parser_error(parser, "did not find expected <stream-start>", token.start_mark)
	***REMOVED***
	parser.state = yaml_PARSE_IMPLICIT_DOCUMENT_START_STATE
	*event = yaml_event_t***REMOVED***
		typ:        yaml_STREAM_START_EVENT,
		start_mark: token.start_mark,
		end_mark:   token.end_mark,
		encoding:   token.encoding,
	***REMOVED***
	skip_token(parser)
	return true
***REMOVED***

// Parse the productions:
// implicit_document    ::= block_node DOCUMENT-END*
//                          *
// explicit_document    ::= DIRECTIVE* DOCUMENT-START block_node? DOCUMENT-END*
//                          *************************
func yaml_parser_parse_document_start(parser *yaml_parser_t, event *yaml_event_t, implicit bool) bool ***REMOVED***

	token := peek_token(parser)
	if token == nil ***REMOVED***
		return false
	***REMOVED***

	// Parse extra document end indicators.
	if !implicit ***REMOVED***
		for token.typ == yaml_DOCUMENT_END_TOKEN ***REMOVED***
			skip_token(parser)
			token = peek_token(parser)
			if token == nil ***REMOVED***
				return false
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if implicit && token.typ != yaml_VERSION_DIRECTIVE_TOKEN &&
		token.typ != yaml_TAG_DIRECTIVE_TOKEN &&
		token.typ != yaml_DOCUMENT_START_TOKEN &&
		token.typ != yaml_STREAM_END_TOKEN ***REMOVED***
		// Parse an implicit document.
		if !yaml_parser_process_directives(parser, nil, nil) ***REMOVED***
			return false
		***REMOVED***
		parser.states = append(parser.states, yaml_PARSE_DOCUMENT_END_STATE)
		parser.state = yaml_PARSE_BLOCK_NODE_STATE

		var head_comment []byte
		if len(parser.head_comment) > 0 ***REMOVED***
			// [Go] Scan the header comment backwards, and if an empty line is found, break
			//      the header so the part before the last empty line goes into the
			//      document header, while the bottom of it goes into a follow up event.
			for i := len(parser.head_comment) - 1; i > 0; i-- ***REMOVED***
				if parser.head_comment[i] == '\n' ***REMOVED***
					if i == len(parser.head_comment)-1 ***REMOVED***
						head_comment = parser.head_comment[:i]
						parser.head_comment = parser.head_comment[i+1:]
						break
					***REMOVED*** else if parser.head_comment[i-1] == '\n' ***REMOVED***
						head_comment = parser.head_comment[:i-1]
						parser.head_comment = parser.head_comment[i+1:]
						break
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED***

		*event = yaml_event_t***REMOVED***
			typ:        yaml_DOCUMENT_START_EVENT,
			start_mark: token.start_mark,
			end_mark:   token.end_mark,

			head_comment: head_comment,
		***REMOVED***

	***REMOVED*** else if token.typ != yaml_STREAM_END_TOKEN ***REMOVED***
		// Parse an explicit document.
		var version_directive *yaml_version_directive_t
		var tag_directives []yaml_tag_directive_t
		start_mark := token.start_mark
		if !yaml_parser_process_directives(parser, &version_directive, &tag_directives) ***REMOVED***
			return false
		***REMOVED***
		token = peek_token(parser)
		if token == nil ***REMOVED***
			return false
		***REMOVED***
		if token.typ != yaml_DOCUMENT_START_TOKEN ***REMOVED***
			yaml_parser_set_parser_error(parser,
				"did not find expected <document start>", token.start_mark)
			return false
		***REMOVED***
		parser.states = append(parser.states, yaml_PARSE_DOCUMENT_END_STATE)
		parser.state = yaml_PARSE_DOCUMENT_CONTENT_STATE
		end_mark := token.end_mark

		*event = yaml_event_t***REMOVED***
			typ:               yaml_DOCUMENT_START_EVENT,
			start_mark:        start_mark,
			end_mark:          end_mark,
			version_directive: version_directive,
			tag_directives:    tag_directives,
			implicit:          false,
		***REMOVED***
		skip_token(parser)

	***REMOVED*** else ***REMOVED***
		// Parse the stream end.
		parser.state = yaml_PARSE_END_STATE
		*event = yaml_event_t***REMOVED***
			typ:        yaml_STREAM_END_EVENT,
			start_mark: token.start_mark,
			end_mark:   token.end_mark,
		***REMOVED***
		skip_token(parser)
	***REMOVED***

	return true
***REMOVED***

// Parse the productions:
// explicit_document    ::= DIRECTIVE* DOCUMENT-START block_node? DOCUMENT-END*
//                                                    ***********
//
func yaml_parser_parse_document_content(parser *yaml_parser_t, event *yaml_event_t) bool ***REMOVED***
	token := peek_token(parser)
	if token == nil ***REMOVED***
		return false
	***REMOVED***

	if token.typ == yaml_VERSION_DIRECTIVE_TOKEN ||
		token.typ == yaml_TAG_DIRECTIVE_TOKEN ||
		token.typ == yaml_DOCUMENT_START_TOKEN ||
		token.typ == yaml_DOCUMENT_END_TOKEN ||
		token.typ == yaml_STREAM_END_TOKEN ***REMOVED***
		parser.state = parser.states[len(parser.states)-1]
		parser.states = parser.states[:len(parser.states)-1]
		return yaml_parser_process_empty_scalar(parser, event,
			token.start_mark)
	***REMOVED***
	return yaml_parser_parse_node(parser, event, true, false)
***REMOVED***

// Parse the productions:
// implicit_document    ::= block_node DOCUMENT-END*
//                                     *************
// explicit_document    ::= DIRECTIVE* DOCUMENT-START block_node? DOCUMENT-END*
//
func yaml_parser_parse_document_end(parser *yaml_parser_t, event *yaml_event_t) bool ***REMOVED***
	token := peek_token(parser)
	if token == nil ***REMOVED***
		return false
	***REMOVED***

	start_mark := token.start_mark
	end_mark := token.start_mark

	implicit := true
	if token.typ == yaml_DOCUMENT_END_TOKEN ***REMOVED***
		end_mark = token.end_mark
		skip_token(parser)
		implicit = false
	***REMOVED***

	parser.tag_directives = parser.tag_directives[:0]

	parser.state = yaml_PARSE_DOCUMENT_START_STATE
	*event = yaml_event_t***REMOVED***
		typ:        yaml_DOCUMENT_END_EVENT,
		start_mark: start_mark,
		end_mark:   end_mark,
		implicit:   implicit,
	***REMOVED***
	yaml_parser_set_event_comments(parser, event)
	if len(event.head_comment) > 0 && len(event.foot_comment) == 0 ***REMOVED***
		event.foot_comment = event.head_comment
		event.head_comment = nil
	***REMOVED***
	return true
***REMOVED***

func yaml_parser_set_event_comments(parser *yaml_parser_t, event *yaml_event_t) ***REMOVED***
	event.head_comment = parser.head_comment
	event.line_comment = parser.line_comment
	event.foot_comment = parser.foot_comment
	parser.head_comment = nil
	parser.line_comment = nil
	parser.foot_comment = nil
	parser.tail_comment = nil
	parser.stem_comment = nil
***REMOVED***

// Parse the productions:
// block_node_or_indentless_sequence    ::=
//                          ALIAS
//                          *****
//                          | properties (block_content | indentless_block_sequence)?
//                            **********  *
//                          | block_content | indentless_block_sequence
//                            *
// block_node           ::= ALIAS
//                          *****
//                          | properties block_content?
//                            ********** *
//                          | block_content
//                            *
// flow_node            ::= ALIAS
//                          *****
//                          | properties flow_content?
//                            ********** *
//                          | flow_content
//                            *
// properties           ::= TAG ANCHOR? | ANCHOR TAG?
//                          *************************
// block_content        ::= block_collection | flow_collection | SCALAR
//                                                               ******
// flow_content         ::= flow_collection | SCALAR
//                                            ******
func yaml_parser_parse_node(parser *yaml_parser_t, event *yaml_event_t, block, indentless_sequence bool) bool ***REMOVED***
	//defer trace("yaml_parser_parse_node", "block:", block, "indentless_sequence:", indentless_sequence)()

	token := peek_token(parser)
	if token == nil ***REMOVED***
		return false
	***REMOVED***

	if token.typ == yaml_ALIAS_TOKEN ***REMOVED***
		parser.state = parser.states[len(parser.states)-1]
		parser.states = parser.states[:len(parser.states)-1]
		*event = yaml_event_t***REMOVED***
			typ:        yaml_ALIAS_EVENT,
			start_mark: token.start_mark,
			end_mark:   token.end_mark,
			anchor:     token.value,
		***REMOVED***
		yaml_parser_set_event_comments(parser, event)
		skip_token(parser)
		return true
	***REMOVED***

	start_mark := token.start_mark
	end_mark := token.start_mark

	var tag_token bool
	var tag_handle, tag_suffix, anchor []byte
	var tag_mark yaml_mark_t
	if token.typ == yaml_ANCHOR_TOKEN ***REMOVED***
		anchor = token.value
		start_mark = token.start_mark
		end_mark = token.end_mark
		skip_token(parser)
		token = peek_token(parser)
		if token == nil ***REMOVED***
			return false
		***REMOVED***
		if token.typ == yaml_TAG_TOKEN ***REMOVED***
			tag_token = true
			tag_handle = token.value
			tag_suffix = token.suffix
			tag_mark = token.start_mark
			end_mark = token.end_mark
			skip_token(parser)
			token = peek_token(parser)
			if token == nil ***REMOVED***
				return false
			***REMOVED***
		***REMOVED***
	***REMOVED*** else if token.typ == yaml_TAG_TOKEN ***REMOVED***
		tag_token = true
		tag_handle = token.value
		tag_suffix = token.suffix
		start_mark = token.start_mark
		tag_mark = token.start_mark
		end_mark = token.end_mark
		skip_token(parser)
		token = peek_token(parser)
		if token == nil ***REMOVED***
			return false
		***REMOVED***
		if token.typ == yaml_ANCHOR_TOKEN ***REMOVED***
			anchor = token.value
			end_mark = token.end_mark
			skip_token(parser)
			token = peek_token(parser)
			if token == nil ***REMOVED***
				return false
			***REMOVED***
		***REMOVED***
	***REMOVED***

	var tag []byte
	if tag_token ***REMOVED***
		if len(tag_handle) == 0 ***REMOVED***
			tag = tag_suffix
			tag_suffix = nil
		***REMOVED*** else ***REMOVED***
			for i := range parser.tag_directives ***REMOVED***
				if bytes.Equal(parser.tag_directives[i].handle, tag_handle) ***REMOVED***
					tag = append([]byte(nil), parser.tag_directives[i].prefix...)
					tag = append(tag, tag_suffix...)
					break
				***REMOVED***
			***REMOVED***
			if len(tag) == 0 ***REMOVED***
				yaml_parser_set_parser_error_context(parser,
					"while parsing a node", start_mark,
					"found undefined tag handle", tag_mark)
				return false
			***REMOVED***
		***REMOVED***
	***REMOVED***

	implicit := len(tag) == 0
	if indentless_sequence && token.typ == yaml_BLOCK_ENTRY_TOKEN ***REMOVED***
		end_mark = token.end_mark
		parser.state = yaml_PARSE_INDENTLESS_SEQUENCE_ENTRY_STATE
		*event = yaml_event_t***REMOVED***
			typ:        yaml_SEQUENCE_START_EVENT,
			start_mark: start_mark,
			end_mark:   end_mark,
			anchor:     anchor,
			tag:        tag,
			implicit:   implicit,
			style:      yaml_style_t(yaml_BLOCK_SEQUENCE_STYLE),
		***REMOVED***
		return true
	***REMOVED***
	if token.typ == yaml_SCALAR_TOKEN ***REMOVED***
		var plain_implicit, quoted_implicit bool
		end_mark = token.end_mark
		if (len(tag) == 0 && token.style == yaml_PLAIN_SCALAR_STYLE) || (len(tag) == 1 && tag[0] == '!') ***REMOVED***
			plain_implicit = true
		***REMOVED*** else if len(tag) == 0 ***REMOVED***
			quoted_implicit = true
		***REMOVED***
		parser.state = parser.states[len(parser.states)-1]
		parser.states = parser.states[:len(parser.states)-1]

		*event = yaml_event_t***REMOVED***
			typ:             yaml_SCALAR_EVENT,
			start_mark:      start_mark,
			end_mark:        end_mark,
			anchor:          anchor,
			tag:             tag,
			value:           token.value,
			implicit:        plain_implicit,
			quoted_implicit: quoted_implicit,
			style:           yaml_style_t(token.style),
		***REMOVED***
		yaml_parser_set_event_comments(parser, event)
		skip_token(parser)
		return true
	***REMOVED***
	if token.typ == yaml_FLOW_SEQUENCE_START_TOKEN ***REMOVED***
		// [Go] Some of the events below can be merged as they differ only on style.
		end_mark = token.end_mark
		parser.state = yaml_PARSE_FLOW_SEQUENCE_FIRST_ENTRY_STATE
		*event = yaml_event_t***REMOVED***
			typ:        yaml_SEQUENCE_START_EVENT,
			start_mark: start_mark,
			end_mark:   end_mark,
			anchor:     anchor,
			tag:        tag,
			implicit:   implicit,
			style:      yaml_style_t(yaml_FLOW_SEQUENCE_STYLE),
		***REMOVED***
		yaml_parser_set_event_comments(parser, event)
		return true
	***REMOVED***
	if token.typ == yaml_FLOW_MAPPING_START_TOKEN ***REMOVED***
		end_mark = token.end_mark
		parser.state = yaml_PARSE_FLOW_MAPPING_FIRST_KEY_STATE
		*event = yaml_event_t***REMOVED***
			typ:        yaml_MAPPING_START_EVENT,
			start_mark: start_mark,
			end_mark:   end_mark,
			anchor:     anchor,
			tag:        tag,
			implicit:   implicit,
			style:      yaml_style_t(yaml_FLOW_MAPPING_STYLE),
		***REMOVED***
		yaml_parser_set_event_comments(parser, event)
		return true
	***REMOVED***
	if block && token.typ == yaml_BLOCK_SEQUENCE_START_TOKEN ***REMOVED***
		end_mark = token.end_mark
		parser.state = yaml_PARSE_BLOCK_SEQUENCE_FIRST_ENTRY_STATE
		*event = yaml_event_t***REMOVED***
			typ:        yaml_SEQUENCE_START_EVENT,
			start_mark: start_mark,
			end_mark:   end_mark,
			anchor:     anchor,
			tag:        tag,
			implicit:   implicit,
			style:      yaml_style_t(yaml_BLOCK_SEQUENCE_STYLE),
		***REMOVED***
		if parser.stem_comment != nil ***REMOVED***
			event.head_comment = parser.stem_comment
			parser.stem_comment = nil
		***REMOVED***
		return true
	***REMOVED***
	if block && token.typ == yaml_BLOCK_MAPPING_START_TOKEN ***REMOVED***
		end_mark = token.end_mark
		parser.state = yaml_PARSE_BLOCK_MAPPING_FIRST_KEY_STATE
		*event = yaml_event_t***REMOVED***
			typ:        yaml_MAPPING_START_EVENT,
			start_mark: start_mark,
			end_mark:   end_mark,
			anchor:     anchor,
			tag:        tag,
			implicit:   implicit,
			style:      yaml_style_t(yaml_BLOCK_MAPPING_STYLE),
		***REMOVED***
		return true
	***REMOVED***
	if len(anchor) > 0 || len(tag) > 0 ***REMOVED***
		parser.state = parser.states[len(parser.states)-1]
		parser.states = parser.states[:len(parser.states)-1]

		*event = yaml_event_t***REMOVED***
			typ:             yaml_SCALAR_EVENT,
			start_mark:      start_mark,
			end_mark:        end_mark,
			anchor:          anchor,
			tag:             tag,
			implicit:        implicit,
			quoted_implicit: false,
			style:           yaml_style_t(yaml_PLAIN_SCALAR_STYLE),
		***REMOVED***
		return true
	***REMOVED***

	context := "while parsing a flow node"
	if block ***REMOVED***
		context = "while parsing a block node"
	***REMOVED***
	yaml_parser_set_parser_error_context(parser, context, start_mark,
		"did not find expected node content", token.start_mark)
	return false
***REMOVED***

// Parse the productions:
// block_sequence ::= BLOCK-SEQUENCE-START (BLOCK-ENTRY block_node?)* BLOCK-END
//                    ********************  *********** *             *********
//
func yaml_parser_parse_block_sequence_entry(parser *yaml_parser_t, event *yaml_event_t, first bool) bool ***REMOVED***
	if first ***REMOVED***
		token := peek_token(parser)
		parser.marks = append(parser.marks, token.start_mark)
		skip_token(parser)
	***REMOVED***

	token := peek_token(parser)
	if token == nil ***REMOVED***
		return false
	***REMOVED***

	if token.typ == yaml_BLOCK_ENTRY_TOKEN ***REMOVED***
		mark := token.end_mark
		prior_head := len(parser.head_comment)
		skip_token(parser)
		token = peek_token(parser)
		if token == nil ***REMOVED***
			return false
		***REMOVED***
		if prior_head > 0 && token.typ == yaml_BLOCK_SEQUENCE_START_TOKEN ***REMOVED***
			// [Go] It's a sequence under a sequence entry, so the former head comment
			//      is for the list itself, not the first list item under it.
			parser.stem_comment = parser.head_comment[:prior_head]
			if len(parser.head_comment) == prior_head ***REMOVED***
				parser.head_comment = nil
			***REMOVED*** else ***REMOVED***
				// Copy suffix to prevent very strange bugs if someone ever appends
				// further bytes to the prefix in the stem_comment slice above.
				parser.head_comment = append([]byte(nil), parser.head_comment[prior_head+1:]...)
			***REMOVED***

		***REMOVED***
		if token.typ != yaml_BLOCK_ENTRY_TOKEN && token.typ != yaml_BLOCK_END_TOKEN ***REMOVED***
			parser.states = append(parser.states, yaml_PARSE_BLOCK_SEQUENCE_ENTRY_STATE)
			return yaml_parser_parse_node(parser, event, true, false)
		***REMOVED*** else ***REMOVED***
			parser.state = yaml_PARSE_BLOCK_SEQUENCE_ENTRY_STATE
			return yaml_parser_process_empty_scalar(parser, event, mark)
		***REMOVED***
	***REMOVED***
	if token.typ == yaml_BLOCK_END_TOKEN ***REMOVED***
		parser.state = parser.states[len(parser.states)-1]
		parser.states = parser.states[:len(parser.states)-1]
		parser.marks = parser.marks[:len(parser.marks)-1]

		*event = yaml_event_t***REMOVED***
			typ:        yaml_SEQUENCE_END_EVENT,
			start_mark: token.start_mark,
			end_mark:   token.end_mark,
		***REMOVED***

		skip_token(parser)
		return true
	***REMOVED***

	context_mark := parser.marks[len(parser.marks)-1]
	parser.marks = parser.marks[:len(parser.marks)-1]
	return yaml_parser_set_parser_error_context(parser,
		"while parsing a block collection", context_mark,
		"did not find expected '-' indicator", token.start_mark)
***REMOVED***

// Parse the productions:
// indentless_sequence  ::= (BLOCK-ENTRY block_node?)+
//                           *********** *
func yaml_parser_parse_indentless_sequence_entry(parser *yaml_parser_t, event *yaml_event_t) bool ***REMOVED***
	token := peek_token(parser)
	if token == nil ***REMOVED***
		return false
	***REMOVED***

	if token.typ == yaml_BLOCK_ENTRY_TOKEN ***REMOVED***
		mark := token.end_mark
		skip_token(parser)
		token = peek_token(parser)
		if token == nil ***REMOVED***
			return false
		***REMOVED***
		if token.typ != yaml_BLOCK_ENTRY_TOKEN &&
			token.typ != yaml_KEY_TOKEN &&
			token.typ != yaml_VALUE_TOKEN &&
			token.typ != yaml_BLOCK_END_TOKEN ***REMOVED***
			parser.states = append(parser.states, yaml_PARSE_INDENTLESS_SEQUENCE_ENTRY_STATE)
			return yaml_parser_parse_node(parser, event, true, false)
		***REMOVED***
		parser.state = yaml_PARSE_INDENTLESS_SEQUENCE_ENTRY_STATE
		return yaml_parser_process_empty_scalar(parser, event, mark)
	***REMOVED***
	parser.state = parser.states[len(parser.states)-1]
	parser.states = parser.states[:len(parser.states)-1]

	*event = yaml_event_t***REMOVED***
		typ:        yaml_SEQUENCE_END_EVENT,
		start_mark: token.start_mark,
		end_mark:   token.start_mark, // [Go] Shouldn't this be token.end_mark?
	***REMOVED***
	return true
***REMOVED***

// Parse the productions:
// block_mapping        ::= BLOCK-MAPPING_START
//                          *******************
//                          ((KEY block_node_or_indentless_sequence?)?
//                            *** *
//                          (VALUE block_node_or_indentless_sequence?)?)*
//
//                          BLOCK-END
//                          *********
//
func yaml_parser_parse_block_mapping_key(parser *yaml_parser_t, event *yaml_event_t, first bool) bool ***REMOVED***
	if first ***REMOVED***
		token := peek_token(parser)
		parser.marks = append(parser.marks, token.start_mark)
		skip_token(parser)
	***REMOVED***

	token := peek_token(parser)
	if token == nil ***REMOVED***
		return false
	***REMOVED***

	// [Go] A tail comment was left from the prior mapping value processed. Emit an event
	//      as it needs to be processed with that value and not the following key.
	if len(parser.tail_comment) > 0 ***REMOVED***
		*event = yaml_event_t***REMOVED***
			typ:          yaml_TAIL_COMMENT_EVENT,
			start_mark:   token.start_mark,
			end_mark:     token.end_mark,
			foot_comment: parser.tail_comment,
		***REMOVED***
		parser.tail_comment = nil
		return true
	***REMOVED***

	if token.typ == yaml_KEY_TOKEN ***REMOVED***
		mark := token.end_mark
		skip_token(parser)
		token = peek_token(parser)
		if token == nil ***REMOVED***
			return false
		***REMOVED***
		if token.typ != yaml_KEY_TOKEN &&
			token.typ != yaml_VALUE_TOKEN &&
			token.typ != yaml_BLOCK_END_TOKEN ***REMOVED***
			parser.states = append(parser.states, yaml_PARSE_BLOCK_MAPPING_VALUE_STATE)
			return yaml_parser_parse_node(parser, event, true, true)
		***REMOVED*** else ***REMOVED***
			parser.state = yaml_PARSE_BLOCK_MAPPING_VALUE_STATE
			return yaml_parser_process_empty_scalar(parser, event, mark)
		***REMOVED***
	***REMOVED*** else if token.typ == yaml_BLOCK_END_TOKEN ***REMOVED***
		parser.state = parser.states[len(parser.states)-1]
		parser.states = parser.states[:len(parser.states)-1]
		parser.marks = parser.marks[:len(parser.marks)-1]
		*event = yaml_event_t***REMOVED***
			typ:        yaml_MAPPING_END_EVENT,
			start_mark: token.start_mark,
			end_mark:   token.end_mark,
		***REMOVED***
		yaml_parser_set_event_comments(parser, event)
		skip_token(parser)
		return true
	***REMOVED***

	context_mark := parser.marks[len(parser.marks)-1]
	parser.marks = parser.marks[:len(parser.marks)-1]
	return yaml_parser_set_parser_error_context(parser,
		"while parsing a block mapping", context_mark,
		"did not find expected key", token.start_mark)
***REMOVED***

// Parse the productions:
// block_mapping        ::= BLOCK-MAPPING_START
//
//                          ((KEY block_node_or_indentless_sequence?)?
//
//                          (VALUE block_node_or_indentless_sequence?)?)*
//                           ***** *
//                          BLOCK-END
//
//
func yaml_parser_parse_block_mapping_value(parser *yaml_parser_t, event *yaml_event_t) bool ***REMOVED***
	token := peek_token(parser)
	if token == nil ***REMOVED***
		return false
	***REMOVED***
	if token.typ == yaml_VALUE_TOKEN ***REMOVED***
		mark := token.end_mark
		skip_token(parser)
		token = peek_token(parser)
		if token == nil ***REMOVED***
			return false
		***REMOVED***
		if token.typ != yaml_KEY_TOKEN &&
			token.typ != yaml_VALUE_TOKEN &&
			token.typ != yaml_BLOCK_END_TOKEN ***REMOVED***
			parser.states = append(parser.states, yaml_PARSE_BLOCK_MAPPING_KEY_STATE)
			return yaml_parser_parse_node(parser, event, true, true)
		***REMOVED***
		parser.state = yaml_PARSE_BLOCK_MAPPING_KEY_STATE
		return yaml_parser_process_empty_scalar(parser, event, mark)
	***REMOVED***
	parser.state = yaml_PARSE_BLOCK_MAPPING_KEY_STATE
	return yaml_parser_process_empty_scalar(parser, event, token.start_mark)
***REMOVED***

// Parse the productions:
// flow_sequence        ::= FLOW-SEQUENCE-START
//                          *******************
//                          (flow_sequence_entry FLOW-ENTRY)*
//                           *                   **********
//                          flow_sequence_entry?
//                          *
//                          FLOW-SEQUENCE-END
//                          *****************
// flow_sequence_entry  ::= flow_node | KEY flow_node? (VALUE flow_node?)?
//                          *
//
func yaml_parser_parse_flow_sequence_entry(parser *yaml_parser_t, event *yaml_event_t, first bool) bool ***REMOVED***
	if first ***REMOVED***
		token := peek_token(parser)
		parser.marks = append(parser.marks, token.start_mark)
		skip_token(parser)
	***REMOVED***
	token := peek_token(parser)
	if token == nil ***REMOVED***
		return false
	***REMOVED***
	if token.typ != yaml_FLOW_SEQUENCE_END_TOKEN ***REMOVED***
		if !first ***REMOVED***
			if token.typ == yaml_FLOW_ENTRY_TOKEN ***REMOVED***
				skip_token(parser)
				token = peek_token(parser)
				if token == nil ***REMOVED***
					return false
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				context_mark := parser.marks[len(parser.marks)-1]
				parser.marks = parser.marks[:len(parser.marks)-1]
				return yaml_parser_set_parser_error_context(parser,
					"while parsing a flow sequence", context_mark,
					"did not find expected ',' or ']'", token.start_mark)
			***REMOVED***
		***REMOVED***

		if token.typ == yaml_KEY_TOKEN ***REMOVED***
			parser.state = yaml_PARSE_FLOW_SEQUENCE_ENTRY_MAPPING_KEY_STATE
			*event = yaml_event_t***REMOVED***
				typ:        yaml_MAPPING_START_EVENT,
				start_mark: token.start_mark,
				end_mark:   token.end_mark,
				implicit:   true,
				style:      yaml_style_t(yaml_FLOW_MAPPING_STYLE),
			***REMOVED***
			skip_token(parser)
			return true
		***REMOVED*** else if token.typ != yaml_FLOW_SEQUENCE_END_TOKEN ***REMOVED***
			parser.states = append(parser.states, yaml_PARSE_FLOW_SEQUENCE_ENTRY_STATE)
			return yaml_parser_parse_node(parser, event, false, false)
		***REMOVED***
	***REMOVED***

	parser.state = parser.states[len(parser.states)-1]
	parser.states = parser.states[:len(parser.states)-1]
	parser.marks = parser.marks[:len(parser.marks)-1]

	*event = yaml_event_t***REMOVED***
		typ:        yaml_SEQUENCE_END_EVENT,
		start_mark: token.start_mark,
		end_mark:   token.end_mark,
	***REMOVED***
	yaml_parser_set_event_comments(parser, event)

	skip_token(parser)
	return true
***REMOVED***

//
// Parse the productions:
// flow_sequence_entry  ::= flow_node | KEY flow_node? (VALUE flow_node?)?
//                                      *** *
//
func yaml_parser_parse_flow_sequence_entry_mapping_key(parser *yaml_parser_t, event *yaml_event_t) bool ***REMOVED***
	token := peek_token(parser)
	if token == nil ***REMOVED***
		return false
	***REMOVED***
	if token.typ != yaml_VALUE_TOKEN &&
		token.typ != yaml_FLOW_ENTRY_TOKEN &&
		token.typ != yaml_FLOW_SEQUENCE_END_TOKEN ***REMOVED***
		parser.states = append(parser.states, yaml_PARSE_FLOW_SEQUENCE_ENTRY_MAPPING_VALUE_STATE)
		return yaml_parser_parse_node(parser, event, false, false)
	***REMOVED***
	mark := token.end_mark
	skip_token(parser)
	parser.state = yaml_PARSE_FLOW_SEQUENCE_ENTRY_MAPPING_VALUE_STATE
	return yaml_parser_process_empty_scalar(parser, event, mark)
***REMOVED***

// Parse the productions:
// flow_sequence_entry  ::= flow_node | KEY flow_node? (VALUE flow_node?)?
//                                                      ***** *
//
func yaml_parser_parse_flow_sequence_entry_mapping_value(parser *yaml_parser_t, event *yaml_event_t) bool ***REMOVED***
	token := peek_token(parser)
	if token == nil ***REMOVED***
		return false
	***REMOVED***
	if token.typ == yaml_VALUE_TOKEN ***REMOVED***
		skip_token(parser)
		token := peek_token(parser)
		if token == nil ***REMOVED***
			return false
		***REMOVED***
		if token.typ != yaml_FLOW_ENTRY_TOKEN && token.typ != yaml_FLOW_SEQUENCE_END_TOKEN ***REMOVED***
			parser.states = append(parser.states, yaml_PARSE_FLOW_SEQUENCE_ENTRY_MAPPING_END_STATE)
			return yaml_parser_parse_node(parser, event, false, false)
		***REMOVED***
	***REMOVED***
	parser.state = yaml_PARSE_FLOW_SEQUENCE_ENTRY_MAPPING_END_STATE
	return yaml_parser_process_empty_scalar(parser, event, token.start_mark)
***REMOVED***

// Parse the productions:
// flow_sequence_entry  ::= flow_node | KEY flow_node? (VALUE flow_node?)?
//                                                                      *
//
func yaml_parser_parse_flow_sequence_entry_mapping_end(parser *yaml_parser_t, event *yaml_event_t) bool ***REMOVED***
	token := peek_token(parser)
	if token == nil ***REMOVED***
		return false
	***REMOVED***
	parser.state = yaml_PARSE_FLOW_SEQUENCE_ENTRY_STATE
	*event = yaml_event_t***REMOVED***
		typ:        yaml_MAPPING_END_EVENT,
		start_mark: token.start_mark,
		end_mark:   token.start_mark, // [Go] Shouldn't this be end_mark?
	***REMOVED***
	return true
***REMOVED***

// Parse the productions:
// flow_mapping         ::= FLOW-MAPPING-START
//                          ******************
//                          (flow_mapping_entry FLOW-ENTRY)*
//                           *                  **********
//                          flow_mapping_entry?
//                          ******************
//                          FLOW-MAPPING-END
//                          ****************
// flow_mapping_entry   ::= flow_node | KEY flow_node? (VALUE flow_node?)?
//                          *           *** *
//
func yaml_parser_parse_flow_mapping_key(parser *yaml_parser_t, event *yaml_event_t, first bool) bool ***REMOVED***
	if first ***REMOVED***
		token := peek_token(parser)
		parser.marks = append(parser.marks, token.start_mark)
		skip_token(parser)
	***REMOVED***

	token := peek_token(parser)
	if token == nil ***REMOVED***
		return false
	***REMOVED***

	if token.typ != yaml_FLOW_MAPPING_END_TOKEN ***REMOVED***
		if !first ***REMOVED***
			if token.typ == yaml_FLOW_ENTRY_TOKEN ***REMOVED***
				skip_token(parser)
				token = peek_token(parser)
				if token == nil ***REMOVED***
					return false
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				context_mark := parser.marks[len(parser.marks)-1]
				parser.marks = parser.marks[:len(parser.marks)-1]
				return yaml_parser_set_parser_error_context(parser,
					"while parsing a flow mapping", context_mark,
					"did not find expected ',' or '***REMOVED***'", token.start_mark)
			***REMOVED***
		***REMOVED***

		if token.typ == yaml_KEY_TOKEN ***REMOVED***
			skip_token(parser)
			token = peek_token(parser)
			if token == nil ***REMOVED***
				return false
			***REMOVED***
			if token.typ != yaml_VALUE_TOKEN &&
				token.typ != yaml_FLOW_ENTRY_TOKEN &&
				token.typ != yaml_FLOW_MAPPING_END_TOKEN ***REMOVED***
				parser.states = append(parser.states, yaml_PARSE_FLOW_MAPPING_VALUE_STATE)
				return yaml_parser_parse_node(parser, event, false, false)
			***REMOVED*** else ***REMOVED***
				parser.state = yaml_PARSE_FLOW_MAPPING_VALUE_STATE
				return yaml_parser_process_empty_scalar(parser, event, token.start_mark)
			***REMOVED***
		***REMOVED*** else if token.typ != yaml_FLOW_MAPPING_END_TOKEN ***REMOVED***
			parser.states = append(parser.states, yaml_PARSE_FLOW_MAPPING_EMPTY_VALUE_STATE)
			return yaml_parser_parse_node(parser, event, false, false)
		***REMOVED***
	***REMOVED***

	parser.state = parser.states[len(parser.states)-1]
	parser.states = parser.states[:len(parser.states)-1]
	parser.marks = parser.marks[:len(parser.marks)-1]
	*event = yaml_event_t***REMOVED***
		typ:        yaml_MAPPING_END_EVENT,
		start_mark: token.start_mark,
		end_mark:   token.end_mark,
	***REMOVED***
	yaml_parser_set_event_comments(parser, event)
	skip_token(parser)
	return true
***REMOVED***

// Parse the productions:
// flow_mapping_entry   ::= flow_node | KEY flow_node? (VALUE flow_node?)?
//                                   *                  ***** *
//
func yaml_parser_parse_flow_mapping_value(parser *yaml_parser_t, event *yaml_event_t, empty bool) bool ***REMOVED***
	token := peek_token(parser)
	if token == nil ***REMOVED***
		return false
	***REMOVED***
	if empty ***REMOVED***
		parser.state = yaml_PARSE_FLOW_MAPPING_KEY_STATE
		return yaml_parser_process_empty_scalar(parser, event, token.start_mark)
	***REMOVED***
	if token.typ == yaml_VALUE_TOKEN ***REMOVED***
		skip_token(parser)
		token = peek_token(parser)
		if token == nil ***REMOVED***
			return false
		***REMOVED***
		if token.typ != yaml_FLOW_ENTRY_TOKEN && token.typ != yaml_FLOW_MAPPING_END_TOKEN ***REMOVED***
			parser.states = append(parser.states, yaml_PARSE_FLOW_MAPPING_KEY_STATE)
			return yaml_parser_parse_node(parser, event, false, false)
		***REMOVED***
	***REMOVED***
	parser.state = yaml_PARSE_FLOW_MAPPING_KEY_STATE
	return yaml_parser_process_empty_scalar(parser, event, token.start_mark)
***REMOVED***

// Generate an empty scalar event.
func yaml_parser_process_empty_scalar(parser *yaml_parser_t, event *yaml_event_t, mark yaml_mark_t) bool ***REMOVED***
	*event = yaml_event_t***REMOVED***
		typ:        yaml_SCALAR_EVENT,
		start_mark: mark,
		end_mark:   mark,
		value:      nil, // Empty
		implicit:   true,
		style:      yaml_style_t(yaml_PLAIN_SCALAR_STYLE),
	***REMOVED***
	return true
***REMOVED***

var default_tag_directives = []yaml_tag_directive_t***REMOVED***
	***REMOVED***[]byte("!"), []byte("!")***REMOVED***,
	***REMOVED***[]byte("!!"), []byte("tag:yaml.org,2002:")***REMOVED***,
***REMOVED***

// Parse directives.
func yaml_parser_process_directives(parser *yaml_parser_t,
	version_directive_ref **yaml_version_directive_t,
	tag_directives_ref *[]yaml_tag_directive_t) bool ***REMOVED***

	var version_directive *yaml_version_directive_t
	var tag_directives []yaml_tag_directive_t

	token := peek_token(parser)
	if token == nil ***REMOVED***
		return false
	***REMOVED***

	for token.typ == yaml_VERSION_DIRECTIVE_TOKEN || token.typ == yaml_TAG_DIRECTIVE_TOKEN ***REMOVED***
		if token.typ == yaml_VERSION_DIRECTIVE_TOKEN ***REMOVED***
			if version_directive != nil ***REMOVED***
				yaml_parser_set_parser_error(parser,
					"found duplicate %YAML directive", token.start_mark)
				return false
			***REMOVED***
			if token.major != 1 || token.minor != 1 ***REMOVED***
				yaml_parser_set_parser_error(parser,
					"found incompatible YAML document", token.start_mark)
				return false
			***REMOVED***
			version_directive = &yaml_version_directive_t***REMOVED***
				major: token.major,
				minor: token.minor,
			***REMOVED***
		***REMOVED*** else if token.typ == yaml_TAG_DIRECTIVE_TOKEN ***REMOVED***
			value := yaml_tag_directive_t***REMOVED***
				handle: token.value,
				prefix: token.prefix,
			***REMOVED***
			if !yaml_parser_append_tag_directive(parser, value, false, token.start_mark) ***REMOVED***
				return false
			***REMOVED***
			tag_directives = append(tag_directives, value)
		***REMOVED***

		skip_token(parser)
		token = peek_token(parser)
		if token == nil ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***

	for i := range default_tag_directives ***REMOVED***
		if !yaml_parser_append_tag_directive(parser, default_tag_directives[i], true, token.start_mark) ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***

	if version_directive_ref != nil ***REMOVED***
		*version_directive_ref = version_directive
	***REMOVED***
	if tag_directives_ref != nil ***REMOVED***
		*tag_directives_ref = tag_directives
	***REMOVED***
	return true
***REMOVED***

// Append a tag directive to the directives stack.
func yaml_parser_append_tag_directive(parser *yaml_parser_t, value yaml_tag_directive_t, allow_duplicates bool, mark yaml_mark_t) bool ***REMOVED***
	for i := range parser.tag_directives ***REMOVED***
		if bytes.Equal(value.handle, parser.tag_directives[i].handle) ***REMOVED***
			if allow_duplicates ***REMOVED***
				return true
			***REMOVED***
			return yaml_parser_set_parser_error(parser, "found duplicate %TAG directive", mark)
		***REMOVED***
	***REMOVED***

	// [Go] I suspect the copy is unnecessary. This was likely done
	// because there was no way to track ownership of the data.
	value_copy := yaml_tag_directive_t***REMOVED***
		handle: make([]byte, len(value.handle)),
		prefix: make([]byte, len(value.prefix)),
	***REMOVED***
	copy(value_copy.handle, value.handle)
	copy(value_copy.prefix, value.prefix)
	parser.tag_directives = append(parser.tag_directives, value_copy)
	return true
***REMOVED***
