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
	"fmt"
)

// Introduction
// ************
//
// The following notes assume that you are familiar with the YAML specification
// (http://yaml.org/spec/1.2/spec.html).  We mostly follow it, although in
// some cases we are less restrictive that it requires.
//
// The process of transforming a YAML stream into a sequence of events is
// divided on two steps: Scanning and Parsing.
//
// The Scanner transforms the input stream into a sequence of tokens, while the
// parser transform the sequence of tokens produced by the Scanner into a
// sequence of parsing events.
//
// The Scanner is rather clever and complicated. The Parser, on the contrary,
// is a straightforward implementation of a recursive-descendant parser (or,
// LL(1) parser, as it is usually called).
//
// Actually there are two issues of Scanning that might be called "clever", the
// rest is quite straightforward.  The issues are "block collection start" and
// "simple keys".  Both issues are explained below in details.
//
// Here the Scanning step is explained and implemented.  We start with the list
// of all the tokens produced by the Scanner together with short descriptions.
//
// Now, tokens:
//
//      STREAM-START(encoding)          # The stream start.
//      STREAM-END                      # The stream end.
//      VERSION-DIRECTIVE(major,minor)  # The '%YAML' directive.
//      TAG-DIRECTIVE(handle,prefix)    # The '%TAG' directive.
//      DOCUMENT-START                  # '---'
//      DOCUMENT-END                    # '...'
//      BLOCK-SEQUENCE-START            # Indentation increase denoting a block
//      BLOCK-MAPPING-START             # sequence or a block mapping.
//      BLOCK-END                       # Indentation decrease.
//      FLOW-SEQUENCE-START             # '['
//      FLOW-SEQUENCE-END               # ']'
//      BLOCK-SEQUENCE-START            # '***REMOVED***'
//      BLOCK-SEQUENCE-END              # '***REMOVED***'
//      BLOCK-ENTRY                     # '-'
//      FLOW-ENTRY                      # ','
//      KEY                             # '?' or nothing (simple keys).
//      VALUE                           # ':'
//      ALIAS(anchor)                   # '*anchor'
//      ANCHOR(anchor)                  # '&anchor'
//      TAG(handle,suffix)              # '!handle!suffix'
//      SCALAR(value,style)             # A scalar.
//
// The following two tokens are "virtual" tokens denoting the beginning and the
// end of the stream:
//
//      STREAM-START(encoding)
//      STREAM-END
//
// We pass the information about the input stream encoding with the
// STREAM-START token.
//
// The next two tokens are responsible for tags:
//
//      VERSION-DIRECTIVE(major,minor)
//      TAG-DIRECTIVE(handle,prefix)
//
// Example:
//
//      %YAML   1.1
//      %TAG    !   !foo
//      %TAG    !yaml!  tag:yaml.org,2002:
//      ---
//
// The correspoding sequence of tokens:
//
//      STREAM-START(utf-8)
//      VERSION-DIRECTIVE(1,1)
//      TAG-DIRECTIVE("!","!foo")
//      TAG-DIRECTIVE("!yaml","tag:yaml.org,2002:")
//      DOCUMENT-START
//      STREAM-END
//
// Note that the VERSION-DIRECTIVE and TAG-DIRECTIVE tokens occupy a whole
// line.
//
// The document start and end indicators are represented by:
//
//      DOCUMENT-START
//      DOCUMENT-END
//
// Note that if a YAML stream contains an implicit document (without '---'
// and '...' indicators), no DOCUMENT-START and DOCUMENT-END tokens will be
// produced.
//
// In the following examples, we present whole documents together with the
// produced tokens.
//
//      1. An implicit document:
//
//          'a scalar'
//
//      Tokens:
//
//          STREAM-START(utf-8)
//          SCALAR("a scalar",single-quoted)
//          STREAM-END
//
//      2. An explicit document:
//
//          ---
//          'a scalar'
//          ...
//
//      Tokens:
//
//          STREAM-START(utf-8)
//          DOCUMENT-START
//          SCALAR("a scalar",single-quoted)
//          DOCUMENT-END
//          STREAM-END
//
//      3. Several documents in a stream:
//
//          'a scalar'
//          ---
//          'another scalar'
//          ---
//          'yet another scalar'
//
//      Tokens:
//
//          STREAM-START(utf-8)
//          SCALAR("a scalar",single-quoted)
//          DOCUMENT-START
//          SCALAR("another scalar",single-quoted)
//          DOCUMENT-START
//          SCALAR("yet another scalar",single-quoted)
//          STREAM-END
//
// We have already introduced the SCALAR token above.  The following tokens are
// used to describe aliases, anchors, tag, and scalars:
//
//      ALIAS(anchor)
//      ANCHOR(anchor)
//      TAG(handle,suffix)
//      SCALAR(value,style)
//
// The following series of examples illustrate the usage of these tokens:
//
//      1. A recursive sequence:
//
//          &A [ *A ]
//
//      Tokens:
//
//          STREAM-START(utf-8)
//          ANCHOR("A")
//          FLOW-SEQUENCE-START
//          ALIAS("A")
//          FLOW-SEQUENCE-END
//          STREAM-END
//
//      2. A tagged scalar:
//
//          !!float "3.14"  # A good approximation.
//
//      Tokens:
//
//          STREAM-START(utf-8)
//          TAG("!!","float")
//          SCALAR("3.14",double-quoted)
//          STREAM-END
//
//      3. Various scalar styles:
//
//          --- # Implicit empty plain scalars do not produce tokens.
//          --- a plain scalar
//          --- 'a single-quoted scalar'
//          --- "a double-quoted scalar"
//          --- |-
//            a literal scalar
//          --- >-
//            a folded
//            scalar
//
//      Tokens:
//
//          STREAM-START(utf-8)
//          DOCUMENT-START
//          DOCUMENT-START
//          SCALAR("a plain scalar",plain)
//          DOCUMENT-START
//          SCALAR("a single-quoted scalar",single-quoted)
//          DOCUMENT-START
//          SCALAR("a double-quoted scalar",double-quoted)
//          DOCUMENT-START
//          SCALAR("a literal scalar",literal)
//          DOCUMENT-START
//          SCALAR("a folded scalar",folded)
//          STREAM-END
//
// Now it's time to review collection-related tokens. We will start with
// flow collections:
//
//      FLOW-SEQUENCE-START
//      FLOW-SEQUENCE-END
//      FLOW-MAPPING-START
//      FLOW-MAPPING-END
//      FLOW-ENTRY
//      KEY
//      VALUE
//
// The tokens FLOW-SEQUENCE-START, FLOW-SEQUENCE-END, FLOW-MAPPING-START, and
// FLOW-MAPPING-END represent the indicators '[', ']', '***REMOVED***', and '***REMOVED***'
// correspondingly.  FLOW-ENTRY represent the ',' indicator.  Finally the
// indicators '?' and ':', which are used for denoting mapping keys and values,
// are represented by the KEY and VALUE tokens.
//
// The following examples show flow collections:
//
//      1. A flow sequence:
//
//          [item 1, item 2, item 3]
//
//      Tokens:
//
//          STREAM-START(utf-8)
//          FLOW-SEQUENCE-START
//          SCALAR("item 1",plain)
//          FLOW-ENTRY
//          SCALAR("item 2",plain)
//          FLOW-ENTRY
//          SCALAR("item 3",plain)
//          FLOW-SEQUENCE-END
//          STREAM-END
//
//      2. A flow mapping:
//
//          ***REMOVED***
//              a simple key: a value,  # Note that the KEY token is produced.
//              ? a complex key: another value,
//          ***REMOVED***
//
//      Tokens:
//
//          STREAM-START(utf-8)
//          FLOW-MAPPING-START
//          KEY
//          SCALAR("a simple key",plain)
//          VALUE
//          SCALAR("a value",plain)
//          FLOW-ENTRY
//          KEY
//          SCALAR("a complex key",plain)
//          VALUE
//          SCALAR("another value",plain)
//          FLOW-ENTRY
//          FLOW-MAPPING-END
//          STREAM-END
//
// A simple key is a key which is not denoted by the '?' indicator.  Note that
// the Scanner still produce the KEY token whenever it encounters a simple key.
//
// For scanning block collections, the following tokens are used (note that we
// repeat KEY and VALUE here):
//
//      BLOCK-SEQUENCE-START
//      BLOCK-MAPPING-START
//      BLOCK-END
//      BLOCK-ENTRY
//      KEY
//      VALUE
//
// The tokens BLOCK-SEQUENCE-START and BLOCK-MAPPING-START denote indentation
// increase that precedes a block collection (cf. the INDENT token in Python).
// The token BLOCK-END denote indentation decrease that ends a block collection
// (cf. the DEDENT token in Python).  However YAML has some syntax pecularities
// that makes detections of these tokens more complex.
//
// The tokens BLOCK-ENTRY, KEY, and VALUE are used to represent the indicators
// '-', '?', and ':' correspondingly.
//
// The following examples show how the tokens BLOCK-SEQUENCE-START,
// BLOCK-MAPPING-START, and BLOCK-END are emitted by the Scanner:
//
//      1. Block sequences:
//
//          - item 1
//          - item 2
//          -
//            - item 3.1
//            - item 3.2
//          -
//            key 1: value 1
//            key 2: value 2
//
//      Tokens:
//
//          STREAM-START(utf-8)
//          BLOCK-SEQUENCE-START
//          BLOCK-ENTRY
//          SCALAR("item 1",plain)
//          BLOCK-ENTRY
//          SCALAR("item 2",plain)
//          BLOCK-ENTRY
//          BLOCK-SEQUENCE-START
//          BLOCK-ENTRY
//          SCALAR("item 3.1",plain)
//          BLOCK-ENTRY
//          SCALAR("item 3.2",plain)
//          BLOCK-END
//          BLOCK-ENTRY
//          BLOCK-MAPPING-START
//          KEY
//          SCALAR("key 1",plain)
//          VALUE
//          SCALAR("value 1",plain)
//          KEY
//          SCALAR("key 2",plain)
//          VALUE
//          SCALAR("value 2",plain)
//          BLOCK-END
//          BLOCK-END
//          STREAM-END
//
//      2. Block mappings:
//
//          a simple key: a value   # The KEY token is produced here.
//          ? a complex key
//          : another value
//          a mapping:
//            key 1: value 1
//            key 2: value 2
//          a sequence:
//            - item 1
//            - item 2
//
//      Tokens:
//
//          STREAM-START(utf-8)
//          BLOCK-MAPPING-START
//          KEY
//          SCALAR("a simple key",plain)
//          VALUE
//          SCALAR("a value",plain)
//          KEY
//          SCALAR("a complex key",plain)
//          VALUE
//          SCALAR("another value",plain)
//          KEY
//          SCALAR("a mapping",plain)
//          BLOCK-MAPPING-START
//          KEY
//          SCALAR("key 1",plain)
//          VALUE
//          SCALAR("value 1",plain)
//          KEY
//          SCALAR("key 2",plain)
//          VALUE
//          SCALAR("value 2",plain)
//          BLOCK-END
//          KEY
//          SCALAR("a sequence",plain)
//          VALUE
//          BLOCK-SEQUENCE-START
//          BLOCK-ENTRY
//          SCALAR("item 1",plain)
//          BLOCK-ENTRY
//          SCALAR("item 2",plain)
//          BLOCK-END
//          BLOCK-END
//          STREAM-END
//
// YAML does not always require to start a new block collection from a new
// line.  If the current line contains only '-', '?', and ':' indicators, a new
// block collection may start at the current line.  The following examples
// illustrate this case:
//
//      1. Collections in a sequence:
//
//          - - item 1
//            - item 2
//          - key 1: value 1
//            key 2: value 2
//          - ? complex key
//            : complex value
//
//      Tokens:
//
//          STREAM-START(utf-8)
//          BLOCK-SEQUENCE-START
//          BLOCK-ENTRY
//          BLOCK-SEQUENCE-START
//          BLOCK-ENTRY
//          SCALAR("item 1",plain)
//          BLOCK-ENTRY
//          SCALAR("item 2",plain)
//          BLOCK-END
//          BLOCK-ENTRY
//          BLOCK-MAPPING-START
//          KEY
//          SCALAR("key 1",plain)
//          VALUE
//          SCALAR("value 1",plain)
//          KEY
//          SCALAR("key 2",plain)
//          VALUE
//          SCALAR("value 2",plain)
//          BLOCK-END
//          BLOCK-ENTRY
//          BLOCK-MAPPING-START
//          KEY
//          SCALAR("complex key")
//          VALUE
//          SCALAR("complex value")
//          BLOCK-END
//          BLOCK-END
//          STREAM-END
//
//      2. Collections in a mapping:
//
//          ? a sequence
//          : - item 1
//            - item 2
//          ? a mapping
//          : key 1: value 1
//            key 2: value 2
//
//      Tokens:
//
//          STREAM-START(utf-8)
//          BLOCK-MAPPING-START
//          KEY
//          SCALAR("a sequence",plain)
//          VALUE
//          BLOCK-SEQUENCE-START
//          BLOCK-ENTRY
//          SCALAR("item 1",plain)
//          BLOCK-ENTRY
//          SCALAR("item 2",plain)
//          BLOCK-END
//          KEY
//          SCALAR("a mapping",plain)
//          VALUE
//          BLOCK-MAPPING-START
//          KEY
//          SCALAR("key 1",plain)
//          VALUE
//          SCALAR("value 1",plain)
//          KEY
//          SCALAR("key 2",plain)
//          VALUE
//          SCALAR("value 2",plain)
//          BLOCK-END
//          BLOCK-END
//          STREAM-END
//
// YAML also permits non-indented sequences if they are included into a block
// mapping.  In this case, the token BLOCK-SEQUENCE-START is not produced:
//
//      key:
//      - item 1    # BLOCK-SEQUENCE-START is NOT produced here.
//      - item 2
//
// Tokens:
//
//      STREAM-START(utf-8)
//      BLOCK-MAPPING-START
//      KEY
//      SCALAR("key",plain)
//      VALUE
//      BLOCK-ENTRY
//      SCALAR("item 1",plain)
//      BLOCK-ENTRY
//      SCALAR("item 2",plain)
//      BLOCK-END
//

