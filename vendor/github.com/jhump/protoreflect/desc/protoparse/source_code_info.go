package protoparse

import (
	"bytes"
	"strings"

	"github.com/golang/protobuf/proto"
	dpb "github.com/golang/protobuf/protoc-gen-go/descriptor"

	"github.com/jhump/protoreflect/desc/internal"
)

func (r *parseResult) generateSourceCodeInfo() *dpb.SourceCodeInfo ***REMOVED***
	if r.nodes == nil ***REMOVED***
		// skip files that do not have AST info (these will be files
		// that came from well-known descriptors, instead of from source)
		return nil
	***REMOVED***

	sci := sourceCodeInfo***REMOVED***commentsUsed: map[*comment]struct***REMOVED******REMOVED******REMOVED******REMOVED******REMOVED***
	path := make([]int32, 0, 10)

	fn := r.getFileNode(r.fd).(*fileNode)
	sci.newLocWithoutComments(fn, nil)

	if fn.syntax != nil ***REMOVED***
		sci.newLoc(fn.syntax, append(path, internal.File_syntaxTag))
	***REMOVED***

	var depIndex, optIndex, msgIndex, enumIndex, extendIndex, svcIndex int32

	for _, child := range fn.decls ***REMOVED***
		switch ***REMOVED***
		case child.imp != nil:
			sci.newLoc(child.imp, append(path, internal.File_dependencyTag, int32(depIndex)))
			depIndex++
		case child.pkg != nil:
			sci.newLoc(child.pkg, append(path, internal.File_packageTag))
		case child.option != nil:
			r.generateSourceCodeInfoForOption(&sci, child.option, false, &optIndex, append(path, internal.File_optionsTag))
		case child.message != nil:
			r.generateSourceCodeInfoForMessage(&sci, child.message, nil, append(path, internal.File_messagesTag, msgIndex))
			msgIndex++
		case child.enum != nil:
			r.generateSourceCodeInfoForEnum(&sci, child.enum, append(path, internal.File_enumsTag, enumIndex))
			enumIndex++
		case child.extend != nil:
			r.generateSourceCodeInfoForExtensions(&sci, child.extend, &extendIndex, &msgIndex, append(path, internal.File_extensionsTag), append(dup(path), internal.File_messagesTag))
		case child.service != nil:
			r.generateSourceCodeInfoForService(&sci, child.service, append(path, internal.File_servicesTag, svcIndex))
			svcIndex++
		***REMOVED***
	***REMOVED***

	return &dpb.SourceCodeInfo***REMOVED***Location: sci.locs***REMOVED***
***REMOVED***

func (r *parseResult) generateSourceCodeInfoForOption(sci *sourceCodeInfo, n *optionNode, compact bool, uninterpIndex *int32, path []int32) ***REMOVED***
	if !compact ***REMOVED***
		sci.newLocWithoutComments(n, path)
	***REMOVED***
	subPath := r.interpretedOptions[n]
	if len(subPath) > 0 ***REMOVED***
		p := path
		if subPath[0] == -1 ***REMOVED***
			// used by "default" and "json_name" field pseudo-options
			// to attribute path to parent element (since those are
			// stored directly on the descriptor, not its options)
			p = make([]int32, len(path)-1)
			copy(p, path)
			subPath = subPath[1:]
		***REMOVED***
		sci.newLoc(n, append(p, subPath...))
		return
	***REMOVED***

	// it's an uninterpreted option
	optPath := append(path, internal.UninterpretedOptionsTag, *uninterpIndex)
	*uninterpIndex++
	sci.newLoc(n, optPath)
	var valTag int32
	switch n.val.(type) ***REMOVED***
	case *compoundIdentNode:
		valTag = internal.Uninterpreted_identTag
	case *intLiteralNode:
		valTag = internal.Uninterpreted_posIntTag
	case *compoundIntNode:
		valTag = internal.Uninterpreted_negIntTag
	case *compoundFloatNode:
		valTag = internal.Uninterpreted_doubleTag
	case *compoundStringNode:
		valTag = internal.Uninterpreted_stringTag
	case *aggregateLiteralNode:
		valTag = internal.Uninterpreted_aggregateTag
	***REMOVED***
	if valTag != 0 ***REMOVED***
		sci.newLoc(n.val, append(optPath, valTag))
	***REMOVED***
	for j, nn := range n.name.parts ***REMOVED***
		optNmPath := append(optPath, internal.Uninterpreted_nameTag, int32(j))
		sci.newLoc(nn, optNmPath)
		sci.newLoc(nn.text, append(optNmPath, internal.UninterpretedName_nameTag))
	***REMOVED***
