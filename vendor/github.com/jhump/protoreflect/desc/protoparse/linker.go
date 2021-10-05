package protoparse

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"github.com/golang/protobuf/proto"
	dpb "github.com/golang/protobuf/protoc-gen-go/descriptor"

	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/internal"
	"github.com/jhump/protoreflect/desc/protoparse/ast"
)

type linker struct ***REMOVED***
	files             map[string]*parseResult
	filenames         []string
	errs              *errorHandler
	descriptorPool    map[*dpb.FileDescriptorProto]map[string]proto.Message
	packageNamespaces map[*dpb.FileDescriptorProto]map[string]struct***REMOVED******REMOVED***
	extensions        map[string]map[int32]string
	usedImports       map[*dpb.FileDescriptorProto]map[string]struct***REMOVED******REMOVED***
***REMOVED***

func newLinker(files *parseResults, errs *errorHandler) *linker ***REMOVED***
	return &linker***REMOVED***files: files.resultsByFilename, filenames: files.filenames, errs: errs***REMOVED***
***REMOVED***

func (l *linker) linkFiles() (map[string]*desc.FileDescriptor, error) ***REMOVED***
	// First, we put all symbols into a single pool, which lets us ensure there
	// are no duplicate symbols and will also let us resolve and revise all type
	// references in next step.
	if err := l.createDescriptorPool(); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// After we've populated the pool, we can now try to resolve all type
	// references. All references must be checked for correct type, any fields
	// with enum types must be corrected (since we parse them as if they are
	// message references since we don't actually know message or enum until
	// link time), and references will be re-written to be fully-qualified
	// references (e.g. start with a dot ".").
	if err := l.resolveReferences(); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if err := l.errs.getError(); err != nil ***REMOVED***
		// we won't be able to create real descriptors if we've encountered
		// errors up to this point, so bail at this point
		return nil, err
	***REMOVED***

	// Now we've validated the descriptors, so we can link them into rich
	// descriptors. This is a little redundant since that step does similar
	// checking of symbols. But, without breaking encapsulation (e.g. exporting
	// a lot of fields from desc package that are currently unexported) or
	// merging this into the same package, we can't really prevent it.
	linked, err := l.createdLinkedDescriptors()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// Now that we have linked descriptors, we can interpret any uninterpreted
	// options that remain.
	for _, r := range l.files ***REMOVED***
		fd := linked[r.fd.GetName()]
		if err := interpretFileOptions(l, r, richFileDescriptorish***REMOVED***FileDescriptor: fd***REMOVED***); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		// we should now have any message_set_wire_format options parsed
		// and can do further validation on tag ranges
		if err := checkExtensionsInFile(fd, r); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	// When Parser calls linkFiles, it does not check errs again, and it expects that linkFiles
	// will return all errors it should process. If the ErrorReporter handles all errors itself
	// and always returns nil, we should get ErrInvalidSource here, and need to propagate this
	if err := l.errs.getError(); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return linked, nil
***REMOVED***

func (l *linker) createDescriptorPool() error ***REMOVED***
	l.descriptorPool = map[*dpb.FileDescriptorProto]map[string]proto.Message***REMOVED******REMOVED***
	l.packageNamespaces = map[*dpb.FileDescriptorProto]map[string]struct***REMOVED******REMOVED******REMOVED******REMOVED***
	for _, filename := range l.filenames ***REMOVED***
		r := l.files[filename]
		fd := r.fd
		pool := map[string]proto.Message***REMOVED******REMOVED***
		l.descriptorPool[fd] = pool
		prefix := fd.GetPackage()
		l.packageNamespaces[fd] = namespacesFromPackage(prefix)
		if prefix != "" ***REMOVED***
			prefix += "."
		***REMOVED***
		for _, md := range fd.MessageType ***REMOVED***
			if err := addMessageToPool(r, pool, l.errs, prefix, md); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
		for _, fld := range fd.Extension ***REMOVED***
			if err := addFieldToPool(r, pool, l.errs, prefix, fld); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
		for _, ed := range fd.EnumType ***REMOVED***
			if err := addEnumToPool(r, pool, l.errs, prefix, ed); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
		for _, sd := range fd.Service ***REMOVED***
			if err := addServiceToPool(r, pool, l.errs, prefix, sd); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***
	// try putting everything into a single pool, to ensure there are no duplicates
	// across files (e.g. same symbol, but declared in two different files)
	type entry struct ***REMOVED***
		file string
		msg  proto.Message
	***REMOVED***
	pool := map[string]entry***REMOVED******REMOVED***
	for _, filename := range l.filenames ***REMOVED***
		f := l.files[filename].fd
		p := l.descriptorPool[f]
		keys := make([]string, 0, len(p))
		for k := range p ***REMOVED***
			keys = append(keys, k)
		***REMOVED***
		sort.Strings(keys) // for deterministic error reporting
		for _, k := range keys ***REMOVED***
			v := p[k]
			if e, ok := pool[k]; ok ***REMOVED***
				desc1 := e.msg
				file1 := e.file
				desc2 := v
				file2 := f.GetName()
				if file2 < file1 ***REMOVED***
					file1, file2 = file2, file1
					desc1, desc2 = desc2, desc1
				***REMOVED***
				node := l.files[file2].nodes[desc2]
				if err := l.errs.handleErrorWithPos(node.Start(), "duplicate symbol %s: already defined as %s in %q", k, descriptorType(desc1), file1); err != nil ***REMOVED***
					return err
				***REMOVED***
			***REMOVED***
			pool[k] = entry***REMOVED***file: f.GetName(), msg: v***REMOVED***
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func namespacesFromPackage(pkg string) map[string]struct***REMOVED******REMOVED*** ***REMOVED***
	if pkg == "" ***REMOVED***
		return nil
	***REMOVED***
	offs := 0
	pkgs := map[string]struct***REMOVED******REMOVED******REMOVED******REMOVED***
	pkgs[pkg] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	for ***REMOVED***
		pos := strings.IndexByte(pkg[offs:], '.')
		if pos == -1 ***REMOVED***
			return pkgs
		***REMOVED***
		pkgs[pkg[:offs+pos]] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
		offs = offs + pos + 1
	***REMOVED***