// Ensure that the buffer contains the required number of characters.
// Return true on success, false on failure (reader error or memory error).
func cache(parser *yaml_parser_t, length int) bool ***REMOVED***
	// [Go] This was inlined: !cache(A, B) -> unread < B && !update(A, B)
	return parser.unread >= length || yaml_parser_update_buffer(parser, length)
***REMOVED***

// Advance the buffer pointer.
func skip(parser *yaml_parser_t) ***REMOVED***
	if !is_blank(parser.buffer, parser.buffer_pos) ***REMOVED***
		parser.newlines = 0
	***REMOVED***
	parser.mark.index++
	parser.mark.column++
	parser.unread--
	parser.buffer_pos += width(parser.buffer[parser.buffer_pos])
***REMOVED***

func skip_line(parser *yaml_parser_t) ***REMOVED***
	if is_crlf(parser.buffer, parser.buffer_pos) ***REMOVED***
		parser.mark.index += 2
		parser.mark.column = 0
		parser.mark.line++
		parser.unread -= 2
		parser.buffer_pos += 2
		parser.newlines++
	***REMOVED*** else if is_break(parser.buffer, parser.buffer_pos) ***REMOVED***
		parser.mark.index++
		parser.mark.column = 0
		parser.mark.line++
		parser.unread--
		parser.buffer_pos += width(parser.buffer[parser.buffer_pos])
		parser.newlines++
	***REMOVED***
***REMOVED***

// Copy a character to a string buffer and advance pointers.
func read(parser *yaml_parser_t, s []byte) []byte ***REMOVED***
	if !is_blank(parser.buffer, parser.buffer_pos) ***REMOVED***
		parser.newlines = 0
	***REMOVED***
	w := width(parser.buffer[parser.buffer_pos])
	if w == 0 ***REMOVED***
		panic("invalid character sequence")
	***REMOVED***
	if len(s) == 0 ***REMOVED***
		s = make([]byte, 0, 32)
	***REMOVED***
	if w == 1 && len(s)+w <= cap(s) ***REMOVED***
		s = s[:len(s)+1]
		s[len(s)-1] = parser.buffer[parser.buffer_pos]
		parser.buffer_pos++
	***REMOVED*** else ***REMOVED***
		s = append(s, parser.buffer[parser.buffer_pos:parser.buffer_pos+w]...)
		parser.buffer_pos += w
	***REMOVED***
	parser.mark.index++
	parser.mark.column++
	parser.unread--
	return s
***REMOVED***

// Copy a line break character to a string buffer and advance pointers.
func read_line(parser *yaml_parser_t, s []byte) []byte ***REMOVED***
	buf := parser.buffer
	pos := parser.buffer_pos
	switch ***REMOVED***
	case buf[pos] == '\r' && buf[pos+1] == '\n':
		// CR LF . LF
		s = append(s, '\n')
		parser.buffer_pos += 2
		parser.mark.index++
		parser.unread--
	case buf[pos] == '\r' || buf[pos] == '\n':
		// CR|LF . LF
		s = append(s, '\n')
		parser.buffer_pos += 1
	case buf[pos] == '\xC2' && buf[pos+1] == '\x85':
		// NEL . LF
		s = append(s, '\n')
		parser.buffer_pos += 2
	case buf[pos] == '\xE2' && buf[pos+1] == '\x80' && (buf[pos+2] == '\xA8' || buf[pos+2] == '\xA9'):
		// LS|PS . LS|PS
		s = append(s, buf[parser.buffer_pos:pos+3]...)
		parser.buffer_pos += 3
	default:
		return s
	***REMOVED***
	parser.mark.index++
	parser.mark.column = 0
	parser.mark.line++
	parser.unread--
	parser.newlines++
	return s
***REMOVED***

// Get the next token.
func yaml_parser_scan(parser *yaml_parser_t, token *yaml_token_t) bool ***REMOVED***
	// Erase the token object.
	*token = yaml_token_t***REMOVED******REMOVED*** // [Go] Is this necessary?

	// No tokens after STREAM-END or error.
	if parser.stream_end_produced || parser.error != yaml_NO_ERROR ***REMOVED***
		return true
	***REMOVED***

	// Ensure that the tokens queue contains enough tokens.
	if !parser.token_available ***REMOVED***
		if !yaml_parser_fetch_more_tokens(parser) ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***

	// Fetch the next token from the queue.
	*token = parser.tokens[parser.tokens_head]
	parser.tokens_head++
	parser.tokens_parsed++
	parser.token_available = false

	if token.typ == yaml_STREAM_END_TOKEN ***REMOVED***
		parser.stream_end_produced = true
	***REMOVED***
	return true
***REMOVED***

// Set the scanner error and return false.
func yaml_parser_set_scanner_error(parser *yaml_parser_t, context string, context_mark yaml_mark_t, problem string) bool ***REMOVED***
	parser.error = yaml_SCANNER_ERROR
	parser.context = context
	parser.context_mark = context_mark
	parser.problem = problem
	parser.problem_mark = parser.mark
	return false
***REMOVED***

func yaml_parser_set_scanner_tag_error(parser *yaml_parser_t, directive bool, context_mark yaml_mark_t, problem string) bool ***REMOVED***
	context := "while parsing a tag"
	if directive ***REMOVED***
		context = "while parsing a %TAG directive"
	***REMOVED***
	return yaml_parser_set_scanner_error(parser, context, context_mark, problem)
***REMOVED***

func trace(args ...interface***REMOVED******REMOVED***) func() ***REMOVED***
	pargs := append([]interface***REMOVED******REMOVED******REMOVED***"+++"***REMOVED***, args...)
	fmt.Println(pargs...)
	pargs = append([]interface***REMOVED******REMOVED******REMOVED***"---"***REMOVED***, args...)
	return func() ***REMOVED*** fmt.Println(pargs...) ***REMOVED***
***REMOVED***

// Ensure that the tokens queue contains at least one token which can be
// returned to the Parser.
func yaml_parser_fetch_more_tokens(parser *yaml_parser_t) bool ***REMOVED***
	// While we need more tokens to fetch, do it.
	for ***REMOVED***
		// [Go] The comment parsing logic requires a lookahead of two tokens
		// so that foot comments may be parsed in time of associating them
		// with the tokens that are parsed before them, and also for line
		// comments to be transformed into head comments in some edge cases.
		if parser.tokens_head < len(parser.tokens)-2 ***REMOVED***
			// If a potential simple key is at the head position, we need to fetch
			// the next token to disambiguate it.
			head_tok_idx, ok := parser.simple_keys_by_tok[parser.tokens_parsed]
			if !ok ***REMOVED***
				break
			***REMOVED*** else if valid, ok := yaml_simple_key_is_valid(parser, &parser.simple_keys[head_tok_idx]); !ok ***REMOVED***
				return false
			***REMOVED*** else if !valid ***REMOVED***
				break
			***REMOVED***
		***REMOVED***
		// Fetch the next token.
		if !yaml_parser_fetch_next_token(parser) ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***

	parser.token_available = true
	return true
***REMOVED***

// The dispatcher for token fetchers.
func yaml_parser_fetch_next_token(parser *yaml_parser_t) (ok bool) ***REMOVED***
	// Ensure that the buffer is initialized.
	if parser.unread < 1 && !yaml_parser_update_buffer(parser, 1) ***REMOVED***
		return false
	***REMOVED***

	// Check if we just started scanning.  Fetch STREAM-START then.
	if !parser.stream_start_produced ***REMOVED***
		return yaml_parser_fetch_stream_start(parser)
	***REMOVED***

	scan_mark := parser.mark

	// Eat whitespaces and comments until we reach the next token.
	if !yaml_parser_scan_to_next_token(parser) ***REMOVED***
		return false
	***REMOVED***

	// [Go] While unrolling indents, transform the head comments of prior
	// indentation levels observed after scan_start into foot comments at
	// the respective indexes.

	// Check the indentation level against the current column.
	if !yaml_parser_unroll_indent(parser, parser.mark.column, scan_mark) ***REMOVED***
		return false
	***REMOVED***

	// Ensure that the buffer contains at least 4 characters.  4 is the length
	// of the longest indicators ('--- ' and '... ').
	if parser.unread < 4 && !yaml_parser_update_buffer(parser, 4) ***REMOVED***
		return false
	***REMOVED***

	// Is it the end of the stream?
	if is_z(parser.buffer, parser.buffer_pos) ***REMOVED***
		return yaml_parser_fetch_stream_end(parser)
	***REMOVED***

	// Is it a directive?
	if parser.mark.column == 0 && parser.buffer[parser.buffer_pos] == '%' ***REMOVED***
		return yaml_parser_fetch_directive(parser)
	***REMOVED***

	buf := parser.buffer
	pos := parser.buffer_pos

	// Is it the document start indicator?
	if parser.mark.column == 0 && buf[pos] == '-' && buf[pos+1] == '-' && buf[pos+2] == '-' && is_blankz(buf, pos+3) ***REMOVED***
		return yaml_parser_fetch_document_indicator(parser, yaml_DOCUMENT_START_TOKEN)
	***REMOVED***

	// Is it the document end indicator?
	if parser.mark.column == 0 && buf[pos] == '.' && buf[pos+1] == '.' && buf[pos+2] == '.' && is_blankz(buf, pos+3) ***REMOVED***
		return yaml_parser_fetch_document_indicator(parser, yaml_DOCUMENT_END_TOKEN)
	***REMOVED***

	comment_mark := parser.mark
	if len(parser.tokens) > 0 && (parser.flow_level == 0 && buf[pos] == ':' || parser.flow_level > 0 && buf[pos] == ',') ***REMOVED***
		// Associate any following comments with the prior token.
		comment_mark = parser.tokens[len(parser.tokens)-1].start_mark
	***REMOVED***
	defer func() ***REMOVED***
		if !ok ***REMOVED***
			return
		***REMOVED***
		if !yaml_parser_scan_line_comment(parser, comment_mark) ***REMOVED***
			ok = false
			return
		***REMOVED***
	***REMOVED***()

	// Is it the flow sequence start indicator?
	if buf[pos] == '[' ***REMOVED***
		return yaml_parser_fetch_flow_collection_start(parser, yaml_FLOW_SEQUENCE_START_TOKEN)
	***REMOVED***

	// Is it the flow mapping start indicator?
	if parser.buffer[parser.buffer_pos] == '***REMOVED***' ***REMOVED***
		return yaml_parser_fetch_flow_collection_start(parser, yaml_FLOW_MAPPING_START_TOKEN)
	***REMOVED***

	// Is it the flow sequence end indicator?
	if parser.buffer[parser.buffer_pos] == ']' ***REMOVED***
		return yaml_parser_fetch_flow_collection_end(parser,
			yaml_FLOW_SEQUENCE_END_TOKEN)
	***REMOVED***

	// Is it the flow mapping end indicator?
	if parser.buffer[parser.buffer_pos] == '***REMOVED***' ***REMOVED***
		return yaml_parser_fetch_flow_collection_end(parser,
			yaml_FLOW_MAPPING_END_TOKEN)
	***REMOVED***

	// Is it the flow entry indicator?
	if parser.buffer[parser.buffer_pos] == ',' ***REMOVED***
		return yaml_parser_fetch_flow_entry(parser)
	***REMOVED***

	// Is it the block entry indicator?
	if parser.buffer[parser.buffer_pos] == '-' && is_blankz(parser.buffer, parser.buffer_pos+1) ***REMOVED***
		return yaml_parser_fetch_block_entry(parser)
	***REMOVED***

	// Is it the key indicator?
	if parser.buffer[parser.buffer_pos] == '?' && (parser.flow_level > 0 || is_blankz(parser.buffer, parser.buffer_pos+1)) ***REMOVED***
		return yaml_parser_fetch_key(parser)
	***REMOVED***

	// Is it the value indicator?
	if parser.buffer[parser.buffer_pos] == ':' && (parser.flow_level > 0 || is_blankz(parser.buffer, parser.buffer_pos+1)) ***REMOVED***
		return yaml_parser_fetch_value(parser)
	***REMOVED***

	// Is it an alias?
	if parser.buffer[parser.buffer_pos] == '*' ***REMOVED***
		return yaml_parser_fetch_anchor(parser, yaml_ALIAS_TOKEN)
	***REMOVED***

	// Is it an anchor?
	if parser.buffer[parser.buffer_pos] == '&' ***REMOVED***
		return yaml_parser_fetch_anchor(parser, yaml_ANCHOR_TOKEN)
	***REMOVED***

	// Is it a tag?
	if parser.buffer[parser.buffer_pos] == '!' ***REMOVED***
		return yaml_parser_fetch_tag(parser)
	***REMOVED***

	// Is it a literal scalar?
	if parser.buffer[parser.buffer_pos] == '|' && parser.flow_level == 0 ***REMOVED***
		return yaml_parser_fetch_block_scalar(parser, true)
	***REMOVED***

	// Is it a folded scalar?
	if parser.buffer[parser.buffer_pos] == '>' && parser.flow_level == 0 ***REMOVED***
		return yaml_parser_fetch_block_scalar(parser, false)
	***REMOVED***

	// Is it a single-quoted scalar?
	if parser.buffer[parser.buffer_pos] == '\'' ***REMOVED***
		return yaml_parser_fetch_flow_scalar(parser, true)
	***REMOVED***

	// Is it a double-quoted scalar?
	if parser.buffer[parser.buffer_pos] == '"' ***REMOVED***
		return yaml_parser_fetch_flow_scalar(parser, false)
	***REMOVED***

	// Is it a plain scalar?
	//
	// A plain scalar may start with any non-blank characters except
	//
	//      '-', '?', ':', ',', '[', ']', '***REMOVED***', '***REMOVED***',
	//      '#', '&', '*', '!', '|', '>', '\'', '\"',
	//      '%', '@', '`'.
	//
	// In the block context (and, for the '-' indicator, in the flow context
	// too), it may also start with the characters
	//
	//      '-', '?', ':'
	//
	// if it is followed by a non-space character.
	//
	// The last rule is more restrictive than the specification requires.
	// [Go] TODO Make this logic more reasonable.
	//switch parser.buffer[parser.buffer_pos] ***REMOVED***
	//case '-', '?', ':', ',', '?', '-', ',', ':', ']', '[', '***REMOVED***', '***REMOVED***', '&', '#', '!', '*', '>', '|', '"', '\'', '@', '%', '-', '`':
	//***REMOVED***
	if !(is_blankz(parser.buffer, parser.buffer_pos) || parser.buffer[parser.buffer_pos] == '-' ||
		parser.buffer[parser.buffer_pos] == '?' || parser.buffer[parser.buffer_pos] == ':' ||
		parser.buffer[parser.buffer_pos] == ',' || parser.buffer[parser.buffer_pos] == '[' ||
		parser.buffer[parser.buffer_pos] == ']' || parser.buffer[parser.buffer_pos] == '***REMOVED***' ||
		parser.buffer[parser.buffer_pos] == '***REMOVED***' || parser.buffer[parser.buffer_pos] == '#' ||
		parser.buffer[parser.buffer_pos] == '&' || parser.buffer[parser.buffer_pos] == '*' ||
		parser.buffer[parser.buffer_pos] == '!' || parser.buffer[parser.buffer_pos] == '|' ||
		parser.buffer[parser.buffer_pos] == '>' || parser.buffer[parser.buffer_pos] == '\'' ||
		parser.buffer[parser.buffer_pos] == '"' || parser.buffer[parser.buffer_pos] == '%' ||
		parser.buffer[parser.buffer_pos] == '@' || parser.buffer[parser.buffer_pos] == '`') ||
		(parser.buffer[parser.buffer_pos] == '-' && !is_blank(parser.buffer, parser.buffer_pos+1)) ||
		(parser.flow_level == 0 &&
			(parser.buffer[parser.buffer_pos] == '?' || parser.buffer[parser.buffer_pos] == ':') &&
			!is_blankz(parser.buffer, parser.buffer_pos+1)) ***REMOVED***
		return yaml_parser_fetch_plain_scalar(parser)
	***REMOVED***

	// If we don't determine the token type so far, it is an error.
	return yaml_parser_set_scanner_error(parser,
		"while scanning for the next token", parser.mark,
		"found character that cannot start any token")
