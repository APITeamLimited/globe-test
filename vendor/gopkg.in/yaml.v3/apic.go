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
	"io"
)

func yaml_insert_token(parser *yaml_parser_t, pos int, token *yaml_token_t) ***REMOVED***
	//fmt.Println("yaml_insert_token", "pos:", pos, "typ:", token.typ, "head:", parser.tokens_head, "len:", len(parser.tokens))

	// Check if we can move the queue at the beginning of the buffer.
	if parser.tokens_head > 0 && len(parser.tokens) == cap(parser.tokens) ***REMOVED***
		if parser.tokens_head != len(parser.tokens) ***REMOVED***
			copy(parser.tokens, parser.tokens[parser.tokens_head:])
		***REMOVED***
		parser.tokens = parser.tokens[:len(parser.tokens)-parser.tokens_head]
		parser.tokens_head = 0
	***REMOVED***
	parser.tokens = append(parser.tokens, *token)
	if pos < 0 ***REMOVED***
		return
	***REMOVED***
	copy(parser.tokens[parser.tokens_head+pos+1:], parser.tokens[parser.tokens_head+pos:])
	parser.tokens[parser.tokens_head+pos] = *token
***REMOVED***

// Create a new parser object.
func yaml_parser_initialize(parser *yaml_parser_t) bool ***REMOVED***
	*parser = yaml_parser_t***REMOVED***
		raw_buffer: make([]byte, 0, input_raw_buffer_size),
		buffer:     make([]byte, 0, input_buffer_size),
	***REMOVED***
	return true
***REMOVED***

// Destroy a parser object.
func yaml_parser_delete(parser *yaml_parser_t) ***REMOVED***
	*parser = yaml_parser_t***REMOVED******REMOVED***
***REMOVED***

// String read handler.
func yaml_string_read_handler(parser *yaml_parser_t, buffer []byte) (n int, err error) ***REMOVED***
	if parser.input_pos == len(parser.input) ***REMOVED***
		return 0, io.EOF
	***REMOVED***
	n = copy(buffer, parser.input[parser.input_pos:])
	parser.input_pos += n
	return n, nil
***REMOVED***

// Reader read handler.
func yaml_reader_read_handler(parser *yaml_parser_t, buffer []byte) (n int, err error) ***REMOVED***
	return parser.input_reader.Read(buffer)
***REMOVED***

// Set a string input.
func yaml_parser_set_input_string(parser *yaml_parser_t, input []byte) ***REMOVED***
	if parser.read_handler != nil ***REMOVED***
		panic("must set the input source only once")
	***REMOVED***
	parser.read_handler = yaml_string_read_handler
	parser.input = input
	parser.input_pos = 0
***REMOVED***

// Set a file input.
func yaml_parser_set_input_reader(parser *yaml_parser_t, r io.Reader) ***REMOVED***
	if parser.read_handler != nil ***REMOVED***
		panic("must set the input source only once")
	***REMOVED***
	parser.read_handler = yaml_reader_read_handler
	parser.input_reader = r
***REMOVED***

// Set the source encoding.
func yaml_parser_set_encoding(parser *yaml_parser_t, encoding yaml_encoding_t) ***REMOVED***
	if parser.encoding != yaml_ANY_ENCODING ***REMOVED***
		panic("must set the encoding only once")
	***REMOVED***
	parser.encoding = encoding
***REMOVED***

// Create a new emitter object.
func yaml_emitter_initialize(emitter *yaml_emitter_t) ***REMOVED***
	*emitter = yaml_emitter_t***REMOVED***
		buffer:     make([]byte, output_buffer_size),
		raw_buffer: make([]byte, 0, output_raw_buffer_size),
		states:     make([]yaml_emitter_state_t, 0, initial_stack_size),
		events:     make([]yaml_event_t, 0, initial_queue_size),
	***REMOVED***
***REMOVED***

// Destroy an emitter object.
func yaml_emitter_delete(emitter *yaml_emitter_t) ***REMOVED***
	*emitter = yaml_emitter_t***REMOVED******REMOVED***
***REMOVED***

// String write handler.
func yaml_string_write_handler(emitter *yaml_emitter_t, buffer []byte) error ***REMOVED***
	*emitter.output_buffer = append(*emitter.output_buffer, buffer...)
	return nil
***REMOVED***

// yaml_writer_write_handler uses emitter.output_writer to write the
// emitted text.
func yaml_writer_write_handler(emitter *yaml_emitter_t, buffer []byte) error ***REMOVED***
	_, err := emitter.output_writer.Write(buffer)
	return err