***REMOVED***

func (r *parseResult) generateSourceCodeInfoForMessage(sci *sourceCodeInfo, n msgDecl, fieldPath []int32, path []int32) ***REMOVED***
	sci.newLoc(n, path)

	var decls []*messageElement
	switch n := n.(type) ***REMOVED***
	case *messageNode:
		decls = n.decls
	case *groupNode:
		decls = n.decls
	case *mapFieldNode:
		// map entry so nothing else to do
		return
	***REMOVED***

	sci.newLoc(n.messageName(), append(path, internal.Message_nameTag))
	// matching protoc, which emits the corresponding field type name (for group fields)
	// right after the source location for the group message name
	if fieldPath != nil ***REMOVED***
		sci.newLoc(n.messageName(), append(fieldPath, internal.Field_typeNameTag))
	***REMOVED***

	var optIndex, fieldIndex, oneOfIndex, extendIndex, nestedMsgIndex int32
	var nestedEnumIndex, extRangeIndex, reservedRangeIndex, reservedNameIndex int32
	for _, child := range decls ***REMOVED***
		switch ***REMOVED***
		case child.option != nil:
			r.generateSourceCodeInfoForOption(sci, child.option, false, &optIndex, append(path, internal.Message_optionsTag))
		case child.field != nil:
			r.generateSourceCodeInfoForField(sci, child.field, append(path, internal.Message_fieldsTag, fieldIndex))
			fieldIndex++
		case child.group != nil:
			fldPath := append(path, internal.Message_fieldsTag, fieldIndex)
			r.generateSourceCodeInfoForField(sci, child.group, fldPath)
			fieldIndex++
			r.generateSourceCodeInfoForMessage(sci, child.group, fldPath, append(dup(path), internal.Message_nestedMessagesTag, nestedMsgIndex))
			nestedMsgIndex++
		case child.mapField != nil:
			r.generateSourceCodeInfoForField(sci, child.mapField, append(path, internal.Message_fieldsTag, fieldIndex))
			fieldIndex++
		case child.oneOf != nil:
			r.generateSourceCodeInfoForOneOf(sci, child.oneOf, &fieldIndex, &nestedMsgIndex, append(path, internal.Message_fieldsTag), append(dup(path), internal.Message_nestedMessagesTag), append(dup(path), internal.Message_oneOfsTag, oneOfIndex))
			oneOfIndex++
		case child.nested != nil:
			r.generateSourceCodeInfoForMessage(sci, child.nested, nil, append(path, internal.Message_nestedMessagesTag, nestedMsgIndex))
			nestedMsgIndex++
		case child.enum != nil:
			r.generateSourceCodeInfoForEnum(sci, child.enum, append(path, internal.Message_enumsTag, nestedEnumIndex))
			nestedEnumIndex++
		case child.extend != nil:
			r.generateSourceCodeInfoForExtensions(sci, child.extend, &extendIndex, &nestedMsgIndex, append(path, internal.Message_extensionsTag), append(dup(path), internal.Message_nestedMessagesTag))
		case child.extensionRange != nil:
			r.generateSourceCodeInfoForExtensionRanges(sci, child.extensionRange, &extRangeIndex, append(path, internal.Message_extensionRangeTag))
		case child.reserved != nil:
			if len(child.reserved.names) > 0 ***REMOVED***
				resPath := append(path, internal.Message_reservedNameTag)
				sci.newLoc(child.reserved, resPath)
				for _, rn := range child.reserved.names ***REMOVED***
					sci.newLoc(rn, append(resPath, reservedNameIndex))
					reservedNameIndex++
				***REMOVED***
			***REMOVED***
			if len(child.reserved.ranges) > 0 ***REMOVED***
				resPath := append(path, internal.Message_reservedRangeTag)
				sci.newLoc(child.reserved, resPath)
				for _, rr := range child.reserved.ranges ***REMOVED***
					r.generateSourceCodeInfoForReservedRange(sci, rr, append(resPath, reservedRangeIndex))
					reservedRangeIndex++
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func (r *parseResult) generateSourceCodeInfoForEnum(sci *sourceCodeInfo, n *enumNode, path []int32) ***REMOVED***
	sci.newLoc(n, path)
	sci.newLoc(n.name, append(path, internal.Enum_nameTag))

	var optIndex, valIndex, reservedNameIndex, reservedRangeIndex int32
	for _, child := range n.decls ***REMOVED***
		switch ***REMOVED***
		case child.option != nil:
			r.generateSourceCodeInfoForOption(sci, child.option, false, &optIndex, append(path, internal.Enum_optionsTag))
		case child.value != nil:
			r.generateSourceCodeInfoForEnumValue(sci, child.value, append(path, internal.Enum_valuesTag, valIndex))
			valIndex++
		case child.reserved != nil:
			if len(child.reserved.names) > 0 ***REMOVED***
				resPath := append(path, internal.Enum_reservedNameTag)
				sci.newLoc(child.reserved, resPath)
				for _, rn := range child.reserved.names ***REMOVED***
					sci.newLoc(rn, append(resPath, reservedNameIndex))
					reservedNameIndex++
				***REMOVED***
			***REMOVED***
			if len(child.reserved.ranges) > 0 ***REMOVED***
				resPath := append(path, internal.Enum_reservedRangeTag)
				sci.newLoc(child.reserved, resPath)
				for _, rr := range child.reserved.ranges ***REMOVED***
					r.generateSourceCodeInfoForReservedRange(sci, rr, append(resPath, reservedRangeIndex))
					reservedRangeIndex++
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func (r *parseResult) generateSourceCodeInfoForEnumValue(sci *sourceCodeInfo, n *enumValueNode, path []int32) ***REMOVED***
	sci.newLoc(n, path)
	sci.newLoc(n.name, append(path, internal.EnumVal_nameTag))
	sci.newLoc(n.getNumber(), append(path, internal.EnumVal_numberTag))

	// enum value options
	if n.options != nil ***REMOVED***
		optsPath := append(path, internal.EnumVal_optionsTag)
		sci.newLoc(n.options, optsPath)
		var optIndex int32
		for _, opt := range n.options.decls ***REMOVED***
			r.generateSourceCodeInfoForOption(sci, opt, true, &optIndex, optsPath)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (r *parseResult) generateSourceCodeInfoForReservedRange(sci *sourceCodeInfo, n *rangeNode, path []int32) ***REMOVED***
	sci.newLoc(n, path)
	sci.newLoc(n.startNode, append(path, internal.ReservedRange_startTag))
	if n.endNode != nil ***REMOVED***
		sci.newLoc(n.endNode, append(path, internal.ReservedRange_endTag))
	***REMOVED***
