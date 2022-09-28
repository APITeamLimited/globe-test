package protoparse

import (
	"bytes"
	"math"
	"strings"
	"unicode"

	"github.com/golang/protobuf/proto"
	dpb "github.com/golang/protobuf/protoc-gen-go/descriptor"

	"github.com/jhump/protoreflect/desc/internal"
	"github.com/jhump/protoreflect/desc/protoparse/ast"
)

func (r *parseResult) createFileDescriptor(filename string, file *ast.FileNode) ***REMOVED***
	fd := &dpb.FileDescriptorProto***REMOVED***Name: proto.String(filename)***REMOVED***
	r.fd = fd
	r.putFileNode(fd, file)

	isProto3 := false
	if file.Syntax != nil ***REMOVED***
		if file.Syntax.Syntax.AsString() == "proto3" ***REMOVED***
			isProto3 = true
		***REMOVED*** else if file.Syntax.Syntax.AsString() != "proto2" ***REMOVED***
			if r.errs.handleErrorWithPos(file.Syntax.Syntax.Start(), `syntax value must be "proto2" or "proto3"`) != nil ***REMOVED***
				return
			***REMOVED***
		***REMOVED***

		// proto2 is the default, so no need to set unless proto3
		if isProto3 ***REMOVED***
			fd.Syntax = proto.String(file.Syntax.Syntax.AsString())
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		r.errs.warn(file.Start(), ErrNoSyntax)
	***REMOVED***

	for _, decl := range file.Decls ***REMOVED***
		if r.errs.err != nil ***REMOVED***
			return
		***REMOVED***
		switch decl := decl.(type) ***REMOVED***
		case *ast.EnumNode:
			fd.EnumType = append(fd.EnumType, r.asEnumDescriptor(decl))
		case *ast.ExtendNode:
			r.addExtensions(decl, &fd.Extension, &fd.MessageType, isProto3)
		case *ast.ImportNode:
			index := len(fd.Dependency)
			fd.Dependency = append(fd.Dependency, decl.Name.AsString())
			if decl.Public != nil ***REMOVED***
				fd.PublicDependency = append(fd.PublicDependency, int32(index))
			***REMOVED*** else if decl.Weak != nil ***REMOVED***
				fd.WeakDependency = append(fd.WeakDependency, int32(index))
			***REMOVED***
		case *ast.MessageNode:
			fd.MessageType = append(fd.MessageType, r.asMessageDescriptor(decl, isProto3))
		case *ast.OptionNode:
			if fd.Options == nil ***REMOVED***
				fd.Options = &dpb.FileOptions***REMOVED******REMOVED***
			***REMOVED***
			fd.Options.UninterpretedOption = append(fd.Options.UninterpretedOption, r.asUninterpretedOption(decl))
		case *ast.ServiceNode:
			fd.Service = append(fd.Service, r.asServiceDescriptor(decl))
		case *ast.PackageNode:
			if fd.Package != nil ***REMOVED***
				if r.errs.handleErrorWithPos(decl.Start(), "files should have only one package declaration") != nil ***REMOVED***
					return
				***REMOVED***
			***REMOVED***
			fd.Package = proto.String(string(decl.Name.AsIdentifier()))
		***REMOVED***
	***REMOVED***
***REMOVED***

func (r *parseResult) asUninterpretedOptions(nodes []*ast.OptionNode) []*dpb.UninterpretedOption ***REMOVED***
	if len(nodes) == 0 ***REMOVED***
		return nil
	***REMOVED***
	opts := make([]*dpb.UninterpretedOption, len(nodes))
	for i, n := range nodes ***REMOVED***
		opts[i] = r.asUninterpretedOption(n)
	***REMOVED***
	return opts
***REMOVED***