***REMOVED***

// Set a string output.
func yaml_emitter_set_output_string(emitter *yaml_emitter_t, output_buffer *[]byte) ***REMOVED***
	if emitter.write_handler != nil ***REMOVED***
		panic("must set the output target only once")
	***REMOVED***
	emitter.write_handler = yaml_string_write_handler
	emitter.output_buffer = output_buffer
***REMOVED***

// Set a file output.
func yaml_emitter_set_output_writer(emitter *yaml_emitter_t, w io.Writer) ***REMOVED***
	if emitter.write_handler != nil ***REMOVED***
		panic("must set the output target only once")
	***REMOVED***
	emitter.write_handler = yaml_writer_write_handler
	emitter.output_writer = w
***REMOVED***

// Set the output encoding.
func yaml_emitter_set_encoding(emitter *yaml_emitter_t, encoding yaml_encoding_t) ***REMOVED***
	if emitter.encoding != yaml_ANY_ENCODING ***REMOVED***
		panic("must set the output encoding only once")
	***REMOVED***
	emitter.encoding = encoding
***REMOVED***

// Set the canonical output style.
func yaml_emitter_set_canonical(emitter *yaml_emitter_t, canonical bool) ***REMOVED***
	emitter.canonical = canonical
***REMOVED***

// Set the indentation increment.
func yaml_emitter_set_indent(emitter *yaml_emitter_t, indent int) ***REMOVED***
	if indent < 2 || indent > 9 ***REMOVED***
		indent = 2
	***REMOVED***
	emitter.best_indent = indent
***REMOVED***

// Set the preferred line width.
func yaml_emitter_set_width(emitter *yaml_emitter_t, width int) ***REMOVED***
	if width < 0 ***REMOVED***
		width = -1
	***REMOVED***
	emitter.best_width = width
***REMOVED***

// Set if unescaped non-ASCII characters are allowed.
func yaml_emitter_set_unicode(emitter *yaml_emitter_t, unicode bool) ***REMOVED***
	emitter.unicode = unicode
***REMOVED***

// Set the preferred line break character.
func yaml_emitter_set_break(emitter *yaml_emitter_t, line_break yaml_break_t) ***REMOVED***
	emitter.line_break = line_break
***REMOVED***

///*
// * Destroy a token object.
// */
//
//YAML_DECLARE(void)
//yaml_token_delete(yaml_token_t *token)
//***REMOVED***
//    assert(token);  // Non-NULL token object expected.
//
//    switch (token.type)
//    ***REMOVED***
//        case YAML_TAG_DIRECTIVE_TOKEN:
//            yaml_free(token.data.tag_directive.handle);
//            yaml_free(token.data.tag_directive.prefix);
//            break;
//
//        case YAML_ALIAS_TOKEN:
//            yaml_free(token.data.alias.value);
//            break;
//
//        case YAML_ANCHOR_TOKEN:
//            yaml_free(token.data.anchor.value);
//            break;
//
//        case YAML_TAG_TOKEN:
//            yaml_free(token.data.tag.handle);
//            yaml_free(token.data.tag.suffix);
//            break;
//
//        case YAML_SCALAR_TOKEN:
//            yaml_free(token.data.scalar.value);
//            break;
//
//        default:
//            break;
//    ***REMOVED***
//
//    memset(token, 0, sizeof(yaml_token_t));
//***REMOVED***
//
///*
// * Check if a string is a valid UTF-8 sequence.
// *
// * Check 'reader.c' for more details on UTF-8 encoding.
// */
//
//static int
//yaml_check_utf8(yaml_char_t *start, size_t length)
//***REMOVED***
//    yaml_char_t *end = start+length;
//    yaml_char_t *pointer = start;
//
//    while (pointer < end) ***REMOVED***
//        unsigned char octet;
//        unsigned int width;
//        unsigned int value;
//        size_t k;
//
//        octet = pointer[0];
//        width = (octet & 0x80) == 0x00 ? 1 :
//                (octet & 0xE0) == 0xC0 ? 2 :
//                (octet & 0xF0) == 0xE0 ? 3 :
//                (octet & 0xF8) == 0xF0 ? 4 : 0;
//        value = (octet & 0x80) == 0x00 ? octet & 0x7F :
//                (octet & 0xE0) == 0xC0 ? octet & 0x1F :
//                (octet & 0xF0) == 0xE0 ? octet & 0x0F :
//                (octet & 0xF8) == 0xF0 ? octet & 0x07 : 0;
//        if (!width) return 0;
//        if (pointer+width > end) return 0;
//        for (k = 1; k < width; k ++) ***REMOVED***
//            octet = pointer[k];
//            if ((octet & 0xC0) != 0x80) return 0;
//            value = (value << 6) + (octet & 0x3F);
//        ***REMOVED***
//        if (!((width == 1) ||
//            (width == 2 && value >= 0x80) ||
//            (width == 3 && value >= 0x800) ||
//            (width == 4 && value >= 0x10000))) return 0;
//
//        pointer += width;
//    ***REMOVED***
//
//    return 1;
//***REMOVED***
//