***REMOVED***

func (r *parseResult) generateSourceCodeInfoForExtensions(sci *sourceCodeInfo, n *extendNode, extendIndex, msgIndex *int32, extendPath, msgPath []int32) ***REMOVED***
	sci.newLoc(n, extendPath)
	for _, decl := range n.decls ***REMOVED***
		switch ***REMOVED***
		case decl.field != nil:
			r.generateSourceCodeInfoForField(sci, decl.field, append(extendPath, *extendIndex))
			*extendIndex++
		case decl.group != nil:
			fldPath := append(extendPath, *extendIndex)
			r.generateSourceCodeInfoForField(sci, decl.group, fldPath)
			*extendIndex++
			r.generateSourceCodeInfoForMessage(sci, decl.group, fldPath, append(msgPath, *msgIndex))
			*msgIndex++
		***REMOVED***
	***REMOVED***
***REMOVED***

func (r *parseResult) generateSourceCodeInfoForOneOf(sci *sourceCodeInfo, n *oneOfNode, fieldIndex, nestedMsgIndex *int32, fieldPath, nestedMsgPath, oneOfPath []int32) ***REMOVED***
	sci.newLoc(n, oneOfPath)
	sci.newLoc(n.name, append(oneOfPath, internal.OneOf_nameTag))

	var optIndex int32
	for _, child := range n.decls ***REMOVED***
		switch ***REMOVED***
		case child.option != nil:
			r.generateSourceCodeInfoForOption(sci, child.option, false, &optIndex, append(oneOfPath, internal.OneOf_optionsTag))
		case child.field != nil:
			r.generateSourceCodeInfoForField(sci, child.field, append(fieldPath, *fieldIndex))
			*fieldIndex++
		case child.group != nil:
			fldPath := append(fieldPath, *fieldIndex)
			r.generateSourceCodeInfoForField(sci, child.group, fldPath)
			*fieldIndex++
			r.generateSourceCodeInfoForMessage(sci, child.group, fldPath, append(nestedMsgPath, *nestedMsgIndex))
			*nestedMsgIndex++
		***REMOVED***
	***REMOVED***