func (r *parseResult) asUninterpretedOption(node *ast.OptionNode) *dpb.UninterpretedOption ***REMOVED***
	opt := &dpb.UninterpretedOption***REMOVED***Name: r.asUninterpretedOptionName(node.Name.Parts)***REMOVED***
	r.putOptionNode(opt, node)

	switch val := node.Val.Value().(type) ***REMOVED***
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
	case ast.Identifier:
		opt.IdentifierValue = proto.String(string(val))
	case []*ast.MessageFieldNode:
		var buf bytes.Buffer
		aggToString(val, &buf)
		aggStr := buf.String()
		opt.AggregateValue = proto.String(aggStr)
		//the grammar does not allow arrays here, so no case for []ast.ValueNode
	***REMOVED***
	return opt
***REMOVED***

func (r *parseResult) asUninterpretedOptionName(parts []*ast.FieldReferenceNode) []*dpb.UninterpretedOption_NamePart ***REMOVED***
	ret := make([]*dpb.UninterpretedOption_NamePart, len(parts))
	for i, part := range parts ***REMOVED***
		np := &dpb.UninterpretedOption_NamePart***REMOVED***
			NamePart:    proto.String(string(part.Name.AsIdentifier())),
			IsExtension: proto.Bool(part.IsExtension()),
		***REMOVED***
		r.putOptionNamePartNode(np, part)
		ret[i] = np
	***REMOVED***
	return ret
***REMOVED***

func (r *parseResult) addExtensions(ext *ast.ExtendNode, flds *[]*dpb.FieldDescriptorProto, msgs *[]*dpb.DescriptorProto, isProto3 bool) ***REMOVED***
	extendee := string(ext.Extendee.AsIdentifier())
	count := 0
	for _, decl := range ext.Decls ***REMOVED***
		switch decl := decl.(type) ***REMOVED***
		case *ast.FieldNode:
			count++
			// use higher limit since we don't know yet whether extendee is messageset wire format
			fd := r.asFieldDescriptor(decl, internal.MaxTag, isProto3)
			fd.Extendee = proto.String(extendee)
			*flds = append(*flds, fd)
		case *ast.GroupNode:
			count++
			// ditto: use higher limit right now
			fd, md := r.asGroupDescriptors(decl, isProto3, internal.MaxTag)
			fd.Extendee = proto.String(extendee)
			*flds = append(*flds, fd)
			*msgs = append(*msgs, md)
		***REMOVED***
	***REMOVED***
	if count == 0 ***REMOVED***
		_ = r.errs.handleErrorWithPos(ext.Start(), "extend sections must define at least one extension")
	***REMOVED***
***REMOVED***

func asLabel(lbl *ast.FieldLabel) *dpb.FieldDescriptorProto_Label ***REMOVED***
	if !lbl.IsPresent() ***REMOVED***
		return nil
	***REMOVED***
	switch ***REMOVED***
	case lbl.Repeated:
		return dpb.FieldDescriptorProto_LABEL_REPEATED.Enum()
	case lbl.Required:
		return dpb.FieldDescriptorProto_LABEL_REQUIRED.Enum()
	default:
		return dpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum()
	***REMOVED***
***REMOVED***

func (r *parseResult) asFieldDescriptor(node *ast.FieldNode, maxTag int32, isProto3 bool) *dpb.FieldDescriptorProto ***REMOVED***
	tag := node.Tag.Val
	if err := checkTag(node.Tag.Start(), tag, maxTag); err != nil ***REMOVED***
		_ = r.errs.handleError(err)
	***REMOVED***
	fd := newFieldDescriptor(node.Name.Val, string(node.FldType.AsIdentifier()), int32(tag), asLabel(&node.Label))
	r.putFieldNode(fd, node)
	if opts := node.Options.GetElements(); len(opts) > 0 ***REMOVED***
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

