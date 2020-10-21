package protoparse

import (
	"bytes"
	"math"
	"strings"
	"unicode"

	"github.com/golang/protobuf/proto"
	dpb "github.com/golang/protobuf/protoc-gen-go/descriptor"

	"github.com/jhump/protoreflect/desc/internal"
)

func (r *parseResult) createFileDescriptor(filename string, file *fileNode) ***REMOVED***
	fd := &dpb.FileDescriptorProto***REMOVED***Name: proto.String(filename)***REMOVED***
	r.fd = fd
	r.putFileNode(fd, file)

	isProto3 := false
	if file.syntax != nil ***REMOVED***
		if file.syntax.syntax.val == "proto3" ***REMOVED***
			isProto3 = true
		***REMOVED*** else if file.syntax.syntax.val != "proto2" ***REMOVED***
			if r.errs.handleErrorWithPos(file.syntax.syntax.start(), `syntax value must be "proto2" or "proto3"`) != nil ***REMOVED***
				return
			***REMOVED***
		***REMOVED***

		// proto2 is the default, so no need to set unless proto3
		if isProto3 ***REMOVED***
			fd.Syntax = proto.String(file.syntax.syntax.val)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		r.errs.warn(file.start(), ErrNoSyntax)
	***REMOVED***

	for _, decl := range file.decls ***REMOVED***
		if r.errs.err != nil ***REMOVED***
			return
		***REMOVED***
		if decl.enum != nil ***REMOVED***
			fd.EnumType = append(fd.EnumType, r.asEnumDescriptor(decl.enum))
		***REMOVED*** else if decl.extend != nil ***REMOVED***
			r.addExtensions(decl.extend, &fd.Extension, &fd.MessageType, isProto3)
		***REMOVED*** else if decl.imp != nil ***REMOVED***
			file.imports = append(file.imports, decl.imp)
			index := len(fd.Dependency)
			fd.Dependency = append(fd.Dependency, decl.imp.name.val)
			if decl.imp.public ***REMOVED***
				fd.PublicDependency = append(fd.PublicDependency, int32(index))
			***REMOVED*** else if decl.imp.weak ***REMOVED***
				fd.WeakDependency = append(fd.WeakDependency, int32(index))
			***REMOVED***
		***REMOVED*** else if decl.message != nil ***REMOVED***
			fd.MessageType = append(fd.MessageType, r.asMessageDescriptor(decl.message, isProto3))
		***REMOVED*** else if decl.option != nil ***REMOVED***
			if fd.Options == nil ***REMOVED***
				fd.Options = &dpb.FileOptions***REMOVED******REMOVED***
			***REMOVED***
			fd.Options.UninterpretedOption = append(fd.Options.UninterpretedOption, r.asUninterpretedOption(decl.option))
		***REMOVED*** else if decl.service != nil ***REMOVED***
			fd.Service = append(fd.Service, r.asServiceDescriptor(decl.service))
		***REMOVED*** else if decl.pkg != nil ***REMOVED***
			if fd.Package != nil ***REMOVED***
				if r.errs.handleErrorWithPos(decl.pkg.start(), "files should have only one package declaration") != nil ***REMOVED***
					return
				***REMOVED***
			***REMOVED***
			fd.Package = proto.String(decl.pkg.name.val)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (r *parseResult) asUninterpretedOptions(nodes []*optionNode) []*dpb.UninterpretedOption ***REMOVED***
	if len(nodes) == 0 ***REMOVED***
		return nil
	***REMOVED***
	opts := make([]*dpb.UninterpretedOption, len(nodes))
	for i, n := range nodes ***REMOVED***
		opts[i] = r.asUninterpretedOption(n)
	***REMOVED***
	return opts
***REMOVED***