***REMOVED***

func (r *parseResult) generateSourceCodeInfoForField(sci *sourceCodeInfo, n fieldDecl, path []int32) ***REMOVED***
	isGroup := false
	var opts *compactOptionsNode
	var extendee *extendNode
	var fieldType string
	switch n := n.(type) ***REMOVED***
	case *fieldNode:
		opts = n.options
		extendee = n.extendee
		fieldType = n.fldType.val
	case *mapFieldNode:
		opts = n.options
	case *groupNode:
		isGroup = true
		extendee = n.extendee
	case *syntheticMapField:
		// shouldn't get here since we don't recurse into fields from a mapNode
		// in generateSourceCodeInfoForMessage... but just in case
		return
	***REMOVED***

	if isGroup ***REMOVED***
		// comments will appear on group message
		sci.newLocWithoutComments(n, path)
		if extendee != nil ***REMOVED***
			sci.newLoc(extendee.extendee, append(path, internal.Field_extendeeTag))
		***REMOVED***
		if n.fieldLabel() != nil ***REMOVED***
			// no comments here either (label is first token for group, so we want
			// to leave the comments to be associated with the group message instead)
			sci.newLocWithoutComments(n.fieldLabel(), append(path, internal.Field_labelTag))
		***REMOVED***
		sci.newLoc(n.fieldType(), append(path, internal.Field_typeTag))
		// let the name comments be attributed to the group name
		sci.newLocWithoutComments(n.fieldName(), append(path, internal.Field_nameTag))
	***REMOVED*** else ***REMOVED***
		sci.newLoc(n, path)
		if extendee != nil ***REMOVED***
			sci.newLoc(extendee.extendee, append(path, internal.Field_extendeeTag))
		***REMOVED***
		if n.fieldLabel() != nil ***REMOVED***
			sci.newLoc(n.fieldLabel(), append(path, internal.Field_labelTag))
		***REMOVED***
		n.fieldType()
		var tag int32
		if _, isScalar := fieldTypes[fieldType]; isScalar ***REMOVED***
			tag = internal.Field_typeTag
		***REMOVED*** else ***REMOVED***
			// this is a message or an enum, so attribute type location
			// to the type name field
			tag = internal.Field_typeNameTag
		***REMOVED***
		sci.newLoc(n.fieldType(), append(path, tag))
		sci.newLoc(n.fieldName(), append(path, internal.Field_nameTag))
	***REMOVED***
	sci.newLoc(n.fieldTag(), append(path, internal.Field_numberTag))

	if opts != nil ***REMOVED***
		optsPath := append(path, internal.Field_optionsTag)
		sci.newLoc(opts, optsPath)
		var optIndex int32
		for _, opt := range opts.decls ***REMOVED***
			r.generateSourceCodeInfoForOption(sci, opt, true, &optIndex, optsPath)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (r *parseResult) generateSourceCodeInfoForExtensionRanges(sci *sourceCodeInfo, n *extensionRangeNode, extRangeIndex *int32, path []int32) ***REMOVED***
	sci.newLoc(n, path)
	for _, child := range n.ranges ***REMOVED***
		path := append(path, *extRangeIndex)
		*extRangeIndex++
		sci.newLoc(child, path)
		sci.newLoc(child.startNode, append(path, internal.ExtensionRange_startTag))
		if child.endNode != nil ***REMOVED***
			sci.newLoc(child.endNode, append(path, internal.ExtensionRange_endTag))
		***REMOVED***
		if n.options != nil ***REMOVED***
			optsPath := append(path, internal.ExtensionRange_optionsTag)
			sci.newLoc(n.options, optsPath)
			var optIndex int32
			for _, opt := range n.options.decls ***REMOVED***
				r.generateSourceCodeInfoForOption(sci, opt, true, &optIndex, optsPath)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func (r *parseResult) generateSourceCodeInfoForService(sci *sourceCodeInfo, n *serviceNode, path []int32) ***REMOVED***
	sci.newLoc(n, path)
	sci.newLoc(n.name, append(path, internal.Service_nameTag))
	var optIndex, rpcIndex int32
	for _, child := range n.decls ***REMOVED***
		switch ***REMOVED***
		case child.option != nil:
			r.generateSourceCodeInfoForOption(sci, child.option, false, &optIndex, append(path, internal.Service_optionsTag))
		case child.rpc != nil:
			r.generateSourceCodeInfoForMethod(sci, child.rpc, append(path, internal.Service_methodsTag, rpcIndex))
			rpcIndex++
		***REMOVED***
	***REMOVED***