func (r *parseResult) asGroupDescriptors(group *ast.GroupNode, isProto3 bool, maxTag int32) (*dpb.FieldDescriptorProto, *dpb.DescriptorProto) ***REMOVED***
	tag := group.Tag.Val
	if err := checkTag(group.Tag.Start(), tag, maxTag); err != nil ***REMOVED***
		_ = r.errs.handleError(err)
	***REMOVED***
	if !unicode.IsUpper(rune(group.Name.Val[0])) ***REMOVED***
		_ = r.errs.handleErrorWithPos(group.Name.Start(), "group %s should have a name that starts with a capital letter", group.Name.Val)
	***REMOVED***
	fieldName := strings.ToLower(group.Name.Val)
	fd := &dpb.FieldDescriptorProto***REMOVED***
		Name:     proto.String(fieldName),
		JsonName: proto.String(internal.JsonName(fieldName)),
		Number:   proto.Int32(int32(tag)),
		Label:    asLabel(&group.Label),
		Type:     dpb.FieldDescriptorProto_TYPE_GROUP.Enum(),
		TypeName: proto.String(group.Name.Val),
	***REMOVED***
	r.putFieldNode(fd, group)
	if opts := group.Options.GetElements(); len(opts) > 0 ***REMOVED***
		fd.Options = &dpb.FieldOptions***REMOVED***UninterpretedOption: r.asUninterpretedOptions(opts)***REMOVED***
	***REMOVED***
	md := &dpb.DescriptorProto***REMOVED***Name: proto.String(group.Name.Val)***REMOVED***
	r.putMessageNode(md, group)
	r.addMessageBody(md, &group.MessageBody, isProto3)
	return fd, md
***REMOVED***

func (r *parseResult) asMapDescriptors(mapField *ast.MapFieldNode, isProto3 bool, maxTag int32) (*dpb.FieldDescriptorProto, *dpb.DescriptorProto) ***REMOVED***
	tag := mapField.Tag.Val
	if err := checkTag(mapField.Tag.Start(), tag, maxTag); err != nil ***REMOVED***
		_ = r.errs.handleError(err)
	***REMOVED***
	var lbl *dpb.FieldDescriptorProto_Label
	if !isProto3 ***REMOVED***
		lbl = dpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum()
	***REMOVED***
	keyFd := newFieldDescriptor("key", mapField.MapType.KeyType.Val, 1, lbl)
	r.putFieldNode(keyFd, mapField.KeyField())
	valFd := newFieldDescriptor("value", string(mapField.MapType.ValueType.AsIdentifier()), 2, lbl)
	r.putFieldNode(valFd, mapField.ValueField())
	entryName := internal.InitCap(internal.JsonName(mapField.Name.Val)) + "Entry"
	fd := newFieldDescriptor(mapField.Name.Val, entryName, int32(tag), dpb.FieldDescriptorProto_LABEL_REPEATED.Enum())
	if opts := mapField.Options.GetElements(); len(opts) > 0 ***REMOVED***
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

func (r *parseResult) asExtensionRanges(node *ast.ExtensionRangeNode, maxTag int32) []*dpb.DescriptorProto_ExtensionRange ***REMOVED***
	opts := r.asUninterpretedOptions(node.Options.GetElements())
	ers := make([]*dpb.DescriptorProto_ExtensionRange, len(node.Ranges))
	for i, rng := range node.Ranges ***REMOVED***
		start, end := getRangeBounds(r, rng, 1, maxTag)
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

func (r *parseResult) asEnumValue(ev *ast.EnumValueNode) *dpb.EnumValueDescriptorProto ***REMOVED***
	num, ok := ast.AsInt32(ev.Number, math.MinInt32, math.MaxInt32)
	if !ok ***REMOVED***
		_ = r.errs.handleErrorWithPos(ev.Number.Start(), "value %d is out of range: should be between %d and %d", ev.Number.Value(), math.MinInt32, math.MaxInt32)
	***REMOVED***
	evd := &dpb.EnumValueDescriptorProto***REMOVED***Name: proto.String(ev.Name.Val), Number: proto.Int32(num)***REMOVED***
	r.putEnumValueNode(evd, ev)
	if opts := ev.Options.GetElements(); len(opts) > 0 ***REMOVED***
		evd.Options = &dpb.EnumValueOptions***REMOVED***UninterpretedOption: r.asUninterpretedOptions(opts)***REMOVED***
	***REMOVED***
	return evd
***REMOVED***