func (r *parseResult) asUninterpretedOption(node *optionNode) *dpb.UninterpretedOption ***REMOVED***
	opt := &dpb.UninterpretedOption***REMOVED***Name: r.asUninterpretedOptionName(node.name.parts)***REMOVED***
	r.putOptionNode(opt, node)

	switch val := node.val.value().(type) ***REMOVED***
	case bool:
		if val ***REMOVED***
			opt.IdentifierValue = proto.String("true")
		***REMOVED*** else ***REMOVED***
			opt.IdentifierValue = proto.String("false")
		***REMOVED***
	case int64:
		opt.NegativeIntValue = proto.Int64(val)
	case uint64:
		opt.PositiveIntValue = proto.Uint64(val)
	case float64:
		opt.DoubleValue = proto.Float64(val)
	case string:
		opt.StringValue = []byte(val)
	case identifier:
		opt.IdentifierValue = proto.String(string(val))
	case []*aggregateEntryNode:
		var buf bytes.Buffer
		aggToString(val, &buf)
		aggStr := buf.String()
		opt.AggregateValue = proto.String(aggStr)
	***REMOVED***
	return opt
***REMOVED***

func (r *parseResult) asUninterpretedOptionName(parts []*optionNamePartNode) []*dpb.UninterpretedOption_NamePart ***REMOVED***
	ret := make([]*dpb.UninterpretedOption_NamePart, len(parts))
	for i, part := range parts ***REMOVED***
		txt := part.text.val
		if !part.isExtension ***REMOVED***
			txt = part.text.val[part.offset : part.offset+part.length]
		***REMOVED***
		np := &dpb.UninterpretedOption_NamePart***REMOVED***
			NamePart:    proto.String(txt),
			IsExtension: proto.Bool(part.isExtension),
		***REMOVED***
		r.putOptionNamePartNode(np, part)
		ret[i] = np
	***REMOVED***
	return ret
***REMOVED***

func (r *parseResult) addExtensions(ext *extendNode, flds *[]*dpb.FieldDescriptorProto, msgs *[]*dpb.DescriptorProto, isProto3 bool) ***REMOVED***
	extendee := ext.extendee.val
	count := 0
	for _, decl := range ext.decls ***REMOVED***
		if decl.field != nil ***REMOVED***
			count++
			decl.field.extendee = ext
			// use higher limit since we don't know yet whether extendee is messageset wire format
			fd := r.asFieldDescriptor(decl.field, internal.MaxTag, isProto3)
			fd.Extendee = proto.String(extendee)
			*flds = append(*flds, fd)
		***REMOVED*** else if decl.group != nil ***REMOVED***
			count++
			decl.group.extendee = ext
			// ditto: use higher limit right now
			fd, md := r.asGroupDescriptors(decl.group, isProto3, internal.MaxTag)
			fd.Extendee = proto.String(extendee)
			*flds = append(*flds, fd)
			*msgs = append(*msgs, md)
		***REMOVED***
	***REMOVED***
	if count == 0 ***REMOVED***
		_ = r.errs.handleErrorWithPos(ext.start(), "extend sections must define at least one extension")
	***REMOVED***
***REMOVED***

func asLabel(lbl *fieldLabel) *dpb.FieldDescriptorProto_Label ***REMOVED***
	if lbl.identNode == nil ***REMOVED***
		return nil
	***REMOVED***
	switch ***REMOVED***
	case lbl.repeated:
		return dpb.FieldDescriptorProto_LABEL_REPEATED.Enum()
	case lbl.required:
		return dpb.FieldDescriptorProto_LABEL_REQUIRED.Enum()
	default:
		return dpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum()
	***REMOVED***
***REMOVED***

func (r *parseResult) asFieldDescriptor(node *fieldNode, maxTag int32, isProto3 bool) *dpb.FieldDescriptorProto ***REMOVED***
	tag := node.tag.val
	if err := checkTag(node.tag.start(), tag, maxTag); err != nil ***REMOVED***
		_ = r.errs.handleError(err)
	***REMOVED***
	fd := newFieldDescriptor(node.name.val, node.fldType.val, int32(tag), asLabel(&node.label))
	r.putFieldNode(fd, node)
	if opts := node.options.Elements(); len(opts) > 0 ***REMOVED***
		fd.Options = &dpb.FieldOptions***REMOVED***UninterpretedOption: r.asUninterpretedOptions(opts)***REMOVED***
	***REMOVED***
	if isProto3 && fd.Label != nil && fd.GetLabel() == dpb.FieldDescriptorProto_LABEL_OPTIONAL ***REMOVED***
		internal.SetProto3Optional(fd)
	***REMOVED***
	return fd