***REMOVED***

func yaml_simple_key_is_valid(parser *yaml_parser_t, simple_key *yaml_simple_key_t) (valid, ok bool) ***REMOVED***
	if !simple_key.possible ***REMOVED***
		return false, true
	***REMOVED***

	// The 1.2 specification says:
	//
	//     "If the ? indicator is omitted, parsing needs to see past the
	//     implicit key to recognize it as such. To limit the amount of
	//     lookahead required, the “:” indicator must appear at most 1024
	//     Unicode characters beyond the start of the key. In addition, the key
	//     is restricted to a single line."
	//
	if simple_key.mark.line < parser.mark.line || simple_key.mark.index+1024 < parser.mark.index ***REMOVED***
		// Check if the potential simple key to be removed is required.
		if simple_key.required ***REMOVED***
			return false, yaml_parser_set_scanner_error(parser,
				"while scanning a simple key", simple_key.mark,
				"could not find expected ':'")
		***REMOVED***
		simple_key.possible = false
		return false, true
	***REMOVED***
	return true, true
***REMOVED***

// Check if a simple key may start at the current position and add it if
// needed.
func yaml_parser_save_simple_key(parser *yaml_parser_t) bool ***REMOVED***
	// A simple key is required at the current position if the scanner is in
	// the block context and the current column coincides with the indentation
	// level.

	required := parser.flow_level == 0 && parser.indent == parser.mark.column

	//
	// If the current position may start a simple key, save it.
	//
	if parser.simple_key_allowed ***REMOVED***
		simple_key := yaml_simple_key_t***REMOVED***
			possible:     true,
			required:     required,
			token_number: parser.tokens_parsed + (len(parser.tokens) - parser.tokens_head),
			mark:         parser.mark,
		***REMOVED***

		if !yaml_parser_remove_simple_key(parser) ***REMOVED***
			return false
		***REMOVED***
		parser.simple_keys[len(parser.simple_keys)-1] = simple_key
		parser.simple_keys_by_tok[simple_key.token_number] = len(parser.simple_keys) - 1
	***REMOVED***
	return true
***REMOVED***

// Remove a potential simple key at the current flow level.
func yaml_parser_remove_simple_key(parser *yaml_parser_t) bool ***REMOVED***
	i := len(parser.simple_keys) - 1
	if parser.simple_keys[i].possible ***REMOVED***
		// If the key is required, it is an error.
		if parser.simple_keys[i].required ***REMOVED***
			return yaml_parser_set_scanner_error(parser,
				"while scanning a simple key", parser.simple_keys[i].mark,
				"could not find expected ':'")
		***REMOVED***
		// Remove the key from the stack.
		parser.simple_keys[i].possible = false
		delete(parser.simple_keys_by_tok, parser.simple_keys[i].token_number)
	***REMOVED***
	return true
***REMOVED***

// max_flow_level limits the flow_level
const max_flow_level = 10000

// Increase the flow level and resize the simple key list if needed.
func yaml_parser_increase_flow_level(parser *yaml_parser_t) bool ***REMOVED***
	// Reset the simple key on the next level.
	parser.simple_keys = append(parser.simple_keys, yaml_simple_key_t***REMOVED***
		possible:     false,
		required:     false,
		token_number: parser.tokens_parsed + (len(parser.tokens) - parser.tokens_head),
		mark:         parser.mark,
	***REMOVED***)

	// Increase the flow level.
	parser.flow_level++
	if parser.flow_level > max_flow_level ***REMOVED***
		return yaml_parser_set_scanner_error(parser,
			"while increasing flow level", parser.simple_keys[len(parser.simple_keys)-1].mark,
			fmt.Sprintf("exceeded max depth of %d", max_flow_level))
	***REMOVED***
	return true
***REMOVED***

// Decrease the flow level.
func yaml_parser_decrease_flow_level(parser *yaml_parser_t) bool ***REMOVED***
	if parser.flow_level > 0 ***REMOVED***
		parser.flow_level--
		last := len(parser.simple_keys) - 1
		delete(parser.simple_keys_by_tok, parser.simple_keys[last].token_number)
		parser.simple_keys = parser.simple_keys[:last]
	***REMOVED***
	return true
***REMOVED***

// max_indents limits the indents stack size
const max_indents = 10000

// Push the current indentation level to the stack and set the new level
// the current column is greater than the indentation level.  In this case,
// append or insert the specified token into the token queue.
func yaml_parser_roll_indent(parser *yaml_parser_t, column, number int, typ yaml_token_type_t, mark yaml_mark_t) bool ***REMOVED***
	// In the flow context, do nothing.
	if parser.flow_level > 0 ***REMOVED***
		return true
	***REMOVED***

	if parser.indent < column ***REMOVED***
		// Push the current indentation level to the stack and set the new
		// indentation level.
		parser.indents = append(parser.indents, parser.indent)
		parser.indent = column
		if len(parser.indents) > max_indents ***REMOVED***
			return yaml_parser_set_scanner_error(parser,
				"while increasing indent level", parser.simple_keys[len(parser.simple_keys)-1].mark,
				fmt.Sprintf("exceeded max depth of %d", max_indents))
		***REMOVED***

		// Create a token and insert it into the queue.
		token := yaml_token_t***REMOVED***
			typ:        typ,
			start_mark: mark,
			end_mark:   mark,
		***REMOVED***
		if number > -1 ***REMOVED***
			number -= parser.tokens_parsed
		***REMOVED***
		yaml_insert_token(parser, number, &token)
	***REMOVED***
	return true
***REMOVED***

// Pop indentation levels from the indents stack until the current level
// becomes less or equal to the column.  For each indentation level, append
// the BLOCK-END token.
func yaml_parser_unroll_indent(parser *yaml_parser_t, column int, scan_mark yaml_mark_t) bool ***REMOVED***
	// In the flow context, do nothing.
	if parser.flow_level > 0 ***REMOVED***
		return true
	***REMOVED***

	block_mark := scan_mark
	block_mark.index--

	// Loop through the indentation levels in the stack.
	for parser.indent > column ***REMOVED***

		// [Go] Reposition the end token before potential following
		//      foot comments of parent blocks. For that, search
		//      backwards for recent comments that were at the same
		//      indent as the block that is ending now.
		stop_index := block_mark.index
		for i := len(parser.comments) - 1; i >= 0; i-- ***REMOVED***
			comment := &parser.comments[i]

			if comment.end_mark.index < stop_index ***REMOVED***
				// Don't go back beyond the start of the comment/whitespace scan, unless column < 0.
				// If requested indent column is < 0, then the document is over and everything else
				// is a foot anyway.
				break
			***REMOVED***
			if comment.start_mark.column == parser.indent+1 ***REMOVED***
				// This is a good match. But maybe there's a former comment
				// at that same indent level, so keep searching.
				block_mark = comment.start_mark
			***REMOVED***

			// While the end of the former comment matches with
			// the start of the following one, we know there's
			// nothing in between and scanning is still safe.
			stop_index = comment.scan_mark.index
		***REMOVED***

		// Create a token and append it to the queue.
		token := yaml_token_t***REMOVED***
			typ:        yaml_BLOCK_END_TOKEN,
			start_mark: block_mark,
			end_mark:   block_mark,
		***REMOVED***
		yaml_insert_token(parser, -1, &token)

		// Pop the indentation level.
		parser.indent = parser.indents[len(parser.indents)-1]
		parser.indents = parser.indents[:len(parser.indents)-1]
	***REMOVED***
	return true
***REMOVED***

// Initialize the scanner and produce the STREAM-START token.
func yaml_parser_fetch_stream_start(parser *yaml_parser_t) bool ***REMOVED***

	// Set the initial indentation.
	parser.indent = -1

	// Initialize the simple key stack.
	parser.simple_keys = append(parser.simple_keys, yaml_simple_key_t***REMOVED******REMOVED***)

	parser.simple_keys_by_tok = make(map[int]int)

	// A simple key is allowed at the beginning of the stream.
	parser.simple_key_allowed = true

	// We have started.
	parser.stream_start_produced = true

	// Create the STREAM-START token and append it to the queue.
	token := yaml_token_t***REMOVED***
		typ:        yaml_STREAM_START_TOKEN,
		start_mark: parser.mark,
		end_mark:   parser.mark,
		encoding:   parser.encoding,
	***REMOVED***
	yaml_insert_token(parser, -1, &token)
	return true
***REMOVED***

// Produce the STREAM-END token and shut down the scanner.
func yaml_parser_fetch_stream_end(parser *yaml_parser_t) bool ***REMOVED***

	// Force new line.
	if parser.mark.column != 0 ***REMOVED***
		parser.mark.column = 0
		parser.mark.line++
	***REMOVED***

	// Reset the indentation level.
	if !yaml_parser_unroll_indent(parser, -1, parser.mark) ***REMOVED***
		return false
	***REMOVED***

	// Reset simple keys.
	if !yaml_parser_remove_simple_key(parser) ***REMOVED***
		return false
	***REMOVED***

	parser.simple_key_allowed = false

	// Create the STREAM-END token and append it to the queue.
	token := yaml_token_t***REMOVED***
		typ:        yaml_STREAM_END_TOKEN,
		start_mark: parser.mark,
		end_mark:   parser.mark,
	***REMOVED***
	yaml_insert_token(parser, -1, &token)
	return true
***REMOVED***

// Produce a VERSION-DIRECTIVE or TAG-DIRECTIVE token.
func yaml_parser_fetch_directive(parser *yaml_parser_t) bool ***REMOVED***
	// Reset the indentation level.
	if !yaml_parser_unroll_indent(parser, -1, parser.mark) ***REMOVED***
		return false
	***REMOVED***

	// Reset simple keys.
	if !yaml_parser_remove_simple_key(parser) ***REMOVED***
		return false
	***REMOVED***

	parser.simple_key_allowed = false

	// Create the YAML-DIRECTIVE or TAG-DIRECTIVE token.
	token := yaml_token_t***REMOVED******REMOVED***
	if !yaml_parser_scan_directive(parser, &token) ***REMOVED***
		return false
	***REMOVED***
	// Append the token to the queue.
	yaml_insert_token(parser, -1, &token)
	return true
***REMOVED***