***REMOVED***

func addMessageToPool(r *parseResult, pool map[string]proto.Message, errs *errorHandler, prefix string, md *dpb.DescriptorProto) error ***REMOVED***
	fqn := prefix + md.GetName()
	if err := addToPool(r, pool, errs, fqn, md); err != nil ***REMOVED***
		return err
	***REMOVED***
	prefix = fqn + "."
	for _, fld := range md.Field ***REMOVED***
		if err := addFieldToPool(r, pool, errs, prefix, fld); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	for _, fld := range md.Extension ***REMOVED***
		if err := addFieldToPool(r, pool, errs, prefix, fld); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	for _, nmd := range md.NestedType ***REMOVED***
		if err := addMessageToPool(r, pool, errs, prefix, nmd); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	for _, ed := range md.EnumType ***REMOVED***
		if err := addEnumToPool(r, pool, errs, prefix, ed); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func addFieldToPool(r *parseResult, pool map[string]proto.Message, errs *errorHandler, prefix string, fld *dpb.FieldDescriptorProto) error ***REMOVED***
	fqn := prefix + fld.GetName()
	return addToPool(r, pool, errs, fqn, fld)
***REMOVED***

func addEnumToPool(r *parseResult, pool map[string]proto.Message, errs *errorHandler, prefix string, ed *dpb.EnumDescriptorProto) error ***REMOVED***
	fqn := prefix + ed.GetName()
	if err := addToPool(r, pool, errs, fqn, ed); err != nil ***REMOVED***
		return err
	***REMOVED***
	for _, evd := range ed.Value ***REMOVED***
		// protobuf name-scoping rules for enum values follow C++ scoping rules:
		// the enum value name is a symbol in the *parent* scope (the one
		// enclosing the enum).
		vfqn := prefix + evd.GetName()
		if err := addToPool(r, pool, errs, vfqn, evd); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func addServiceToPool(r *parseResult, pool map[string]proto.Message, errs *errorHandler, prefix string, sd *dpb.ServiceDescriptorProto) error ***REMOVED***
	fqn := prefix + sd.GetName()
	if err := addToPool(r, pool, errs, fqn, sd); err != nil ***REMOVED***
		return err
	***REMOVED***
	for _, mtd := range sd.Method ***REMOVED***
		mfqn := fqn + "." + mtd.GetName()
		if err := addToPool(r, pool, errs, mfqn, mtd); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func addToPool(r *parseResult, pool map[string]proto.Message, errs *errorHandler, fqn string, dsc proto.Message) error ***REMOVED***
	if d, ok := pool[fqn]; ok ***REMOVED***
		node := r.nodes[dsc]
		_, additionIsEnumVal := dsc.(*dpb.EnumValueDescriptorProto)
		_, existingIsEnumVal := d.(*dpb.EnumValueDescriptorProto)
		// because of weird scoping for enum values, provide more context in error message
		// if this conflict is with an enum value
		var suffix string
		if additionIsEnumVal || existingIsEnumVal ***REMOVED***
			suffix = "; protobuf uses C++ scoping rules for enum values, so they exist in the scope enclosing the enum"
		***REMOVED***
		// TODO: also include the source location for the conflicting symbol
		if err := errs.handleErrorWithPos(node.Start(), "duplicate symbol %s: already defined as %s%s", fqn, descriptorType(d), suffix); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	pool[fqn] = dsc
	return nil
