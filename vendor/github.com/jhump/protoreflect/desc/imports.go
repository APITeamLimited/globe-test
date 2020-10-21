package desc

import (
	"fmt"
	"path/filepath"
	"reflect"
	"strings"
	"sync"

	"github.com/golang/protobuf/proto"
	dpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
)

var (
	globalImportPathConf map[string]string
	globalImportPathMu   sync.RWMutex
)

// RegisterImportPath registers an alternate import path for a given registered
// proto file path. For more details on why alternate import paths may need to
// be configured, see ImportResolver.
//
// This method panics if provided invalid input. An empty importPath is invalid.
// An un-registered registerPath is also invalid. For example, if an attempt is
// made to register the import path "foo/bar.proto" as "bar.proto", but there is
// no "bar.proto" registered in the Go protobuf runtime, this method will panic.
// This method also panics if an attempt is made to register the same import
// path more than once.
//
// This function works globally, applying to all descriptors loaded by this
// package. If you instead want more granular support for handling alternate
// import paths -- such as for a single invocation of a function in this
// package or when the alternate path is only used from one file (so you don't
// want the alternate path used when loading every other file), use an
// ImportResolver instead.
func RegisterImportPath(registerPath, importPath string) ***REMOVED***
	if len(importPath) == 0 ***REMOVED***
		panic("import path cannot be empty")
	***REMOVED***
	desc := proto.FileDescriptor(registerPath)
	if len(desc) == 0 ***REMOVED***
		panic(fmt.Sprintf("path %q is not a registered proto file", registerPath))
	***REMOVED***
	globalImportPathMu.Lock()
	defer globalImportPathMu.Unlock()
	if reg := globalImportPathConf[importPath]; reg != "" ***REMOVED***
		panic(fmt.Sprintf("import path %q already registered for %s", importPath, reg))
	***REMOVED***
	if globalImportPathConf == nil ***REMOVED***
		globalImportPathConf = map[string]string***REMOVED******REMOVED***
	***REMOVED***
	globalImportPathConf[importPath] = registerPath
***REMOVED***

// ResolveImport resolves the given import path. If it has been registered as an
// alternate via RegisterImportPath, the registered path is returned. Otherwise,
// the given import path is returned unchanged.
func ResolveImport(importPath string) string ***REMOVED***
	importPath = clean(importPath)
	globalImportPathMu.RLock()
	defer globalImportPathMu.RUnlock()
	reg := globalImportPathConf[importPath]
	if reg == "" ***REMOVED***
		return importPath
	***REMOVED***
	return reg
***REMOVED***

// ImportResolver lets you work-around linking issues that are caused by
// mismatches between how a particular proto source file is registered in the Go
// protobuf runtime and how that same file is imported by other files. The file
// is registered using the same relative path given to protoc when the file is
// compiled (i.e. when Go code is generated). So if any file tries to import
// that source file, but using a different relative path, then a link error will
// occur when this package tries to load a descriptor for the importing file.
//
// For example, let's say we have two proto source files: "foo/bar.proto" and
// "fubar/baz.proto". The latter imports the former using a line like so:
//    import "foo/bar.proto";
// However, when protoc is invoked, the command-line args looks like so:
//    protoc -Ifoo/ --go_out=foo/ bar.proto
//    protoc -I./ -Ifubar/ --go_out=fubar/ baz.proto
// Because the path given to protoc is just "bar.proto" and "baz.proto", this is
// how they are registered in the Go protobuf runtime. So, when loading the
// descriptor for "fubar/baz.proto", we'll see an import path of "foo/bar.proto"
// but will find no file registered with that path:
//    fd, err := desc.LoadFileDescriptor("baz.proto")
//    // err will be non-nil, complaining that there is no such file
//    // found named "foo/bar.proto"
//
// This can be remedied by registering alternate import paths using an
// ImportResolver. Continuing with the example above, the code below would fix
// any link issue:
//    var r desc.ImportResolver
//    r.RegisterImportPath("bar.proto", "foo/bar.proto")
//    fd, err := r.LoadFileDescriptor("baz.proto")
//    // err will be nil; descriptor successfully loaded!
//
// If there are files that are *always* imported using a different relative
// path then how they are registered, consider using the global
// RegisterImportPath function, so you don't have to use an ImportResolver for
// every file that imports it.
type ImportResolver struct ***REMOVED***
	children    map[string]*ImportResolver
	importPaths map[string]string

	// By default, an ImportResolver will fallback to consulting any paths
	// registered via the top-level RegisterImportPath function. Setting this
	// field to true will cause the ImportResolver to skip that fallback and
	// only examine its own locally registered paths.
	SkipFallbackRules bool