// Produce the DOCUMENT-START or DOCUMENT-END token.
func yaml_parser_fetch_document_indicator(parser *yaml_parser_t, typ yaml_token_type_t) bool ***REMOVED***
	// Reset the indentation level.
	if !yaml_parser_unroll_indent(parser, -1, parser.mark) ***REMOVED***
		return false
	***REMOVED***

	// Reset simple keys.
	if !yaml_parser_remove_simple_key(parser) ***REMOVED***
		return false
	***REMOVED***

	parser.simple_key_allowed = false

	// Consume the token.
	start_mark := parser.mark

	skip(parser)
	skip(parser)
	skip(parser)

	end_mark := parser.mark

	// Create the DOCUMENT-START or DOCUMENT-END token.
	token := yaml_token_t***REMOVED***
		typ:        typ,
		start_mark: start_mark,
		end_mark:   end_mark,
	***REMOVED***
	// Append the token to the queue.
	yaml_insert_token(parser, -1, &token)
	return true
***REMOVED***

// Produce the FLOW-SEQUENCE-START or FLOW-MAPPING-START token.
func yaml_parser_fetch_flow_collection_start(parser *yaml_parser_t, typ yaml_token_type_t) bool ***REMOVED***

	// The indicators '[' and '***REMOVED***' may start a simple key.
	if !yaml_parser_save_simple_key(parser) ***REMOVED***
		return false
	***REMOVED***

	// Increase the flow level.
	if !yaml_parser_increase_flow_level(parser) ***REMOVED***
		return false
	***REMOVED***

	// A simple key may follow the indicators '[' and '***REMOVED***'.
	parser.simple_key_allowed = true

	// Consume the token.
	start_mark := parser.mark
	skip(parser)
	end_mark := parser.mark

	// Create the FLOW-SEQUENCE-START of FLOW-MAPPING-START token.
	token := yaml_token_t***REMOVED***
		typ:        typ,
		start_mark: start_mark,
		end_mark:   end_mark,
	***REMOVED***
	// Append the token to the queue.
	yaml_insert_token(parser, -1, &token)
	return true
***REMOVED***

// Produce the FLOW-SEQUENCE-END or FLOW-MAPPING-END token.
func yaml_parser_fetch_flow_collection_end(parser *yaml_parser_t, typ yaml_token_type_t) bool ***REMOVED***
	// Reset any potential simple key on the current flow level.
	if !yaml_parser_remove_simple_key(parser) ***REMOVED***
		return false
	***REMOVED***

	// Decrease the flow level.
	if !yaml_parser_decrease_flow_level(parser) ***REMOVED***
		return false
	***REMOVED***

	// No simple keys after the indicators ']' and '***REMOVED***'.
	parser.simple_key_allowed = false

	// Consume the token.

	start_mark := parser.mark
	skip(parser)
	end_mark := parser.mark

	// Create the FLOW-SEQUENCE-END of FLOW-MAPPING-END token.
	token := yaml_token_t***REMOVED***
		typ:        typ,
		start_mark: start_mark,
		end_mark:   end_mark,
	***REMOVED***
	// Append the token to the queue.
	yaml_insert_token(parser, -1, &token)
	return true
***REMOVED***

// Produce the FLOW-ENTRY token.
func yaml_parser_fetch_flow_entry(parser *yaml_parser_t) bool ***REMOVED***
	// Reset any potential simple keys on the current flow level.
	if !yaml_parser_remove_simple_key(parser) ***REMOVED***
		return false
	***REMOVED***

	// Simple keys are allowed after ','.
	parser.simple_key_allowed = true

	// Consume the token.
	start_mark := parser.mark
	skip(parser)
	end_mark := parser.mark

	// Create the FLOW-ENTRY token and append it to the queue.
	token := yaml_token_t***REMOVED***
		typ:        yaml_FLOW_ENTRY_TOKEN,
		start_mark: start_mark,
		end_mark:   end_mark,
	***REMOVED***
	yaml_insert_token(parser, -1, &token)
	return true
***REMOVED***

// Produce the BLOCK-ENTRY token.
func yaml_parser_fetch_block_entry(parser *yaml_parser_t) bool ***REMOVED***
	// Check if the scanner is in the block context.
	if parser.flow_level == 0 ***REMOVED***
		// Check if we are allowed to start a new entry.
		if !parser.simple_key_allowed ***REMOVED***
			return yaml_parser_set_scanner_error(parser, "", parser.mark,
				"block sequence entries are not allowed in this context")
		***REMOVED***
		// Add the BLOCK-SEQUENCE-START token if needed.
		if !yaml_parser_roll_indent(parser, parser.mark.column, -1, yaml_BLOCK_SEQUENCE_START_TOKEN, parser.mark) ***REMOVED***
			return false
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		// It is an error for the '-' indicator to occur in the flow context,
		// but we let the Parser detect and report about it because the Parser
		// is able to point to the context.
	***REMOVED***

	// Reset any potential simple keys on the current flow level.
	if !yaml_parser_remove_simple_key(parser) ***REMOVED***
		return false
	***REMOVED***

	// Simple keys are allowed after '-'.
	parser.simple_key_allowed = true

	// Consume the token.
	start_mark := parser.mark
	skip(parser)
	end_mark := parser.mark

	// Create the BLOCK-ENTRY token and append it to the queue.
	token := yaml_token_t***REMOVED***
		typ:        yaml_BLOCK_ENTRY_TOKEN,
		start_mark: start_mark,
		end_mark:   end_mark,
	***REMOVED***
	yaml_insert_token(parser, -1, &token)
	return true
***REMOVED***

// Produce the KEY token.
func yaml_parser_fetch_key(parser *yaml_parser_t) bool ***REMOVED***

	// In the block context, additional checks are required.
	if parser.flow_level == 0 ***REMOVED***
		// Check if we are allowed to start a new key (not nessesary simple).
		if !parser.simple_key_allowed ***REMOVED***
			return yaml_parser_set_scanner_error(parser, "", parser.mark,
				"mapping keys are not allowed in this context")
		***REMOVED***
		// Add the BLOCK-MAPPING-START token if needed.
		if !yaml_parser_roll_indent(parser, parser.mark.column, -1, yaml_BLOCK_MAPPING_START_TOKEN, parser.mark) ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***

	// Reset any potential simple keys on the current flow level.
	if !yaml_parser_remove_simple_key(parser) ***REMOVED***
		return false
	***REMOVED***

	// Simple keys are allowed after '?' in the block context.
	parser.simple_key_allowed = parser.flow_level == 0

	// Consume the token.
	start_mark := parser.mark
	skip(parser)
	end_mark := parser.mark

	// Create the KEY token and append it to the queue.
	token := yaml_token_t***REMOVED***
		typ:        yaml_KEY_TOKEN,
		start_mark: start_mark,
		end_mark:   end_mark,
	***REMOVED***
	yaml_insert_token(parser, -1, &token)
	return true
***REMOVED***

// Produce the VALUE token.
func yaml_parser_fetch_value(parser *yaml_parser_t) bool ***REMOVED***

	simple_key := &parser.simple_keys[len(parser.simple_keys)-1]

	// Have we found a simple key?
	if valid, ok := yaml_simple_key_is_valid(parser, simple_key); !ok ***REMOVED***
		return false

	***REMOVED*** else if valid ***REMOVED***

		// Create the KEY token and insert it into the queue.
		token := yaml_token_t***REMOVED***
			typ:        yaml_KEY_TOKEN,
			start_mark: simple_key.mark,
			end_mark:   simple_key.mark,
		***REMOVED***
		yaml_insert_token(parser, simple_key.token_number-parser.tokens_parsed, &token)

		// In the block context, we may need to add the BLOCK-MAPPING-START token.
		if !yaml_parser_roll_indent(parser, simple_key.mark.column,
			simple_key.token_number,
			yaml_BLOCK_MAPPING_START_TOKEN, simple_key.mark) ***REMOVED***
			return false
		***REMOVED***

		// Remove the simple key.
		simple_key.possible = false
		delete(parser.simple_keys_by_tok, simple_key.token_number)

		// A simple key cannot follow another simple key.
		parser.simple_key_allowed = false

	***REMOVED*** else ***REMOVED***
		// The ':' indicator follows a complex key.

		// In the block context, extra checks are required.
		if parser.flow_level == 0 ***REMOVED***

			// Check if we are allowed to start a complex value.
			if !parser.simple_key_allowed ***REMOVED***
				return yaml_parser_set_scanner_error(parser, "", parser.mark,
					"mapping values are not allowed in this context")
			***REMOVED***

			// Add the BLOCK-MAPPING-START token if needed.
			if !yaml_parser_roll_indent(parser, parser.mark.column, -1, yaml_BLOCK_MAPPING_START_TOKEN, parser.mark) ***REMOVED***
				return false
			***REMOVED***
		***REMOVED***

		// Simple keys after ':' are allowed in the block context.
		parser.simple_key_allowed = parser.flow_level == 0
	***REMOVED***

	// Consume the token.
	start_mark := parser.mark
	skip(parser)
	end_mark := parser.mark

	// Create the VALUE token and append it to the queue.
	token := yaml_token_t***REMOVED***
		typ:        yaml_VALUE_TOKEN,
		start_mark: start_mark,
		end_mark:   end_mark,
	***REMOVED***
	yaml_insert_token(parser, -1, &token)
	return true
***REMOVED***

// Produce the ALIAS or ANCHOR token.
func yaml_parser_fetch_anchor(parser *yaml_parser_t, typ yaml_token_type_t) bool ***REMOVED***
	// An anchor or an alias could be a simple key.
	if !yaml_parser_save_simple_key(parser) ***REMOVED***
		return false
	***REMOVED***

	// A simple key cannot follow an anchor or an alias.
	parser.simple_key_allowed = false

	// Create the ALIAS or ANCHOR token and append it to the queue.
	var token yaml_token_t
	if !yaml_parser_scan_anchor(parser, &token, typ) ***REMOVED***
		return false
	***REMOVED***
	yaml_insert_token(parser, -1, &token)
	return true
***REMOVED***

// Produce the TAG token.
func yaml_parser_fetch_tag(parser *yaml_parser_t) bool ***REMOVED***
	// A tag could be a simple key.
	if !yaml_parser_save_simple_key(parser) ***REMOVED***
		return false
	***REMOVED***

	// A simple key cannot follow a tag.
	parser.simple_key_allowed = false

	// Create the TAG token and append it to the queue.
	var token yaml_token_t
	if !yaml_parser_scan_tag(parser, &token) ***REMOVED***
		return false
	***REMOVED***
	yaml_insert_token(parser, -1, &token)
	return true
***REMOVED***

// Produce the SCALAR(...,literal) or SCALAR(...,folded) tokens.
func yaml_parser_fetch_block_scalar(parser *yaml_parser_t, literal bool) bool ***REMOVED***
	// Remove any potential simple keys.
	if !yaml_parser_remove_simple_key(parser) ***REMOVED***
		return false
	***REMOVED***

	// A simple key may follow a block scalar.
	parser.simple_key_allowed = true

	// Create the SCALAR token and append it to the queue.
	var token yaml_token_t
	if !yaml_parser_scan_block_scalar(parser, &token, literal) ***REMOVED***
		return false
	***REMOVED***
	yaml_insert_token(parser, -1, &token)
	return true
***REMOVED***

// Produce the SCALAR(...,single-quoted) or SCALAR(...,double-quoted) tokens.
func yaml_parser_fetch_flow_scalar(parser *yaml_parser_t, single bool) bool ***REMOVED***
	// A plain scalar could be a simple key.
	if !yaml_parser_save_simple_key(parser) ***REMOVED***
		return false
	***REMOVED***

	// A simple key cannot follow a flow scalar.
	parser.simple_key_allowed = false

	// Create the SCALAR token and append it to the queue.
	var token yaml_token_t
	if !yaml_parser_scan_flow_scalar(parser, &token, single) ***REMOVED***
		return false
	***REMOVED***
	yaml_insert_token(parser, -1, &token)
	return true
***REMOVED***

// Produce the SCALAR(...,plain) token.
func yaml_parser_fetch_plain_scalar(parser *yaml_parser_t) bool ***REMOVED***
	// A plain scalar could be a simple key.
	if !yaml_parser_save_simple_key(parser) ***REMOVED***
		return false
	***REMOVED***

	// A simple key cannot follow a flow scalar.
	parser.simple_key_allowed = false

	// Create the SCALAR token and append it to the queue.
	var token yaml_token_t
	if !yaml_parser_scan_plain_scalar(parser, &token) ***REMOVED***
		return false
	***REMOVED***
	yaml_insert_token(parser, -1, &token)
	return true
***REMOVED***

