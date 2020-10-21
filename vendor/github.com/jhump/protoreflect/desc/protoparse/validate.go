package protoparse

import (
	"fmt"
	"sort"

	"github.com/golang/protobuf/proto"

	dpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
)

func validateBasic(res *parseResult, containsErrors bool) ***REMOVED***
	fd := res.fd
	isProto3 := fd.GetSyntax() == "proto3"

	for _, md := range fd.MessageType ***REMOVED***
		if validateMessage(res, isProto3, "", md, containsErrors) != nil ***REMOVED***
			return
		***REMOVED***
	***REMOVED***

	for _, ed := range fd.EnumType ***REMOVED***
		if validateEnum(res, isProto3, "", ed, containsErrors) != nil ***REMOVED***
			return
		***REMOVED***
	***REMOVED***

	for _, fld := range fd.Extension ***REMOVED***
		if validateField(res, isProto3, "", fld) != nil ***REMOVED***
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

func validateMessage(res *parseResult, isProto3 bool, prefix string, md *dpb.DescriptorProto, containsErrors bool) error ***REMOVED***
	nextPrefix := md.GetName() + "."

	for _, fld := range md.Field ***REMOVED***
		if err := validateField(res, isProto3, nextPrefix, fld); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	for _, fld := range md.Extension ***REMOVED***
		if err := validateField(res, isProto3, nextPrefix, fld); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	for _, ed := range md.EnumType ***REMOVED***
		if err := validateEnum(res, isProto3, nextPrefix, ed, containsErrors); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	for _, nmd := range md.NestedType ***REMOVED***
		if err := validateMessage(res, isProto3, nextPrefix, nmd, containsErrors); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	scope := fmt.Sprintf("message %s%s", prefix, md.GetName())

	if isProto3 && len(md.ExtensionRange) > 0 ***REMOVED***
		n := res.getExtensionRangeNode(md.ExtensionRange[0])
		if err := res.errs.handleErrorWithPos(n.start(), "%s: extension ranges are not allowed in proto3", scope); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	if index, err := findOption(res, scope, md.Options.GetUninterpretedOption(), "map_entry"); err != nil ***REMOVED***
		return err
	***REMOVED*** else if index >= 0 ***REMOVED***
		opt := md.Options.UninterpretedOption[index]
		optn := res.getOptionNode(opt)
		md.Options.UninterpretedOption = removeOption(md.Options.UninterpretedOption, index)
		valid := false
		if opt.IdentifierValue != nil ***REMOVED***
			if opt.GetIdentifierValue() == "true" ***REMOVED***
				valid = true
				if err := res.errs.handleErrorWithPos(optn.getValue().start(), "%s: map_entry option should not be set explicitly; use map type instead", scope); err != nil ***REMOVED***
					return err
				***REMOVED***
			***REMOVED*** else if opt.GetIdentifierValue() == "false" ***REMOVED***
				valid = true
				md.Options.MapEntry = proto.Bool(false)
			***REMOVED***
		***REMOVED***
		if !valid ***REMOVED***
			if err := res.errs.handleErrorWithPos(optn.getValue().start(), "%s: expecting bool value for map_entry option", scope); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// reserved ranges should not overlap
	rsvd := make(tagRanges, len(md.ReservedRange))
	for i, r := range md.ReservedRange ***REMOVED***
		n := res.getMessageReservedRangeNode(r)
		rsvd[i] = tagRange***REMOVED***start: r.GetStart(), end: r.GetEnd(), node: n***REMOVED***

	***REMOVED***
	sort.Sort(rsvd)
	for i := 1; i < len(rsvd); i++ ***REMOVED***
		if rsvd[i].start < rsvd[i-1].end ***REMOVED***
			if err := res.errs.handleErrorWithPos(rsvd[i].node.start(), "%s: reserved ranges overlap: %d to %d and %d to %d", scope, rsvd[i-1].start, rsvd[i-1].end-1, rsvd[i].start, rsvd[i].end-1); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// extensions ranges should not overlap
	exts := make(tagRanges, len(md.ExtensionRange))
	for i, r := range md.ExtensionRange ***REMOVED***
		n := res.getExtensionRangeNode(r)
		exts[i] = tagRange***REMOVED***start: r.GetStart(), end: r.GetEnd(), node: n***REMOVED***
	***REMOVED***
	sort.Sort(exts)
	for i := 1; i < len(exts); i++ ***REMOVED***
		if exts[i].start < exts[i-1].end ***REMOVED***
			if err := res.errs.handleErrorWithPos(exts[i].node.start(), "%s: extension ranges overlap: %d to %d and %d to %d", scope, exts[i-1].start, exts[i-1].end-1, exts[i].start, exts[i].end-1); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// see if any extension range overlaps any reserved range
	var i, j int // i indexes rsvd; j indexes exts
	for i < len(rsvd) && j < len(exts) ***REMOVED***
		if rsvd[i].start >= exts[j].start && rsvd[i].start < exts[j].end ||
			exts[j].start >= rsvd[i].start && exts[j].start < rsvd[i].end ***REMOVED***

			var pos *SourcePos
			if rsvd[i].start >= exts[j].start && rsvd[i].start < exts[j].end ***REMOVED***
				pos = rsvd[i].node.start()
			***REMOVED*** else ***REMOVED***
				pos = exts[j].node.start()
			***REMOVED***
			// ranges overlap
			if err := res.errs.handleErrorWithPos(pos, "%s: extension range %d to %d overlaps reserved range %d to %d", scope, exts[j].start, exts[j].end-1, rsvd[i].start, rsvd[i].end-1); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
		if rsvd[i].start < exts[j].start ***REMOVED***
			i++
		***REMOVED*** else ***REMOVED***
			j++
		***REMOVED***
	***REMOVED***

	// now, check that fields don't re-use tags and don't try to use extension
	// or reserved ranges or reserved names
	rsvdNames := map[string]struct***REMOVED******REMOVED******REMOVED******REMOVED***
	for _, n := range md.ReservedName ***REMOVED***
		rsvdNames[n] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	***REMOVED***
	fieldTags := map[int32]string***REMOVED******REMOVED***
	for _, fld := range md.Field ***REMOVED***
		fn := res.getFieldNode(fld)
		if _, ok := rsvdNames[fld.GetName()]; ok ***REMOVED***
			if err := res.errs.handleErrorWithPos(fn.fieldName().start(), "%s: field %s is using a reserved name", scope, fld.GetName()); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
		if existing := fieldTags[fld.GetNumber()]; existing != "" ***REMOVED***
			if err := res.errs.handleErrorWithPos(fn.fieldTag().start(), "%s: fields %s and %s both have the same tag %d", scope, existing, fld.GetName(), fld.GetNumber()); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
		fieldTags[fld.GetNumber()] = fld.GetName()
		// check reserved ranges
		r := sort.Search(len(rsvd), func(index int) bool ***REMOVED*** return rsvd[index].end > fld.GetNumber() ***REMOVED***)
		if r < len(rsvd) && rsvd[r].start <= fld.GetNumber() ***REMOVED***
			if err := res.errs.handleErrorWithPos(fn.fieldTag().start(), "%s: field %s is using tag %d which is in reserved range %d to %d", scope, fld.GetName(), fld.GetNumber(), rsvd[r].start, rsvd[r].end-1); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
		// and check extension ranges
		e := sort.Search(len(exts), func(index int) bool ***REMOVED*** return exts[index].end > fld.GetNumber() ***REMOVED***)
		if e < len(exts) && exts[e].start <= fld.GetNumber() ***REMOVED***
			if err := res.errs.handleErrorWithPos(fn.fieldTag().start(), "%s: field %s is using tag %d which is in extension range %d to %d", scope, fld.GetName(), fld.GetNumber(), exts[e].start, exts[e].end-1); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func validateEnum(res *parseResult, isProto3 bool, prefix string, ed *dpb.EnumDescriptorProto, containsErrors bool) error ***REMOVED***
	scope := fmt.Sprintf("enum %s%s", prefix, ed.GetName())

	if !containsErrors && len(ed.Value) == 0 ***REMOVED***
		// we only check this if file parsing had no errors; otherwise, the file may have
		// had an enum value, but the parser encountered an error processing it, in which
		// case the value would be absent from the descriptor. In such a case, this error
		// would be confusing and incorrect, so we just skip this check.
		enNode := res.getEnumNode(ed)
		if err := res.errs.handleErrorWithPos(enNode.start(), "%s: enums must define at least one value", scope); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	allowAlias := false
	if index, err := findOption(res, scope, ed.Options.GetUninterpretedOption(), "allow_alias"); err != nil ***REMOVED***
		return err
	***REMOVED*** else if index >= 0 ***REMOVED***
		opt := ed.Options.UninterpretedOption[index]
		valid := false
		if opt.IdentifierValue != nil ***REMOVED***
			if opt.GetIdentifierValue() == "true" ***REMOVED***
				allowAlias = true
				valid = true
			***REMOVED*** else if opt.GetIdentifierValue() == "false" ***REMOVED***
				valid = true
			***REMOVED***
		***REMOVED***
		if !valid ***REMOVED***
			optNode := res.getOptionNode(opt)
			if err := res.errs.handleErrorWithPos(optNode.getValue().start(), "%s: expecting bool value for allow_alias option", scope); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if isProto3 && len(ed.Value) > 0 && ed.Value[0].GetNumber() != 0 ***REMOVED***
		evNode := res.getEnumValueNode(ed.Value[0])
		if err := res.errs.handleErrorWithPos(evNode.getNumber().start(), "%s: proto3 requires that first value in enum have numeric value of 0", scope); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	if !allowAlias ***REMOVED***
		// make sure all value numbers are distinct
		vals := map[int32]string***REMOVED******REMOVED***
		for _, evd := range ed.Value ***REMOVED***
			if existing := vals[evd.GetNumber()]; existing != "" ***REMOVED***
				evNode := res.getEnumValueNode(evd)
				if err := res.errs.handleErrorWithPos(evNode.getNumber().start(), "%s: values %s and %s both have the same numeric value %d; use allow_alias option if intentional", scope, existing, evd.GetName(), evd.GetNumber()); err != nil ***REMOVED***
					return err
				***REMOVED***
			***REMOVED***
			vals[evd.GetNumber()] = evd.GetName()
		***REMOVED***
	***REMOVED***

	// reserved ranges should not overlap
	rsvd := make(tagRanges, len(ed.ReservedRange))
	for i, r := range ed.ReservedRange ***REMOVED***
		n := res.getEnumReservedRangeNode(r)
		rsvd[i] = tagRange***REMOVED***start: r.GetStart(), end: r.GetEnd(), node: n***REMOVED***
	***REMOVED***
	sort.Sort(rsvd)
	for i := 1; i < len(rsvd); i++ ***REMOVED***
		if rsvd[i].start <= rsvd[i-1].end ***REMOVED***
			if err := res.errs.handleErrorWithPos(rsvd[i].node.start(), "%s: reserved ranges overlap: %d to %d and %d to %d", scope, rsvd[i-1].start, rsvd[i-1].end, rsvd[i].start, rsvd[i].end); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// now, check that fields don't re-use tags and don't try to use extension
	// or reserved ranges or reserved names
	rsvdNames := map[string]struct***REMOVED******REMOVED******REMOVED******REMOVED***
	for _, n := range ed.ReservedName ***REMOVED***
		rsvdNames[n] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	***REMOVED***
	for _, ev := range ed.Value ***REMOVED***
		evn := res.getEnumValueNode(ev)
		if _, ok := rsvdNames[ev.GetName()]; ok ***REMOVED***
			if err := res.errs.handleErrorWithPos(evn.getName().start(), "%s: value %s is using a reserved name", scope, ev.GetName()); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
		// check reserved ranges
		r := sort.Search(len(rsvd), func(index int) bool ***REMOVED*** return rsvd[index].end >= ev.GetNumber() ***REMOVED***)
		if r < len(rsvd) && rsvd[r].start <= ev.GetNumber() ***REMOVED***
			if err := res.errs.handleErrorWithPos(evn.getNumber().start(), "%s: value %s is using number %d which is in reserved range %d to %d", scope, ev.GetName(), ev.GetNumber(), rsvd[r].start, rsvd[r].end); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func validateField(res *parseResult, isProto3 bool, prefix string, fld *dpb.FieldDescriptorProto) error ***REMOVED***
	scope := fmt.Sprintf("field %s%s", prefix, fld.GetName())

	node := res.getFieldNode(fld)
	if isProto3 ***REMOVED***
		if fld.GetType() == dpb.FieldDescriptorProto_TYPE_GROUP ***REMOVED***
			n := node.(*groupNode)
			if err := res.errs.handleErrorWithPos(n.groupKeyword.start(), "%s: groups are not allowed in proto3", scope); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED*** else if fld.Label != nil && fld.GetLabel() == dpb.FieldDescriptorProto_LABEL_REQUIRED ***REMOVED***
			if err := res.errs.handleErrorWithPos(node.fieldLabel().start(), "%s: label 'required' is not allowed in proto3", scope); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED*** else if fld.Extendee != nil && fld.Label != nil && fld.GetLabel() == dpb.FieldDescriptorProto_LABEL_OPTIONAL ***REMOVED***
			if err := res.errs.handleErrorWithPos(node.fieldLabel().start(), "%s: label 'optional' is not allowed on extensions in proto3", scope); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
		if index, err := findOption(res, scope, fld.Options.GetUninterpretedOption(), "default"); err != nil ***REMOVED***
			return err
		***REMOVED*** else if index >= 0 ***REMOVED***
			optNode := res.getOptionNode(fld.Options.GetUninterpretedOption()[index])
			if err := res.errs.handleErrorWithPos(optNode.getName().start(), "%s: default values are not allowed in proto3", scope); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if fld.Label == nil && fld.OneofIndex == nil ***REMOVED***
			if err := res.errs.handleErrorWithPos(node.fieldName().start(), "%s: field has no label; proto2 requires explicit 'optional' label", scope); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
		if fld.GetExtendee() != "" && fld.Label != nil && fld.GetLabel() == dpb.FieldDescriptorProto_LABEL_REQUIRED ***REMOVED***
			if err := res.errs.handleErrorWithPos(node.fieldLabel().start(), "%s: extension fields cannot be 'required'", scope); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// finally, set any missing label to optional
	if fld.Label == nil ***REMOVED***
		fld.Label = dpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum()
	***REMOVED***

	return nil
***REMOVED***

type tagRange struct ***REMOVED***
	start int32
	end   int32
	node  rangeDecl
***REMOVED***

type tagRanges []tagRange

func (r tagRanges) Len() int ***REMOVED***
	return len(r)
***REMOVED***

func (r tagRanges) Less(i, j int) bool ***REMOVED***
	return r[i].start < r[j].start ||
		(r[i].start == r[j].start && r[i].end < r[j].end)
***REMOVED***

func (r tagRanges) Swap(i, j int) ***REMOVED***
	r[i], r[j] = r[j], r[i]
***REMOVED***