***REMOVED***

// ResolveImport resolves the given import path in the context of the given
// source file. If a matching alternate has been registered with this resolver
// via a call to RegisterImportPath or RegisterImportPathFrom, then the
// registered path is returned. Otherwise, the given import path is returned
// unchanged.
func (r *ImportResolver) ResolveImport(source, importPath string) string ***REMOVED***
	if r != nil ***REMOVED***
		res := r.resolveImport(clean(source), clean(importPath))
		if res != "" ***REMOVED***
			return res
		***REMOVED***
		if r.SkipFallbackRules ***REMOVED***
			return importPath
		***REMOVED***
	***REMOVED***
	return ResolveImport(importPath)
***REMOVED***

func (r *ImportResolver) resolveImport(source, importPath string) string ***REMOVED***
	if source == "" ***REMOVED***
		return r.importPaths[importPath]
	***REMOVED***
	var car, cdr string
	idx := strings.IndexRune(source, filepath.Separator)
	if idx < 0 ***REMOVED***
		car, cdr = source, ""
	***REMOVED*** else ***REMOVED***
		car, cdr = source[:idx], source[idx+1:]
	***REMOVED***
	ch := r.children[car]
	if ch != nil ***REMOVED***
		if reg := ch.resolveImport(cdr, importPath); reg != "" ***REMOVED***
			return reg
		***REMOVED***
	***REMOVED***
	return r.importPaths[importPath]
***REMOVED***

// RegisterImportPath registers an alternate import path for a given registered
// proto file path with this resolver. Any appearance of the given import path
// when linking files will instead try to link the given registered path. If the
// registered path cannot be located, then linking will fallback to the actual
// imported path.
//
// This method will panic if given an empty path or if the same import path is
// registered more than once.
//
// To constrain the contexts where the given import path is to be re-written,
// use RegisterImportPathFrom instead.
func (r *ImportResolver) RegisterImportPath(registerPath, importPath string) ***REMOVED***
	r.RegisterImportPathFrom(registerPath, importPath, "")
***REMOVED***

// RegisterImportPathFrom registers an alternate import path for a given
// registered proto file path with this resolver, but only for imports in the
// specified source context.
//
// The source context can be the name of a folder or a proto source file. Any
// appearance of the given import path in that context will instead try to link
// the given registered path. To be in context, the file that is being linked
// (i.e. the one whose import statement is being resolved) must be the same
// relative path of the source context or be a sub-path (i.e. a descendant of
// the source folder).
//
// If the registered path cannot be located, then linking will fallback to the
// actual imported path.
//
// This method will panic if given an empty path. The source context, on the
// other hand, is allowed to be blank. A blank source matches all files. This
// method also panics if the same import path is registered in the same source
// context more than once.
func (r *ImportResolver) RegisterImportPathFrom(registerPath, importPath, source string) ***REMOVED***
	importPath = clean(importPath)
	if len(importPath) == 0 ***REMOVED***
		panic("import path cannot be empty")
	***REMOVED***
	registerPath = clean(registerPath)
	if len(registerPath) == 0 ***REMOVED***
		panic("registered path cannot be empty")
	***REMOVED***
	r.registerImportPathFrom(registerPath, importPath, clean(source))
***REMOVED***

func (r *ImportResolver) registerImportPathFrom(registerPath, importPath, source string) ***REMOVED***
	if source == "" ***REMOVED***
		if r.importPaths == nil ***REMOVED***
			r.importPaths = map[string]string***REMOVED******REMOVED***
		***REMOVED*** else if reg := r.importPaths[importPath]; reg != "" ***REMOVED***
			panic(fmt.Sprintf("already registered import path %q as %q", importPath, registerPath))
		***REMOVED***
		r.importPaths[importPath] = registerPath
		return
	***REMOVED***
	var car, cdr string
	idx := strings.IndexRune(source, filepath.Separator)
	if idx < 0 ***REMOVED***
		car, cdr = source, ""
	***REMOVED*** else ***REMOVED***
		car, cdr = source[:idx], source[idx+1:]
	***REMOVED***
	ch := r.children[car]
	if ch == nil ***REMOVED***
		if r.children == nil ***REMOVED***
			r.children = map[string]*ImportResolver***REMOVED******REMOVED***
		***REMOVED***
		ch = &ImportResolver***REMOVED******REMOVED***
		r.children[car] = ch
	***REMOVED***
	ch.registerImportPathFrom(registerPath, importPath, cdr)