// Eat whitespaces and comments until the next token is found.
func yaml_parser_scan_to_next_token(parser *yaml_parser_t) bool ***REMOVED***

	scan_mark := parser.mark

	// Until the next token is not found.
	for ***REMOVED***
		// Allow the BOM mark to start a line.
		if parser.unread < 1 && !yaml_parser_update_buffer(parser, 1) ***REMOVED***
			return false
		***REMOVED***
		if parser.mark.column == 0 && is_bom(parser.buffer, parser.buffer_pos) ***REMOVED***
			skip(parser)
		***REMOVED***

		// Eat whitespaces.
		// Tabs are allowed:
		//  - in the flow context
		//  - in the block context, but not at the beginning of the line or
		//  after '-', '?', or ':' (complex value).
		if parser.unread < 1 && !yaml_parser_update_buffer(parser, 1) ***REMOVED***
			return false
		***REMOVED***

		for parser.buffer[parser.buffer_pos] == ' ' || ((parser.flow_level > 0 || !parser.simple_key_allowed) && parser.buffer[parser.buffer_pos] == '\t') ***REMOVED***
			skip(parser)
			if parser.unread < 1 && !yaml_parser_update_buffer(parser, 1) ***REMOVED***
				return false
			***REMOVED***
		***REMOVED***

		// Check if we just had a line comment under a sequence entry that
		// looks more like a header to the following content. Similar to this:
		//
		// - # The comment
		//   - Some data
		//
		// If so, transform the line comment to a head comment and reposition.
		if len(parser.comments) > 0 && len(parser.tokens) > 1 ***REMOVED***
			tokenA := parser.tokens[len(parser.tokens)-2]
			tokenB := parser.tokens[len(parser.tokens)-1]
			comment := &parser.comments[len(parser.comments)-1]
			if tokenA.typ == yaml_BLOCK_SEQUENCE_START_TOKEN && tokenB.typ == yaml_BLOCK_ENTRY_TOKEN && len(comment.line) > 0 && !is_break(parser.buffer, parser.buffer_pos) ***REMOVED***
				// If it was in the prior line, reposition so it becomes a
				// header of the follow up token. Otherwise, keep it in place
				// so it becomes a header of the former.
				comment.head = comment.line
				comment.line = nil
				if comment.start_mark.line == parser.mark.line-1 ***REMOVED***
					comment.token_mark = parser.mark
				***REMOVED***
			***REMOVED***
		***REMOVED***

		// Eat a comment until a line break.
		if parser.buffer[parser.buffer_pos] == '#' ***REMOVED***
			if !yaml_parser_scan_comments(parser, scan_mark) ***REMOVED***
				return false
			***REMOVED***
		***REMOVED***

		// If it is a line break, eat it.
		if is_break(parser.buffer, parser.buffer_pos) ***REMOVED***
			if parser.unread < 2 && !yaml_parser_update_buffer(parser, 2) ***REMOVED***
				return false
			***REMOVED***
			skip_line(parser)

			// In the block context, a new line may start a simple key.
			if parser.flow_level == 0 ***REMOVED***
				parser.simple_key_allowed = true
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			break // We have found a token.
		***REMOVED***
	***REMOVED***

	return true
***REMOVED***

// Scan a YAML-DIRECTIVE or TAG-DIRECTIVE token.
//
// Scope:
//      %YAML    1.1    # a comment \n
//      ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
//      %TAG    !yaml!  tag:yaml.org,2002:  \n
//      ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
//
func yaml_parser_scan_directive(parser *yaml_parser_t, token *yaml_token_t) bool ***REMOVED***
	// Eat '%'.
	start_mark := parser.mark
	skip(parser)

	// Scan the directive name.
	var name []byte
	if !yaml_parser_scan_directive_name(parser, start_mark, &name) ***REMOVED***
		return false
	***REMOVED***

	// Is it a YAML directive?
	if bytes.Equal(name, []byte("YAML")) ***REMOVED***
		// Scan the VERSION directive value.
		var major, minor int8
		if !yaml_parser_scan_version_directive_value(parser, start_mark, &major, &minor) ***REMOVED***
			return false
		***REMOVED***
		end_mark := parser.mark

		// Create a VERSION-DIRECTIVE token.
		*token = yaml_token_t***REMOVED***
			typ:        yaml_VERSION_DIRECTIVE_TOKEN,
			start_mark: start_mark,
			end_mark:   end_mark,
			major:      major,
			minor:      minor,
		***REMOVED***

		// Is it a TAG directive?
	***REMOVED*** else if bytes.Equal(name, []byte("TAG")) ***REMOVED***
		// Scan the TAG directive value.
		var handle, prefix []byte
		if !yaml_parser_scan_tag_directive_value(parser, start_mark, &handle, &prefix) ***REMOVED***
			return false
		***REMOVED***
		end_mark := parser.mark

		// Create a TAG-DIRECTIVE token.
		*token = yaml_token_t***REMOVED***
			typ:        yaml_TAG_DIRECTIVE_TOKEN,
			start_mark: start_mark,
			end_mark:   end_mark,
			value:      handle,
			prefix:     prefix,
		***REMOVED***

		// Unknown directive.
	***REMOVED*** else ***REMOVED***
		yaml_parser_set_scanner_error(parser, "while scanning a directive",
			start_mark, "found unknown directive name")
		return false
	***REMOVED***

	// Eat the rest of the line including any comments.
	if parser.unread < 1 && !yaml_parser_update_buffer(parser, 1) ***REMOVED***
		return false
	***REMOVED***

	for is_blank(parser.buffer, parser.buffer_pos) ***REMOVED***
		skip(parser)
		if parser.unread < 1 && !yaml_parser_update_buffer(parser, 1) ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***

	if parser.buffer[parser.buffer_pos] == '#' ***REMOVED***
		// [Go] Discard this inline comment for the time being.
		//if !yaml_parser_scan_line_comment(parser, start_mark) ***REMOVED***
		//	return false
		//***REMOVED***
		for !is_breakz(parser.buffer, parser.buffer_pos) ***REMOVED***
			skip(parser)
			if parser.unread < 1 && !yaml_parser_update_buffer(parser, 1) ***REMOVED***
				return false
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// Check if we are at the end of the line.
	if !is_breakz(parser.buffer, parser.buffer_pos) ***REMOVED***
		yaml_parser_set_scanner_error(parser, "while scanning a directive",
			start_mark, "did not find expected comment or line break")
		return false
	***REMOVED***

	// Eat a line break.
	if is_break(parser.buffer, parser.buffer_pos) ***REMOVED***
		if parser.unread < 2 && !yaml_parser_update_buffer(parser, 2) ***REMOVED***
			return false
		***REMOVED***
		skip_line(parser)
	***REMOVED***

	return true
***REMOVED***

// Scan the directive name.
//
// Scope:
//      %YAML   1.1     # a comment \n
//       ^^^^
//      %TAG    !yaml!  tag:yaml.org,2002:  \n
//       ^^^
//
func yaml_parser_scan_directive_name(parser *yaml_parser_t, start_mark yaml_mark_t, name *[]byte) bool ***REMOVED***
	// Consume the directive name.
	if parser.unread < 1 && !yaml_parser_update_buffer(parser, 1) ***REMOVED***
		return false
	***REMOVED***

	var s []byte
	for is_alpha(parser.buffer, parser.buffer_pos) ***REMOVED***
		s = read(parser, s)
		if parser.unread < 1 && !yaml_parser_update_buffer(parser, 1) ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***

	// Check if the name is empty.
	if len(s) == 0 ***REMOVED***
		yaml_parser_set_scanner_error(parser, "while scanning a directive",
			start_mark, "could not find expected directive name")
		return false
	***REMOVED***

	// Check for an blank character after the name.
	if !is_blankz(parser.buffer, parser.buffer_pos) ***REMOVED***
		yaml_parser_set_scanner_error(parser, "while scanning a directive",
			start_mark, "found unexpected non-alphabetical character")
		return false
	***REMOVED***
	*name = s
	return true
***REMOVED***

// Scan the value of VERSION-DIRECTIVE.
//
// Scope:
//      %YAML   1.1     # a comment \n
//           ^^^^^^
func yaml_parser_scan_version_directive_value(parser *yaml_parser_t, start_mark yaml_mark_t, major, minor *int8) bool ***REMOVED***
	// Eat whitespaces.
	if parser.unread < 1 && !yaml_parser_update_buffer(parser, 1) ***REMOVED***
		return false
	***REMOVED***
	for is_blank(parser.buffer, parser.buffer_pos) ***REMOVED***
		skip(parser)
		if parser.unread < 1 && !yaml_parser_update_buffer(parser, 1) ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***

	// Consume the major version number.
	if !yaml_parser_scan_version_directive_number(parser, start_mark, major) ***REMOVED***
		return false
	***REMOVED***

	// Eat '.'.
	if parser.buffer[parser.buffer_pos] != '.' ***REMOVED***
		return yaml_parser_set_scanner_error(parser, "while scanning a %YAML directive",
			start_mark, "did not find expected digit or '.' character")
	***REMOVED***

	skip(parser)

	// Consume the minor version number.
	if !yaml_parser_scan_version_directive_number(parser, start_mark, minor) ***REMOVED***
		return false
	***REMOVED***
	return true
***REMOVED***

const max_number_length = 2

// Scan the version number of VERSION-DIRECTIVE.
//
// Scope:
//      %YAML   1.1     # a comment \n
//              ^
//      %YAML   1.1     # a comment \n
//                ^
func yaml_parser_scan_version_directive_number(parser *yaml_parser_t, start_mark yaml_mark_t, number *int8) bool ***REMOVED***

	// Repeat while the next character is digit.
	if parser.unread < 1 && !yaml_parser_update_buffer(parser, 1) ***REMOVED***
		return false
	***REMOVED***
	var value, length int8
	for is_digit(parser.buffer, parser.buffer_pos) ***REMOVED***
		// Check if the number is too long.
		length++
		if length > max_number_length ***REMOVED***
			return yaml_parser_set_scanner_error(parser, "while scanning a %YAML directive",
				start_mark, "found extremely long version number")
		***REMOVED***
		value = value*10 + int8(as_digit(parser.buffer, parser.buffer_pos))
		skip(parser)
		if parser.unread < 1 && !yaml_parser_update_buffer(parser, 1) ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***

	// Check if the number was present.
	if length == 0 ***REMOVED***
		return yaml_parser_set_scanner_error(parser, "while scanning a %YAML directive",
			start_mark, "did not find expected version number")
	***REMOVED***
	*number = value
	return true
***REMOVED***

// Scan the value of a TAG-DIRECTIVE token.
//
// Scope:
//      %TAG    !yaml!  tag:yaml.org,2002:  \n
//          ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
//
func yaml_parser_scan_tag_directive_value(parser *yaml_parser_t, start_mark yaml_mark_t, handle, prefix *[]byte) bool ***REMOVED***
	var handle_value, prefix_value []byte

	// Eat whitespaces.
	if parser.unread < 1 && !yaml_parser_update_buffer(parser, 1) ***REMOVED***
		return false
	***REMOVED***

	for is_blank(parser.buffer, parser.buffer_pos) ***REMOVED***
		skip(parser)
		if parser.unread < 1 && !yaml_parser_update_buffer(parser, 1) ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***

	// Scan a handle.
	if !yaml_parser_scan_tag_handle(parser, true, start_mark, &handle_value) ***REMOVED***
		return false
	***REMOVED***

	// Expect a whitespace.
	if parser.unread < 1 && !yaml_parser_update_buffer(parser, 1) ***REMOVED***
		return false
	***REMOVED***
	if !is_blank(parser.buffer, parser.buffer_pos) ***REMOVED***
		yaml_parser_set_scanner_error(parser, "while scanning a %TAG directive",
			start_mark, "did not find expected whitespace")
		return false
	***REMOVED***

	// Eat whitespaces.
	for is_blank(parser.buffer, parser.buffer_pos) ***REMOVED***
		skip(parser)
		if parser.unread < 1 && !yaml_parser_update_buffer(parser, 1) ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***

	// Scan a prefix.
	if !yaml_parser_scan_tag_uri(parser, true, nil, start_mark, &prefix_value) ***REMOVED***
		return false
	***REMOVED***

	// Expect a whitespace or line break.
	if parser.unread < 1 && !yaml_parser_update_buffer(parser, 1) ***REMOVED***
		return false
	***REMOVED***
	if !is_blankz(parser.buffer, parser.buffer_pos) ***REMOVED***
		yaml_parser_set_scanner_error(parser, "while scanning a %TAG directive",
			start_mark, "did not find expected whitespace or line break")
		return false
	***REMOVED***

	*handle = handle_value
	*prefix = prefix_value
	return true
***REMOVED***

func yaml_parser_scan_anchor(parser *yaml_parser_t, token *yaml_token_t, typ yaml_token_type_t) bool ***REMOVED***
	var s []byte

	// Eat the indicator character.
	start_mark := parser.mark
	skip(parser)

	// Consume the value.
	if parser.unread < 1 && !yaml_parser_update_buffer(parser, 1) ***REMOVED***
		return false
	***REMOVED***

	for is_alpha(parser.buffer, parser.buffer_pos) ***REMOVED***
		s = read(parser, s)
		if parser.unread < 1 && !yaml_parser_update_buffer(parser, 1) ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***

	end_mark := parser.mark

	/*
	 * Check if length of the anchor is greater than 0 and it is followed by
	 * a whitespace character or one of the indicators:
	 *
	 *      '?', ':', ',', ']', '***REMOVED***', '%', '@', '`'.
	 */

	if len(s) == 0 ||
		!(is_blankz(parser.buffer, parser.buffer_pos) || parser.buffer[parser.buffer_pos] == '?' ||
			parser.buffer[parser.buffer_pos] == ':' || parser.buffer[parser.buffer_pos] == ',' ||
			parser.buffer[parser.buffer_pos] == ']' || parser.buffer[parser.buffer_pos] == '***REMOVED***' ||
			parser.buffer[parser.buffer_pos] == '%' || parser.buffer[parser.buffer_pos] == '@' ||
			parser.buffer[parser.buffer_pos] == '`') ***REMOVED***
		context := "while scanning an alias"
		if typ == yaml_ANCHOR_TOKEN ***REMOVED***
			context = "while scanning an anchor"
		***REMOVED***
		yaml_parser_set_scanner_error(parser, context, start_mark,
			"did not find expected alphabetic or numeric character")
		return false
	***REMOVED***

	// Create a token.
	*token = yaml_token_t***REMOVED***
		typ:        typ,
		start_mark: start_mark,
		end_mark:   end_mark,
		value:      s,
	***REMOVED***

	return true