***REMOVED***

func descriptorType(m proto.Message) string ***REMOVED***
	switch m := m.(type) ***REMOVED***
	case *dpb.DescriptorProto:
		return "message"
	case *dpb.DescriptorProto_ExtensionRange:
		return "extension range"
	case *dpb.FieldDescriptorProto:
		if m.GetExtendee() == "" ***REMOVED***
			return "field"
		***REMOVED*** else ***REMOVED***
			return "extension"
		***REMOVED***
	case *dpb.EnumDescriptorProto:
		return "enum"
	case *dpb.EnumValueDescriptorProto:
		return "enum value"
	case *dpb.ServiceDescriptorProto:
		return "service"
	case *dpb.MethodDescriptorProto:
		return "method"
	case *dpb.FileDescriptorProto:
		return "file"
	default:
		// shouldn't be possible
		return fmt.Sprintf("%T", m)
	***REMOVED***
***REMOVED***

func (l *linker) resolveReferences() error ***REMOVED***
	l.extensions = map[string]map[int32]string***REMOVED******REMOVED***
	l.usedImports = map[*dpb.FileDescriptorProto]map[string]struct***REMOVED******REMOVED******REMOVED******REMOVED***
	for _, filename := range l.filenames ***REMOVED***
		r := l.files[filename]
		fd := r.fd
		prefix := fd.GetPackage()
		scopes := []scope***REMOVED***fileScope(fd, l)***REMOVED***
		if prefix != "" ***REMOVED***
			prefix += "."
		***REMOVED***
		if fd.Options != nil ***REMOVED***
			if err := l.resolveOptions(r, fd, "file", fd.GetName(), proto.MessageName(fd.Options), fd.Options.UninterpretedOption, scopes); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
		for _, md := range fd.MessageType ***REMOVED***
			if err := l.resolveMessageTypes(r, fd, prefix, md, scopes); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
		for _, fld := range fd.Extension ***REMOVED***
			if err := l.resolveFieldTypes(r, fd, prefix, fld, scopes); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
		for _, ed := range fd.EnumType ***REMOVED***
			if err := l.resolveEnumTypes(r, fd, prefix, ed, scopes); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
		for _, sd := range fd.Service ***REMOVED***
			if err := l.resolveServiceTypes(r, fd, prefix, sd, scopes); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (l *linker) resolveEnumTypes(r *parseResult, fd *dpb.FileDescriptorProto, prefix string, ed *dpb.EnumDescriptorProto, scopes []scope) error ***REMOVED***
	enumFqn := prefix + ed.GetName()
	if ed.Options != nil ***REMOVED***
		if err := l.resolveOptions(r, fd, "enum", enumFqn, proto.MessageName(ed.Options), ed.Options.UninterpretedOption, scopes); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	for _, evd := range ed.Value ***REMOVED***
		if evd.Options != nil ***REMOVED***
			evFqn := enumFqn + "." + evd.GetName()
			if err := l.resolveOptions(r, fd, "enum value", evFqn, proto.MessageName(evd.Options), evd.Options.UninterpretedOption, scopes); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (l *linker) resolveMessageTypes(r *parseResult, fd *dpb.FileDescriptorProto, prefix string, md *dpb.DescriptorProto, scopes []scope) error ***REMOVED***
	fqn := prefix + md.GetName()
	scope := messageScope(fqn, isProto3(fd), l, fd)
	scopes = append(scopes, scope)
	prefix = fqn + "."

	if md.Options != nil ***REMOVED***
		if err := l.resolveOptions(r, fd, "message", fqn, proto.MessageName(md.Options), md.Options.UninterpretedOption, scopes); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	for _, nmd := range md.NestedType ***REMOVED***
		if err := l.resolveMessageTypes(r, fd, prefix, nmd, scopes); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	for _, ned := range md.EnumType ***REMOVED***
		if err := l.resolveEnumTypes(r, fd, prefix, ned, scopes); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	for _, fld := range md.Field ***REMOVED***
		if err := l.resolveFieldTypes(r, fd, prefix, fld, scopes); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	for _, ood := range md.OneofDecl ***REMOVED***
		if ood.Options != nil ***REMOVED***
			ooName := fmt.Sprintf("%s.%s", fqn, ood.GetName())
			if err := l.resolveOptions(r, fd, "oneof", ooName, proto.MessageName(ood.Options), ood.Options.UninterpretedOption, scopes); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***
	for _, fld := range md.Extension ***REMOVED***
		if err := l.resolveFieldTypes(r, fd, prefix, fld, scopes); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	for _, er := range md.ExtensionRange ***REMOVED***
		if er.Options != nil ***REMOVED***
			erName := fmt.Sprintf("%s:%d-%d", fqn, er.GetStart(), er.GetEnd()-1)
			if err := l.resolveOptions(r, fd, "extension range", erName, proto.MessageName(er.Options), er.Options.UninterpretedOption, scopes); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (l *linker) resolveFieldTypes(r *parseResult, fd *dpb.FileDescriptorProto, prefix string, fld *dpb.FieldDescriptorProto, scopes []scope) error ***REMOVED***
	thisName := prefix + fld.GetName()
	scope := fmt.Sprintf("field %s", thisName)
	node := r.getFieldNode(fld)
	elemType := "field"
	if fld.GetExtendee() != "" ***REMOVED***
		elemType = "extension"
		fqn, dsc, _ := l.resolve(fd, fld.GetExtendee(), true, scopes)
		if dsc == nil ***REMOVED***
			return l.errs.handleErrorWithPos(node.FieldExtendee().Start(), "unknown extendee type %s", fld.GetExtendee())
		***REMOVED***
		if dsc == sentinelMissingSymbol ***REMOVED***
			return l.errs.handleErrorWithPos(node.FieldExtendee().Start(), "unknown extendee type %s; resolved to %s which is not defined; consider using a leading dot", fld.GetExtendee(), fqn)
		***REMOVED***
		extd, ok := dsc.(*dpb.DescriptorProto)
		if !ok ***REMOVED***
			otherType := descriptorType(dsc)
			return l.errs.handleErrorWithPos(node.FieldExtendee().Start(), "extendee is invalid: %s is a %s, not a message", fqn, otherType)
		***REMOVED***
		fld.Extendee = proto.String("." + fqn)
		// make sure the tag number is in range
		found := false
		tag := fld.GetNumber()
		for _, rng := range extd.ExtensionRange ***REMOVED***
			if tag >= rng.GetStart() && tag < rng.GetEnd() ***REMOVED***
				found = true
				break
			***REMOVED***
		***REMOVED***
		if !found ***REMOVED***
			if err := l.errs.handleErrorWithPos(node.FieldTag().Start(), "%s: tag %d is not in valid range for extended type %s", scope, tag, fqn); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			// make sure tag is not a duplicate
			usedExtTags := l.extensions[fqn]
			if usedExtTags == nil ***REMOVED***
				usedExtTags = map[int32]string***REMOVED******REMOVED***
				l.extensions[fqn] = usedExtTags
			***REMOVED***
			if other := usedExtTags[fld.GetNumber()]; other != "" ***REMOVED***
				if err := l.errs.handleErrorWithPos(node.FieldTag().Start(), "%s: duplicate extension: %s and %s are both using tag %d", scope, other, thisName, fld.GetNumber()); err != nil ***REMOVED***
					return err
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				usedExtTags[fld.GetNumber()] = thisName
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if fld.Options != nil ***REMOVED***
		if err := l.resolveOptions(r, fd, elemType, thisName, proto.MessageName(fld.Options), fld.Options.UninterpretedOption, scopes); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	if fld.GetTypeName() == "" ***REMOVED***
		// scalar type; no further resolution required
		return nil
	***REMOVED***

	fqn, dsc, proto3 := l.resolve(fd, fld.GetTypeName(), true, scopes)
	if dsc == nil ***REMOVED***
		return l.errs.handleErrorWithPos(node.FieldType().Start(), "%s: unknown type %s", scope, fld.GetTypeName())
	***REMOVED***
	if dsc == sentinelMissingSymbol ***REMOVED***
		return l.errs.handleErrorWithPos(node.FieldType().Start(), "%s: unknown type %s; resolved to %s which is not defined; consider using a leading dot", scope, fld.GetTypeName(), fqn)
	***REMOVED***
	switch dsc := dsc.(type) ***REMOVED***
	case *dpb.DescriptorProto:
		fld.TypeName = proto.String("." + fqn)
		// if type was tentatively unset, we now know it's actually a message
		if fld.Type == nil ***REMOVED***
			fld.Type = dpb.FieldDescriptorProto_TYPE_MESSAGE.Enum()
		***REMOVED***
	case *dpb.EnumDescriptorProto:
		if fld.GetExtendee() == "" && isProto3(fd) && !proto3 ***REMOVED***
			// fields in a proto3 message cannot refer to proto2 enums
			return l.errs.handleErrorWithPos(node.FieldType().Start(), "%s: cannot use proto2 enum %s in a proto3 message", scope, fld.GetTypeName())
		***REMOVED***
		fld.TypeName = proto.String("." + fqn)
		// the type was tentatively unset, but now we know it's actually an enum
		fld.Type = dpb.FieldDescriptorProto_TYPE_ENUM.Enum()
	default:
		otherType := descriptorType(dsc)
		return l.errs.handleErrorWithPos(node.FieldType().Start(), "%s: invalid type: %s is a %s, not a message or enum", scope, fqn, otherType)
	***REMOVED***
	return nil