***REMOVED***

// LoadFileDescriptor is the same as the package function of the same name, but
// any alternate paths configured in this resolver are used when linking the
// given descriptor proto.
func (r *ImportResolver) LoadFileDescriptor(filePath string) (*FileDescriptor, error) ***REMOVED***
	return loadFileDescriptor(filePath, r)
***REMOVED***

// LoadMessageDescriptor is the same as the package function of the same name,
// but any alternate paths configured in this resolver are used when linking
// files for the returned descriptor.
func (r *ImportResolver) LoadMessageDescriptor(msgName string) (*MessageDescriptor, error) ***REMOVED***
	return loadMessageDescriptor(msgName, r)
***REMOVED***

// LoadMessageDescriptorForMessage is the same as the package function of the
// same name, but any alternate paths configured in this resolver are used when
// linking files for the returned descriptor.
func (r *ImportResolver) LoadMessageDescriptorForMessage(msg proto.Message) (*MessageDescriptor, error) ***REMOVED***
	return loadMessageDescriptorForMessage(msg, r)
***REMOVED***

// LoadMessageDescriptorForType is the same as the package function of the same
// name, but any alternate paths configured in this resolver are used when
// linking files for the returned descriptor.
func (r *ImportResolver) LoadMessageDescriptorForType(msgType reflect.Type) (*MessageDescriptor, error) ***REMOVED***
	return loadMessageDescriptorForType(msgType, r)
***REMOVED***

// LoadEnumDescriptorForEnum is the same as the package function of the same
// name, but any alternate paths configured in this resolver are used when
// linking files for the returned descriptor.
func (r *ImportResolver) LoadEnumDescriptorForEnum(enum protoEnum) (*EnumDescriptor, error) ***REMOVED***
	return loadEnumDescriptorForEnum(enum, r)
***REMOVED***

// LoadEnumDescriptorForType is the same as the package function of the same
// name, but any alternate paths configured in this resolver are used when
// linking files for the returned descriptor.
func (r *ImportResolver) LoadEnumDescriptorForType(enumType reflect.Type) (*EnumDescriptor, error) ***REMOVED***
	return loadEnumDescriptorForType(enumType, r)
***REMOVED***

// LoadFieldDescriptorForExtension is the same as the package function of the
// same name, but any alternate paths configured in this resolver are used when
// linking files for the returned descriptor.
func (r *ImportResolver) LoadFieldDescriptorForExtension(ext *proto.ExtensionDesc) (*FieldDescriptor, error) ***REMOVED***
	return loadFieldDescriptorForExtension(ext, r)
***REMOVED***

// CreateFileDescriptor is the same as the package function of the same name,
// but any alternate paths configured in this resolver are used when linking the
// given descriptor proto.
func (r *ImportResolver) CreateFileDescriptor(fdp *dpb.FileDescriptorProto, deps ...*FileDescriptor) (*FileDescriptor, error) ***REMOVED***
	return createFileDescriptor(fdp, deps, r)
***REMOVED***

// CreateFileDescriptors is the same as the package function of the same name,
// but any alternate paths configured in this resolver are used when linking the
// given descriptor protos.
func (r *ImportResolver) CreateFileDescriptors(fds []*dpb.FileDescriptorProto) (map[string]*FileDescriptor, error) ***REMOVED***
	return createFileDescriptors(fds, r)
***REMOVED***

// CreateFileDescriptorFromSet is the same as the package function of the same
// name, but any alternate paths configured in this resolver are used when
// linking the descriptor protos in the given set.
func (r *ImportResolver) CreateFileDescriptorFromSet(fds *dpb.FileDescriptorSet) (*FileDescriptor, error) ***REMOVED***
	return createFileDescriptorFromSet(fds, r)
***REMOVED***

// CreateFileDescriptorsFromSet is the same as the package function of the same
// name, but any alternate paths configured in this resolver are used when
// linking the descriptor protos in the given set.
func (r *ImportResolver) CreateFileDescriptorsFromSet(fds *dpb.FileDescriptorSet) (map[string]*FileDescriptor, error) ***REMOVED***
	return createFileDescriptorsFromSet(fds, r)
***REMOVED***

const dotPrefix = "." + string(filepath.Separator)

func clean(path string) string ***REMOVED***
	if path == "" ***REMOVED***
		return ""
	***REMOVED***
	path = filepath.Clean(path)
	if path == "." ***REMOVED***
		return ""
	***REMOVED***
	return strings.TrimPrefix(path, dotPrefix)
***REMOVED***