***REMOVED***

/*
 * Scan a TAG token.
 */

func yaml_parser_scan_tag(parser *yaml_parser_t, token *yaml_token_t) bool ***REMOVED***
	var handle, suffix []byte

	start_mark := parser.mark

	// Check if the tag is in the canonical form.
	if parser.unread < 2 && !yaml_parser_update_buffer(parser, 2) ***REMOVED***
		return false
	***REMOVED***

	if parser.buffer[parser.buffer_pos+1] == '<' ***REMOVED***
		// Keep the handle as ''

		// Eat '!<'
		skip(parser)
		skip(parser)

		// Consume the tag value.
		if !yaml_parser_scan_tag_uri(parser, false, nil, start_mark, &suffix) ***REMOVED***
			return false
		***REMOVED***

		// Check for '>' and eat it.
		if parser.buffer[parser.buffer_pos] != '>' ***REMOVED***
			yaml_parser_set_scanner_error(parser, "while scanning a tag",
				start_mark, "did not find the expected '>'")
			return false
		***REMOVED***

		skip(parser)
	***REMOVED*** else ***REMOVED***
		// The tag has either the '!suffix' or the '!handle!suffix' form.

		// First, try to scan a handle.
		if !yaml_parser_scan_tag_handle(parser, false, start_mark, &handle) ***REMOVED***
			return false
		***REMOVED***

		// Check if it is, indeed, handle.
		if handle[0] == '!' && len(handle) > 1 && handle[len(handle)-1] == '!' ***REMOVED***
			// Scan the suffix now.
			if !yaml_parser_scan_tag_uri(parser, false, nil, start_mark, &suffix) ***REMOVED***
				return false
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			// It wasn't a handle after all.  Scan the rest of the tag.
			if !yaml_parser_scan_tag_uri(parser, false, handle, start_mark, &suffix) ***REMOVED***
				return false
			***REMOVED***

			// Set the handle to '!'.
			handle = []byte***REMOVED***'!'***REMOVED***

			// A special case: the '!' tag.  Set the handle to '' and the
			// suffix to '!'.
			if len(suffix) == 0 ***REMOVED***
				handle, suffix = suffix, handle
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// Check the character which ends the tag.
	if parser.unread < 1 && !yaml_parser_update_buffer(parser, 1) ***REMOVED***
		return false
	***REMOVED***
	if !is_blankz(parser.buffer, parser.buffer_pos) ***REMOVED***
		yaml_parser_set_scanner_error(parser, "while scanning a tag",
			start_mark, "did not find expected whitespace or line break")
		return false
	***REMOVED***

	end_mark := parser.mark

	// Create a token.
	*token = yaml_token_t***REMOVED***
		typ:        yaml_TAG_TOKEN,
		start_mark: start_mark,
		end_mark:   end_mark,
		value:      handle,
		suffix:     suffix,
	***REMOVED***
	return true
***REMOVED***

// Scan a tag handle.
func yaml_parser_scan_tag_handle(parser *yaml_parser_t, directive bool, start_mark yaml_mark_t, handle *[]byte) bool ***REMOVED***
	// Check the initial '!' character.
	if parser.unread < 1 && !yaml_parser_update_buffer(parser, 1) ***REMOVED***
		return false
	***REMOVED***
	if parser.buffer[parser.buffer_pos] != '!' ***REMOVED***
		yaml_parser_set_scanner_tag_error(parser, directive,
			start_mark, "did not find expected '!'")
		return false
	***REMOVED***

	var s []byte

	// Copy the '!' character.
	s = read(parser, s)

	// Copy all subsequent alphabetical and numerical characters.
	if parser.unread < 1 && !yaml_parser_update_buffer(parser, 1) ***REMOVED***
		return false
	***REMOVED***
	for is_alpha(parser.buffer, parser.buffer_pos) ***REMOVED***
		s = read(parser, s)
		if parser.unread < 1 && !yaml_parser_update_buffer(parser, 1) ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***

	// Check if the trailing character is '!' and copy it.
	if parser.buffer[parser.buffer_pos] == '!' ***REMOVED***
		s = read(parser, s)
	***REMOVED*** else ***REMOVED***
		// It's either the '!' tag or not really a tag handle.  If it's a %TAG
		// directive, it's an error.  If it's a tag token, it must be a part of URI.
		if directive && string(s) != "!" ***REMOVED***
			yaml_parser_set_scanner_tag_error(parser, directive,
				start_mark, "did not find expected '!'")
			return false
		***REMOVED***
	***REMOVED***

	*handle = s
	return true
***REMOVED***

// Scan a tag.
func yaml_parser_scan_tag_uri(parser *yaml_parser_t, directive bool, head []byte, start_mark yaml_mark_t, uri *[]byte) bool ***REMOVED***
	//size_t length = head ? strlen((char *)head) : 0
	var s []byte
	hasTag := len(head) > 0

	// Copy the head if needed.
	//
	// Note that we don't copy the leading '!' character.
	if len(head) > 1 ***REMOVED***
		s = append(s, head[1:]...)
	***REMOVED***

	// Scan the tag.
	if parser.unread < 1 && !yaml_parser_update_buffer(parser, 1) ***REMOVED***
		return false
	***REMOVED***

	// The set of characters that may appear in URI is as follows:
	//
	//      '0'-'9', 'A'-'Z', 'a'-'z', '_', '-', ';', '/', '?', ':', '@', '&',
	//      '=', '+', '$', ',', '.', '!', '~', '*', '\'', '(', ')', '[', ']',
	//      '%'.
	// [Go] TODO Convert this into more reasonable logic.
	for is_alpha(parser.buffer, parser.buffer_pos) || parser.buffer[parser.buffer_pos] == ';' ||
		parser.buffer[parser.buffer_pos] == '/' || parser.buffer[parser.buffer_pos] == '?' ||
		parser.buffer[parser.buffer_pos] == ':' || parser.buffer[parser.buffer_pos] == '@' ||
		parser.buffer[parser.buffer_pos] == '&' || parser.buffer[parser.buffer_pos] == '=' ||
		parser.buffer[parser.buffer_pos] == '+' || parser.buffer[parser.buffer_pos] == '$' ||
		parser.buffer[parser.buffer_pos] == ',' || parser.buffer[parser.buffer_pos] == '.' ||
		parser.buffer[parser.buffer_pos] == '!' || parser.buffer[parser.buffer_pos] == '~' ||
		parser.buffer[parser.buffer_pos] == '*' || parser.buffer[parser.buffer_pos] == '\'' ||
		parser.buffer[parser.buffer_pos] == '(' || parser.buffer[parser.buffer_pos] == ')' ||
		parser.buffer[parser.buffer_pos] == '[' || parser.buffer[parser.buffer_pos] == ']' ||
		parser.buffer[parser.buffer_pos] == '%' ***REMOVED***
		// Check if it is a URI-escape sequence.
		if parser.buffer[parser.buffer_pos] == '%' ***REMOVED***
			if !yaml_parser_scan_uri_escapes(parser, directive, start_mark, &s) ***REMOVED***
				return false
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			s = read(parser, s)
		***REMOVED***
		if parser.unread < 1 && !yaml_parser_update_buffer(parser, 1) ***REMOVED***
			return false
		***REMOVED***
		hasTag = true
	***REMOVED***

	if !hasTag ***REMOVED***
		yaml_parser_set_scanner_tag_error(parser, directive,
			start_mark, "did not find expected tag URI")
		return false
	***REMOVED***
	*uri = s
	return true
***REMOVED***

// Decode an URI-escape sequence corresponding to a single UTF-8 character.
func yaml_parser_scan_uri_escapes(parser *yaml_parser_t, directive bool, start_mark yaml_mark_t, s *[]byte) bool ***REMOVED***

	// Decode the required number of characters.
	w := 1024
	for w > 0 ***REMOVED***
		// Check for a URI-escaped octet.
		if parser.unread < 3 && !yaml_parser_update_buffer(parser, 3) ***REMOVED***
			return false
		***REMOVED***

		if !(parser.buffer[parser.buffer_pos] == '%' &&
			is_hex(parser.buffer, parser.buffer_pos+1) &&
			is_hex(parser.buffer, parser.buffer_pos+2)) ***REMOVED***
			return yaml_parser_set_scanner_tag_error(parser, directive,
				start_mark, "did not find URI escaped octet")
		***REMOVED***

		// Get the octet.
		octet := byte((as_hex(parser.buffer, parser.buffer_pos+1) << 4) + as_hex(parser.buffer, parser.buffer_pos+2))

		// If it is the leading octet, determine the length of the UTF-8 sequence.
		if w == 1024 ***REMOVED***
			w = width(octet)
			if w == 0 ***REMOVED***
				return yaml_parser_set_scanner_tag_error(parser, directive,
					start_mark, "found an incorrect leading UTF-8 octet")
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			// Check if the trailing octet is correct.
			if octet&0xC0 != 0x80 ***REMOVED***
				return yaml_parser_set_scanner_tag_error(parser, directive,
					start_mark, "found an incorrect trailing UTF-8 octet")
			***REMOVED***
		***REMOVED***

		// Copy the octet and move the pointers.
		*s = append(*s, octet)
		skip(parser)
		skip(parser)
		skip(parser)
		w--
	***REMOVED***
	return true
***REMOVED***

// Scan a block scalar.
func yaml_parser_scan_block_scalar(parser *yaml_parser_t, token *yaml_token_t, literal bool) bool ***REMOVED***
	// Eat the indicator '|' or '>'.
	start_mark := parser.mark
	skip(parser)

	// Scan the additional block scalar indicators.
	if parser.unread < 1 && !yaml_parser_update_buffer(parser, 1) ***REMOVED***
		return false
	***REMOVED***

	// Check for a chomping indicator.
	var chomping, increment int
	if parser.buffer[parser.buffer_pos] == '+' || parser.buffer[parser.buffer_pos] == '-' ***REMOVED***
		// Set the chomping method and eat the indicator.
		if parser.buffer[parser.buffer_pos] == '+' ***REMOVED***
			chomping = +1
		***REMOVED*** else ***REMOVED***
			chomping = -1
		***REMOVED***
		skip(parser)

		// Check for an indentation indicator.
		if parser.unread < 1 && !yaml_parser_update_buffer(parser, 1) ***REMOVED***
			return false
		***REMOVED***
		if is_digit(parser.buffer, parser.buffer_pos) ***REMOVED***
			// Check that the indentation is greater than 0.
			if parser.buffer[parser.buffer_pos] == '0' ***REMOVED***
				yaml_parser_set_scanner_error(parser, "while scanning a block scalar",
					start_mark, "found an indentation indicator equal to 0")
				return false
			***REMOVED***

			// Get the indentation level and eat the indicator.
			increment = as_digit(parser.buffer, parser.buffer_pos)
			skip(parser)
		***REMOVED***

	***REMOVED*** else if is_digit(parser.buffer, parser.buffer_pos) ***REMOVED***
		// Do the same as above, but in the opposite order.

		if parser.buffer[parser.buffer_pos] == '0' ***REMOVED***
			yaml_parser_set_scanner_error(parser, "while scanning a block scalar",
				start_mark, "found an indentation indicator equal to 0")
			return false
		***REMOVED***
		increment = as_digit(parser.buffer, parser.buffer_pos)
		skip(parser)

		if parser.unread < 1 && !yaml_parser_update_buffer(parser, 1) ***REMOVED***
			return false
		***REMOVED***
		if parser.buffer[parser.buffer_pos] == '+' || parser.buffer[parser.buffer_pos] == '-' ***REMOVED***
			if parser.buffer[parser.buffer_pos] == '+' ***REMOVED***
				chomping = +1
			***REMOVED*** else ***REMOVED***
				chomping = -1
			***REMOVED***
			skip(parser)
		***REMOVED***
	***REMOVED***

	// Eat whitespaces and comments to the end of the line.
	if parser.unread < 1 && !yaml_parser_update_buffer(parser, 1) ***REMOVED***
		return false
	***REMOVED***
	for is_blank(parser.buffer, parser.buffer_pos) ***REMOVED***
		skip(parser)
		if parser.unread < 1 && !yaml_parser_update_buffer(parser, 1) ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	if parser.buffer[parser.buffer_pos] == '#' ***REMOVED***
		// TODO Test this and then re-enable it.
		//if !yaml_parser_scan_line_comment(parser, start_mark) ***REMOVED***
		//	return false
		//***REMOVED***
		for !is_breakz(parser.buffer, parser.buffer_pos) ***REMOVED***
			skip(parser)
			if parser.unread < 1 && !yaml_parser_update_buffer(parser, 1) ***REMOVED***
				return false
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// Check if we are at the end of the line.
	if !is_breakz(parser.buffer, parser.buffer_pos) ***REMOVED***
		yaml_parser_set_scanner_error(parser, "while scanning a block scalar",
			start_mark, "did not find expected comment or line break")
		return false
	***REMOVED***

	// Eat a line break.
	if is_break(parser.buffer, parser.buffer_pos) ***REMOVED***
		if parser.unread < 2 && !yaml_parser_update_buffer(parser, 2) ***REMOVED***
			return false
		***REMOVED***
		skip_line(parser)
	***REMOVED***

	end_mark := parser.mark

	// Set the indentation level if it was specified.
	var indent int
	if increment > 0 ***REMOVED***
		if parser.indent >= 0 ***REMOVED***
			indent = parser.indent + increment
		***REMOVED*** else ***REMOVED***
			indent = increment
		***REMOVED***
	***REMOVED***

	// Scan the leading line breaks and determine the indentation level if needed.
	var s, leading_break, trailing_breaks []byte
	if !yaml_parser_scan_block_scalar_breaks(parser, &indent, &trailing_breaks, start_mark, &end_mark) ***REMOVED***
		return false
	***REMOVED***

	// Scan the block scalar content.
	if parser.unread < 1 && !yaml_parser_update_buffer(parser, 1) ***REMOVED***
		return false
	***REMOVED***
	var leading_blank, trailing_blank bool
	for parser.mark.column == indent && !is_z(parser.buffer, parser.buffer_pos) ***REMOVED***
		// We are at the beginning of a non-empty line.

		// Is it a trailing whitespace?
		trailing_blank = is_blank(parser.buffer, parser.buffer_pos)

		// Check if we need to fold the leading line break.
		if !literal && !leading_blank && !trailing_blank && len(leading_break) > 0 && leading_break[0] == '\n' ***REMOVED***
			// Do we need to join the lines by space?
			if len(trailing_breaks) == 0 ***REMOVED***
				s = append(s, ' ')
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			s = append(s, leading_break...)
		***REMOVED***
		leading_break = leading_break[:0]

		// Append the remaining line breaks.
		s = append(s, trailing_breaks...)
		trailing_breaks = trailing_breaks[:0]

		// Is it a leading whitespace?
		leading_blank = is_blank(parser.buffer, parser.buffer_pos)

		// Consume the current line.
		for !is_breakz(parser.buffer, parser.buffer_pos) ***REMOVED***
			s = read(parser, s)
			if parser.unread < 1 && !yaml_parser_update_buffer(parser, 1) ***REMOVED***
				return false
			***REMOVED***
		***REMOVED***

		// Consume the line break.
		if parser.unread < 2 && !yaml_parser_update_buffer(parser, 2) ***REMOVED***
			return false
		***REMOVED***

		leading_break = read_line(parser, leading_break)

		// Eat the following indentation spaces and line breaks.
		if !yaml_parser_scan_block_scalar_breaks(parser, &indent, &trailing_breaks, start_mark, &end_mark) ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***

	// Chomp the tail.
	if chomping != -1 ***REMOVED***
		s = append(s, leading_break...)
	***REMOVED***
	if chomping == 1 ***REMOVED***
		s = append(s, trailing_breaks...)
	***REMOVED***

	// Create a token.
	*token = yaml_token_t***REMOVED***
		typ:        yaml_SCALAR_TOKEN,
		start_mark: start_mark,
		end_mark:   end_mark,
		value:      s,
		style:      yaml_LITERAL_SCALAR_STYLE,
	***REMOVED***
	if !literal ***REMOVED***
		token.style = yaml_FOLDED_SCALAR_STYLE
	***REMOVED***
	return true