// Create STREAM-START.
func yaml_stream_start_event_initialize(event *yaml_event_t, encoding yaml_encoding_t) ***REMOVED***
	*event = yaml_event_t***REMOVED***
		typ:      yaml_STREAM_START_EVENT,
		encoding: encoding,
	***REMOVED***
***REMOVED***

// Create STREAM-END.
func yaml_stream_end_event_initialize(event *yaml_event_t) ***REMOVED***
	*event = yaml_event_t***REMOVED***
		typ: yaml_STREAM_END_EVENT,
	***REMOVED***
***REMOVED***

// Create DOCUMENT-START.
func yaml_document_start_event_initialize(
	event *yaml_event_t,
	version_directive *yaml_version_directive_t,
	tag_directives []yaml_tag_directive_t,
	implicit bool,
) ***REMOVED***
	*event = yaml_event_t***REMOVED***
		typ:               yaml_DOCUMENT_START_EVENT,
		version_directive: version_directive,
		tag_directives:    tag_directives,
		implicit:          implicit,
	***REMOVED***
***REMOVED***

// Create DOCUMENT-END.
func yaml_document_end_event_initialize(event *yaml_event_t, implicit bool) ***REMOVED***
	*event = yaml_event_t***REMOVED***
		typ:      yaml_DOCUMENT_END_EVENT,
		implicit: implicit,
	***REMOVED***
***REMOVED***

// Create ALIAS.
func yaml_alias_event_initialize(event *yaml_event_t, anchor []byte) bool ***REMOVED***
	*event = yaml_event_t***REMOVED***
		typ:    yaml_ALIAS_EVENT,
		anchor: anchor,
	***REMOVED***
	return true
***REMOVED***

// Create SCALAR.
func yaml_scalar_event_initialize(event *yaml_event_t, anchor, tag, value []byte, plain_implicit, quoted_implicit bool, style yaml_scalar_style_t) bool ***REMOVED***
	*event = yaml_event_t***REMOVED***
		typ:             yaml_SCALAR_EVENT,
		anchor:          anchor,
		tag:             tag,
		value:           value,
		implicit:        plain_implicit,
		quoted_implicit: quoted_implicit,
		style:           yaml_style_t(style),
	***REMOVED***
	return true
***REMOVED***

// Create SEQUENCE-START.
func yaml_sequence_start_event_initialize(event *yaml_event_t, anchor, tag []byte, implicit bool, style yaml_sequence_style_t) bool ***REMOVED***
	*event = yaml_event_t***REMOVED***
		typ:      yaml_SEQUENCE_START_EVENT,
		anchor:   anchor,
		tag:      tag,
		implicit: implicit,
		style:    yaml_style_t(style),
	***REMOVED***
	return true
***REMOVED***

// Create SEQUENCE-END.
func yaml_sequence_end_event_initialize(event *yaml_event_t) bool ***REMOVED***
	*event = yaml_event_t***REMOVED***
		typ: yaml_SEQUENCE_END_EVENT,
	***REMOVED***
	return true
***REMOVED***

// Create MAPPING-START.
func yaml_mapping_start_event_initialize(event *yaml_event_t, anchor, tag []byte, implicit bool, style yaml_mapping_style_t) ***REMOVED***
	*event = yaml_event_t***REMOVED***
		typ:      yaml_MAPPING_START_EVENT,
		anchor:   anchor,
		tag:      tag,
		implicit: implicit,
		style:    yaml_style_t(style),
	***REMOVED***
***REMOVED***

// Create MAPPING-END.
func yaml_mapping_end_event_initialize(event *yaml_event_t) ***REMOVED***
	*event = yaml_event_t***REMOVED***
		typ: yaml_MAPPING_END_EVENT,
	***REMOVED***
***REMOVED***