***REMOVED***

var fieldTypes = map[string]dpb.FieldDescriptorProto_Type***REMOVED***
	"double":   dpb.FieldDescriptorProto_TYPE_DOUBLE,
	"float":    dpb.FieldDescriptorProto_TYPE_FLOAT,
	"int32":    dpb.FieldDescriptorProto_TYPE_INT32,
	"int64":    dpb.FieldDescriptorProto_TYPE_INT64,
	"uint32":   dpb.FieldDescriptorProto_TYPE_UINT32,
	"uint64":   dpb.FieldDescriptorProto_TYPE_UINT64,
	"sint32":   dpb.FieldDescriptorProto_TYPE_SINT32,
	"sint64":   dpb.FieldDescriptorProto_TYPE_SINT64,
	"fixed32":  dpb.FieldDescriptorProto_TYPE_FIXED32,
	"fixed64":  dpb.FieldDescriptorProto_TYPE_FIXED64,
	"sfixed32": dpb.FieldDescriptorProto_TYPE_SFIXED32,
	"sfixed64": dpb.FieldDescriptorProto_TYPE_SFIXED64,
	"bool":     dpb.FieldDescriptorProto_TYPE_BOOL,
	"string":   dpb.FieldDescriptorProto_TYPE_STRING,
	"bytes":    dpb.FieldDescriptorProto_TYPE_BYTES,
***REMOVED***

func newFieldDescriptor(name string, fieldType string, tag int32, lbl *dpb.FieldDescriptorProto_Label) *dpb.FieldDescriptorProto ***REMOVED***
	fd := &dpb.FieldDescriptorProto***REMOVED***
		Name:     proto.String(name),
		JsonName: proto.String(internal.JsonName(name)),
		Number:   proto.Int32(tag),
		Label:    lbl,
	***REMOVED***
	t, ok := fieldTypes[fieldType]
	if ok ***REMOVED***
		fd.Type = t.Enum()
	***REMOVED*** else ***REMOVED***
		// NB: we don't have enough info to determine whether this is an enum
		// or a message type, so we'll leave Type nil and set it later
		// (during linking)
		fd.TypeName = proto.String(fieldType)
	***REMOVED***
	return fd
***REMOVED***

func (r *parseResult) asGroupDescriptors(group *groupNode, isProto3 bool, maxTag int32) (*dpb.FieldDescriptorProto, *dpb.DescriptorProto) ***REMOVED***
	tag := group.tag.val
	if err := checkTag(group.tag.start(), tag, maxTag); err != nil ***REMOVED***
		_ = r.errs.handleError(err)
	***REMOVED***
	if !unicode.IsUpper(rune(group.name.val[0])) ***REMOVED***
		_ = r.errs.handleErrorWithPos(group.name.start(), "group %s should have a name that starts with a capital letter", group.name.val)
	***REMOVED***
	fieldName := strings.ToLower(group.name.val)
	fd := &dpb.FieldDescriptorProto***REMOVED***
		Name:     proto.String(fieldName),
		JsonName: proto.String(internal.JsonName(fieldName)),
		Number:   proto.Int32(int32(tag)),
		Label:    asLabel(&group.label),
		Type:     dpb.FieldDescriptorProto_TYPE_GROUP.Enum(),
		TypeName: proto.String(group.name.val),
	***REMOVED***
	r.putFieldNode(fd, group)
	if opts := group.options.Elements(); len(opts) > 0 ***REMOVED***
		fd.Options = &dpb.FieldOptions***REMOVED***UninterpretedOption: r.asUninterpretedOptions(opts)***REMOVED***
	***REMOVED***
	md := &dpb.DescriptorProto***REMOVED***Name: proto.String(group.name.val)***REMOVED***
	r.putMessageNode(md, group)
	r.addMessageDecls(md, group.decls, isProto3)
	return fd, md