***REMOVED***

func (r *parseResult) generateSourceCodeInfoForMethod(sci *sourceCodeInfo, n *methodNode, path []int32) ***REMOVED***
	sci.newLoc(n, path)
	sci.newLoc(n.name, append(path, internal.Method_nameTag))
	if n.input.streamKeyword != nil ***REMOVED***
		sci.newLoc(n.input.streamKeyword, append(path, internal.Method_inputStreamTag))
	***REMOVED***
	sci.newLoc(n.input.msgType, append(path, internal.Method_inputTag))
	if n.output.streamKeyword != nil ***REMOVED***
		sci.newLoc(n.output.streamKeyword, append(path, internal.Method_outputStreamTag))
	***REMOVED***
	sci.newLoc(n.output.msgType, append(path, internal.Method_outputTag))

	optsPath := append(path, internal.Method_optionsTag)
	var optIndex int32
	for _, opt := range n.options ***REMOVED***
		r.generateSourceCodeInfoForOption(sci, opt, false, &optIndex, optsPath)
	***REMOVED***
***REMOVED***

type sourceCodeInfo struct ***REMOVED***
	locs         []*dpb.SourceCodeInfo_Location
	commentsUsed map[*comment]struct***REMOVED******REMOVED***
***REMOVED***

func (sci *sourceCodeInfo) newLocWithoutComments(n node, path []int32) ***REMOVED***
	dup := make([]int32, len(path))
	copy(dup, path)
	sci.locs = append(sci.locs, &dpb.SourceCodeInfo_Location***REMOVED***
		Path: dup,
		Span: makeSpan(n.start(), n.end()),
	***REMOVED***)
***REMOVED***