func (r *parseResult) asMethodDescriptor(node *ast.RPCNode) *dpb.MethodDescriptorProto ***REMOVED***
	md := &dpb.MethodDescriptorProto***REMOVED***
		Name:       proto.String(node.Name.Val),
		InputType:  proto.String(string(node.Input.MessageType.AsIdentifier())),
		OutputType: proto.String(string(node.Output.MessageType.AsIdentifier())),
	***REMOVED***
	r.putMethodNode(md, node)
	if node.Input.Stream != nil ***REMOVED***
		md.ClientStreaming = proto.Bool(true)
	***REMOVED***
	if node.Output.Stream != nil ***REMOVED***
		md.ServerStreaming = proto.Bool(true)
	***REMOVED***
	// protoc always adds a MethodOptions if there are brackets
	// We do the same to match protoc as closely as possible
	// https://github.com/protocolbuffers/protobuf/blob/0c3f43a6190b77f1f68b7425d1b7e1a8257a8d0c/src/google/protobuf/compiler/parser.cc#L2152
	if node.OpenBrace != nil ***REMOVED***
		md.Options = &dpb.MethodOptions***REMOVED******REMOVED***
		for _, decl := range node.Decls ***REMOVED***
			switch decl := decl.(type) ***REMOVED***
			case *ast.OptionNode:
				md.Options.UninterpretedOption = append(md.Options.UninterpretedOption, r.asUninterpretedOption(decl))
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return md
***REMOVED***

func (r *parseResult) asEnumDescriptor(en *ast.EnumNode) *dpb.EnumDescriptorProto ***REMOVED***
	ed := &dpb.EnumDescriptorProto***REMOVED***Name: proto.String(en.Name.Val)***REMOVED***
	r.putEnumNode(ed, en)
	for _, decl := range en.Decls ***REMOVED***
		switch decl := decl.(type) ***REMOVED***
		case *ast.OptionNode:
			if ed.Options == nil ***REMOVED***
				ed.Options = &dpb.EnumOptions***REMOVED******REMOVED***
			***REMOVED***
			ed.Options.UninterpretedOption = append(ed.Options.UninterpretedOption, r.asUninterpretedOption(decl))
		case *ast.EnumValueNode:
			ed.Value = append(ed.Value, r.asEnumValue(decl))
		case *ast.ReservedNode:
			for _, n := range decl.Names ***REMOVED***
				ed.ReservedName = append(ed.ReservedName, n.AsString())
			***REMOVED***
			for _, rng := range decl.Ranges ***REMOVED***
				ed.ReservedRange = append(ed.ReservedRange, r.asEnumReservedRange(rng))
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return ed
***REMOVED***

func (r *parseResult) asEnumReservedRange(rng *ast.RangeNode) *dpb.EnumDescriptorProto_EnumReservedRange ***REMOVED***
	start, end := getRangeBounds(r, rng, math.MinInt32, math.MaxInt32)
	rr := &dpb.EnumDescriptorProto_EnumReservedRange***REMOVED***
		Start: proto.Int32(start),
		End:   proto.Int32(end),
	***REMOVED***
	r.putEnumReservedRangeNode(rr, rng)
	return rr
***REMOVED***

func (r *parseResult) asMessageDescriptor(node *ast.MessageNode, isProto3 bool) *dpb.DescriptorProto ***REMOVED***
	msgd := &dpb.DescriptorProto***REMOVED***Name: proto.String(node.Name.Val)***REMOVED***
	r.putMessageNode(msgd, node)
	r.addMessageBody(msgd, &node.MessageBody, isProto3)
	return msgd
***REMOVED***