***REMOVED***

func (r *parseResult) asMapDescriptors(mapField *mapFieldNode, isProto3 bool, maxTag int32) (*dpb.FieldDescriptorProto, *dpb.DescriptorProto) ***REMOVED***
	tag := mapField.tag.val
	if err := checkTag(mapField.tag.start(), tag, maxTag); err != nil ***REMOVED***
		_ = r.errs.handleError(err)
	***REMOVED***
	var lbl *dpb.FieldDescriptorProto_Label
	if !isProto3 ***REMOVED***
		lbl = dpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum()
	***REMOVED***
	keyFd := newFieldDescriptor("key", mapField.mapType.keyType.val, 1, lbl)
	r.putFieldNode(keyFd, mapField.keyField())
	valFd := newFieldDescriptor("value", mapField.mapType.valueType.val, 2, lbl)
	r.putFieldNode(valFd, mapField.valueField())
	entryName := internal.InitCap(internal.JsonName(mapField.name.val)) + "Entry"
	fd := newFieldDescriptor(mapField.name.val, entryName, int32(tag), dpb.FieldDescriptorProto_LABEL_REPEATED.Enum())
	if opts := mapField.options.Elements(); len(opts) > 0 ***REMOVED***
		fd.Options = &dpb.FieldOptions***REMOVED***UninterpretedOption: r.asUninterpretedOptions(opts)***REMOVED***
	***REMOVED***
	r.putFieldNode(fd, mapField)
	md := &dpb.DescriptorProto***REMOVED***
		Name:    proto.String(entryName),
		Options: &dpb.MessageOptions***REMOVED***MapEntry: proto.Bool(true)***REMOVED***,
		Field:   []*dpb.FieldDescriptorProto***REMOVED***keyFd, valFd***REMOVED***,
	***REMOVED***
	r.putMessageNode(md, mapField)
	return fd, md
***REMOVED***

func (r *parseResult) asExtensionRanges(node *extensionRangeNode, maxTag int32) []*dpb.DescriptorProto_ExtensionRange ***REMOVED***
	opts := r.asUninterpretedOptions(node.options.Elements())
	ers := make([]*dpb.DescriptorProto_ExtensionRange, len(node.ranges))
	for i, rng := range node.ranges ***REMOVED***
		start, end := getRangeBounds(r, rng, 0, maxTag)
		er := &dpb.DescriptorProto_ExtensionRange***REMOVED***
			Start: proto.Int32(start),
			End:   proto.Int32(end + 1),
		***REMOVED***
		if len(opts) > 0 ***REMOVED***
			er.Options = &dpb.ExtensionRangeOptions***REMOVED***UninterpretedOption: opts***REMOVED***
		***REMOVED***
		r.putExtensionRangeNode(er, rng)
		ers[i] = er
	***REMOVED***
	return ers
***REMOVED***

func (r *parseResult) asEnumValue(ev *enumValueNode) *dpb.EnumValueDescriptorProto ***REMOVED***
	num, ok := ev.number.asInt32(math.MinInt32, math.MaxInt32)
	if !ok ***REMOVED***
		_ = r.errs.handleErrorWithPos(ev.number.start(), "value %d is out of range: should be between %d and %d", ev.number.value(), math.MinInt32, math.MaxInt32)
	***REMOVED***
	evd := &dpb.EnumValueDescriptorProto***REMOVED***Name: proto.String(ev.name.val), Number: proto.Int32(num)***REMOVED***
	r.putEnumValueNode(evd, ev)
	if opts := ev.options.Elements(); len(opts) > 0 ***REMOVED***
		evd.Options = &dpb.EnumValueOptions***REMOVED***UninterpretedOption: r.asUninterpretedOptions(opts)***REMOVED***
	***REMOVED***
	return evd
***REMOVED***