func (sci *sourceCodeInfo) newLoc(n node, path []int32) ***REMOVED***
	leadingComments := n.leadingComments()
	trailingComments := n.trailingComments()
	if sci.commentUsed(leadingComments) ***REMOVED***
		leadingComments = nil
	***REMOVED***
	if sci.commentUsed(trailingComments) ***REMOVED***
		trailingComments = nil
	***REMOVED***
	detached := groupComments(leadingComments)
	var trail *string
	if str, ok := combineComments(trailingComments); ok ***REMOVED***
		trail = proto.String(str)
	***REMOVED***
	var lead *string
	if len(leadingComments) > 0 && leadingComments[len(leadingComments)-1].end.Line >= n.start().Line-1 ***REMOVED***
		lead = proto.String(detached[len(detached)-1])
		detached = detached[:len(detached)-1]
	***REMOVED***
	dup := make([]int32, len(path))
	copy(dup, path)
	sci.locs = append(sci.locs, &dpb.SourceCodeInfo_Location***REMOVED***
		LeadingDetachedComments: detached,
		LeadingComments:         lead,
		TrailingComments:        trail,
		Path:                    dup,
		Span:                    makeSpan(n.start(), n.end()),
	***REMOVED***)
***REMOVED***

func makeSpan(start, end *SourcePos) []int32 ***REMOVED***
	if start.Line == end.Line ***REMOVED***
		return []int32***REMOVED***int32(start.Line) - 1, int32(start.Col) - 1, int32(end.Col) - 1***REMOVED***
	***REMOVED***
	return []int32***REMOVED***int32(start.Line) - 1, int32(start.Col) - 1, int32(end.Line) - 1, int32(end.Col) - 1***REMOVED***
***REMOVED***

func (sci *sourceCodeInfo) commentUsed(c []comment) bool ***REMOVED***
	if len(c) == 0 ***REMOVED***
		return false
	***REMOVED***
	if _, ok := sci.commentsUsed[&c[0]]; ok ***REMOVED***
		return true
	***REMOVED***

	sci.commentsUsed[&c[0]] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	return false
***REMOVED***

func groupComments(comments []comment) []string ***REMOVED***
	if len(comments) == 0 ***REMOVED***
		return nil
	***REMOVED***

	var groups []string
	singleLineStyle := comments[0].text[:2] == "//"
	line := comments[0].end.Line
	start := 0
	for i := 1; i < len(comments); i++ ***REMOVED***
		c := comments[i]
		prevSingleLine := singleLineStyle
		singleLineStyle = strings.HasPrefix(comments[i].text, "//")
		if !singleLineStyle || prevSingleLine != singleLineStyle || c.start.Line > line+1 ***REMOVED***
			// new group!
			if str, ok := combineComments(comments[start:i]); ok ***REMOVED***
				groups = append(groups, str)
			***REMOVED***
			start = i
		***REMOVED***
		line = c.end.Line
	***REMOVED***
	// don't forget last group
	if str, ok := combineComments(comments[start:]); ok ***REMOVED***
		groups = append(groups, str)
	***REMOVED***
	return groups
***REMOVED***

func combineComments(comments []comment) (string, bool) ***REMOVED***
	if len(comments) == 0 ***REMOVED***
		return "", false
	***REMOVED***
	var buf bytes.Buffer
	for _, c := range comments ***REMOVED***
		if c.text[:2] == "//" ***REMOVED***
			buf.WriteString(c.text[2:])
		***REMOVED*** else ***REMOVED***
			lines := strings.Split(c.text[2:len(c.text)-2], "\n")
			first := true
			for _, l := range lines ***REMOVED***
				if first ***REMOVED***
					first = false
				***REMOVED*** else ***REMOVED***
					buf.WriteByte('\n')
				***REMOVED***

				// strip a prefix of whitespace followed by '*'
				j := 0
				for j < len(l) ***REMOVED***
					if l[j] != ' ' && l[j] != '\t' ***REMOVED***
						break
					***REMOVED***
					j++
				***REMOVED***
				if j == len(l) ***REMOVED***
					l = ""
				***REMOVED*** else if l[j] == '*' ***REMOVED***
					l = l[j+1:]
				***REMOVED*** else if j > 0 ***REMOVED***
					l = " " + l[j:]
				***REMOVED***

				buf.WriteString(l)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return buf.String(), true
***REMOVED***

func dup(p []int32) []int32 ***REMOVED***
	return append(([]int32)(nil), p...)
***REMOVED***
