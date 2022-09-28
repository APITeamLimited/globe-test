package desc

import (
	"errors"
	"fmt"
	"strings"

	dpb "github.com/golang/protobuf/protoc-gen-go/descriptor"

	"github.com/jhump/protoreflect/desc/internal"
	intn "github.com/jhump/protoreflect/internal"
)

// CreateFileDescriptor instantiates a new file descriptor for the given descriptor proto.
// The file's direct dependencies must be provided. If the given dependencies do not include
// all of the file's dependencies or if the contents of the descriptors are internally
// inconsistent (e.g. contain unresolvable symbols) then an error is returned.
func CreateFileDescriptor(fd *dpb.FileDescriptorProto, deps ...*FileDescriptor) (*FileDescriptor, error) ***REMOVED***
	return createFileDescriptor(fd, deps, nil)
***REMOVED***

func createFileDescriptor(fd *dpb.FileDescriptorProto, deps []*FileDescriptor, r *ImportResolver) (*FileDescriptor, error) ***REMOVED***
	ret := &FileDescriptor***REMOVED***
		proto:      fd,
		symbols:    map[string]Descriptor***REMOVED******REMOVED***,
		fieldIndex: map[string]map[int32]*FieldDescriptor***REMOVED******REMOVED***,
	***REMOVED***
	pkg := fd.GetPackage()

	// populate references to file descriptor dependencies
	files := map[string]*FileDescriptor***REMOVED******REMOVED***
	for _, f := range deps ***REMOVED***
		files[f.proto.GetName()] = f
	***REMOVED***
	ret.deps = make([]*FileDescriptor, len(fd.GetDependency()))
	for i, d := range fd.GetDependency() ***REMOVED***
		resolved := r.ResolveImport(fd.GetName(), d)
		ret.deps[i] = files[resolved]
		if ret.deps[i] == nil ***REMOVED***
			if resolved != d ***REMOVED***
				ret.deps[i] = files[d]
			***REMOVED***
			if ret.deps[i] == nil ***REMOVED***
				return nil, intn.ErrNoSuchFile(d)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	ret.publicDeps = make([]*FileDescriptor, len(fd.GetPublicDependency()))
	for i, pd := range fd.GetPublicDependency() ***REMOVED***
		ret.publicDeps[i] = ret.deps[pd]
	***REMOVED***
	ret.weakDeps = make([]*FileDescriptor, len(fd.GetWeakDependency()))
	for i, wd := range fd.GetWeakDependency() ***REMOVED***
		ret.weakDeps[i] = ret.deps[wd]
	***REMOVED***
	ret.isProto3 = fd.GetSyntax() == "proto3"

	// populate all tables of child descriptors
	for _, m := range fd.GetMessageType() ***REMOVED***
		md, n := createMessageDescriptor(ret, ret, pkg, m, ret.symbols)
		ret.symbols[n] = md
		ret.messages = append(ret.messages, md)
	***REMOVED***
	for _, e := range fd.GetEnumType() ***REMOVED***
		ed, n := createEnumDescriptor(ret, ret, pkg, e, ret.symbols)
		ret.symbols[n] = ed
		ret.enums = append(ret.enums, ed)
	***REMOVED***
	for _, ex := range fd.GetExtension() ***REMOVED***
		exd, n := createFieldDescriptor(ret, ret, pkg, ex)
		ret.symbols[n] = exd
		ret.extensions = append(ret.extensions, exd)
	***REMOVED***
	for _, s := range fd.GetService() ***REMOVED***
		sd, n := createServiceDescriptor(ret, pkg, s, ret.symbols)
		ret.symbols[n] = sd
		ret.services = append(ret.services, sd)
	***REMOVED***

	ret.sourceInfo = internal.CreateSourceInfoMap(fd)
	ret.sourceInfoRecomputeFunc = ret.recomputeSourceInfo

	// now we can resolve all type references and source code info
	scopes := []scope***REMOVED***fileScope(ret)***REMOVED***
	path := make([]int32, 1, 8)
	path[0] = internal.File_messagesTag
	for i, md := range ret.messages ***REMOVED***
		if err := md.resolve(append(path, int32(i)), scopes); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	path[0] = internal.File_enumsTag
	for i, ed := range ret.enums ***REMOVED***
		ed.resolve(append(path, int32(i)))
	***REMOVED***
	path[0] = internal.File_extensionsTag
	for i, exd := range ret.extensions ***REMOVED***
		if err := exd.resolve(append(path, int32(i)), scopes); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	path[0] = internal.File_servicesTag
	for i, sd := range ret.services ***REMOVED***
		if err := sd.resolve(append(path, int32(i)), scopes); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	return ret, nil
***REMOVED***

// CreateFileDescriptors constructs a set of descriptors, one for each of the
// given descriptor protos. The given set of descriptor protos must include all
// transitive dependencies for every file.
func CreateFileDescriptors(fds []*dpb.FileDescriptorProto) (map[string]*FileDescriptor, error) ***REMOVED***
	return createFileDescriptors(fds, nil)
***REMOVED***