func (r *parseResult) asMethodDescriptor(node *methodNode) *dpb.MethodDescriptorProto ***REMOVED***
	md := &dpb.MethodDescriptorProto***REMOVED***
		Name:       proto.String(node.name.val),
		InputType:  proto.String(node.input.msgType.val),
		OutputType: proto.String(node.output.msgType.val),
	***REMOVED***
	r.putMethodNode(md, node)
	if node.input.streamKeyword != nil ***REMOVED***
		md.ClientStreaming = proto.Bool(true)
	***REMOVED***
	if node.output.streamKeyword != nil ***REMOVED***
		md.ServerStreaming = proto.Bool(true)
	***REMOVED***
	// protoc always adds a MethodOptions if there are brackets
	// We have a non-nil node.options if there are brackets
	// We do the same to match protoc as closely as possible
	// https://github.com/protocolbuffers/protobuf/blob/0c3f43a6190b77f1f68b7425d1b7e1a8257a8d0c/src/google/protobuf/compiler/parser.cc#L2152
	if node.options != nil ***REMOVED***
		md.Options = &dpb.MethodOptions***REMOVED***UninterpretedOption: r.asUninterpretedOptions(node.options)***REMOVED***
	***REMOVED***
	return md
***REMOVED***

func (r *parseResult) asEnumDescriptor(en *enumNode) *dpb.EnumDescriptorProto ***REMOVED***
	ed := &dpb.EnumDescriptorProto***REMOVED***Name: proto.String(en.name.val)***REMOVED***
	r.putEnumNode(ed, en)
	for _, decl := range en.decls ***REMOVED***
		if decl.option != nil ***REMOVED***
			if ed.Options == nil ***REMOVED***
				ed.Options = &dpb.EnumOptions***REMOVED******REMOVED***
			***REMOVED***
			ed.Options.UninterpretedOption = append(ed.Options.UninterpretedOption, r.asUninterpretedOption(decl.option))
		***REMOVED*** else if decl.value != nil ***REMOVED***
			ed.Value = append(ed.Value, r.asEnumValue(decl.value))
		***REMOVED*** else if decl.reserved != nil ***REMOVED***
			for _, n := range decl.reserved.names ***REMOVED***
				ed.ReservedName = append(ed.ReservedName, n.val)
			***REMOVED***
			for _, rng := range decl.reserved.ranges ***REMOVED***
				ed.ReservedRange = append(ed.ReservedRange, r.asEnumReservedRange(rng))
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return ed
***REMOVED***

func (r *parseResult) asEnumReservedRange(rng *rangeNode) *dpb.EnumDescriptorProto_EnumReservedRange ***REMOVED***
	start, end := getRangeBounds(r, rng, math.MinInt32, math.MaxInt32)
	rr := &dpb.EnumDescriptorProto_EnumReservedRange***REMOVED***
		Start: proto.Int32(start),
		End:   proto.Int32(end),
	***REMOVED***
	r.putEnumReservedRangeNode(rr, rng)
	return rr
***REMOVED***

func (r *parseResult) asMessageDescriptor(node *messageNode, isProto3 bool) *dpb.DescriptorProto ***REMOVED***
	msgd := &dpb.DescriptorProto***REMOVED***Name: proto.String(node.name.val)***REMOVED***
	r.putMessageNode(msgd, node)
	r.addMessageDecls(msgd, node.decls, isProto3)
	return msgd
***REMOVED***