***REMOVED***

func (l *linker) resolveServiceTypes(r *parseResult, fd *dpb.FileDescriptorProto, prefix string, sd *dpb.ServiceDescriptorProto, scopes []scope) error ***REMOVED***
	svcFqn := prefix + sd.GetName()
	if sd.Options != nil ***REMOVED***
		if err := l.resolveOptions(r, fd, "service", svcFqn, proto.MessageName(sd.Options), sd.Options.UninterpretedOption, scopes); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	// not a message, but same scoping rules for nested elements as if it were
	scope := messageScope(svcFqn, isProto3(fd), l, fd)
	scopes = append(scopes, scope)

	for _, mtd := range sd.Method ***REMOVED***
		if mtd.Options != nil ***REMOVED***
			if err := l.resolveOptions(r, fd, "method", svcFqn+"."+mtd.GetName(), proto.MessageName(mtd.Options), mtd.Options.UninterpretedOption, scopes); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
		scope := fmt.Sprintf("method %s.%s", svcFqn, mtd.GetName())
		node := r.getMethodNode(mtd)
		fqn, dsc, _ := l.resolve(fd, mtd.GetInputType(), false, scopes)
		if dsc == nil ***REMOVED***
			if err := l.errs.handleErrorWithPos(node.GetInputType().Start(), "%s: unknown request type %s", scope, mtd.GetInputType()); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED*** else if dsc == sentinelMissingSymbol ***REMOVED***
			if err := l.errs.handleErrorWithPos(node.GetInputType().Start(), "%s: unknown request type %s; resolved to %s which is not defined; consider using a leading dot", scope, mtd.GetInputType(), fqn); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED*** else if _, ok := dsc.(*dpb.DescriptorProto); !ok ***REMOVED***
			otherType := descriptorType(dsc)
			if err := l.errs.handleErrorWithPos(node.GetInputType().Start(), "%s: invalid request type: %s is a %s, not a message", scope, fqn, otherType); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			mtd.InputType = proto.String("." + fqn)
		***REMOVED***

		// TODO: make input and output type resolution more DRY
		fqn, dsc, _ = l.resolve(fd, mtd.GetOutputType(), false, scopes)
		if dsc == nil ***REMOVED***
			if err := l.errs.handleErrorWithPos(node.GetOutputType().Start(), "%s: unknown response type %s", scope, mtd.GetOutputType()); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED*** else if dsc == sentinelMissingSymbol ***REMOVED***
			if err := l.errs.handleErrorWithPos(node.GetOutputType().Start(), "%s: unknown response type %s; resolved to %s which is not defined; consider using a leading dot", scope, mtd.GetOutputType(), fqn); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED*** else if _, ok := dsc.(*dpb.DescriptorProto); !ok ***REMOVED***
			otherType := descriptorType(dsc)
			if err := l.errs.handleErrorWithPos(node.GetOutputType().Start(), "%s: invalid response type: %s is a %s, not a message", scope, fqn, otherType); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			mtd.OutputType = proto.String("." + fqn)
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (l *linker) resolveOptions(r *parseResult, fd *dpb.FileDescriptorProto, elemType, elemName, optType string, opts []*dpb.UninterpretedOption, scopes []scope) error ***REMOVED***
	var scope string
	if elemType != "file" ***REMOVED***
		scope = fmt.Sprintf("%s %s: ", elemType, elemName)
	***REMOVED***