***REMOVED***

// Scan indentation spaces and line breaks for a block scalar.  Determine the
// indentation level if needed.
func yaml_parser_scan_block_scalar_breaks(parser *yaml_parser_t, indent *int, breaks *[]byte, start_mark yaml_mark_t, end_mark *yaml_mark_t) bool ***REMOVED***
	*end_mark = parser.mark

	// Eat the indentation spaces and line breaks.
	max_indent := 0
	for ***REMOVED***
		// Eat the indentation spaces.
		if parser.unread < 1 && !yaml_parser_update_buffer(parser, 1) ***REMOVED***
			return false
		***REMOVED***
		for (*indent == 0 || parser.mark.column < *indent) && is_space(parser.buffer, parser.buffer_pos) ***REMOVED***
			skip(parser)
			if parser.unread < 1 && !yaml_parser_update_buffer(parser, 1) ***REMOVED***
				return false
			***REMOVED***
		***REMOVED***
		if parser.mark.column > max_indent ***REMOVED***
			max_indent = parser.mark.column
		***REMOVED***

		// Check for a tab character messing the indentation.
		if (*indent == 0 || parser.mark.column < *indent) && is_tab(parser.buffer, parser.buffer_pos) ***REMOVED***
			return yaml_parser_set_scanner_error(parser, "while scanning a block scalar",
				start_mark, "found a tab character where an indentation space is expected")
		***REMOVED***

		// Have we found a non-empty line?
		if !is_break(parser.buffer, parser.buffer_pos) ***REMOVED***
			break
		***REMOVED***

		// Consume the line break.
		if parser.unread < 2 && !yaml_parser_update_buffer(parser, 2) ***REMOVED***
			return false
		***REMOVED***
		// [Go] Should really be returning breaks instead.
		*breaks = read_line(parser, *breaks)
		*end_mark = parser.mark
	***REMOVED***

	// Determine the indentation level if needed.
	if *indent == 0 ***REMOVED***
		*indent = max_indent
		if *indent < parser.indent+1 ***REMOVED***
			*indent = parser.indent + 1
		***REMOVED***
		if *indent < 1 ***REMOVED***
			*indent = 1
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

// Scan a quoted scalar.
func yaml_parser_scan_flow_scalar(parser *yaml_parser_t, token *yaml_token_t, single bool) bool ***REMOVED***
	// Eat the left quote.
	start_mark := parser.mark
	skip(parser)

	// Consume the content of the quoted scalar.
	var s, leading_break, trailing_breaks, whitespaces []byte
	for ***REMOVED***
		// Check that there are no document indicators at the beginning of the line.
		if parser.unread < 4 && !yaml_parser_update_buffer(parser, 4) ***REMOVED***
			return false
		***REMOVED***

		if parser.mark.column == 0 &&
			((parser.buffer[parser.buffer_pos+0] == '-' &&
				parser.buffer[parser.buffer_pos+1] == '-' &&
				parser.buffer[parser.buffer_pos+2] == '-') ||
				(parser.buffer[parser.buffer_pos+0] == '.' &&
					parser.buffer[parser.buffer_pos+1] == '.' &&
					parser.buffer[parser.buffer_pos+2] == '.')) &&
			is_blankz(parser.buffer, parser.buffer_pos+3) ***REMOVED***
			yaml_parser_set_scanner_error(parser, "while scanning a quoted scalar",
				start_mark, "found unexpected document indicator")
			return false
		***REMOVED***

		// Check for EOF.
		if is_z(parser.buffer, parser.buffer_pos) ***REMOVED***
			yaml_parser_set_scanner_error(parser, "while scanning a quoted scalar",
				start_mark, "found unexpected end of stream")
			return false
		***REMOVED***

		// Consume non-blank characters.
		leading_blanks := false
		for !is_blankz(parser.buffer, parser.buffer_pos) ***REMOVED***
			if single && parser.buffer[parser.buffer_pos] == '\'' && parser.buffer[parser.buffer_pos+1] == '\'' ***REMOVED***
				// Is is an escaped single quote.
				s = append(s, '\'')
				skip(parser)
				skip(parser)

			***REMOVED*** else if single && parser.buffer[parser.buffer_pos] == '\'' ***REMOVED***
				// It is a right single quote.
				break
			***REMOVED*** else if !single && parser.buffer[parser.buffer_pos] == '"' ***REMOVED***
				// It is a right double quote.
				break

			***REMOVED*** else if !single && parser.buffer[parser.buffer_pos] == '\\' && is_break(parser.buffer, parser.buffer_pos+1) ***REMOVED***
				// It is an escaped line break.
				if parser.unread < 3 && !yaml_parser_update_buffer(parser, 3) ***REMOVED***
					return false
				***REMOVED***
				skip(parser)
				skip_line(parser)
				leading_blanks = true
				break

			***REMOVED*** else if !single && parser.buffer[parser.buffer_pos] == '\\' ***REMOVED***
				// It is an escape sequence.
				code_length := 0

				// Check the escape character.
				switch parser.buffer[parser.buffer_pos+1] ***REMOVED***
				case '0':
					s = append(s, 0)
				case 'a':
					s = append(s, '\x07')
				case 'b':
					s = append(s, '\x08')
				case 't', '\t':
					s = append(s, '\x09')
				case 'n':
					s = append(s, '\x0A')
				case 'v':
					s = append(s, '\x0B')
				case 'f':
					s = append(s, '\x0C')
				case 'r':
					s = append(s, '\x0D')
				case 'e':
					s = append(s, '\x1B')
				case ' ':
					s = append(s, '\x20')
				case '"':
					s = append(s, '"')
				case '\'':
					s = append(s, '\'')
				case '\\':
					s = append(s, '\\')
				case 'N': // NEL (#x85)
					s = append(s, '\xC2')
					s = append(s, '\x85')
				case '_': // #xA0
					s = append(s, '\xC2')
					s = append(s, '\xA0')
				case 'L': // LS (#x2028)
					s = append(s, '\xE2')
					s = append(s, '\x80')
					s = append(s, '\xA8')
				case 'P': // PS (#x2029)
					s = append(s, '\xE2')
					s = append(s, '\x80')
					s = append(s, '\xA9')
				case 'x':
					code_length = 2
				case 'u':
					code_length = 4
				case 'U':
					code_length = 8
				default:
					yaml_parser_set_scanner_error(parser, "while parsing a quoted scalar",
						start_mark, "found unknown escape character")
					return false
				***REMOVED***

				skip(parser)
				skip(parser)

				// Consume an arbitrary escape code.
				if code_length > 0 ***REMOVED***
					var value int

					// Scan the character value.
					if parser.unread < code_length && !yaml_parser_update_buffer(parser, code_length) ***REMOVED***
						return false
					***REMOVED***
					for k := 0; k < code_length; k++ ***REMOVED***
						if !is_hex(parser.buffer, parser.buffer_pos+k) ***REMOVED***
							yaml_parser_set_scanner_error(parser, "while parsing a quoted scalar",
								start_mark, "did not find expected hexdecimal number")
							return false
						***REMOVED***
						value = (value << 4) + as_hex(parser.buffer, parser.buffer_pos+k)
					***REMOVED***

					// Check the value and write the character.
					if (value >= 0xD800 && value <= 0xDFFF) || value > 0x10FFFF ***REMOVED***
						yaml_parser_set_scanner_error(parser, "while parsing a quoted scalar",
							start_mark, "found invalid Unicode character escape code")
						return false
					***REMOVED***
					if value <= 0x7F ***REMOVED***
						s = append(s, byte(value))
					***REMOVED*** else if value <= 0x7FF ***REMOVED***
						s = append(s, byte(0xC0+(value>>6)))
						s = append(s, byte(0x80+(value&0x3F)))
					***REMOVED*** else if value <= 0xFFFF ***REMOVED***
						s = append(s, byte(0xE0+(value>>12)))
						s = append(s, byte(0x80+((value>>6)&0x3F)))
						s = append(s, byte(0x80+(value&0x3F)))
					***REMOVED*** else ***REMOVED***
						s = append(s, byte(0xF0+(value>>18)))
						s = append(s, byte(0x80+((value>>12)&0x3F)))
						s = append(s, byte(0x80+((value>>6)&0x3F)))
						s = append(s, byte(0x80+(value&0x3F)))
					***REMOVED***

					// Advance the pointer.
					for k := 0; k < code_length; k++ ***REMOVED***
						skip(parser)
					***REMOVED***
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				// It is a non-escaped non-blank character.
				s = read(parser, s)
			***REMOVED***
			if parser.unread < 2 && !yaml_parser_update_buffer(parser, 2) ***REMOVED***
				return false
			***REMOVED***
		***REMOVED***

		if parser.unread < 1 && !yaml_parser_update_buffer(parser, 1) ***REMOVED***
			return false
		***REMOVED***

		// Check if we are at the end of the scalar.
		if single ***REMOVED***
			if parser.buffer[parser.buffer_pos] == '\'' ***REMOVED***
				break
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			if parser.buffer[parser.buffer_pos] == '"' ***REMOVED***
				break
			***REMOVED***
		***REMOVED***

		// Consume blank characters.
		for is_blank(parser.buffer, parser.buffer_pos) || is_break(parser.buffer, parser.buffer_pos) ***REMOVED***
			if is_blank(parser.buffer, parser.buffer_pos) ***REMOVED***
				// Consume a space or a tab character.
				if !leading_blanks ***REMOVED***
					whitespaces = read(parser, whitespaces)
				***REMOVED*** else ***REMOVED***
					skip(parser)
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				if parser.unread < 2 && !yaml_parser_update_buffer(parser, 2) ***REMOVED***
					return false
				***REMOVED***

				// Check if it is a first line break.
				if !leading_blanks ***REMOVED***
					whitespaces = whitespaces[:0]
					leading_break = read_line(parser, leading_break)
					leading_blanks = true
				***REMOVED*** else ***REMOVED***
					trailing_breaks = read_line(parser, trailing_breaks)
				***REMOVED***
			***REMOVED***
			if parser.unread < 1 && !yaml_parser_update_buffer(parser, 1) ***REMOVED***
				return false
			***REMOVED***
		***REMOVED***

		// Join the whitespaces or fold line breaks.
		if leading_blanks ***REMOVED***
			// Do we need to fold line breaks?
			if len(leading_break) > 0 && leading_break[0] == '\n' ***REMOVED***
				if len(trailing_breaks) == 0 ***REMOVED***
					s = append(s, ' ')
				***REMOVED*** else ***REMOVED***
					s = append(s, trailing_breaks...)
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				s = append(s, leading_break...)
				s = append(s, trailing_breaks...)
			***REMOVED***
			trailing_breaks = trailing_breaks[:0]
			leading_break = leading_break[:0]
		***REMOVED*** else ***REMOVED***
			s = append(s, whitespaces...)
			whitespaces = whitespaces[:0]
		***REMOVED***
	***REMOVED***

	// Eat the right quote.
	skip(parser)
	end_mark := parser.mark

	// Create a token.
	*token = yaml_token_t***REMOVED***
		typ:        yaml_SCALAR_TOKEN,
		start_mark: start_mark,
		end_mark:   end_mark,
		value:      s,
		style:      yaml_SINGLE_QUOTED_SCALAR_STYLE,
	***REMOVED***
	if !single ***REMOVED***
		token.style = yaml_DOUBLE_QUOTED_SCALAR_STYLE
	***REMOVED***
	return true