func (r *parseResult) addMessageDecls(msgd *dpb.DescriptorProto, decls []*messageElement, isProto3 bool) ***REMOVED***
	// first process any options
	for _, decl := range decls ***REMOVED***
		if decl.option != nil ***REMOVED***
			if msgd.Options == nil ***REMOVED***
				msgd.Options = &dpb.MessageOptions***REMOVED******REMOVED***
			***REMOVED***
			msgd.Options.UninterpretedOption = append(msgd.Options.UninterpretedOption, r.asUninterpretedOption(decl.option))
		***REMOVED***
	***REMOVED***

	// now that we have options, we can see if this uses messageset wire format, which
	// impacts how we validate tag numbers in any fields in the message
	maxTag := int32(internal.MaxNormalTag)
	if isMessageSet, err := isMessageSetWireFormat(r, "message "+msgd.GetName(), msgd); err != nil ***REMOVED***
		return
	***REMOVED*** else if isMessageSet ***REMOVED***
		maxTag = internal.MaxTag // higher limit for messageset wire format
	***REMOVED***

	rsvdNames := map[string]int***REMOVED******REMOVED***

	// now we can process the rest
	for _, decl := range decls ***REMOVED***
		if decl.enum != nil ***REMOVED***
			msgd.EnumType = append(msgd.EnumType, r.asEnumDescriptor(decl.enum))
		***REMOVED*** else if decl.extend != nil ***REMOVED***
			r.addExtensions(decl.extend, &msgd.Extension, &msgd.NestedType, isProto3)
		***REMOVED*** else if decl.extensionRange != nil ***REMOVED***
			msgd.ExtensionRange = append(msgd.ExtensionRange, r.asExtensionRanges(decl.extensionRange, maxTag)...)
		***REMOVED*** else if decl.field != nil ***REMOVED***
			fd := r.asFieldDescriptor(decl.field, maxTag, isProto3)
			msgd.Field = append(msgd.Field, fd)
		***REMOVED*** else if decl.mapField != nil ***REMOVED***
			fd, md := r.asMapDescriptors(decl.mapField, isProto3, maxTag)
			msgd.Field = append(msgd.Field, fd)
			msgd.NestedType = append(msgd.NestedType, md)
		***REMOVED*** else if decl.group != nil ***REMOVED***
			fd, md := r.asGroupDescriptors(decl.group, isProto3, maxTag)
			msgd.Field = append(msgd.Field, fd)
			msgd.NestedType = append(msgd.NestedType, md)
		***REMOVED*** else if decl.oneOf != nil ***REMOVED***
			oodIndex := len(msgd.OneofDecl)
			ood := &dpb.OneofDescriptorProto***REMOVED***Name: proto.String(decl.oneOf.name.val)***REMOVED***
			r.putOneOfNode(ood, decl.oneOf)
			msgd.OneofDecl = append(msgd.OneofDecl, ood)
			ooFields := 0
			for _, oodecl := range decl.oneOf.decls ***REMOVED***
				if oodecl.option != nil ***REMOVED***
					if ood.Options == nil ***REMOVED***
						ood.Options = &dpb.OneofOptions***REMOVED******REMOVED***
					***REMOVED***
					ood.Options.UninterpretedOption = append(ood.Options.UninterpretedOption, r.asUninterpretedOption(oodecl.option))
				***REMOVED*** else if oodecl.field != nil ***REMOVED***
					fd := r.asFieldDescriptor(oodecl.field, maxTag, isProto3)
					fd.OneofIndex = proto.Int32(int32(oodIndex))
					msgd.Field = append(msgd.Field, fd)
					ooFields++
				***REMOVED*** else if oodecl.group != nil ***REMOVED***
					fd, md := r.asGroupDescriptors(oodecl.group, isProto3, maxTag)
					fd.OneofIndex = proto.Int32(int32(oodIndex))
					msgd.Field = append(msgd.Field, fd)
					msgd.NestedType = append(msgd.NestedType, md)
					ooFields++
				***REMOVED***
			***REMOVED***
			if ooFields == 0 ***REMOVED***
				_ = r.errs.handleErrorWithPos(decl.oneOf.start(), "oneof must contain at least one field")
			***REMOVED***
		***REMOVED*** else if decl.nested != nil ***REMOVED***
			msgd.NestedType = append(msgd.NestedType, r.asMessageDescriptor(decl.nested, isProto3))
		***REMOVED*** else if decl.reserved != nil ***REMOVED***
			for _, n := range decl.reserved.names ***REMOVED***
				count := rsvdNames[n.val]
				if count == 1 ***REMOVED*** // already seen
					_ = r.errs.handleErrorWithPos(n.start(), "name %q is reserved multiple times", n.val)
				***REMOVED***
				rsvdNames[n.val] = count + 1
				msgd.ReservedName = append(msgd.ReservedName, n.val)
			***REMOVED***
			for _, rng := range decl.reserved.ranges ***REMOVED***
				msgd.ReservedRange = append(msgd.ReservedRange, r.asMessageReservedRange(rng, maxTag))
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// process any proto3_optional fields
	if isProto3 ***REMOVED***
		internal.ProcessProto3OptionalFields(msgd)
	***REMOVED***