opts:
	for _, opt := range opts ***REMOVED***
		for _, nm := range opt.Name ***REMOVED***
			if nm.GetIsExtension() ***REMOVED***
				node := r.getOptionNamePartNode(nm)
				fqn, dsc, _ := l.resolve(fd, nm.GetNamePart(), false, scopes)
				if dsc == nil ***REMOVED***
					if err := l.errs.handleErrorWithPos(node.Start(), "%sunknown extension %s", scope, nm.GetNamePart()); err != nil ***REMOVED***
						return err
					***REMOVED***
					continue opts
				***REMOVED***
				if dsc == sentinelMissingSymbol ***REMOVED***
					if err := l.errs.handleErrorWithPos(node.Start(), "%sunknown extension %s; resolved to %s which is not defined; consider using a leading dot", scope, nm.GetNamePart(), fqn); err != nil ***REMOVED***
						return err
					***REMOVED***
					continue opts
				***REMOVED***
				if ext, ok := dsc.(*dpb.FieldDescriptorProto); !ok ***REMOVED***
					otherType := descriptorType(dsc)
					if err := l.errs.handleErrorWithPos(node.Start(), "%sinvalid extension: %s is a %s, not an extension", scope, nm.GetNamePart(), otherType); err != nil ***REMOVED***
						return err
					***REMOVED***
					continue opts
				***REMOVED*** else if ext.GetExtendee() == "" ***REMOVED***
					if err := l.errs.handleErrorWithPos(node.Start(), "%sinvalid extension: %s is a field but not an extension", scope, nm.GetNamePart()); err != nil ***REMOVED***
						return err
					***REMOVED***
					continue opts
				***REMOVED***
				nm.NamePart = proto.String("." + fqn)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (l *linker) resolve(fd *dpb.FileDescriptorProto, name string, onlyTypes bool, scopes []scope) (fqn string, element proto.Message, proto3 bool) ***REMOVED***
	if strings.HasPrefix(name, ".") ***REMOVED***
		// already fully-qualified
		d, proto3 := l.findSymbol(fd, name[1:])
		if d != nil ***REMOVED***
			return name[1:], d, proto3
		***REMOVED***
		return "", nil, false
	***REMOVED***
	// unqualified, so we look in the enclosing (last) scope first and move
	// towards outermost (first) scope, trying to resolve the symbol
	pos := strings.IndexByte(name, '.')
	firstName := name
	if pos > 0 ***REMOVED***
		firstName = name[:pos]
	***REMOVED***
	var bestGuess proto.Message
	var bestGuessFqn string
	var bestGuessProto3 bool
	for i := len(scopes) - 1; i >= 0; i-- ***REMOVED***
		fqn, d, proto3 := scopes[i](firstName, name)
		if d != nil ***REMOVED***
			if !onlyTypes || isType(d) ***REMOVED***
				return fqn, d, proto3
			***REMOVED*** else if bestGuess == nil ***REMOVED***
				bestGuess = d
				bestGuessFqn = fqn
				bestGuessProto3 = proto3
			***REMOVED***
		***REMOVED***
	***REMOVED***
	// we return best guess, even though it was not an allowed kind of
	// descriptor, so caller can print a better error message (e.g.
	// indicating that the name was found but that it's the wrong type)
	return bestGuessFqn, bestGuess, bestGuessProto3