***REMOVED***

// Scan a plain scalar.
func yaml_parser_scan_plain_scalar(parser *yaml_parser_t, token *yaml_token_t) bool ***REMOVED***

	var s, leading_break, trailing_breaks, whitespaces []byte
	var leading_blanks bool
	var indent = parser.indent + 1

	start_mark := parser.mark
	end_mark := parser.mark

	// Consume the content of the plain scalar.
	for ***REMOVED***
		// Check for a document indicator.
		if parser.unread < 4 && !yaml_parser_update_buffer(parser, 4) ***REMOVED***
			return false
		***REMOVED***
		if parser.mark.column == 0 &&
			((parser.buffer[parser.buffer_pos+0] == '-' &&
				parser.buffer[parser.buffer_pos+1] == '-' &&
				parser.buffer[parser.buffer_pos+2] == '-') ||
				(parser.buffer[parser.buffer_pos+0] == '.' &&
					parser.buffer[parser.buffer_pos+1] == '.' &&
					parser.buffer[parser.buffer_pos+2] == '.')) &&
			is_blankz(parser.buffer, parser.buffer_pos+3) ***REMOVED***
			break
		***REMOVED***

		// Check for a comment.
		if parser.buffer[parser.buffer_pos] == '#' ***REMOVED***
			break
		***REMOVED***

		// Consume non-blank characters.
		for !is_blankz(parser.buffer, parser.buffer_pos) ***REMOVED***

			// Check for indicators that may end a plain scalar.
			if (parser.buffer[parser.buffer_pos] == ':' && is_blankz(parser.buffer, parser.buffer_pos+1)) ||
				(parser.flow_level > 0 &&
					(parser.buffer[parser.buffer_pos] == ',' ||
						parser.buffer[parser.buffer_pos] == '?' || parser.buffer[parser.buffer_pos] == '[' ||
						parser.buffer[parser.buffer_pos] == ']' || parser.buffer[parser.buffer_pos] == '***REMOVED***' ||
						parser.buffer[parser.buffer_pos] == '***REMOVED***')) ***REMOVED***
				break
			***REMOVED***

			// Check if we need to join whitespaces and breaks.
			if leading_blanks || len(whitespaces) > 0 ***REMOVED***
				if leading_blanks ***REMOVED***
					// Do we need to fold line breaks?
					if leading_break[0] == '\n' ***REMOVED***
						if len(trailing_breaks) == 0 ***REMOVED***
							s = append(s, ' ')
						***REMOVED*** else ***REMOVED***
							s = append(s, trailing_breaks...)
						***REMOVED***
					***REMOVED*** else ***REMOVED***
						s = append(s, leading_break...)
						s = append(s, trailing_breaks...)
					***REMOVED***
					trailing_breaks = trailing_breaks[:0]
					leading_break = leading_break[:0]
					leading_blanks = false
				***REMOVED*** else ***REMOVED***
					s = append(s, whitespaces...)
					whitespaces = whitespaces[:0]
				***REMOVED***
			***REMOVED***

			// Copy the character.
			s = read(parser, s)

			end_mark = parser.mark
			if parser.unread < 2 && !yaml_parser_update_buffer(parser, 2) ***REMOVED***
				return false
			***REMOVED***
		***REMOVED***

		// Is it the end?
		if !(is_blank(parser.buffer, parser.buffer_pos) || is_break(parser.buffer, parser.buffer_pos)) ***REMOVED***
			break
		***REMOVED***

		// Consume blank characters.
		if parser.unread < 1 && !yaml_parser_update_buffer(parser, 1) ***REMOVED***
			return false
		***REMOVED***

		for is_blank(parser.buffer, parser.buffer_pos) || is_break(parser.buffer, parser.buffer_pos) ***REMOVED***
			if is_blank(parser.buffer, parser.buffer_pos) ***REMOVED***

				// Check for tab characters that abuse indentation.
				if leading_blanks && parser.mark.column < indent && is_tab(parser.buffer, parser.buffer_pos) ***REMOVED***
					yaml_parser_set_scanner_error(parser, "while scanning a plain scalar",
						start_mark, "found a tab character that violates indentation")
					return false
				***REMOVED***

				// Consume a space or a tab character.
				if !leading_blanks ***REMOVED***
					whitespaces = read(parser, whitespaces)
				***REMOVED*** else ***REMOVED***
					skip(parser)
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				if parser.unread < 2 && !yaml_parser_update_buffer(parser, 2) ***REMOVED***
					return false
				***REMOVED***

				// Check if it is a first line break.
				if !leading_blanks ***REMOVED***
					whitespaces = whitespaces[:0]
					leading_break = read_line(parser, leading_break)
					leading_blanks = true
				***REMOVED*** else ***REMOVED***
					trailing_breaks = read_line(parser, trailing_breaks)
				***REMOVED***
			***REMOVED***
			if parser.unread < 1 && !yaml_parser_update_buffer(parser, 1) ***REMOVED***
				return false
			***REMOVED***
		***REMOVED***

		// Check indentation level.
		if parser.flow_level == 0 && parser.mark.column < indent ***REMOVED***
			break
		***REMOVED***
	***REMOVED***

	// Create a token.
	*token = yaml_token_t***REMOVED***
		typ:        yaml_SCALAR_TOKEN,
		start_mark: start_mark,
		end_mark:   end_mark,
		value:      s,
		style:      yaml_PLAIN_SCALAR_STYLE,
	***REMOVED***

	// Note that we change the 'simple_key_allowed' flag.
	if leading_blanks ***REMOVED***
		parser.simple_key_allowed = true
	***REMOVED***
	return true
***REMOVED***

func yaml_parser_scan_line_comment(parser *yaml_parser_t, token_mark yaml_mark_t) bool ***REMOVED***
	if parser.newlines > 0 ***REMOVED***
		return true
	***REMOVED***

	var start_mark yaml_mark_t
	var text []byte

	for peek := 0; peek < 512; peek++ ***REMOVED***
		if parser.unread < peek+1 && !yaml_parser_update_buffer(parser, peek+1) ***REMOVED***
			break
		***REMOVED***
		if is_blank(parser.buffer, parser.buffer_pos+peek) ***REMOVED***
			continue
		***REMOVED***
		if parser.buffer[parser.buffer_pos+peek] == '#' ***REMOVED***
			seen := parser.mark.index+peek
			for ***REMOVED***
				if parser.unread < 1 && !yaml_parser_update_buffer(parser, 1) ***REMOVED***
					return false
				***REMOVED***
				if is_breakz(parser.buffer, parser.buffer_pos) ***REMOVED***
					if parser.mark.index >= seen ***REMOVED***
						break
					***REMOVED***
					if parser.unread < 2 && !yaml_parser_update_buffer(parser, 2) ***REMOVED***
						return false
					***REMOVED***
					skip_line(parser)
				***REMOVED*** else ***REMOVED***
					if parser.mark.index >= seen ***REMOVED***
						if len(text) == 0 ***REMOVED***
							start_mark = parser.mark
						***REMOVED***
						text = append(text, parser.buffer[parser.buffer_pos])
					***REMOVED***
					skip(parser)
				***REMOVED***
			***REMOVED***
		***REMOVED***
		break
	***REMOVED***
	if len(text) > 0 ***REMOVED***
		parser.comments = append(parser.comments, yaml_comment_t***REMOVED***
			token_mark: token_mark,
			start_mark: start_mark,
			line: text,
		***REMOVED***)
	***REMOVED***
	return true
***REMOVED***

func yaml_parser_scan_comments(parser *yaml_parser_t, scan_mark yaml_mark_t) bool ***REMOVED***
	token := parser.tokens[len(parser.tokens)-1]

	if token.typ == yaml_FLOW_ENTRY_TOKEN && len(parser.tokens) > 1 ***REMOVED***
		token = parser.tokens[len(parser.tokens)-2]
	***REMOVED***

	var token_mark = token.start_mark
	var start_mark yaml_mark_t

	var recent_empty = false
	var first_empty = parser.newlines <= 1

	var line = parser.mark.line
	var column = parser.mark.column

	var text []byte

	// The foot line is the place where a comment must start to
	// still be considered as a foot of the prior content.
	// If there's some content in the currently parsed line, then
	// the foot is the line below it.
	var foot_line = -1
	if scan_mark.line > 0 ***REMOVED***
		foot_line = parser.mark.line-parser.newlines+1
		if parser.newlines == 0 && parser.mark.column > 1 ***REMOVED***
			foot_line++
		***REMOVED***
	***REMOVED***

	var peek = 0
	for ; peek < 512; peek++ ***REMOVED***
		if parser.unread < peek+1 && !yaml_parser_update_buffer(parser, peek+1) ***REMOVED***
			break
		***REMOVED***
		column++
		if is_blank(parser.buffer, parser.buffer_pos+peek) ***REMOVED***
			continue
		***REMOVED***
		c := parser.buffer[parser.buffer_pos+peek]
		if is_breakz(parser.buffer, parser.buffer_pos+peek) || parser.flow_level > 0 && (c == ']' || c == '***REMOVED***') ***REMOVED***
			// Got line break or terminator.
			if !recent_empty ***REMOVED***
				if first_empty && (start_mark.line == foot_line || start_mark.column-1 < parser.indent) ***REMOVED***
					// This is the first empty line and there were no empty lines before,
					// so this initial part of the comment is a foot of the prior token
					// instead of being a head for the following one. Split it up.
					if len(text) > 0 ***REMOVED***
						if start_mark.column-1 < parser.indent ***REMOVED***
							// If dedented it's unrelated to the prior token.
							token_mark = start_mark
						***REMOVED***
						parser.comments = append(parser.comments, yaml_comment_t***REMOVED***
							scan_mark:  scan_mark,
							token_mark: token_mark,
							start_mark: start_mark,
							end_mark:   yaml_mark_t***REMOVED***parser.mark.index + peek, line, column***REMOVED***,
							foot:       text,
						***REMOVED***)
						scan_mark = yaml_mark_t***REMOVED***parser.mark.index + peek, line, column***REMOVED***
						token_mark = scan_mark
						text = nil
					***REMOVED***
				***REMOVED*** else ***REMOVED***
					if len(text) > 0 && parser.buffer[parser.buffer_pos+peek] != 0 ***REMOVED***
						text = append(text, '\n')
					***REMOVED***
				***REMOVED***
			***REMOVED***
			if !is_break(parser.buffer, parser.buffer_pos+peek) ***REMOVED***
				break
			***REMOVED***
			first_empty = false
			recent_empty = true
			column = 0
			line++
			continue
		***REMOVED***

		if len(text) > 0 && column < parser.indent+1 && column != start_mark.column ***REMOVED***
			// The comment at the different indentation is a foot of the
			// preceding data rather than a head of the upcoming one.
			parser.comments = append(parser.comments, yaml_comment_t***REMOVED***
				scan_mark:  scan_mark,
				token_mark: token_mark,
				start_mark: start_mark,
				end_mark:   yaml_mark_t***REMOVED***parser.mark.index + peek, line, column***REMOVED***,
				foot:       text,
			***REMOVED***)
			scan_mark = yaml_mark_t***REMOVED***parser.mark.index + peek, line, column***REMOVED***
			token_mark = scan_mark
			text = nil
		***REMOVED***

		if parser.buffer[parser.buffer_pos+peek] != '#' ***REMOVED***
			break
		***REMOVED***

		if len(text) == 0 ***REMOVED***
			start_mark = yaml_mark_t***REMOVED***parser.mark.index + peek, line, column***REMOVED***
		***REMOVED*** else ***REMOVED***
			text = append(text, '\n')
		***REMOVED***

		recent_empty = false

		// Consume until after the consumed comment line.
		seen := parser.mark.index+peek
		for ***REMOVED***
			if parser.unread < 1 && !yaml_parser_update_buffer(parser, 1) ***REMOVED***
				return false
			***REMOVED***
			if is_breakz(parser.buffer, parser.buffer_pos) ***REMOVED***
				if parser.mark.index >= seen ***REMOVED***
					break
				***REMOVED***
				if parser.unread < 2 && !yaml_parser_update_buffer(parser, 2) ***REMOVED***
					return false
				***REMOVED***
				skip_line(parser)
			***REMOVED*** else ***REMOVED***
				if parser.mark.index >= seen ***REMOVED***
					text = append(text, parser.buffer[parser.buffer_pos])
				***REMOVED***
				skip(parser)
			***REMOVED***
		***REMOVED***

		peek = 0
		column = 0
		line = parser.mark.line
	***REMOVED***

	if len(text) > 0 ***REMOVED***
		parser.comments = append(parser.comments, yaml_comment_t***REMOVED***
			scan_mark:  scan_mark,
			token_mark: start_mark,
			start_mark: start_mark,
			end_mark:   yaml_mark_t***REMOVED***parser.mark.index + peek - 1, line, column***REMOVED***,
			head:       text,
		***REMOVED***)
	***REMOVED***
	return true
***REMOVED***