***REMOVED***

func isMessageSetWireFormat(res *parseResult, scope string, md *dpb.DescriptorProto) (bool, error) ***REMOVED***
	uo := md.GetOptions().GetUninterpretedOption()
	index, err := findOption(res, scope, uo, "message_set_wire_format")
	if err != nil ***REMOVED***
		return false, err
	***REMOVED***
	if index == -1 ***REMOVED***
		// no such option, so default to false
		return false, nil
	***REMOVED***

	opt := uo[index]
	optNode := res.getOptionNode(opt)

	switch opt.GetIdentifierValue() ***REMOVED***
	case "true":
		return true, nil
	case "false":
		return false, nil
	default:
		return false, res.errs.handleErrorWithPos(optNode.getValue().start(), "%s: expecting bool value for message_set_wire_format option", scope)
	***REMOVED***
***REMOVED***

func (r *parseResult) asMessageReservedRange(rng *rangeNode, maxTag int32) *dpb.DescriptorProto_ReservedRange ***REMOVED***
	start, end := getRangeBounds(r, rng, 0, maxTag)
	rr := &dpb.DescriptorProto_ReservedRange***REMOVED***
		Start: proto.Int32(start),
		End:   proto.Int32(end + 1),
	***REMOVED***
	r.putMessageReservedRangeNode(rr, rng)
	return rr
***REMOVED***

func getRangeBounds(res *parseResult, rng *rangeNode, minVal, maxVal int32) (int32, int32) ***REMOVED***
	checkOrder := true
	start, ok := rng.startValueAsInt32(minVal, maxVal)
	if !ok ***REMOVED***
		checkOrder = false
		_ = res.errs.handleErrorWithPos(rng.startNode.start(), "range start %d is out of range: should be between %d and %d", rng.startValue(), minVal, maxVal)
	***REMOVED***

	end, ok := rng.endValueAsInt32(minVal, maxVal)
	if !ok ***REMOVED***
		checkOrder = false
		if rng.endNode != nil ***REMOVED***
			_ = res.errs.handleErrorWithPos(rng.endNode.start(), "range end %d is out of range: should be between %d and %d", rng.endValue(), minVal, maxVal)
		***REMOVED***
	***REMOVED***

	if checkOrder && start > end ***REMOVED***
		_ = res.errs.handleErrorWithPos(rng.rangeStart().start(), "range, %d to %d, is invalid: start must be <= end", start, end)
	***REMOVED***

	return start, end
***REMOVED***

func (r *parseResult) asServiceDescriptor(svc *serviceNode) *dpb.ServiceDescriptorProto ***REMOVED***
	sd := &dpb.ServiceDescriptorProto***REMOVED***Name: proto.String(svc.name.val)***REMOVED***
	r.putServiceNode(sd, svc)
	for _, decl := range svc.decls ***REMOVED***
		if decl.option != nil ***REMOVED***
			if sd.Options == nil ***REMOVED***
				sd.Options = &dpb.ServiceOptions***REMOVED******REMOVED***
			***REMOVED***
			sd.Options.UninterpretedOption = append(sd.Options.UninterpretedOption, r.asUninterpretedOption(decl.option))
		***REMOVED*** else if decl.rpc != nil ***REMOVED***
			sd.Method = append(sd.Method, r.asMethodDescriptor(decl.rpc))
		***REMOVED***
	***REMOVED***
	return sd
***REMOVED***