***REMOVED***

func isType(m proto.Message) bool ***REMOVED***
	switch m.(type) ***REMOVED***
	case *dpb.DescriptorProto, *dpb.EnumDescriptorProto:
		return true
	***REMOVED***
	return false
***REMOVED***

// scope represents a lexical scope in a proto file in which messages and enums
// can be declared.
type scope func(firstName, fullName string) (fqn string, element proto.Message, proto3 bool)

func fileScope(fd *dpb.FileDescriptorProto, l *linker) scope ***REMOVED***
	// we search symbols in this file, but also symbols in other files that have
	// the same package as this file or a "parent" package (in protobuf,
	// packages are a hierarchy like C++ namespaces)
	prefixes := internal.CreatePrefixList(fd.GetPackage())
	querySymbol := func(n string) (d proto.Message, isProto3 bool) ***REMOVED***
		return l.findSymbol(fd, n)
	***REMOVED***
	return func(firstName, fullName string) (string, proto.Message, bool) ***REMOVED***
		for _, prefix := range prefixes ***REMOVED***
			var n1, n string
			if prefix == "" ***REMOVED***
				// exhausted all prefixes, so it must be in this one
				n1, n = fullName, fullName
			***REMOVED*** else ***REMOVED***
				n = prefix + "." + fullName
				n1 = prefix + "." + firstName
			***REMOVED***
			d, proto3 := findSymbolRelative(n1, n, querySymbol)
			if d != nil ***REMOVED***
				return n, d, proto3
			***REMOVED***
		***REMOVED***
		return "", nil, false
	***REMOVED***
***REMOVED***

func messageScope(messageName string, proto3 bool, l *linker, fd *dpb.FileDescriptorProto) scope ***REMOVED***
	querySymbol := func(n string) (d proto.Message, isProto3 bool) ***REMOVED***
		return l.findSymbolInFile(n, fd), false
	***REMOVED***
	return func(firstName, fullName string) (string, proto.Message, bool) ***REMOVED***
		n1 := messageName + "." + firstName
		n := messageName + "." + fullName
		d, _ := findSymbolRelative(n1, n, querySymbol)
		if d != nil ***REMOVED***
			return n, d, proto3
		***REMOVED***
		return "", nil, false
	***REMOVED***