// Destroy an event object.
func yaml_event_delete(event *yaml_event_t) ***REMOVED***
	*event = yaml_event_t***REMOVED******REMOVED***
***REMOVED***

///*
// * Create a document object.
// */
//
//YAML_DECLARE(int)
//yaml_document_initialize(document *yaml_document_t,
//        version_directive *yaml_version_directive_t,
//        tag_directives_start *yaml_tag_directive_t,
//        tag_directives_end *yaml_tag_directive_t,
//        start_implicit int, end_implicit int)
//***REMOVED***
//    struct ***REMOVED***
//        error yaml_error_type_t
//    ***REMOVED*** context
//    struct ***REMOVED***
//        start *yaml_node_t
//        end *yaml_node_t
//        top *yaml_node_t
//    ***REMOVED*** nodes = ***REMOVED*** NULL, NULL, NULL ***REMOVED***
//    version_directive_copy *yaml_version_directive_t = NULL
//    struct ***REMOVED***
//        start *yaml_tag_directive_t
//        end *yaml_tag_directive_t
//        top *yaml_tag_directive_t
//    ***REMOVED*** tag_directives_copy = ***REMOVED*** NULL, NULL, NULL ***REMOVED***
//    value yaml_tag_directive_t = ***REMOVED*** NULL, NULL ***REMOVED***
//    mark yaml_mark_t = ***REMOVED*** 0, 0, 0 ***REMOVED***
//
//    assert(document) // Non-NULL document object is expected.
//    assert((tag_directives_start && tag_directives_end) ||
//            (tag_directives_start == tag_directives_end))
//                            // Valid tag directives are expected.
//
//    if (!STACK_INIT(&context, nodes, INITIAL_STACK_SIZE)) goto error
//
//    if (version_directive) ***REMOVED***
//        version_directive_copy = yaml_malloc(sizeof(yaml_version_directive_t))
//        if (!version_directive_copy) goto error
//        version_directive_copy.major = version_directive.major
//        version_directive_copy.minor = version_directive.minor
//    ***REMOVED***
//
//    if (tag_directives_start != tag_directives_end) ***REMOVED***
//        tag_directive *yaml_tag_directive_t
//        if (!STACK_INIT(&context, tag_directives_copy, INITIAL_STACK_SIZE))
//            goto error
//        for (tag_directive = tag_directives_start
//                tag_directive != tag_directives_end; tag_directive ++) ***REMOVED***
//            assert(tag_directive.handle)
//            assert(tag_directive.prefix)
//            if (!yaml_check_utf8(tag_directive.handle,
//                        strlen((char *)tag_directive.handle)))
//                goto error
//            if (!yaml_check_utf8(tag_directive.prefix,
//                        strlen((char *)tag_directive.prefix)))
//                goto error
//            value.handle = yaml_strdup(tag_directive.handle)
//            value.prefix = yaml_strdup(tag_directive.prefix)
//            if (!value.handle || !value.prefix) goto error
//            if (!PUSH(&context, tag_directives_copy, value))
//                goto error
//            value.handle = NULL
//            value.prefix = NULL
//        ***REMOVED***
//    ***REMOVED***
//
//    DOCUMENT_INIT(*document, nodes.start, nodes.end, version_directive_copy,
//            tag_directives_copy.start, tag_directives_copy.top,
//            start_implicit, end_implicit, mark, mark)
//
//    return 1
//
//error:
//    STACK_DEL(&context, nodes)
//    yaml_free(version_directive_copy)
//    while (!STACK_EMPTY(&context, tag_directives_copy)) ***REMOVED***
//        value yaml_tag_directive_t = POP(&context, tag_directives_copy)
//        yaml_free(value.handle)
//        yaml_free(value.prefix)
//    ***REMOVED***
//    STACK_DEL(&context, tag_directives_copy)
//    yaml_free(value.handle)
//    yaml_free(value.prefix)
//
//    return 0
//***REMOVED***
//
///*
// * Destroy a document object.
// */
//
//YAML_DECLARE(void)
//yaml_document_delete(document *yaml_document_t)
//***REMOVED***
//    struct ***REMOVED***
//        error yaml_error_type_t
//    ***REMOVED*** context
//    tag_directive *yaml_tag_directive_t
//
//    context.error = YAML_NO_ERROR // Eliminate a compiler warning.
//
//    assert(document) // Non-NULL document object is expected.
//
//    while (!STACK_EMPTY(&context, document.nodes)) ***REMOVED***
//        node yaml_node_t = POP(&context, document.nodes)
//        yaml_free(node.tag)
//        switch (node.type) ***REMOVED***
//            case YAML_SCALAR_NODE:
//                yaml_free(node.data.scalar.value)
//                break
//            case YAML_SEQUENCE_NODE:
//                STACK_DEL(&context, node.data.sequence.items)
//                break
//            case YAML_MAPPING_NODE:
//                STACK_DEL(&context, node.data.mapping.pairs)
//                break
//            default:
//                assert(0) // Should not happen.
//        ***REMOVED***
//    ***REMOVED***
//    STACK_DEL(&context, document.nodes)
//
//    yaml_free(document.version_directive)
//    for (tag_directive = document.tag_directives.start
//            tag_directive != document.tag_directives.end
//            tag_directive++) ***REMOVED***
//        yaml_free(tag_directive.handle)
//        yaml_free(tag_directive.prefix)
//    ***REMOVED***
//    yaml_free(document.tag_directives.start)
//
//    memset(document, 0, sizeof(yaml_document_t))
//***REMOVED***
//
///**
// * Get a document node.
// */
//
//YAML_DECLARE(yaml_node_t *)
//yaml_document_get_node(document *yaml_document_t, index int)
//***REMOVED***
//    assert(document) // Non-NULL document object is expected.
//
//    if (index > 0 && document.nodes.start + index <= document.nodes.top) ***REMOVED***
//        return document.nodes.start + index - 1
//    ***REMOVED***
//    return NULL
//***REMOVED***
//
///**
// * Get the root object.
// */
//
//YAML_DECLARE(yaml_node_t *)
//yaml_document_get_root_node(document *yaml_document_t)
//***REMOVED***
//    assert(document) // Non-NULL document object is expected.
//
//    if (document.nodes.top != document.nodes.start) ***REMOVED***
//        return document.nodes.start
//    ***REMOVED***
//    return NULL
//***REMOVED***
//
///*
// * Add a scalar node to a document.
// */
//
//YAML_DECLARE(int)
//yaml_document_add_scalar(document *yaml_document_t,
//        tag *yaml_char_t, value *yaml_char_t, length int,
//        style yaml_scalar_style_t)
//***REMOVED***
//    struct ***REMOVED***
//        error yaml_error_type_t
//    ***REMOVED*** context
//    mark yaml_mark_t = ***REMOVED*** 0, 0, 0 ***REMOVED***
//    tag_copy *yaml_char_t = NULL
//    value_copy *yaml_char_t = NULL
//    node yaml_node_t
//
//    assert(document) // Non-NULL document object is expected.
//    assert(value) // Non-NULL value is expected.
//
//    if (!tag) ***REMOVED***
//        tag = (yaml_char_t *)YAML_DEFAULT_SCALAR_TAG
//    ***REMOVED***
//
//    if (!yaml_check_utf8(tag, strlen((char *)tag))) goto error
//    tag_copy = yaml_strdup(tag)
//    if (!tag_copy) goto error
//
//    if (length < 0) ***REMOVED***
//        length = strlen((char *)value)
//    ***REMOVED***
//
//    if (!yaml_check_utf8(value, length)) goto error
//    value_copy = yaml_malloc(length+1)
//    if (!value_copy) goto error
//    memcpy(value_copy, value, length)
//    value_copy[length] = '\0'
//
//    SCALAR_NODE_INIT(node, tag_copy, value_copy, length, style, mark, mark)
//    if (!PUSH(&context, document.nodes, node)) goto error
//
//    return document.nodes.top - document.nodes.start
//
//error:
//    yaml_free(tag_copy)
//    yaml_free(value_copy)
//
//    return 0
//***REMOVED***
//
///*
// * Add a sequence node to a document.
// */
//
//YAML_DECLARE(int)
//yaml_document_add_sequence(document *yaml_document_t,
//        tag *yaml_char_t, style yaml_sequence_style_t)
//***REMOVED***
//    struct ***REMOVED***
//        error yaml_error_type_t
//    ***REMOVED*** context
//    mark yaml_mark_t = ***REMOVED*** 0, 0, 0 ***REMOVED***
//    tag_copy *yaml_char_t = NULL
//    struct ***REMOVED***
//        start *yaml_node_item_t
//        end *yaml_node_item_t
//        top *yaml_node_item_t
//    ***REMOVED*** items = ***REMOVED*** NULL, NULL, NULL ***REMOVED***
//    node yaml_node_t
//
//    assert(document) // Non-NULL document object is expected.
//
//    if (!tag) ***REMOVED***
//        tag = (yaml_char_t *)YAML_DEFAULT_SEQUENCE_TAG
//    ***REMOVED***
//
//    if (!yaml_check_utf8(tag, strlen((char *)tag))) goto error
//    tag_copy = yaml_strdup(tag)
//    if (!tag_copy) goto error
//
//    if (!STACK_INIT(&context, items, INITIAL_STACK_SIZE)) goto error
//
//    SEQUENCE_NODE_INIT(node, tag_copy, items.start, items.end,
//            style, mark, mark)
//    if (!PUSH(&context, document.nodes, node)) goto error
//
//    return document.nodes.top - document.nodes.start
//
//error:
//    STACK_DEL(&context, items)
//    yaml_free(tag_copy)
//
//    return 0
//***REMOVED***
//
///*
// * Add a mapping node to a document.
// */
//
//YAML_DECLARE(int)
//yaml_document_add_mapping(document *yaml_document_t,
//        tag *yaml_char_t, style yaml_mapping_style_t)
//***REMOVED***
//    struct ***REMOVED***
//        error yaml_error_type_t
//    ***REMOVED*** context
//    mark yaml_mark_t = ***REMOVED*** 0, 0, 0 ***REMOVED***
//    tag_copy *yaml_char_t = NULL
//    struct ***REMOVED***
//        start *yaml_node_pair_t
//        end *yaml_node_pair_t
//        top *yaml_node_pair_t
//    ***REMOVED*** pairs = ***REMOVED*** NULL, NULL, NULL ***REMOVED***
//    node yaml_node_t
//
//    assert(document) // Non-NULL document object is expected.
//
//    if (!tag) ***REMOVED***
//        tag = (yaml_char_t *)YAML_DEFAULT_MAPPING_TAG
//    ***REMOVED***
//
//    if (!yaml_check_utf8(tag, strlen((char *)tag))) goto error
//    tag_copy = yaml_strdup(tag)
//    if (!tag_copy) goto error
//
//    if (!STACK_INIT(&context, pairs, INITIAL_STACK_SIZE)) goto error
//
//    MAPPING_NODE_INIT(node, tag_copy, pairs.start, pairs.end,
//            style, mark, mark)
//    if (!PUSH(&context, document.nodes, node)) goto error
//
//    return document.nodes.top - document.nodes.start
//
//error:
//    STACK_DEL(&context, pairs)
//    yaml_free(tag_copy)
//
//    return 0
//***REMOVED***
//
///*
// * Append an item to a sequence node.
// */
//
//YAML_DECLARE(int)
//yaml_document_append_sequence_item(document *yaml_document_t,
//        sequence int, item int)
//***REMOVED***
//    struct ***REMOVED***
//        error yaml_error_type_t
//    ***REMOVED*** context
//
//    assert(document) // Non-NULL document is required.
//    assert(sequence > 0
//            && document.nodes.start + sequence <= document.nodes.top)
//                            // Valid sequence id is required.
//    assert(document.nodes.start[sequence-1].type == YAML_SEQUENCE_NODE)
//                            // A sequence node is required.
//    assert(item > 0 && document.nodes.start + item <= document.nodes.top)
//                            // Valid item id is required.
//
//    if (!PUSH(&context,
//                document.nodes.start[sequence-1].data.sequence.items, item))
//        return 0
//
//    return 1
//***REMOVED***
//
///*
// * Append a pair of a key and a value to a mapping node.
// */
//
//YAML_DECLARE(int)
//yaml_document_append_mapping_pair(document *yaml_document_t,
//        mapping int, key int, value int)
//***REMOVED***
//    struct ***REMOVED***
//        error yaml_error_type_t
//    ***REMOVED*** context
//
//    pair yaml_node_pair_t
//
//    assert(document) // Non-NULL document is required.
//    assert(mapping > 0
//            && document.nodes.start + mapping <= document.nodes.top)
//                            // Valid mapping id is required.
//    assert(document.nodes.start[mapping-1].type == YAML_MAPPING_NODE)
//                            // A mapping node is required.
//    assert(key > 0 && document.nodes.start + key <= document.nodes.top)
//                            // Valid key id is required.
//    assert(value > 0 && document.nodes.start + value <= document.nodes.top)
//                            // Valid value id is required.
//
//    pair.key = key
//    pair.value = value
//
//    if (!PUSH(&context,
//                document.nodes.start[mapping-1].data.mapping.pairs, pair))
//        return 0
//
//    return 1
//***REMOVED***
//
//