func (r *parseResult) addMessageBody(msgd *dpb.DescriptorProto, body *ast.MessageBody, isProto3 bool) ***REMOVED***
	// first process any options
	for _, decl := range body.Decls ***REMOVED***
		if opt, ok := decl.(*ast.OptionNode); ok ***REMOVED***
			if msgd.Options == nil ***REMOVED***
				msgd.Options = &dpb.MessageOptions***REMOVED******REMOVED***
			***REMOVED***
			msgd.Options.UninterpretedOption = append(msgd.Options.UninterpretedOption, r.asUninterpretedOption(opt))
		***REMOVED***
	***REMOVED***

	// now that we have options, we can see if this uses messageset wire format, which
	// impacts how we validate tag numbers in any fields in the message
	maxTag := int32(internal.MaxNormalTag)
	messageSetOpt, err := isMessageSetWireFormat(r, "message "+msgd.GetName(), msgd)
	if err != nil ***REMOVED***
		return
	***REMOVED*** else if messageSetOpt != nil ***REMOVED***
		maxTag = internal.MaxTag // higher limit for messageset wire format
	***REMOVED***

	rsvdNames := map[string]int***REMOVED******REMOVED***

	// now we can process the rest
	for _, decl := range body.Decls ***REMOVED***
		switch decl := decl.(type) ***REMOVED***
		case *ast.EnumNode:
			msgd.EnumType = append(msgd.EnumType, r.asEnumDescriptor(decl))
		case *ast.ExtendNode:
			r.addExtensions(decl, &msgd.Extension, &msgd.NestedType, isProto3)
		case *ast.ExtensionRangeNode:
			msgd.ExtensionRange = append(msgd.ExtensionRange, r.asExtensionRanges(decl, maxTag)...)
		case *ast.FieldNode:
			fd := r.asFieldDescriptor(decl, maxTag, isProto3)
			msgd.Field = append(msgd.Field, fd)
		case *ast.MapFieldNode:
			fd, md := r.asMapDescriptors(decl, isProto3, maxTag)
			msgd.Field = append(msgd.Field, fd)
			msgd.NestedType = append(msgd.NestedType, md)
		case *ast.GroupNode:
			fd, md := r.asGroupDescriptors(decl, isProto3, maxTag)
			msgd.Field = append(msgd.Field, fd)
			msgd.NestedType = append(msgd.NestedType, md)
		case *ast.OneOfNode:
			oodIndex := len(msgd.OneofDecl)
			ood := &dpb.OneofDescriptorProto***REMOVED***Name: proto.String(decl.Name.Val)***REMOVED***
			r.putOneOfNode(ood, decl)
			msgd.OneofDecl = append(msgd.OneofDecl, ood)
			ooFields := 0
			for _, oodecl := range decl.Decls ***REMOVED***
				switch oodecl := oodecl.(type) ***REMOVED***
				case *ast.OptionNode:
					if ood.Options == nil ***REMOVED***
						ood.Options = &dpb.OneofOptions***REMOVED******REMOVED***
					***REMOVED***
					ood.Options.UninterpretedOption = append(ood.Options.UninterpretedOption, r.asUninterpretedOption(oodecl))
				case *ast.FieldNode:
					fd := r.asFieldDescriptor(oodecl, maxTag, isProto3)
					fd.OneofIndex = proto.Int32(int32(oodIndex))
					msgd.Field = append(msgd.Field, fd)
					ooFields++
				case *ast.GroupNode:
					fd, md := r.asGroupDescriptors(oodecl, isProto3, maxTag)
					fd.OneofIndex = proto.Int32(int32(oodIndex))
					msgd.Field = append(msgd.Field, fd)
					msgd.NestedType = append(msgd.NestedType, md)
					ooFields++
				***REMOVED***
			***REMOVED***
			if ooFields == 0 ***REMOVED***
				_ = r.errs.handleErrorWithPos(decl.Start(), "oneof must contain at least one field")
			***REMOVED***
		case *ast.MessageNode:
			msgd.NestedType = append(msgd.NestedType, r.asMessageDescriptor(decl, isProto3))
		case *ast.ReservedNode:
			for _, n := range decl.Names ***REMOVED***
				count := rsvdNames[n.AsString()]
				if count == 1 ***REMOVED*** // already seen
					_ = r.errs.handleErrorWithPos(n.Start(), "name %q is reserved multiple times", n.AsString())
				***REMOVED***
				rsvdNames[n.AsString()] = count + 1
				msgd.ReservedName = append(msgd.ReservedName, n.AsString())
			***REMOVED***
			for _, rng := range decl.Ranges ***REMOVED***
				msgd.ReservedRange = append(msgd.ReservedRange, r.asMessageReservedRange(rng, maxTag))
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if messageSetOpt != nil ***REMOVED***
		if len(msgd.Field) > 0 ***REMOVED***
			node := r.getFieldNode(msgd.Field[0])
			_ = r.errs.handleErrorWithPos(node.Start(), "messages with message-set wire format cannot contain non-extension fields")
		***REMOVED***
		if len(msgd.ExtensionRange) == 0 ***REMOVED***
			node := r.getOptionNode(messageSetOpt)
			_ = r.errs.handleErrorWithPos(node.Start(), "messages with message-set wire format must contain at least one extension range")
		***REMOVED***
	***REMOVED***

	// process any proto3_optional fields
	if isProto3 ***REMOVED***
		internal.ProcessProto3OptionalFields(msgd)
	***REMOVED***