***REMOVED***

func findSymbolRelative(firstName, fullName string, query func(name string) (d proto.Message, isProto3 bool)) (d proto.Message, isProto3 bool) ***REMOVED***
	d, proto3 := query(firstName)
	if d == nil ***REMOVED***
		return nil, false
	***REMOVED***
	if firstName == fullName ***REMOVED***
		return d, proto3
	***REMOVED***
	if !isAggregateDescriptor(d) ***REMOVED***
		// can't possibly find the rest of full name if
		// the first name indicated a leaf descriptor
		return nil, false
	***REMOVED***
	d, proto3 = query(fullName)
	if d == nil ***REMOVED***
		return sentinelMissingSymbol, false
	***REMOVED***
	return d, proto3
***REMOVED***

func (l *linker) findSymbolInFile(name string, fd *dpb.FileDescriptorProto) proto.Message ***REMOVED***
	d, ok := l.descriptorPool[fd][name]
	if ok ***REMOVED***
		return d
	***REMOVED***
	_, ok = l.packageNamespaces[fd][name]
	if ok ***REMOVED***
		// this sentinel means the name is a valid namespace but
		// does not refer to a descriptor
		return sentinelMissingSymbol
	***REMOVED***
	return nil
***REMOVED***

func (l *linker) markUsed(entryPoint, used *dpb.FileDescriptorProto) ***REMOVED***
	importsForFile := l.usedImports[entryPoint]
	if importsForFile == nil ***REMOVED***
		importsForFile = map[string]struct***REMOVED******REMOVED******REMOVED******REMOVED***
		l.usedImports[entryPoint] = importsForFile
	***REMOVED***
	importsForFile[used.GetName()] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
***REMOVED***

func isAggregateDescriptor(m proto.Message) bool ***REMOVED***
	if m == sentinelMissingSymbol ***REMOVED***
		// this indicates the name matched a package, not a
		// descriptor, but a package is an aggregate so
		// we return true
		return true
	***REMOVED***
	switch m.(type) ***REMOVED***
	case *dpb.DescriptorProto, *dpb.EnumDescriptorProto, *dpb.ServiceDescriptorProto:
		return true
	default:
		return false
	***REMOVED***
***REMOVED***

// This value is a bogus/nil value, but results in a non-nil
// proto.Message interface value. So we use it as a sentinel
// to indicate "stop searching for symbol... because it
// definitively does not exist".
var sentinelMissingSymbol = (*dpb.DescriptorProto)(nil)

func (l *linker) findSymbol(fd *dpb.FileDescriptorProto, name string) (element proto.Message, proto3 bool) ***REMOVED***
	return l.findSymbolRecursive(fd, fd, name, false, map[*dpb.FileDescriptorProto]struct***REMOVED******REMOVED******REMOVED******REMOVED***)
***REMOVED***

func (l *linker) findSymbolRecursive(entryPoint, fd *dpb.FileDescriptorProto, name string, public bool, checked map[*dpb.FileDescriptorProto]struct***REMOVED******REMOVED***) (element proto.Message, proto3 bool) ***REMOVED***
	if _, ok := checked[fd]; ok ***REMOVED***
		// already checked this one
		return nil, false
	***REMOVED***
	checked[fd] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	d := l.findSymbolInFile(name, fd)
	if d != nil ***REMOVED***
		return d, isProto3(fd)
	***REMOVED***

	// When public = false, we are searching only directly imported symbols. But we
	// also need to search transitive public imports due to semantics of public imports.
	if public ***REMOVED***
		for _, depIndex := range fd.PublicDependency ***REMOVED***
			dep := fd.Dependency[depIndex]
			depres := l.files[dep]
			if depres == nil ***REMOVED***
				// we'll catch this error later
				continue
			***REMOVED***
			if d, proto3 := l.findSymbolRecursive(entryPoint, depres.fd, name, true, checked); d != nil ***REMOVED***
				l.markUsed(entryPoint, depres.fd)
				return d, proto3
			***REMOVED***
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		for _, dep := range fd.Dependency ***REMOVED***
			depres := l.files[dep]
			if depres == nil ***REMOVED***
				// we'll catch this error later
				continue
			***REMOVED***
			if d, proto3 := l.findSymbolRecursive(entryPoint, depres.fd, name, true, checked); d != nil ***REMOVED***
				l.markUsed(entryPoint, depres.fd)
				return d, proto3
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return nil, false
***REMOVED***

func isProto3(fd *dpb.FileDescriptorProto) bool ***REMOVED***
	return fd.GetSyntax() == "proto3"
***REMOVED***