func createFileDescriptors(fds []*dpb.FileDescriptorProto, r *ImportResolver) (map[string]*FileDescriptor, error) ***REMOVED***
	if len(fds) == 0 ***REMOVED***
		return nil, nil
	***REMOVED***
	files := map[string]*dpb.FileDescriptorProto***REMOVED******REMOVED***
	resolved := map[string]*FileDescriptor***REMOVED******REMOVED***
	var name string
	for _, fd := range fds ***REMOVED***
		name = fd.GetName()
		files[name] = fd
	***REMOVED***
	for _, fd := range fds ***REMOVED***
		_, err := createFromSet(fd.GetName(), r, nil, files, resolved)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	return resolved, nil
***REMOVED***

// ToFileDescriptorSet creates a FileDescriptorSet proto that contains all of the given
// file descriptors and their transitive dependencies. The files are topologically sorted
// so that a file will always appear after its dependencies.
func ToFileDescriptorSet(fds ...*FileDescriptor) *dpb.FileDescriptorSet ***REMOVED***
	var fdps []*dpb.FileDescriptorProto
	addAllFiles(fds, &fdps, map[string]struct***REMOVED******REMOVED******REMOVED******REMOVED***)
	return &dpb.FileDescriptorSet***REMOVED***File: fdps***REMOVED***
***REMOVED***

func addAllFiles(src []*FileDescriptor, results *[]*dpb.FileDescriptorProto, seen map[string]struct***REMOVED******REMOVED***) ***REMOVED***
	for _, fd := range src ***REMOVED***
		if _, ok := seen[fd.GetName()]; ok ***REMOVED***
			continue
		***REMOVED***
		seen[fd.GetName()] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
		addAllFiles(fd.GetDependencies(), results, seen)
		*results = append(*results, fd.AsFileDescriptorProto())
	***REMOVED***
***REMOVED***

// CreateFileDescriptorFromSet creates a descriptor from the given file descriptor set. The
// set's *last* file will be the returned descriptor. The set's remaining files must comprise
// the full set of transitive dependencies of that last file. This is the same format and
// order used by protoc when emitting a FileDescriptorSet file with an invocation like so:
//    protoc --descriptor_set_out=./test.protoset --include_imports -I. test.proto
func CreateFileDescriptorFromSet(fds *dpb.FileDescriptorSet) (*FileDescriptor, error) ***REMOVED***
	return createFileDescriptorFromSet(fds, nil)
***REMOVED***

func createFileDescriptorFromSet(fds *dpb.FileDescriptorSet, r *ImportResolver) (*FileDescriptor, error) ***REMOVED***
	result, err := createFileDescriptorsFromSet(fds, r)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	files := fds.GetFile()
	lastFilename := files[len(files)-1].GetName()
	return result[lastFilename], nil
***REMOVED***

// CreateFileDescriptorsFromSet creates file descriptors from the given file descriptor set.
// The returned map includes all files in the set, keyed b name. The set must include the
// full set of transitive dependencies for all files therein or else a link error will occur
// and be returned instead of the slice of descriptors. This is the same format used by
// protoc when a FileDescriptorSet file with an invocation like so:
//    protoc --descriptor_set_out=./test.protoset --include_imports -I. test.proto
func CreateFileDescriptorsFromSet(fds *dpb.FileDescriptorSet) (map[string]*FileDescriptor, error) ***REMOVED***
	return createFileDescriptorsFromSet(fds, nil)
***REMOVED***

func createFileDescriptorsFromSet(fds *dpb.FileDescriptorSet, r *ImportResolver) (map[string]*FileDescriptor, error) ***REMOVED***
	files := fds.GetFile()
	if len(files) == 0 ***REMOVED***
		return nil, errors.New("file descriptor set is empty")
	***REMOVED***
	return createFileDescriptors(files, r)
***REMOVED***

// createFromSet creates a descriptor for the given filename. It recursively
// creates descriptors for the given file's dependencies.
func createFromSet(filename string, r *ImportResolver, seen []string, files map[string]*dpb.FileDescriptorProto, resolved map[string]*FileDescriptor) (*FileDescriptor, error) ***REMOVED***
	for _, s := range seen ***REMOVED***
		if filename == s ***REMOVED***
			return nil, fmt.Errorf("cycle in imports: %s", strings.Join(append(seen, filename), " -> "))
		***REMOVED***
	***REMOVED***
	seen = append(seen, filename)

	if d, ok := resolved[filename]; ok ***REMOVED***
		return d, nil
	***REMOVED***
	fdp := files[filename]
	if fdp == nil ***REMOVED***
		return nil, intn.ErrNoSuchFile(filename)
	***REMOVED***
	deps := make([]*FileDescriptor, len(fdp.GetDependency()))
	for i, depName := range fdp.GetDependency() ***REMOVED***
		resolvedDep := r.ResolveImport(filename, depName)
		dep, err := createFromSet(resolvedDep, r, seen, files, resolved)
		if _, ok := err.(intn.ErrNoSuchFile); ok && resolvedDep != depName ***REMOVED***
			dep, err = createFromSet(depName, r, seen, files, resolved)
		***REMOVED***
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		deps[i] = dep
	***REMOVED***
	d, err := createFileDescriptor(fdp, deps, r)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	resolved[filename] = d
	return d, nil
***REMOVED***