***REMOVED***

func isMessageSetWireFormat(res *parseResult, scope string, md *dpb.DescriptorProto) (*dpb.UninterpretedOption, error) ***REMOVED***
	uo := md.GetOptions().GetUninterpretedOption()
	index, err := findOption(res, scope, uo, "message_set_wire_format")
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if index == -1 ***REMOVED***
		// no such option
		return nil, nil
	***REMOVED***

	opt := uo[index]

	switch opt.GetIdentifierValue() ***REMOVED***
	case "true":
		return opt, nil
	case "false":
		return nil, nil
	default:
		optNode := res.getOptionNode(opt)
		return nil, res.errs.handleErrorWithPos(optNode.GetValue().Start(), "%s: expecting bool value for message_set_wire_format option", scope)
	***REMOVED***
***REMOVED***

func (r *parseResult) asMessageReservedRange(rng *ast.RangeNode, maxTag int32) *dpb.DescriptorProto_ReservedRange ***REMOVED***
	start, end := getRangeBounds(r, rng, 1, maxTag)
	rr := &dpb.DescriptorProto_ReservedRange***REMOVED***
		Start: proto.Int32(start),
		End:   proto.Int32(end + 1),
	***REMOVED***
	r.putMessageReservedRangeNode(rr, rng)
	return rr
***REMOVED***

func getRangeBounds(res *parseResult, rng *ast.RangeNode, minVal, maxVal int32) (int32, int32) ***REMOVED***
	checkOrder := true
	start, ok := rng.StartValueAsInt32(minVal, maxVal)
	if !ok ***REMOVED***
		checkOrder = false
		_ = res.errs.handleErrorWithPos(rng.StartVal.Start(), "range start %d is out of range: should be between %d and %d", rng.StartValue(), minVal, maxVal)
	***REMOVED***

	end, ok := rng.EndValueAsInt32(minVal, maxVal)
	if !ok ***REMOVED***
		checkOrder = false
		if rng.EndVal != nil ***REMOVED***
			_ = res.errs.handleErrorWithPos(rng.EndVal.Start(), "range end %d is out of range: should be between %d and %d", rng.EndValue(), minVal, maxVal)
		***REMOVED***
	***REMOVED***

	if checkOrder && start > end ***REMOVED***
		_ = res.errs.handleErrorWithPos(rng.RangeStart().Start(), "range, %d to %d, is invalid: start must be <= end", start, end)
	***REMOVED***

	return start, end
***REMOVED***

func (r *parseResult) asServiceDescriptor(svc *ast.ServiceNode) *dpb.ServiceDescriptorProto ***REMOVED***
	sd := &dpb.ServiceDescriptorProto***REMOVED***Name: proto.String(svc.Name.Val)***REMOVED***
	r.putServiceNode(sd, svc)
	for _, decl := range svc.Decls ***REMOVED***
		switch decl := decl.(type) ***REMOVED***
		case *ast.OptionNode:
			if sd.Options == nil ***REMOVED***
				sd.Options = &dpb.ServiceOptions***REMOVED******REMOVED***
			***REMOVED***
			sd.Options.UninterpretedOption = append(sd.Options.UninterpretedOption, r.asUninterpretedOption(decl))
		case *ast.RPCNode:
			sd.Method = append(sd.Method, r.asMethodDescriptor(decl))
		***REMOVED***
	***REMOVED***
	return sd
***REMOVED***