func (l *linker) createdLinkedDescriptors() (map[string]*desc.FileDescriptor, error) ***REMOVED***
	names := make([]string, 0, len(l.files))
	for name := range l.files ***REMOVED***
		names = append(names, name)
	***REMOVED***
	sort.Strings(names)
	linked := map[string]*desc.FileDescriptor***REMOVED******REMOVED***
	for _, name := range names ***REMOVED***
		if _, err := l.linkFile(name, nil, nil, linked); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	return linked, nil
***REMOVED***

func (l *linker) linkFile(name string, rootImportLoc *SourcePos, seen []string, linked map[string]*desc.FileDescriptor) (*desc.FileDescriptor, error) ***REMOVED***
	// check for import cycle
	for _, s := range seen ***REMOVED***
		if name == s ***REMOVED***
			var msg bytes.Buffer
			first := true
			for _, s := range seen ***REMOVED***
				if first ***REMOVED***
					first = false
				***REMOVED*** else ***REMOVED***
					msg.WriteString(" -> ")
				***REMOVED***
				_, _ = fmt.Fprintf(&msg, "%q", s)
			***REMOVED***
			_, _ = fmt.Fprintf(&msg, " -> %q", name)
			return nil, ErrorWithSourcePos***REMOVED***
				Underlying: fmt.Errorf("cycle found in imports: %s", msg.String()),
				Pos:        rootImportLoc,
			***REMOVED***
		***REMOVED***
	***REMOVED***
	seen = append(seen, name)

	if lfd, ok := linked[name]; ok ***REMOVED***
		// already linked
		return lfd, nil
	***REMOVED***
	r := l.files[name]
	if r == nil ***REMOVED***
		importer := seen[len(seen)-2] // len-1 is *this* file, before that is the one that imported it
		return nil, fmt.Errorf("no descriptor found for %q, imported by %q", name, importer)
	***REMOVED***
	var deps []*desc.FileDescriptor
	if rootImportLoc == nil ***REMOVED***
		// try to find a source location for this "root" import
		decl := r.getFileNode(r.fd)
		fnode, ok := decl.(*ast.FileNode)
		if ok ***REMOVED***
			for _, decl := range fnode.Decls ***REMOVED***
				if dep, ok := decl.(*ast.ImportNode); ok ***REMOVED***
					ldep, err := l.linkFile(dep.Name.AsString(), dep.Name.Start(), seen, linked)
					if err != nil ***REMOVED***
						return nil, err
					***REMOVED***
					deps = append(deps, ldep)
				***REMOVED***
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			// no AST? just use the descriptor
			for _, dep := range r.fd.Dependency ***REMOVED***
				ldep, err := l.linkFile(dep, decl.Start(), seen, linked)
				if err != nil ***REMOVED***
					return nil, err
				***REMOVED***
				deps = append(deps, ldep)
			***REMOVED***
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		// we can just use the descriptor since we don't need source location
		// (we'll just attribute any import cycles found to the "root" import)
		for _, dep := range r.fd.Dependency ***REMOVED***
			ldep, err := l.linkFile(dep, rootImportLoc, seen, linked)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			deps = append(deps, ldep)
		***REMOVED***
	***REMOVED***
	lfd, err := desc.CreateFileDescriptor(r.fd, deps...)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("error linking %q: %s", name, err)
	***REMOVED***
	linked[name] = lfd
	return lfd, nil
***REMOVED***

func (l *linker) checkForUnusedImports(filename string) ***REMOVED***
	r := l.files[filename]
	usedImports := l.usedImports[r.fd]
	node := r.nodes[r.fd]
	fileNode, _ := node.(*ast.FileNode)
	for i, dep := range r.fd.Dependency ***REMOVED***
		if _, ok := usedImports[dep]; !ok ***REMOVED***
			isPublic := false
			// it's fine if it's a public import
			for _, j := range r.fd.PublicDependency ***REMOVED***
				if i == int(j) ***REMOVED***
					isPublic = true
					break
				***REMOVED***
			***REMOVED***
			if isPublic ***REMOVED***
				break
			***REMOVED***
			var pos *SourcePos
			if fileNode != nil ***REMOVED***
				for _, decl := range fileNode.Decls ***REMOVED***
					imp, ok := decl.(*ast.ImportNode)
					if !ok ***REMOVED***
						continue
					***REMOVED***
					if imp.Name.AsString() == dep ***REMOVED***
						pos = imp.Start()
					***REMOVED***
				***REMOVED***
			***REMOVED***
			if pos == nil ***REMOVED***
				pos = ast.UnknownPos(r.fd.GetName())
			***REMOVED***
			r.errs.warn(pos, errUnusedImport(dep))
		***REMOVED***
	***REMOVED***
***REMOVED***
